package models

type Tax struct {
	GstinNo            string        `json:"gstin_no"`
	Page               int           `json:"page"`
	Name               string        `json:"name"`
	BillNo             string        `json:"bill_no"`
	Date               string        `json:"date"`
	Time               string        `json:"time"`
	PhoneNumber        string        `json:"phone_number"`
	Invoice            []*TaxInvoice `json:"invoice"`
	SalesQty           int           `json:"sale_qty"`
	SubTotal           float64       `json:"subtotal"`
	Cgst               float64       `json:"cgst"`
	Sgst               float64       `json:"sgst"`
	BankName           string        `json:"bank_name"`
	AcNo               string        `json:"ac_no"`
	Ifsc               string        `json:"ifsc"`
	BillType           string        `json:"bill_type"`
	NetTotal           float64       `json:"net_total"`
	For                string        `json:"for"`
	TermsAndConditions []string      `json:"term_and_conditions"`
}

type TaxInvoice struct {
	No       int     `json:"no"`
	ItemName string  `json:"item_name"`
	Price    int     `json:"price"`
	Hsn      string  `json:"hsn,omitempty"`
	Qty      int     `json:"qty"`
	Rate     float64 `json:"rate"`
	Gst      int     `json:"gst"`
	Amount   float64 `json:"amount"`
}
