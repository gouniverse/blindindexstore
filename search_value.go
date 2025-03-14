package blindindexstore

import (
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
)

// == CLASS ==================================================================

type SearchValue struct {
	dataobject.DataObject
}

// == CONSTRUCTORS ===========================================================

func NewSearchValue() *SearchValue {
	d := (&SearchValue{}).
		SetID(uid.HumanUid()).
		SetSourceReferenceID("").
		SetSearchValue("").
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetDeletedAt(sb.MAX_DATETIME)

	return d
}

func NewSearchValueFromExistingData(data map[string]string) *SearchValue {
	o := &SearchValue{}
	o.Hydrate(data)
	return o
}

// == METHODS ================================================================

// == SETTERS AND GETTERS ====================================================
func (d *SearchValue) CreatedAt() string {
	return d.Get(COLUMN_CREATED_AT)
}

func (d *SearchValue) CreatedAtCarbon() carbon.Carbon {
	createdAt := d.CreatedAt()
	return carbon.Parse(createdAt)
}

func (d *SearchValue) SetCreatedAt(createdAt string) *SearchValue {
	d.Set(COLUMN_CREATED_AT, createdAt)
	return d
}

func (d *SearchValue) DeletedAt() string {
	return d.Get(COLUMN_DELETED_AT)
}

func (d *SearchValue) DeletedAtCarbon() carbon.Carbon {
	deletedAt := d.DeletedAt()
	return carbon.Parse(deletedAt)
}

func (d *SearchValue) SetDeletedAt(deletedAt string) *SearchValue {
	d.Set(COLUMN_DELETED_AT, deletedAt)
	return d
}

// ID returns the ID of the exam
func (o *SearchValue) ID() string {
	return o.Get(COLUMN_ID)
}

// SetID sets the ID of the exam
func (o *SearchValue) SetID(id string) *SearchValue {
	o.Set(COLUMN_ID, id)
	return o
}

func (d *SearchValue) SourceReferenceID() string {
	return d.Get(COLUMN_SOURCE_REFERENCE_ID)
}

func (d *SearchValue) SetSourceReferenceID(objectID string) *SearchValue {
	d.Set(COLUMN_SOURCE_REFERENCE_ID, objectID)
	return d
}

func (d *SearchValue) SearchValue() string {
	return d.Get(COLUMN_SEARCH_VALUE)
}

func (d *SearchValue) SetSearchValue(value string) *SearchValue {
	d.Set(COLUMN_SEARCH_VALUE, value)
	return d
}

func (d *SearchValue) UpdatedAt() string {
	return d.Get(COLUMN_UPDATED_AT)
}

func (d *SearchValue) UpdatedAtCarbon() carbon.Carbon {
	updatedAt := d.UpdatedAt()
	return carbon.Parse(updatedAt)
}

func (d *SearchValue) SetUpdatedAt(updatedAt string) *SearchValue {
	d.Set(COLUMN_UPDATED_AT, updatedAt)
	return d
}
