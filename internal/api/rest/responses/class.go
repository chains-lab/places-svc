package responses

import (
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/resources"
)

func Class(m models.ClassWithLocale) resources.Class {
	resp := resources.Class{
		Data: resources.ClassData{
			Id:   m.Data.Code,
			Type: resources.ClassType,
			Attributes: resources.ClassDataAttributes{
				Name:      m.Locale.Name,
				Status:    m.Data.Status,
				Icon:      m.Data.Icon,
				CreatedAt: m.Data.CreatedAt,
				UpdatedAt: m.Data.UpdatedAt,
			},
		},
	}

	if m.Data.Parent != nil {
		resp.Data.Relationships = &resources.ClassRelationships{
			Parent: resources.ClassRelationshipsParent{
				Data: &resources.RelationshipDataObject{
					Type: resources.ClassType,
					Id:   *m.Data.Parent,
				},
			},
		}
	}
	return resp
}

func ClassesCollection(ms []models.ClassWithLocale, pag pagi.Response) resources.ClassesCollection {
	resp := resources.ClassesCollection{
		Data: make([]resources.ClassData, 0, len(ms)),
		Links: resources.PaginationData{
			PageNumber: int64(pag.Page),
			PageSize:   int64(pag.Size),
			TotalItems: int64(pag.Total),
		},
	}

	for _, m := range ms {
		class := Class(m).Data

		resp.Data = append(resp.Data, class)
	}

	return resp
}

func ClassLocale(m models.ClassLocale) resources.ClassLocale {
	return resources.ClassLocale{
		Data: resources.ClassLocaleData{
			Id:   m.Class,
			Type: resources.ClassLocaleType,
			Attributes: resources.ClassLocaleDataAttributes{
				Name: m.Name,
			},
		},
	}
}

func ClassLocalesCollection(ms []models.ClassLocale, pag pagi.Response) resources.ClassLocalesCollection {
	resp := resources.ClassLocalesCollection{
		Data:     make([]resources.RelationshipDataObject, 0, len(ms)),
		Included: make([]resources.ClassLocaleData, 0, len(ms)),
	}

	for _, m := range ms {
		classLocale := ClassLocale(m).Data

		resp.Data = append(resp.Data, resources.RelationshipDataObject{
			Type: resources.ClassLocaleType,
			Id:   classLocale.Id,
		})
		resp.Included = append(resp.Included, classLocale)
	}

	return resp
}
