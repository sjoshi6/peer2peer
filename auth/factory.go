package controllers

import "github.com/dgrijalva/jwt-go"

// NewAuthManager is a factory for auth managers
func NewAuthManager(key []byte, configs ...Config) *Manager {

	/* The function uses ... config to ensure it can be 0 or more */
	var c Config

	if len(configs) == 0 {
		// If no configs provided
		c = Config{}
	} else {
		// Set the config value to first amongst arr
		c = configs[0]
	}
	m := &Manager{
		key:    key,
		method: c.Method,
		ttl:    c.TTL,
	}
	m.setDefaults()
	return m
}

// setDefaults is a local function to assign defaults to auth-manager object
func (m *Manager) setDefaults() {
	if m.method == nil {
		m.method = jwt.SigningMethodHS256
	}
	if m.ttl == 0 {
		m.ttl = defaultTTL
	}
}
