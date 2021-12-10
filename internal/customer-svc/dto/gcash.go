package dto

type GCash struct {
	TotalSaldo  int       `json:"totalSaldo"`
	Va          []GCashVa `json:"va"`
	VaAvailable []string  `json:"vaAvailable"`
}

type GCashVa struct {
	ID             string `json:"id"`
	UserAIID       string `json:"user_AIID"`
	Amount         string `json:"amount"`
	KodeBank       string `json:"kodeBank"`
	TrxID          string `json:"trxId"`
	TglExpired     string `json:"tglExpired"`
	VirtualAccount string `json:"virtualAccount"`
	VaNumber       string `json:"vaNumber"`
	CreatedAt      string `json:"createdAt"`
	LastUpdate     string `json:"lastUpdate"`
	NamaBank       string `json:"namaBank"`
	Thumbnail      string `json:"thumbnail"`
}
