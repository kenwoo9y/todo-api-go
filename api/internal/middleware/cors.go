package middleware

import (
	"net/http"

	"github.com/kenwoo9y/todo-api-go/api/internal/config"
)

type CORSConfig struct {
	Origins []string
}

func NewCORSConfig(cfg *config.Config) *CORSConfig {
	return &CORSConfig{
		Origins: cfg.CORSOrigins,
	}
}

func (c *CORSConfig) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			// 許可されたオリジンのリストに含まれているか、または*が設定されているかチェック
			allowed := false
			for _, allowedOrigin := range c.Origins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowed = true
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
			if !allowed {
				http.Error(w, "Not allowed origin", http.StatusForbidden)
				return
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// OPTIONSリクエストの場合は、ここで処理を終了
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
