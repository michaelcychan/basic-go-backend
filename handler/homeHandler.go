package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/michaelcychan/basic-go-backend/model"
)

func IndexHandler(c *fiber.Ctx) error {
	home := []model.Room{
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
	roomName := c.Params("roomname")
	for _, room := range home {
		if roomName == room.Name {
			return c.SendString(fmt.Sprintf("%d", room.Temp))
		}
	}
	return c.SendString("room not found")
}
