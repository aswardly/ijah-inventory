//Package datamapper provides the definitions of datamapper
package datamapper

import (
	"database/sql"
	"fmt"

	"github.com/go-errors/errors"
	_ "github.com/mattn/go-sqlite3"

	"ijah-inventory/repository/inventory/domain/inventory/model"
)

//Custom error used for masking error type from specific sql driver
var (
	ErrNotFound = fmt.Errorf("Record not found")
)

//Stock is a struct of datamapper for stock domain model
type Stock struct {
	db *sql.DB
}

//NewStock creates a new Stock datamapper and returns a pointer to it
func NewStock(dbSession *sql.DB) *Stock {
	return &Stock{
		db: dbSession,
	}
}

//FindByID is a function for finding a record by id
func (s *Stock) FindByID(id string) (model.Model, *errors.Error) {
	stmt, err := s.db.Prepare("SELECT SKU, NAME, QUANTITY, BUY_PRICE, SELL_PRICE FROM stock WHERE SKU = ?")

	if err != nil {
		return nil, errors.Wrap(err, 0)
	}
	defer stmt.Close()

	var sku, name sql.NullString
	var quantity sql.NullInt64
	var buyPrice, sellPrice sql.NullFloat64

	row := stmt.QueryRow(id)
	err = row.Scan(&sku, &name, &quantity, &buyPrice, &sellPrice)
	if err != nil {
		var returnedErr error
		if err == sql.ErrNoRows {
			returnedErr = ErrNotFound
		} else {
			returnedErr = err
		}
		return nil, errors.Wrap(returnedErr, 0)
	}

	skuValue := sku.String
	nameValue := name.String
	quantityValue := quantity.Int64
	buyPriceValue := buyPrice.Float64
	sellPriceValue := sellPrice.Float64

	stockModel := &model.Stock{
		Sku:       skuValue,
		Name:      nameValue,
		Quantity:  quantityValue,
		BuyPrice:  buyPriceValue,
		SellPrice: sellPriceValue,
	}
	stockModel.SetLoadedFromStorage(true)

	return stockModel, nil
}

//FindAll is a function for finding all records
func (s *Stock) FindAll() ([]model.Model, *errors.Error) {
	rows, err := s.db.Query("SELECT SKU, NAME, QUANTITY, BUY_PRICE, SELL_PRICE FROM stock ORDER BY SKU ASC")
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}
	defer rows.Close()

	var sku, name sql.NullString
	var quantity sql.NullInt64
	var buyPrice, sellPrice sql.NullFloat64

	var returnedRow []model.Model
	var firstScan = true
	for rows.Next() {
		err := rows.Scan(&sku, &name, &quantity, &buyPrice, &sellPrice)
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
		skuValue := sku.String
		nameValue := name.String
		quantityValue := quantity.Int64
		buyPriceValue := buyPrice.Float64
		sellPriceValue := sellPrice.Float64

		stockModel := &model.Stock{
			Sku:       skuValue,
			Name:      nameValue,
			Quantity:  quantityValue,
			BuyPrice:  buyPriceValue,
			SellPrice: sellPriceValue,
		}
		stockModel.SetLoadedFromStorage(true)

		returnedRow = append(returnedRow, stockModel)
	}
	return returnedRow, nil
}

//Insert is a function for inserting a record
func (s *Stock) Insert(stockModel model.Model) *errors.Error {
	stockModelObj, ok := stockModel.(*model.Stock)
	if false == ok {
		return errors.Wrap(fmt.Errorf("Failed asserting to *model.Stock"), 0)
	}
	foundModel, _ := s.FindByID(stockModel.GetID())
	if foundModel != nil {
		return errors.Wrap(fmt.Errorf("cannot insert, model with id: %v already exists", stockModel.GetID()), 0)
	}
	stmt, err := s.db.Prepare("INSERT INTO stock(SKU, NAME, QUANTITY, BUY_PRICE, SELL_PRICE) values(?,?,?,?,?)")
	if err != nil {
		return errors.Wrap(err, 0)
	}
	_, err = stmt.Exec(stockModelObj.Sku, stockModelObj.Name, stockModelObj.Quantity, stockModelObj.BuyPrice, stockModelObj.SellPrice)
	if err != nil {
		return errors.Wrap(err, 0)
	}
	return nil
}

//Update is a function for updating record
func (s *Stock) Update(stockModel model.Model) *errors.Error {
	stockModelObj, ok := stockModel.(*model.Stock)
	if false == ok {
		return errors.Wrap(fmt.Errorf("Failed asserting to *model.Stock"), 0)
	}
	_, errs := s.FindByID(stockModel.GetID())
	if errs != nil && errs.Err == sql.ErrNoRows {
		return errors.Wrap(fmt.Errorf("cannot update, model with id: %v doesn't exist", stockModel.GetID()), 0)
	}
	stmt, err := s.db.Prepare("UPDATE stock SET SKU=?, NAME=?, QUANTITY=?, BUY_PRICE=?, SELL_PRICE=? WHERE SKU=?")
	if err != nil {
		return errors.Wrap(err, 0)
	}
	_, err = stmt.Exec(stockModelObj.Sku, stockModelObj.Name, stockModelObj.Quantity, stockModelObj.BuyPrice, stockModelObj.SellPrice, stockModelObj.Sku)
	if err != nil {
		return errors.Wrap(err, 0)
	}
	return nil
}

//Delete is a function for deleting record
func (s *Stock) Delete(stockModel model.Model) *errors.Error {
	stockModelObj, ok := stockModel.(*model.Stock)
	if false == ok {
		return errors.Wrap(fmt.Errorf("Failed asserting to *model.Stock"), 0)
	}
	_, errs := s.FindByID(stockModel.GetID())
	if errs.Err == sql.ErrNoRows {
		return errors.Wrap(fmt.Errorf("cannot delete, model with id: %v doesn't exist", stockModel.GetID()), 0)
	}
	stmt, err := s.db.Prepare("DELETE FROM stock WHERE SKU=?")
	if err != nil {
		return errors.Wrap(err, 0)
	}
	_, err = stmt.Exec(stockModelObj.Sku)
	if err != nil {
		return errors.Wrap(err, 0)
	}
	return nil
}

//Save is a function for persisting a model object to db
func (s *Stock) Save(stockModel model.Model) *errors.Error {
	var err *errors.Error
	if true == stockModel.GetLoadedFromStorage() {
		//update operation
		err = s.Update(stockModel)
	} else {
		//insert operation
		err = s.Insert(stockModel)
	}
	return err
}

//StartUp allows the datamapper to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (s *Stock) StartUp() {
	//Note: Perform any initialization or bootstrapping here
}

//Shutdown allows the datamapper to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (s *Stock) Shutdown() {
	//Note: perform any cleanup here
}
