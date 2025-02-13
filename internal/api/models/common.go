package models

type CommonInvoice struct {
	No          int     `json:"no"`
	Barcode     string  `json:"barcode"`
	ProductName string  `json:"product_name"`
	Qty         float64 `json:"qty"`
	FocQty      float64 `json:"foc_qty"`
	Uom         string  `json:"uom"`
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount"`
	DiscAmount  float64 `json:"disc_amount"`
	SubTotal    float64 `json:"subtotal"`
}
