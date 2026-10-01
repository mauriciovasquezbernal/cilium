// SPDX-License-Identifier: (GPL-2.0-only OR BSD-2-Clause)
/* Copyright Authors of Cilium */

/*
 * Hybrid routing off: the ipcache's skip_tunnel flag is the only input to the
 * tunnel decision, and the subnet map must be ignored even when it is
 * populated. See skip_tunnel_nodeport_nat_hybrid.c for the flag-cleared,
 * subnet-map-driven cases.
 */
#define ENABLE_HYBRID_ROUTING 0
#include "skip_tunnel_nodeport_nat_common.h"

PKTGEN(PROG_TYPE, "01_ipv4_nodeport_egress_no_flags")
int ipv4_nodeport_egress_no_flags_pktgen(struct __ctx_buff *ctx)
{
	return pktgen(ctx, true);
}

SETUP(PROG_TYPE, "01_ipv4_nodeport_egress_no_flags")
int ipv4_nodeport_egress_no_flags_setup(struct __ctx_buff *ctx)
{
	return setup(ctx, true, false);
}

CHECK(PROG_TYPE, "01_ipv4_nodeport_egress_no_flags")
int ipv4_nodeport_egress_no_flags_check(__maybe_unused const struct __ctx_buff *ctx)
{
	return check_ctx(ctx, true, CTX_ACT_REDIRECT);
}

PKTGEN(PROG_TYPE, "02_ipv4_nodeport_egress_skip_tunnel")
int ipv4_nodeport_egress_skip_tunnel_pktgen(struct __ctx_buff *ctx)
{
	return pktgen(ctx, true);
}

SETUP(PROG_TYPE, "02_ipv4_nodeport_egress_skip_tunnel")
int ipv4_nodeport_egress_skip_tunnel_setup(struct __ctx_buff *ctx)
{
	return setup(ctx, true, true);
}

CHECK(PROG_TYPE, "02_ipv4_nodeport_egress_skip_tunnel")
int ipv4_nodeport_egress_skip_tunnel_check(__maybe_unused const struct __ctx_buff *ctx)
{
	return check_ctx(ctx, true, CTX_ACT_DROP);
}

PKTGEN(PROG_TYPE, "03_ipv6_nodeport_egress_no_flags")
int ipv6_nodeport_egress_no_flags_pktgen(struct __ctx_buff *ctx)
{
	return pktgen(ctx, false);
}

SETUP(PROG_TYPE, "03_ipv6_nodeport_egress_no_flags")
int ipv6_nodeport_egress_no_flags_setup(struct __ctx_buff *ctx)
{
	return setup(ctx, false, false);
}

CHECK(PROG_TYPE, "03_ipv6_nodeport_egress_no_flags")
int ipv6_nodeport_egress_no_flags_check(__maybe_unused const struct __ctx_buff *ctx)
{
	return check_ctx(ctx, false, CTX_ACT_REDIRECT);
}

PKTGEN(PROG_TYPE, "04_ipv6_nodeport_egress_skip_tunnel")
int ipv6_nodeport_egress_skip_tunnel_pktgen(struct __ctx_buff *ctx)
{
	return pktgen(ctx, false);
}

SETUP(PROG_TYPE, "04_ipv6_nodeport_egress_skip_tunnel")
int ipv6_nodeport_egress_skip_tunnel_setup(struct __ctx_buff *ctx)
{
	return setup(ctx, false, true);
}

CHECK(PROG_TYPE, "04_ipv6_nodeport_egress_skip_tunnel")
int ipv6_nodeport_egress_skip_tunnel_check(__maybe_unused const struct __ctx_buff *ctx)
{
	return check_ctx(ctx, false, CTX_ACT_DROP);
}

/*
 * With hybrid routing off the map must not be consulted at all, so the tunnel
 * is still taken.
 */
PKTGEN(PROG_TYPE, "05_ipv4_nodeport_egress_same_subnet_no_hybrid")
int ipv4_nodeport_egress_same_subnet_no_hybrid_pktgen(struct __ctx_buff *ctx)
{
	return pktgen(ctx, true);
}

SETUP(PROG_TYPE, "05_ipv4_nodeport_egress_same_subnet_no_hybrid")
int ipv4_nodeport_egress_same_subnet_no_hybrid_setup(struct __ctx_buff *ctx)
{
	setup_subnet_table(true, SAME_SUBNET_ID, SAME_SUBNET_ID);
	return setup(ctx, true, false);
}

CHECK(PROG_TYPE, "05_ipv4_nodeport_egress_same_subnet_no_hybrid")
int ipv4_nodeport_egress_same_subnet_no_hybrid_check(__maybe_unused const struct __ctx_buff *ctx)
{
	return check_ctx(ctx, true, CTX_ACT_REDIRECT);
}

PKTGEN(PROG_TYPE, "06_ipv6_nodeport_egress_same_subnet_no_hybrid")
int ipv6_nodeport_egress_same_subnet_no_hybrid_pktgen(struct __ctx_buff *ctx)
{
	return pktgen(ctx, false);
}

SETUP(PROG_TYPE, "06_ipv6_nodeport_egress_same_subnet_no_hybrid")
int ipv6_nodeport_egress_same_subnet_no_hybrid_setup(struct __ctx_buff *ctx)
{
	setup_subnet_table(false, SAME_SUBNET_ID, SAME_SUBNET_ID);
	return setup(ctx, false, false);
}

CHECK(PROG_TYPE, "06_ipv6_nodeport_egress_same_subnet_no_hybrid")
int ipv6_nodeport_egress_same_subnet_no_hybrid_check(__maybe_unused const struct __ctx_buff *ctx)
{
	return check_ctx(ctx, false, CTX_ACT_REDIRECT);
}
