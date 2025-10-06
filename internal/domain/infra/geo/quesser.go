package geo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/paulmach/orb"
)

type Address struct {
	Formatted   string `json:"formatted"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	Suburb      string `json:"suburb"`
	Street      string `json:"road"`
	HouseNumber string `json:"house_number"`
	Postcode    string `json:"postcode"`
}

type nominatimResp struct {
	DisplayName string            `json:"display_name"`
	Address     map[string]string `json:"address"`
}

type Guesser struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string
}

func NewGuesser() *Guesser {
	return &Guesser{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		baseURL:    "https://nominatim.openstreetmap.org/reverse",
		userAgent:  "chains-lab-places-svc",
	}
}

func (g *Guesser) Guess(ctx context.Context, pt orb.Point) (string, error) {
	q := url.Values{}
	q.Set("lat", fmt.Sprintf("%f", pt[1]))
	q.Set("lon", fmt.Sprintf("%f", pt[0]))
	q.Set("format", "json")
	q.Set("accept-language", "en")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.baseURL+"?"+q.Encode(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", g.userAgent)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("geocode failed: %s", resp.Status)
	}

	var raw nominatimResp
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return "", err
	}

	addr := Address{
		Formatted:   raw.DisplayName,
		Country:     raw.Address["country"],
		CountryCode: raw.Address["country_code"],
		State:       raw.Address["state"],
		City:        firstNonEmpty(raw.Address["city"], raw.Address["town"], raw.Address["village"]),
		Suburb:      raw.Address["suburb"],
		Street:      raw.Address["road"],
		HouseNumber: raw.Address["house_number"],
		Postcode:    raw.Address["postcode"],
	}

	return fmt.Sprintf("%+v\n", addr), nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
