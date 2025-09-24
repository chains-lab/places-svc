package responses

import (
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/resources"
)

func Class(m models.Class) resources.Class {
	resp := resources.Class{
		Data: resources.ClassData{
			Id:   m.Code,
			Type: resources.ClassType,
			Attributes: resources.ClassDataAttributes{
				Name:      m.Name,
				Status:    m.Status,
				Icon:      m.Icon,
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
			},
		},
	}

	if m.Parent != nil {
		resp.Data.Relationships = &resources.ClassRelationships{
			Parent: resources.ClassRelationshipsParent{
				Data: &resources.RelationshipDataObject{
					Type: resources.ClassType,
					Id:   *m.Parent,
				},
			},
		}
	}
	return resp
}

func ClassesCollection(ms models.ClassesCollection) resources.ClassesCollection {
	resp := resources.ClassesCollection{
		Data: make([]resources.ClassData, 0, len(ms.Data)),
		Links: resources.PaginationData{
			PageNumber: int64(ms.Page),
			PageSize:   int64(ms.Size),
			TotalItems: int64(ms.Total),
		},
	}

	for _, m := range ms.Data {
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

func ClassLocalesCollection(ms models.ClassLocaleCollection) resources.ClassLocalesCollection {
	resp := resources.ClassLocalesCollection{
		Data:     make([]resources.RelationshipDataObject, 0, len(ms.Data)),
		Included: make([]resources.ClassLocaleData, 0, len(ms.Data)),
		Links: resources.PaginationData{
			PageNumber: int64(ms.Page),
			PageSize:   int64(ms.Size),
			TotalItems: int64(ms.Total),
		},
	}

	for _, m := range ms.Data {
		classLocale := ClassLocale(m).Data

		resp.Data = append(resp.Data, resources.RelationshipDataObject{
			Type: resources.ClassLocaleType,
			Id:   classLocale.Id,
		})
		resp.Included = append(resp.Included, classLocale)
	}

	return resp
}
