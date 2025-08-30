package translations

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/config"
	"github.com/gin-gonic/gin"
)

type Store struct {
	mu    sync.RWMutex
	store map[string]map[string]string
}

func NewStore() *Store {
	return &Store{store: make(map[string]map[string]string)}
}

func (s *Store) Load(lang, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var messages map[string]string
	if err := json.Unmarshal(data, &messages); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.store[lang] == nil {
		s.store[lang] = make(map[string]string)
	}
	for k, v := range messages {
		s.store[lang][k] = v
	}

	return nil
}

func (s *Store) LoadAllModules(baseDir string) error {
	langs, err := os.ReadDir(baseDir)
	if err != nil {
		return err
	}

	for _, langDir := range langs {
		if !langDir.IsDir() {
			continue
		}
		lang := langDir.Name()
		files, err := os.ReadDir(filepath.Join(baseDir, lang))
		if err != nil {
			return err
		}

		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".json") {
				continue
			}
			path := filepath.Join(baseDir, lang, f.Name())
			if err := s.Load(lang, path); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Store) TL(locale Locale, key Key, data *map[string]any, fallback Locale) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	msg := s.store[string(locale)][string(key)]
	if msg == "" {
		msg = s.store[string(fallback)][string(key)]
		if msg == "" {
			return string(key)
		}
	}

	tmpl, _ := template.New("msg").Parse(msg)
	var sb strings.Builder
	_ = tmpl.Execute(&sb, data)

	return sb.String()
}

func (s *Store) TGin(c *gin.Context, key Key, vars *map[string]any) string {
	locale := GetLocaleFromGin(c, Locale(config.NewEnv().App.DefaultLocale))

	// fallback default: env
	return s.TL(locale, key, vars, LocaleEN)
}

func (s *Store) TContext(ctx context.Context, key Key, vars *map[string]any) string {
	locale := GetLocaleFromContext(ctx, LocaleEN) // fallback default

	msg := s.store[string(locale)][string(key)]
	if msg == "" {
		msg = s.store[string(LocaleEN)][string(key)]
		if msg == "" {
			return string(key)
		}
	}

	tmpl, _ := template.New("msg").Parse(msg)
	var sb strings.Builder
	_ = tmpl.Execute(&sb, vars)

	return sb.String()
}
