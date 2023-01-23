package models

import "github.com/labstack/echo/v4"

type UserModel struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	RememberToken string `json:"rememberToken"`
}

type DisplayUserModel struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RegisterModel struct {
	Name            string `form:"name" json:"name" validate:"required"`
	Email           string `form:"email" json:"email" validate:"required"`
	Password        string `form:"password" json:"password" validate:"required"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" validate:"required"`
}

type LoginModel struct {
	Email    string `form:"email" json:"email" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
}

type UserContext struct {
	echo.Context
	User DisplayUserModel
}
