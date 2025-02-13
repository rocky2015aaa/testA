package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	ErrInvoiceTypeNotFound = errors.New("invoice type not found")

	invoiceTypes = map[string]func(pdfText []string) (interface{}, error){
		invoiceType1: parseTaxInvoice,
		invoiceType2: parseBakertyInvoice,
		invoiceType3: parseWinnersInvoice,
		invoiceType4: parseTambaramInvoice,
		invoiceType5: parseRPInvoice,
	}
)

const (
	invoiceType1 = "tax"
	invoiceType2 = "bakerty"
	invoiceType3 = "winners"
	invoiceType4 = "tambaram"
	invoiceType5 = "rp"

	responseSuccess     = "success"
	responseData        = "data"
	responseError       = "error"
	responseDescription = "description"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func getResponse(success bool, data interface{}, err, description string) gin.H {
	return gin.H{
		responseSuccess:     success,
		responseData:        data,
		responseError:       err,
		responseDescription: description,
	}
}

func (h *Handler) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, getResponse(true, nil, "", "pong"))
}

func (h *Handler) CreateJsonFromPdf(ctx *gin.Context) {
	invoiceType := ctx.Query("invoice_type")
	invoiceDataParser, exists := invoiceTypes[invoiceType]
	if !exists {
		err := fmt.Errorf("not a valid invoice type")
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, getResponse(false, nil, err.Error(), err.Error()))
		return
	}
	file, err := ctx.FormFile("file")
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, getResponse(true, nil, err.Error(), "Getting the file has failed"))
		return
	}

	// Log the file details (optional)
	log.Printf("Uploaded File: %+v", file.Filename)
	log.Printf("File Size: %+v", file.Size)
	log.Printf("MIME Header: %+v", file.Header)

	// Define the path where the file will be saved
	savePath := filepath.Join(".", file.Filename)

	// Save the file to the specified path
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, getResponse(true, nil, err.Error(), "Saving the file has failed"))
		return
	}

	lines, err := getTextFromfile(ctx, file.Filename)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, getResponse(true, nil, err.Error(), "Extracting text from the file has failed"))
		return
	}
	if err := os.Remove(savePath); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, getResponse(true, nil, err.Error(), "Remove the file has failed"))
		return
	}
	invoiceData, err := invoiceDataParser(lines)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, getResponse(true, nil, err.Error(), "Getting the invoice data has failed"))
		return
	}
	// Respond with success
	ctx.JSON(http.StatusOK, getResponse(true, invoiceData, "", "Getting the invoice data has succeeded"))
}
