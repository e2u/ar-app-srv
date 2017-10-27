package middle

import (
	"net/http"

	"e2u.io/ar-app-srv/controllers"
)

// 初始化 controller
func InitController(c *controllers.Controller) func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if render := c.InitWithMiddle(w, r.WithContext(ctx)); render {
				return
			}
			h.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
	return f
}
