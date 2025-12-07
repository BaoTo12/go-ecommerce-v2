package security

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

var (
	ErrSecretNotFound = errors.New("secret not found")
	ErrSecretExpired  = errors.New("secret expired")
)

// Secret represents a secret value
type Secret struct {
	Key       string    `json:"key"`
	Value     string    `json:"-"` // Never serialize the actual value
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// SecretStore interface for secret management
type SecretStore interface {
	Get(key string) (*Secret, error)
	Set(key, value string, expiresIn time.Duration) error
	Delete(key string) error
	List() ([]string, error)
	Rotate(key string, newValue string) error
}

// MemorySecretStore stores secrets in memory (dev only)
type MemorySecretStore struct {
	secrets map[string]*Secret
	mu      sync.RWMutex
}

// NewMemorySecretStore creates an in-memory secret store
func NewMemorySecretStore() *MemorySecretStore {
	return &MemorySecretStore{
		secrets: make(map[string]*Secret),
	}
}

func (s *MemorySecretStore) Get(key string) (*Secret, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	secret, ok := s.secrets[key]
	if !ok {
		return nil, ErrSecretNotFound
	}

	if !secret.ExpiresAt.IsZero() && time.Now().After(secret.ExpiresAt) {
		return nil, ErrSecretExpired
	}

	return secret, nil
}

func (s *MemorySecretStore) Set(key, value string, expiresIn time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	version := 1
	if existing, ok := s.secrets[key]; ok {
		version = existing.Version + 1
	}

	secret := &Secret{
		Key:       key,
		Value:     value,
		Version:   version,
		CreatedAt: time.Now(),
	}

	if expiresIn > 0 {
		secret.ExpiresAt = time.Now().Add(expiresIn)
	}

	s.secrets[key] = secret
	return nil
}

func (s *MemorySecretStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.secrets, key)
	return nil
}

func (s *MemorySecretStore) List() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.secrets))
	for k := range s.secrets {
		keys = append(keys, k)
	}
	return keys, nil
}

func (s *MemorySecretStore) Rotate(key string, newValue string) error {
	return s.Set(key, newValue, 0)
}

// EnvSecretStore reads secrets from environment variables
type EnvSecretStore struct {
	prefix string
	cache  sync.Map
}

// NewEnvSecretStore creates an environment-based secret store
func NewEnvSecretStore(prefix string) *EnvSecretStore {
	return &EnvSecretStore{prefix: prefix}
}

func (s *EnvSecretStore) Get(key string) (*Secret, error) {
	envKey := s.prefix + key
	value := os.Getenv(envKey)
	if value == "" {
		return nil, ErrSecretNotFound
	}

	return &Secret{
		Key:       key,
		Value:     value,
		Version:   1,
		CreatedAt: time.Now(),
	}, nil
}

func (s *EnvSecretStore) Set(key, value string, expiresIn time.Duration) error {
	return os.Setenv(s.prefix+key, value)
}

func (s *EnvSecretStore) Delete(key string) error {
	return os.Unsetenv(s.prefix + key)
}

func (s *EnvSecretStore) List() ([]string, error) {
	return nil, errors.New("not supported for env store")
}

func (s *EnvSecretStore) Rotate(key string, newValue string) error {
	return s.Set(key, newValue, 0)
}

// SecretManager manages secrets with caching and rotation
type SecretManager struct {
	store    SecretStore
	cache    sync.Map
	cacheTTL time.Duration
}

// NewSecretManager creates a secret manager
func NewSecretManager(store SecretStore, cacheTTL time.Duration) *SecretManager {
	return &SecretManager{
		store:    store,
		cacheTTL: cacheTTL,
	}
}

type cachedSecret struct {
	secret    *Secret
	expiresAt time.Time
}

func (m *SecretManager) Get(key string) (string, error) {
	// Check cache
	if cached, ok := m.cache.Load(key); ok {
		cs := cached.(*cachedSecret)
		if time.Now().Before(cs.expiresAt) {
			return cs.secret.Value, nil
		}
		m.cache.Delete(key)
	}

	// Fetch from store
	secret, err := m.store.Get(key)
	if err != nil {
		return "", err
	}

	// Cache it
	m.cache.Store(key, &cachedSecret{
		secret:    secret,
		expiresAt: time.Now().Add(m.cacheTTL),
	})

	return secret.Value, nil
}

func (m *SecretManager) GetBytes(key string) ([]byte, error) {
	value, err := m.Get(key)
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(value)
}

func (m *SecretManager) Set(key, value string) error {
	return m.store.Set(key, value, 0)
}

func (m *SecretManager) Rotate(key, newValue string) error {
	if err := m.store.Rotate(key, newValue); err != nil {
		return err
	}
	m.cache.Delete(key)
	return nil
}

// EncryptedSecretStore wraps a store with encryption
type EncryptedSecretStore struct {
	inner     SecretStore
	key       []byte
}

// NewEncryptedSecretStore creates an encrypted secret store
func NewEncryptedSecretStore(inner SecretStore, encryptionKey []byte) *EncryptedSecretStore {
	return &EncryptedSecretStore{
		inner: inner,
		key:   encryptionKey,
	}
}

// Note: Actual encryption implementation would use crypto/aes + crypto/cipher
// This is a placeholder showing the pattern

func (s *EncryptedSecretStore) Get(key string) (*Secret, error) {
	secret, err := s.inner.Get(key)
	if err != nil {
		return nil, err
	}
	// In real impl: decrypt secret.Value here
	return secret, nil
}

func (s *EncryptedSecretStore) Set(key, value string, expiresIn time.Duration) error {
	// In real impl: encrypt value here
	return s.inner.Set(key, value, expiresIn)
}

func (s *EncryptedSecretStore) Delete(key string) error {
	return s.inner.Delete(key)
}

func (s *EncryptedSecretStore) List() ([]string, error) {
	return s.inner.List()
}

func (s *EncryptedSecretStore) Rotate(key string, newValue string) error {
	// In real impl: encrypt newValue here
	return s.inner.Rotate(key, newValue)
}

// LoadSecretsFromJSON loads secrets from a JSON file
func LoadSecretsFromJSON(store SecretStore, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var secrets map[string]string
	if err := json.Unmarshal(data, &secrets); err != nil {
		return err
	}

	for key, value := range secrets {
		if err := store.Set(key, value, 0); err != nil {
			return err
		}
	}

	return nil
}
