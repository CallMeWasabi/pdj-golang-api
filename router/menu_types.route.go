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

func GetMenuType(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	includes := c.Get("Includes")

	iter := client.Collection("menu_types").Documents(ctx)
	var MenuTypesData []models.MenuType
	var buffer models.MenuType
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err := doc.DataTo(&buffer); err != nil {
			log.Fatalln("Failed to convert data: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if strings.ToLower(includes) == "true" {
			var MenusData []models.Menu
			var MenuData models.Menu
			for i := 0; i < len(buffer.MenusId); i++ {
				snapshot, err := client.Collection("menus").Doc(buffer.MenusId[i]).Get(ctx)
				if err != nil {
					log.Fatalln("Failed to get document: ", err)
					break
				}
				if err := snapshot.DataTo(&MenuData); err != nil {
					log.Fatalln("Failed to convert data: ", err)
					break
				}

				var OptionsData []models.Option
				var OptionData models.Option
				for i := 0; i < len(MenuData.OptionsId); i++ {
					snapshot, err := client.Collection("options").Doc(MenuData.OptionsId[i]).Get(ctx)
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

				MenuData.Options = OptionsData

				MenusData = append(MenusData, MenuData)
			}

			buffer.Menus = MenusData
		}
		MenuTypesData = append(MenuTypesData, buffer)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"result":  MenuTypesData,
	})
}

func GetMenuTypeByID(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	id := c.Params("id")

	iter := client.Collection("menu_types").Where("id", "==", id).Documents(ctx)
	buffer := models.MenuType{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalln("error iterate doc: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err := doc.DataTo(&buffer); err != nil {
			log.Fatalln("Failed to convert data: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		var MenusData []models.Menu
		var MenuData models.Menu
		for i := 0; i < len(buffer.MenusId); i++ {
			snapshot, err := client.Collection("menus").Doc(buffer.MenusId[i]).Get(ctx)
			if err != nil {
				log.Fatalln("Failed to get document: ", err)
				break
			}
			if err := snapshot.DataTo(&MenuData); err != nil {
				log.Fatalln("Failed to convert data: ", err)
				break
			}
			MenusData = append(MenusData, MenuData)
		}

		buffer.Menus = MenusData
	}

	return c.JSON(fiber.Map{
		"success": true,
		"result":  buffer,
	})
}

func CreateMenuType(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	uid := uuid.New()
	splitID := strings.Split(uid.String(), "-")
	id := splitID[0] + splitID[1] + splitID[2] + splitID[3] + splitID[4]

	var newMenuType = models.MenuType{ID: id, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := c.BodyParser(&newMenuType); err != nil {
		log.Fatalf("error parse new menu type : %v", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	_, _, err := client.Collection("menu_types").Add(ctx, newMenuType)
	if err != nil {
		log.Fatalf("error create new menu type : %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"result":  newMenuType,
	})
}

func UpdateMenuType(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	id := c.Params("id")

	var newMenuType = models.MenuType{ID: id, UpdatedAt: time.Now()}
	if err := c.BodyParser(&newMenuType); err != nil {
		log.Fatalf("Error parse new menu_type : %v", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	var refId string
	var currentMenuType models.MenuType
	iter := client.Collection("menu_types").Where("id", "==", id).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalln("Failed to iterate menu_type: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err := doc.DataTo(&currentMenuType); err != nil {
			log.Fatalln("Failed to convert data: ", err)
			break
		}

		refId = doc.Ref.ID
	}

	currentMenuType.Name = newMenuType.Name
	currentMenuType.Status = newMenuType.Status

	if _, err := client.Collection("menu_types").Doc(refId).Set(ctx, currentMenuType); err != nil {
		log.Fatalln("An error has occurred: ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func DeleteMenuType(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	id := c.Params("id")

	var menuType models.MenuType
	var refId string
	iter := client.Collection("menu_types").Where("id", "==", id).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalf("Failed to iterate menu_type: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err := doc.DataTo(&menuType); err != nil {
			log.Fatalln("An error has occurred: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		refId = doc.Ref.ID
	}

	if len(menuType.MenusId) > 0 {
		for i := 0; i < len(menuType.MenusId); i++ {
			menuRefId := menuType.MenusId[i]
			if _, err := client.Collection("menus").Doc(menuRefId).Delete(ctx); err != nil {
				log.Fatalf("An error has occurred: %s", err)
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}
	}

	if _, err := client.Collection("menu_types").Doc(refId).Delete(ctx); err != nil {
		log.Fatalf("An error has occurred: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}
