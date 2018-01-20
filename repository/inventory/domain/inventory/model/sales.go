//Package model provides the domain model definitions
package model

import (
	"time"
)

//SalesStatusDraft is const for 'draft' Sales status
const SalesStatusDraft string = "D"

//SalesStatusCanceled is const for 'canceled' Sales status
const SalesStatusCanceled string = "C"

//SalesStatusDone is const for 'Success/Done' Sales status
const SalesStatusDone string = "S"

//Sales is business domain model definition of a sale
type Sales struct {
	InvoiceID         string
	Date              time.Time
	Status            string
	Note              string
	Items             map[string]*SaleItem
	loadedFromStorage bool //flag indicating whether the model object was loaded from storage or not
}

//GetID is a function for returning id of the model
func (s *Sales) GetID() string {
	return s.InvoiceID
}

//GetLoadedFromStorage is a function for returning loaded from storage flag value
func (s *Sales) GetLoadedFromStorage() bool {
	return s.loadedFromStorage
}

//SetLoadedFromStorage is a function for setting loaded from storage flag value
func (s *Sales) SetLoadedFromStorage(flagValue bool) {
	s.loadedFromStorage = flagValue
}

//SaleItem is a business domain model definition of a sale item
type SaleItem struct {
	id                int64
	Sku               string
	Quantity          int64
	BuyPrice          float64
	SellPrice         float64
	loadedFromStorage bool //flag indicating whether the model object was loaded from storage or not
}

//GetID is a function for returning id of the model
func (si *SaleItem) GetID() int64 {
	return si.id
}

//SetID is a function for setting id of the model
func (si *SaleItem) SetID(id int64) {
	si.id = id
}

//GetLoadedFromStorage is a function for returning loaded from storage flag value
func (si *SaleItem) GetLoadedFromStorage() bool {
	return si.loadedFromStorage
}

//SetLoadedFromStorage is a function for setting loaded from storage flag value
func (si *SaleItem) SetLoadedFromStorage(flagValue bool) {
	si.loadedFromStorage = flagValue
}
