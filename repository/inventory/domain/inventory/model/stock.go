//Package model provides the domain model definitions
package model

//Stock is business domain model definition of item stock
type Stock struct {
	Sku               string
	Name              string
	Quantity          int64
	BuyPrice          float64
	SellPrice         float64
	loadedFromStorage bool //flag indicating whether the model object was loaded from storage or not
}

//GetID is a function for returning id of the model
func (s *Stock) GetID() string {
	return s.Sku
}

//GetLoadedFromStorage is a function for returning loaded from storage flag value
func (s *Stock) GetLoadedFromStorage() bool {
	return s.loadedFromStorage
}

//SetLoadedFromStorage is a function for setting loaded from storage flag value
func (s *Stock) SetLoadedFromStorage(flagValue bool) {
	s.loadedFromStorage = flagValue
}
