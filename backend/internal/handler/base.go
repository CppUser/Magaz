package handler

import (
	"Magaz/backend/internal/config"
	"Magaz/backend/internal/handler/templates"
	"Magaz/backend/internal/system/sse"
	ws "Magaz/backend/internal/system/websocket"
	"Magaz/backend/pkg/bot/telegram"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	TmplCache *templates.TemplateCache
	Session   *sessions.CookieStore
	SSES      *sse.SSEHub //Server side event system
	WS        *ws.Manager //WebSockets system
}

func NewHandler(api *config.APIConfig) *Handler {
	handler := &Handler{
		Api:    api,
		Logger: &zap.Logger{},
		Redis:  &redis.Client{},
		DB:     &gorm.DB{},
		Bot:    &telegram.Bot{},

		Session: &sessions.CookieStore{},
		WS:      &ws.Manager{},
	}

	handler.TmplCache, _ = templates.NewTemplateCache(api.Tmpl.Layouts, api.Tmpl.Pages, api.Tmpl.Components)

	return handler
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
			CheckOrigin:     checkOrigin,
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

// TODO: IMPORTANT CORS implement in env to allow who can connect
func checkOrigin(r *http.Request) bool { //TODO: Adjust in env variable to allow origin connections
	origin := r.Header.Get("Origin")

	switch origin {
	case "https://9e4f-73-192-67-43.ngrok-free.app": //TODO:retrieve from env
		return true
	default:
		return false
	}
}
