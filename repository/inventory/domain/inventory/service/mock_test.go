//service_test provides unit tests for user domain service layer
package service_test

import (
	"fmt"
	"ijah-inventory/repository/inventory/domain/inventory/datamapper"
	"ijah-inventory/repository/inventory/domain/inventory/model"

	"time"

	"github.com/go-errors/errors"
)

//Mock object for stock datamapper (successful responses)
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
	stockSlice := make([]model.Model, 0)
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

//Mock object for purchase datamapper (successful responses)
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

//Mock object for sales datamapper (successful responses)
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

//Mock object for stock datamapper (failed responses)
type MockFailedStockMapper struct {
}

func (m *MockFailedStockMapper) FindByID(id string) (model.Model, *errors.Error) {
	return nil, errors.Wrap(datamapper.ErrNotFound, 0)
}

func (m *MockFailedStockMapper) FindAll() ([]model.Model, *errors.Error) {
	return nil, errors.Wrap(datamapper.ErrNotFound, 0)
}

func (m *MockFailedStockMapper) Insert(model model.Model) *errors.Error {
	return errors.Wrap(fmt.Errorf("dummy failure for insert"), 0)
}

func (m *MockFailedStockMapper) Update(model model.Model) *errors.Error {
	return errors.Wrap(fmt.Errorf("dummy failure for update"), 0)
}

func (m *MockFailedStockMapper) Delete(model model.Model) *errors.Error {
	return errors.Wrap(fmt.Errorf("dummy failure for delete"), 0)
}

func (m *MockFailedStockMapper) Save(model model.Model) *errors.Error {
	return errors.Wrap(fmt.Errorf("dummy failure for save"), 0)
}

//Mock object for purchase datamapper (failed responses)
type MockFailedPurchaseMapper struct {
}

func (m *MockFailedPurchaseMapper) FindByID(id string) (model.Model, *errors.Error) {
	return nil, errors.Wrap(datamapper.ErrNotFound, 0)
}

func (m *MockFailedPurchaseMapper) FindAll() ([]model.Model, *errors.Error) {
	return nil, errors.Wrap(datamapper.ErrNotFound, 0)
}

func (m *MockFailedPurchaseMapper) Insert(model model.Model) *errors.Error {
	return errors.Wrap(fmt.Errorf("dummy failure for insert"), 0)
}

func (m *MockFailedPurchaseMapper) Update(model model.Model) *errors.Error {
	return errors.Wrap(fmt.Errorf("dummy failure for update"), 0)
}

func (m *MockFailedPurchaseMapper) Delete(model model.Model) *errors.Error {
	return errors.Wrap(fmt.Errorf("dummy failure for delete"), 0)
}

func (m *MockFailedPurchaseMapper) Save(model model.Model) *errors.Error {
	return errors.Wrap(fmt.Errorf("dummy failure for save"), 0)
}

//Mock object for sales datamapper (failed responses)
type MockFailedSalesMapper struct {
}

func (m *MockFailedSalesMapper) FindByID(id string) (model.Model, *errors.Error) {
	return nil, errors.Wrap(datamapper.ErrNotFound, 0)
}

func (m *MockFailedSalesMapper) FindAll() ([]model.Model, *errors.Error) {
	return nil, errors.Wrap(datamapper.ErrNotFound, 0)
}

func (m *MockFailedSalesMapper) FindByDoneStatusAndDateRange(time.Time, time.Time) ([]model.Model, *errors.Error) {
	return nil, errors.Wrap(datamapper.ErrNotFound, 0)
}

func (m *MockFailedSalesMapper) Insert(model model.Model) *errors.Error {
	return nil
}

func (m *MockFailedSalesMapper) Update(model model.Model) *errors.Error {
	return errors.Wrap(fmt.Errorf("dummy failure for rollback"), 0)
}

func (m *MockFailedSalesMapper) Delete(model model.Model) *errors.Error {
	return errors.Wrap(fmt.Errorf("dummy failure for rollback"), 0)
}

func (m *MockFailedSalesMapper) Save(model model.Model) *errors.Error {
	return errors.Wrap(fmt.Errorf("dummy failure for rollback"), 0)
}

//Special Mock object for CreateSale test (combination of successful and failed cases)
type MockCreateSalesMapper struct {
}

//must return err (failed) for testing successful case of CreateSale
func (m *MockCreateSalesMapper) FindByID(id string) (model.Model, *errors.Error) {
	return nil, errors.Wrap(datamapper.ErrNotFound, 0)
}

func (m *MockCreateSalesMapper) FindAll() ([]model.Model, *errors.Error) {
	modelSlice := make([]model.Model, 0)

	modelObj := dummySalesModel2
	items := make(map[string]*model.SaleItem, 0)
	items["dummySku"] = dummySalesItem1
	modelObj.Items = items

	modelSlice = append(modelSlice, modelObj)

	anotherModelObj, _ := m.FindByID("not used")
	modelSlice = append(modelSlice, anotherModelObj)

	return modelSlice, nil
}

func (m *MockCreateSalesMapper) FindByDoneStatusAndDateRange(time.Time, time.Time) ([]model.Model, *errors.Error) {
	return m.FindAll()
}

func (m *MockCreateSalesMapper) Insert(model model.Model) *errors.Error {
	return nil
}

func (m *MockCreateSalesMapper) Update(model model.Model) *errors.Error {
	return nil
}

func (m *MockCreateSalesMapper) Delete(model model.Model) *errors.Error {
	return nil
}

func (m *MockCreateSalesMapper) Save(model model.Model) *errors.Error {
	return nil
}
