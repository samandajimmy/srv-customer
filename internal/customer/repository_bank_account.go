package customer

import (
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/op"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/pq/query"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/statement"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	nsqlDep "repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

// Define filters
var bankAccountFilters = map[string]nsql.FilterParser{
	constant.AccountNameKey:   query.LikeFilter("accountName", op.LikeSubString, option.Schema(statement.BankAccountSchema)),
	constant.AccountNumberKey: query.LikeFilter("accountNumber", op.LikeSubString, option.Schema(statement.BankAccountSchema)),
}

func (rc *RepositoryContext) ListBankAccount(customerID int64, params *dto.ListPayload) (*model.ListBankAccountResult, error) {
	// Init query builder
	b := query.From(statement.BankAccountSchema)

	// Set where
	filters := query.NewFilter(params.Filters, bankAccountFilters)
	b.Where(
		query.And(
			query.Equal(query.Column("customerId", query.Schema(statement.BankAccountSchema)), query.BindVar()),
		),
		filters.Conditions(),
	)

	// Set order by
	switch params.SortBy {
	case constant.SortByCreated:
		b.OrderBy("createdAt")
	default:
		params.SortBy = constant.SortByLastCreated
		b.OrderBy("createdAt", option.SortDirection(op.Descending))
	}

	// Create select query
	selectQuery := b.Select(query.Column("*")).Limit(params.Limit).Skip(params.Skip).Build()
	selectQuery = rc.conn.Rebind(selectQuery)

	// Create count query
	b.ResetOrderBy().ResetSkip().ResetLimit()
	countQuery := b.Select(query.Count("*", option.Schema(statement.BankAccountSchema), option.As("count"))).Build()
	countQuery = rc.conn.Rebind(countQuery)

	// Combine arguments with customerId from filters
	args := append([]interface{}{customerID}, filters.Args()...)

	// Execute query
	var rows []model.BankAccount
	err := rc.conn.SelectContext(rc.ctx, &rows, selectQuery, args...)
	if err != nil {
		return nil, errx.Trace(err)
	}

	var count int64
	err = rc.conn.GetContext(rc.ctx, &count, countQuery, args...)
	if err != nil {
		return nil, errx.Trace(err)
	}

	// Prepare result
	result := model.ListBankAccountResult{
		Rows:  rows,
		Count: count,
	}
	return &result, err
}

func (rc *RepositoryContext) CreateBankAccount(row *model.BankAccount) error {
	_, err := rc.stmt.BankAccount.Insert.ExecContext(rc.ctx, &row)
	if err != nil {
		return err
	}
	return nil
}

func (rc *RepositoryContext) FindBankAccountByXID(xid string) (*model.BankAccount, error) {
	var result model.BankAccount
	err := rc.stmt.BankAccount.FindByXID.GetContext(rc.ctx, &result, xid)
	return &result, err
}

func (rc *RepositoryContext) UpdateBankAccount(bankAccount *model.BankAccount) error {
	result, err := rc.stmt.BankAccount.Update.ExecContext(rc.ctx, bankAccount)
	if err != nil {
		return errx.Trace(err)
	}
	if !nsqlDep.IsUpdated(result) {
		return constant.BankAccountNotFoundError
	}
	return nil
}

func (rc *RepositoryContext) DeleteBankAccountByXID(xid string) error {
	result, err := rc.stmt.BankAccount.DeleteByXID.ExecContext(rc.ctx, xid)
	if err != nil {
		return errx.Trace(err)
	}
	if !nsqlDep.IsUpdated(result) {
		return constant.BankAccountNotFoundError
	}
	return nil
}

func (rc *RepositoryContext) FindBankAccountByAccountNumberAndCustomerID(accountNumber string, customerID int64) (*model.BankAccount, error) {
	var result model.BankAccount
	err := rc.stmt.BankAccount.FindByAccountNumberAndCustomerID.GetContext(rc.ctx, &result, accountNumber, customerID)
	return &result, err
}
