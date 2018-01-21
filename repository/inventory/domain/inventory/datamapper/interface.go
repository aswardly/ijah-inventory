//Package datamapper provides the definitions of data mapper
package datamapper

import (
	"ijah-inventory/repository/inventory/domain/inventory/model"
	"time"

	"github.com/go-errors/errors"
)

//DataMapper is an interface for common data mapper
type DataMapper interface {
	FindByID(id string) (model.Model, *errors.Error)
	FindAll() ([]model.Model, *errors.Error)
	Insert(model.Model) *errors.Error
	Update(model.Model) *errors.Error
	Delete(model.Model) *errors.Error
	Save(model.Model) *errors.Error
	BeginTransaction() *errors.Error
	Commit() *errors.Error
	Rollback() *errors.Error
}

//SalesDataMapper is an interface for data mapper specific for sales
type SalesDataMapper interface {
	DataMapper

	FindByDoneStatusAndDateRange(time.Time, time.Time) ([]model.Model, *errors.Error)
}
