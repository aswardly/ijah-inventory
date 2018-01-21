package handler

import (
	"ijah-inventory/repository/inventory/domain/inventory/service"
	"net/http"
	"strconv"
)

//UpdateSKUHandler is a specific http handler for adding new sku
type UpdateSKUHandler struct {
	Handler
	InventoryService *service.Inventory `inject:"inventoryService"`
}

//UpdateSKUHandle is the implementation of http handler for a UpdateSKUHandler object
func (h *UpdateSKUHandler) UpdateSKUHandle(w http.ResponseWriter, r *http.Request) error {
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

	addErr := h.InventoryService.UpdateSKU(sku, quantityParam, buyPriceParam, sellPriceParam)
	if addErr != nil {
		return composeError(addErr)
	}
	//compose successful response
	response := SimpleResponseStruct{}
	response.Code = ErrCodeSuccessful
	response.Message = "Update successful"

	successfulResponse, statusError := composeJSONResponse(response)
	if statusError != nil {
		return statusError
	}
	//else no problem in json marshalling the response
	w.Write([]byte(successfulResponse))
	return nil
}

//StartUp allows the handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (h *UpdateSKUHandler) StartUp() {
	//Note: perform initialization/bootstrapping here
}

//Shutdown allows the handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (h *UpdateSKUHandler) Shutdown() {
	//Note: perform any cleanup here
}
