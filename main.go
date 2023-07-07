package main

import (
	"go-payment/controllers"
	auth_controller "go-payment/controllers"
	"go-payment/middlewares"
	"go-payment/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	r := gin.New()
	r.Use(gin.Logger())
	gin.SetMode(gin.DebugMode)

	models.ConnectDatabase()

	public := r
	// public
	public.GET("/", controllers.Index)
	public.POST("/register", auth_controller.Register)
	public.POST("/login", auth_controller.Login)
	public.GET("/logout", auth_controller.Logout)

	protected := r.Group("/api/v1")
	protected.Use(middlewares.JwtAuthMiddleware())

	// transaction
	protected.GET("/transactions", controllers.FindTransactions)
	protected.GET("/transaction/:id", controllers.FindTransactionByID)
	protected.POST("/transaction", controllers.CreateTransaction)
	protected.POST("/transaction/approve", controllers.ApproveTransaction)

	// product
	protected.GET("/products", controllers.FindProducts)
	protected.GET("/product/:id", controllers.FindProductById)
	protected.POST("/product", controllers.CreateProduct)
	protected.PUT("/product/:id", controllers.UpdateProduct)
	protected.DELETE("/product/:id", controllers.DeleteProduct)

	// user
	protected.GET("/users", controllers.FindUsers)
	protected.GET("/user/:id", controllers.FindUserById)
	protected.POST("/user", controllers.CreateUser)
	protected.PUT("/user/:id", controllers.UpdateUser)
	protected.DELETE("/user/:id", controllers.DeleteUser)

	// auth get user login
	protected.GET("/user", auth_controller.CurrentUser)

	r.Run(":9898")
}
