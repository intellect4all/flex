package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberAdapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
)

var fiberLambda *fiberAdapter.FiberLambda

func init() {
	// get the environment variables and initialize the application
	initResponse, err := InitializationHandler()
	if err != nil {
		panic(err)
	}

	var app *fiber.App
	app = fiber.New()

	// register the authentication routes
	RegisterAuth(initResponse.mongoDbClient, context.Background(), app)

	if !initResponse.environmentIsLocal {
		fiberLambda = fiberAdapter.New(app)
		return
	}

	err = app.Listen(":3000")
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	if fiberLambda != nil {
		lambda.Start(Handler)
	}

}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return fiberLambda.ProxyWithContext(ctx, req)
}
