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
	"github.com/ymanshur/synasishouse/order/internal/usecase"
	"github.com/ymanshur/synasishouse/pkg/util"
)

// Router
type Router struct {
	orderUseCase usecase.Orderer
}

func NewRouter(
	orderUseCase usecase.Orderer,
) Router {
	return Router{
		orderUseCase: orderUseCase,
	}
}

// Route preparing Gin router and return HTTP handler
func (r *Router) Route() http.Handler {
	config := appctx.LoadConfig()

	if util.TranslateEnv(config.Environment) == "production" {
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
	orderHandler := handler.NewOrder(r.orderUseCase)

	ApiRoutes := router.Group("/api")
	ApiRoutes.GET("/health", healthHandler.Check)
	ApiRoutes.POST("/orders", orderHandler.Create)

	return router
}
