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

func Register(c echo.Context) error {
	register := &models.RegisterModel{}

	//validate request
	if err := c.Bind(register); err != nil {
		return utilities.HandleError(c, err.Error(), http.StatusBadRequest)
	}

	//password = confirm_password
	if register.Password != register.ConfirmPassword {
		return utilities.HandleError(c, "Password mismatch", http.StatusBadRequest)
	}

	// hash password
	hPass, err := utilities.HashPassword(register.Password)
	if err != nil {
		return utilities.HandleError(c, err.Error(), http.StatusInternalServerError)
	}

	//generate token
	token, err := utilities.GenerateToken()
	if err != nil {
		return utilities.HandleError(c, err.Error(), http.StatusInternalServerError)
	}

	//database transaction
	tx, err := db.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return utilities.HandleError(c, err.Error(), http.StatusInternalServerError)
	}

	//insert user
	query, err := tx.Prepare(`INSERT INTO public.users 
	("name", email, "password", remember_token, created_at) 
	VALUES ($1, $2, $3, $4, now());`)

	if err != nil {
		tx.Rollback()
		return utilities.HandleError(c, "Prepare:"+err.Error(), http.StatusInternalServerError)
	}

	_, err = query.Exec(register.Name, register.Email, hPass, token)
	switch {
	case err == nil:
		//get user
		row := tx.QueryRow(`SELECT id,name, email, remember_token FROM users WHERE email=$1`, register.Email)
		user := models.UserModel{}
		err = row.Scan(&user.Id, &user.Name, &user.Email, &user.RememberToken)
		if err != nil {
			tx.Rollback()
			return utilities.HandleError(c, "QueryRow:"+err.Error(), http.StatusInternalServerError)
		}

		rows, err := tx.Query(`SELECT id, type, name FROM currencies WHERE type='coin'`)

		if err != nil {
			tx.Rollback()
			return utilities.HandleError(c, "Query:"+err.Error(), http.StatusInternalServerError)
		}

		//create user wallet
		currencies := []models.CurrencyModel{}
		for rows.Next() {
			var currency models.CurrencyModel
			err := rows.Scan(&currency.Id, &currency.Type, &currency.Name)
			if err != nil {
				tx.Rollback()
				return utilities.HandleError(c, "Scan:"+err.Error(), http.StatusInternalServerError)
			}
			currencies = append(currencies, currency)
		}

		if rows.Err(); err != nil {
			tx.Rollback()
			return utilities.HandleError(c, "rows.Err:"+err.Error(), http.StatusInternalServerError)
		}

		queryString := `INSERT INTO wallets 
		(user_id, coin_id, total, in_order, created_at) 
		VALUES `

		//doing multiple insert queries
		vals := []interface{}{}
		count := 1
		for _, cur := range currencies {
			value1 := "$" + strconv.Itoa(count)
			count++
			queryString += "(" + value1 + ", $" + strconv.Itoa(count) + ", 0, 0, now()),"
			vals = append(vals, user.Id, cur.Id)
			count++
		}

		queryString = queryString[0 : len(queryString)-1]

		stmt, err := tx.Prepare(queryString)
		if err != nil {
			tx.Rollback()
			return utilities.HandleError(c, "stmt:"+err.Error(), http.StatusInternalServerError)
		}

		_, err = stmt.Exec(vals...)
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

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"success": true,
			"data":    user,
		})
	default:
		return utilities.HandleError(c, "Exec:"+err.Error(), http.StatusInternalServerError)
	}
}

func Login(c echo.Context) error {
	login := &models.LoginModel{}
	user := models.UserModel{}
	var password string

	if err := c.Bind(login); err != nil {
		return c.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"message": err.Error(),
				"success": false})
	}

	//get user data
	row := db.Db.QueryRow(`SELECT id,name, email, remember_token, password FROM users WHERE email=$1`, login.Email)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.RememberToken, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return utilities.HandleError(c, "user not found", http.StatusNotFound)
		}
		return utilities.HandleError(c, "QueryRow:"+err.Error(), http.StatusInternalServerError)
	}

	//check password match
	if utilities.CheckPasswordHash(login.Password, password) {
		//generate token and update user to database
		token, _ := utilities.GenerateToken()
		if err != nil {
			return utilities.HandleError(c, err.Error(), http.StatusInternalServerError)
		}
		stmt, err := db.Db.Prepare(`UPDATE users SET remember_token=$1 WHERE id=$2`)
		if err != nil {
			return utilities.HandleError(c, err.Error(), http.StatusInternalServerError)
		}
		_, err = stmt.Exec(token, user.Id)
		if err != nil {
			return utilities.HandleError(c, err.Error(), http.StatusInternalServerError)
		}
		user.RememberToken = token
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"success": true,
			"data":    user,
		})
	} else {
		return utilities.HandleError(c, "password are not correct.", http.StatusUnauthorized)
	}
}
