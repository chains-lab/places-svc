package domaintest

import (
	"context"
	"testing"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/domain/services/class"
)

func TestClassLocales_SetGetAndList(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}
	cleanDb(t)

	ctx := context.Background()

	// Создаём класс сразу с EN-локалью (так устроена бизнес-логика создания)
	food := CreateClass(s, t, "Food", "food", nil)

	// 1) Фоллбэк: запрашиваем UK, но её ещё нет → должны получить EN
	gotUKBefore, err := s.domain.class.Get(ctx, food.Code, enum.LocaleUK)
	if err != nil {
		t.Fatalf("GetClass (fallback to EN): %v", err)
	}
	if gotUKBefore.Locale != enum.LocaleEN {
		t.Fatalf("fallback locale mismatch: want %s, got %s",
			enum.LocaleEN, gotUKBefore.Locale)
	}
	if gotUKBefore.Name != "Food" {
		t.Fatalf("fallback name mismatch: want %q, got %q", "Food", gotUKBefore.Name)
	}

	// 2) Добавляем локали UK и RU
	err = s.domain.class.SetLocale(ctx, food.Code,
		class.SetLocaleParams{Locale: enum.LocaleUK, Name: "Food UK"},
		class.SetLocaleParams{Locale: enum.LocaleRU, Name: "Еда RU"},
	)
	if err != nil {
		t.Fatalf("SetClassLocales: %v", err)
	}

	// Теперь UK должен отдаваться как UK
	gotUKAfter, err := s.domain.class.Get(ctx, food.Code, enum.LocaleUK)
	if err != nil {
		t.Fatalf("GetClass (after set UK): %v", err)
	}
	if gotUKAfter.Locale != enum.LocaleUK || gotUKAfter.Name != "Food UK" {
		t.Fatalf("GetClass UK mismatch: got locale=%s name=%q; want %s / %q",
			gotUKAfter.Locale, gotUKAfter.Name, enum.LocaleUK, "Food UK")
	}

	// Запросим несуществующую локаль (DE) → снова должен быть фоллбэк EN
	gotDE, err := s.domain.class.Get(ctx, food.Code, "de")
	if err != nil {
		t.Fatalf("GetClass (fallback from DE to EN): %v", err)
	}
	if gotDE.Locale != enum.LocaleEN {
		t.Fatalf("fallback from DE: want %s, got %s", enum.LocaleEN, gotDE.Locale)
	}
	if gotDE.Name != "Food" {
		t.Fatalf("fallback from DE: want name %q, got %q", "Food", gotDE.Name)
	}

	locs, err := s.domain.class.LocalesList(ctx, food.Code, 1, 10)
	if err != nil {
		t.Fatalf("ListClassLocales: %v", err)
	}
	// ожидаем минимум EN, UK, RU
	if locs.Total < 3 {
		t.Fatalf("ListClassLocales total: want >=3, got %d", locs.Total)
	}
	if len(locs.Data) < 3 {
		t.Fatalf("ListClassLocales len: want >=3, got %d", len(locs.Data))
	}

	// Проверим, что три ожидаемые локали присутствуют
	var haveEN, haveUK, haveRU bool
	for _, l := range locs.Data {
		switch l.Locale {
		case enum.LocaleEN:
			haveEN = true
		case enum.LocaleUK:
			haveUK = true
		case enum.LocaleRU:
			haveRU = true
		}
	}
	if !(haveEN && haveUK && haveRU) {
		t.Fatalf("missing locales in list: EN=%v UK=%v RU=%v", haveEN, haveUK, haveRU)
	}
}

func TestClassLocales_Pagination(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}
	cleanDb(t)

	ctx := context.Background()

	// База: EN создаётся вместе с классом
	c := CreateClass(s, t, "Clothes", "clothes", nil)

	// Добавим ещё 2 локали: всего будет 3 (en, uk, ru)
	if err = s.domain.class.SetLocale(ctx, c.Code,
		class.SetLocaleParams{Locale: enum.LocaleUK, Name: "Clothes UK"},
		class.SetLocaleParams{Locale: enum.LocaleRU, Name: "Одежда RU"},
	); err != nil {
		t.Fatalf("SetClassLocales: %v", err)
	}

	// Страница 1, размер 2
	page1, err := s.domain.class.LocalesList(ctx, c.Code, 1, 2)
	if err != nil {
		t.Fatalf("ListClassLocales page1: %v", err)
	}
	if page1.Total != 3 {
		t.Fatalf("total mismatch (page1): want 3, got %d", page1.Total)
	}
	if len(page1.Data) != 2 {
		t.Fatalf("page1 len mismatch: want 2, got %d", len(page1.Data))
	}

	// Страница 2, размер 2 (должна вернуть 1 запись)
	page2, err := s.domain.class.LocalesList(ctx, c.Code, 2, 2)
	if err != nil {
		t.Fatalf("ListClassLocales page2: %v", err)
	}
	if page2.Total != 3 {
		t.Fatalf("total mismatch (page2): want 3, got %d", page2.Total)
	}
	if len(page2.Data) != 1 {
		t.Fatalf("page2 len mismatch: want 1, got %d", len(page2.Data))
	}

	// Проверим, что среди полученных локалей есть en, uk, ru (независимо от порядка)
	got := map[string]bool{}
	for _, l := range append(page1.Data, page2.Data...) {
		got[l.Locale] = true
	}
	want := []string{enum.LocaleEN, enum.LocaleUK, enum.LocaleRU}
	for _, w := range want {
		if !got[w] {
			t.Fatalf("missing locale %q in paginated results; got=%v", w, got)
		}
	}
}
