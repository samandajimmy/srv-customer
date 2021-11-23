package contract

var AliasDisabledError = new(aliasDisabledError)

type aliasDisabledError struct{}

func (e *aliasDisabledError) Error() string {
	return "alias path is disabled"
}
