package repository

import (
	"fmt"

	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/repository/statement"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nsql"
)

type UserExternal struct {
	db   *nsql.DB
	stmt *statement.UserStatement
}

func (a *UserExternal) HasInitialized() bool {
	return true
}

func (a *UserExternal) Init(dataSources DataSourceMap, _ contract.RepositoryMap) error {
	a.db = dataSources.DBExternal
	a.stmt = statement.NewUserStatement(a.db)
	return nil
}

func (a *UserExternal) FindByEmailOrPhone(email string) (*model.User, error) {
	var row model.User
	err := a.stmt.FindByEmailOrPhone.Get(&row, email, email)
	return &row, err
}

func (a *UserExternal) FindAddressByCustomerId(id int64) (*model.AddressExternal, error) {

	from := `user`
	columns := `user.user_AIID, user.alamat, user.kodepos, kel.nama_kelurahan as kelurahan, kec.nama_kecamatan as kecamatan, kab.nama_kabupaten as kabupaten, prov.nama_provinsi as provinsi, `
	columns += `kel.id as id_kelurahan, kec.id as id_kecamatan, kab.id as id_kabupaten, prov.id as id_provinsi`
	joinTables := `LEFT JOIN ref_kelurahan kel ON user.id_kelurahan=kel.kode_kelurahan `
	joinTables += `LEFT JOIN ref_kecamatan kec ON kel.kode_kecamatan=kec.kode_kecamatan `
	joinTables += `LEFT JOIN ref_kabupaten kab ON kec.kode_kabupaten=kab.kode_kabupaten `
	joinTables += `LEFT JOIN ref_provinsi prov ON kab.kode_provinsi=prov.kode_provinsi `

	q := fmt.Sprintf("SELECT %s FROM %s %s WHERE user.user_AIID = %v", columns, from, joinTables, id)
	// Execute query
	q = a.db.Conn.Rebind(q)
	var row []model.AddressExternal
	err := a.db.Conn.Select(&row, q)
	if err != nil {
		return nil, ncore.TraceError(err)
	}
	return &row[0], nil
}
