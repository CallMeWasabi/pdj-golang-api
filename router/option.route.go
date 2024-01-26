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

func GetOption(c *fiber.Ctx) error {
	var optionsData []models.Option

	ctx := db.Provider.Ctx
	client := db.Provider.Client

	iter := client.Collection("options").Documents(ctx)
	defer iter.Stop()
	var buffer models.Option
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalln("Failed to iterate over option: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err := doc.DataTo(&buffer); err != nil {
			log.Fatalln("Failed to convert data option: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		optionsData = append(optionsData, buffer)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"result":  optionsData,
	})
}

func GetOptionByID(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	id := c.Params("id")

	iter := client.Collection("options").Where("id", "==", id).Documents(ctx)
	defer iter.Stop()
	var buffer models.Option
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalln("Failed to iterate over options: ", err)
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

	doc, err := client.Collection("options").Doc(id).Get(ctx)
	if err != nil {
		log.Fatalln("Failed to get document option: ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if id == doc.Ref.ID {
		if err := doc.DataTo(&buffer); err != nil {
			log.Fatalln("Failed to convert data option: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(fiber.Map{
			"success": true,
			"result":  buffer,
		})
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func GetOptionRefID(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	id := c.Params("id")

	iter := client.Collection("options").Where("id", "==", id).Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalln("Failed to iterate over options: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		} else {
			return c.JSON(fiber.Map{
				"success": true,
				"result":  doc.Ref.ID,
			})
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func CreateOption(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	uid := uuid.New()
	splitID := strings.Split(uid.String(), "-")
	id := splitID[0] + splitID[1] + splitID[2] + splitID[3] + splitID[4]

	newOption := models.Option{ID: id, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := c.BodyParser(&newOption); err != nil {
		log.Fatalln("Error parse new option: ", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	_, _, err := client.Collection("options").Add(ctx, newOption)
	if err != nil {
		log.Fatalln("Error create new option: ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"result":  newOption,
	})
}
