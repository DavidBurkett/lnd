package kvdb

import (
	"github.com/btcsuite/btclog"
	"github.com/ltcsuite/lnd/kvdb/sqlbase"
)

// log is a logger that is initialized as disabled.  This means the package will
// not perform any logging by default until a logger is set.
var log = btclog.Disabled

// UseLogger uses a specified Logger to output package logging info.
func UseLogger(logger btclog.Logger) {
	log = logger

	sqlbase.UseLogger(log)
}
