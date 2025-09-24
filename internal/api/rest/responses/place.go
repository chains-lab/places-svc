package responses

import (
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/resources"
)

func Place(m models.PlaceWithDetails) resources.Place {
	resp := resources.Place{
		Data: resources.PlaceData{
			Id:   m.ID,
			Type: resources.PlaceType,
			Attributes: resources.PlaceDataAttributes{
				CityId:   m.CityID,
				Class:    m.Class,
				Status:   m.Status,
				Verified: m.Verified,
				Point: resources.Point{
					Lon: m.Point[0],
					Lat: m.Point[1],
				},
				Locale:      m.Locale,
				Name:        m.Name,
				Address:     m.Address,
				Description: m.Description,
				CreatedAt:   m.CreatedAt,
				UpdatedAt:   m.UpdatedAt,
			},
		},
	}

	if m.DistributorID != nil {
		resp.Data.Attributes.DistributorId = m.DistributorID
	}
	if m.Website != nil {
		resp.Data.Attributes.Website = m.Website
	}
	if m.Phone != nil {
		resp.Data.Attributes.Phone = m.Phone
	}

	if m.Timetable.Table != nil {
		resp.Included = make([]resources.TimetableData, 0, 1)
		resp.Included = append(resp.Included, Timetable(m.Timetable).Data)
	}

	return resp
}

func PlacesCollection(ms []models.PlaceWithDetails, pag pagi.Response) resources.PlacesCollection {
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
