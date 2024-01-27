package router

import (
	db "demo-go-firebase/firebase"
	"demo-go-firebase/models"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/api/iterator"
)

func CreateToken(c *fiber.Ctx) error {
	ctx := db.Provider.Ctx
	client := db.Provider.Client
	uuid := c.Params("uuid")

	iter := client.Collection("tables").Where("access_uuid", "==", uuid).Documents(ctx)
	var tableData models.Tables
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalln("Failed to iterate over table: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err := doc.DataTo(&tableData); err != nil {
			log.Fatalln("Failed to convert table: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	if tableData.Status != "IN_SERVICE" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Header["alg"] = "HS256"
	claims := token.Claims.(jwt.MapClaims)
	claims["table_id"] = tableData.ID
	claims["table_name"] = tableData.Name
	claims["table_status"] = tableData.Status

	t, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"success":    true,
		"table_name": tableData.Name,
		"table_id":   tableData.ID,
		"token":      t,
	})
}

func RecheckStatus(c *fiber.Ctx) error {
	tokenString := c.Get("Token")
	if tokenString == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		log.Fatalln("Failed to parse jwt: ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	iter := db.Provider.Client.Collection("tables").Where("id", "==", claims["table_id"].(string)).Documents(db.Provider.Ctx)
	var tableData models.Tables
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatalln("Failed to iterate over table: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err := doc.DataTo(&tableData); err != nil {
			log.Fatalln("Failed to conver data table: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	if tableData.Status != "IN_SERVICE" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.SendStatus(fiber.StatusOK)
}
