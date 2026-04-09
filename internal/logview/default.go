package logview

import "sync"

var (
	defaultMu      sync.RWMutex
	defaultManager *Manager
)

func SetDefault(m *Manager) {
	defaultMu.Lock()
	defaultManager = m
	defaultMu.Unlock()
}

func Default() *Manager {
	defaultMu.RLock()
	defer defaultMu.RUnlock()
	return defaultManager
}
