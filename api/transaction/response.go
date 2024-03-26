package transaction

type SingleBalanceRes struct {
	Amount       int    `json:"balance"`
	CurrencyCode string `json:"currencyCode"`
}

type GetBalanceResponse []SingleBalanceRes

type SingleHistoryRes struct {
	TransactionId     string `json:"transactionId"`
	Balance           int    `json:"balance"`
	CurrencyCode      string `json:"currency"`
	TransferProof     string `json:"transferProofImg"`
	CreatedAt         int64  `json:"createdAt"`
	SourceTransaction `json:"source"`
}

type SourceTransaction struct {
	BankAccountNumber string `json:"bankAccountNumber"`
	BankName          string `json:"bankName"`
}
type GetHistoryResponse []SingleHistoryRes

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}
