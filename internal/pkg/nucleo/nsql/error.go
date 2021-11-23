package nsql

var RowNotUpdatedError = new(rowNotUpdatedError)

type rowNotUpdatedError struct{}

func (n *rowNotUpdatedError) Error() string {
	return "nsql: row is not updated"
}
