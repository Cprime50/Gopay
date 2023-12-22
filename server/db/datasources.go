package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"context"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type dataSources struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

// InitDS establishes connections to fields in dataSources
func initDS() (*dataSources, error) {
	log.Printf("Initializing data sources\n")

	dsn := os.Getenv("DATABASE_URL")

	log.Printf("Connecting to Mysql\n")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
	}})

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Mysql database: %w", err)
	}

	//Enable pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Verify database connection is working
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	// Initialize redis connection
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	log.Printf("Connecting to Redis\n")
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "",
		DB:       0,
	})

	// verify redis connection
	_, err = rdb.Ping(context.Background()).Result()

	if err != nil {
		return nil, fmt.Errorf("error connecting to redis: %w", err)
	}

	return &dataSources{
		DB:          db,
		RedisClient: rdb,
	}, nil
}

// close to be used in graceful server shutdown
func (ds *dataSources) close() error {
	sqlDB, err := ds.DB()
	if err != nil {
		fmt.Error(err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("error closing Mysql: %w", err)
	}

	if err := ds.RedisClient.Close(); err != nil {
		return fmt.Errorf("error closing Redis Client: %w", err)
	}

	return nil
}
