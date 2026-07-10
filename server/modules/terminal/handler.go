package terminal

import (
	"net/http"

	"github.com/aipanel/aipanel/server/config"
	"github.com/aipanel/aipanel/server/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	cfg    *config.Config
	audit  *service.AuditService
	secret string
}

func NewHandler(cfg *config.Config, audit *service.AuditService) *Handler {
	return &Handler{cfg: cfg, audit: audit, secret: cfg.JWT.Secret}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) WS(c *gin.Context) {
	if !h.cfg.Terminal.Enabled {
		c.JSON(http.StatusForbidden, gin.H{"error": "terminal is disabled"})
		return
	}

	token := c.Query("token")
	claims, err := service.ParseToken(h.secret, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	c.Set("user_id", claims.UserID)
	c.Set("username", claims.Username)
	h.audit.Record(c, "terminal.login", h.cfg.Terminal.Shell)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	sess, err := newSession(h.cfg.Terminal.Shell)
	if err != nil {
		_ = conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}
	defer sess.Close()

	_ = conn.WriteMessage(websocket.TextMessage, []byte("Welcome to AIPanel Terminal\r\n"))

	done := make(chan struct{})
	go func() {
		defer close(done)
		buffer := make([]byte, 4096)
		for {
			n, err := sess.Read(buffer)
			if err != nil {
				return
			}
			if n > 0 {
				if err := conn.WriteMessage(websocket.BinaryMessage, buffer[:n]); err != nil {
					return
				}
			}
		}
	}()

	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if len(payload) > 0 {
			if _, err := sess.Write(payload); err != nil {
				break
			}
		}
	}

	_ = sess.Close()
	<-done
}
