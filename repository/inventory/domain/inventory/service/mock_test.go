//service_test provides unit tests for user domain service layer
package service_test

import (
	"ijah-inventory/repository/inventory/domain/inventory/model"
	"time"

	"github.com/go-errors/errors"
)

//Mock object for stock datamapper
type MockStockMapper struct {
}

var dummyStockModel1 = &model.Stock{
	Sku:       "dummySku",
	Name:      "dummyItem",
	BuyPrice:  50000,
	SellPrice: 55000,
	Quantity:  250,
}

var dummyStockModel2 = &model.Stock{
	Sku:       "dummySku2",
	Name:      "dummyItem2",
	BuyPrice:  60000,
	SellPrice: 63000,
	Quantity:  120,
}

func (m *MockStockMapper) FindByID(id string) (model.Model, *errors.Error) {
	return dummyStockModel1, nil
}

func (m *MockStockMapper) FindAll() ([]model.Model, *errors.Error) {
	var stockSlice []model.Model
	stockSlice = append(stockSlice, dummyStockModel1)
	stockSlice = append(stockSlice, dummyStockModel2)
	return stockSlice, nil
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

func (m *MockStockMapper) BeginTransaction() *errors.Error {
	return nil
}

func (m *MockStockMapper) Commit() *errors.Error {
	return nil
}

func (m *MockStockMapper) Rollback() *errors.Error {
	return nil
}

//Mock object for purchase datamapper
type MockPurchaseMapper struct {
}

var dummyPurchaseModel1 = &model.Purchase{
	PurchaseID: "dummyPurchaseId",
	Date:       time.Now(),
	Status:     model.PurchaseStatusDraft,
}
var dummyPurchaseItem1 = &model.PurchaseItem{
	Sku:      "dummySku",
	Quantity: 12,
	BuyPrice: 48000,
	Note:     "dummySku purchase",
}
var dummyPurchaseItem2 = &model.PurchaseItem{
	Sku:      "dummySku2",
	Quantity: 5,
	BuyPrice: 62000,
	Note:     "dummySku2 purchase",
}
var dummyPurchaseModel2 = &model.Purchase{
	PurchaseID: "dummyPurchaseIdX",
	Date:       time.Now(),
	Status:     model.PurchaseStatusDone,
}

func (m *MockPurchaseMapper) FindByID(id string) (model.Model, *errors.Error) {
	modelObj := dummyPurchaseModel1
	items := make(map[string]*model.PurchaseItem, 0)
	items["dummySku"] = dummyPurchaseItem1
	items["dummySku2"] = dummyPurchaseItem2
	modelObj.Items = items
	return modelObj, nil
}

func (m *MockPurchaseMapper) FindAll() ([]model.Model, *errors.Error) {
	var modelSlice []model.Model

	modelObj := dummyPurchaseModel2
	items := make(map[string]*model.PurchaseItem, 0)
	items["dummySku"] = dummyPurchaseItem1
	items["dummySku2"] = dummyPurchaseItem2
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

func (m *MockPurchaseMapper) BeginTransaction() *errors.Error {
	return nil
}

func (m *MockPurchaseMapper) Commit() *errors.Error {
	return nil
}

func (m *MockPurchaseMapper) Rollback() *errors.Error {
	return nil
}

//Mock object for sales datamapper
type MockSalesMapper struct {
}

var dummySalesModel1 = &model.Sales{
	InvoiceID: "dummyInvoice",
	Date:      time.Now(),
	Status:    model.SalesStatusDone,
	Note:      "Dummy invoice",
}
var dummySalesItem1 = &model.SaleItem{
	Sku:       "dummySku",
	Quantity:  3,
	BuyPrice:  50000,
	SellPrice: 55000,
}
var dummySalesModel2 = &model.Sales{
	InvoiceID: "dummyInvoic2e",
	Date:      time.Now(),
	Status:    model.SalesStatusDone,
	Note:      "Dummy invoice 2",
}

func (m *MockSalesMapper) FindByID(id string) (model.Model, *errors.Error) {
	modelObj := dummySalesModel1
	items := make(map[string]*model.SaleItem, 0)
	items["dummySku"] = dummySalesItem1
	modelObj.Items = items
	return modelObj, nil
}

func (m *MockSalesMapper) FindAll() ([]model.Model, *errors.Error) {
	var modelSlice []model.Model

	modelObj := dummySalesModel2
	items := make(map[string]*model.SaleItem, 0)
	items["dummySku"] = dummySalesItem1
	modelObj.Items = items

	modelSlice = append(modelSlice, modelObj)

	anotherModelObj, _ := m.FindByID("not used")
	modelSlice = append(modelSlice, anotherModelObj)

	return modelSlice, nil
}

func (m *MockSalesMapper) FindByDoneStatusAndDateRange(time.Time, time.Time) ([]model.Model, *errors.Error) {
	return m.FindAll()
}

func (m *MockSalesMapper) Insert(model model.Model) *errors.Error {
	return nil
}

func (m *MockSalesMapper) Update(model model.Model) *errors.Error {
	return nil
}

func (m *MockSalesMapper) Delete(model model.Model) *errors.Error {
	return nil
}

func (m *MockSalesMapper) Save(model model.Model) *errors.Error {
	return nil
}

func (m *MockSalesMapper) BeginTransaction() *errors.Error {
	return nil
}

func (m *MockSalesMapper) Commit() *errors.Error {
	return nil
}

func (m *MockSalesMapper) Rollback() *errors.Error {
	return nil
}
