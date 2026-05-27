package repositories

import (
	"context"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{
		db: db,
	}
}

func (repo *BaseRepository[T]) DB() *gorm.DB {
	return repo.db
}

func (repo *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return repo.db.
		WithContext(ctx).
		Create(entity).
		Error
}

func (repo *BaseRepository[T]) CreateInBatches(ctx context.Context, entities *[]T, batchSize int) error {
	return repo.db.
		WithContext(ctx).
		CreateInBatches(entities, batchSize).
		Error
}

func (repo *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return repo.db.
		WithContext(ctx).
		Save(entity).
		Error
}

func (repo *BaseRepository[T]) Delete(ctx context.Context, entity *T) error {
	return repo.db.
		WithContext(ctx).
		Delete(entity).
		Error
}

func (repo *BaseRepository[T]) DeleteByID(ctx context.Context, id uint) error {
	return repo.db.
		WithContext(ctx).
		Delete(new(T), id).
		Error
}

func (repo *BaseRepository[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	var entity T

	err := repo.db.
		WithContext(ctx).
		First(&entity, id).
		Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (repo *BaseRepository[T]) FindOne(ctx context.Context, query any, args ...any) (*T, error) {
	var entity T

	err := repo.db.
		WithContext(ctx).
		Where(query, args...).
		First(&entity).
		Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (repo *BaseRepository[T]) FindAll(ctx context.Context, query any, args ...any) ([]T, error) {
	var entities []T

	err := repo.db.
		WithContext(ctx).
		Where(query, args...).
		Find(&entities).
		Error

	return entities, err
}

func (repo *BaseRepository[T]) Exists(ctx context.Context, query any, args ...any) (bool, error) {
	var count int64

	err := repo.db.
		WithContext(ctx).
		Model(new(T)).
		Where(query, args...).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *BaseRepository[T]) Count(ctx context.Context, query any, args ...any) (int64, error) {
	var count int64

	err := repo.db.
		WithContext(ctx).
		Model(new(T)).
		Where(query, args...).
		Count(&count).
		Error

	return count, err
}

func (repo *BaseRepository[T]) UpdateColumns(ctx context.Context, model *T, values map[string]any) error {
	return repo.db.
		WithContext(ctx).
		Model(model).
		Updates(values).
		Error
}

func (repo *BaseRepository[T]) FirstOrCreate(ctx context.Context, entity *T, conditions any) error {
	return repo.db.
		WithContext(ctx).
		Where(conditions).
		FirstOrCreate(entity).
		Error
}

func (repo *BaseRepository[T]) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return repo.db.
		WithContext(ctx).
		Transaction(fn)
}
