package models

import (
	"github.com/go-playground/validator/v10"
)

type Model interface{
	UniqueID() string
	Validate() error
}

type ReadingItem struct{
	ID string `json:"id" fauna:"id" validate:"required"`
	Title string `json:"title" fauna:"title" validate:"required"`
	Link string `json:"link" fauna:"link" validate:"omitempty,url"`
	Type string `json:"type" fauna:"type" validate:"required"`
	Author string `json:"author" fauna:"author" validate:"required"`
}

func NewReadingItem(id, title, link, itemType, author string) ReadingItem {
	return ReadingItem{
		id,
		title,
		link,
		itemType,
		author,
	}
}

func (r ReadingItem) Validate() error {
	return validator.New().Struct(r)
}

func (r ReadingItem) UniqueID() string {
	return r.ID
}





