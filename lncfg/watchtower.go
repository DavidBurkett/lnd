package lncfg

import "github.com/ltcsuite/lnd/watchtower"

// Watchtower holds the daemon specific configuration parameters for running a
// watchtower that shares resources with the daemon.
//
//nolint:lll
type Watchtower struct {
	Active bool `long:"active" description:"If the watchtower should be active or not"`

	TowerDir string `long:"towerdir" description:"Directory of the watchtower.db"`

	watchtower.Conf
}

// DefaultWatchtowerCfg creates a Watchtower with some default values filled
// out.
func DefaultWatchtowerCfg(defaultTowerDir string) *Watchtower {
	conf := watchtower.DefaultConf()

	return &Watchtower{
		TowerDir: defaultTowerDir,
		Conf:     *conf,
	}
}
