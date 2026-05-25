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

func (r *BaseRepository[T]) DB() *gorm.DB {
	return r.db
}

func (r *BaseRepository[T]) Create(
	ctx context.Context,
	entity *T,
) error {
	return r.db.
		WithContext(ctx).
		Create(entity).
		Error
}

func (r *BaseRepository[T]) CreateInBatches(
	ctx context.Context,
	entities *[]T,
	batchSize int,
) error {
	return r.db.
		WithContext(ctx).
		CreateInBatches(entities, batchSize).
		Error
}

func (r *BaseRepository[T]) Update(
	ctx context.Context,
	entity *T,
) error {
	return r.db.
		WithContext(ctx).
		Save(entity).
		Error
}

func (r *BaseRepository[T]) Delete(
	ctx context.Context,
	entity *T,
) error {
	return r.db.
		WithContext(ctx).
		Delete(entity).
		Error
}

func (r *BaseRepository[T]) DeleteByID(
	ctx context.Context,
	id uint,
) error {
	return r.db.
		WithContext(ctx).
		Delete(new(T), id).
		Error
}

func (r *BaseRepository[T]) FindByID(
	ctx context.Context,
	id uint,
) (*T, error) {
	var entity T

	err := r.db.
		WithContext(ctx).
		First(&entity, id).
		Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *BaseRepository[T]) FindOne(
	ctx context.Context,
	query any,
	args ...any,
) (*T, error) {
	var entity T

	err := r.db.
		WithContext(ctx).
		Where(query, args...).
		First(&entity).
		Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *BaseRepository[T]) FindAll(
	ctx context.Context,
	query any,
	args ...any,
) ([]T, error) {
	var entities []T

	err := r.db.
		WithContext(ctx).
		Where(query, args...).
		Find(&entities).
		Error

	return entities, err
}

func (r *BaseRepository[T]) Exists(
	ctx context.Context,
	query any,
	args ...any,
) (bool, error) {
	var count int64

	err := r.db.
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

func (r *BaseRepository[T]) Count(
	ctx context.Context,
	query any,
	args ...any,
) (int64, error) {
	var count int64

	err := r.db.
		WithContext(ctx).
		Model(new(T)).
		Where(query, args...).
		Count(&count).
		Error

	return count, err
}

func (r *BaseRepository[T]) UpdateColumns(
	ctx context.Context,
	model *T,
	values map[string]any,
) error {
	return r.db.
		WithContext(ctx).
		Model(model).
		Updates(values).
		Error
}

func (r *BaseRepository[T]) FirstOrCreate(
	ctx context.Context,
	entity *T,
	conditions any,
) error {
	return r.db.
		WithContext(ctx).
		Where(conditions).
		FirstOrCreate(entity).
		Error
}

func (r *BaseRepository[T]) Transaction(
	ctx context.Context,
	fn func(tx *gorm.DB) error,
) error {
	return r.db.
		WithContext(ctx).
		Transaction(fn)
}
