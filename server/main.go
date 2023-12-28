package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	db "github.com/Cprime50/Gopay/db"
	"github.com/Cprime50/Gopay/migrations"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	serveApp()

}

func serveApp() {
	//Load env
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	log.Println(".env file loaded successfully")

	// load database
	dataSources := &db.DataSources{}
	_, _err := dataSources.InitDS(context.Background())
	if _err != nil {
		log.Fatalf("Unable to initialize data sources: %v\n", err)
	}

	//run migrations
	migrate := &migrations.Database{}
	migrate.Migrate()

	router := gin.Default()

	srv := &http.Server{
		Addr:         ":8082", // Good practice to set timeouts to avoid Slow-loris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
		//log.Println("server running on https://localhost:8082")
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 2 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 2 seconds.")
	}
	log.Println("Server exiting")
}
