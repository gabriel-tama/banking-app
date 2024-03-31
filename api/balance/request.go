package balance

import (
	"errors"
	"regexp"
)

type AddBalancePayload struct {
	Sender        string `json:"senderBankAccountNumber" binding:"required,min=5,max=50"`
	SenderBank    string `json:"senderBankName" binding:"required,min=5,max=50"`
	AddBalance    int    `json:"addedBalance" binding:"required"`
	Currency      string `json:"currency" binding:"required" validate:"required,iso4217"`
	TransferProof string `json:"transferProofImg" binding:"required" validate:"required,url"`
}

func (p *AddBalancePayload) Validate() error {
	// Validate TransferProof URL
	if p.TransferProof != "" {
		match, _ := regexp.MatchString(`^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(?:/[^/#?]+)+\.(?:jpg|jpeg|png)$`, p.TransferProof)
		if !match {
			return errors.New("invalid transfer proof URL")
		}
	}

	return nil
}

type ReduceBalancePayload struct {
	Sender     string `json:"recipientBankAccountNumber" binding:"required,min=5,max=50"`
	SenderBank string `json:"recipientBankName" binding:"required,min=5,max=50"`
	AddBalance int    `json:"balances" binding:"required"`
	Currency   string `json:"fromCurrency" binding:"required" validate:"required,iso4217"`
}

type GetHistoryPayload struct {
	Limit  int `form:"limit,default=5" binding:"min=0"`
	Offset int `form:"offset,default=0" binding:"min=0"`
}
