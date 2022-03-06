package customer

type Verification struct {
	*Handler
}

func NewVerification(h *Handler) *Verification {
	return &Verification{h}
}
