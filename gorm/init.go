package gorm

import (
	"github.com/leon-yc/gopher/qulibs"

	// for mysql
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// TODO: We should use mysql cluster proxy instead of local DefaultMgr!!!
var (
	DefaultMgr *Manager
)

func init() {
	DefaultMgr = NewManager(nil)
}

// NewClient returns a valid gorm instance for given name registered in DefaultMgr
func NewClient(name string, log qulibs.Logger) (*Client, error) {
	return DefaultMgr.NewClientWithLogger(name, log)
}
