package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"gostubproject/stub/drivers"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	ctx := context.Background()

	config := drivers.LoadConfiguration()
	db, dbClose := drivers.InitDB(ctx, config)
	defer func() {
		log.Println(" >. Closing Database Connection")
		dbClose()
	}()

	drivers.RunMigrate(ctx, db, "")
	// TODO remove this
	err := db.Ping(ctx)
	if err != nil {
		panic(err)
	}

	var server *http.Server
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	go func() {
		router.GET("/", func(c *gin.Context) { c.JSON(200, "Hello") })
		server = &http.Server{
			Addr:    ":8080",
			Handler: router,
		}
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println("http server closed:", err)
		}
	}()
	<-done
	log.Println(" >. Stopping Http Server")

	if err := server.Shutdown(ctx); err != nil {
		log.Println("http server shutdown:", err)
	}
}
