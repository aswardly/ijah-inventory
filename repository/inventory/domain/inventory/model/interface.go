//Package model provides the definitions of domain model
package model

//Model is an interface of a domain model
type Model interface {
	GetID() string
	GetLoadedFromStorage() bool
	SetLoadedFromStorage(bool)
}
