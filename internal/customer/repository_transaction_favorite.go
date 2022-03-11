package customer

import (
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/op"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/pq/query"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/statement"
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
