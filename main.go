package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/secret", func(c *fiber.Ctx) error {
		secretName := "Encryptor-B2B" // Replace with your secret name

		secretValue, err := getSecret(secretName)
		if err != nil {
			return c.Status(500).SendString("Error retrieving secret")
		}

		return c.SendString(fmt.Sprintf("Secret Value: %s", secretValue))
	})

	err := app.Listen(":80")
	if err != nil {
		panic(err)
	}
}

func getSecret(secretName string) (string, error) {
	region := "ap-southeast-3" // Replace with your AWS region
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return "", err
	}

	svc := secretsmanager.New(sess)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return "", err
	}

	return *result.SecretString, nil
}
