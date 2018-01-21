package handler

import (
	"ijah-inventory/repository/inventory/domain/inventory/service"
	"time"

	"bytes"
	"encoding/csv"
	"net/http"
	"strconv"
)

//ExportStockCSVHandler is a specific http handler for creating sale
type ExportStockCSVHandler struct {
	Handler
	InventoryService *service.Inventory `inject:"inventoryService"`
}

//csvDateLayout is the layout used for formatting time.Time object to string in csv output
const csvDateLayout = "2006-01-02"

//ExportStockCSVHandle is the implementation of http handler for a ExportStockCSVHandler object
func (h *ExportStockCSVHandler) ExportStockCSVHandle(w http.ResponseWriter, r *http.Request) error {

	stockData, err := h.InventoryService.GetAllStockValue()
	if err != nil {
		//compose failed response
		return composeError(err)
	}

	//compose the csv data
	var csvString [][]string
	csvString = make([][]string, 0)

	//1st row is for summary data
	firstRow := make([]string, 0)
	firstRow = append(firstRow, stockData.Date.Format(csvDateLayout))
	firstRow = append(firstRow, strconv.Itoa(stockData.TotalItemKind))
	firstRow = append(firstRow, strconv.FormatInt(stockData.TotalQuantity, 10))
	firstRow = append(firstRow, strconv.FormatFloat(stockData.TotalAmount, 'f', 2, 64))
	csvString = append(csvString, firstRow)
	//the remaining rows are for the items
	for _, val := range stockData.Items {
		quantityStr := strconv.FormatInt(val.Quantity, 10)
		buyPriceStr := strconv.FormatFloat(val.BuyPrice, 'f', 2, 64)
		totalAmountStr := strconv.FormatFloat(val.TotalAmount, 'f', 2, 64)

		newRow := make([]string, 0)
		newRow = append(newRow, val.Sku)
		newRow = append(newRow, quantityStr)
		newRow = append(newRow, buyPriceStr)
		newRow = append(newRow, totalAmountStr)
		csvString = append(csvString, newRow)
	}

	//fmt.Printf("csv data: %+v\n", csvString)

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
	now := time.Now()
	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Disposition", "attachment; filename=StockValue"+now.Format(csvDateLayout)+".csv")
	_, errOutput := buff.WriteTo(w)
	if errOutput != nil {
		return composeError(errOutput)
	}
	return nil
	/*
		record := []string{"test1", "test2", "test3"} // just some test data to use for the wr.Writer() method below.


		wr := csv.NewWriter(b)     // creates a csv writer that uses the io buffer.
		for i := 0; i < 100; i++ { // make a loop for 100 rows just for testing purposes
			wr.Write(record) // converts array of string to comma seperated values for 1 row.
		}
		wr.Flush() // writes the csv writer data to  the buffered data io writer(b(bytes.buffer))

		ctx.ResponseWriter.Header().Set("Content-Type", "text/csv") // setting the content type header to text/csv

		ctx.ResponseWriter.Header().Set("Content-Type", "text/csv")
		ctx.ResponseWriter.Header().Set("Content-Disposition", "attachment;filename=TheCSVFileName.csv")
		ctx.ResponseWriter.Write(b.Bytes())

		//compose successful response
		response := SimpleResponseStruct{}
		response.Code = ErrCodeSuccessful
		response.Message = "Inquiry successful"
		response.Data = stockObj

		successfulResponse, statusError := composeJSONResponse(response)
		if statusError != nil {
			return statusError
		}
		//else no problem in json marshalling the response
		w.Write([]byte(successfulResponse))
	*/
}

//StartUp allows the handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (h *ExportStockCSVHandler) StartUp() {
	//Note: perform initialization/bootstrapping here
}

//Shutdown allows the handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (h *ExportStockCSVHandler) Shutdown() {
	//Note: perform any cleanup here
}
