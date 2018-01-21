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

//datetime format used in converting datetime string to time object and vice versa
const timeFormat = "2006-01-02 15:04:05"
const dateFormat = "2006-01-02"

//Purchase is a struct of datamapper for purchase domain model
type Purchase struct {
	db *sql.DB
	tx *sql.Tx
}

//NewPurchase creates a new Purchase datamapper and returns a pointer to it
func NewPurchase(dbSession *sql.DB) *Purchase {
	return &Purchase{
		db: dbSession,
	}
}

//FindByID is a function for finding a record by id
func (p *Purchase) FindByID(id string) (model.Model, *errors.Error) {
	stmt, err := p.db.Prepare("SELECT PURCHASE_ID, DATETIME(PURCHASE_DATE), STATUS, NOTE FROM purchase WHERE PURCHASE_ID = ?")
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}
	defer stmt.Close()

	var purchaseID, date, status, note sql.NullString

	row := stmt.QueryRow(id)
	err = row.Scan(&purchaseID, &date, &status, &note)
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

	purchaseIDValue := purchaseID.String
	dateValue := date.String
	dateTimeValue, err := time.Parse(timeFormat, dateValue)

	statusValue := status.String
	noteValue := note.String

	purchaseModel := &model.Purchase{
		PurchaseID: purchaseIDValue,
		Date:       dateTimeValue,
		Status:     statusValue,
		Note:       noteValue,
	}
	purchaseModel.SetLoadedFromStorage(true)

	//load purchase items
	itemStmt, err := p.db.Prepare("SELECT ID, SKU, QUANTITY, BUY_PRICE, NOTE FROM purchase_items WHERE PURCHASE_ID = ?")
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}
	defer itemStmt.Close()

	var itemID int64
	var sku, itemNote sql.NullString
	var quantity sql.NullInt64
	var buyPrice sql.NullFloat64

	rows, err := itemStmt.Query(id)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}
	defer rows.Close()

	itemsRow := make(map[string]*model.PurchaseItem, 5)
	for rows.Next() {
		err := rows.Scan(&itemID, &sku, &quantity, &buyPrice, &itemNote)
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}
		skuValue := sku.String
		quantityValue := quantity.Int64
		buyPriceValue := buyPrice.Float64
		itemNoteValue := itemNote.String

		purchaseItemModel := &model.PurchaseItem{
			Sku:      skuValue,
			Quantity: quantityValue,
			BuyPrice: buyPriceValue,
			Note:     itemNoteValue,
		}
		purchaseItemModel.SetID(itemID)
		purchaseItemModel.SetLoadedFromStorage(true)

		itemsRow[skuValue] = purchaseItemModel
	}
	purchaseModel.Items = itemsRow

	return purchaseModel, nil
}

//FindAll is a function for finding all records
func (p *Purchase) FindAll() ([]model.Model, *errors.Error) {
	rows, err := p.db.Query("SELECT PURCHASE_ID, DATETIME(PURCHASE_DATE), STATUS, NOTE FROM purchase ORDER BY PURCHASE_ID ASC")
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}
	defer rows.Close()

	var purchaseID, date, status, note sql.NullString

	var itemID int64
	var sku, itemNote sql.NullString
	var quantity sql.NullInt64
	var buyPrice sql.NullFloat64

	var returnedRow []model.Model
	var firstScan = true
	for rows.Next() {
		err := rows.Scan(&purchaseID, &date, &status, &note)
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
		purchaseIDValue := purchaseID.String
		dateValue := date.String
		dateTimeValue, err := time.Parse(timeFormat, dateValue)
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}

		statusValue := status.String
		noteValue := note.String

		purchaseModel := &model.Purchase{
			PurchaseID: purchaseIDValue,
			Date:       dateTimeValue,
			Status:     statusValue,
			Note:       noteValue,
		}
		purchaseModel.SetLoadedFromStorage(true)

		//load purchase items
		itemStmt, err := p.db.Prepare("SELECT ID, SKU, QUANTITY, BUY_PRICE, NOTE FROM purchase_items WHERE PURCHASE_ID = ?")
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}
		defer itemStmt.Close()

		itemRows, err := itemStmt.Query(purchaseIDValue)
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}
		defer itemRows.Close()

		itemsRow := make(map[string]*model.PurchaseItem, 5)
		for itemRows.Next() {
			err := itemRows.Scan(&itemID, &sku, &quantity, &buyPrice, &itemNote)
			if err != nil {
				return nil, errors.Wrap(err, 0)
			}
			skuValue := sku.String
			quantityValue := quantity.Int64
			buyPriceValue := buyPrice.Float64
			itemNoteValue := itemNote.String

			purchaseItemModel := &model.PurchaseItem{
				Sku:      skuValue,
				Quantity: quantityValue,
				BuyPrice: buyPriceValue,
				Note:     itemNoteValue,
			}
			purchaseItemModel.SetID(itemID)
			purchaseItemModel.SetLoadedFromStorage(true)

			itemsRow[skuValue] = purchaseItemModel
		}
		purchaseModel.Items = itemsRow

		returnedRow = append(returnedRow, purchaseModel)
	}
	return returnedRow, nil
}

//Insert is a function for inserting a record
func (p *Purchase) Insert(purchaseModel model.Model) *errors.Error {
	purchaseModelObj, ok := purchaseModel.(*model.Purchase)
	if false == ok {
		return errors.Wrap(fmt.Errorf("Failed asserting to *model.Purchase"), 0)
	}

	foundModel, _ := p.FindByID(purchaseModel.GetID())
	if foundModel != nil {
		return errors.Wrap(fmt.Errorf("cannot insert, model with id: %v already exists", purchaseModel.GetID()), 0)
	}

	//start db transaction
	tx, err := p.db.Begin()
	if err != nil {
		return errors.Wrap(err, 0)
	}

	stmt, err := p.db.Prepare("INSERT INTO purchase(PURCHASE_ID, PURCHASE_DATE, STATUS, NOTE) values(?,?,?,?)")
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, 0)
	}
	dateString := purchaseModelObj.Date.Format(timeFormat)
	_, err = stmt.Exec(purchaseModelObj.PurchaseID, dateString, purchaseModelObj.Status, purchaseModelObj.Note)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, 0)
	}

	//insert the items
	for _, val := range purchaseModelObj.Items {
		itemStmt, err := p.db.Prepare("INSERT INTO purchase_items(PURCHASE_ID, SKU, QUANTITY, BUY_PRICE, NOTE) values(?,?,?,?,?)")
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, 0)
		}
		_, err = itemStmt.Exec(purchaseModelObj.PurchaseID, val.Sku, val.Quantity, val.BuyPrice, val.Note)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, 0)
		}
	}
	tx.Commit()

	return nil
}

//Update is a function for updating record
func (p *Purchase) Update(purchaseModel model.Model) *errors.Error {
	purchaseModelObj, ok := purchaseModel.(*model.Purchase)
	if false == ok {
		return errors.Wrap(fmt.Errorf("Failed asserting to *model.Purchase"), 0)
	}

	_, errs := p.FindByID(purchaseModel.GetID())
	if errs != nil && errs.Err == sql.ErrNoRows {
		return errors.Wrap(fmt.Errorf("cannot update, model with id: %v doesn't exist", purchaseModel.GetID()), 0)
	}

	//start db transaction
	tx, err := p.db.Begin()
	if err != nil {
		return errors.Wrap(err, 0)
	}

	stmt, err := p.db.Prepare("UPDATE purchase SET DATE=?, STATUS=?, NOTE=? WHERE PURCHASE_ID=?")
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, 0)
	}
	dateString := purchaseModelObj.Date.Format(timeFormat)
	_, err = stmt.Exec(dateString, purchaseModelObj.Status, purchaseModelObj.Note, purchaseModelObj.PurchaseID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, 0)
	}

	//update items
	for _, val := range purchaseModelObj.Items {
		var itemStmt *sql.Stmt
		if false == val.GetLoadedFromStorage() {
			itemStmt, err = p.db.Prepare("INSERT INTO purchase_items(PURCHASE_ID, SKU, QUANTITY, BUY_PRICE, NOTE) values(?,?,?,?,?)")
			if err != nil {
				tx.Rollback()
				return errors.Wrap(err, 0)
			}
			_, err = itemStmt.Exec(purchaseModelObj.PurchaseID, val.Sku, val.Quantity, val.BuyPrice, val.Note)
			if err != nil {
				tx.Rollback()
				return errors.Wrap(err, 0)
			}
		} else {
			itemStmt, err = p.db.Prepare("UPDATE purchase_items SET QUANTITY=?, BUY_PRICE=?, NOTE=? WHERE ID=?")
			if err != nil {
				tx.Rollback()
				return errors.Wrap(err, 0)
			}
			_, err = itemStmt.Exec(val.Quantity, val.BuyPrice, val.Note, val.GetID())
			if err != nil {
				tx.Rollback()
				return errors.Wrap(err, 0)
			}
		}
	}
	tx.Commit()

	return nil
}

//Delete is a function for deleting record
func (p *Purchase) Delete(purchaseModel model.Model) *errors.Error {
	purchaseModelObj, ok := purchaseModel.(*model.Purchase)
	if false == ok {
		return errors.Wrap(fmt.Errorf("Failed asserting to *model.Purchase"), 0)
	}

	_, errs := p.FindByID(purchaseModel.GetID())
	if errs.Err == sql.ErrNoRows {
		return errors.Wrap(fmt.Errorf("cannot delete, model with id: %v doesn't exist", purchaseModel.GetID()), 0)
	}

	//start db transaction
	tx, err := p.db.Begin()
	if err != nil {
		return errors.Wrap(err, 0)
	}

	//delete items
	itemStmt, err := p.db.Prepare("DELETE FROM purchase_items WHERE PURCHASE_ID=?")
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, 0)
	}
	_, err = itemStmt.Exec(purchaseModelObj.PurchaseID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, 0)
	}

	//delete the model
	stmt, err := p.db.Prepare("DELETE FROM purchase WHERE PURCHASE_ID=?")
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, 0)
	}
	_, err = stmt.Exec(purchaseModelObj.PurchaseID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, 0)
	}
	tx.Commit()

	return nil
}

//Save is a function for persisting a model object to db
func (p *Purchase) Save(purchaseModel model.Model) *errors.Error {
	var err *errors.Error
	if true == purchaseModel.GetLoadedFromStorage() {
		//update operation
		err = p.Update(purchaseModel)
	} else {
		//insert operation
		err = p.Insert(purchaseModel)
	}
	return err
}

//BeginTransaction starts a transaction on the connected session
func (p *Purchase) BeginTransaction() *errors.Error {
	tx, err := p.db.Begin()
	if err != nil {
		return errors.Wrap(err, 0)
	}
	p.tx = tx
	return nil
}

//Commit commits the transaction
func (p *Purchase) Commit() *errors.Error {
	if p.tx == nil {
		return errors.Wrap(fmt.Errorf("Can't commit, no transaction has been started"), 0)
	}
	err := p.tx.Commit()
	if err != nil {
		return errors.Wrap(err, 0)
	}
	p.tx = nil
	return nil
}

//Rollback cancels the transaction
func (p *Purchase) Rollback() *errors.Error {
	if p.tx == nil {
		return errors.Wrap(fmt.Errorf("Can't rollback, no transaction has been started"), 0)
	}
	err := p.tx.Rollback()
	if err != nil {
		return errors.Wrap(err, 0)
	}
	p.tx = nil
	return nil
}

//StartUp allows the datamapper to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (p *Purchase) StartUp() {
	//Note: Perform any initialization or bootstrapping here
}

//Shutdown allows the datamapper to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (p *Purchase) Shutdown() {
	//Note: perform any cleanup here
}
