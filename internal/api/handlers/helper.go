package handlers

import (
	"bytes"
	"context"
	"os/exec"
	"strconv"
	"strings"

	"github.com/shmoulana/hashocr/internal/api/models"
	"github.com/shmoulana/hashocr/internal/pkg/utils"
)

func getTextFromfile(ctx context.Context, fileName string) ([]string, error) {
	// See "man pdftotext" for more options.
	args := []string{
		"-layout",  // Maintain (as best as possible) the original physical layout of the text.
		"-nopgbrk", // Don't insert page breaks (form feed characters) between pages.
		fileName,
		"-", // Send the output to stdout.
	}
	cmd := exec.CommandContext(ctx, "pdftotext", args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return strings.Split(buf.String(), "\n"), nil
}

func parseTaxInvoice(pdfText []string) (interface{}, error) {
	tempExtraData := []string{}
	tempInvoice := []string{}
	// // Filter out empty lines and trim whitespace
	isInvoice := false
	for idx, line := range pdfText {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			if idx == 11 {
				isInvoice = true
			}
			if strings.HasPrefix(trimmedLine, "Sales Qty") {
				isInvoice = false
			}
			if isInvoice {
				tempInvoice = append(tempInvoice, trimmedLine)
			} else {
				tempExtraData = append(tempExtraData, trimmedLine)
			}
		}
	}

	invoiceData := [][]string{}
	itemCount := 0
	hasHsn := false
	for _, line := range tempInvoice[1:] {
		if len(line) > 50 {
			items := utils.SplitByMoreThanNSpaces(line, 2)
			if strings.Contains(items[3], "hsn_identifier") {
				extendedItems := make([]string, 0, len(line)+1)
				extendedItems = append(extendedItems, items[:3]...)
				extendedItems = append(extendedItems, "")
				extendedItems = append(extendedItems, items[3:]...)
				invoiceData = append(invoiceData, extendedItems)
				hasHsn = true
			} else {
				invoiceData = append(invoiceData, items)
			}
			itemCount++
		} else {
			invoiceData[itemCount-1][1] += " " + strings.TrimSpace(line)
		}
	}
	invoices := []*models.TaxInvoice{}
	for _, item := range invoiceData {
		no, _ := strconv.Atoi(strings.TrimSpace(item[0]))
		price, _ := strconv.Atoi(strings.TrimSpace(item[2]))
		var qty, gst int
		var rate, amount float64
		var qtyRate []string
		if hasHsn {
			qtyRate = strings.Split(item[4], " ")
			if len(qtyRate) > 1 {
				qty, _ = strconv.Atoi(strings.TrimSpace(qtyRate[0]))
				rate, _ = strconv.ParseFloat(strings.TrimSpace(qtyRate[1]), 64)
				gst, _ = strconv.Atoi(strings.TrimSpace(item[5]))
				amount, _ = strconv.ParseFloat(strings.TrimSpace(item[6]), 64)
			} else {
				qty, _ = strconv.Atoi(strings.TrimSpace(item[4]))
				rate, _ = strconv.ParseFloat(item[5], 64)
				gst, _ = strconv.Atoi(strings.TrimSpace(item[6]))
				amount, _ = strconv.ParseFloat(strings.TrimSpace(item[7]), 64)
			}
		} else {
			qtyRate = strings.Split(item[3], " ")
			if len(qtyRate) > 1 {
				qty, _ = strconv.Atoi(strings.TrimSpace(qtyRate[0]))
				rate, _ = strconv.ParseFloat(strings.TrimSpace(qtyRate[1]), 64)
				gst, _ = strconv.Atoi(strings.TrimSpace(item[4]))
				amount, _ = strconv.ParseFloat(strings.TrimSpace(item[5]), 64)
			} else {
				qty, _ = strconv.Atoi(strings.TrimSpace(item[3]))
				rate, _ = strconv.ParseFloat(item[4], 64)
				gst, _ = strconv.Atoi(strings.TrimSpace(item[5]))
				amount, _ = strconv.ParseFloat(strings.TrimSpace(item[6]), 64)
			}
		}

		invoice := models.TaxInvoice{
			No:       no,
			ItemName: strings.TrimSpace(item[1]),
			Price:    price,
			Qty:      qty,
			Rate:     rate,
			Gst:      gst,
			Amount:   amount,
		}
		if hasHsn {
			invoice.Hsn = strings.TrimSpace(item[3])
		}
		invoices = append(invoices, &invoice)
	}

	extraData := []string{}
	for _, line := range tempExtraData {
		extraData = append(extraData, utils.SplitByMoreThanNSpaces(line, 6)...)
	}
	page, _ := strconv.Atoi(strings.TrimSpace(strings.Split(extraData[2], ":")[1]))
	salesQty, _ := strconv.Atoi(strings.TrimSpace(strings.Split(extraData[9], ":")[1]))
	subTotal, _ := strconv.ParseFloat(strings.TrimSpace(extraData[11]), 64)
	cgst, _ := strconv.ParseFloat(strings.TrimSpace(extraData[14]), 64)
	sgst, _ := strconv.ParseFloat(strings.TrimSpace(extraData[17]), 64)
	netTotal, _ := strconv.ParseFloat(strings.TrimSpace(extraData[23]), 64)
	taxInvoice := models.Tax{
		GstinNo:     strings.TrimSpace(strings.Split(extraData[0], ":")[1]),
		Page:        page,
		Name:        strings.TrimSpace(extraData[3]),
		BillNo:      strings.TrimSpace(strings.Split(extraData[4], ":")[1]),
		Date:        strings.TrimSpace(strings.Split(extraData[5], ":")[1]),
		Time:        strings.TrimSpace(strings.SplitN(extraData[6], ":", 2)[1]),
		PhoneNumber: strings.TrimSpace(strings.Split(extraData[7], ":")[1]),
		Invoice:     invoices,
		SalesQty:    salesQty,
		SubTotal:    subTotal,
		Cgst:        cgst,
		Sgst:        sgst,
		BankName:    strings.TrimSpace(strings.Split(extraData[18], ":")[1]),
		AcNo:        strings.TrimSpace(strings.Split(extraData[19], ":")[1]),
		Ifsc:        strings.TrimSpace(strings.Split(extraData[20], ":")[1]),
		BillType:    strings.TrimSpace(extraData[21]),
		NetTotal:    netTotal,
		For:         strings.TrimSpace(strings.Split(extraData[25], ":")[1]),
		TermsAndConditions: []string{
			strings.TrimSpace(extraData[26]),
			strings.TrimSpace(extraData[27]),
			strings.TrimSpace(extraData[28]),
			strings.TrimSpace(extraData[29]),
			strings.TrimSpace(extraData[30]),
		},
	}

	return taxInvoice, nil
}

func parseBakertyInvoice(pdfText []string) (interface{}, error) {
	extraData := []string{}
	tempInvoice := []string{}
	// Filter out empty lines and trim whitespace
	isInvoice := false
	for idx, line := range pdfText {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			if idx == 12 {
				isInvoice = true
			}
			if strings.HasPrefix(trimmedLine, "Sub Total") {
				isInvoice = false
			}
			if isInvoice {
				tempInvoice = append(tempInvoice, trimmedLine)
			} else {
				extraData = append(extraData, trimmedLine)
			}
		}
	}

	invoiceData := [][]string{}
	itemCount := 0
	for _, line := range tempInvoice[2:] {
		if len(line) > 50 {
			items := utils.SplitByMoreThanNSpaces(line, 2)
			invoiceData = append(invoiceData, items)
			itemCount++
		} else {
			invoiceData[itemCount-1][1] += strings.TrimSpace(line)
		}
	}
	invoices := []*models.BakerytInvoice{}

	for _, item := range invoiceData {
		no, _ := strconv.Atoi(strings.TrimSpace(item[0]))
		qty, _ := strconv.ParseFloat(strings.TrimSpace(item[3]), 64)
		focQty, _ := strconv.ParseFloat(strings.TrimSpace(item[4]), 64)
		price, _ := strconv.ParseFloat(strings.TrimSpace(item[6]), 64)
		discount, _ := strconv.ParseFloat(strings.TrimSpace(item[7]), 64)
		discountAmount, _ := strconv.ParseFloat(strings.TrimSpace(item[8]), 64)
		subTotal, _ := strconv.ParseFloat(strings.TrimSpace(item[9]), 64)

		invoice := models.BakerytInvoice{
			No:          no,
			Code:        strings.TrimSpace(item[1]),
			Description: strings.TrimSpace(item[2]),
			Qty:         qty,
			FocQty:      focQty,
			Uom:         strings.TrimSpace(item[5]),
			Price:       price,
			Discount:    discount,
			DiscAmount:  discountAmount,
			SubTotal:    subTotal,
		}
		invoices = append(invoices, &invoice)
	}

	commonData := utils.SplitByMoreThanNSpaces(extraData[6], 6)
	subTotal, _ := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(strings.Fields(extraData[9])[2], ",", "")), 64)
	tax, _ := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(extraData[10], ",", "")), 64)
	netTotal, _ := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(strings.Fields(extraData[11])[2], ",", "")), 64)
	total, _ := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(strings.Fields(extraData[12])[3], ",", "")), 64)
	bakerytInvoice := models.Bakeryt{
		Name:           strings.TrimSpace(extraData[0]),
		BakerytAddress: strings.TrimSpace(extraData[1]) + " " + strings.TrimSpace(extraData[2]) + " " + strings.TrimSpace(extraData[3]),
		AddressTo:      strings.TrimSpace(commonData[0]),
		InvoiceNo:      strings.TrimSpace(utils.SplitByMoreThanNSpaces(extraData[4], 6)[1]) + strings.TrimSpace(commonData[1]),
		Date:           strings.TrimSpace(strings.Fields(extraData[7])[1]),
		Invoice:        invoices,
		SubTotal:       subTotal,
		Tax:            tax,
		NetTotal:       netTotal,
		Total:          total,
		Remarks:        strings.TrimSpace(extraData[14]),
	}

	return bakerytInvoice, nil
}

func parseWinnersInvoice(pdfText []string) (interface{}, error) {
	extraData := []string{}
	tempInvoice := []string{}
	// Filter out empty lines and trim whitespace
	isInvoice := false
	for idx, line := range pdfText {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			if idx == 11 {
				isInvoice = true
			}
			if strings.HasPrefix(trimmedLine, "Sub Total") {
				isInvoice = false
			}
			if isInvoice {
				tempInvoice = append(tempInvoice, trimmedLine)
			} else {
				extraData = append(extraData, trimmedLine)
			}
		}
	}

	invoiceData := [][]string{}
	itemCount := 0
	for _, line := range tempInvoice[2:] {
		if len(line) > 50 {
			items := utils.SplitByMoreThanNSpaces(line, 2)
			invoiceData = append(invoiceData, items)
			itemCount++
		} else {
			invoiceData[itemCount-1][2] += " " + strings.TrimSpace(line)
		}
	}
	invoices := []*models.CommonInvoice{}

	for _, item := range invoiceData {
		no, _ := strconv.Atoi(strings.TrimSpace(item[0]))
		qty, _ := strconv.ParseFloat(strings.TrimSpace(item[3]), 64)
		focQty, _ := strconv.ParseFloat(strings.TrimSpace(item[4]), 64)
		price, _ := strconv.ParseFloat(strings.TrimSpace(item[6]), 64)
		discount, _ := strconv.ParseFloat(strings.TrimSpace(item[7]), 64)
		discountAmount, _ := strconv.ParseFloat(strings.TrimSpace(item[8]), 64)
		subTotal, _ := strconv.ParseFloat(strings.TrimSpace(item[9]), 64)

		invoice := models.CommonInvoice{
			No:          no,
			Barcode:     strings.TrimSpace(item[1]),
			ProductName: strings.TrimSpace(item[2]),
			Qty:         qty,
			FocQty:      focQty,
			Uom:         strings.TrimSpace(item[5]),
			Price:       price,
			Discount:    discount,
			DiscAmount:  discountAmount,
			SubTotal:    subTotal,
		}
		invoices = append(invoices, &invoice)
	}

	vendorInvoiceNo, _ := strconv.Atoi(strings.TrimSpace(strings.Split(extraData[3], ":")[1]))
	subTotal, _ := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(strings.Fields(extraData[8])[2], ",", "")), 64)
	totalQty, _ := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(strings.Fields(extraData[9])[2], ",", "")), 64)
	tax, _ := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(strings.Fields(extraData[10])[2], ",", "")), 64)
	netTotal, _ := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(strings.Fields(extraData[11])[2], ",", "")), 64)

	winnersInvoice := models.Winners{
		PurchaseInvoiceNo: strings.TrimSpace(strings.Split(extraData[0], ":")[1]),
		Date:              strings.TrimSpace(strings.Split(extraData[1], ":")[1]),
		Name:              strings.TrimSpace(extraData[2]),
		VendorInvoiceNo:   vendorInvoiceNo,
		Address:           strings.TrimSpace(extraData[4]),
		City:              strings.TrimSpace(strings.Split(extraData[5], ":")[1]),
		AddressTo:         strings.TrimSpace(extraData[7]),
		Invoice:           invoices,
		SubTotal:          subTotal,
		TotalQty:          totalQty,
		Tax:               tax,
		NetTotal:          netTotal,
	}

	return winnersInvoice, nil
}

func parseTambaramInvoice(pdfText []string) (interface{}, error) {
	extraData := []string{}
	tempInvoice := []string{}
	// Filter out empty lines and trim whitespace
	isInvoice := false
	for idx, line := range pdfText {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			if idx == 11 {
				isInvoice = true
			}
			if strings.HasPrefix(trimmedLine, "Sub Total") {
				isInvoice = false
			}
			if isInvoice {
				tempInvoice = append(tempInvoice, trimmedLine)
			} else {
				extraData = append(extraData, trimmedLine)
			}
		}
	}

	invoiceData := [][]string{}
	itemCount := 0
	for _, line := range tempInvoice[2:] {
		if len(line) > 50 {
			items := utils.SplitByMoreThanNSpaces(line, 2)
			invoiceData = append(invoiceData, items)
			itemCount++
		} else {
			invoiceData[itemCount-1][2] += " " + strings.TrimSpace(line)
		}
	}
	invoices := []*models.CommonInvoice{}

	for _, item := range invoiceData {
		no, _ := strconv.Atoi(strings.TrimSpace(item[0]))
		qty, _ := strconv.ParseFloat(strings.TrimSpace(item[3]), 64)
		focQty, _ := strconv.ParseFloat(strings.TrimSpace(item[4]), 64)
		price, _ := strconv.ParseFloat(strings.TrimSpace(item[6]), 64)
		discount, _ := strconv.ParseFloat(strings.TrimSpace(item[7]), 64)
		discountAmount, _ := strconv.ParseFloat(strings.TrimSpace(item[8]), 64)
		subTotal, _ := strconv.ParseFloat(strings.TrimSpace(item[9]), 64)

		invoice := models.CommonInvoice{
			No:          no,
			Barcode:     strings.TrimSpace(item[1]),
			ProductName: strings.TrimSpace(item[2]),
			Qty:         qty,
			FocQty:      focQty,
			Uom:         strings.TrimSpace(item[5]),
			Price:       price,
			Discount:    discount,
			DiscAmount:  discountAmount,
			SubTotal:    subTotal,
		}
		invoices = append(invoices, &invoice)
	}

	vendorInvoiceNo, _ := strconv.Atoi(strings.TrimSpace(strings.Fields(extraData[4])[3]))
	subTotal, _ := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(strings.Fields(extraData[7])[2], ",", "")), 64)
	totalQty, _ := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(strings.Fields(extraData[8])[2], ",", "")), 64)
	tax, _ := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(strings.Fields(extraData[9])[3], ",", "")), 64)
	roundedAmount, _ := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(strings.Fields(extraData[10])[2], ",", "")), 64)
	netTotal, _ := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(strings.Fields(extraData[11])[2], ",", "")), 64)

	tambaramInvoice := models.Tambaram{
		PurchaseInvoiceNo: strings.TrimSpace(strings.Fields(extraData[0])[2]) + strings.TrimSpace(strings.Fields(extraData[1])[1]),
		Date:              strings.TrimSpace(strings.Fields(extraData[2])[1]),
		Name:              strings.TrimSpace(extraData[3]),
		VendorInvoiceNo:   vendorInvoiceNo,
		AddressTo:         strings.TrimSpace(extraData[6]),
		Invoice:           invoices,
		SubTotal:          subTotal,
		TotalQty:          totalQty,
		Tax:               tax,
		RoundedAmount:     roundedAmount,
		NetTotal:          netTotal,
	}

	return tambaramInvoice, nil
}

func parseRPInvoice(pdfText []string) (interface{}, error) {
	extraData := []string{}
	tempInvoice := []string{}
	// Filter out empty lines and trim whitespace
	isInvoice := false
	for idx, line := range pdfText {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			if idx == 14 {
				isInvoice = true
			}
			if strings.HasPrefix(trimmedLine, "Qty:") {
				isInvoice = false
			}
			if isInvoice {
				tempInvoice = append(tempInvoice, trimmedLine)
			} else {
				extraData = append(extraData, trimmedLine)
			}
		}
	}

	invoiceData := [][]string{}
	itemCount := 0
	for _, line := range tempInvoice[1:] {
		if len(line) > 50 {
			items := utils.SplitByMoreThanNSpaces(line, 2)
			invoiceData = append(invoiceData, items)
			itemCount++
		} else {
			invoiceData[itemCount-1][2] += " " + strings.TrimSpace(line)
		}
	}
	invoices := []*models.RPInvoice{}
	for _, item := range invoiceData {
		slno, _ := strconv.Atoi(strings.TrimSpace(item[0]))
		qty, _ := strconv.Atoi(strings.TrimSpace(item[3]))
		cost, _ := strconv.ParseFloat(strings.TrimSpace(item[5]), 64)
		var disc, foc, amount float64
		if len(item) == 9 {
			disc, _ = strconv.ParseFloat(strings.TrimSpace(item[6]), 64)
			foc, _ = strconv.ParseFloat(strings.TrimSpace(item[7]), 64)
			amount, _ = strconv.ParseFloat(strings.TrimSpace(item[8]), 64)
		} else {
			amount, _ = strconv.ParseFloat(strings.TrimSpace(item[6]), 64)
		}

		invoice := models.RPInvoice{
			SlNo:        slno,
			ItemCode:    strings.TrimSpace(item[1]),
			Description: strings.TrimSpace(item[2]),
			Qty:         qty,
			Uom:         strings.TrimSpace(item[4]),
			Cost:        cost,
			Amount:      amount,
		}
		if disc > 0 {
			invoice.Disc = disc
		}
		if foc > 0 {
			invoice.Foc = foc
		}
		invoices = append(invoices, &invoice)
	}

	commonInfo := strings.Fields(extraData[1])
	innerInfoAttention := ""
	innerInfoAttentionVal := strings.Fields(extraData[5])
	if len(innerInfoAttentionVal) > 1 {
		innerInfoAttention = strings.TrimSpace(innerInfoAttentionVal[1])
	}
	commonInfo2 := utils.SplitByMoreThanNSpaces(extraData[8], 5)
	innerInfoPhone := ""
	innerInfoFax := ""
	page := strings.TrimSpace(commonInfo2[3])
	if commonInfo2[1] != "Fax:" {
		innerInfoPhone = strings.TrimSpace(commonInfo2[1])
		if commonInfo2[3] != "Page:" {
			innerInfoFax = strings.TrimSpace(commonInfo2[3])
			page = strings.TrimSpace(commonInfo2[5])
		}
	} else {
		if commonInfo2[2] != "Page:" {
			innerInfoFax = strings.TrimSpace(commonInfo2[2])
			page = strings.TrimSpace(commonInfo2[4])
		}
	}
	innerInfo := models.InnerInfo{
		Name: strings.TrimSpace(extraData[3]),
	}
	if innerInfoAttention != "" {
		innerInfo.Attention = innerInfoAttention
	}
	if innerInfoPhone != "" {
		innerInfo.Phone = innerInfoPhone
	}
	if innerInfoFax != "" {
		innerInfo.Fax = innerInfoFax
	}
	terms := strings.Fields(extraData[7])
	commonInfo3 := strings.Fields(extraData[9])
	qty, _ := strconv.ParseFloat(strings.TrimSpace(commonInfo3[1]), 64)
	total, _ := strconv.ParseFloat(strings.TrimSpace(commonInfo3[3]), 64)
	remarks := ""
	remarksVal := strings.Fields(extraData[10])
	if len(remarksVal) > 1 {
		remarks = strings.TrimSpace(remarksVal[1])
	}
	discount, _ := strconv.ParseFloat(strings.TrimSpace(strings.Fields(extraData[11])[2]), 64)
	cgst, _ := strconv.ParseFloat(strings.TrimSpace(strings.Fields(extraData[12])[2]), 64)
	sgst, _ := strconv.ParseFloat(strings.TrimSpace(strings.Fields(extraData[13])[2]), 64)
	netTotal, _ := strconv.ParseFloat(strings.TrimSpace(strings.Fields(extraData[14])[3]), 64)
	orderBy := utils.SplitByMoreThanNSpaces(extraData[15], 5)
	rpInvoice := models.RP{
		Name:         strings.TrimSpace(extraData[0]),
		Address:      strings.TrimSpace(commonInfo[0]) + " " + strings.TrimSpace(commonInfo[1]),
		Tel:          strings.TrimSpace(commonInfo[3]) + " " + strings.TrimSpace(commonInfo[1]),
		Fax:          strings.TrimSpace(commonInfo[6]),
		InvoiceName:  strings.TrimSpace(commonInfo[7]) + " " + strings.TrimSpace(commonInfo[8]),
		PurchaseNo:   strings.TrimSpace(strings.Fields(extraData[2])[3]),
		Date:         strings.TrimSpace(strings.Fields(extraData[4])[1]),
		DeliveryDate: strings.TrimSpace(strings.Fields(extraData[6])[3]),
		Terms:        strings.TrimSpace(terms[2]) + " " + strings.TrimSpace(terms[3]),
		Page:         page,
		InnerInfo:    &innerInfo,
		Invoice:      invoices,
		Qty:          qty,
		Total:        total,
		Remarks:      remarks,
		Discount:     discount,
		Sgst:         sgst,
		Cgst:         cgst,
		NetTotal:     netTotal,
		OrdererBy:    strings.TrimSpace(orderBy[1]),
	}

	return rpInvoice, nil
}
