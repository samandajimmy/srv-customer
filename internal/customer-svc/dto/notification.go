package dto

type NotificationPayload struct {
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Image string            `json:"image"`
	Token string            `json:"token"`
	Data  map[string]string `json:"data"`
}
