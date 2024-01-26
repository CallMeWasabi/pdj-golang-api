package router

import (
	db "demo-go-firebase/firebase"
	"demo-go-firebase/models"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/iterator"
)

func UpdateAllOrder(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	tableId := c.Params("table_id")
	var orders models.OrderQuery

	if err := c.BodyParser(&orders); err != nil {
		log.Fatalln("Failed parse new order: ", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var tableRefId string
	var tableData models.Tables
	iter := client.Collection("tables").Where("id", "==", tableId).Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalln("Failed to iterate over table: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err := doc.DataTo(&tableData); err != nil {
			log.Fatalln("Failed to convert data table: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		tableRefId = doc.Ref.ID
	}

	tableData.Orders = orders.Orders
	tableData.UpdatedAt = time.Now()

	if _, err := client.Collection("tables").Doc(tableRefId).Set(ctx, tableData); err != nil {
		log.Fatalln("Failed to save data table: ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func UpdateStatusOneOrder(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	tableId := c.Params("table_id")
	orderUuid := c.Params("order_uuid")
	newStatus := c.Get("New-Status")

	var tableRefId string
	var tableData models.Tables
	iter := client.Collection("tables").Where("id", "==", tableId).Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalln("Failed to iterate over table: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err := doc.DataTo(&tableData); err != nil {
			log.Fatalln("Failed to convert data table: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		tableRefId = doc.Ref.ID
	}

	for i := 0; i < len(tableData.Orders); i++ {
		if tableData.Orders[i].Uuid == orderUuid {
			tableData.Orders[i].Status = newStatus
		}
	}

	if _, err := client.Collection("tables").Doc(tableRefId).Set(ctx, tableData); err != nil {
		log.Fatalln("Failed to save data table: ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}
