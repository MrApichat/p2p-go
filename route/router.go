package route

import (
	"net/http"

	"github.com/MrApichat/p2p-go/handlers"
	"github.com/MrApichat/p2p-go/middlewares"
	"github.com/MrApichat/p2p-go/models"
	"github.com/labstack/echo/v4"
)

func Router(e *echo.Echo) {
	g := e.Group("/api", middlewares.Auth)
	g.GET("/hello", func(c echo.Context) error {
		cc := c.(*models.UserContext)
		return c.JSON(http.StatusOK, map[string]interface{}{"message": cc.User})
	})
	g.POST("/register", handlers.Register)
	g.POST("/login", handlers.Login)
	g.GET("/wallets", handlers.GetWallets)
	g.POST("/transfer_orders", handlers.CreateTransfer)
	g.GET("/transfer_orders", handlers.ShowTransfer)
	g.POST("/merchant_orders", handlers.CreateMerchant)
	g.GET("/merchant_orders", handlers.ShowMerchant)
}
