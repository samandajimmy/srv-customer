package customer

import "github.com/nbs-go/nlogger/v2"

var log nlogger.Logger

func init() {
	log = nlogger.Get().NewChild(nlogger.WithNamespace("pds/customer"))
}
