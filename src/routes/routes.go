package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yopaz-huytc/go-crud/src/controllers"
)

// Routes is a function that defines all the routes of the application
func Routes() {
	router := gin.Default()

	router.POST("/todo", controllers.CreateTodo)
	router.GET("/todo", controllers.GetAllTodos)
	router.PUT("/todo/:idTodo", controllers.UpdateTodo)
	router.DELETE("/todo/:idTodo", controllers.DeleteTodo)
	router.POST("/upload", controllers.UploadFile)

	err := router.Run()
	if err != nil {
		return
	}
}
