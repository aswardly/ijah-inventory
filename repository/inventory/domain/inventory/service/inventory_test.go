//service_test provides unit tests for user domain service layer
package service_test

import (
	"ijah-inventory/repository/inventory/domain/inventory/model"
	"ijah-inventory/repository/inventory/domain/inventory/service"
	"time"

	"os"
	"testing"
)

var inventoryService *service.Inventory

func TestMain(m *testing.M) {
	//test setup
	inventoryService = &service.Inventory{
		StockDatamapper:    &MockStockMapper{},
		PurchaseDatamapper: &MockPurchaseMapper{},
		SalesDatamapper:    &MockSalesMapper{},
	}

	//run tests
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGetItemInfo(t *testing.T) {
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
}

func TestAddSKU(t *testing.T) {
	err := inventoryService.AddSKU("dummyNewSku", 250, 55000, 60000)
	t.Run("err return must be nil", func(t *testing.T) {
		if err != nil {
			t.Errorf("expected nil but got %v", err)
		}
	})
}

func TestUpdateSKU(t *testing.T) {
	err := inventoryService.UpdateSKU("dummySku", 500, 60000, 65000)
	t.Run("err return must be nil", func(t *testing.T) {
		if err != nil {
			t.Errorf("expected nil but got %v", err)
		}
	})
}

func TestCreateSale(t *testing.T) {
	saleItem := service.SaleItem{
		Sku:      "dummySku",
		Quantity: 10,
	}
	ok, err := inventoryService.CreateSale("newInvoice01", "dummy new invoice", saleItem)

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
}

func TestUpdateSale(t *testing.T) {
	ok, err := inventoryService.UpdateSale("dummyInvoice", model.SalesStatusDone)

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
}

func TestGetAllStockValue(t *testing.T) {
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
}

func TestGetAllSalesValue(t *testing.T) {
	saleValue, err := inventoryService.GetAllSalesValue(time.Now(), time.Now()) //on dummy sales mapper these params are ignored anyway

	t.Run("GetAllSalesValue return must be sale value object", func(t *testing.T) {
		if getType(saleValue) != "*SaleValue" {
			t.Errorf("expected *SaleValue but got %v", getType(saleValue))
		}
	})
	t.Run("err return must be nil", func(t *testing.T) {
		if err != nil {
			t.Errorf("expected nil but got %v", err)
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
}
