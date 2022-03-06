package customer

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nbs-go/errx"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"time"
)

func NewCommonController(startTime time.Time, manifest ncore.Manifest) *CommonController {
	h := CommonController{
		startTime: startTime,
		manifest:  manifest,
	}
	return &h
}

type CommonController struct {
	startTime time.Time
	manifest  ncore.Manifest
}

func (c *CommonController) GetAPIStatus(_ *nhttp.Request) (*nhttp.Response, error) {
	res := nhttp.Success().
		SetData(map[string]string{
			"appVersion":     c.manifest.AppVersion,
			"buildSignature": c.manifest.BuildSignature,
			"uptime":         time.Since(c.startTime).String(),
		})
	return res, nil
}

func (c *CommonController) ValidateClient(r *nhttp.Request) (*nhttp.Response, error) {
	// Get subject from headers
	subjectID := r.Header.Get(constant.SubjectIDHeader)
	subjectRefID, ok := nval.ParseInt64(subjectID)
	if !ok {
		log.Error("x-subject-id is required")
		return nil, errors.New("x-subject-id is required")
	}

	// Get subject role
	subjectRole := r.Header.Get(constant.SubjectRoleHeader)
	role := constant.AdminModifierRole
	if subjectRole != constant.AdminModifierRole {
		role = constant.UserModifierRole
	}

	subject := dto.Subject{
		SubjectID:    subjectID,
		SubjectRefID: subjectRefID,
		SubjectType:  constant.UserSubjectType,
		SubjectRole:  role,
		ModifiedBy: dto.Modifier{
			ID:       subjectID,
			Role:     role,
			FullName: r.Header.Get(constant.SubjectNameHeader),
		},
		Metadata: nil,
	}

	r.SetContextValue(constant.SubjectContextKey, &subject)

	return nhttp.Continue(), nil
}

func GetSubject(rx *nhttp.Request) (*dto.Subject, error) {
	v := rx.GetContextValue(constant.SubjectContextKey)
	subject, ok := v.(*dto.Subject)
	if !ok {
		return nil, errx.Trace(errors.New("no subject found in request context"))
	}
	return subject, nil
}

func GetRequestID(rx *nhttp.Request) string {
	reqID, ok := rx.GetContextValue(nhttp.RequestIDContextKey).(string)
	if !ok {
		// Generate new request id
		id, err := uuid.NewUUID()
		if err != nil {
			panic(fmt.Errorf("unable to generate new request id. %w", err))
		}
		return id.String()
	}

	return reqID
}
