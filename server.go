package emailer

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ContactFormData struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone"`
	Message   string `json:"message" binding:"required"`
}

type HandlerDependencies struct {
	Sender EmailSender
}

func CreateContactHandler(deps HandlerDependencies, recipientEmail string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var formData ContactFormData

		if err := c.ShouldBindJSON(&formData); err != nil {
			log.Printf("Error: Failed to Bind to JSON: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid form data provided. Please check required fields and email format.",
			})
			return
		}

		subject := fmt.Sprintf("Contact Form : %v %v", formData.FirstName, formData.LastName)

		body := fmt.Sprintf("Name: %v %v\nEmail: %v\nPhone: %v\nMessage:%v", formData.FirstName, formData.LastName, formData.Email, formData.Phone, formData.Message)

		err := deps.Sender.Send(recipientEmail, subject, body)
		if err != nil {
			log.Printf("Error: Failed to send Email via sender : %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed due to internal server error",
			})
			return
		}
		log.Printf("Info: Success!! contact form processed from %v", formData.Email)
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "message received succesfully",
		})
	}
}

func SetupRouter(deps HandlerDependencies, recipientEmail string) *gin.Engine {
	router := gin.Default()

	router.Use(cors.Default())

	router.POST("/contact", CreateContactHandler(deps, recipientEmail))
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	return router
}
