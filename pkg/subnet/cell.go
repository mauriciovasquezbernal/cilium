// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package subnet

import (
	"context"

	"github.com/cilium/hive/cell"
	"github.com/cilium/hive/job"

	"github.com/cilium/cilium/pkg/endpoint/regeneration"
	k8sResources "github.com/cilium/cilium/pkg/k8s"
	"github.com/cilium/cilium/pkg/option"
)

// Cell provides the subnet watcher functionality
var Cell = cell.Module(
	"subnet",
	"Subnet watcher and management",

	cell.Provide(
		newSubnetWatcher,
		k8sResources.CiliumSubnetTopologyResource,
	),

	cell.Invoke(
		registerSubnetWatcher,
	),
)

func registerSubnetWatcher(cfg *option.DaemonConfig, fence regeneration.Fence, sw *SubnetWatcher) {
	if cfg.RoutingMode != option.RoutingModeHybrid {
		sw.logger.Debug("Routing mode is not hybrid, skipping subnet watcher")
		return
	}
	if sw.topologyResource == nil {
		sw.logger.Debug("Subnet topology resource is unavailable, skipping subnet watcher")
		return
	}

	synced := make(chan struct{})

	// Add a fence to ensure that subnet map is ready before starting endpoint regeneration.
	fence.Add("subnet-map", func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-synced:
			sw.logger.Info("Subnet topology synced")
			return nil
		}
	})

	sw.jobGroup.Add(job.OneShot("subnet-watcher", func(ctx context.Context, health cell.Health) error {
		return sw.run(ctx, health, synced)
	}))
}
