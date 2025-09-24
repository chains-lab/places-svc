package apptest

import (
	"context"
	"testing"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/domain"
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
	gotUKBefore, err := s.app.GetClass(ctx, food.Data.Code, enum.LocaleUK)
	if err != nil {
		t.Fatalf("GetClass (fallback to EN): %v", err)
	}
	if gotUKBefore.Locale.Locale != enum.LocaleEN {
		t.Fatalf("fallback locale mismatch: want %s, got %s",
			enum.LocaleEN, gotUKBefore.Locale.Locale)
	}
	if gotUKBefore.Locale.Name != "Food" {
		t.Fatalf("fallback name mismatch: want %q, got %q", "Food", gotUKBefore.Locale.Name)
	}

	// 2) Добавляем локали UK и RU
	err = s.app.SetClassLocales(ctx, food.Data.Code,
		app.SetClassLocaleParams{Locale: enum.LocaleUK, Name: "Food UK"},
		app.SetClassLocaleParams{Locale: enum.LocaleRU, Name: "Еда RU"},
	)
	if err != nil {
		t.Fatalf("SetClassLocales: %v", err)
	}

	// Теперь UK должен отдаваться как UK
	gotUKAfter, err := s.app.GetClass(ctx, food.Data.Code, enum.LocaleUK)
	if err != nil {
		t.Fatalf("GetClass (after set UK): %v", err)
	}
	if gotUKAfter.Locale.Locale != enum.LocaleUK || gotUKAfter.Locale.Name != "Food UK" {
		t.Fatalf("GetClass UK mismatch: got locale=%s name=%q; want %s / %q",
			gotUKAfter.Locale.Locale, gotUKAfter.Locale.Name, enum.LocaleUK, "Food UK")
	}

	// Запросим несуществующую локаль (DE) → снова должен быть фоллбэк EN
	gotDE, err := s.app.GetClass(ctx, food.Data.Code, "de")
	if err != nil {
		t.Fatalf("GetClass (fallback from DE to EN): %v", err)
	}
	if gotDE.Locale.Locale != enum.LocaleEN {
		t.Fatalf("fallback from DE: want %s, got %s", enum.LocaleEN, gotDE.Locale.Locale)
	}
	if gotDE.Locale.Name != "Food" {
		t.Fatalf("fallback from DE: want name %q, got %q", "Food", gotDE.Locale.Name)
	}

	// 3) Список локалей + total
	locs, pr, err := s.app.ListClassLocales(ctx, food.Data.Code, pagi.Request{Page: 1, Size: 10})
	if err != nil {
		t.Fatalf("ListClassLocales: %v", err)
	}
	// ожидаем минимум EN, UK, RU
	if pr.Total < 3 {
		t.Fatalf("ListClassLocales total: want >=3, got %d", pr.Total)
	}
	if len(locs) < 3 {
		t.Fatalf("ListClassLocales len: want >=3, got %d", len(locs))
	}

	// Проверим, что три ожидаемые локали присутствуют
	var haveEN, haveUK, haveRU bool
	for _, l := range locs {
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
	class := CreateClass(s, t, "Clothes", "clothes", nil)

	// Добавим ещё 2 локали: всего будет 3 (en, uk, ru)
	if err = s.app.SetClassLocales(ctx, class.Data.Code,
		app.SetClassLocaleParams{Locale: enum.LocaleUK, Name: "Clothes UK"},
		app.SetClassLocaleParams{Locale: enum.LocaleRU, Name: "Одежда RU"},
	); err != nil {
		t.Fatalf("SetClassLocales: %v", err)
	}

	// Страница 1, размер 2
	page1, pr1, err := s.app.ListClassLocales(ctx, class.Data.Code, pagi.Request{Page: 1, Size: 2})
	if err != nil {
		t.Fatalf("ListClassLocales page1: %v", err)
	}
	if pr1.Total != 3 {
		t.Fatalf("total mismatch (page1): want 3, got %d", pr1.Total)
	}
	if len(page1) != 2 {
		t.Fatalf("page1 len mismatch: want 2, got %d", len(page1))
	}

	// Страница 2, размер 2 (должна вернуть 1 запись)
	page2, pr2, err := s.app.ListClassLocales(ctx, class.Data.Code, pagi.Request{Page: 2, Size: 2})
	if err != nil {
		t.Fatalf("ListClassLocales page2: %v", err)
	}
	if pr2.Total != 3 {
		t.Fatalf("total mismatch (page2): want 3, got %d", pr2.Total)
	}
	if len(page2) != 1 {
		t.Fatalf("page2 len mismatch: want 1, got %d", len(page2))
	}

	// Проверим, что среди полученных локалей есть en, uk, ru (независимо от порядка)
	got := map[string]bool{}
	for _, l := range append(page1, page2...) {
		got[l.Locale] = true
	}
	want := []string{enum.LocaleEN, enum.LocaleUK, enum.LocaleRU}
	for _, w := range want {
		if !got[w] {
			t.Fatalf("missing locale %q in paginated results; got=%v", w, got)
		}
	}
}
