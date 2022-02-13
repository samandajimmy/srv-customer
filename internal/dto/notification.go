package dto

type NotificationPayload struct {
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Image string            `json:"image"`
	Token string            `json:"token"`
	Data  map[string]string `json:"data"`
}

type NotificationRegister struct {
	Customer     interface{} // *model.Customer
	Verification interface{} // *model.Verification
	RegisterOTP  interface{} // *model.VerificationOTP
	Payload      RegisterNewCustomer
}

type NotificationBlock struct {
	Customer     interface{} // *model.Customer
	Message      string
	LastTryLogin string
}
