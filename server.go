package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Room struct {
	Name string
	Temp int
}

type FullMonarchJson struct {
	Name        string `json: name`
	YearOfBirth int    `json: birth_year`
	YearOfDeath int    `json: death_year`
	ReignStart  int    `json: reign_start`
	ReignEnd    int    `json: reign_end`
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

func findMonarchHandler(c *fiber.Ctx, db *sql.DB, monarch string) []FullMonarchJson {
	query := fmt.Sprintf("SELECT * FROM monarch WHERE name = '%s';", monarch)
	queryResult, err := db.Query(query)

	if err != nil {
		log.Fatalf("Query error: %s", err)
		return []FullMonarchJson{}
	}

	defer queryResult.Close()

	var result []FullMonarchJson
	for queryResult.Next() {
		var mon FullMonarchJson
		queryResult.Scan(&mon.Name, &mon.YearOfBirth, &mon.YearOfDeath, &mon.ReignStart, &mon.ReignEnd)
		result = append(result, mon)
	}
	return result
}

func main() {

	dbUser := goDotEnvVariable("DBUSER")
	dbPassword := goDotEnvVariable("DBPASSWORD")
	dbURL := goDotEnvVariable("DBURL")
	dbPort := goDotEnvVariable("DBPORT")
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/kingdom?sslmode=disable", dbUser, dbPassword, dbURL, dbPort)

	db, dbErr := sql.Open("postgres", connStr)
	if dbErr != nil {
		log.Fatalf("Database error: %s", dbErr)
	}

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
		return c.JSON("Welcome to the Kingdom")
	})

	kingdomPath.Get("/get-monarch/:monarch", func(c *fiber.Ctx) error {
		unescapedParam, errUnescape := url.PathUnescape(c.Params("monarch"))
		if errUnescape != nil {
			log.Fatalf("unescape error: %s", errUnescape)
		}
		data := findMonarchHandler(c, db, unescapedParam)
		fmt.Printf("outside findMonarchHandler %v\n", data)
		return c.JSON(data)
	})

	app.Listen(":" + serverPort)
}
