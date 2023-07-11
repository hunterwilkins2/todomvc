package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

type TemplateCache struct {
	cache map[string]*template.Template
}

func New(dir string) (*TemplateCache, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	partials, err := filepath.Glob(filepath.Join(dir, "*.partial.html"))
	if err != nil {
		return nil, err
	}

	for _, partial := range partials {
		name := filepath.Base(partial)
		ts, err := template.New(name).ParseFiles(partial)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return &TemplateCache{cache: cache}, nil
}

func (t *TemplateCache) Render(w http.ResponseWriter, name string, data any) error {
	ts, ok := t.cache[name]
	if !ok {
		return fmt.Errorf("template not found: %s", name)
	}

	buff := new(bytes.Buffer)
	err := ts.Execute(buff, data)
	if err != nil {
		return err
	}

	if _, err = buff.WriteTo(w); err != nil {
		return err
	}

	return nil
}
