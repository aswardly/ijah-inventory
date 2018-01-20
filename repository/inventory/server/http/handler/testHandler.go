package handler

import (
	"fmt"
	"net/http"

	"ijah-inventory/repository/inventory/domain/inventory/datamapper"
	"ijah-inventory/repository/inventory/domain/inventory/model"
)

//TestHandler is a specific http handler for testing
type TestHandler struct {
	Handler
	StockMapper    *datamapper.Stock    `inject:"stockDatamapper"`
	PurchaseMapper *datamapper.Purchase `inject:"purchaseDatamapper"`
}

/*
//NewTestHandler creates a new test handler and returns a pointer to it
func NewTestHandler(sc gocontainer.ServiceContainer) *TestHandler {
	newHandler := new(TestHandler)
	newHandler.SetContainer(sc)
	newHandler.Handle = newHandler.TestHandle

	return newHandler
}
*/

//TestHandle is the implementation of http handler for a TestHandler object
func (th *TestHandler) TestHandle(w http.ResponseWriter, r *http.Request) error {
	fmt.Printf("datamapper: %v\n", th.StockMapper)
	//===========
	stockModel, err := th.StockMapper.FindByID("SSI-D00791015-LL-BWH")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("stock: %v\n", stockModel)
	//===========
	stockModels, err := th.StockMapper.FindAll()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("stock: %T\n", stockModels)
	for key, val := range stockModels {
		fmt.Printf("key: %v, type:%T, val:%v\n", key, val, val)
	}
	//===========
	/*stockModel.Name = "Updated name"
	err = th.StockMapper.Update(stockModel)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	stockModel, err = th.StockMapper.FindByID("SSI-D00791015-LL-BWH")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("updated stock: %v\n", stockModel)
	*/
	//===========
	/*
		err = th.StockMapper.Delete(stockModel)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
	*/
	stockModels, err = th.StockMapper.FindAll()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("after deletion stock: %T\n", stockModels)
	for key, val := range stockModels {
		fmt.Printf("key: %v, type:%T, val:%v\n", key, val, val)
	}

	//===========
	purchaseModel, err := th.PurchaseMapper.FindByID("PO01")
	purchaseModelObj, ok := purchaseModel.(*model.Purchase)
	if false == ok {
		fmt.Printf("Assertion failed!")
	}
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("purchase: %v\n", purchaseModelObj)
	for key, val := range purchaseModelObj.Items {
		fmt.Printf("key: %v, purchaseItem: %v\n", key, val)
	}

	//---------------
	response := SimpleResponseStruct{
		Code:    ErrCodeSuccessful,
		Message: "User successfully activated",
	}
	responseJSON, marshallErr := composeJSONResponse(response)
	if marshallErr != nil {
		return marshallErr
	}
	w.Write([]byte(responseJSON))
	return nil
}

//StartUp allows the handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (th *TestHandler) StartUp() {
	//perform initialization/bootstrapping here
	var stockMapper *datamapper.Stock
	serviceObj, found := th.Sc.GetService("stockDatamapper")
	if false == found {
		panic("service 'stockDatamapper' is not found")
	}
	stockMapper, ok := serviceObj.(*datamapper.Stock)
	if false == ok {
		panic("Failed asserting service 'stockDatamapper'")
	}
	th.StockMapper = stockMapper
}

//Shutdown allows the handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (th *TestHandler) Shutdown() {
	//TODO: perform any cleanup here
}
