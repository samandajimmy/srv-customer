package nvalidate

import (
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ErrMessageVO struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var (
	reDigit = regexp.MustCompile("^[0-9]+$")
)

func Init() {
	customMessage()
}

func customMessage() {
	// required
	validation.ErrRequired = validation.ErrRequired.SetMessage("harus diisi")

	// email
	is.ErrEmail = is.ErrEmail.SetMessage("alamat e-mail tidak valid")
	is.Email = validation.NewStringRuleWithError(govalidator.IsExistingEmail, is.ErrEmail)

	// digit
	is.ErrDigit = validation.NewError("validation_is_digit", "harus berisi angka saja")
	is.Digit = validation.NewStringRuleWithError(isDigit, is.ErrDigit)
}

func Message(err string, additional ...*ErrMessageVO) interface{} {
	splitErrMessage := strings.Split(err, "; ")

	var messages []*ErrMessageVO

	for _, errMessage := range splitErrMessage {

		splitMessage := strings.Split(errMessage, ": ")
		message := strings.Replace(splitMessage[1], ".", "", 1)
		item := &ErrMessageVO{
			Field:   splitMessage[0],
			Message: message,
		}

		messages = append(messages, item)
	}

	if additional != nil && len(additional) != 0 {
		for _, m := range additional {
			messages = append(messages, m)
		}
	}

	return messages
}

func isDigit(value string) bool {
	return reDigit.MatchString(value)
}
