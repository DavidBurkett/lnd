package cluster

import (
	"github.com/btcsuite/btclog"
	"github.com/ltcsuite/lnd/build"
)

// Subsystem defines the logging code for this subsystem.
const Subsystem = "CLUS"

// log is a logger that is initialized with the btclog.Disabled logger.
var log btclog.Logger

// The default amount of logging is none.
func init() {
	UseLogger(build.NewSubLogger(Subsystem, nil))
}

// DisableLog disables all logging output.
func DisableLog() {
	UseLogger(btclog.Disabled)
}

// UseLogger uses a specified Logger to output package logging info.
func UseLogger(logger btclog.Logger) {
	log = logger
}
