package helper

import (
	"context"
	"fmt"

	"github.com/alpinnz/go-rest-api-boilerplate/pkg/errors"
	"gorm.io/gorm"
)

func WithGormTransaction(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) (interface{}, error)) (interface{}, error) {
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, errors.NewDatabaseError(
			fmt.Sprintf("failed to begin transaction: %v\n", tx.Error),
			tx.Error,
		)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
	}()

	res, err := fn(tx)
	if err != nil {
		_ = tx.Rollback()

		if be, ok := err.(*errors.BaseError); ok {
			return nil, be
		}

		// tampilkan pesan asli agar trace jelas
		return nil, errors.NewInternalError(
			fmt.Sprintf("transaction failed: %v\n", err),
			err,
		)
	}

	if commitErr := tx.Commit().Error; commitErr != nil {
		return nil, errors.NewDatabaseError(
			fmt.Sprintf("commit failed: %v\n", commitErr),
			commitErr,
		)
	}

	return res, nil
}
