package main

type RescueRequestModel struct {
	GarageId    int    `json:"garageId"`
	Description string `json:"description"`
}

type ChangeRescueStatusRequestModel struct {
	CustomerId int `json:"customerId"`
	Status     int `json:"status"`
}

type GarageRejectRequestModel struct {
	CustomerId   int    `json:"customerId"`
	RejectReason string `json:"rejectReason"`
}

type CustomerCancelRequestModel struct {
	GarageId     int    `json:"garageId"`
	RejectReason string `json:"rejectReason"`
}
