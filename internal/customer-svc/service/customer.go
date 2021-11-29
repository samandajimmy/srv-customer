package service

import (
	"github.com/rs/xid"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/convert"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"strings"
)

type Customer struct {
	customerRepo contract.CustomerRepository
	response     *ncore.ResponseMap
}

func (c *Customer) HasInitialized() bool {
	return true
}

func (c *Customer) Init(app *contract.PdsApp) error {
	c.customerRepo = app.Repositories.Customer
	c.response = app.Responses
	return nil
}

func (c *Customer) Register(payload dto.RegisterNewCustomer) (*dto.NewRegisterResponse, error) {

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
	userId, err := c.customerRepo.Insert(insert)
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

	return &dto.NewRegisterResponse{
		Token:  customerXID,
		ReffId: userId,
	}, nil
}
