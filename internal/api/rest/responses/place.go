package responses

import (
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/resources"
)

func Place(m models.PlaceWithDetails) resources.Place {
	resp := resources.Place{
		Data: resources.PlaceData{
			Id:   m.Place.ID.String(),
			Type: resources.PlaceType,
			Attributes: resources.PlaceDataAttributes{
				CityId:   m.Place.CityID.String(),
				Class:    m.Place.Class,
				Status:   m.Place.Status,
				Verified: m.Place.Verified,
				Point: resources.Point{
					Lon: m.Place.Point[0],
					Lat: m.Place.Point[1],
				},
				Name:        m.Locale.Name,
				Address:     m.Place.Address,
				Description: m.Locale.Description,
				CreatedAt:   m.Place.CreatedAt,
				UpdatedAt:   m.Place.UpdatedAt,
			},
		},
	}

	if m.Place.DistributorID != nil {
		disID := m.Place.DistributorID.String()
		resp.Data.Attributes.DistributorId = &disID
	}
	if m.Place.Website != nil {
		resp.Data.Attributes.Website = m.Place.Website
	}
	if m.Place.Phone != nil {
		resp.Data.Attributes.Phone = m.Place.Phone
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
