package utilities

import (
	"strings"

	"github.com/MrApichat/p2p-go/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func HandleError(c echo.Context, message string, status int) error {
	return c.JSON(status,
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

func ValidationError(err error) string {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return err.Error()
	}

	val := []string{}
	for _, err := range err.(validator.ValidationErrors) {
		val = append(val, err.Field())
	}

	verb := " is required."
	if len(val) > 1 {
		verb = " are required."
	}
	return strings.Join(val[:], ", ") +verb
}
