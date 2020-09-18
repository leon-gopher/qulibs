package gorm

// defaults
const (
	MaxDialTimeout   = 1000 // millisecond
	MaxReadTimeout   = 3000 // millisecond
	MaxWriteTimeout  = 5000 // millisecond
	MaxOpenConn      = 128
	MaxIdleConn      = 16
	MaxLifecycleConn = 300 // in second
)
