package models

import "github.com/labstack/echo/v4"

type UserModel struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	RememberToken string `json:"rememberToken"`
}

type DisplayUserModel struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
}

type RegisterModel struct {
	Name            string `form:"name" json:"name"`
	Email           string `form:"email" json:"email"`
	Password        string `form:"password" json:"password"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password"`
}

type LoginModel struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

type UserContext struct {
	echo.Context
	User DisplayUserModel
}
