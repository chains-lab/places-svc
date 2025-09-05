package entities

//type Classificator struct {
//	kindQ       dbx.KindsQ
//	categoriesQ dbx.CategoryQ
//}
//
//func NewClassificator(db *sql.DB) Classificator {
//	return Classificator{
//		kindQ:       dbx.NewPlaceKindsQ(db),
//		categoriesQ: dbx.NewCategoryQ(db),
//	}
//}
//
//type PlaceCategoryCreateParams struct {
//	Code string
//	Name string
//}
//
//func (c Classificator) CreateNewCategory(ctx context.Context, params PlaceCategoryCreateParams) error {
//	return c.categoriesQ.New().Insert(ctx, dbx.PlaceCategory{
//		Code:      params.Code,
//		Name:      params.Name,
//		CreatedAt: time.Now().UTC(),
//		UpdatedAt: time.Now().UTC(),
//	})
//}
//
//type PlaceKindCreateParams struct {
//	Code       string
//	CategoryID string
//	Name       string
//}
//
//func (c Classificator) CreateNewKind(ctx context.Context, params PlaceKindCreateParams) error {
//	return c.kindQ.New().Insert(ctx, dbx.PlaceKind{
//		Code:         params.Code,
//		CategoryCode: params.CategoryID,
//		Name:         params.Name,
//		CreatedAt:    time.Now().UTC(),
//		UpdatedAt:    time.Now().UTC(),
//	})
//}
//
//func (c Classificator) GetCategory(ctx context.Context, categoryCode string) (models.Category, error) {
//	res, err := c.categoriesQ.New().FilterCode(categoryCode).Get(ctx)
//	if err != nil {
//		switch {
//		case errors.Is(err, sql.ErrNoRows):
//			return models.Category{}, errx.ErrorCategoryNotFound.Raise(fmt.Errorf("category %s not found", categoryCode))
//		default:
//			return models.Category{}, errx.ErrorInternal.Raise(fmt.Errorf("more than one row found for category %s", categoryCode))
//		}
//	}
//
//	return models.Category{
//		Code:      res.Code,
//		Name:      res.Name,
//		CreatedAt: res.CreatedAt,
//		UpdatedAt: res.UpdatedAt,
//	}, err
//}
//
//func (c Classificator) GetKind(ctx context.Context, kindCode string) (models.Kind, error) {
//	res, err := c.kindQ.New().FilterCode(kindCode).Get(ctx)
//	if err != nil {
//		switch {
//		case errors.Is(err, sql.ErrNoRows):
//			return models.Kind{}, errx.ErrorKindNotFound.Raise(fmt.Errorf("kind %s not found", kindCode))
//		default:
//			return models.Kind{}, errx.ErrorInternal.Raise(fmt.Errorf("more than one row found for kind %s", kindCode))
//		}
//	}
//
//	return models.Kind{
//		Code:       res.Code,
//		CategoryID: res.CategoryCode,
//		Name:       res.Name,
//		CreatedAt:  res.CreatedAt,
//		UpdatedAt:  res.UpdatedAt,
//	}, err
//}
//
//func (c Classificator) ListCategories(ctx context.Context) ([]models.Category, error) {
//	res, err := c.categoriesQ.New().Select(ctx)
//	if err != nil {
//		return nil, errx.ErrorInternal.Raise(fmt.Errorf("faidel to get categories, cause: %w", err))
//	}
//
//	var out []models.Category
//	for _, cat := range res {
//		out = append(out, toModelCategory(cat))
//	}
//
//	return out, nil
//}
//
//func (c Classificator) ListKinds(ctx context.Context) ([]models.Kind, error) {
//	res, err := c.kindQ.New().Select(ctx)
//	if err != nil {
//		return nil, errx.ErrorInternal.Raise(fmt.Errorf("faidel to get kinds, cause: %w", err))
//	}
//
//	var out []models.Kind
//	for _, k := range res {
//		out = append(out, toModelKind(k))
//	}
//
//	return out, nil
//}
//
//func (c Classificator) ListKindsForCategory(ctx context.Context, categoryCode string) ([]models.Kind, error) {
//	res, err := c.kindQ.New().FilterCategoryCode(categoryCode).Select(ctx)
//	if err != nil {
//		return nil, errx.ErrorInternal.Raise(
//			fmt.Errorf("faidel to get kinds for category %s, cause: %w", categoryCode, err),
//		)
//	}
//
//	var out []models.Kind
//	for _, k := range res {
//		out = append(out, toModelKind(k))
//	}
//
//	return out, nil
//}
//
//type UpdateCategoryParams struct {
//	Name string
//}
//
//func (c Classificator) UpdateCategory(ctx context.Context, categoryCOde string, params UpdateCategoryParams) error {
//	return c.categoriesQ.New().FilterCode(categoryCOde).Update(ctx, dbx.UpdatePlaceCategoryParams{
//		Name:      &params.Name,
//		UpdatedAt: time.Now().UTC(),
//	})
//}
//
//type UpdateKindParams struct {
//	Name string
//}
//
//func (c Classificator) UpdateKind(ctx context.Context, kindCode string, params UpdateKindParams) error {
//	return c.kindQ.New().FilterCode(kindCode).Update(ctx, dbx.PlaceUpdateParams{
//		Name:      &params.Name,
//		UpdatedAt: time.Now().UTC(),
//	})
//}
//
//func (c Classificator) ActivateCategory(ctx context.Context, categoryCode string) error {
//	status := constant.PlaceCategoryStatusActive
//
//	err := c.categoriesQ.New().FilterCode(categoryCode).Update(ctx, dbx.UpdatePlaceCategoryParams{
//		Status:    &status,
//		UpdatedAt: time.Now().UTC(),
//	})
//	if err != nil {
//		return fmt.Errorf("failed to activate category %s: %w", categoryCode, err)
//	}
//
//	return nil
//}
//
//func (c Classificator) DeactivateCategory(ctx context.Context, categoryCode string) error {
//	kindStatus := constant.PlaceCategoryStatusInactive
//	categoryStatus := constant.PlaceCategoryStatusInactive
//
//	err := c.categoriesQ.New().FilterCode(categoryCode).Update(ctx, dbx.UpdatePlaceCategoryParams{
//		Status:    &categoryStatus,
//		UpdatedAt: time.Now().UTC(),
//	})
//	if err != nil {
//		return fmt.Errorf("failed to activate category %s: %w", categoryCode, err)
//	}
//
//	err = c.kindQ.New().FilterCategoryCode(categoryCode).Update(ctx, dbx.PlaceUpdateParams{
//		Status:    &kindStatus,
//		UpdatedAt: time.Now().UTC(),
//	})
//	if err != nil {
//		return fmt.Errorf("failed to activate kinds for category %s: %w", categoryCode, err)
//	}
//
//	return nil
//}
//
//func (c Classificator) ActivateKind(ctx context.Context, kindCode string) error {
//	status := constant.PlaceKindStatusActive
//
//	err := c.kindQ.New().FilterCode(kindCode).Update(ctx, dbx.PlaceUpdateParams{
//		Status:    &status,
//		UpdatedAt: time.Now().UTC(),
//	})
//	if err != nil {
//		return fmt.Errorf("failed to activate kind %s: %w", kindCode, err)
//	}
//
//	return nil
//}
//
//func (c Classificator) DeactivateKind(ctx context.Context, kindCode string) error {
//	status := constant.PlaceKindStatusInactive
//
//	err := c.kindQ.New().FilterCode(kindCode).Update(ctx, dbx.PlaceUpdateParams{
//		Status:    &status,
//		UpdatedAt: time.Now().UTC(),
//	})
//	if err != nil {
//		return fmt.Errorf("failed to deactivate kind %s: %w", kindCode, err)
//	}
//
//	return nil
//}
//
//func toModelCategory(dbCategory dbx.PlaceCategory) models.Category {
//	return models.Category{
//		Code:      dbCategory.Code,
//		Name:      dbCategory.Name,
//		CreatedAt: dbCategory.CreatedAt,
//		UpdatedAt: dbCategory.UpdatedAt,
//	}
//}
//
//func toModelKind(dbKind dbx.PlaceKind) models.Kind {
//	return models.Kind{
//		Code:       dbKind.Code,
//		CategoryID: dbKind.CategoryCode,
//		Name:       dbKind.Name,
//		CreatedAt:  dbKind.CreatedAt,
//		UpdatedAt:  dbKind.UpdatedAt,
//	}
//}
