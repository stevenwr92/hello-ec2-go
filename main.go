package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/secret", func(c *fiber.Ctx) error {
		secretName := "Encryptor-B2B"
		region := "ap-southeast-3"

		config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		if err != nil {
			log.Fatal(err)
		}

		// Create Secrets Manager client
		svc := secretsmanager.NewFromConfig(config)

		input := &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String(secretName),
			VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
		}

		result, err := svc.GetSecretValue(context.TODO(), input)
		if err != nil {
			log.Fatal(err.Error())
		}

		var secretString string = *result.SecretString
		fmt.Println(secretString)
		return c.SendString(secretString)
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		// Check the database connection
		db, err := sql.Open("mysql", "devDBAdm:DB-B3azy*@tcp(devdbac.cluster-cry4cich1xjp.ap-southeast-3.rds.amazonaws.com:3306)/b2b_ximply_production")
		if err != nil {
			log.Fatal(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		defer db.Close()

		err = db.Ping()
		if err != nil {
			log.Fatal(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Retrieve additional information about the database
		var version string
		err = db.QueryRow("SELECT VERSION()").Scan(&version)
		if err != nil {
			log.Fatal(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Build response with database information
		response := fmt.Sprintf("Database connection is healthy!\nDatabase Version: %s", version)

		return c.SendString(response)
	})

	err := app.Listen(":80")
	if err != nil {
		panic(err)
	}
}
