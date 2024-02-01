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

func GetMenu(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	includes := c.Get("Includes")

	iter := client.Collection("menus").Documents(ctx)
	defer iter.Stop()
	var MenusData []models.Menu
	var buffer models.Menu
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalln("Failed to iterate over menu: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err := doc.DataTo(&buffer); err != nil {
			log.Fatalln("Failed to convert menu: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if strings.ToLower(includes) == "true" {
			var OptionsData []models.Option
			var OptionData models.Option
			for i := 0; i < len(buffer.OptionsId); i++ {
				snapshot, err := client.Collection("options").Doc(buffer.OptionsId[i]).Get(ctx)
				if err != nil {
					log.Fatalln("Failed to get document: ", err)
					break
				}
				if err := snapshot.DataTo(&OptionData); err != nil {
					log.Fatalln("Failed to convert data: ", err)
					break
				}
				OptionsData = append(OptionsData, OptionData)
			}

			buffer.Options = OptionsData
		}

		MenusData = append(MenusData, buffer)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"result":  MenusData,
	})
}

func GetMenuByID(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	id := c.Params("id")
	includes := c.Get("Includes")

	iter := client.Collection("menus").Where("id", "==", id).Documents(ctx)
	menuData := models.Menu{}
	for {
		menuDoc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalln("Failed to iterate menu: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err := menuDoc.DataTo(&menuData); err != nil {
			log.Fatalln("Failed to convert menu: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if strings.ToLower(includes) == "true" {
			var OptionsData []models.Option
			for i := 0; i < len(menuData.OptionsId); i++ {
				var OptionData models.Option
				optionDoc, err := client.Collection("options").Doc(menuData.OptionsId[i]).Get(ctx)
				if err != nil {
					log.Fatalln("Failed to get document: ", err)
					break
				}
				if err := optionDoc.DataTo(&OptionData); err != nil {
					log.Fatalln("Failed to convert data: ", err)
					break
				}
				OptionsData = append(OptionsData, OptionData)
			}

			menuData.Options = OptionsData
		}

		if id == menuData.ID {
			return c.JSON(fiber.Map{
				"success": true,
				"result":  menuData,
			})
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func CreateMenu(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	uid := uuid.New()
	splitID := strings.Split(uid.String(), "-")
	id := splitID[0] + splitID[1] + splitID[2] + splitID[3] + splitID[4]

	newMenu := models.Menu{ID: id, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := c.BodyParser(&newMenu); err != nil {
		log.Fatalln("error parse new menu : ", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Get menu type ref id
	iter := client.Collection("menu_types").Where("id", "==", newMenu.MenuTypeId).Documents(ctx)
	var menuTypeDoc models.MenuType
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalln("Failed to iterate over menu-types: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if err := doc.DataTo(&menuTypeDoc); err != nil {
			log.Fatalln("Failed to convert data menu-types: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		newMenu.MenuTypeId = doc.Ref.ID
	}

	docRef, _, err := client.Collection("menus").Add(ctx, newMenu)
	if err != nil {
		log.Fatalln("error create new menu : ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	menuTypeDoc.MenusId = append(menuTypeDoc.MenusId, docRef.ID)

	_, err = client.Collection("menu_types").Doc(newMenu.MenuTypeId).Set(ctx, menuTypeDoc)
	if err != nil {
		log.Fatalln("An error has occurred: ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	for i := 0; i < len(newMenu.OptionsId); i++ {
		optionDocRef, err := client.Collection("options").Doc(newMenu.OptionsId[i]).Get(ctx)
		if err != nil {
			log.Fatalln("Failed to get doc option: ", err)
		}

		var optionData models.Option
		if err := optionDocRef.DataTo(&optionData); err != nil {
			log.Fatalln("Failed to convert data option: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		optionData.MenusId = append(optionData.MenusId, docRef.ID)
		if _, err = client.Collection("options").Doc(optionDocRef.Ref.ID).Set(ctx, optionData); err != nil {
			log.Fatalln("Failed to update option: ", err)
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"result":  newMenu,
	})
}

func UpdateMenu(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	id := c.Params("id")

	menuData := models.Menu{ID: id, UpdatedAt: time.Now()}
	if err := c.BodyParser(&menuData); err != nil {
		log.Fatalln("Error parse new menu : ", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var refId string
	iter := client.Collection("menus").Where("id", "==", id).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalln("Failed to iterate : ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		refId = doc.Ref.ID
	}

	if _, err := client.Collection("menus").Doc(refId).Set(ctx, menuData); err != nil {
		log.Fatalln("An error has occurred: ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func DeleteMenu(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	id := c.Params("id")

	iter := client.Collection("menus").Where("id", "==", id).Documents(ctx)
	defer iter.Stop()
	var menuRefId string
	var menuDoc models.Menu
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalln("Failed to iterate menu: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err := doc.DataTo(&menuDoc); err != nil {
			log.Fatalln("Failed to convert data menu: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		menuRefId = doc.Ref.ID
	}

	var menuTypeDoc models.MenuType
	refMenuType, err := client.Collection("menu_types").Doc(menuDoc.MenuTypeId).Get(ctx)
	if err != nil {
		log.Fatalln("Failed to get document menu_types: ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if err := refMenuType.DataTo(&menuTypeDoc); err != nil {
		log.Fatalln("Failed to convert data menu_types: ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	for i := 0; i < len(menuTypeDoc.MenusId); i++ {
		if menuTypeDoc.MenusId[i] == menuRefId {
			removedMenusId := append(menuTypeDoc.MenusId[:i], menuTypeDoc.MenusId[i+1:]...)
			menuTypeDoc.MenusId = removedMenusId
			break
		}
	}

	_, err = client.Collection("menu_types").Doc(menuDoc.MenuTypeId).Set(ctx, menuTypeDoc)
	if err != nil {
		log.Fatalln("Failed to update menu_type: ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	for i := 0; i < len(menuDoc.OptionsId); i++ {
		var optionDoc models.Option
		refOption, err := client.Collection("options").Doc(menuDoc.OptionsId[i]).Get(ctx)
		if err != nil {
			log.Fatalln("Failed to get document option: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if err := refOption.DataTo(&optionDoc); err != nil {
			log.Fatalln("Failed to convert data option: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		for j := 0; j < len(optionDoc.MenusId); j++ {
			if optionDoc.MenusId[j] == menuRefId {
				removedOptionsId := append(optionDoc.MenusId[:j], optionDoc.MenusId[j+1:]...)
				optionDoc.MenusId = removedOptionsId
				break
			}
		}

		_, err = client.Collection("options").Doc(menuDoc.OptionsId[i]).Set(ctx, optionDoc)
		if err != nil {
			log.Fatalln("Failed to update options: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	_, err = client.Collection("menus").Doc(menuRefId).Delete(ctx)
	if err != nil {
		log.Fatalln("Failed to delete menu: ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}
