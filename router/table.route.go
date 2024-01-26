package router

import (
	db "demo-go-firebase/firebase"
	"demo-go-firebase/models"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

func GetTable(c *fiber.Ctx) error {
	var tablesData []models.Tables

	ctx := db.Provider.Ctx
	client := db.Provider.Client

	iter := client.Collection("tables").Documents(ctx)
	defer iter.Stop()
	for {
		var buffer models.Tables
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalln("Failed to iterate over menu: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if err := doc.DataTo(&buffer); err != nil {
			log.Fatalln("Failed to convert data menu: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		tablesData = append(tablesData, buffer)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"result":  tablesData,
	})
}

func GetTableByID(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	id := c.Params("id")

	iter := client.Collection("tables").Where("id", "==", id).Documents(ctx)
	defer iter.Stop()
	buffer := models.Tables{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalln("Failed to iterate menu: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err := doc.DataTo(&buffer); err != nil {
			log.Fatalln("Failed to convert data option: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if id == buffer.ID {
			return c.JSON(fiber.Map{
				"success": true,
				"result":  buffer,
			})
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func CreateTable(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	uid := uuid.New()
	splitID := strings.Split(uid.String(), "-")
	id := splitID[0] + splitID[1] + splitID[2] + splitID[3] + splitID[4]

	uid = uuid.New()
	splitID = strings.Split(uid.String(), "-")
	tableUuid := splitID[0] + splitID[1] + splitID[2] + splitID[3] + splitID[4]

	newTable := models.Tables{ID: id, CreatedAt: time.Now(), UpdatedAt: time.Now(), AccessUuid: tableUuid}
	if err := c.BodyParser(&newTable); err != nil {
		log.Fatalln("Error parse new menu : ", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	_, _, err := client.Collection("tables").Add(ctx, newTable)
	if err != nil {
		log.Fatalln("Error create new menu : ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"result":  newTable,
	})
}

func UpdateTable(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	id := c.Params("id")

	var tableData models.Tables
	if err := c.BodyParser(&tableData); err != nil {
		log.Fatalln("Error parse new tableData: ", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	tableData.UpdatedAt = time.Now()

	var tableRefId string
	iter := client.Collection("tables").Where("id", "==", id).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalln("Failed to iterate over table : ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		tableRefId = doc.Ref.ID
	}

	if _, err := client.Collection("tables").Doc(tableRefId).Set(ctx, tableData); err != nil {
		log.Fatalln("An error has occurred: ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}
