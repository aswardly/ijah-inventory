package handler

import (
	"ijah-inventory/repository/inventory/domain/inventory/service"
	"net/http"
	"regexp"
	"strconv"
)

//CreateSaleHandler is a specific http handler for creating sale
type CreateSaleHandler struct {
	Handler
	InventoryService *service.Inventory `inject:"inventoryService"`
}

//CreateSaleHandle is the implementation of http handler for a CreateSaleHandler object
func (h *CreateSaleHandler) CreateSaleHandle(w http.ResponseWriter, r *http.Request) error {
	//read the following POST data:
	// - invoiceNo
	// - note
	//repeating items
	// - sku[x]
	// - quantity[x]

	//errForm := r.ParseMultipartForm(131072) //128kb memory max
	errForm := r.ParseForm()
	if errForm != nil {
		return composeError(errForm)
	}

	var invoiceID, note string
	var itemsSku, itemsQuantity map[string]string

	itemsSku = make(map[string]string, 0)
	itemsQuantity = make(map[string]string, 0)

	//regex for parsing items in form post data
	skuRegxp := regexp.MustCompile(`sku\[(?P<sku>\d+)\]`)
	quantityRegxp := regexp.MustCompile(`quantity\[(?P<quantity>\d+)\]`)

	//parse through all post data
	for key, val := range r.PostForm {
		if key == "invoiceId" {
			invoiceID = val[0]
		}
		if key == "note" {
			note = val[0]
		}
		skuFound := skuRegxp.FindStringSubmatch(key)
		if len(skuFound) > 0 {
			//found "sku[x]" pattern in post data
			itemsSku[skuFound[1]] = val[0]
		}
		quantityFound := quantityRegxp.FindStringSubmatch(key)
		if len(quantityFound) > 0 {
			//found "quantity[x]" pattern in post data
			itemsQuantity[quantityFound[1]] = val[0]
		}
	}

	//parse obtained sku and quantity
	var saleItemSlice []service.SaleItem
	saleItemSlice = make([]service.SaleItem, 0)

	for skuKey, skuVal := range itemsSku {
		if itemsQuantity[skuKey] != "" {
			theQuantity, err := strconv.ParseInt(itemsQuantity[skuKey], 10, 64)
			if err == nil {
				newSaleItem := service.SaleItem{
					Sku:      skuVal,
					Quantity: theQuantity,
				}
				saleItemSlice = append(saleItemSlice, newSaleItem)
			}
		}
	}

	_, errc := h.InventoryService.CreateSale(invoiceID, note, saleItemSlice)
	if errc != nil {
		return composeError(errc)
	}

	//compose successful response
	response := SimpleResponseStruct{}
	response.Code = ErrCodeSuccessful
	response.Message = "Sale created successfully"

	successfulResponse, statusError := composeJSONResponse(response)
	if statusError != nil {
		return statusError
	}
	//else no problem in json marshalling the response
	w.Write([]byte(successfulResponse))
	return nil
}

//StartUp allows the handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (h *CreateSaleHandler) StartUp() {
	//TODO: perform initialization/bootstrapping here
}

//Shutdown allows the handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (h *CreateSaleHandler) Shutdown() {
	//TODO: perform any cleanup here
}
