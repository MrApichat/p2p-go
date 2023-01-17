package utilities

import (
	"net/http"

	"github.com/MrApichat/p2p-go/models"
	"github.com/labstack/echo/v4"
)

func HandleError(c echo.Context, message string, status int) error {
	return c.JSON(http.StatusInternalServerError,
		map[string]interface{}{
			"message": message,
			"success": false})
}

func IsLogin(c echo.Context) (mo *models.UserContext, boo bool) {
	cc := c.(*models.UserContext)

	if cc.User.Id == 0 {
		return cc, false
	} else {
		return cc, true
	}
}
