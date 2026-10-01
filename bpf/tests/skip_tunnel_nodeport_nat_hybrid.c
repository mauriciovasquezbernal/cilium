// SPDX-License-Identifier: (GPL-2.0-only OR BSD-2-Clause)
/* Copyright Authors of Cilium */

/*
 * Hybrid routing on. The ipcache's skip_tunnel flag is CLEARED in every case
 * below, so the subnet map is the only thing that may suppress encapsulation:
 * tail_nodeport_nat_egress_ipv{4,6}() must not encapsulate towards a backend
 * that shares a subnet group with this node.
 *
 * See skip_tunnel_nodeport_nat.c for the same scenarios with hybrid routing
 * off, where the flag alone decides.
 */
#define ENABLE_HYBRID_ROUTING 1
#include "skip_tunnel_nodeport_nat_common.h"

PKTGEN(PROG_TYPE, "01_ipv4_nodeport_hybrid_same_subnet")
int ipv4_nodeport_hybrid_same_subnet_pktgen(struct __ctx_buff *ctx)
{
	return pktgen(ctx, true);
}

SETUP(PROG_TYPE, "01_ipv4_nodeport_hybrid_same_subnet")
int ipv4_nodeport_hybrid_same_subnet_setup(struct __ctx_buff *ctx)
{
	setup_subnet_table(true, SAME_SUBNET_ID, SAME_SUBNET_ID);
	return setup(ctx, true, false);
}

CHECK(PROG_TYPE, "01_ipv4_nodeport_hybrid_same_subnet")
int ipv4_nodeport_hybrid_same_subnet_check(__maybe_unused const struct __ctx_buff *ctx)
{
	return check_ctx(ctx, true, CTX_ACT_DROP);
}

PKTGEN(PROG_TYPE, "02_ipv4_nodeport_hybrid_other_subnet")
int ipv4_nodeport_hybrid_other_subnet_pktgen(struct __ctx_buff *ctx)
{
	return pktgen(ctx, true);
}

SETUP(PROG_TYPE, "02_ipv4_nodeport_hybrid_other_subnet")
int ipv4_nodeport_hybrid_other_subnet_setup(struct __ctx_buff *ctx)
{
	setup_subnet_table(true, SAME_SUBNET_ID, OTHER_SUBNET_ID);
	return setup(ctx, true, false);
}

CHECK(PROG_TYPE, "02_ipv4_nodeport_hybrid_other_subnet")
int ipv4_nodeport_hybrid_other_subnet_check(__maybe_unused const struct __ctx_buff *ctx)
{
	return check_ctx(ctx, true, CTX_ACT_REDIRECT);
}

PKTGEN(PROG_TYPE, "03_ipv4_nodeport_hybrid_unmapped_dst")
int ipv4_nodeport_hybrid_unmapped_dst_pktgen(struct __ctx_buff *ctx)
{
	return pktgen(ctx, true);
}

SETUP(PROG_TYPE, "03_ipv4_nodeport_hybrid_unmapped_dst")
int ipv4_nodeport_hybrid_unmapped_dst_setup(struct __ctx_buff *ctx)
{
	setup_subnet_table(true, SAME_SUBNET_ID, 0);
	return setup(ctx, true, false);
}

CHECK(PROG_TYPE, "03_ipv4_nodeport_hybrid_unmapped_dst")
int ipv4_nodeport_hybrid_unmapped_dst_check(__maybe_unused const struct __ctx_buff *ctx)
{
	return check_ctx(ctx, true, CTX_ACT_REDIRECT);
}

/* IPv6 mirrors of the three cases above. */
PKTGEN(PROG_TYPE, "04_ipv6_nodeport_hybrid_same_subnet")
int ipv6_nodeport_hybrid_same_subnet_pktgen(struct __ctx_buff *ctx)
{
	return pktgen(ctx, false);
}

SETUP(PROG_TYPE, "04_ipv6_nodeport_hybrid_same_subnet")
int ipv6_nodeport_hybrid_same_subnet_setup(struct __ctx_buff *ctx)
{
	setup_subnet_table(false, SAME_SUBNET_ID, SAME_SUBNET_ID);
	return setup(ctx, false, false);
}

CHECK(PROG_TYPE, "04_ipv6_nodeport_hybrid_same_subnet")
int ipv6_nodeport_hybrid_same_subnet_check(__maybe_unused const struct __ctx_buff *ctx)
{
	return check_ctx(ctx, false, CTX_ACT_DROP);
}

PKTGEN(PROG_TYPE, "05_ipv6_nodeport_hybrid_other_subnet")
int ipv6_nodeport_hybrid_other_subnet_pktgen(struct __ctx_buff *ctx)
{
	return pktgen(ctx, false);
}

SETUP(PROG_TYPE, "05_ipv6_nodeport_hybrid_other_subnet")
int ipv6_nodeport_hybrid_other_subnet_setup(struct __ctx_buff *ctx)
{
	setup_subnet_table(false, SAME_SUBNET_ID, OTHER_SUBNET_ID);
	return setup(ctx, false, false);
}

CHECK(PROG_TYPE, "05_ipv6_nodeport_hybrid_other_subnet")
int ipv6_nodeport_hybrid_other_subnet_check(__maybe_unused const struct __ctx_buff *ctx)
{
	return check_ctx(ctx, false, CTX_ACT_REDIRECT);
}

PKTGEN(PROG_TYPE, "06_ipv6_nodeport_hybrid_unmapped_dst")
int ipv6_nodeport_hybrid_unmapped_dst_pktgen(struct __ctx_buff *ctx)
{
	return pktgen(ctx, false);
}

SETUP(PROG_TYPE, "06_ipv6_nodeport_hybrid_unmapped_dst")
int ipv6_nodeport_hybrid_unmapped_dst_setup(struct __ctx_buff *ctx)
{
	setup_subnet_table(false, SAME_SUBNET_ID, 0);
	return setup(ctx, false, false);
}

CHECK(PROG_TYPE, "06_ipv6_nodeport_hybrid_unmapped_dst")
int ipv6_nodeport_hybrid_unmapped_dst_check(__maybe_unused const struct __ctx_buff *ctx)
{
	return check_ctx(ctx, false, CTX_ACT_REDIRECT);
}
