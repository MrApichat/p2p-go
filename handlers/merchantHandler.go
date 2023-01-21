package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/MrApichat/p2p-go/db"
	"github.com/MrApichat/p2p-go/models"
	"github.com/MrApichat/p2p-go/utilities"
	"github.com/labstack/echo/v4"
)

func CreateMerchant(c echo.Context) error {
	var request = &models.MerchantRequest{}

	//validate
	if err := c.Bind(&request); err != nil {
		return utilities.HandleError(c, err.Error(), http.StatusBadRequest)
	}

	//check login
	cc, isLogin := utilities.IsLogin(c)
	if isLogin == false {
		return utilities.HandleError(cc, "Please Login", http.StatusUnauthorized)
	}

	tx, err := db.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return utilities.HandleError(c, err.Error(), http.StatusInternalServerError)
	}

	rows, err := tx.Query(`select id, name, type from currencies c where name=$1 or name=$2`, request.Coin, request.Fiat)
	if err != nil {
		tx.Rollback()
		return utilities.HandleError(c, "Query1:"+err.Error(), http.StatusInternalServerError)
	}

	coin := models.CurrencyModel{}
	fiat := models.CurrencyModel{}
	for rows.Next() {
		var currency models.CurrencyModel
		err := rows.Scan(&currency.Id, &currency.Name, &currency.Type)
		if err != nil {
			tx.Rollback()
			return utilities.HandleError(c, "Scan:"+err.Error(), http.StatusInternalServerError)
		}
		if currency.Type == "coin" {
			//find coin
			coin = currency
		} else {
			//find fiat
			fiat = currency
		}
	}

	//check available coin is less than total when is buy type
	if request.Type == "buy" {
		wallet := models.WalletModel{}
		row := tx.QueryRow(`SELECT id, total, in_order FROM wallets WHERE user_id=$1 AND coin_id=$2`, cc.User.Id, coin.Id)
		err = row.Scan(&wallet.Id, &wallet.Total, &wallet.InOrder)
		if err != nil {
			tx.Rollback()
			return utilities.HandleError(c, "QueryRow:"+err.Error(), http.StatusInternalServerError)
		} else if wallet.Total-wallet.InOrder < request.AvailableCoin {
			return utilities.HandleError(c, "Your balance has not enough to create order", http.StatusBadRequest)
		}

		//change in_order wallet for buy type
		wallet.InOrder = wallet.InOrder + request.AvailableCoin
		stmt, err := tx.Prepare(`UPDATE wallets SET in_order=$1 WHERE id=$2`)
		if err != nil {
			return utilities.HandleError(c, err.Error(), http.StatusInternalServerError)
		}
		_, err = stmt.Exec(wallet.InOrder, wallet.Id)
		if err != nil {
			tx.Rollback()
			return utilities.HandleError(c, err.Error(), http.StatusInternalServerError)
		}

	}
	//lower limit is possible
	if request.AvailableCoin*request.Price < request.LowerLimit {
		tx.Rollback()
		return utilities.HandleError(c, "Your lower limit is impossible", http.StatusBadRequest)
	}

	//fiat coin must unique in merchant order that status open
	duplicate := models.MerchantOrderModel{}
	err = tx.QueryRow(`
	SELECT id 
	FROM merchant_orders 
	WHERE fiat_id=$1 and coin_id=$2 and type=$3 and merchant_id=$4`,
		fiat.Id, coin.Id, request.Type, cc.User.Id).Scan(&duplicate.Id)

	if err != sql.ErrNoRows {
		tx.Rollback()
		return utilities.HandleError(c, "You already have open order for this coin and fiat", http.StatusBadRequest)
	} else if err == sql.ErrNoRows {
	} else if err != nil {
		tx.Rollback()
		return utilities.HandleError(c, "QueryRow:"+err.Error(), http.StatusInternalServerError)
	}
	//create order
	order := models.MerchantOrderModel{}
	order.AvailableCoin = request.AvailableCoin
	order.Coin = coin
	order.Fiat = fiat
	order.Type = request.Type
	order.Merchant = cc.User
	order.Price = request.Price
	order.LowerLimit = float64(request.LowerLimit)
	order.Status = "start"

	q, err := tx.Prepare(`INSERT INTO merchant_orders 
	("type", fiat_id, coin_id, merchant_id, price, available_coin, lower_limit, status, created_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now()) RETURNING id`)
	if err != nil {
		tx.Rollback()
		return utilities.HandleError(c, "Prepare:"+err.Error(), http.StatusInternalServerError)
	}

	err = q.QueryRow(order.Type, order.Fiat.Id,
		order.Coin.Id, order.Merchant.Id, order.Price,
		order.AvailableCoin, order.LowerLimit, order.Status).Scan(&order.Id)
	if err != nil {
		tx.Rollback()
		return utilities.HandleError(c, "QueryRow:"+err.Error(), http.StatusInternalServerError)
	}

	//add payment method to order
	qString := `SELECT id, name FROM payment_methods WHERE `
	i := 1
	vals := []interface{}{}
	for _, v := range request.PaymentMethods {
		qString = qString + `name=$` + strconv.Itoa(i) + ` OR `
		i++
		vals = append(vals, v)
	}

	qString = qString[0 : len(qString)-3]

	payQuery, err := tx.Query(qString, vals...)
	if err != nil {
		tx.Rollback()
		return utilities.HandleError(c, "Query2:"+err.Error(), http.StatusInternalServerError)
	}

	for payQuery.Next() {
		var payment models.PaymentMethodModel
		err = payQuery.Scan(&payment.Id, &payment.Name)
		if err != nil {
			tx.Rollback()
			return utilities.HandleError(c, "Scan:"+err.Error(), http.StatusInternalServerError)
		}

		order.PaymentMethods = append(order.PaymentMethods, payment)
	}

	qString = `INSERT INTO merchant_orders_payment_methods (payment_method_id, merchant_order_id) VALUES `
	plist := []interface{}{}
	i = 1
	for _, v := range order.PaymentMethods {
		qString = qString + `($` + strconv.Itoa(i) + ", " + strconv.Itoa(int(order.Id)) + ") ,"
		i++
		plist = append(plist, v.Id)
	}

	qString = qString[0 : len(qString)-1]

	stmt, err := tx.Prepare(qString)
	if err != nil {
		tx.Rollback()
		return utilities.HandleError(c, "stmt:"+err.Error(), http.StatusInternalServerError)
	}

	_, err = stmt.Exec(plist...)
	if err != nil {
		tx.Rollback()
		return utilities.HandleError(c, "stmtExec:"+err.Error(), http.StatusInternalServerError)
	}

	//commit database
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return utilities.HandleError(c, "commit:"+err.Error(), http.StatusInternalServerError)
	}

	//response
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    order,
	})
}
