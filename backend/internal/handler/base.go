package handler

import (
	"Magaz/backend/internal/config"
	"Magaz/backend/internal/system/sse"
	ws "Magaz/backend/internal/system/websocket"
	"Magaz/backend/pkg/bot/telegram"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

// TODO: Needs to move to more global space so everyone who need can access it
// Mb call Repository (API)?
type Handler struct {
	Api       *config.APIConfig
	Logger    *zap.Logger
	Bot       *telegram.Bot
	Redis     *redis.Client
	DB        *gorm.DB
	TmplCache map[string]*template.Template
	Session   *sessions.CookieStore
	SSES      *sse.SSEHub //Server side event system
	WS        *ws.Manager //WebSockets system
}

func NewHandler(api *config.APIConfig) *Handler {
	return &Handler{
		Api:       api,
		Logger:    &zap.Logger{},
		Redis:     &redis.Client{},
		DB:        &gorm.DB{},
		Bot:       &telegram.Bot{},
		TmplCache: make(map[string]*template.Template),
		Session:   &sessions.CookieStore{},
		WS:        &ws.Manager{},
	}
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

func (h *Handler) ServeImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		img := c.Param("image")

		possibleExtensions := []string{".jpeg", ".jpg", ".jfif", ".pjpeg", "pjp", ".png", ".webp", ".svg"}

		cwd, _ := os.Getwd()
		var imagePath string
		fileFound := false
		for _, ext := range possibleExtensions {
			imagePath = filepath.Join(cwd, "backend", "storage", "images", img+ext)
			if _, err := os.Stat(imagePath); err == nil {
				fileFound = true
				break
			}
		}

		if !fileFound {
			h.Logger.Error("Image not found", zap.String("img", img))
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}

		c.File(imagePath)
	}
}

// Upgrade serve client request to upgrade regular connection into websocket
func (h *Handler) Upgrade() gin.HandlerFunc {
	return func(c *gin.Context) {

		websocketUpgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		conn, err := websocketUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			h.Logger.Error("websocket upgrade error", zap.Error(err))
		}

		client := ws.NewClient(conn, h.WS)
		h.WS.AddClient(client)

		//start read and write process
		go client.ReadMessage()
		go client.WriteMessage()

		//conn.SetCloseHandler(func(code int, text string) error {
		//	h.Logger.Info("websocket connection closed", zap.String("code", strconv.Itoa(code)), zap.String("text", text))
		//	return nil
		//})

	}
}
