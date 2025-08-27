package utils

import (
	"gorm.io/gorm"
)

func QueryParamsApply(db *gorm.DB, params QueryParams) *gorm.DB {
	if params.Filter != "" {
		db = db.Where("name LIKE ?", "%"+params.Filter+"%") // Ganti "name" dengan kolom yang sesuai
	}

	if params.Limit > 0 {
		db = db.Limit(params.Limit)
	}

	if params.Offset > 0 {
		db = db.Offset(params.Offset)
	}

	if params.Sort != "" {
		db = db.Order(params.Sort)
	}

	return db
}
