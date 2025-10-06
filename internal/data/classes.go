package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/data/pgdb"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/internal/domain/services/class"
)

func (d Database) CreateClass(ctx context.Context, class models.Class) error {
	return d.sql.classes.Insert(ctx, classModelToSchema(class))
}

func (d Database) GetClassByCode(ctx context.Context, code string) (models.Class, error) {
	schema, err := d.sql.classes.New().FilterCode(code).Get(ctx)
	switch {
	case err == sql.ErrNoRows:
		return models.Class{}, nil
	case err != nil:
		return models.Class{}, err
	}

	return classSchemaToModel(schema), nil
}

func (d Database) ClassIsExistByCode(ctx context.Context, code string) (bool, error) {
	return d.sql.classes.New().FilterCode(code).Exists(ctx)
}
func (d Database) ClassIsExistByName(ctx context.Context, name string) (bool, error) {
	return d.sql.classes.New().FilterName(name).Exists(ctx)
}

func (d Database) CountClassChildren(ctx context.Context, parentCode string) (uint64, error) {
	return d.sql.classes.New().FilterParent(sql.NullString{
		String: parentCode,
		Valid:  true,
	}).Count(ctx)
}

func (d Database) CountPlacesByClass(ctx context.Context, classCode string) (uint64, error) {
	return d.sql.places.New().FilterClass(classCode).Count(ctx)
}

func (d Database) UpdateClass(ctx context.Context, code string, params class.UpdateParams, updateAt time.Time) error {
	query := d.sql.classes.New()

	if params.Name != nil {
		query = query.UpdateName(*params.Name)
	}
	if params.Icon != nil {
		query = query.UpdateIcon(*params.Icon)
	}
	if params.Parent != nil {
		if *params.Parent == "" {
			query = query.UpdateParent(sql.NullString{Valid: false})
		} else {
			query = query.UpdateParent(sql.NullString{String: *params.Parent, Valid: true})
		}
	}

	return query.Update(ctx, updateAt)
}
func (d Database) UpdateClassStatus(ctx context.Context, code string, status string, updateAt time.Time) error {
	return d.sql.classes.New().FilterCode(code).UpdateStatus(status).Update(ctx, updateAt)
}
func (d Database) ReplaceClassInPlaces(ctx context.Context, oldCode, newCode string, updateAt time.Time) error {
	return d.sql.places.New().FilterClass(oldCode).UpdateClass(newCode).Update(ctx, updateAt)
}

func (d Database) FilterClasses(
	ctx context.Context,
	filter class.FilterParams,
	page, size uint64,
) (models.ClassesCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	query := d.sql.classes.New()

	if filter.Status != nil {
		query = query.FilterStatus(*filter.Status)
	}
	if filter.Parent != nil {
		if filter.ParentCycle {
			query = query.FilterParentCycle(*filter.Parent)
		} else {
			if *filter.Parent == "" {
				query = query.FilterParent(sql.NullString{Valid: false})
			} else {
				query = query.FilterParent(sql.NullString{String: *filter.Parent, Valid: true})
			}
		}
	}

	total, err := query.Count(ctx)
	if err != nil {
		return models.ClassesCollection{}, err
	}

	rows, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return models.ClassesCollection{}, err
	}

	collection := make([]models.Class, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, classSchemaToModel(row))
	}

	return models.ClassesCollection{
		Data:  collection,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func (d Database) CheckParentCycle(ctx context.Context, classCode, parentCode string) (bool, error) {
	return d.sql.classes.New().FilterCode(classCode).FilterParentCycle(parentCode).Exists(ctx)
}

func (d Database) DeleteClass(ctx context.Context, code string) error {
	return d.sql.classes.New().FilterCode(code).Delete(ctx)
}

func classModelToSchema(model models.Class) pgdb.Class {
	res := pgdb.Class{
		Code:      model.Code,
		Name:      model.Name,
		Status:    model.Status,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}

	if model.Parent != nil {
		res.Parent = sql.NullString{String: *model.Parent, Valid: true}
	}

	return res
}

func classSchemaToModel(schema pgdb.Class) models.Class {
	res := models.Class{
		Code:      schema.Code,
		Name:      schema.Name,
		Status:    schema.Status,
		CreatedAt: schema.CreatedAt,
		UpdatedAt: schema.UpdatedAt,
	}

	if schema.Parent.Valid {
		res.Parent = &schema.Parent.String
	}

	return res
}
