package handler

import (
	"fmt"
	"ijah-inventory/repository/inventory/domain/inventory/service"
	"time"

	"bytes"
	"encoding/csv"
	"net/http"
	"strconv"
)

//ExportSalesCSVHandler is a specific http handler for creating sale
type ExportSalesCSVHandler struct {
	Handler
	InventoryService *service.Inventory `inject:"inventoryService"`
}

//ExportSalesCSVHandle is the implementation of http handler for a ExportSalesCSVHandler object
func (h *ExportSalesCSVHandler) ExportSalesCSVHandle(w http.ResponseWriter, r *http.Request) error {
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

	salesData, errs := h.InventoryService.GetAllSalesValue(startTimeObj, endTimeObj)
	if errs != nil {
		//compose failed response
		return composeError(errs)
	}

	//compose the csv data
	var csvString [][]string
	csvString = make([][]string, 0)

	//1st row is for summary data
	//summary data order:
	//start date, end date, total item kind, total item quantity, sale count, total profit, total sales turnover (omset)
	firstRow := make([]string, 0)
	firstRow = append(firstRow, salesData.StartDate.Format(csvDateLayout))
	firstRow = append(firstRow, salesData.EndDate.Format(csvDateLayout))
	firstRow = append(firstRow, strconv.Itoa(salesData.TotalItemKind))
	firstRow = append(firstRow, strconv.FormatInt(salesData.TotalQuantity, 10))
	firstRow = append(firstRow, strconv.Itoa(salesData.SaleCount))
	firstRow = append(firstRow, strconv.FormatFloat(salesData.Profit, 'f', 2, 64))
	firstRow = append(firstRow, strconv.FormatFloat(salesData.SalesTurnOver, 'f', 2, 64))
	csvString = append(csvString, firstRow)

	//the remaining rows are for the items
	//data order:
	//sku, quantity, buy price, total amount (buy price * quantity)
	for _, val := range salesData.Items {
		quantityStr := strconv.FormatInt(val.Quantity, 10)
		buyPriceStr := strconv.FormatFloat(val.BuyPrice, 'f', 2, 64)
		sellPriceStr := strconv.FormatFloat(val.SellPrice, 'f', 2, 64)
		profitStr := strconv.FormatFloat(val.Profit, 'f', 2, 64)

		newRow := make([]string, 0)
		newRow = append(newRow, val.Sku)
		newRow = append(newRow, quantityStr)
		newRow = append(newRow, buyPriceStr)
		newRow = append(newRow, sellPriceStr)
		newRow = append(newRow, profitStr)
		csvString = append(csvString, newRow)
	}

	//create csv writer
	buff := &bytes.Buffer{} //placeholder buffer
	csvWriter := csv.NewWriter(buff)
	for _, val := range csvString {
		errw := csvWriter.Write(val)
		if errw != nil {
			return composeError(errw)
		}
	}
	csvWriter.Flush() //flush to buffer

	//output the csv
	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Disposition", "attachment; filename=SalesValue_"+startTimeObj.Format(csvDateLayout)+"-"+endTimeObj.Format(csvDateLayout)+".csv")
	_, errOutput := buff.WriteTo(w)
	if errOutput != nil {
		return composeError(errOutput)
	}
	return nil
}

//StartUp allows the handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (h *ExportSalesCSVHandler) StartUp() {
	//Note: perform initialization/bootstrapping here
}

//Shutdown allows the handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (h *ExportSalesCSVHandler) Shutdown() {
	//Note: perform any cleanup here
}
