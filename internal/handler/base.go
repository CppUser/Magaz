package handler

import (
	"Magaz/internal/config"
	"Magaz/pkg/bot/telegram"
	"github.com/gorilla/sessions"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"html/template"
	"path/filepath"
)

// TODO: Needs to move to more global space so everyone who need can access it
// Mb call Repository ?
type Handler struct {
	Api       *config.APIConfig
	Logger    *zap.Logger
	Bot       *telegram.Bot
	Redis     *redis.Client
	DB        *gorm.DB
	TmplCache map[string]*template.Template
	Session   *sessions.CookieStore
}

//func (h *Handler) NewHandler() *Handler {
//	return &Handler{
//		Api:    api,
//		Logger: log,
//		Bot:    bot,
//		Redis:  rd,
//		DB:     db,
//	}
//}

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
