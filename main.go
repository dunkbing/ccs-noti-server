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

func NewRescueRequest(c *gin.Context) {
	var rescueRequest RescueRequestModel
	if err := c.BindJSON(&rescueRequest); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	tokens, err := getDeviceTokens(firestoreClient, c.Request.Context(), MANAGER_DEVICE_TOKENS, fmt.Sprintf("%v", rescueRequest.GarageId))

	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: "Yêu cầu cứu hộ mới",
			Body:  fmt.Sprintf("%v", rescueRequest.Description),
		},
		Data: map[string]string{
			"type": RESCUE,
		},
	}

	br, err := msgClient.SendMulticast(context.Background(), message)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("success: %v", br.SuccessCount),
	})
}

func GarageRejectRequest(c *gin.Context) {
	var rejectRequest GarageRejectRequestModel
	if err := c.BindJSON(&rejectRequest); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	tokens, err := getDeviceTokens(firestoreClient, c.Request.Context(), CUSTOMER_DEVICE_TOKENS, fmt.Sprintf("%v", rejectRequest.CustomerId))

	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: "Yêu cầu của quý khách đã bị từ chối",
			Body:  fmt.Sprintf("%v", rejectRequest.RejectReason),
		},
		Data: map[string]string{
			"type": GARAGE_REJECT_REQUEST,
		},
	}

	br, err := msgClient.SendMulticast(context.Background(), message)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("success: %v", br.SuccessCount),
	})
}

func CustomerCancelRequest(c *gin.Context) {
	var rejectRequest CustomerCancelRequestModel
	if err := c.BindJSON(&rejectRequest); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	tokens, err := getDeviceTokens(firestoreClient, c.Request.Context(), MANAGER_DEVICE_TOKENS, fmt.Sprintf("%v", rejectRequest.GarageId))

	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: "Khách hàng đã từ chối yêu cầu",
			Body:  fmt.Sprintf("%v", rejectRequest.RejectReason),
		},
		Data: map[string]string{
			"type": CUSTOMER_CANCEL_REQUEST,
		},
	}

	br, err := msgClient.SendMulticast(context.Background(), message)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("success: %v", br.SuccessCount),
	})
}

func ChangeRescueStatus(c *gin.Context) {
	var rescueRequest ChangeRescueStatusRequestModel
	if err := c.BindJSON(&rescueRequest); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	tokens, err := getDeviceTokens(firestoreClient, c.Request.Context(), CUSTOMER_DEVICE_TOKENS, fmt.Sprintf("%v", rescueRequest.CustomerId))
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "success",
			"message": err.Error(),
		})
		return
	}

	body := map[int]string{
		PENDING:  "Đang chờ",
		ACCEPTED: "Garage đã chấp nhận yêu cầu của bạn",
		ARRIVING: "Nhân viên cứu hộ đang đến",
		ARRIVED:  "Cứu hộ đã đến nơi",
		WORKING:  "Đang tiến hành sửa chữa",
		DONE:     "Đã hoàn thành sửa chữa",
		REJECTED: "Garage đã từ chối yêu cầu của bạn",
	}[rescueRequest.Status]

	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: "Tình trạng cứu hộ",
			Body:  body,
		},
		Data: map[string]string{
			"status": fmt.Sprintf("%v", rescueRequest.Status),
		},
	}

	br, err := msgClient.SendMulticast(context.Background(), message)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("success: %v", br.SuccessCount),
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
	router.POST("/rescues", NewRescueRequest)
	router.PUT("/rescues/status", ChangeRescueStatus)
	router.PUT("rescues/garage-reject", GarageRejectRequest)
	router.PUT("rescues/customer-cancel", CustomerCancelRequest)
	router.Run(":8080")
}
