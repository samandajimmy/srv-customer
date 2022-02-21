package customer

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
	"time"
)

func NewCommon(startTime time.Time, manifest ncore.Manifest) *Common {
	h := Common{
		startTime: startTime,
		manifest:  manifest,
	}
	return &h
}

type Common struct {
	startTime time.Time
	manifest  ncore.Manifest
}

func (h *Common) GetAPIStatus(_ *nhttp.Request) (*nhttp.Response, error) {

	res := nhttp.Success().
		SetData(map[string]string{
			"appVersion":     h.manifest.AppVersion,
			"buildSignature": h.manifest.BuildSignature,
			"uptime":         time.Since(h.startTime).String(),
		})
	return res, nil
}

func (h *Common) ValidateClient(r *nhttp.Request) (*nhttp.Response, error) {

	// Get subject from headers
	subjectID := r.Header.Get(constant.SubjectIDHeader)
	subjectRefID, ok := nval.ParseInt64(subjectID)
	if !ok {
		log.Error("x-subject-id is required")
		return nil, errors.New("x-subject-id is required")
	}

	//Get subject role
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

	r.SetContextValue(constant.SubjectKey, &subject)

	return nhttp.Continue(), nil
}

func GetSubject(rx *nhttp.Request) (*dto.Subject, error) {
	v := rx.GetContextValue(constant.SubjectKey)
	subject, ok := v.(*dto.Subject)
	if !ok {
		return nil, ncore.NewError("no subject found in request context")
	}
	return subject, nil
}

func GetRequestId(rx *nhttp.Request) string {
	reqId, ok := rx.GetContextValue(nhttp.RequestIdKey).(string)
	if !ok {
		// Generate new request id
		id, err := uuid.NewUUID()
		if err != nil {
			panic(fmt.Errorf("unable to generate new request id. %w", err))
		}
		return id.String()
	}

	return reqId
}
