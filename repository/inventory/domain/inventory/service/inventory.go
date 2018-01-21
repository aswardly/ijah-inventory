//Package service provide definitions for inventory service layer
package service

import (
	"fmt"
	"time"

	"github.com/go-errors/errors"

	"ijah-inventory/repository/inventory/domain/inventory/datamapper"
	"ijah-inventory/repository/inventory/domain/inventory/model"
)

//SaleItem is a definition of items in a sale
type SaleItem struct {
	Sku      string `json:"sku"`
	Quantity int64  `json:"quantity"`
}

//StockValue is a struct containing stock value information
type StockValue struct {
	Date          time.Time                  `json:"date"`
	TotalQuantity int64                      `json:"totalQuantity"`
	TotalAmount   float64                    `json:"totalAmount"`
	TotalItemKind int                        `json:"totalItemKind"`
	Items         map[string]*StockValueItem `json:"items"`
}

//StockValueItem is a struct containing stock value for a specific Sku
type StockValueItem struct {
	Sku         string  `json:"sku"`
	Quantity    int64   `json:"quantity"`
	BuyPrice    float64 `json:"buyPrice"`
	TotalAmount float64 `json:"totalAmount"`
}

//SaleValue is a struct containing sales value information
type SaleValue struct {
	StartDate     time.Time        `json:"startDate"`
	EndDate       time.Time        `json:"endDate"`
	TotalQuantity int64            `json:"totalQuantity"`
	TotalItemKind int              `json:"totalItemKind"`
	SaleCount     int              `json:"saleCount"`
	SalesTurnOver float64          `json:"omzet"`
	Profit        float64          `json:"totalProfit"`
	Items         []*SaleValueItem `json:"items"`
}

//SaleValueItem is a struct containing sales value for a specific Sku
type SaleValueItem struct {
	Sku       string  `json:"sku"`
	Quantity  int64   `json:"quantity"`
	BuyPrice  float64 `json:"buyPrice"`
	SellPrice float64 `json:"sellPrice"`
	Profit    float64 `json:"profit"`
}

//NewInventory returns a new inventory service object
func NewInventory(stockMapper, purchaseMapper, salesMapper datamapper.DataMapper) *Inventory {
	return &Inventory{
		StockDatamapper:    stockMapper,
		PurchaseDatamapper: purchaseMapper,
		SalesDatamapper:    salesMapper,
	}
}

//Inventory is a service object dealing with inventory business domain
type Inventory struct {
	StockDatamapper    datamapper.DataMapper `inject:"stockDatamapper"`
	PurchaseDatamapper datamapper.DataMapper `inject:"purchaseDatamapper"`
	SalesDatamapper    datamapper.DataMapper `inject:"salesDatamapper"`
}

//GetItemInfo is a function for obtaining information of an item
func (i *Inventory) GetItemInfo(sku string) (*model.Stock, *errors.Error) {
	foundItem, err := i.StockDatamapper.FindByID(sku)
	if err != nil {
		return nil, err
	}
	foundItemObj, ok := foundItem.(*model.Stock)
	if false == ok {
		return nil, errors.Wrap(fmt.Errorf("Failed asserting returned model"), 0)
	}

	return foundItemObj, nil
}

//AddSKU is a function for adding a new item type to inventory
func (i *Inventory) AddSKU(sku string, quantity int64, buyPrice, sellPrice float64) *errors.Error {
	//compose stock model object
	newSku := &model.Stock{
		Sku:       sku,
		Quantity:  quantity,
		BuyPrice:  buyPrice,
		SellPrice: sellPrice,
	}
	return i.StockDatamapper.Insert(newSku)
}

//UpdateSKU is a function for updating SKU info
func (i *Inventory) UpdateSKU(sku string, quantity int64, buyPrice, sellPrice float64) *errors.Error {
	stockObj, err := i.GetItemInfo(sku)
	if err != nil {
		return err
	}
	stockObj.Quantity = quantity
	stockObj.BuyPrice = buyPrice
	stockObj.SellPrice = sellPrice

	return i.StockDatamapper.Update(stockObj)
}

//CreateSale is a fucntion for creating a new sale
func (i *Inventory) CreateSale(invoiceNo, note string, items ...SaleItem) (bool, *errors.Error) {
	//compose sale domain model
	newSale := &model.Sales{
		InvoiceID: invoiceNo,
		Date:      time.Now(),
		Note:      note,
		Status:    model.SalesStatusDraft,
	}
	newSalesItems := make(map[string]*model.SaleItem, 0)
	for _, val := range items {
		//get buy and sell price of the sku
		foundItem, err := i.StockDatamapper.FindByID(val.Sku)
		if err != nil {
			if err.Err == datamapper.ErrNotFound {
				//invalid sku, cannot continue
				return false, errors.Wrap(fmt.Errorf("Cannot create sale. Sku %v is not valid item", val.Sku), 0)
			}
			return false, errors.Wrap(err, 0)
		}
		foundItemObj, ok := foundItem.(*model.Stock)
		if false == ok {
			return false, errors.Wrap(fmt.Errorf("Failed asserting returned model"), 0)
		}
		//check whether sale quantity is enough
		if val.Quantity > foundItemObj.Quantity {
			return false, errors.Wrap(fmt.Errorf("Cannot create sale. Not enough stock for Sku %v", val.Sku), 0)
		}
		//compose sale item
		newItem := &model.SaleItem{
			Sku:       val.Sku,
			Quantity:  val.Quantity,
			BuyPrice:  foundItemObj.BuyPrice,
			SellPrice: foundItemObj.SellPrice,
		}
		newSalesItems[val.Sku] = newItem
	}
	newSale.Items = newSalesItems

	//save new sale
	err := i.SalesDatamapper.Insert(newSale)
	if err != nil {
		return false, errors.Wrap(err, 0)
	}
	return true, nil
}

//UpdateSale is a function for updating sale status
func (i *Inventory) UpdateSale(invoiceNo, status string) (bool, *errors.Error) {
	//validation, check whether given status is valid
	if status != model.SalesStatusDraft &&
		status != model.SalesStatusDone &&
		status != model.SalesStatusCanceled {
		return false, errors.Wrap(fmt.Errorf("Invalid status %v from param", status), 0)
	}
	//check whether the sale exists or not
	foundSale, err := i.SalesDatamapper.FindByID(invoiceNo)
	if err != nil {
		if err.Err == datamapper.ErrNotFound {
			//sale not found
			return false, errors.Wrap(fmt.Errorf("Sale %v is not found", invoiceNo), 0)
		}
		return false, errors.Wrap(err, 0)
	}
	foundSaleObj, ok := foundSale.(*model.Sales)
	if false == ok {
		return false, errors.Wrap(fmt.Errorf("Failed asserting returned model"), 0)
	}

	//note: stock datamapper and sale datamapper uses the same db session
	err = i.SalesDatamapper.BeginTransaction()
	if err != nil {
		return false, errors.Wrap(err, 0)
	}
	//sale status updated to Done from Other status
	if status == model.SalesStatusDone && foundSaleObj.Status != model.SalesStatusDone {
		for _, val := range foundSaleObj.Items {
			//update stock quantity
			saleItem, err := i.StockDatamapper.FindByID(val.Sku)
			if err != nil {
				i.SalesDatamapper.Rollback()
				return false, errors.Wrap(err, 0)
			}

			saleItemObj, ok := saleItem.(*model.Stock)
			if false == ok {
				i.SalesDatamapper.Rollback()
				return false, errors.Wrap(fmt.Errorf("Failed asserting stock"), 0)
			}

			if saleItemObj.Quantity < val.Quantity {
				i.SalesDatamapper.Rollback()
				return false, errors.Wrap(fmt.Errorf("Sku: %v doesn't have enough stock", saleItemObj.Sku), 0)
			}
			//update stock
			saleItemObj.Quantity -= val.Quantity
			err = i.StockDatamapper.Update(saleItemObj)
			if err != nil {
				i.SalesDatamapper.Rollback()
				return false, errors.Wrap(fmt.Errorf("Sku: %v stock update failed", saleItemObj.Sku), 0)
			}
		}
	}
	//update sale
	foundSaleObj.Status = status

	err = i.SalesDatamapper.Save(foundSaleObj)
	if err != nil {
		i.SalesDatamapper.Rollback()
		return false, errors.Wrap(err, 0)
	}
	i.SalesDatamapper.Commit()
	return true, nil
}

//GetAllStockValue is a function for obtaining current stock value
func (i *Inventory) GetAllStockValue() (*StockValue, *errors.Error) {
	currentStock, err := i.StockDatamapper.FindAll()
	if err != nil {
		if err.Err == datamapper.ErrNotFound {
			//no stock at all
			return nil, errors.Wrap(fmt.Errorf("No Sku available"), 0)
		}
		return nil, errors.Wrap(err, 0)
	}

	//compose stock value
	stockValue := &StockValue{
		Date: time.Now(),
	}
	var kind int            //total kind of sku available
	var totalAmount float64 //total amount of sku value, accumulate buy price * quantity for every sku
	var totalQuantity int64 //total quantity of all sku, accumulate quantity for every sku
	stockValueItems := make(map[string]*StockValueItem, 0)
	for _, val := range currentStock {
		valObj, ok := val.(*model.Stock)
		if false == ok {
			return nil, errors.Wrap(fmt.Errorf("Failed asserting returned model"), 0)
		}
		newStockValueItem := &StockValueItem{
			Sku:         valObj.Sku,
			Quantity:    valObj.Quantity,
			BuyPrice:    valObj.BuyPrice,
			TotalAmount: valObj.BuyPrice * float64(valObj.Quantity),
		}
		stockValueItems[valObj.Sku] = newStockValueItem
		totalAmount += newStockValueItem.TotalAmount
		totalQuantity += newStockValueItem.Quantity
		kind++
	}
	stockValue.Items = stockValueItems
	stockValue.TotalItemKind = kind
	stockValue.TotalAmount = totalAmount
	stockValue.TotalQuantity = totalQuantity

	return stockValue, nil
}

//GetAllSalesValue is a function for obtaining sales value of all sku
func (i *Inventory) GetAllSalesValue(startTime, endTime time.Time) (*SaleValue, *errors.Error) {
	//validate start and end date
	if !(startTime.Before(endTime) || startTime.Equal(endTime)) {
		return nil, errors.Wrap(fmt.Errorf("Invalid start date %v and end date %v from param", startTime.Format("2006-01-02"), endTime.Format("2006-01-02")), 0)
	}

	salesValue := &SaleValue{
		StartDate: startTime,
		EndDate:   endTime,
	}
	//get sales data from db

	//TODO: refactor this,
	//we need to type assert i.salesDataMapper to datamapper.saleMapper, since we need to execute method FindByDoneStatusAndDateRange
	//but got problem with unit test (a struct MockSaleMapper that implements DataMapper fails type assertion, since it is asserted as MockSaleMapper, not datamapper.Sale)
	salesDatamapper := i.SalesDatamapper
	//fmt.Printf("%T, %v\n", i.SalesDatamapper, i.SalesDatamapper)
	/*salesDatamapper, ok := i.SalesDatamapper.(*datamapper.Sale)
	if false == ok {
		return nil, errors.Wrap(fmt.Errorf("Failed asserting sales datamapper"), 0)
	}
	*/
	//salesData, err := salesDatamapper.FindByDoneStatusAndDateRange(startTime, endTime)
	salesData, err := salesDatamapper.FindAll()
	if err != nil {
		if err.Err == datamapper.ErrNotFound {
			//no data
			return nil, errors.Wrap(fmt.Errorf("No data available"), 0)
		}
		return nil, errors.Wrap(err, 0)
	}

	var totalProfit float64   //total profit = sellprice - buyprice, accumulate for every sale item
	var salesTurnover float64 //salesTurnover (omset), sellprice * quantity, accumulate for every sale item
	var totalQuantity int64   //total quantity of all items, accumulate quantity for every sale item
	var totalKind int         //total kind of sku sold during the given period
	var saleCount int         //total count of sales during the given period

	var tempSku = make(map[string]bool) //temporary storage for counting total kind of sku
	saleValueItems := make([]*SaleValueItem, 0)
	for _, val := range salesData {
		valObj, ok := val.(*model.Sales)
		if false == ok {
			return nil, errors.Wrap(fmt.Errorf("Failed asserting returned model"), 0)
		}
		saleCount++

		//get the sales items
		for _, itemVal := range valObj.Items {
			if _, exists := tempSku[itemVal.Sku]; false == exists {
				totalKind++
				tempSku[itemVal.Sku] = true
			}
			totalQuantity += itemVal.Quantity
			totalProfit += (itemVal.SellPrice - itemVal.BuyPrice) * float64(itemVal.Quantity)
			salesTurnover += itemVal.SellPrice * float64(itemVal.Quantity)

			saleValueItem := &SaleValueItem{
				Sku:       itemVal.Sku,
				BuyPrice:  itemVal.BuyPrice,
				SellPrice: itemVal.SellPrice,
				Quantity:  itemVal.Quantity,
				Profit:    (itemVal.SellPrice - itemVal.BuyPrice) * float64(itemVal.Quantity),
			}
			saleValueItems = append(saleValueItems, saleValueItem)
		}
	}
	salesValue.Items = saleValueItems
	salesValue.TotalQuantity = totalQuantity
	salesValue.TotalItemKind = totalKind
	salesValue.Profit = totalProfit
	salesValue.SalesTurnOver = salesTurnover
	salesValue.SaleCount = saleCount
	return salesValue, nil
}
