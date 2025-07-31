package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	cr "gostubproject/stub/controller"
	"gostubproject/stub/drivers"
	tr "gostubproject/stub/repos/tasksrepo"
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

	taskDBAccessor := tr.NewTaskDbAccessorImpl(db)
	taskRepo := tr.NewTasksRepo(taskDBAccessor)
	appController := cr.NewAppController(taskRepo)

	go func() {
		router.GET("/", func(c *gin.Context) { c.JSON(200, "Hello") })
		router.GET("/ping", appController.Ping)
		router.GET("/task/:id", appController.GetTaskByID)
		router.POST("/task/create", appController.CreateTask)
		router.PUT("/task/update", appController.UpdateTask)
		router.DELETE("/task/delete/:id", appController.DeleteTask)
		router.POST("/tasks", appController.ListTasksPaginated)

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
