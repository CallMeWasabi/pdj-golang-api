package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func InitializeRoutes(app *fiber.App) {

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Adjust this to be more restrictive if needed
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept, Token, New-Status, Includes",
	}))

	app.Get("/auth/:uuid", CreateToken)
	app.Get("/auth", RecheckStatus)

	app.Get("/menus", GetMenu)
	app.Get("/menus/:id", GetMenuByID)
	app.Post("/menus", CreateMenu)
	app.Put("/menus/:id", UpdateMenu)
	app.Delete("/menus/:id", DeleteMenu)

	app.Get("/menu-types", GetMenuType)
	app.Get("/menu-types/:id", GetMenuTypeByID)
	app.Post("/menu-types", CreateMenuType)
	app.Put("/menu-types/:id", UpdateMenuType)
	app.Delete("/menu-types/:id", DeleteMenuType)

	app.Get("/tables", GetTable)
	app.Get("/tables/:id", GetTableByID)
	app.Post("/tables", CreateTable)
	app.Put("/tables/:id", UpdateTable)
	app.Delete("/tables/:id", DeleteTable)

	app.Get("/options", GetOption)
	app.Get("/options/:id", GetOptionByID)
	app.Get("/options/ref/:id", GetOptionRefID)
	app.Post("/options", CreateOption)
	app.Put("/options/:id", UpdateOption)
	app.Delete("/options/:id", DeleteOption)

	app.Put("/orders/:table_id", UpdateAllOrder)
	app.Put("/orders/:table_id/:order_uuid", UpdateStatusOneOrder)
}
