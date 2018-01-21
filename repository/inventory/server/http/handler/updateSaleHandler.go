package handler

import (
	"ijah-inventory/repository/inventory/domain/inventory/service"
	"net/http"
)

//UpdateSaleHandler is a specific http handler for creating sale
type UpdateSaleHandler struct {
	Handler
	InventoryService *service.Inventory `inject:"inventoryService"`
}

//UpdateSaleHandle is the implementation of http handler for a CreateSaleHandler object
func (h *UpdateSaleHandler) UpdateSaleHandle(w http.ResponseWriter, r *http.Request) error {
	//read the following POST data:
	// - invoiceId
	// - status
	invoiceNo := r.PostFormValue("invoiceId")
	status := r.PostFormValue("status")

	_, err := h.InventoryService.UpdateSale(invoiceNo, status)
	if err != nil {
		return composeError(err)
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
func (h *UpdateSaleHandler) StartUp() {
	//TODO: perform initialization/bootstrapping here
}

//Shutdown allows the handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (h *UpdateSaleHandler) Shutdown() {
	//TODO: perform any cleanup here
}
