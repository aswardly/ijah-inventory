package handler

import (
	"ijah-inventory/repository/inventory/domain/inventory/service"
	"net/http"
	"strconv"
)

//AddSKUHandler is a specific http handler for adding new sku
type AddSKUHandler struct {
	Handler
	InventoryService *service.Inventory `inject:"inventoryService"`
}

//AddSKUHandle is the implementation of http handler for a AddSKUHandler object
func (h *AddSKUHandler) AddSKUHandle(w http.ResponseWriter, r *http.Request) error {
	//read the following POST data:
	// - sku
	// - quantity
	// - buyPrice
	// - sellPrice
	sku := r.PostFormValue("sku")
	quantity := r.PostFormValue("quantity")
	buyPrice := r.PostFormValue("buyPrice")
	sellPrice := r.PostFormValue("sellPrice")

	quantityParam, err := strconv.ParseInt(quantity, 10, 64)
	if err != nil {
		return composeError(err)
	}
	buyPriceParam, err := strconv.ParseFloat(buyPrice, 64)
	if err != nil {
		return composeError(err)
	}
	sellPriceParam, err := strconv.ParseFloat(sellPrice, 64)
	if err != nil {
		return composeError(err)
	}

	err = h.InventoryService.AddSKU(sku, quantityParam, buyPriceParam, sellPriceParam)
	if err != nil {
		//compose failed response
		response := SimpleResponseStruct{}
		response.Code = ErrCodeFailed
		response.Message = "Error: " + err.Error()
		statusErr := composeJSONError(response)
		return statusErr
	}
	//compose successful response
	response := SimpleResponseStruct{}
	response.Code = ErrCodeSuccessful
	response.Message = "Addition successful"

	successfulResponse, statusError := composeJSONResponse(response)
	if statusError != nil {
		return statusError
	}
	//else no problem in json marshalling the response
	w.Write([]byte(successfulResponse))
	return nil
}

//StartUp allows the handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (h *AddSKUHandler) StartUp() {
	//Note: perform initialization/bootstrapping here
}

//Shutdown allows the handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (h *AddSKUHandler) Shutdown() {
	//Note: perform any cleanup here
}
