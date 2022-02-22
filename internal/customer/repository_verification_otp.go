package customer

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

func (rc *RepositoryContext) CreateVerificationOTP(row *model.VerificationOTP) (int64, error) {
	var lastInsertID int64
	err := rc.stmt.VerificationOTP.Insert.QueryRowContext(rc.ctx, &row).Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}
	return lastInsertID, nil
}

func (rc *RepositoryContext) FindVerificationOTPByRegistrationIDAndPhone(id string, phone string) (*model.VerificationOTP, error) {
	var row model.VerificationOTP
	err := rc.stmt.VerificationOTP.FindByRegistrationID.GetContext(rc.ctx, &row, id, phone)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (rc *RepositoryContext) DeleteVerificationOTP(id string, phone string) error {
	result, err := rc.stmt.VerificationOTP.Delete.ExecContext(rc.ctx, id, phone)
	if err != nil {
		return ncore.TraceError("failed to delete verification otp", err)
	}
	if !nsql.IsUpdated(result) {
		return constant.ResourceNotFoundError
	}
	return nil
}
