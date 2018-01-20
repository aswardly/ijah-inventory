//service_test provides unit tests for user domain service layer
package service_test

import (
	"ijah-inventory/repository/inventory/domain/inventory/model"
	"time"

	"github.com/go-errors/errors"

	"reflect"
)

//Mock object for stock datamapper
type MockStockMapper struct {
}

func (m *MockStockMapper) FindByID(id string) (model.Model, *errors.Error) {
	return &model.Stock{
		Sku:       "dummySku",
		Name:      "dummyItem",
		BuyPrice:  50000,
		SellPrice: 55000,
		Quantity:  250,
	}, nil
}

func (m *MockStockMapper) FindAll() ([]model.Model, *errors.Error) {
	var modelSlice []model.Model
	modelSlice = append(modelSlice, &model.Stock{
		Sku:       "dummySku",
		Name:      "dummyItem",
		BuyPrice:  50000,
		SellPrice: 55000,
		Quantity:  250,
	})
	modelSlice = append(modelSlice, &model.Stock{
		Sku:       "dummySku2",
		Name:      "dummyItem2",
		BuyPrice:  60000,
		SellPrice: 63000,
		Quantity:  120,
	})
	return modelSlice, nil
}

func (m *MockStockMapper) Insert(model model.Model) *errors.Error {
	return nil
}

func (m *MockStockMapper) Update(model model.Model) *errors.Error {
	return nil
}

func (m *MockStockMapper) Delete(model model.Model) *errors.Error {
	return nil
}

func (m *MockStockMapper) Save(model model.Model) *errors.Error {
	return nil
}

//Mock object for purchase datamapper
type MockPurchaseMapper struct {
}

func (m *MockPurchaseMapper) FindByID(id string) (model.Model, *errors.Error) {
	modelObj := &model.Purchase{
		PurchaseID: "dummyPurchaseId",
		Date:       time.Now(),
		Status:     model.PurchaseStatusDraft,
	}
	items := make(map[string]*model.PurchaseItem, 2)
	items["dummySku"] = &model.PurchaseItem{
		Sku:      "dummySku",
		Quantity: 12,
		BuyPrice: 48000,
		Note:     "dummySku purchase",
	}
	items["dummySku2"] = &model.PurchaseItem{
		Sku:      "dummySku2",
		Quantity: 5,
		BuyPrice: 62000,
		Note:     "dummySku2 purchase",
	}
	modelObj.Items = items
	return modelObj, nil
}

func (m *MockPurchaseMapper) FindAll() ([]model.Model, *errors.Error) {
	var modelSlice []model.Model

	modelObj := &model.Purchase{
		PurchaseID: "dummyPurchaseIdX",
		Date:       time.Now(),
		Status:     model.PurchaseStatusDone,
	}
	items := make(map[string]*model.PurchaseItem, 2)
	items["dummySku"] = &model.PurchaseItem{
		Sku:      "dummySku",
		Quantity: 10,
		BuyPrice: 49000,
		Note:     "dummySku purchase X",
	}
	items["dummySku2"] = &model.PurchaseItem{
		Sku:      "dummySku2",
		Quantity: 2,
		BuyPrice: 63000,
		Note:     "dummySku2 purchase X",
	}
	modelObj.Items = items
	modelSlice = append(modelSlice, modelObj)
	anotherModelObj, _ := m.FindByID("not used")
	modelSlice = append(modelSlice, anotherModelObj)
	return modelSlice, nil
}

func (m *MockPurchaseMapper) Insert(model model.Model) *errors.Error {
	return nil
}

func (m *MockPurchaseMapper) Update(model model.Model) *errors.Error {
	return nil
}

func (m *MockPurchaseMapper) Delete(model model.Model) *errors.Error {
	return nil
}

func (m *MockPurchaseMapper) Save(model model.Model) *errors.Error {
	return nil
}

//Mock object for sales datamapper
type MockSalesMapper struct {
}

func (m *MockSalesMapper) FindByID(id string) (model.Model, *errors.Error) {
	modelObj := &model.Sales{
		InvoiceID: "dummyInvoice",
		Date:      time.Now(),
		Status:    model.SalesStatusDraft,
		Note:      "Dummy invoice",
	}
	items := make(map[string]*model.SaleItem, 2)
	items["dummySku"] = &model.SaleItem{
		Sku:       "dummySku",
		Quantity:  3,
		BuyPrice:  55000,
		SellPrice: 60000,
	}
	modelObj.Items = items
	return modelObj, nil
}

func (m *MockSalesMapper) FindAll() ([]model.Model, *errors.Error) {
	var modelSlice []model.Model

	modelObj := &model.Sales{
		InvoiceID: "dummyInvoice",
		Date:      time.Now(),
		Status:    model.SalesStatusDraft,
		Note:      "Dummy invoice",
	}
	items := make(map[string]*model.SaleItem, 2)
	items["dummySku"] = &model.SaleItem{
		Sku:       "dummySku",
		Quantity:  3,
		BuyPrice:  55000,
		SellPrice: 60000,
	}
	modelObj.Items = items

	modelSlice = append(modelSlice, modelObj)

	anotherModelObj, _ := m.FindByID("not used")
	modelSlice = append(modelSlice, anotherModelObj)

	return modelSlice, nil
}

//getType is a function to get type of something (without package name)
//see: https://stackoverflow.com/questions/35790935/using-reflection-in-go-to-get-the-name-of-a-struct
func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
