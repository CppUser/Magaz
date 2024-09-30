package templates

import (
	"Magaz/backend/pkg/utils/parser"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type PageMetadata struct {
	Page       string   `mapstructure:"page"`       // Page file
	Layout     string   `mapstructure:"layout"`     // General layout file
	Components []string `mapstructure:"components"` // List of component files
}

type Metadata struct {
	Pages []PageMetadata `mapstructure:"pages"` // List of page metadata entries
}

type TemplateCache struct {
	cache         map[string]*template.Template
	mu            sync.RWMutex
	layoutDir     string
	pagesDir      string
	componentsDir string
	metadata      *Metadata
}

func NewTemplateCache(layoutDir, pagesDir, componentsDir string) (*TemplateCache, error) {

	metadata, err := LoadTemplateConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load template metadata: %w", err)
	}

	return &TemplateCache{
		cache:         make(map[string]*template.Template),
		layoutDir:     layoutDir,
		pagesDir:      pagesDir,
		componentsDir: componentsDir,
		metadata:      metadata,
	}, nil
}

func LoadTemplateConfig() (*Metadata, error) {
	var cfg Metadata

	configPaths := []string{
		".",
		"backend/config/",
		"backend/config/frontend",
	}

	if err := parser.Load("page_config", "yaml", configPaths, &cfg); err != nil {
		log.Fatalf("Failed to load API configs: %v", err)
	}

	return &cfg, nil
}

func (tm *TemplateCache) GetTemplate(name string, components ...string) (*template.Template, error) {

	tm.mu.RLock()
	if tmpl, ok := tm.cache[name]; ok {
		tm.mu.RUnlock()
		return tmpl, nil
	}
	tm.mu.RUnlock()

	return tm.loadCachedTemplate(name, components...)
}

func (tc *TemplateCache) loadCachedTemplate(name string, components ...string) (*template.Template, error) {
	tc.mu.Lock() // Lock when modifying the cache
	defer tc.mu.Unlock()

	// Check if the template is already cached
	if tmpl, ok := tc.cache[name]; ok {
		return tmpl, nil
	}

	var pageMetadata *PageMetadata
	for _, entry := range tc.metadata.Pages {
		if entry.Page == name {
			pageMetadata = &entry
			break
		}
	}

	if pageMetadata == nil {
		return nil, fmt.Errorf("no metadata found for page: %s", name)
	}

	// Load layout files
	layoutPath := filepath.Join(tc.layoutDir, pageMetadata.Layout)
	if _, err := os.Stat(layoutPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("layout file %s not found", layoutPath)
	}

	// Load component files
	var componentPaths []string
	for _, component := range pageMetadata.Components {
		componentPath := filepath.Join(tc.componentsDir, component)
		if _, err := os.Stat(componentPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("component file %s not found", componentPath)
		}
		componentPaths = append(componentPaths, componentPath)
	}

	// Load the main page file
	pageFile := filepath.Join(tc.pagesDir, name)
	if _, err := os.Stat(pageFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("page file %s not found", pageFile)
	}

	// Parse all files together: layout, components, and page
	filesToParse := append([]string{layoutPath}, componentPaths...)
	filesToParse = append(filesToParse, pageFile)

	// Debug: Log the files being parsed
	fmt.Println("All files to parse:", filesToParse)

	// Create a new template and add the dict function to the function map
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"dict": dict, // Add the dict function
	}).ParseFiles(filesToParse...) // Parse the template files with the function map
	if err != nil {
		return nil, err
	}

	// Cache the parsed template for future use
	tc.cache[name] = tmpl

	return tmpl, nil
}

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("invalid dict call, must have an even number of arguments")
	}
	dict := make(map[string]interface{})
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func CreateTemplateCache(layoutDir string, pagesDir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Load all layout files
	layouts, err := filepath.Glob(filepath.Join(layoutDir, "*layout.gohtml"))
	if err != nil {
		return nil, err
	}

	// Load all page files
	pages, err := filepath.Glob(filepath.Join(pagesDir, "*.gohtml"))
	if err != nil {
		return nil, err
	}

	// Parse each page with the layout
	for _, page := range pages {
		// Extract the template name
		name := filepath.Base(page)

		// Parse the layout files and the page file
		tmpl, err := template.ParseFiles(append(layouts, page)...)
		if err != nil {
			return nil, err
		}

		// Store the parsed template in the cache
		cache[name] = tmpl
	}

	return cache, nil
}
