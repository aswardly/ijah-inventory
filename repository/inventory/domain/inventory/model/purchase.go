//Package model provides the domain model definitions
package model

import (
	"time"
)

//PurchaseStatusDraft is const for 'draft' purchase status
const PurchaseStatusDraft string = "D"

//PurchaseStatusCanceled is const for 'canceled' purchase status
const PurchaseStatusCanceled string = "C"

//PurchaseStatusDone is const for 'Success/Done' purchase status
const PurchaseStatusDone string = "S"

//Purchase is business domain model definition of a purchase
type Purchase struct {
	PurchaseID        string
	Date              time.Time
	Status            string
	Note              string
	Items             map[string]*PurchaseItem
	loadedFromStorage bool //flag indicating whether the model object was loaded from storage or not
}

//GetID is a function for returning id of the model
func (p *Purchase) GetID() string {
	return p.PurchaseID
}

//GetLoadedFromStorage is a function for returning loaded from storage flag value
func (p *Purchase) GetLoadedFromStorage() bool {
	return p.loadedFromStorage
}

//SetLoadedFromStorage is a function for setting loaded from storage flag value
func (p *Purchase) SetLoadedFromStorage(flagValue bool) {
	p.loadedFromStorage = flagValue
}

//PurchaseItem is a business domain model definition of a purchase item
type PurchaseItem struct {
	id                int64
	Sku               string
	Quantity          int64
	BuyPrice          float64
	Note              string
	loadedFromStorage bool //flag indicating whether the model object was loaded from storage or not
}

//GetID is a function for returning id of the model
func (pi *PurchaseItem) GetID() int64 {
	return pi.id
}

//SetID is a function for setting id of the model
func (pi *PurchaseItem) SetID(id int64) {
	pi.id = id
}

//GetLoadedFromStorage is a function for returning loaded from storage flag value
func (pi *PurchaseItem) GetLoadedFromStorage() bool {
	return pi.loadedFromStorage
}

//SetLoadedFromStorage is a function for setting loaded from storage flag value
func (pi *PurchaseItem) SetLoadedFromStorage(flagValue bool) {
	pi.loadedFromStorage = flagValue
}
