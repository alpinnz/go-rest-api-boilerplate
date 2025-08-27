package errors

import (
	"errors"

	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
	"gorm.io/gorm"

	"github.com/lib/pq" // Postgres
)

// BaseGormErrorType defines custom GORM/Postgres errors types
type BaseGormErrorType string

const (
	ErrNotFound        BaseGormErrorType = "NotFound"
	ErrDuplicateKey    BaseGormErrorType = "DuplicateKey"
	ErrForeignKey      BaseGormErrorType = "ForeignKeyViolation"
	ErrCheckConstraint BaseGormErrorType = "CheckConstraint"
	ErrNotNull         BaseGormErrorType = "NotNull"
	ErrTableNotExist   BaseGormErrorType = "TableNotExist"
	ErrDBConnection    BaseGormErrorType = "DBConnection"
	ErrTransaction     BaseGormErrorType = "Transaction"
	ErrInvalidData     BaseGormErrorType = "InvalidData"
	ErrOther           BaseGormErrorType = "Other"
)

// invalidDataErrorsMap = GORM internal errors related to invalid data/query
var invalidDataErrorsMap = map[error]struct{}{
	gorm.ErrModelValueRequired:            {},
	gorm.ErrPrimaryKeyRequired:            {},
	gorm.ErrModelAccessibleFieldsRequired: {},
	gorm.ErrSubQueryRequired:              {},
	gorm.ErrInvalidData:                   {},
	gorm.ErrUnsupportedDriver:             {},
	gorm.ErrEmptySlice:                    {},
	gorm.ErrDryRunModeUnsupported:         {},
	gorm.ErrInvalidDB:                     {},
	gorm.ErrInvalidValue:                  {},
	gorm.ErrInvalidValueOfLength:          {},
	gorm.ErrPreloadNotAllowed:             {},
}

// HandleGormError classifies GORM and PostgreSQL errors
func HandleGormError(err error) (BaseGormErrorType, error) {
	if err == nil {
		return "", nil
	}

	// 1️⃣ GORM internal errors
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrNotFound, err
	case errors.Is(err, gorm.ErrInvalidTransaction):
		return ErrTransaction, err
	default:
		if _, ok := invalidDataErrorsMap[err]; ok {
			return ErrInvalidData, err
		}
	}

	// 2️⃣ PostgreSQL-specific errors
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return mapPostgresError(pqErr), err
	}

	// 3️⃣ Fallback for unknown errors
	return ErrOther, err
}

// mapPostgresError maps pq.Error codes to BaseGormErrorType
func mapPostgresError(pqErr *pq.Error) BaseGormErrorType {
	switch pqErr.Code.Name() {
	case "unique_violation":
		return ErrDuplicateKey
	case "foreign_key_violation":
		return ErrForeignKey
	case "check_violation":
		return ErrCheckConstraint
	case "not_null_violation":
		return ErrNotNull
	case "undefined_table":
		return ErrTableNotExist
	case "connection_exception", "connection_does_not_exist":
		return ErrDBConnection
	case "serialization_failure", "deadlock_detected":
		return ErrTransaction
	default:
		return ErrOther
	}
}

// ConvertGormError converts BaseGormErrorType into BaseError (user-friendly)
func ConvertGormError(baseType BaseGormErrorType, err error) BaseError {
	switch baseType {
	case ErrNotFound:
		return *NewNotFound("Data tidak ditemukan", err.Error())
	case ErrDuplicateKey:
		return *NewBadRequest("Data sudah ada", err.Error())
	case ErrForeignKey:
		return *NewBadRequest("Data terkait tidak ditemukan", err.Error())
	case ErrCheckConstraint:
		return *NewBadRequest("Data tidak valid", err.Error())
	case ErrNotNull:
		return *NewBadRequest("Field wajib diisi", err.Error())
	case ErrTableNotExist:
		return *NewBadRequest("Data yang diminta tidak ditemukan", err.Error())
	case ErrTransaction:
		return *NewInternalError("Terjadi kesalahan saat memproses data", err.Error())
	case ErrDBConnection:
		return *NewInternalError("Gagal menghubungkan ke database", err.Error())
	case ErrInvalidData:
		return *NewBadRequest("Data tidak valid", err.Error())
	default:
		return *NewInternalError("Terjadi kesalahan database", err.Error())
	}
}

// HandleRepoError wraps any errors from GORM/Postgres into BaseError
func HandleRepoError[T any](result T, err error) (T, error) {
	if err == nil {
		return result, nil
	}

	// klasifikasi errors GORM / Postgres
	baseType, _ := HandleGormError(err)
	baseErr := ConvertGormError(baseType, err)

	return result, &baseErr
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	// Cek kalau error dari GORM
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}

	// Cek kalau error custom kita
	var be *BaseError
	if errors.As(err, &be) && be.Code == response.CodeNotFound {
		return true
	}

	return false
}
