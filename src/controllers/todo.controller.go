package controllers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"
	"github.com/yopaz-huytc/go-crud/src/config"
	"github.com/yopaz-huytc/go-crud/src/models"
	"gorm.io/gorm"
)

var db *gorm.DB = config.ConnectDB()
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Todo struct for request body
type todoRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	IsDone      int    `json:"is_done" validate:"gte=0,lte=1"`
}

// Defining the struct for the response body
type todoResponse struct {
	todoRequest
	ID uint `json:"id"`
}

// CreateTodo Create todo data to database
func CreateTodo(context *gin.Context) {
	var data todoRequest
	// Binding request body json to request body struct
	if err := context.ShouldBindJSON(&data); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//validate request body
	validationErr := validate.Struct(data)
	if validationErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return

	}
	// Matching todo models struct with todo request struct
	todo := models.Todo{}
	todo.Name = data.Name
	todo.Description = data.Description
	todo.IsDone = 0

	// Querying to database
	result := db.Create(&todo)
	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	// Matching result to create response
	var response todoResponse
	response.ID = todo.ID
	response.Name = todo.Name
	response.Description = todo.Description
	response.IsDone = todo.IsDone

	//create http response
	context.JSON(http.StatusCreated, response)
}

func GetAllTodos(context *gin.Context) {
	var todos []models.Todo
	// Querying to find todo data
	err := db.Find(&todos)
	if err.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error.Error()})
		return
	}
	// Creating http response
	context.JSON(http.StatusOK, gin.H{
		"status":  "200",
		"message": "Success",
		"data":    todos,
	})
}

func UpdateTodo(context *gin.Context) {
	var data todoRequest
	// Defining request parameter to get todo id
	reqParamId := context.Param("idTodo")
	idTodo := cast.ToUint(reqParamId)
	// Binding request body json to request body struct
	if err := context.ShouldBindJSON(&data); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//validate request body
	validationErr := validate.Struct(data)
	if validationErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	// Initiate models todo
	todo := models.Todo{}
	// Querying to find todo data by id from request parameter
	result := db.Where("id = ?", idTodo).First(&todo)
	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	// Updating todo with data from request
	todo.Name = data.Name
	todo.Description = data.Description
	todo.IsDone = data.IsDone
	// Saving updated todo to the database
	result = db.Save(&todo)
	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	// Matching todo request with todo models
	var response todoResponse
	response.ID = todo.ID
	response.Name = todo.Name
	response.Description = todo.Description
	response.IsDone = todo.IsDone
	//  Creating http response
	context.JSON(http.StatusOK, response)
}

// DeleteTodo Delete todo data by id
func DeleteTodo(context *gin.Context) {
	// Initiate todo models
	todo := models.Todo{}
	// getting request parameter id
	reqParamId := context.Param("idTodo")
	idTodo := cast.ToUint(reqParamId)
	// Querying to delete todo data by id
	result := db.Where("id = ?", idTodo).First(&todo)
	if result.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"status":  "404",
			"message": "Data not found",
			"error":   result.Error.Error(),
		})
		return
	}
	result = db.Delete(&todo)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status":  "500",
			"message": "Error deleting record",
			"error":   result.Error.Error(),
		})
		return
	}
	fmt.Println(result) // print the result of the Delete operation
	// Creating http response
	context.JSON(http.StatusOK, gin.H{
		"status":  "200",
		"message": "Success",
		"data":    idTodo,
	})
}

func UploadFile(context *gin.Context) {
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if file.Size > 10*1024*1024 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "File size should not exceed 10MB"})
		return
	}

	// Create the uploads directory if it doesn't exist
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if the file already exists
	filePath := "./uploads/" + file.Filename
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "File already exists"})
		return
	}

	// Create a new file in the uploads directory
	out, err := os.Create(filePath)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}(out)

	// Open the uploaded file
	in, err := file.Open()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func(in multipart.File) {
		err := in.Close()
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}(in)

	// Copy the uploaded file to the new file
	_, err = io.Copy(out, in)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file_path": filePath})
}
