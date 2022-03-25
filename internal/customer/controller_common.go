package customer

import (
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

const (
	AnonymousUserID       = "ANON"
	AnonymousUserRefID    = 0
	AnonymousUserFullName = "Anonymous User"
)

func (c *CommonController) ParseSubject(r *nhttp.Request) (*nhttp.Response, error) {
	// Get subject from headers
	id := r.Header.Get(constant.SubjectIDHeader)
	if id == "" {
		id = AnonymousUserID
	}

	// Get subject reference id
	refID, ok := nval.ParseInt64(id)
	if !ok {
		refID = AnonymousUserRefID
	}

	// Get subject role and determine subject type
	role := r.Header.Get(constant.SubjectRoleHeader)
	var subjectType constant.SubjectType
	switch role {
	case constant.AdminModifierRole, constant.UserModifierRole:
		subjectType = constant.UserSubjectType
	case constant.SystemModifierRole:
		subjectType = constant.SystemSubjectType
	default:
		// Fallback to anonymous user
		subjectType = constant.UserSubjectType
		role = constant.UserModifierRole
	}

	// Get name
	fullName := r.Header.Get(constant.SubjectNameHeader)
	if fullName == "" {
		fullName = AnonymousUserFullName
	}

	subject := dto.Subject{
		ID:          id,
		RefID:       refID,
		Role:        role,
		FullName:    fullName,
		SubjectType: subjectType,
		Metadata:    map[string]string{},
		SessionID:   0,
	}

	r.SetContextValue(constant.SubjectContextKey, &subject)

	return nhttp.Continue(), nil
}

func GetSubject(rx *nhttp.Request) *dto.Subject {
	v := rx.GetContextValue(constant.SubjectContextKey)
	subject, ok := v.(*dto.Subject)
	if !ok {
		// Return anonymous subject
		return &dto.Subject{
			ID:          AnonymousUserID,
			RefID:       AnonymousUserRefID,
			Role:        constant.UserModifierRole,
			FullName:    AnonymousUserFullName,
			SubjectType: constant.UserSubjectType,
			SessionID:   0,
			Metadata:    map[string]string{},
		}
	}
	return subject
}
