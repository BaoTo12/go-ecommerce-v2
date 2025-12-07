package config

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// HotConfig supports runtime configuration updates
type HotConfig struct {
	path      string
	data      atomic.Value
	watchers  []func(interface{})
	hash      string
	mu        sync.RWMutex
	unmarshal func([]byte) (interface{}, error)
}

// NewHotConfig creates a hot-reloadable config
func NewHotConfig(path string, unmarshal func([]byte) (interface{}, error)) (*HotConfig, error) {
	hc := &HotConfig{
		path:      path,
		watchers:  make([]func(interface{}), 0),
		unmarshal: unmarshal,
	}

	// Initial load
	if err := hc.reload(); err != nil {
		return nil, err
	}

	return hc, nil
}

// Get returns the current config
func (hc *HotConfig) Get() interface{} {
	return hc.data.Load()
}

// OnChange registers a callback for config changes
func (hc *HotConfig) OnChange(fn func(interface{})) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.watchers = append(hc.watchers, fn)
}

// Watch starts watching for file changes
func (hc *HotConfig) Watch(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			hc.checkAndReload()
		}
	}
}

func (hc *HotConfig) checkAndReload() {
	data, err := os.ReadFile(hc.path)
	if err != nil {
		return
	}

	hash := hashBytes(data)
	hc.mu.RLock()
	changed := hash != hc.hash
	hc.mu.RUnlock()

	if changed {
		hc.reload()
	}
}

func (hc *HotConfig) reload() error {
	data, err := os.ReadFile(hc.path)
	if err != nil {
		return err
	}

	config, err := hc.unmarshal(data)
	if err != nil {
		return err
	}

	hc.mu.Lock()
	hc.hash = hashBytes(data)
	watchers := make([]func(interface{}), len(hc.watchers))
	copy(watchers, hc.watchers)
	hc.mu.Unlock()

	hc.data.Store(config)

	// Notify watchers
	for _, fn := range watchers {
		go fn(config)
	}

	return nil
}

func hashBytes(data []byte) string {
	h := sha256.Sum256(data)
	return string(h[:])
}

// FeatureFlags manages runtime feature toggles
type FeatureFlags struct {
	flags map[string]*Flag
	mu    sync.RWMutex
}

// Flag represents a feature flag
type Flag struct {
	Name        string         `json:"name"`
	Enabled     bool           `json:"enabled"`
	Percentage  int            `json:"percentage,omitempty"`  // Rollout percentage
	Whitelist   []string       `json:"whitelist,omitempty"`   // Always enabled for these
	Blacklist   []string       `json:"blacklist,omitempty"`   // Always disabled for these
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// NewFeatureFlags creates a new feature flag manager
func NewFeatureFlags() *FeatureFlags {
	return &FeatureFlags{
		flags: make(map[string]*Flag),
	}
}

// Set sets a feature flag
func (ff *FeatureFlags) Set(flag *Flag) {
	ff.mu.Lock()
	defer ff.mu.Unlock()
	ff.flags[flag.Name] = flag
}

// IsEnabled checks if a flag is enabled
func (ff *FeatureFlags) IsEnabled(name string) bool {
	ff.mu.RLock()
	defer ff.mu.RUnlock()

	flag, exists := ff.flags[name]
	if !exists {
		return false
	}
	return flag.Enabled
}

// IsEnabledFor checks if flag is enabled for a specific user
func (ff *FeatureFlags) IsEnabledFor(name, userID string) bool {
	ff.mu.RLock()
	defer ff.mu.RUnlock()

	flag, exists := ff.flags[name]
	if !exists {
		return false
	}

	// Check blacklist first
	for _, id := range flag.Blacklist {
		if id == userID {
			return false
		}
	}

	// Check whitelist
	for _, id := range flag.Whitelist {
		if id == userID {
			return true
		}
	}

	// Check percentage rollout
	if flag.Percentage > 0 && flag.Percentage < 100 {
		hash := hashString(userID + name)
		return hash%100 < flag.Percentage
	}

	return flag.Enabled
}

func hashString(s string) int {
	h := 0
	for _, c := range s {
		h = 31*h + int(c)
	}
	if h < 0 {
		h = -h
	}
	return h
}

// LoadFromJSON loads flags from JSON
func (ff *FeatureFlags) LoadFromJSON(data []byte) error {
	var flags map[string]*Flag
	if err := json.Unmarshal(data, &flags); err != nil {
		return err
	}

	ff.mu.Lock()
	defer ff.mu.Unlock()
	ff.flags = flags
	return nil
}

// All returns all flags
func (ff *FeatureFlags) All() map[string]*Flag {
	ff.mu.RLock()
	defer ff.mu.RUnlock()

	result := make(map[string]*Flag)
	for k, v := range ff.flags {
		result[k] = v
	}
	return result
}

// DynamicConfig combines hot config with feature flags
type DynamicConfig struct {
	HotConfig    *HotConfig
	FeatureFlags *FeatureFlags
}

// NewDynamicConfig creates a dynamic config manager
func NewDynamicConfig(configPath string) (*DynamicConfig, error) {
	hc, err := NewHotConfig(configPath, func(data []byte) (interface{}, error) {
		var config map[string]interface{}
		err := json.Unmarshal(data, &config)
		return config, err
	})
	if err != nil {
		return nil, err
	}

	ff := NewFeatureFlags()

	dc := &DynamicConfig{
		HotConfig:    hc,
		FeatureFlags: ff,
	}

	// Sync feature flags from config
	hc.OnChange(func(cfg interface{}) {
		if m, ok := cfg.(map[string]interface{}); ok {
			if flags, ok := m["features"].(map[string]interface{}); ok {
				for name, enabled := range flags {
					if e, ok := enabled.(bool); ok {
						ff.Set(&Flag{Name: name, Enabled: e})
					}
				}
			}
		}
	})

	return dc, nil
}
