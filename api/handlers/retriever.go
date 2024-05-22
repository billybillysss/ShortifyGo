package handlers

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/billybillysss/ShortifyGo/database"
	"github.com/gofiber/fiber/v2"
)

// RetrieveUrl is a handler function for GET requests to the API.
// It is responsible for retrieving the URL associated with a given
// short ID and redirecting the client to that URL.
//
// The handler function first retrieves the URL associated with the
// given short ID from the Redis database. If the short ID is not
// found in the database, the handler returns a 404 error.
//
// Next, the handler updates the TTL for the short ID in the Redis
// database. If the TTL renewal fails, the handler returns an error.
//
// Finally, the handler redirects the client to the retrieved URL
// with a 301 status code.
func RetrieveUrl(c *fiber.Ctx) error {
	shortID := c.Params("short_id")
	expiration, _ := strconv.Atoi(os.Getenv("EXPIRATION"))

	err := database.RdsValid.Ping(database.Ctx).Err()
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":      "error",
			"description": fmt.Sprintf("%v", err),
		})
	}

	// Get the URL associated with the short ID from the Redis database
	url, _ := database.RdsUrl.Get(database.Ctx, shortID).Result()
	if url == "" {
		// If the short ID is not found in the database, return a 404 error
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":      "error",
			"description": "URL not found",
		})
	}

	// Update the TTL for the short ID in the Redis database
	err = database.RdsUrl.Expire(database.Ctx, shortID, time.Duration(expiration)*time.Minute).Err()

	if err != nil {
		// If the TTL renewal fails, return an error
		log.Printf("Failed to renew TTL for short ID %v: %v", shortID, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":      "error",
			"description": "failed to renew TTL",
		})
	}
	// Redirect the client to the retrieved URL
	log.Printf("Redirecting to URL %v", url)
	return c.Redirect(url, 301)
}
