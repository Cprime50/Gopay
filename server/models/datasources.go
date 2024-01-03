package models

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Cprime50/Gopay/helper"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DataSources struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

// GetDB returns the gorm.DB instance
func (ds *DataSources) GetDB() *gorm.DB {
	return ds.DB
}

// InitDS establishes connections to fields in dataSources
func InitDS() (*DataSources, error) {
	log.Printf("Initializing data sources\n")

	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, username, password, dbname, port)

	log.Printf("Connecting to Postgres\n")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
	}})

	if err != nil {
		log.Fatalf("Failed to connect to Postgres database: %s", err)
		return nil, helper.NewInternal()
	}

	//Enable pooling
	// sqlDB, err := db.DB()
	// if err != nil {
	// 	return nil, err
	// }
	// sqlDB.SetMaxIdleConns(10)
	// sqlDB.SetMaxOpenConns(100)
	fmt.Println("Connected to postgres successfully")

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
	fmt.Println("Connected to redis successfully")

	ds := &DataSources{
		DB:          db,
		RedisClient: rdb,
	}
	return ds, nil
}

// close to be used in graceful server shutdown
func (ds *DataSources) Close() error {
	sqlDB, err := ds.DB.DB()
	if err != nil {
		log.Fatal(err)
		return helper.NewInternal()
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatal("error closing Postgresql: %w", err)
		return helper.NewInternal()
	}

	if err := ds.RedisClient.Close(); err != nil {
		log.Fatal("error closing Redis Client: %w", err)
		return helper.NewInternal()
	}

	return nil
}
