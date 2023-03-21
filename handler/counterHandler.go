/**
To avoid data race or other simultaneous updates to the database,
use database transactions.
A transaction is a sequence of database operations that are treated as a single unit of work.
When using transactions, changes made to the database are isolated from other transactions until they are committed.
**/

package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/michaelcychan/basic-go-backend/database"
)

func AddToCounter(c *fiber.Ctx) error {
	// Get the name of the counter from the request parameters
	counterName := c.Params("name")

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}
	defer db.Close()

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Error starting transaction: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}

	// Prepare the update statement
	stmt, err := tx.Prepare("UPDATE counter SET counter = counter + 1 WHERE name = $1;")
	if err != nil {
		log.Fatalf("Error preparing statement: %v", err)
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}

	// Execute the update statement
	_, err = stmt.Exec(counterName)
	if err != nil {
		log.Fatalf("Error executing statement: %v", err)
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Error committing transaction: %v", err)
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"message": "server error"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "counter updated successfully"})
}
