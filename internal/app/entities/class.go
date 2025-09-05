package entities

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/chains-lab/places-svc/internal/dbx"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/pkg/errors"
)

type Classificator struct {
	query   dbx.ClassesQ
	localeQ dbx.ClassLocaleQ
}

func NewClassificator(db *sql.DB) Classificator {
	return Classificator{
		query:   dbx.NewClassesQ(db),
		localeQ: dbx.NewClassLocaleQ(db),
	}
}

type CreateClassParams struct {
	Code   string
	Father *string
	Icon   string
	Path   string
	Locale ClassLocaleCreateParams
}

type ClassLocaleCreateParams struct {
	Locale string
	Name   string
}

func (c Classificator) CreateClass(ctx context.Context, params CreateClassParams) (models.PlaceClass, error) {
	_, err := c.query.New().FilterCode(params.Code).Get(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.PlaceClass{}, errx.ErrorInternal.Raise(fmt.Errorf("failed to check class existence, cause: %w", err))
	}
	if err == nil {
		return models.PlaceClass{}, errx.ErrorClassAlreadyExists.Raise(
			fmt.Errorf("class with code %s already exists", params.Code),
		)
	}

	fatherValue := sql.NullString{}
	if params.Father != nil {
		fatherValue = sql.NullString{String: *params.Father, Valid: true}
	}

	now := time.Now().UTC()

	err = c.query.New().Insert(ctx, dbx.InsertPlaceClass{
		Code:      params.Code,
		Father:    fatherValue,
		Status:    constant.PlaceClassStatusesInactive,
		Icon:      params.Icon,
		Path:      params.Path,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return models.PlaceClass{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create class, cause: %w", err),
		)
	}

	err = c.localeQ.Insert(ctx, dbx.PlaceClassLocale{
		Class:  params.Code,
		Locale: params.Locale.Locale,
		Name:   params.Locale.Name,
	})
	if err != nil {
		return models.PlaceClass{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create class locale, cause: %w", err),
		)
	}

	resp := models.PlaceClass{
		Code:      params.Code,
		Status:    constant.PlaceClassStatusesInactive,
		Icon:      params.Icon,
		Path:      params.Path,
		Name:      params.Locale.Name,
		Locale:    params.Locale.Locale,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if params.Father != nil {
		resp.Father = params.Father
	}

	return resp, nil
}

func (c Classificator) GetClass(
	ctx context.Context,
	code, locale string,
) (models.PlaceClass, error) {
	err := constant.IsLocaleSupported(locale)
	if err != nil {
		locale = constant.LocaleEN
	}

	class, err := c.query.New().FilterCode(code).WithLocale(locale).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.PlaceClass{}, errx.ErrorClassNotFound.Raise(
				fmt.Errorf("class with code %s not found, cause: %w", code, err),
			)
		default:
			return models.PlaceClass{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get class with code %s, cause: %w", code, err),
			)
		}
	}

	return classModelFromDB(class), nil
}

func (c Classificator) DeactivateClass(
	ctx context.Context,
	code, locale string,
) (models.PlaceClass, error) {
	class, err := c.GetClass(ctx, code, locale)
	if err != nil {
		return models.PlaceClass{}, err
	}

	if class.Status == constant.PlaceClassStatusesInactive {
		return class, nil
	}

	now := time.Now().UTC()
	status := constant.PlaceClassStatusesInactive
	err = c.query.New().FilterFatherCycle(code).Update(ctx, dbx.UpdatePlaceClassParams{
		Status:    &status,
		UpdatedAt: now,
	})
	if err != nil {
		return models.PlaceClass{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to deactivate class with code %s and its descendants, cause: %w", code, err),
		)
	}

	err = c.query.New().FilterCode(code).Update(ctx, dbx.UpdatePlaceClassParams{
		Status:    &status,
		UpdatedAt: now,
	})
	if err != nil {
		return models.PlaceClass{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to deactivate class with code %s, cause: %w", code, err),
		)
	}

	return class, nil
}

func (c Classificator) ActivateClass(
	ctx context.Context,
	code, locale string,
) (models.PlaceClass, error) {
	class, err := c.GetClass(ctx, code, locale)
	if err != nil {
		return models.PlaceClass{}, err
	}

	status := constant.PlaceClassStatusesActive
	now := time.Now().UTC()
	err = c.query.New().FilterCode(code).Update(ctx, dbx.UpdatePlaceClassParams{
		Status:    &status,
		UpdatedAt: now,
	})
	if err != nil {
		return models.PlaceClass{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to activate class with code %s, cause: %w", code, err),
		)
	}

	class.Status = status
	class.UpdatedAt = now
	return class, nil
}

type FilterClassesParams struct {
	Father      *string
	FatherCycle *bool
	Status      *string
	Locale      *string
}

func (c Classificator) ListClasses(
	ctx context.Context,
	filter FilterClassesParams,
	pag pagi.Response,
) ([]models.PlaceClass, pagi.Response, error) {
	if pag.Page == 0 {
		pag.Page = 1
	}
	if pag.Size == 0 {
		pag.Size = 20
	}
	if pag.Size > 100 {
		pag.Size = 100
	}

	limit := pag.Size + 1
	offset := (pag.Page - 1) * pag.Size

	query := c.query.New()

	if filter.Father != nil {
		if filter.FatherCycle != nil && *filter.FatherCycle {
			query = query.FilterFatherCycle(*filter.Father)
		}
		query = query.FilterFather(sql.NullString{
			String: *filter.Father,
			Valid:  true,
		})
	}
	if filter.Status != nil {
		query = query.FilterStatus(*filter.Status)
	}
	if filter.Locale != nil {
		query = query.WithLocale(*filter.Locale)
	}

	query = query.Page(limit, offset)

	rows, err := query.Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to select classes: %w", err),
		)
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}

	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
	}

	classes := make([]models.PlaceClass, 0, len(rows))
	for _, r := range rows {
		classes = append(classes, classModelFromDB(r))
	}

	return classes, pagi.Response{
		Page:  pag.Page,
		Size:  pag.Size,
		Total: count,
	}, nil
}

func classModelFromDB(dbClass dbx.PlaceClass) models.PlaceClass {
	resp := models.PlaceClass{
		Code:      dbClass.Code,
		Status:    dbClass.Status,
		Icon:      dbClass.Icon,
		Path:      dbClass.Path,
		CreatedAt: dbClass.CreatedAt,
		UpdatedAt: dbClass.UpdatedAt,
	}
	if dbClass.Father.Valid {
		resp.Father = &dbClass.Father.String
	}
	if dbClass.Name != nil {
		resp.Name = *dbClass.Name
	}
	if dbClass.Locale != nil {
		resp.Locale = *dbClass.Locale
	}

	return resp
}
