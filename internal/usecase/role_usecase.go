package usecase

import (
	"context"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/config"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/repositories"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/interfaces/http/dto"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/interfaces/http/mapper"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/helper"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/translations"
	"gorm.io/gorm"
)

type RoleUsecase interface {
	FindAll(ctx context.Context) ([]dto.Role, error)
	FindAllWithPagination(ctx context.Context, limit, offset int) (dto.PaginationData, error)
}

type RoleUsecaseImpl struct {
	Env            *config.Env
	DB             *gorm.DB
	Tr             *translations.Store
	RoleRepository repositories.RoleRepository
}

func NewRoleUsecase(env *config.Env, db *gorm.DB, tr *translations.Store, repo repositories.RoleRepository) RoleUsecase {
	return &RoleUsecaseImpl{
		Env:            env,
		DB:             db,
		Tr:             tr,
		RoleRepository: repo,
	}
}

func (r *RoleUsecaseImpl) FindAll(ctx context.Context) ([]dto.Role, error) {
	res, err := helper.WithGormTransaction(ctx, r.DB, func(tx *gorm.DB) (interface{}, error) {
		roles, err := r.RoleRepository.List(ctx, tx, 1000, 0)
		if err != nil {
			return nil, err
		}
		return mapper.ToRoleResponses(roles, true), nil
	})
	if err != nil {
		return nil, err
	}
	return res.([]dto.Role), nil
}

func (r *RoleUsecaseImpl) FindAllWithPagination(ctx context.Context, limit, offset int) (dto.PaginationData, error) {
	res, err := helper.WithGormTransaction(ctx, r.DB, func(tx *gorm.DB) (interface{}, error) {
		roles, err := r.RoleRepository.List(ctx, tx, limit, offset)
		if err != nil {
			return nil, err
		}
		total, err := r.RoleRepository.Count(ctx, tx)
		if err != nil {
			return nil, err
		}
		return dto.PaginationData{
			TotalData: int(total),
			Data:      mapper.ToRoleResponses(roles, true),
		}, nil
	})
	if err != nil {
		return dto.PaginationData{}, err
	}
	return res.(dto.PaginationData), nil
}
