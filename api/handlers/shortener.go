package handlers

import (
	"encoding/json"
	"log"
	"time"

	"fmt"

	"os"
	"strconv"

	"github.com/billybillysss/ShortifyGo/database"
	"github.com/billybillysss/ShortifyGo/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type postRequest struct {
	URL   string `json:"URL"`
	Short string `json:"Short"`
}

type postResponse struct {
	Status        string    `json:"Status"`
	URL           string    `json:"URL"`
	Short         string    `json:"Short"`
	RemainRequest int       `json:"RemainRequest"`
	LimitRestTime time.Time `json:"LimitRestTime"`
}

// ShortenURL is the handler function for the POST /shorten endpoint
// It takes in a request body containing a URL and returns a response
// containing the short ID, the URL and the number of requests remaining
// and the time left until the request limit resets
func ShortenUrl(c *fiber.Ctx) error {
	// Get the request body and the IP of the client
	var req postRequest
	var expirationTime time.Time
	var remainReqInt int

	requestBody := c.Body()
	requestIP := c.IP()

	// Get the expiration time in minutes from the environment variable
	expiration, _ := strconv.Atoi(os.Getenv("EXPIRATION"))

	log.Printf("Received request from IP %v with body %v", requestIP, requestBody)

	// Unmarshal the request body to the postRequest struct
	json.Unmarshal(requestBody, &req)

	// Ping the Redis database
	err := database.RdsValid.Ping(database.Ctx).Err()
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":      "error",
			"description": fmt.Sprintf("%v", err),
		})
	}

	// Limit the number of request from the same IP
	remainReq, _ := database.RdsValid.Get(database.Ctx, requestIP).Result()

	if remainReq == "" {
		// Set the request limit for the IP
		_ = database.RdsValid.Set(database.Ctx, requestIP, os.Getenv("REQUEST_LIMIT"), 60*24*time.Minute).Err()
		remainReqInt, _ = strconv.Atoi(os.Getenv("REQUEST_LIMIT"))
		remainReqInt -= 1
		expirationTime = time.Now().Add(60 * 24 * time.Minute)

	} else {
		remainReqInt, _ = strconv.Atoi(remainReq)
		log.Printf("Received request from IP %v with %s requests remaining", requestIP, remainReq)
		// Get the TTL of the IP in the Redis database
		ttl, _ := database.RdsValid.TTL(database.Ctx, requestIP).Result()

		// Check if the key exists and has an expiration time
		if ttl > 0 {
			// Calculate the expiration time by subtracting the TTL from the current time
			expirationTime = time.Now().Add(-ttl)
			log.Printf("IP %v will be reset in %v hours %v minutes", requestIP, int(ttl.Hours()), int(ttl.Minutes())%60)

		} else if ttl == -1 {
			// If the TTL is -1, the key does not exist
			log.Printf("IP %v has no TTL", requestIP)
			expirationTime = time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)
		}

		ttlHours := int(ttl.Hours())
		ttlMinutes := int(ttl.Minutes()) % 60

		if remainReqInt <= 0 {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":      "error",
				"description": fmt.Sprintf("request limit exceed. It will be reset in %v hours  %v minutes", ttlHours, ttlMinutes),
			})
		}
	}

	// Validate the URL
	if !utils.IsUrl(req.URL) {
		log.Printf("URL %v is invalid", req.URL)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"description": "URL is invalid",
		})
	}

	for {
		// Generate a short ID and check if it already exists in the Redis database
		shortID := uuid.New().String()[:8]
		if res, _ := database.RdsUrl.Get(database.Ctx, shortID).Result(); res != "" {
			log.Printf("Short ID %v is already taken", shortID)
			continue
		} else {
			log.Printf("Short ID %v is not taken", shortID)
			req.Short = shortID
			break
		}

	}

	log.Printf("Generated Short ID %v for URL %v", req.Short, req.URL)

	// Set the short ID and URL in the Redis database
	_ = database.RdsUrl.Set(database.Ctx, req.Short, req.URL, time.Duration(expiration)*time.Minute).Err()

	// Decrement the request count for the IP
	_ = database.RdsValid.Decr(database.Ctx, requestIP)

	shortUrl := utils.GenerateUrl(os.Getenv("DOMAIN"), req.Short)

	resp := postResponse{
		Status:        "success",
		URL:           req.URL,
		Short:         shortUrl,
		RemainRequest: remainReqInt,
		LimitRestTime: expirationTime,
	}

	log.Printf("Returning response %v", resp)

	return c.Status(fiber.StatusOK).JSON(resp)
}
