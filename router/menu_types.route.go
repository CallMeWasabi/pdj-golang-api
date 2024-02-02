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
	var menuTypesData []models.MenuType
	for {
		var menuTypeData models.MenuType
		menuTypeDoc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err := menuTypeDoc.DataTo(&menuTypeData); err != nil {
			log.Fatalln("Failed to convert data: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if strings.ToLower(includes) == "true" {
			var menusData []models.Menu
			for i := 0; i < len(menuTypeData.MenusId); i++ {
				var menuData models.Menu
				menuDoc, err := client.Collection("menus").Doc(menuTypeData.MenusId[i]).Get(ctx)
				if err != nil {
					log.Fatalln("Failed to get document: ", err)
					break
				}
				if err := menuDoc.DataTo(&menuData); err != nil {
					log.Fatalln("Failed to convert data: ", err)
					break
				}

				var optionsData []models.Option
				for i := 0; i < len(menuData.OptionsId); i++ {
					var optionData models.Option
					snapshot, err := client.Collection("options").Doc(menuData.OptionsId[i]).Get(ctx)
					if err != nil {
						log.Fatalln("Failed to get document: ", err)
						break
					}
					if err := snapshot.DataTo(&optionData); err != nil {
						log.Fatalln("Failed to convert data: ", err)
						break
					}
					optionsData = append(optionsData, optionData)
				}

				menuData.Options = optionsData

				menusData = append(menusData, menuData)
			}

			menuTypeData.Menus = menusData
		}
		menuTypesData = append(menuTypesData, menuTypeData)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"result":  menuTypesData,
	})
}

func GetMenuTypeByID(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	id := c.Params("id")

	iter := client.Collection("menu_types").Where("id", "==", id).Documents(ctx)
	menuTypeData := models.MenuType{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalln("error iterate doc: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err := doc.DataTo(&menuTypeData); err != nil {
			log.Fatalln("Failed to convert data: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		var menusData []models.Menu
		for i := 0; i < len(menuTypeData.MenusId); i++ {
			var menuData models.Menu
			menuDoc, err := client.Collection("menus").Doc(menuTypeData.MenusId[i]).Get(ctx)
			if err != nil {
				log.Fatalln("Failed to get document: ", err)
				break
			}
			if err := menuDoc.DataTo(&menuData); err != nil {
				log.Fatalln("Failed to convert data: ", err)
				break
			}
			menusData = append(menusData, menuData)
		}

		menuTypeData.Menus = menusData
	}

	return c.JSON(fiber.Map{
		"success": true,
		"result":  menuTypeData,
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

	var newMenuType models.MenuType
	if err := c.BodyParser(&newMenuType); err != nil {
		log.Fatalf("Error parse new menu_type : %v", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	newMenuType.UpdatedAt = time.Now()

	var refId string
	iter := client.Collection("menu_types").Where("id", "==", id).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalln("Failed to iterate menu_type: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		refId = doc.Ref.ID
	}

	if _, err := client.Collection("menu_types").Doc(refId).Set(ctx, newMenuType); err != nil {
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
