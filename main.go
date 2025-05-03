package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tablestar/porto-emailer/emailer"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("WARN: Environment variable %v not set, using default: %v\n", key, fallback)
	return fallback
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Printf("ERROR loading env file : %v.\nrelying on externally set env.\n", err)
	} else {
		log.Println("Successfully loaded .env")
	}

	smtpHost := getEnv("SMTP_HOST", "")
	smtpPort := getEnv("SMTP_PORT", "465")
	smtpUser := getEnv("SMTP_USERNAME", "")
	smtpPass := getEnv("SMTP_PASS", "")
	smtpFrom := getEnv("SMTP_USERNAME", smtpUser)

	recipientEmail := getEnv("RECIPIENT_EMAIL", "")
	serverPort := getEnv("PORT", "8080")

	if smtpHost == "" || smtpUser == "" || smtpPass == "" || recipientEmail == "" {
		log.Fatal("FATAL: Missing required environment variables: SMTP_HOST, SMTP_USERNAME, SMTP_PASS, RECIPIENT_EMAIL must be set either in .env or system environment.")
	}
	log.Println("INFO loaded successfully")

	emailSender := emailer.NewSmtpSender(smtpHost, smtpPort, smtpUser, smtpPass, smtpFrom)

	log.Println("Email Sender created")

	deps := emailer.HandlerDependencies{
		Sender: emailSender,
	}

	router := emailer.SetupRouter(deps, recipientEmail)

	serverAddress := ":" + serverPort
	log.Printf("Starting http server on port %v", serverAddress)

	if err := router.Run(serverAddress); err != nil {
		log.Fatalf("FATAL: WHY GIN NOT START\nError:%v", err)
	}

}
