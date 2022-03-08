package statement

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

var CustomerSchema = schema.New(schema.FromModelRef(model.Customer{}))

type Customer struct {
	Insert             *sqlx.NamedStmt
	UpdateByPhone      *sqlx.NamedStmt
	UpdateByUserRefID  *sqlx.NamedStmt
	FindByRefID        *sqlx.Stmt
	FindByID           *sqlx.Stmt
	FindByPhoneOrCIF   *sqlx.Stmt
	FindByPhone        *sqlx.Stmt
	FindByEmail        *sqlx.Stmt
	FindByEmailOrPhone *sqlx.Stmt
	ReferralCodeExist  *sqlx.Stmt
	UpdateByCIF        *sqlx.NamedStmt
	EmailIsExists      *sqlx.Stmt
	PhoneNumberIsExists *sqlx.Stmt
}

func NewCustomer(db *nsql.DatabaseContext) *Customer {
	// Init query Schema Builder
	sb := query.Schema(CustomerSchema)

	// Init query
	insert := fmt.Sprintf(`%s ON CONFLICT DO NOTHING RETURNING "id"`, sb.Insert())

	findByID := query.Select(query.Column("*")).
		From(CustomerSchema).
		Where(query.Equal(query.Column(CustomerSchema.PrimaryKey()))).
		Build()

	findByRefID := query.Select(query.Column("*")).
		From(CustomerSchema).
		Where(query.Equal(query.Column("userRefId"))).
		Build()

	findByPhone := query.Select(query.Column("*")).
		From(CustomerSchema).
		Where(query.Equal(query.Column("phone"))).
		Build()

	findByPhoneOrCIF := query.Select(query.Column("*")).
		From(CustomerSchema).
		Where(
			query.Or(
				query.Equal(query.Column("cif")),
				query.Equal(query.Column("phone")),
			),
		).
		Build()

	findByEmail := query.Select(query.Column("*")).
		From(CustomerSchema).
		Where(query.Equal(query.Column("email"))).
		Build()

	referralCodeExist := query.Select(query.Column("*")).
		From(CustomerSchema).
		Where(query.Equal(query.Column("referralCode"))).
		Limit(1).
		Build()

	findByEmailOrPhone := query.
		Select(query.Column("*")).
		From(CustomerSchema).
		Where(
			query.Or(
				query.Equal(query.Column("phone")),
				query.Equal(query.Column("email")),
			),
		).
		Build()

	updateByCif := query.
		Update(CustomerSchema, "*").
		Where(query.Equal(query.Column("cif"))).
		Build()

	updateByUserRefID := query.
		Update(CustomerSchema, "*").
		Where(query.Equal(query.Column("userRefId"))).
		Build()

	updateByPhone := query.
		Update(CustomerSchema, "*").
		Where(query.Equal(query.Column("phone"))).
		Build()

	emailIsExists := query.
		Select(query.GreaterThan(query.Count("id"), query.IntVar(0), option.As("isExists"))).
		From(CustomerSchema).
		Where(query.Equal(query.Column("email"))).
		Build()

	phoneNumberIsExist := query.
		Select(query.GreaterThan(query.Count("id"), query.IntVar(0), option.As("isExists"))).
		Where(query.Equal(query.Column("phone"))).
		From(CustomerSchema).
		Build()

	return &Customer{
		Insert:             db.PrepareNamedFmtRebind(insert),
		FindByID:           db.PrepareFmtRebind(findByID),
		FindByRefID:        db.PrepareFmtRebind(findByRefID),
		FindByPhone:        db.PrepareFmtRebind(findByPhone),
		FindByPhoneOrCIF:   db.PrepareFmtRebind(findByPhoneOrCIF),
		FindByEmail:        db.PrepareFmtRebind(findByEmail),
		ReferralCodeExist:  db.PrepareFmtRebind(referralCodeExist),
		FindByEmailOrPhone: db.PrepareFmtRebind(findByEmailOrPhone),
		UpdateByCIF:        db.PrepareNamedFmtRebind(updateByCif),
		UpdateByUserRefID:  db.PrepareNamedFmtRebind(updateByUserRefID),
		UpdateByPhone:      db.PrepareNamedFmtRebind(updateByPhone),
		EmailIsExists:      db.PrepareFmtRebind(emailIsExists),
		PhoneNumberIsExists: db.PrepareFmtRebind(phoneNumberIsExist),
	}
}
