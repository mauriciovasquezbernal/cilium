// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package subnet

import (
	"context"
	"fmt"
	"log/slog"
	"net/netip"

	"github.com/cilium/hive/cell"
	"github.com/cilium/hive/job"
	"github.com/cilium/statedb"

	"github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2alpha1"
	"github.com/cilium/cilium/pkg/k8s/resource"
	"github.com/cilium/cilium/pkg/logging/logfields"
	subnetTable "github.com/cilium/cilium/pkg/maps/subnet"
	"github.com/cilium/cilium/pkg/node"
)

type watcherParams struct {
	cell.In

	Logger           *slog.Logger
	SubnetTable      statedb.RWTable[subnetTable.SubnetTableEntry]
	DB               *statedb.DB
	JobGroup         job.Group
	NodeWriter       *node.Writer                                      `optional:"true"`
	TopologyResource resource.Resource[*v2alpha1.CiliumSubnetTopology] `optional:"true"`
}

type SubnetWatcher struct {
	logger           *slog.Logger
	subnetTable      statedb.RWTable[subnetTable.SubnetTableEntry]
	db               *statedb.DB
	jobGroup         job.Group
	nodeWriter       *node.Writer
	topologyResource resource.Resource[*v2alpha1.CiliumSubnetTopology]
}

func newSubnetWatcher(params watcherParams) *SubnetWatcher {
	return &SubnetWatcher{
		logger:           params.Logger,
		subnetTable:      params.SubnetTable,
		db:               params.DB,
		jobGroup:         params.JobGroup,
		nodeWriter:       params.NodeWriter,
		topologyResource: params.TopologyResource,
	}
}

// SingletonTopologyName is the well-known name of the cluster-global
// CiliumSubnetTopology object that the agent reads. Objects with any other
// name are ignored (convention-only singleton).
const SingletonTopologyName = "default"

// run keeps the subnet table in sync with the singleton CiliumSubnetTopology,
// which is the only source of truth: it is applied while it exists, and its
// absence leaves the table empty. Objects with any other name are ignored
// (convention-only singleton).
func (w *SubnetWatcher) run(ctx context.Context, health cell.Health, synced chan<- struct{}) error {
	w.logger.Info("Starting subnet topology watcher")

	events := w.topologyResource.Events(ctx)

	// cr is the singleton object, or nil while it does not exist.
	var cr *v2alpha1.CiliumSubnetTopology

	storeSynced := false

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case ev, ok := <-events:
			if !ok {
				return ctx.Err()
			}

			matchesName := ev.Key.Name == SingletonTopologyName
			switch ev.Kind {
			case resource.Sync:
				storeSynced = true
			case resource.Upsert:
				if matchesName {
					cr = ev.Object
				} else {
					w.logger.Info("Ignoring non-singleton CiliumSubnetTopology object", logfields.Name, ev.Key.Name)
				}
			case resource.Delete:
				if matchesName {
					w.logger.Info("CiliumSubnetTopology singleton object deleted, clearing subnet topology")
					cr = nil
				}
			}
			ev.Done(nil)

			if !matchesName && ev.Kind != resource.Sync {
				continue
			}
		}

		err := w.applyEntries(ctx, cr)
		if err != nil {
			w.logger.Error("Failed to process subnet topology", logfields.Error, err)
			health.Degraded("Failed to process subnet topology", err)
		} else {
			health.OK("subnet topology processed successfully")
		}

		// The fence tracks the first processing attempt, not its outcome: a
		// topology we cannot decode must not hold up the datapath.
		if storeSynced && synced != nil {
			close(synced)
			synced = nil
		}
	}
}

// topologyEntry is a subnet prefix and the identity assigned to it, decoded
// from a CiliumSubnetTopology and ready to be written to the subnet table.
type topologyEntry struct {
	prefix   netip.Prefix
	identity uint32
}

// entriesFromTopology converts a CiliumSubnetTopology spec into positional
// subnet identities. A nil object yields no entries, which clears the table.
func entriesFromTopology(cst *v2alpha1.CiliumSubnetTopology) ([]topologyEntry, error) {
	if cst == nil {
		return nil, nil
	}
	var entries []topologyEntry
	for groupIdx, g := range cst.Spec.SubnetGroups {
		for _, cidr := range g.CIDRs {
			if !cidr.Prefix.IsValid() {
				return nil, fmt.Errorf("invalid CIDR in subnet group %q", g.Name)
			}
			// Identity is groupIdx + 1 to avoid using identity 0.
			entries = append(entries, topologyEntry{
				prefix:   cidr.Prefix,
				identity: uint32(groupIdx + 1),
			})
		}
	}
	return entries, nil
}

// applyEntries resets the subnet-identities table to exactly the identities
// derived from cst and triggers a re-evaluation of node routes.
func (w *SubnetWatcher) applyEntries(ctx context.Context, cst *v2alpha1.CiliumSubnetTopology) error {
	entries, err := entriesFromTopology(cst)
	if err != nil {
		return err
	}

	wTx := w.db.WriteTxn(w.subnetTable)
	defer wTx.Abort()

	if err := w.subnetTable.DeleteAll(wTx); err != nil {
		return fmt.Errorf("failed to reset subnet table: %w", err)
	}
	for _, e := range entries {
		entry := subnetTable.NewSubnetEntry(e.prefix, e.identity)
		if _, _, err := w.subnetTable.Insert(wTx, entry); err != nil {
			return fmt.Errorf("failed to upsert subnet entry %v: %w", entry, err)
		}
	}
	wTx.Commit()

	// Trigger re-evaluation of all node routes based on new topology.
	if w.nodeWriter != nil {
		if err := w.nodeWriter.Refresh(ctx, node.LinuxNodeReconciler); err != nil {
			return fmt.Errorf("refreshing nodes after subnet topology change: %w", err)
		}
	}

	return nil
}
