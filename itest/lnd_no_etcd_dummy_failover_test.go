//go:build !kvdb_etcd
// +build !kvdb_etcd

package itest

import "github.com/ltcsuite/lnd/lntest"

// testEtcdFailover is an empty itest when LND is not compiled with etcd
// support.
func testEtcdFailover(ht *lntest.HarnessTest) {}
