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
)

// Define filters
var favoriteTransactionFilters = map[string]nsql.FilterParser{
	constant.TypeKey: query.LikeFilter("type", op.LikeSubString, option.Schema(statement.TransactionFavoriteSchema)),
}

func (rc *RepositoryContext) CreateFavorite(row *model.TransactionFavorite) error {
	_, err := rc.stmt.TransactionFavorite.Insert.ExecContext(rc.ctx, &row)
	if err != nil {
		return err
	}
	return nil
}

func (rc *RepositoryContext) FindFavoriteByAccountNumberAndCustomerID(accountNumber string, customerID int64) (*model.TransactionFavorite, error) {
	var result model.TransactionFavorite
	err := rc.stmt.TransactionFavorite.FindByAccountNumberAndCustomerID.GetContext(rc.ctx, &result, accountNumber, customerID)
	return &result, err
}

func (rc *RepositoryContext) ListFavorite(customerID int64, params *dto.ListPayload) (*model.ListTransactionFavoriteResult, error) {
	// Init query builder
	b := query.From(statement.TransactionFavoriteSchema)

	// Set where
	filters := query.NewFilter(params.Filters, favoriteTransactionFilters)
	b.Where(
		query.And(
			query.Equal(query.Column("customerId", query.Schema(statement.TransactionFavoriteSchema)), query.BindVar()),
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
	countQuery := b.Select(query.Count("*", option.Schema(statement.TransactionFavoriteSchema), option.As("count"))).Build()
	countQuery = rc.conn.Rebind(countQuery)

	// Combine arguments with customerId from filters
	args := append([]interface{}{customerID}, filters.Args()...)

	// Execute query
	var rows []model.TransactionFavorite
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
	result := model.ListTransactionFavoriteResult{
		Rows:  rows,
		Count: count,
	}
	return &result, err
}
