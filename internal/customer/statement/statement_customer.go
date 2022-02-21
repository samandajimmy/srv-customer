package statement

import (
	"fmt"
	"github.com/jmoiron/sqlx"
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
	FindByRefId        *sqlx.Stmt
	FindById           *sqlx.Stmt
	FindByPhoneOrCIF   *sqlx.Stmt
	FindByPhone        *sqlx.Stmt
	FindByEmail        *sqlx.Stmt
	FindByEmailOrPhone *sqlx.Stmt
	ReferralCodeExist  *sqlx.Stmt
	UpdateByCIF        *sqlx.NamedStmt
}

func NewCustomer(db *nsql.DatabaseContext) *Customer {
	// Init query Schema Builder
	bs := query.Schema(CredentialSchema)

	// Init query
	insert := fmt.Sprintf(`%s ON CONFLICT DO NOTHING RETURNING "id"`, bs.Insert())

	findById := query.Select(query.Column("*")).
		From(CustomerSchema).
		Where(query.Equal(query.Column(CustomerSchema.PrimaryKey()))).
		Build()

	findByRefId := query.Select(query.Column("*")).
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

	updateByUserRefId := query.
		Update(CustomerSchema, "*").
		Where(query.Equal(query.Column("userRefId"))).
		Build()

	updateByPhone := query.
		Update(CustomerSchema, "*").
		Where(query.Equal(query.Column("phone"))).
		Build()

	return &Customer{
		Insert:             db.PrepareNamedFmtRebind(insert),
		FindById:           db.PrepareFmtRebind(findById),
		FindByRefId:        db.PrepareFmtRebind(findByRefId),
		FindByPhone:        db.PrepareFmtRebind(findByPhone),
		FindByPhoneOrCIF:   db.PrepareFmtRebind(findByPhoneOrCIF),
		FindByEmail:        db.PrepareFmtRebind(findByEmail),
		ReferralCodeExist:  db.PrepareFmtRebind(referralCodeExist),
		FindByEmailOrPhone: db.PrepareFmtRebind(findByEmailOrPhone),
		UpdateByCIF:        db.PrepareNamedFmtRebind(updateByCif),
		UpdateByUserRefID:  db.PrepareNamedFmtRebind(updateByUserRefId),
		UpdateByPhone:      db.PrepareNamedFmtRebind(updateByPhone),
	}
}
