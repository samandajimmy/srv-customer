package dto

type BankAccount struct {
	XID           string `json:"xid"`
	CustomerID    int64  `json:"customerId"`
	AccountNumber string `json:"accountNumber"`
	AccountName   string `json:"accountName"`
	Bank          *Bank  `json:"bank"`
	*BaseField
}

type Bank struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Code      string `json:"code"`
	Thumbnail string `json:"thumbnail"`
}

type ListBankAccountResult struct {
	Rows     []BankAccount `json:"rows"`
	Metadata *ListMetadata `json:"metadata"`
}

type CreateBankAccountPayload struct {
	Subject       *Subject `json:"-"`
	RequestID     string   `json:"-"`
	UserRefID     string   `json:"-"`
	AccountNumber string   `json:"accountNumber"`
	AccountName   string   `json:"accountName"`
	Bank          *Bank    `json:"bank"`
}

type GetDetailBankAccountPayload struct {
	RequestID string `json:"-"`
	XID       string `json:"xid"`
}

type GetDetailBankAccountResult struct {
	XID           string `json:"xid"`
	AccountNumber string `json:"accountNumber"`
	AccountName   string `json:"accountName"`
	Bank          *Bank  `json:"bank"`
	*BaseField
}

type UpdateBankAccountPayload struct {
	Subject       *Subject `json:"-"`
	RequestID     string   `json:"-"`
	UserRefID     string   `json:"-"`
	XID           string   `json:"xid"`
	AccountNumber string   `json:"accountNumber"`
	AccountName   string   `json:"accountName"`
	Status        int8     `json:"status"`
	Bank          *Bank    `json:"bank"`
	Version       int64    `json:"version"`
}
