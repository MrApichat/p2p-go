package middlewares

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/MrApichat/p2p-go/db"
	"github.com/MrApichat/p2p-go/models"
	"github.com/MrApichat/p2p-go/utilities"
	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var header []string
		header = c.Request().Header["Authorization"]
		user := models.UserModel{}
		if header != nil {
			author := c.Request().Header["Authorization"][0]
			authToken := strings.Split(author, "Bearer ")[1]

			//get user data
			row := db.Db.QueryRow(`SELECT id,name, email, remember_token FROM users WHERE remember_token=$1`, authToken)
			err := row.Scan(&user.Id, &user.Name, &user.Email, &user.RememberToken)
			if err != nil {
				if err == sql.ErrNoRows {
					if err == sql.ErrNoRows {
						return utilities.HandleError(c, "user not found", http.StatusNotFound)
					}
				}
				return utilities.HandleError(c, "mwAuth:"+err.Error(), http.StatusInternalServerError)
			}
		}
		cc := &models.UserContext{
			Context: c,
			User:    user,
		}
		return next(cc)
	}
}
