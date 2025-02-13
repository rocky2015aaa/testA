package models

type Bakeryt struct {
	Name           string            `json:"name"`
	BakerytAddress string            `json:"bakeryt_address"`
	AddressTo      string            `json:"address_to"`
	InvoiceNo      string            `json:"invoice_no"`
	Date           string            `json:"date"`
	Invoice        []*BakerytInvoice `json:"invoice"`
	SubTotal       float64           `json:"sub_total"`
	Tax            float64           `json:"tax"`
	NetTotal       float64           `json:"net_total"`
	Total          float64           `json:"total"`
	Remarks        string            `json:"remarks"`
}

type BakerytInvoice struct {
	No          int     `json:"no"`
	Code        string  `json:"code"`
	Description string  `json:"description"`
	Qty         float64 `json:"qty"`
	FocQty      float64 `json:"foc_qty"`
	Uom         string  `json:"uom"`
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount"`
	DiscAmount  float64 `json:"disc_amount"`
	SubTotal    float64 `json:"sub_total"`
}
