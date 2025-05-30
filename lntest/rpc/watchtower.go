package rpc

import (
	"context"

	"github.com/ltcsuite/lnd/lnrpc/watchtowerrpc"
	"github.com/ltcsuite/lnd/lnrpc/wtclientrpc"
)

// =====================
// WatchtowerClient and WatchtowerClientClient related RPCs.
// =====================

// GetInfoWatchtower makes a RPC call to the watchtower of the given node and
// asserts.
func (h *HarnessRPC) GetInfoWatchtower() *watchtowerrpc.GetInfoResponse {
	ctxt, cancel := context.WithTimeout(h.runCtx, DefaultTimeout)
	defer cancel()

	req := &watchtowerrpc.GetInfoRequest{}
	info, err := h.Watchtower.GetInfo(ctxt, req)
	h.NoError(err, "GetInfo from Watchtower")

	return info
}

// GetTowerInfo makes an RPC call to the watchtower client of the given node and
// asserts.
func (h *HarnessRPC) GetTowerInfo(
	req *wtclientrpc.GetTowerInfoRequest) *wtclientrpc.Tower {

	ctxt, cancel := context.WithTimeout(h.runCtx, DefaultTimeout)
	defer cancel()

	info, err := h.WatchtowerClient.GetTowerInfo(ctxt, req)
	h.NoError(err, "GetTowerInfo from WatchtowerClient")

	return info
}

// AddTower makes a RPC call to the WatchtowerClient of the given node and
// asserts.
func (h *HarnessRPC) AddTower(
	req *wtclientrpc.AddTowerRequest) *wtclientrpc.AddTowerResponse {

	ctxt, cancel := context.WithTimeout(h.runCtx, DefaultTimeout)
	defer cancel()

	resp, err := h.WatchtowerClient.AddTower(ctxt, req)
	h.NoError(err, "AddTower")

	return resp
}

// WatchtowerStats makes a RPC call to the WatchtowerClient of the given node
// and asserts.
func (h *HarnessRPC) WatchtowerStats() *wtclientrpc.StatsResponse {
	ctxt, cancel := context.WithTimeout(h.runCtx, DefaultTimeout)
	defer cancel()

	req := &wtclientrpc.StatsRequest{}
	resp, err := h.WatchtowerClient.Stats(ctxt, req)
	h.NoError(err, "Stats from Watchtower")

	return resp
}
