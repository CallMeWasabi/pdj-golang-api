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
	for {
		var optionData models.Option
		optionDoc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalln("Failed to iterate over option: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err := optionDoc.DataTo(&optionData); err != nil {
			log.Fatalln("Failed to convert data option: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		optionsData = append(optionsData, optionData)
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

func DeleteOption(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	id := c.Params("id")

	iter := client.Collection("options").Where("id", "==", id).Documents(ctx)
	defer iter.Stop()
	var optionRefId string
	var optionDoc models.Option
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalln("Failed to iterate menu: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err := doc.DataTo(&optionDoc); err != nil {
			log.Fatalln("Failed to convert data menu: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		optionRefId = doc.Ref.ID
	}

	for i := 0; i < len(optionDoc.MenusId); i++ {
		var menuData models.Menu
		menuDoc, err := client.Collection("menus").Doc(optionDoc.MenusId[i]).Get(ctx)
		if err != nil {
			log.Fatalln("Failed to get document menu: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if err := menuDoc.DataTo(&menuData); err != nil {
			log.Fatalln("Failed to conver data menu: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		for j := 0; j < len(menuData.OptionsId); j++ {
			if menuData.OptionsId[j] == optionRefId {
				removedMenusId := append(menuData.OptionsId[:j], menuData.OptionsId[j+1:]...)
				menuData.OptionsId = removedMenusId
				break
			}
		}

		_, err = client.Collection("menus").Doc(optionDoc.MenusId[i]).Set(ctx, menuData)
		if err != nil {
			log.Fatalln("Failed to update menus: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	_, err := client.Collection("options").Doc(optionRefId).Delete(ctx)
	if err != nil {
		log.Fatalln("Failed to delete option: ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"success": true,
	})

}

func UpdateOption(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	id := c.Params("id")

	var optionData models.Option
	if err := c.BodyParser(&optionData); err != nil {
		log.Fatalln("Failed to parse new option data: ", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	optionData.UpdatedAt = time.Now()

	var optionRefId string
	iter := client.Collection("options").Where("id", "==", id).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalln("Failed to iterate over option : ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		optionRefId = doc.Ref.ID
	}

	if _, err := client.Collection("options").Doc(optionRefId).Set(ctx, optionData); err != nil {
		log.Fatalln("An error has occurred: ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}
