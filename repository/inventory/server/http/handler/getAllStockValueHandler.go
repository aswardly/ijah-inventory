package handler

import (
	"ijah-inventory/repository/inventory/domain/inventory/service"
	"net/http"
)

//GetAllStockValueHandler is a specific http handler for creating sale
type GetAllStockValueHandler struct {
	Handler
	InventoryService *service.Inventory `inject:"inventoryService"`
}

//GetAllStockValueHandle is the implementation of http handler for a GetAllStockValueHandler object
func (h *GetAllStockValueHandler) GetAllStockValueHandle(w http.ResponseWriter, r *http.Request) error {

	stockValueObj, err := h.InventoryService.GetAllStockValue()
	if err != nil {
		//compose failed response
		return composeError(err)
	}
	//compose successful response
	response := SimpleResponseStruct{}
	response.Code = ErrCodeSuccessful
	response.Message = "Inquiry successful"
	response.Data = stockValueObj

	successfulResponse, statusError := composeJSONResponse(response)
	if statusError != nil {
		return statusError
	}
	//else no problem in json marshalling the response
	w.Write([]byte(successfulResponse))
	return nil
}

//StartUp allows the handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (h *GetAllStockValueHandler) StartUp() {
	//Note: perform initialization/bootstrapping here
}

//Shutdown allows the handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (h *GetAllStockValueHandler) Shutdown() {
	//Note: perform any cleanup here
}
