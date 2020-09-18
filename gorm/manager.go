package gorm

import (
	"sync"

	"github.com/leon-yc/gopher/qulibs"
)

// Manager manages multi client instances of gorm for easy usage.
type Manager struct {
	clients sync.Map
	configs sync.Map
}

// NewManager creates a new manager store of redis with configs.
func NewManager(configs *ManagerConfig) *Manager {
	mgr := &Manager{}
	mgr.Load(configs)

	return mgr
}

// NewClient finds or creates a mysql client registered with the name. It uses a dummy logger for default.
func (mgr *Manager) NewClient(name string) (client *Client, err error) {
	return mgr.NewClientWithLogger(name, qulibs.NewDummyLogger())
}

// NewClientWithLogger finds or creates a mysql client registered with the name and logger given
func (mgr *Manager) NewClientWithLogger(name string, log qulibs.Logger) (client *Client, err error) {
	// first, try clients store
	mgrclient, ok := mgr.clients.Load(name)
	if ok {
		client, ok := mgrclient.(*Client)
		if ok {
			client.log = log

			return client, nil
		}
	}

	// second, try creating a new client with config registered with the name.
	config, err := mgr.Config(name)
	if err != nil {
		return nil, err
	}

	// 1, create a new client
	client, err = NewWithLogger(config, log)
	if err != nil {
		return nil, err
	}

	// 2, store the client with the name
	mgr.clients.Store(name, client)

	return client, nil
}

// Config returns a config registered with the name given
func (mgr *Manager) Config(name string) (config *Config, err error) {
	mgrconfig, ok := mgr.configs.Load(name)
	if ok {
		config, ok := mgrconfig.(*Config)
		if ok {
			return config, nil
		}

		return nil, ErrInvalidConfig
	}

	return nil, ErrNotFoundConfig
}

// Add registers a new config of redis with the name given.
//
// NOTE: It will remove client related to the name if existed.
func (mgr *Manager) Add(name string, config *Config) {
	config.FillWithDefaults()

	// store new config
	mgr.configs.Store(name, config)

	// remove old client
	mgr.clients.Delete(name)
}

// Del removes both client and config of redis registered with the name given.
func (mgr *Manager) Del(name string) {
	mgr.configs.Delete(name)
	mgr.clients.Delete(name)
}

// Load registers all configs with its name defined by ManagerConfig
func (mgr *Manager) Load(configs *ManagerConfig) {
	if configs == nil {
		return
	}

	for name, config := range *configs {
		mgr.Add(name, config)
	}
}
