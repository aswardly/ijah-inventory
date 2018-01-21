package handler

import (
	"ijah-inventory/repository/inventory/domain/inventory/service"
	"net/http"
)

//GetItemInfoHandler is a specific http handler for creating sale
type GetItemInfoHandler struct {
	Handler
	InventoryService *service.Inventory `inject:"inventoryService"`
}

//GetItemInfoHandle is the implementation of http handler for a ActivateHandler object
func (h *GetItemInfoHandler) GetItemInfoHandle(w http.ResponseWriter, r *http.Request) error {
	//read the following GET data:
	// - sku
	sku := r.URL.Query().Get("sku")
	stockObj, err := h.InventoryService.GetItemInfo(sku)
	if err != nil {
		//compose failed response
		return composeError(err)
	}
	//compose successful response
	response := SimpleResponseStruct{}
	response.Code = ErrCodeSuccessful
	response.Message = "Inquiry successful"
	response.Data = stockObj

	successfulResponse, statusError := composeJSONResponse(response)
	if statusError != nil {
		return statusError
	}
	//else no problem in json marshalling the response
	w.Write([]byte(successfulResponse))
	return nil
}

//StartUp allows the handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (h *GetItemInfoHandler) StartUp() {
	//Note: perform initialization/bootstrapping here
}

//Shutdown allows the handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (h *GetItemInfoHandler) Shutdown() {
	//Note: perform any cleanup here
}
