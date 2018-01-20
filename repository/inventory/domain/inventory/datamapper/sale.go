//Package datamapper provides the definitions of datamapper
package datamapper

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-errors/errors"
	_ "github.com/mattn/go-sqlite3"

	"ijah-inventory/repository/inventory/domain/inventory/model"
)

//Sale is a struct of datamapper for purchase domain model
type Sale struct {
	db *sql.DB
	tx *sql.Tx
}

//NewSale creates a new Purchase datamapper and returns a pointer to it
func NewSale(dbSession *sql.DB) *Sale {
	return &Sale{
		db: dbSession,
	}
}

//FindByID is a function for finding a record by id
func (s *Sale) FindByID(id string) (model.Model, *errors.Error) {
	stmt, err := s.db.Prepare("SELECT INVOICE_ID, DATETIME(SALE_DATE), STATUS, NOTE FROM sales WHERE INVOICE_ID = ?")
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}
	defer stmt.Close()

	var invoiceID, date, status, note sql.NullString

	row := stmt.QueryRow(id)
	err = row.Scan(&invoiceID, &date, &status, &note)
	if err != nil {
		var returnedErr error
		if err == sql.ErrNoRows {
			returnedErr = ErrNotFound
		} else {
			returnedErr = err
		}
		return nil, errors.Wrap(returnedErr, 0)
	}
	//how to check for no rows on returned error (on the service layer):
	//if err.Err != sql.ErrNoRows

	invoiceIDValue := invoiceID.String
	dateValue := date.String
	dateTimeValue, err := time.Parse(timeFormat, dateValue)

	statusValue := status.String
	noteValue := note.String

	salesModel := &model.Sales{
		InvoiceID: invoiceIDValue,
		Date:      dateTimeValue,
		Status:    statusValue,
		Note:      noteValue,
	}
	salesModel.SetLoadedFromStorage(true)

	//load purchase items
	itemStmt, err := s.db.Prepare("SELECT ID, SKU, QUANTITY, BUY_PRICE, SELL_PRICE FROM sales_items WHERE INVOICE_ID = ?")
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}
	defer itemStmt.Close()

	var itemID int64
	var sku sql.NullString
	var quantity sql.NullInt64
	var buyPrice, sellPrice sql.NullFloat64

	rows, err := itemStmt.Query(id)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}
	defer rows.Close()

	itemsRow := make(map[string]*model.SaleItem, 5)
	for rows.Next() {
		err := rows.Scan(&itemID, &sku, &quantity, &buyPrice, &sellPrice)
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}
		skuValue := sku.String
		quantityValue := quantity.Int64
		buyPriceValue := buyPrice.Float64
		sellPriceValue := sellPrice.Float64

		saleItemModel := &model.SaleItem{
			Sku:       skuValue,
			Quantity:  quantityValue,
			BuyPrice:  buyPriceValue,
			SellPrice: sellPriceValue,
		}
		saleItemModel.SetID(itemID)
		saleItemModel.SetLoadedFromStorage(true)

		itemsRow[skuValue] = saleItemModel
	}
	salesModel.Items = itemsRow

	return salesModel, nil
}

//FindAll is a function for finding all records
func (s *Sale) FindAll() ([]model.Model, *errors.Error) {
	rows, err := s.db.Query("SELECT INVOICE_ID, DATETIME(SALE_DATE), STATUS, NOTE FROM sales ORDER BY INVOICE_ID ASC")
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}
	defer rows.Close()

	var invoiceID, date, status, note sql.NullString

	var itemID int64
	var sku sql.NullString
	var quantity sql.NullInt64
	var buyPrice, sellPrice sql.NullFloat64

	var returnedRow []model.Model
	var firstScan = true
	for rows.Next() {
		err := rows.Scan(&invoiceID, &date, &status, &note)
		firstScan = false
		if err != nil {
			var returnedErr error
			if firstScan && err == sql.ErrNoRows {
				returnedErr = ErrNotFound
			} else {
				returnedErr = err
			}
			return nil, errors.Wrap(returnedErr, 0)
		}
		invoiceIDValue := invoiceID.String
		dateValue := date.String
		dateTimeValue, err := time.Parse(timeFormat, dateValue)
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}

		statusValue := status.String
		noteValue := note.String

		salesModel := &model.Sales{
			InvoiceID: invoiceIDValue,
			Date:      dateTimeValue,
			Status:    statusValue,
			Note:      noteValue,
		}
		salesModel.SetLoadedFromStorage(true)

		//load purchase items
		itemStmt, err := s.db.Prepare("SELECT ID, SKU, QUANTITY, BUY_PRICE, SELL_PRICE FROM sales_items WHERE INVOICE_ID = ?")
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}
		defer itemStmt.Close()

		itemRows, err := itemStmt.Query(invoiceIDValue)
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}
		defer itemRows.Close()

		itemsRow := make(map[string]*model.SaleItem, 5)
		for itemRows.Next() {
			err := itemRows.Scan(&itemID, &sku, &quantity, &buyPrice, &sellPrice)
			if err != nil {
				return nil, errors.Wrap(err, 0)
			}
			skuValue := sku.String
			quantityValue := quantity.Int64
			buyPriceValue := buyPrice.Float64
			sellPriceValue := sellPrice.Float64

			salesItemModel := &model.SaleItem{
				Sku:       skuValue,
				Quantity:  quantityValue,
				BuyPrice:  buyPriceValue,
				SellPrice: sellPriceValue,
			}
			salesItemModel.SetID(itemID)
			salesItemModel.SetLoadedFromStorage(true)

			itemsRow[skuValue] = salesItemModel
		}
		salesModel.Items = itemsRow

		returnedRow = append(returnedRow, salesModel)
	}
	return returnedRow, nil
}

//FindByDateRange is a function for finding sale record based on date range
//TODO: Refactor this and FindAll
func (s *Sale) FindByDateRange(startDate, endDate time.Time) ([]model.Model, *errors.Error) {
	startDateString := startDate.Format(dateFormat)
	endDateString := endDate.Format(dateFormat)

	stmt, err := s.db.Prepare("SELECT INVOICE_ID, DATETIME(SALE_DATE), STATUS, NOTE FROM sales WHERE SALE_DATE BETWEEN ? AND ? ORDER BY INVOICE_ID ASC")
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}
	rows, err := stmt.Query(startDateString, endDateString)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}
	defer rows.Close()

	var invoiceID, date, status, note sql.NullString

	var itemID int64
	var sku sql.NullString
	var quantity sql.NullInt64
	var buyPrice, sellPrice sql.NullFloat64

	var returnedRow []model.Model
	var firstScan = true
	for rows.Next() {
		err := rows.Scan(&invoiceID, &date, &status, &note)
		firstScan = false
		if err != nil {
			var returnedErr error
			if firstScan && err == sql.ErrNoRows {
				returnedErr = ErrNotFound
			} else {
				returnedErr = err
			}
			return nil, errors.Wrap(returnedErr, 0)
		}
		invoiceIDValue := invoiceID.String
		dateValue := date.String
		dateTimeValue, err := time.Parse(timeFormat, dateValue)
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}

		statusValue := status.String
		noteValue := note.String

		salesModel := &model.Sales{
			InvoiceID: invoiceIDValue,
			Date:      dateTimeValue,
			Status:    statusValue,
			Note:      noteValue,
		}
		salesModel.SetLoadedFromStorage(true)

		//load purchase items
		itemStmt, err := s.db.Prepare("SELECT ID, SKU, QUANTITY, BUY_PRICE, SELL_PRICE FROM sales_items WHERE INVOICE_ID = ?")
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}
		defer itemStmt.Close()

		itemRows, err := itemStmt.Query(invoiceIDValue)
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}
		defer itemRows.Close()

		itemsRow := make(map[string]*model.SaleItem, 5)
		for itemRows.Next() {
			err := itemRows.Scan(&itemID, &sku, &quantity, &buyPrice, &sellPrice)
			if err != nil {
				return nil, errors.Wrap(err, 0)
			}
			skuValue := sku.String
			quantityValue := quantity.Int64
			buyPriceValue := buyPrice.Float64
			sellPriceValue := sellPrice.Float64

			salesItemModel := &model.SaleItem{
				Sku:       skuValue,
				Quantity:  quantityValue,
				BuyPrice:  buyPriceValue,
				SellPrice: sellPriceValue,
			}
			salesItemModel.SetID(itemID)
			salesItemModel.SetLoadedFromStorage(true)

			itemsRow[skuValue] = salesItemModel
		}
		salesModel.Items = itemsRow

		returnedRow = append(returnedRow, salesModel)
	}
	return returnedRow, nil
}

//Insert is a function for inserting a record
func (s *Sale) Insert(salesModel model.Model) *errors.Error {
	salesModelObj, ok := salesModel.(*model.Sales)
	if false == ok {
		return errors.Wrap(fmt.Errorf("Failed asserting to *model.Sales"), 0)
	}

	foundModel, _ := s.FindByID(salesModel.GetID())
	if foundModel != nil {
		return errors.Wrap(fmt.Errorf("cannot insert, model with id: %v already exists", salesModel.GetID()), 0)
	}

	//start db transaction
	if s.tx == nil {
		tx, err := s.db.Begin()
		if err != nil {
			return errors.Wrap(err, 0)
		}
		s.tx = tx
	}
	stmt, err := s.db.Prepare("INSERT INTO sales(INVOICE_ID, SALE_DATE, STATUS, NOTE) values(?,?,?,?)")
	if err != nil {
		s.tx.Rollback()
		return errors.Wrap(err, 0)
	}
	dateString := salesModelObj.Date.Format(timeFormat)
	_, err = stmt.Exec(salesModelObj.InvoiceID, dateString, salesModelObj.Status, salesModelObj.Note)
	if err != nil {
		s.tx.Rollback()
		return errors.Wrap(err, 0)
	}

	//insert the items
	for _, val := range salesModelObj.Items {
		itemStmt, err := s.db.Prepare("INSERT INTO sales_items(INVOICE_ID, SKU, QUANTITY, BUY_PRICE, SELL_PRICE) values(?,?,?,?,?)")
		if err != nil {
			s.tx.Rollback()
			return errors.Wrap(err, 0)
		}
		_, err = itemStmt.Exec(salesModelObj.InvoiceID, val.Sku, val.Quantity, val.BuyPrice, val.SellPrice)
		if err != nil {
			s.tx.Rollback()
			return errors.Wrap(err, 0)
		}
	}
	s.tx.Commit()

	return nil
}

//Update is a function for updating record
func (s *Sale) Update(salesModel model.Model) *errors.Error {
	salesModelObj, ok := salesModel.(*model.Sales)
	if false == ok {
		return errors.Wrap(fmt.Errorf("Failed asserting to *model.Sales"), 0)
	}

	_, errs := s.FindByID(salesModel.GetID())
	if errs.Err == sql.ErrNoRows {
		return errors.Wrap(fmt.Errorf("cannot update, model with id: %v doesn't exist", salesModel.GetID()), 0)
	}

	//start db transaction
	if s.tx == nil {
		tx, err := s.db.Begin()
		if err != nil {
			return errors.Wrap(err, 0)
		}
		s.tx = tx
	}

	stmt, err := s.db.Prepare("UPDATE sales SET SALE_DATE=?, STATUS=?, NOTE=? WHERE INVOICE_ID=?")
	if err != nil {
		s.tx.Rollback()
		return errors.Wrap(err, 0)
	}
	dateString := salesModelObj.Date.Format(timeFormat)
	_, err = stmt.Exec(dateString, salesModelObj.Status, salesModelObj.Note, salesModelObj.InvoiceID)
	if err != nil {
		s.tx.Rollback()
		return errors.Wrap(err, 0)
	}

	//update items
	for _, val := range salesModelObj.Items {
		var itemStmt *sql.Stmt
		if false == val.GetLoadedFromStorage() {
			itemStmt, err = s.db.Prepare("INSERT INTO sales_items(INVOICE_ID, SKU, QUANTITY, BUY_PRICE, SELL_PRICE) values(?,?,?,?,?)")
			if err != nil {
				s.tx.Rollback()
				return errors.Wrap(err, 0)
			}
			_, err = itemStmt.Exec(salesModelObj.InvoiceID, val.Sku, val.Quantity, val.BuyPrice, val.SellPrice)
			if err != nil {
				s.tx.Rollback()
				return errors.Wrap(err, 0)
			}
		} else {
			itemStmt, err = s.db.Prepare("UPDATE sales_items SET QUANTITY=?, BUY_PRICE=?, SELL_PRICE=? WHERE ID=?")
			if err != nil {
				s.tx.Rollback()
				return errors.Wrap(err, 0)
			}
			_, err = itemStmt.Exec(val.Quantity, val.BuyPrice, val.SellPrice, val.GetID())
			if err != nil {
				s.tx.Rollback()
				return errors.Wrap(err, 0)
			}
		}
	}
	s.tx.Commit()

	return nil
}

//Delete is a function for deleting record
func (s *Sale) Delete(salesModel model.Model) *errors.Error {
	salesModelObj, ok := salesModel.(*model.Sales)
	if false == ok {
		return errors.Wrap(fmt.Errorf("Failed asserting to *model.Sales"), 0)
	}
	_, errs := s.FindByID(salesModel.GetID())
	if errs.Err == sql.ErrNoRows {
		return errors.Wrap(fmt.Errorf("cannot update, model with id: %v doesn't exist", salesModel.GetID()), 0)
	}

	//start db transaction
	//start db transaction
	if s.tx == nil {
		tx, err := s.db.Begin()
		if err != nil {
			return errors.Wrap(err, 0)
		}
		s.tx = tx
	}

	//delete items
	itemStmt, err := s.db.Prepare("DELETE FROM sales_items WHERE INVOICE_ID=?")
	if err != nil {
		s.tx.Rollback()
		return errors.Wrap(err, 0)
	}
	_, err = itemStmt.Exec(salesModelObj.InvoiceID)
	if err != nil {
		s.tx.Rollback()
		return errors.Wrap(err, 0)
	}

	//delete the model
	stmt, err := s.db.Prepare("DELETE FROM sales WHERE INVOICE_ID=?")
	if err != nil {
		s.tx.Rollback()
		return errors.Wrap(err, 0)
	}
	_, err = stmt.Exec(salesModelObj.InvoiceID)
	if err != nil {
		s.tx.Rollback()
		return errors.Wrap(err, 0)
	}
	s.tx.Commit()

	return nil
}

//Save is a function for persisting a model object to db
func (s *Sale) Save(salesModel model.Model) *errors.Error {
	var err *errors.Error
	if true == salesModel.GetLoadedFromStorage() {
		//update operation
		err = s.Update(salesModel)
	} else {
		//insert operation
		err = s.Insert(salesModel)
	}
	return err
}

//BeginTransaction starts a transaction on the connected session
func (s *Sale) BeginTransaction() *errors.Error {
	tx, err := s.db.Begin()
	if err != nil {
		return errors.Wrap(err, 0)
	}
	s.tx = tx
	return nil
}

//Commit commits the transaction
func (s *Sale) Commit() *errors.Error {
	if s.tx == nil {
		return errors.Wrap(fmt.Errorf("Can't commit, no transaction has been started"), 0)
	}
	err := s.tx.Commit()
	if err != nil {
		return errors.Wrap(err, 0)
	}
	s.tx = nil
	return nil
}

//Rollback cancels the transaction
func (s *Sale) Rollback() *errors.Error {
	if s.tx == nil {
		return errors.Wrap(fmt.Errorf("Can't rollback, no transaction has been started"), 0)
	}
	err := s.tx.Rollback()
	if err != nil {
		return errors.Wrap(err, 0)
	}
	s.tx = nil
	return nil
}

//StartUp allows the datamapper to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (s *Sale) StartUp() {
	//Note: Perform any initialization or bootstrapping here
}

//Shutdown allows the datamapper to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (s *Sale) Shutdown() {
	//Note: perform any cleanup here
}
