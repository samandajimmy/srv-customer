package customer

import "repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"

// Tabungan Emas Service

func (s *Service) validateCIF(cif string, id int64) (bool, error) {
	c, err := s.repo.FindCustomerByRefID(id)
	if err != nil {
		return false, err
	}

	if cif != c.Cif {
		return false, constant.DefaultError.Trace()
	}

	return true, nil
}
