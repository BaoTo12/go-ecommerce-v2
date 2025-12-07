package featureflags

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"sync"
	"time"
)

/*
FEATURE FLAGS WITH A/B TESTING

Implements feature flags with:
- Percentage rollouts
- User targeting
- A/B testing variants
- Kill switches
- Scheduled releases
*/

// Flag represents a feature flag
type Flag struct {
	Key          string
	Description  string
	Enabled      bool
	Percentage   int           // 0-100 for gradual rollout
	Variants     []Variant     // For A/B testing
	Targeting    []TargetRule  // User targeting rules
	StartTime    *time.Time    // Scheduled start
	EndTime      *time.Time    // Scheduled end
	KillSwitch   bool          // Emergency disable
	Metadata     map[string]string
}

// Variant represents an A/B test variant
type Variant struct {
	Key        string
	Weight     int // Relative weight (e.g., 50/50 or 80/20)
	Payload    map[string]interface{}
}

// TargetRule for user targeting
type TargetRule struct {
	Attribute string   // e.g., "country", "plan", "user_id"
	Operator  string   // "eq", "neq", "in", "contains"
	Values    []string // Target values
}

// EvaluationContext contains context for flag evaluation
type EvaluationContext struct {
	UserID     string
	Attributes map[string]string
}

// FlagStore stores feature flags
type FlagStore interface {
	Get(ctx context.Context, key string) (*Flag, error)
	GetAll(ctx context.Context) ([]*Flag, error)
	Set(ctx context.Context, flag *Flag) error
	Delete(ctx context.Context, key string) error
}

// Client evaluates feature flags
type Client struct {
	store   FlagStore
	cache   map[string]*Flag
	cacheTTL time.Duration
	lastSync time.Time
	mu      sync.RWMutex
}

// NewClient creates a feature flag client
func NewClient(store FlagStore, cacheTTL time.Duration) *Client {
	return &Client{
		store:    store,
		cache:    make(map[string]*Flag),
		cacheTTL: cacheTTL,
	}
}

// IsEnabled checks if a flag is enabled for a user
func (c *Client) IsEnabled(ctx context.Context, key string, evalCtx *EvaluationContext) bool {
	flag, err := c.getFlag(ctx, key)
	if err != nil || flag == nil {
		return false
	}
	return c.evaluate(flag, evalCtx)
}

// GetVariant returns the variant for a user
func (c *Client) GetVariant(ctx context.Context, key string, evalCtx *EvaluationContext) *Variant {
	flag, err := c.getFlag(ctx, key)
	if err != nil || flag == nil {
		return nil
	}

	if !c.evaluate(flag, evalCtx) {
		return nil
	}

	if len(flag.Variants) == 0 {
		return nil
	}

	// Consistent hashing for variant selection
	return c.selectVariant(flag, evalCtx.UserID)
}

func (c *Client) getFlag(ctx context.Context, key string) (*Flag, error) {
	c.mu.RLock()
	if time.Since(c.lastSync) < c.cacheTTL {
		if flag, ok := c.cache[key]; ok {
			c.mu.RUnlock()
			return flag, nil
		}
	}
	c.mu.RUnlock()

	// Fetch from store
	flag, err := c.store.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.cache[key] = flag
	c.lastSync = time.Now()
	c.mu.Unlock()

	return flag, nil
}

func (c *Client) evaluate(flag *Flag, evalCtx *EvaluationContext) bool {
	// Kill switch
	if flag.KillSwitch {
		return false
	}

	// Base enabled check
	if !flag.Enabled {
		return false
	}

	// Time-based check
	now := time.Now()
	if flag.StartTime != nil && now.Before(*flag.StartTime) {
		return false
	}
	if flag.EndTime != nil && now.After(*flag.EndTime) {
		return false
	}

	// Targeting rules
	if len(flag.Targeting) > 0 && evalCtx != nil {
		if !c.matchesTargeting(flag.Targeting, evalCtx) {
			return false
		}
	}

	// Percentage rollout
	if flag.Percentage < 100 && flag.Percentage > 0 && evalCtx != nil {
		if !c.inPercentage(flag.Key, evalCtx.UserID, flag.Percentage) {
			return false
		}
	}

	return true
}

func (c *Client) matchesTargeting(rules []TargetRule, evalCtx *EvaluationContext) bool {
	for _, rule := range rules {
		value := evalCtx.Attributes[rule.Attribute]
		if evalCtx.UserID != "" && rule.Attribute == "user_id" {
			value = evalCtx.UserID
		}

		matched := false
		switch rule.Operator {
		case "eq":
			matched = len(rule.Values) > 0 && value == rule.Values[0]
		case "neq":
			matched = len(rule.Values) == 0 || value != rule.Values[0]
		case "in":
			for _, v := range rule.Values {
				if value == v {
					matched = true
					break
				}
			}
		case "contains":
			for _, v := range rule.Values {
				if containsString(value, v) {
					matched = true
					break
				}
			}
		}

		if !matched {
			return false
		}
	}
	return true
}

func (c *Client) inPercentage(flagKey, userID string, percentage int) bool {
	// Consistent hashing for percentage
	hash := sha256.Sum256([]byte(flagKey + userID))
	bucket := int(binary.BigEndian.Uint32(hash[:4]) % 100)
	return bucket < percentage
}

func (c *Client) selectVariant(flag *Flag, userID string) *Variant {
	if len(flag.Variants) == 0 {
		return nil
	}

	// Calculate total weight
	totalWeight := 0
	for _, v := range flag.Variants {
		totalWeight += v.Weight
	}

	if totalWeight == 0 {
		return &flag.Variants[0]
	}

	// Consistent hashing for variant
	hash := sha256.Sum256([]byte(flag.Key + "_variant_" + userID))
	bucket := int(binary.BigEndian.Uint32(hash[:4]) % uint32(totalWeight))

	// Select variant
	cumulative := 0
	for i := range flag.Variants {
		cumulative += flag.Variants[i].Weight
		if bucket < cumulative {
			return &flag.Variants[i]
		}
	}

	return &flag.Variants[0]
}

// Sync syncs all flags from store
func (c *Client) Sync(ctx context.Context) error {
	flags, err := c.store.GetAll(ctx)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*Flag)
	for _, flag := range flags {
		c.cache[flag.Key] = flag
	}
	c.lastSync = time.Now()

	return nil
}

// InMemoryFlagStore is an in-memory implementation
type InMemoryFlagStore struct {
	flags map[string]*Flag
	mu    sync.RWMutex
}

func NewInMemoryFlagStore() *InMemoryFlagStore {
	return &InMemoryFlagStore{
		flags: make(map[string]*Flag),
	}
}

func (s *InMemoryFlagStore) Get(ctx context.Context, key string) (*Flag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.flags[key], nil
}

func (s *InMemoryFlagStore) GetAll(ctx context.Context) ([]*Flag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	flags := make([]*Flag, 0, len(s.flags))
	for _, f := range s.flags {
		flags = append(flags, f)
	}
	return flags, nil
}

func (s *InMemoryFlagStore) Set(ctx context.Context, flag *Flag) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.flags[flag.Key] = flag
	return nil
}

func (s *InMemoryFlagStore) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.flags, key)
	return nil
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ABTest creates a simple A/B test flag
func ABTest(key, description string, controlWeight, treatmentWeight int) *Flag {
	return &Flag{
		Key:         key,
		Description: description,
		Enabled:     true,
		Percentage:  100,
		Variants: []Variant{
			{Key: "control", Weight: controlWeight},
			{Key: "treatment", Weight: treatmentWeight},
		},
	}
}

// GradualRollout creates a gradual rollout flag
func GradualRollout(key, description string, percentage int) *Flag {
	return &Flag{
		Key:         key,
		Description: description,
		Enabled:     true,
		Percentage:  percentage,
	}
}

// KillSwitch creates a kill switch flag
func KillSwitch(key, description string, enabled bool) *Flag {
	return &Flag{
		Key:         key,
		Description: description,
		Enabled:     enabled,
		KillSwitch:  !enabled,
	}
}
