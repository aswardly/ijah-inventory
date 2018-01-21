//service_test provides unit tests for user domain service layer
package service_test

import (
	"database/sql"
	"fmt"
	"ijah-inventory/repository/inventory/domain/inventory/model"
	"ijah-inventory/repository/inventory/domain/inventory/service"

	"os"
	"reflect"
	"testing"
	"time"

	sqlMock "github.com/DATA-DOG/go-sqlmock"
)

var inventoryService, failedInventoryService, successfulCreateSaleInventoryService *service.Inventory
var dummyDb *sql.DB
var dbMock sqlMock.Sqlmock

//getType is a function to get type of something (without package name)
//see: https://stackoverflow.com/questions/35790935/using-reflection-in-go-to-get-the-name-of-a-struct
func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func TestMain(m *testing.M) {
	//test setup

	//mock sql
	dummyDb, dbMock, _ = sqlMock.New()

	//successful case inventory service object
	inventoryService = &service.Inventory{
		StockDatamapper:    &MockStockMapper{},
		PurchaseDatamapper: &MockPurchaseMapper{},
		SalesDatamapper:    &MockSalesMapper{},
		DB:                 dummyDb, //Note: do not use a *sql.DB that connects to production database
	}

	//failed case inventory service object
	failedInventoryService = &service.Inventory{
		StockDatamapper:    &MockFailedStockMapper{},
		PurchaseDatamapper: &MockFailedPurchaseMapper{},
		SalesDatamapper:    &MockFailedSalesMapper{},
		DB:                 dummyDb, //Note: do not use a *sql.DB that connects to production database
	}

	//special case for createSale (combination of successful and failed datamapper)
	successfulCreateSaleInventoryService = &service.Inventory{
		StockDatamapper:    &MockStockMapper{},
		PurchaseDatamapper: &MockFailedPurchaseMapper{},
		SalesDatamapper:    &MockCreateSalesMapper{},
		DB:                 dummyDb, //Note: do not use a *sql.DB that connects to production database
	}

	//run tests
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGetItemInfo(t *testing.T) {
	//successful case
	itemInfo, err := inventoryService.GetItemInfo("dummySku")

	t.Run("GetItemInfo return must be stock model object", func(t *testing.T) {
		if getType(itemInfo) != "*Stock" {
			t.Errorf("expected *Stock but got %v", getType(itemInfo))
		}
	})

	t.Run("err returned must be nil", func(t *testing.T) {
		if err != nil {
			t.Errorf("expected nil but got %v", err)
		}
	})

	t.Run("check returned item", func(t *testing.T) {
		if itemInfo.Sku != dummyStockModel1.Sku {
			t.Errorf("expected %v but got %v", dummyStockModel1.Sku, itemInfo.Sku)
		}
	})

	//failed case
	failedItemInfo, failedErr := failedInventoryService.GetItemInfo("dummySku")

	t.Run("Failed GetItemInfo return must be nil", func(t *testing.T) {
		if failedItemInfo != nil {
			t.Errorf("expected nil but got %v", failedItemInfo)
		}
	})

	t.Run("err returned must type must be correct", func(t *testing.T) {
		if getType(failedErr) != "*Error" {
			t.Errorf("expected *Error but got %v", getType(failedErr))
		}
	})

}

func TestAddSKU(t *testing.T) {
	//successful case
	err := inventoryService.AddSKU("dummyNewSku", 250, 55000, 60000)
	t.Run("err return must be nil", func(t *testing.T) {
		if err != nil {
			t.Errorf("expected nil but got %v", err)
		}
	})

	//failed case
	failedErr := failedInventoryService.AddSKU("dummyNewSku", 250, 55000, 60000)
	t.Run("Failed err returned must type must be correct", func(t *testing.T) {
		if getType(failedErr) != "*Error" {
			t.Errorf("expected *Error but got %v", getType(failedErr))
		}
	})
}

func TestUpdateSKU(t *testing.T) {
	//successful case
	dbMock.ExpectBegin()
	dbMock.ExpectCommit()
	err := inventoryService.UpdateSKU("dummySku", 250, 50000, 55000)
	t.Run("err return must be nil", func(t *testing.T) {
		if err != nil {
			t.Errorf("expected nil but got %v", err)
		}
	})

	//failed case
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()
	failedErr := failedInventoryService.UpdateSKU("dummySku", 250, 50000, 55000)
	t.Run("Failed err returned must type must be correct", func(t *testing.T) {
		if getType(failedErr) != "*Error" {
			t.Errorf("expected *Error but got %v", getType(failedErr))
		}
	})
}

func TestCreateSale(t *testing.T) {
	//successful case
	saleItem := service.SaleItem{
		Sku:      "dummySku",
		Quantity: 10,
	}
	saleItemSlice := make([]service.SaleItem, 0)
	saleItemSlice = append(saleItemSlice, saleItem)

	dbMock.ExpectBegin()
	dbMock.ExpectCommit()
	ok, errt := successfulCreateSaleInventoryService.CreateSale("newInvoiceId", "dummy new invoice", saleItemSlice)
	fmt.Printf("%+v\n", errt)
	fmt.Printf("%+v\n", ok)
	t.Run("return must be true", func(t *testing.T) {
		if true != ok {
			t.Errorf("expected true but got %v", ok)
		}
	})
	t.Run("err return must be nil", func(t *testing.T) {
		if errt != nil {
			t.Errorf("expected nil but got %v", errt)
		}
	})

	//failed case
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()
	failedOk, failedErr := failedInventoryService.CreateSale("newInvoiceId", "dummy new invoice", saleItemSlice)
	t.Run("Failed return must be true", func(t *testing.T) {
		if false != failedOk {
			t.Errorf("expected false but got %v", failedOk)
		}
	})
	t.Run("Failed err returned must type must be correct", func(t *testing.T) {
		if getType(failedErr) != "*Error" {
			t.Errorf("expected *Error but got %v", getType(failedErr))
		}
	})
}

func TestUpdateSale(t *testing.T) {
	//successful case
	dbMock.ExpectBegin()
	dbMock.ExpectCommit()
	ok, err := inventoryService.UpdateSale("dummyInvoice", model.SalesStatusDone)
	fmt.Printf("errx :%+v\n", err)
	t.Run("return must be true", func(t *testing.T) {
		if true != ok {
			t.Errorf("expected true but got %v", ok)
		}
	})
	t.Run("err return must be nil", func(t *testing.T) {
		if err != nil {
			t.Errorf("expected nil but got %v", err)
		}
	})

	//failed case
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()
	failedOk, failedErr := failedInventoryService.UpdateSale("dummyInvoice", model.SalesStatusDone)
	t.Run("Failed return must be true", func(t *testing.T) {
		if false != failedOk {
			t.Errorf("expected false but got %v", failedOk)
		}
	})
	t.Run("Failed err returned must type must be correct", func(t *testing.T) {
		if getType(failedErr) != "*Error" {
			t.Errorf("expected *Error but got %v", getType(failedErr))
		}
	})
}

func TestGetAllStockValue(t *testing.T) {
	//successful case
	stockValue, err := inventoryService.GetAllStockValue()
	t.Run("GetAllStockValue return must be stock value object", func(t *testing.T) {
		if getType(stockValue) != "*StockValue" {
			t.Errorf("expected *StockValue but got %v", getType(stockValue))
		}
	})

	t.Run("err return must be nil", func(t *testing.T) {
		if err != nil {
			t.Errorf("expected nil but got %v", err)
		}
	})

	t.Run("check stock value properties", func(t *testing.T) {
		if stockValue.TotalAmount != 19700000 {
			t.Errorf("expected totalAmount %v but got %v", 19700000, stockValue.TotalAmount)
		}
		if stockValue.TotalItemKind != 2 {
			t.Errorf("expected totalItemKind %v but got %v", 2, stockValue.TotalItemKind)
		}
		if stockValue.TotalQuantity != 370 {
			t.Errorf("expected totalQuantity %v but got %v", 370, stockValue.TotalQuantity)
		}

		for key, val := range stockValue.Items {
			if key == "dummySku" {
				if val.TotalAmount != 12500000 {
					t.Errorf("sku %v expected totalAmount %v but got %v", key, 12500000, val.TotalAmount)
				}
				if val.BuyPrice != 50000 {
					t.Errorf("sku %v expected buyPrice %v but got %v", key, 50000, val.BuyPrice)
				}
				if val.Quantity != 250 {
					t.Errorf("sku %v expected quantity %v but got %v", key, 250, val.Quantity)
				}
				if val.Sku != "dummySku" {
					t.Errorf("sku %v expected quantity %v but got %v", key, "dummySku", val.Sku)
				}
			} else { //key == "dummySku2"
				if val.TotalAmount != 7200000 {
					t.Errorf("sku %v expected totalAmount %v but got %v", key, 7200000, val.TotalAmount)
				}
				if val.BuyPrice != 60000 {
					t.Errorf("sku %v expected buyPrice %v but got %v", key, 60000, val.BuyPrice)
				}
				if val.Quantity != 120 {
					t.Errorf("sku %v expected quantity %v but got %v", key, 120, val.Quantity)
				}
				if val.Sku != "dummySku2" {
					t.Errorf("sku %v expected quantity %v but got %v", key, "dummySku2", val.Sku)
				}
			}
		}
	})

	//failed case
	failedStockValue, failedErr := failedInventoryService.GetAllStockValue()
	t.Run("Failed GetAllStockValue return must be nil", func(t *testing.T) {
		if failedStockValue != nil {
			t.Errorf("expected nil but got %v", failedStockValue)
		}
	})
	t.Run("Failed err returned must type must be correct", func(t *testing.T) {
		if getType(failedErr) != "*Error" {
			t.Errorf("expected *Error but got %v", getType(failedErr))
		}
	})
}

func TestGetAllSalesValue(t *testing.T) {
	//successful case
	saleValue, errs := inventoryService.GetAllSalesValue(time.Now(), time.Now()) //on dummy sales mapper these params are ignored
	t.Run("GetAllSalesValue return must be sale value object", func(t *testing.T) {
		if getType(saleValue) != "*SaleValue" {
			t.Errorf("expected *SaleValue but got %v", getType(saleValue))
		}
	})
	//TODO: this will fail since type assertion from datamapper.DataMapper to *datamapper.Sale in the inventory service fails, since we initiate the sales datamapper with a mock (type MockSSaleMapper)
	t.Run("err returned must be nil", func(t *testing.T) {
		if errs != nil {
			t.Errorf("expected nil but got %v", errs)
		}
	})

	t.Run("check sale value properties", func(t *testing.T) {
		if saleValue.TotalQuantity != 6 {
			t.Errorf("expected totalQuantity %v but got %v", 34, saleValue.TotalQuantity)
		}
		if saleValue.SaleCount != 2 {
			t.Errorf("expected saleCount %v but got %v", 2, saleValue.SaleCount)
		}
		if saleValue.TotalItemKind != 1 {
			t.Errorf("expected totalItemKind %v but got %v", 2, saleValue.TotalItemKind)
		}
		if saleValue.Profit != 30000 {
			t.Errorf("expected profit %v but got %v", 30000, saleValue.Profit)
		}
		if saleValue.SalesTurnOver != 330000 {
			t.Errorf("expected salesTurnOver %v but got %v", 330000, saleValue.SalesTurnOver)
		}
		for _, val := range saleValue.Items {
			if val.Sku == "dummySku" {
				if val.Profit != 15000 {
					t.Errorf("sku dummySku expected profit %v but got %v", 15000, val.Profit)
				}
			}
		}
	})

	//failed case
	failedSaleValue, failedErr := failedInventoryService.GetAllSalesValue(time.Now(), time.Now()) //on dummy sales mapper these params are ignored
	t.Run("Failed GetAllSalesValue return must be nil", func(t *testing.T) {
		if failedSaleValue != nil {
			t.Errorf("expected nil but got %v", getType(failedSaleValue))
		}
	})
	t.Run("Failed err returned must type must be correct", func(t *testing.T) {
		if getType(failedErr) != "*Error" {
			t.Errorf("expected *Error but got %v", getType(failedErr))
		}
	})
}
