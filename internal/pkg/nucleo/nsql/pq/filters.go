package pq

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"time"
)

func ParseQuery(filters map[string]string, filterKey string, column string, cond Operator) (string, string, bool) {
	fv, ok := filters[filterKey]
	if !ok || fv == "" {
		return "", "", false
	}

	// If condition is In, then return error
	if cond == InOperator || cond == LikeOperator || cond == ILikeOperator {
		return "", "", false
	}

	return fmt.Sprintf(`%s %s ?`, column, cond), fv, true
}

func ParseQueryTime(filters map[string]string, filterKey string, column string, cond Operator, timeLayout string) (string, time.Time, bool) {
	q, strVal, ok := ParseQuery(filters, filterKey, column, cond)
	if !ok {
		return "", time.Time{}, false
	}

	// If string value is equal to now
	if strings.ToLower(strVal) == "now" {
		return q, time.Now(), true
	}

	// If numeric, then parse as epoch
	isEpoch := govalidator.IsNumeric(strVal)
	if isEpoch {
		epoch, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			return "", time.Time{}, false
		}
		return q, time.Unix(epoch, 0), true
	}

	// Else parse with layout
	if timeLayout == "" {
		timeLayout = time.RFC3339
	}
	t, err := time.Parse(timeLayout, strVal)
	if err != nil {
		return "", time.Time{}, false
	}

	return q, t, true
}

func ParseQueryArray(filters map[string]string, filterKey string, column string) (string, []interface{}, error) {
	fv, ok := filters[filterKey]
	if !ok {
		return "", nil, fmt.Errorf("filter not found")
	}

	// Split values by comma
	tmp := strings.Split(fv, ",")

	// Convert to interface args
	values := make([]interface{}, len(tmp))
	for i, v := range tmp {
		values[i] = v
	}

	// Prepare query
	raw := fmt.Sprintf(`%s %s (?)`, column, InOperator)

	// Bind query
	q, args, err := sqlx.In(raw, values)
	if err != nil {
		return "", nil, err
	}

	return q, args, nil
}

func ParseQueryLike(filters map[string]string, filterKey string, column string, cond Operator, valueFmt LikeValueFormat) (string, string, bool) {
	fv, ok := filters[filterKey]
	if !ok || fv == "" {
		return "", "", false
	}

	switch valueFmt {
	case LikeFormat, LikePrefixFormat, LikeSuffixFormat:
		break
	default:
		valueFmt = LikeFormat
	}

	// If condition isn't like or ilike operator, then return error
	switch cond {
	case LikeOperator, ILikeOperator:
		return fmt.Sprintf(`%s %s ?`, column, cond), fmt.Sprintf(valueFmt, fv), true
	default:
		return "", "", false
	}
}
