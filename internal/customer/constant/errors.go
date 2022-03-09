package constant

import (
	"github.com/nbs-go/errx"
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
)

var b = errx.NewBuilder("customer-svc", errx.FallbackError(
	errx.NewError("500", "An error has occurred, please try again later",
		nhttp.WithStatus(http.StatusInternalServerError),
	),
))

// Common Errors

var ResourceNotFoundError = b.NewError("E_COMM_1", "Resource not found",
	nhttp.WithStatus(http.StatusNotFound))
var StaleResourceError = b.NewError("E_COMM_2", "Cannot update stale resource",
	nhttp.WithStatus(http.StatusConflict))
var InactiveError = b.NewError("E_COMM_3", "Resource is not active",
	nhttp.WithStatus(http.StatusNotFound))
var BeforePeriodError = b.NewError("E_COMM_4", "Period is not started",
	nhttp.WithStatus(http.StatusNotFound))
var AfterPeriodError = b.NewError("E_COMM_5", "Period has finished",
	nhttp.WithStatus(http.StatusNotFound))
var DuplicatedError = b.NewError("E_COMM_6", "Duplicated resource",
	nhttp.WithStatus(http.StatusConflict))
var UnknownError = b.NewError("E_COMM_7", "An error has occurred, please try again later",
	nhttp.WithStatus(http.StatusConflict))

// Authentication Errors

var InvalidCredentialError = b.NewError("E_AUTH_1", "Invalid credentials",
	nhttp.WithStatus(http.StatusBadRequest))
var UserSuspendedError = b.NewError("E_AUTH_2", "User is suspended",
	nhttp.WithStatus(http.StatusUnauthorized))
var AuthTokenExpiredError = b.NewError("E_AUTH_3", "Token has expired",
	nhttp.WithStatus(http.StatusUnauthorized))
var InvalidChangePasswordFormatError = b.NewError("E_AUTH_4", "Old and new password must be different",
	nhttp.WithStatus(http.StatusBadRequest))
var UsedPhoneNumberError = b.NewError("E_AUTH_5", "Phone number has been used",
	nhttp.WithStatus(http.StatusBadRequest))
var InvalidPhoneInput1Error = b.NewError("E_AUTH_6", "Masukan nomor handphone dan password yang sesuai. 1 kesempatan lagi sebelum akun terkunci selama 1 jam.",
	nhttp.WithStatus(http.StatusBadRequest))
var InvalidPhoneInput2Error = b.NewError("E_AUTH_7", "Masukan nomor handphone dan password yang sesuai. 1 kesempatan lagi sebelum akun terkunci selama 24 jam.",
	nhttp.WithStatus(http.StatusBadRequest))
var InvalidEmailPassInputError = b.NewError("E_AUTH_8", "Masukan email dan password yang sesuai. Pastikan email telah diverifikasi.",
	nhttp.WithStatus(http.StatusBadRequest))
var NoPhoneEmailError = b.NewError("E_AUTH_10", "No HP atau Email tidak terdaftar.",
	nhttp.WithStatus(http.StatusBadRequest))
var InvalidTokenError = b.NewError("E_AUTH_11", "Invalid token",
	nhttp.WithStatus(http.StatusBadRequest))
var InvalidPasswordError = b.NewError("E_AUTH_12", "Password tidak sesuai",
	nhttp.WithStatus(http.StatusBadRequest))

var InvalidFormatError = b.NewError("E_COMM_5", "Invalid Format")

// Asset Errors

var UnknownAssetTypeError = b.NewError("E_AST_1", "Unknown asset type",
	nhttp.WithStatus(http.StatusBadRequest))
var NoFileError = b.NewError("E_AST_2", "No file in request body",
	nhttp.WithStatus(http.StatusBadRequest))

// Registration Errors

var RegistrationFailedError = b.NewError("E_REG_1", "Registration failed",
	nhttp.WithStatus(http.StatusBadRequest))
var EmailHasBeenRegisteredError = b.NewError("E_REG_2", "Email sudah terdaftar",
	nhttp.WithStatus(http.StatusConflict))

// OTP Errors

var IncorrectOTPError = b.NewError("E_OTP_1", "incorrect otp",
	nhttp.WithStatus(http.StatusBadRequest))
var OTPReachResendLimitError = b.NewError("E_OTP_3", "Mohon tunggu 300 detik lagi",
	nhttp.WithStatus(http.StatusBadRequest))
var ExpiredOTPError = b.NewError("E_OTP_4", "otp has been expired",
	nhttp.WithStatus(http.StatusBadRequest))
