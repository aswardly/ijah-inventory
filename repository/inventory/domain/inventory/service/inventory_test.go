//service_test provides unit tests for user domain service layer
package service_test

import (
	"ijah-inventory/repository/inventory/domain/inventory/service"

	"os"
	"testing"
)

var inventoryService *service.Inventory

func TestMain(m *testing.M) {
	//test setup
	inventoryService = service.Inventory{
		StockDatamapper:    &MockStockMapper{},
		PurchaseDatamapper: &MockPurchaseMapper{},
		SalesDatamapper:    &MockSalesMapper{},
	}

	//run tests
	exitCode := m.Run()
	os.Exit(exitCode)
}
