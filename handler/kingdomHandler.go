package handler

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/michaelcychan/basic-go-backend/database"
	"github.com/michaelcychan/basic-go-backend/model"
)

func KingdomRootHandler(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"message": "Welcome to the Kingdom"})
}

func FindMonarchHandler(c *fiber.Ctx) error {
	db, errConnect := database.Connect()
	if errConnect != nil {
		log.Fatalln(errConnect)
	}

	unescapedMonarch, errUnescape := url.PathUnescape(c.Params("monarch"))
	if errUnescape != nil {
		log.Fatalf("unescape error: %s", errUnescape)
	}

	query := fmt.Sprintf("SELECT * FROM monarch WHERE name = '%s';", unescapedMonarch)
	queryResult, err := db.Query(query)

	if err != nil {
		log.Fatalf("Query error: %s", err)
		return c.Status(500).JSON(&fiber.Map{"message": "server error"})
	}

	defer queryResult.Close()

	var mon model.FullMonarchJson

	for queryResult.Next() {
		queryResult.Scan(&mon.Name, &mon.YearOfBirth, &mon.YearOfDeath, &mon.ReignStart, &mon.ReignEnd)
	}

	var empty = model.FullMonarchJson{}
	if mon == empty {
		return c.Status(404).JSON(&fiber.Map{"message": "not found"})
	}
	return c.Status(200).JSON(mon)
}

func FindAllMonarchHandler(c *fiber.Ctx) error {

	db, errConnect := database.Connect()
	if errConnect != nil {
		log.Fatalln(errConnect)
	}
	queryResult, err := db.Query("SELECT * FROM monarch;")

	if err != nil {
		log.Fatalf("Query error: %s", err)
		c.Status(500).JSON(&fiber.Map{"message": "server error"})
	}

	defer queryResult.Close()

	var result []model.FullMonarchJson

	for queryResult.Next() {
		var mon model.FullMonarchJson
		queryResult.Scan(&mon.Name, &mon.YearOfBirth, &mon.YearOfDeath, &mon.ReignStart, &mon.ReignEnd)
		result = append(result, mon)
	}
	return c.Status(200).JSON(result)
}

func FindLongestReignMonarch(c *fiber.Ctx) error {
	db, err := database.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// using current year for reigning monarch when year_of_reign_end is null.
	query := "SELECT name, " +
		"CASE WHEN year_of_reign_end IS NULL THEN EXTRACT(YEAR FROM NOW()) - year_of_reign_start " +
		"ELSE year_of_reign_end - year_of_reign_start END AS total_reign " +
		"FROM monarch " +
		"ORDER BY total_reign DESC " +
		"LIMIT 1;"

	var name string
	var totalReign int
	err = db.QueryRow(query).Scan(&name, &totalReign)

	if err != nil {
		log.Fatalf("Query error: %s", err)
		return c.Status(500).JSON(&fiber.Map{"message": "server error"})
	}

	result := fiber.Map{"name": name, "total_reign": totalReign}
	return c.Status(200).JSON(result)
}

func AddMonarch(c *fiber.Ctx) error {
	var mon model.FullMonarchJson
	if err := c.BodyParser(&mon); err != nil {
		log.Fatalf("Error parsing request body: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"message": "bad request"})
	}
	db, errConnect := database.Connect()
	if errConnect != nil {
		log.Fatalln(errConnect)
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Error starting transaction: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}

	stmt, err := tx.Prepare("INSERT INTO monarch (name, year_of_birth, year_of_death, year_of_reign_start, year_of_reign_end) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		log.Fatalf("Error preparing transaction: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}

	_, err = stmt.Exec(mon.Name, mon.YearOfBirth, mon.YearOfDeath, mon.ReignStart, mon.ReignEnd)
	if err != nil {
		log.Fatalf("Error executing transaction: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Error committing transaction: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}

	return c.Status(200).JSON(&fiber.Map{"message": "inserted"})
}

func UpdateMonarch(c *fiber.Ctx) error {
	var mon model.FullMonarchJson
	if err := c.BodyParser(&mon); err != nil {
		log.Fatalf("Error parsing request body: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"message": "bad request"})
	}
	db, errConnect := database.Connect()
	if errConnect != nil {
		log.Fatalln(errConnect)
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Error starting transaction: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}

	stmt, err := tx.Prepare("UPDATE monarch SET year_of_birth=$2, year_of_death=$3, year_of_reign_start=$4, year_of_reign_end=$5 WHERE name=$1")
	if err != nil {
		log.Fatalf("Error preparing transaction: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}

	_, err = stmt.Exec(mon.Name, mon.YearOfBirth, mon.YearOfDeath, mon.ReignStart, mon.ReignEnd)
	if err != nil {
		log.Fatalf("Error executing transaction: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Error committing transaction: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}

	return c.Status(200).JSON(&fiber.Map{"message": "updated"})
}

func DeleteMonarch(c *fiber.Ctx) error {
	type deleteMonarch struct {
		Name string
	}
	var monarchName deleteMonarch
	if err := c.BodyParser(&monarchName); err != nil {
		log.Fatalf("Error parsing request body: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"message": "bad request"})
	}
	name := monarchName.Name
	db, errConnect := database.Connect()
	if errConnect != nil {
		log.Fatalln(errConnect)
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Error starting transaction: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}

	stmt, err := tx.Prepare("DELETE FROM monarch WHERE name=$1")
	if err != nil {
		log.Fatalf("Error preparing transaction: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}

	_, err = stmt.Exec(name)
	if err != nil {
		log.Fatalf("Error executing transaction: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Error committing transaction: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}

	return c.Status(200).JSON(&fiber.Map{"message": "deleted"})
}
