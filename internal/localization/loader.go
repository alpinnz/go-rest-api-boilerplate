package localization

import (
	"encoding/json"
	"os"
	"sync"
)

// Bundle manages multiple language dictionaries with thread-safe access.
type Bundle struct {
	defaultLang string
	mu          sync.RWMutex
	data        map[string]Dictionary
}

// NewBundle creates a new language bundle with specified default language.
func NewBundle(defaultLang string) *Bundle {
	return &Bundle{
		defaultLang: defaultLang,
		data:        map[string]Dictionary{},
	}
}

// Load reads and parses JSON language file into the bundle.
// Returns error if file cannot be read or JSON is invalid.
func (b *Bundle) Load(langCode, filePath string) error {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var dict Dictionary
	if err := json.Unmarshal(raw, &dict); err != nil {
		return err
	}

	b.mu.Lock()
	b.data[langCode] = dict
	b.mu.Unlock()

	return nil
}

// Get retrieves dictionary for specified language code.
// Falls back to default language if requested language not found.
func (b *Bundle) Get(langCode string) Dictionary {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if d, ok := b.data[langCode]; ok {
		return d
	}
	return b.data[b.defaultLang]
}

// DefaultLang returns the default language code.
func (b *Bundle) DefaultLang() string {
	return b.defaultLang
}
