package nhttp

/// HookFunc represents function that will be called after handler.
/// The function does not implements http.Handler and is not
type HookFunc func(r *Request) error

func ChainHooks(fnArr ...HookFunc) HookFunc {
	return func(r *Request) error {
		var err error

		for _, fn := range fnArr {
			// Call function
			err = fn(r)

			// If error occurred, then break hook chain
			if err != nil {
				return err
			}
		}

		return err
	}
}
