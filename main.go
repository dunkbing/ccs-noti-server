package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

var (
	msgClient       *messaging.Client
	firestoreClient *firestore.Client
)

func initializeAppWithServiceAccount() *firebase.App {
	// [START initialize_app_service_account_golang]
	opt := option.WithCredentialsFile("service-account.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	// [END initialize_app_service_account_golang]

	return app
}

func hello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "hello"})
}

func newRescueRequest(c *gin.Context) {
	var rescueRequest RescueRequest
	if err := c.BindJSON(&rescueRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("%+v\n", rescueRequest)
	tokens, err := getDeviceTokens(firestoreClient, c.Request.Context(), "garage-device-tokens", fmt.Sprintf("%v", rescueRequest.GarageId))

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	message := &messaging.MulticastMessage{
		Data: map[string]string{
			"type": "rescue",
		},
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: "Yêu cầu cứu hộ mới",
			Body:  fmt.Sprintf("%v", rescueRequest.Description),
		},
	}

	br, err := msgClient.SendMulticast(context.Background(), message)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// See the BatchResponse reference documentation
	// for the contents of response.
	fmt.Printf("%d messages were sent successfully\n", br.SuccessCount)
	// [END send_multicast]
	c.JSON(200, gin.H{
		"message": "success",
	})
}

func main() {
	ctx := context.Background()
	app := initializeAppWithServiceAccount()
	var err error
	firestoreClient, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing firestore: %v\n", err)
	}
	msgClient, err = app.Messaging(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	router := gin.Default()
	router.GET("/", hello)
	router.POST("/rescues", newRescueRequest)
	router.Run(":8080")
}
