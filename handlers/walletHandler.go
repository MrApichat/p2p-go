package handlers

import (
	"net/http"

	"github.com/MrApichat/p2p-go/db"
	"github.com/MrApichat/p2p-go/models"
	"github.com/MrApichat/p2p-go/utilities"
	"github.com/labstack/echo/v4"
)

func GetWallets(c echo.Context) error {
	wallets := []models.WalletModel{}
	cc, isLogin := utilities.IsLogin(c)
	if isLogin == false {
		return utilities.HandleError(cc, "Please Login", http.StatusUnauthorized)
	}

	rows, err := db.Db.Query(`
	SELECT w.id, w.coin_id, w.total, w.in_order , c.type, c.name 
	FROM wallets w 
	LEFT JOIN currencies c 
	ON w.coin_id =c.id WHERE w.user_id = $1;`, cc.User.Id)

	if err != nil {
		return utilities.HandleError(cc, "query:"+err.Error(), http.StatusInternalServerError)
	}

	defer rows.Close()
	for rows.Next() {
		var w models.WalletModel
		err := rows.Scan(&w.Id, &w.Coin.Id, &w.Total, &w.InOrder, &w.Coin.Type, &w.Coin.Name)
		if err != nil {
			return utilities.HandleError(cc, "scan:"+err.Error(), http.StatusInternalServerError)
		}
		wallets = append(wallets, w)
	}

	if rows.Err(); err != nil {
		return utilities.HandleError(c, "rows.Err:"+err.Error(), http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"list":    wallets,
	})

}
