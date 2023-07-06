package models

type Payment struct {
	MerchantID  string `json:"merchant_id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
	CallbackURL string `json:"callback_url"`
	MetaData    []byte `json:"metaData,omitempty"`
	Mobile      string `json:"mobile,omitempty"`
	Email       string `json:"email,omitempty"`
}

type PaymentResponse struct {
	Data   Data     `json:"data"`
	Errors []string `json:"errors"`
}
type Data struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Authority string `json:"authority"`
	FeeType   string `json:"fee_type"`
	Fee       int    `json:"fee"`
}

type AuthorityPair struct {
	ID        uint
	OrderID   uint
	Authority string `gorm:"size:255"`
}

type VerifyRequest struct {
	MerchantID string `json:"merchant_id"`
	Amount     int    `json:"amount"`
	Authority  string `json:"authority"`
}
type VerifyResponse struct {
	Data   VerifyData `json:"data"`
	Errors []string   `json:"errors"`
}
type VerifyData struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	CardHash string `json:"card_hash"`
	CardPan  string `json:"card_pan"`
	RefID    int    `json:"ref_id"`
	FeeType  string `json:"fee_type"`
	Fee      int    `json:"fee"`
}
