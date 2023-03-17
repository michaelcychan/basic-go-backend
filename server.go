package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Room struct {
	Name string
	Temp int
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	return os.Getenv(key)
}

func indexHandler(c *fiber.Ctx, rooms []Room, roomName string) string {
	for _, room := range rooms {
		if roomName == room.Name {
			return fmt.Sprintf("%d", room.Temp)
		}
	}
	return "room not found"
}

func main() {
	app := fiber.New()

	serverPort := goDotEnvVariable("PORT")
	if serverPort == "" {
		serverPort = "3001"
	}

	homePath := app.Group("/home")
	kingdomPath := app.Group("/kingdom")

	home := []Room{
		{
			Name: "Master Bedroom",
			Temp: 20,
		}, {
			Name: "Medium Bedroom",
			Temp: 20,
		}, {
			Name: "Small Bedroom",
			Temp: 17,
		}, {
			Name: "Office",
			Temp: 18,
		}, {
			Name: "Reception",
			Temp: 19,
		}, {
			Name: "Bathroom",
			Temp: 20,
		},
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	homePath.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome Home!")
	})

	homePath.Get("/get-temp/:roomname", func(c *fiber.Ctx) error {
		return c.SendString(indexHandler(c, home, c.Params("roomname")))
	})

	kingdomPath.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Kingdom")
	})

	app.Listen(":" + serverPort)
}
