package customer

import (
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
)

func newEqualFilter(s *schema.Schema, col string) nsql.FilterParser {
	return func(qv string) (nsql.WhereWriter, []interface{}) {
		w := query.Equal(query.Column(col, option.Schema(s)))
		return w, []interface{}{qv}
	}
}
