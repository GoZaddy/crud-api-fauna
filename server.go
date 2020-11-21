package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gozaddy/crud-api-fauna/controllers"
	"github.com/gozaddy/crud-api-fauna/customerrors"
	"github.com/gozaddy/crud-api-fauna/database"
	"github.com/joho/godotenv"
	"log"
	"os"
)


func handle(f func(c *gin.Context) error) gin.HandlerFunc {
	return func(context *gin.Context){
		if err := f(context); err != nil{
			if ae, ok := err.(customerrors.AppError); ok{
				context.JSON(ae.StatusCode, gin.H{
					"message": ae.ErrorText,
				})
			} else {
				log.Println(err.Error())
				context.JSON(500, gin.H{
					"message": "Internal server error",
				})
			}
		}

	}
}

func init(){
	//load .env file with godotenv so we can access our FaunaDB secret
	err := godotenv.Load()
	if err != nil{
		log.Fatalln("Error loading .env file:"+ err.Error())
	}
}

func main(){
	fdb := database.NewFaunaDB(os.Getenv("FAUNA_DB_SERVER_SECRET"))
	err := fdb.Init()
	if err != nil{
		log.Fatalln(err)
	}


	controller := controllers.NewController(fdb)
	router := gin.Default()
	router.GET("/api/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome!",
		})
	})
	router.GET("/api/items", handle(controller.GetAllReadingItems))
	router.POST("/api/items", handle(controller.AddReadingItem))
	router.GET("/api/items/:id", handle(controller.GetOneReadingItem))
	router.PATCH("/api/items/:id", handle(controller.UpdateOneReadingItem))
	router.DELETE("/api/items/:id", handle(controller.DeleteOneReadingItem))

	//run our server on port 4000
	_ = router.Run(":4000")
}

