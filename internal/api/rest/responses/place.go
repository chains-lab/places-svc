package responses

import (
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/resources"
)

func Place(m models.PlaceWithLocale) resources.Place {
	resp := resources.Place{
		Data: resources.PlaceData{
			Id:   m.Data.ID.String(),
			Type: resources.PlaceType,
			Attributes: resources.PlaceDataAttributes{
				CityId:   m.Data.CityID.String(),
				Class:    m.Data.Class,
				Status:   m.Data.Status,
				Verified: m.Data.Verified,
				Point: resources.Point{
					Lon: m.Data.Point[0],
					Lat: m.Data.Point[1],
				},
				Name:        m.Locale.Name,
				Address:     m.Data.Address,
				Description: m.Locale.Description,
				CreatedAt:   m.Data.CreatedAt,
				UpdatedAt:   m.Data.UpdatedAt,
			},
		},
	}

	if m.Data.DistributorID != nil {
		disID := m.Data.DistributorID.String()
		resp.Data.Attributes.DistributorId = &disID
	}
	if m.Data.Website != nil {
		resp.Data.Attributes.Website = m.Data.Website
	}
	if m.Data.Phone != nil {
		resp.Data.Attributes.Phone = m.Data.Phone
	}

	return resp
}

func PlacesCollection(ms []models.PlaceWithLocale, pag pagi.Response) resources.PlacesCollection {
	resp := resources.PlacesCollection{
		Data: make([]resources.PlaceData, 0, len(ms)),
		Links: resources.PaginationData{
			PageNumber: int64(pag.Page),
			PageSize:   int64(pag.Size),
			TotalItems: int64(pag.Total),
		},
	}

	for _, m := range ms {
		place := Place(m).Data

		resp.Data = append(resp.Data, place)
	}

	return resp
}

func PlaceLocale(m models.PlaceLocale) resources.PlaceLocale {
	return resources.PlaceLocale{
		Data: resources.PlaceLocaleData{
			Id:   m.PlaceID.String() + ":" + m.Locale,
			Type: resources.PlaceLocaleType,
			Attributes: resources.PlaceLocaleDataAttributes{
				Name:        m.Name,
				Description: m.Description,
			},
		},
	}
}

func PlaceLocalesCollection(ms []models.PlaceLocale, pag pagi.Response) resources.PlaceLocalesCollection {
	resp := resources.PlaceLocalesCollection{
		Data:     make([]resources.RelationshipDataObject, 0, len(ms)),
		Included: make([]resources.PlaceLocaleData, 0, len(ms)),
	}

	for _, m := range ms {
		locale := PlaceLocale(m).Data

		resp.Data = append(resp.Data, resources.RelationshipDataObject{
			Type: resources.PlaceLocaleType,
			Id:   locale.Id,
		})
		resp.Included = append(resp.Included, locale)
	}

	return resp
}
