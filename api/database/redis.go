package database

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var RdsValid *redis.Client
var RdsUrl *redis.Client

// CreateRedisClient creates a new Redis client with the specified database number
// The client is configured with the address, password, and pool size specified
// in the environment variables DB_ADDR, DB_PASSWORD, and CLIENT_POOL_SIZE
func CreateRedisClient(dbNo int) *redis.Client {
	log.Println("CreateRedisClient started")

	// Get the pool size from the environment variable
	poolSize, err := strconv.Atoi(os.Getenv("CLIENT_POOL_SIZE"))
	if err != nil {
		// If the pool size is not a valid integer, return an error
		return nil
	}

	// Create a new Redis client with the specified options
	client := redis.NewClient(&redis.Options{
		// Address of the Redis instance
		Addr: os.Getenv("DB_ADDR"),
		// Password for the Redis instance
		Password: os.Getenv("DB_PASSWORD"),
		// Database number to use
		DB: dbNo,
		// Number of connections in the pool
		PoolSize: poolSize,
	})
	log.Printf("client %v", client)
	return client
}

func init() {
	// Load .env file
	if err := godotenv.Load("/app/.env"); err != nil {
		log.Println("No .env file found")
	}

	RdsValid = CreateRedisClient(1)
	if RdsValid == nil {
		log.Fatal("Failed to create Redis client 1")
	} else {
		log.Println("Redis client 2 created successfully")
	}

	RdsUrl = CreateRedisClient(2)
	if RdsUrl == nil {
		log.Fatal("Failed to create Redis client 2")
	} else {
		log.Println("Redis client 2 created successfully")
	}

}
