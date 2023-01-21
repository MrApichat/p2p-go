package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/MrApichat/p2p-go/db"
	"github.com/MrApichat/p2p-go/models"
	"github.com/MrApichat/p2p-go/utilities"
	"github.com/labstack/echo/v4"
)

func CreateTransfer(c echo.Context) error {
	var order = &models.TransferOrderModel{}
	request := models.TransferRequest{}

	//validate
	if err := c.Bind(&request); err != nil {
		return utilities.HandleError(c, err.Error(), http.StatusBadRequest)
	}

	//check login
	cc, isLogin := utilities.IsLogin(c)
	if isLogin == false {
		return utilities.HandleError(cc, "Please Login", http.StatusUnauthorized)
	}

	if cc.User.Email == request.ReceiverEmail {
		return utilities.HandleError(cc, "You cannot send coins to yourself.", http.StatusBadRequest)
	}

	//get coin
	co := db.Db.QueryRow(`SELECT id, type, name FROM currencies WHERE type='coin' AND name=$1;`, request.Coin)

	err := co.Scan(&order.Coin.Id, &order.Coin.Type, &order.Coin.Name)
	if err != nil {
		return utilities.HandleError(c, "QueryCoin:"+err.Error(), http.StatusInternalServerError)
	}

	//get receiver
	rec := db.Db.QueryRow(`SELECT id,name, email FROM users WHERE email=$1`, request.ReceiverEmail)

	err = rec.Scan(&order.Receiver.Id, &order.Receiver.Name, &order.Receiver.Email)
	if err != nil {
		return utilities.HandleError(c, "QueryRec:"+err.Error(), http.StatusInternalServerError)
	}

	swallet := models.WalletModel{}
	rwallet := models.WalletModel{}

	//sender wallet
	send := db.Db.QueryRow(`SELECT w.id, w.total, w.in_order 
	FROM wallets w WHERE w.coin_id = $1 AND w.user_id = $2;`, order.Coin.Id, cc.User.Id)

	err = send.Scan(&swallet.Id, &swallet.Total, &swallet.InOrder)
	if err != nil {
		return utilities.HandleError(cc, "QuerySend:"+err.Error(), http.StatusInternalServerError)
	}

	if swallet.Total-swallet.InOrder < request.Amount {
		return utilities.HandleError(cc, "Your balance has not enough to send.", http.StatusBadRequest)
	}

	//receiver wallet
	recw := db.Db.QueryRow(`SELECT w.id, w.total, w.in_order 
	FROM wallets w WHERE w.coin_id = $1 AND w.user_id = $2;`, order.Coin.Id, order.Receiver.Id)

	err = recw.Scan(&rwallet.Id, &rwallet.Total, &rwallet.InOrder)
	if err != nil {
		return utilities.HandleError(cc, "QueryRW:"+err.Error(), http.StatusInternalServerError)
	}

	//create order as processing status
	query, err := db.Db.Prepare(`INSERT INTO transfer_orders
	(coin_id, sender_id, receiver_id, amount, status, created_at) 
	VALUES ($1, $2, $3, $4, 'processing',now()) RETURNING id;`)

	if err != nil {
		return utilities.HandleError(c, "Prepare:"+err.Error(), http.StatusInternalServerError)
	}

	err = query.QueryRow(order.Coin.Id, cc.User.Id, order.Receiver.Id, request.Amount).Scan(&order.Id)
	if err != nil {
		return utilities.HandleError(c, "Exec:"+err.Error(), http.StatusInternalServerError)
	}

	defer query.Close()

	order.Sender = models.UserModel{
		Id:    cc.User.Id,
		Name:  cc.User.Name,
		Email: cc.User.Email,
	}
	order.Amount = request.Amount
	order.Status = "processing"

	//database transaction
	tx, err := db.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return utilities.HandleError(c, err.Error(), http.StatusInternalServerError)
	}

	fail := func() {
		tx.Rollback()
		_, _ = db.Db.Exec(`UPDATE transfer_orders SET status=$1 WHERE id=$2`, "failed", order.Id)
	}

	//decease amount from sender wallet
	swallet.Total = swallet.Total - order.Amount
	_, err = tx.Exec(`UPDATE wallets SET total=$1 WHERE id=$2`, swallet.Total, swallet.Id)
	if err != nil {
		fail()
		return utilities.HandleError(c, err.Error(), http.StatusInternalServerError)
	}

	// try to send to receiver wallet
	rwallet.Total = rwallet.Total + order.Amount
	_, err = tx.Exec(`UPDATE wallets SET total=$1 WHERE id=$2`, rwallet.Total, rwallet.Id)
	if err != nil {
		fail()
		return utilities.HandleError(c, err.Error(), http.StatusInternalServerError)
	}

	//finish process
	order.Status = "success"
	_, err = db.Db.Exec(`UPDATE transfer_orders SET status=$1 WHERE id=$2`, order.Status, order.Id)
	if err != nil {
		fail()
		return utilities.HandleError(c, err.Error(), http.StatusInternalServerError)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return utilities.HandleError(c, "commit:"+err.Error(), http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    order,
		"rwal":    rwallet,
		"swal":    swallet,
	})
}

func ShowTransfer(c echo.Context) error {
	request := &models.TransferFilter{}

	if err := c.Bind(&request); err != nil {
		return utilities.HandleError(c, err.Error(), http.StatusBadRequest)
	}

	//check login
	cc, isLogin := utilities.IsLogin(c)
	if isLogin == false {
		return utilities.HandleError(cc, "Please Login", http.StatusUnauthorized)
	}

	q := `select to2.id, to2.amount, to2.status,
	to2.sender_id, u.name as sender_name, u.email as sender_email,
	to2.receiver_id, u2.name as receiver_name, u2.email as receiver_email,
	to2.coin_id , c."type" as coin_type, c."name"  as coin_name
	from transfer_orders to2 
	left join users u on to2.sender_id = u.id 
	left join users u2 on to2.receiver_id = u2.id 
	left join currencies c on c.id = to2.coin_id`

	i := 1
	vals := []interface{}{}
	switch request.Type {
	case "sender":
		q = q + " WHERE to2.sender_id=$1"
		i++
		vals = append(vals, cc.User.Id)

	case "receiver":
		q = q + " WHERE to2.receiver_id=$1"
		i++
		vals = append(vals, cc.User.Id)
	default:
		q = q + " WHERE to2.sender_id=$1 OR to2.receiver_id=$2"
		i = i + 2
		vals = append(vals, cc.User.Id, cc.User.Id)
	}

	if request.Coin != "" {
		q = q + "AND c.name=$" + strconv.Itoa(i)
		i++
		vals = append(vals, request.Coin)
	}

	if request.Status != "" {
		q = q + "AND to2.status=$" + strconv.Itoa(i)
		vals = append(vals, request.Status)
	}

	rows, err := db.Db.Query(q, vals...)

	if err != nil {
		return utilities.HandleError(c, "Query:"+err.Error(), http.StatusInternalServerError)
	}

	orders := []models.TransferOrderModel{}
	for rows.Next() {
		var order models.TransferOrderModel
		err := rows.Scan(&order.Id, &order.Amount, &order.Status,
			&order.Sender.Id, &order.Sender.Name, &order.Sender.Email,
			&order.Receiver.Id, &order.Receiver.Name, &order.Receiver.Email,
			&order.Coin.Id, &order.Coin.Type, &order.Coin.Name,
		)
		if err != nil {
			return utilities.HandleError(c, "Scan:"+err.Error(), http.StatusInternalServerError)
		}
		orders = append(orders, order)
	}

	if rows.Err(); err != nil {
		return utilities.HandleError(c, "rows.Err:"+err.Error(), http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"list":    orders,
	})
}
