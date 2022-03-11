package dto

type Favorite struct {
	XID             string `json:"xid"`
	Type            string `json:"type"`
	TypeTransaction string `json:"typeTransaction"`
	AccountName     string `json:"accountName"`
	AccountNumber   string `json:"accountNumber"`
	BankName        string `json:"bankName"`
	BankCode        string `json:"bankCode"`
	GroupMPO        string `json:"groupMpo"`
	ServiceCodeMPO  string `json:"serviceCodeMpo"`
	*BaseField
}

type CreateFavoritePayload struct {
	Subject         *Subject `json:"-"`
	RequestID       string   `json:"-"`
	UserRefID       string   `json:"-"`
	Type            string   `json:"type"`
	TypeTransaction string   `json:"typeTransaction"`
	AccountName     string   `json:"accountName"`
	AccountNumber   string   `json:"accountNumber"`
	BankName        string   `json:"bankName"`
	BankCode        string   `json:"bankCode"`
	GroupMPO        string   `json:"groupMpo"`
	ServiceCodeMPO  string   `json:"serviceCodeMpo"`
}

type ListFavoriteResult struct {
	Rows     []Favorite    `json:"rows"`
	Metadata *ListMetadata `json:"metadata"`
}

type GetDetailFavoritePayload struct {
	RequestID string `json:"-"`
	UserRefID string `json:"-"`
	XID       string `json:"xid"`
}
