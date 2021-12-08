package nvalidate

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ErrMessageVO struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Init() {
	customMessage()
}

func customMessage() {
	validation.ErrRequired = validation.ErrRequired.SetMessage("harus diisi")
}

func Message(err string) interface{} {
	splitErrMessage := strings.Split(err, "; ")

	var messages []*ErrMessageVO

	for _, errMessage := range splitErrMessage {

		splitMessage := strings.Split(errMessage, ": ")
		item := &ErrMessageVO{
			Field:   splitMessage[0],
			Message: splitMessage[1],
		}

		messages = append(messages, item)
	}

	return messages
}
