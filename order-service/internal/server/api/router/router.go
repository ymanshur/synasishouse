package router

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/ymanshur/synasishouse/order/internal/appctx"
	"github.com/ymanshur/synasishouse/order/internal/server/api/handler"
	"github.com/ymanshur/synasishouse/pkg/util"
)

// Router
type Router struct {
	config *appctx.Config
}

func NewRouter(
	config *appctx.Config,
) Router {
	return Router{
		config: config,
	}
}

// Route preparing Gin router and return HTTP handler
func (r *Router) Route() http.Handler {
	if util.TranslateEnv(r.config.Environment) == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Register field name tag func to use json tag as field names in errors
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" { // ignore fields with `json:"-"`
				return ""
			}
			return name
		})
	}

	router := gin.New()
	router.Use(gin.Recovery())

	healthHandler := handler.NewHealth()

	router.GET("/health", healthHandler.Check)

	return router
}
