package pq

type Operator = string

const (
	EqualOperator            = Operator("=")
	GreaterThanOperator      = Operator(">")
	GreaterThanEqualOperator = Operator(">=")
	LessThanOperator         = Operator("<")
	LessThanEqualOperator    = Operator("<=")
	LikeOperator             = Operator("LIKE")
	ILikeOperator            = Operator("ILIKE")
	InOperator               = Operator("IN")
)

type LikeValueFormat = string

const (
	LikeFormat       = LikeValueFormat("%%%s%%")
	LikeSuffixFormat = LikeValueFormat("%%%s")
	LikePrefixFormat = LikeValueFormat("%s%%")
)
