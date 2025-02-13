package models

type RP struct {
	Name         string       `json:"name"`
	Address      string       `json:"address"`
	Tel          string       `json:"tel"`
	Fax          string       `json:"fax"`
	InvoiceName  string       `json:"invoice_name"`
	PurchaseNo   string       `json:"purchase_no"`
	Date         string       `json:"date"`
	DeliveryDate string       `json:"delivery_date"`
	Terms        string       `json:"terms"`
	Page         string       `json:"page"`
	InnerInfo    *InnerInfo   `json:"inner_info"`
	Invoice      []*RPInvoice `json:"invoice"`
	Qty          float64      `json:"qty"`
	Total        float64      `json:"total"`
	Remarks      string       `json:"remarks"`
	Discount     float64      `json:"discount"`
	Sgst         float64      `json:"sgst"`
	Cgst         float64      `json:"cgst"`
	NetTotal     float64      `json:"net_total"`
	OrdererBy    string       `json:"ordered_by"`
}

type InnerInfo struct {
	Name      string `json:"name"`
	Attention string `json:"attention,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Fax       string `json:"fax,omitempty"`
}

type RPInvoice struct {
	SlNo        int     `json:"sl_no"`
	ItemCode    string  `json:"item_code"`
	Description string  `json:"description"`
	Qty         int     `json:"qty"`
	Uom         string  `json:"uom"`
	Cost        float64 `json:"cost"`
	Disc        float64 `json:"disc,omitempty"`
	Foc         float64 `json:"foc,omitempty"`
	Amount      float64 `json:"amount"`
}
