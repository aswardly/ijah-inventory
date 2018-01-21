package handler

import (
	"fmt"
	"ijah-inventory/repository/inventory/domain/inventory/service"

	"net/http"
	"time"
)

//GetAllSalesValueHandler is a specific http handler for creating sale
type GetAllSalesValueHandler struct {
	Handler
	InventoryService *service.Inventory `inject:"inventoryService"`
}

//inputDateLayout is the layout used for parsing input as time.Time object
const inputDateLayout = "2006-01-02"

//GetAllSalesValueHandle is the implementation of http handler for a GetAllSalesValueHandler object
func (h *GetAllSalesValueHandler) GetAllSalesValueHandle(w http.ResponseWriter, r *http.Request) error {

	//read the following GET data:
	// - starttime
	// - endTime
	startTime := r.URL.Query().Get("startTime")
	endTime := r.URL.Query().Get("endTime")

	startTimeObj, err := time.Parse(inputDateLayout, startTime)
	if err != nil {
		return composeError(fmt.Errorf("startTime format is invalid (should be YYYY-MM-DD)"))
	}
	endTimeObj, err := time.Parse(inputDateLayout, endTime)
	if err != nil {
		return composeError(fmt.Errorf("endTime format is invalid (should be YYYY-MM-DD)"))
	}

	saleValueObj, errs := h.InventoryService.GetAllSalesValue(startTimeObj, endTimeObj)
	fmt.Printf("saleObj is %+v\n", saleValueObj)
	fmt.Printf("err is %+v\n", errs)

	if errs != nil {
		return composeError(errs)
	}
	//compose successful response
	response := SimpleResponseStruct{}
	response.Code = ErrCodeSuccessful
	response.Message = "Inquiry successful"
	response.Data = saleValueObj

	successfulResponse, statusError := composeJSONResponse(response)
	if statusError != nil {
		return statusError
	}
	//else no problem in json marshalling the response
	w.Write([]byte(successfulResponse))
	return nil
}

//StartUp allows the handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (h *GetAllSalesValueHandler) StartUp() {
	//Note: perform initialization/bootstrapping here
}

//Shutdown allows the handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (h *GetAllSalesValueHandler) Shutdown() {
	//Note: perform any cleanup here
}
