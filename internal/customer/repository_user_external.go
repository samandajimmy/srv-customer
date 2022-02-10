package customer

import (
	"fmt"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
)

func (rc *RepositoryContext) FindUserExternalByEmailOrPhone(email string) (*model.User, error) {
	var row model.User
	err := rc.stmt.User.FindByEmailOrPhone.GetContext(rc.ctx, &row, email, email)
	return &row, err
}

func (rc *RepositoryContext) FindUserExternalAddressByCustomerID(id int64) (*model.AddressExternal, error) {

	from := `user`
	columns := `user.user_AIID, user.alamat, user.kodepos, kel.nama_kelurahan as kelurahan, kec.nama_kecamatan as kecamatan, kab.nama_kabupaten as kabupaten, prov.nama_provinsi as provinsi, `
	columns += `kel.id as id_kelurahan, kec.id as id_kecamatan, kab.id as id_kabupaten, prov.id as id_provinsi`
	joinTables := `LEFT JOIN ref_kelurahan kel ON user.id_kelurahan=kel.kode_kelurahan `
	joinTables += `LEFT JOIN ref_kecamatan kec ON kel.kode_kecamatan=kec.kode_kecamatan `
	joinTables += `LEFT JOIN ref_kabupaten kab ON kec.kode_kabupaten=kab.kode_kabupaten `
	joinTables += `LEFT JOIN ref_provinsi prov ON kab.kode_provinsi=prov.kode_provinsi `

	q := fmt.Sprintf("SELECT %s FROM %s %s WHERE user.user_AIID = %v", columns, from, joinTables, id)
	// Execute query
	q = rc.conn.Rebind(q)
	var row []model.AddressExternal
	err := rc.conn.SelectContext(rc.ctx, &row, q)
	if err != nil {
		return nil, ncore.TraceError("error when find user external", err)
	}
	return &row[0], nil
}
