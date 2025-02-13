package models

type Winners struct {
	PurchaseInvoiceNo string           `json:"purchase_invoice_no"`
	Date              string           `json:"date"`
	Name              string           `json:"name"`
	VendorInvoiceNo   int              `json:"vendor_invoice_no"`
	Address           string           `json:"address"`
	City              string           `json:"city"`
	AddressTo         string           `json:"address_to"`
	Invoice           []*CommonInvoice `json:"invoice"`
	SubTotal          float64          `json:"sub_total"`
	TotalQty          float64          `json:"total_qty"`
	Tax               float64          `json:"tax"`
	NetTotal          float64          `json:"net_total"`
}
