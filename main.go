package main

import (
	"log"

	database "github.com/MrApichat/p2p-go/db"
	"github.com/MrApichat/p2p-go/route"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	database.Con()

	e := echo.New()
	e.Use(middleware.Logger())
	route.Router(e)

	port := "2565"
	log.Println("starting... port:", port)

	log.Fatal(e.Start(":" + port))

}
