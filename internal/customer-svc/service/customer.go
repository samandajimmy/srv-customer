package service

import (
	"strings"
	"time"

	"github.com/rs/xid"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/convert"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nclient"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
)

type Customer struct {
	customerRepo        contract.CustomerRepository
	verificationOTPRepo contract.VerificationOTPRepository
	otpService          contract.OTPService
	response            *ncore.ResponseMap
}

func (c *Customer) HasInitialized() bool {
	return true
}

func (c *Customer) Init(app *contract.PdsApp) error {
	c.customerRepo = app.Repositories.Customer
	c.verificationOTPRepo = app.Repositories.VerificationOTP
	c.otpService = app.Services.OTP
	c.response = app.Responses
	return nil
}

func (c *Customer) Login(payload dto.LoginRequest) (*dto.CustomerVO, error) {

	// check user exists
	customer := c.customerRepo.FindByEmailOrPhone(payload.Email)
	if customer == nil {
		return nil, c.response.GetError("E_AUTH_8")
	}

	// counter wrong password count
	customer.WrongPasswordCount += 1

	if customer.WrongPasswordCount == 2 {
		return nil, c.response.GetError("E_AUTH_6")
	} else if customer.WrongPasswordCount == 4 {
		return nil, c.response.GetError("E_AUTH_7")
	}

	// TODO Check account is first login or not

	// TODO Store audit login

	// TODO: update user_model -> try_login_date = now()

	// if password doesn't match
	// 	cek setBlockedUser function
	//    if blocked_to_date > now()
	//    	return err_account_locked message

	if customer.Password != payload.Password {
		//
	}

	// if password is matched
	//    update user_model
	//   		set blocked_date = null
	//        set blocked_to_date = null
	//        wrong_password_count = 0

	// check user account is blocked or not

	// set token authentication

	// get user data

	// get tabungan emas service

	// check is force update password

	// return response user and token

	return nil, nil
}

func (c *Customer) Register(payload dto.RegisterNewCustomer) (*dto.RegisterNewCustomerResponse, error) {

	// find registerID

	// Get data user
	customer := c.customerRepo.FindByPhone(payload.PhoneNumber)

	// Check if user is active
	if customer.Status == 1 {
		log.Errorf("Phone has been used. Phone Number : %s", payload.PhoneNumber)
		return nil, c.response.GetError("E_AUTH_5")
	}

	// If user exist and inactive
	if customer.Phone != "" && customer.Status == 0 {
		// TODO: update fullname on customer
	}

	// Register customer
	customerXID := strings.ToUpper(xid.New().String())
	insert := &model.Customer{
		CustomerXID:    customerXID,
		FullName:       payload.Name,
		Phone:          payload.PhoneNumber,
		Status:         0,
		Email:          "",
		IdentityType:   0,
		IdentityNumber: "",
		UserRefId:      0,
		Photos:         []byte("{}"),
		Profile:        []byte("{}"),
		Cif:            "",
		Sid:            "",
		ReferralCode:   "",
		Metadata:       []byte("{}"),
		ItemMetadata:   model.NewItemMetadata(convert.ModifierDTOToModel(dto.Modifier{ID: "", Role: "", FullName: ""})),
	}
	_, err := c.customerRepo.Insert(insert)
	if err != nil {
		log.Errorf("Error when persist customer : %s", payload.Name)
		return nil, ncore.TraceError(err)
	}

	// Save Code OTP and userId to VerficationOTP
	// TODO: Save Code OTP and userId to VerificationOTP

	// Send OTP
	// TODO: Send OTP to coreService via client http

	// If send otp is success
	// TODO: Update customer data

	return &dto.RegisterNewCustomerResponse{
		Token: customerXID,
	}, nil
}

func (c *Customer) RegisterStepOne(payload dto.RegisterStepOne) (*dto.RegisterStepOneResponse, error) {

	// TODO Validate Phone Number If Exist

	// TODO Validate Email If Exist

	// Set request
	request := dto.SendOTPRequest{
		PhoneNumber: payload.PhoneNumber,
		RequestType: "register",
	}

	// Send OTP To Phone Number
	resp, err := c.otpService.SendOTP(request)
	if err != nil {
		return nil, ncore.TraceError(err)
	}

	// Extract response from server
	data, err := nclient.GetResponseData(resp)

	return &dto.RegisterStepOneResponse{
		Action: data.Message,
	}, nil
}

func (c *Customer) RegisterStepTwo(payload dto.RegisterStepTwo) (*dto.RegisterStepTwoResponse, error) {
	// Set request
	request := dto.VerifyOTPRequest{
		PhoneNumber: payload.PhoneNumber,
		Token:       payload.OTP,
		RequestType: "register",
	}

	// Verify OTP To Phone Number
	resp, err := c.otpService.VerifyOTP(request)
	if err != nil {
		return nil, ncore.TraceError(err)
	}

	// Extract response from server
	data, err := nclient.GetResponseData(resp)

	// wrong otp handle
	if data.ResponseCode != "00" {
		log.Errorf("Wrong OTP. Phone Number : %s", payload.PhoneNumber)
		return nil, c.response.GetError("E_OTP_1")
	}

	registrationId := xid.New().String()
	// Check OTP Wrong
	insert := &model.VerificationOTP{
		CreatedAt:      time.Now(),
		Phone:          payload.PhoneNumber,
		RegistrationId: registrationId,
	}

	_, err = c.verificationOTPRepo.Insert(insert)
	if err != nil {
		log.Errorf("Error when persist verificationOTP. Phone Number: %s", payload.PhoneNumber)
		return nil, ncore.TraceError(err)
	}

	return &dto.RegisterStepTwoResponse{
		RegisterId: registrationId,
	}, nil
}

func (c *Customer) RegisterResendOTP(payload dto.RegisterResendOTP) (*dto.RegisterResendOTPResponse, error) {
	// Set request
	request := dto.SendOTPRequest{
		PhoneNumber: payload.PhoneNumber,
		RequestType: "register",
	}

	// Send OTP To Phone Number
	resp, err := c.otpService.SendOTP(request)
	if err != nil {
		return nil, ncore.TraceError(err)
	}

	// Extract response from server
	data, err := nclient.GetResponseData(resp)

	// wrong otp handle
	if data.ResponseCode != "00" {
		log.Errorf("Wrong OTP. Phone Number : %s", payload.PhoneNumber)
		return nil, c.response.GetError("E_OTP_2")
	}

	return &dto.RegisterResendOTPResponse{
		Action: data.Message,
	}, nil
}

func (c *Customer) RegisterSubmit(payload dto.RegisterResendOTP) (*dto.RegisterResendOTPResponse, error) {

	return nil, nil
}
