package core

import (
	"sync"
	"sync/atomic"

	"github.com/homepi/homepi/pkg/libstr"
)

const (
	CORSDisabled uint32 = iota
	CORSEnabled
)

func (ctx *Context) CORSConfig() *CORSConfig {
	return &CORSConfig{
		Enabled: new(uint32),
	}
}

type CORSConfig struct {
	Enabled        *uint32  `json:"enabled"`
	AllowedOrigins []string `json:"allowed_origins,omitempty"`
	AllowedHeaders []string `json:"allowed_headers,omitempty"`
	sync.RWMutex   `json:"-"`
}

// IsEnabled returns the value of CORSConfig.isEnabled
func (c *CORSConfig) IsEnabled() bool {
	return atomic.LoadUint32(c.Enabled) == CORSEnabled
}

// IsValidOrigin determines if the origin of the request is allowed to make
// cross-origin requests based on the CORSConfig.
func (c *CORSConfig) IsValidOrigin(origin string) bool {
	// If we aren't enabling CORS then all origins are valid
	if !c.IsEnabled() {
		return true
	}

	c.RLock()
	defer c.RUnlock()

	if len(c.AllowedOrigins) == 0 {
		return false
	}

	if len(c.AllowedOrigins) == 1 && (c.AllowedOrigins)[0] == "*" {
		return true
	}

	return libstr.StrListContains(c.AllowedOrigins, origin)
}