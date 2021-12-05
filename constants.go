package main

const (
	PENDING  = iota
	ACCEPTED = iota
	ARRIVING = iota
	ARRIVED  = iota
	WORKING  = iota
	REJECTED = iota
	DONE     = iota
)

const (
	MANAGER_DEVICE_TOKENS  = "manager-device-tokens"
	CUSTOMER_DEVICE_TOKENS = "customer-device-tokens"
)

const (
	RESCUE                  = "rescue"
	CUSTOMER_CANCEL_REQUEST = "customer-cancel-request"
	GARAGE_REJECT_REQUEST   = "garage-reject-request"
)
