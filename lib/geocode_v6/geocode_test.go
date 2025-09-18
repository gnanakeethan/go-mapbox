/**
 * go-mapbox Geocoding Module Tests
 * Wraps the mapbox geocoding API for server side use
 * See https://docs.mapbox.com/api/search/geocoding/ for API information
 *
 * https://github.com/ryankurte/go-mapbox
 * https://github.com/gnanakeethan/go-mapbox
 * Copyright 2017-2025 Ryan Kurte
 * Copyright 2025 Gnanakeethan Balasubramaniam
 */

package geocode

import (
	"os"
	"testing"

	"github.com/gnanakeethan/go-mapbox/lib/base"
)

func TestGeocoder(t *testing.T) {
	b, err := base.NewBase(os.Getenv("MAPBOX_TOKEN"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	geocode := NewGeocode(b)

	t.Run("Can geocode", func(t *testing.T) {
		var reqOpt ForwardRequestOpts
		reqOpt.Limit = 1

		place := "2 lincoln memorial circle nw"

		res, err := geocode.Forward(place, &reqOpt)
		if err != nil {
			t.Error(err)
		}

		if res.Type != "FeatureCollection" {
			t.Errorf("Invalid response type: %s", res.Type)
		}

		if len(res.Features) == 0 {
			t.Error("No features returned")
		}
	})

	t.Run("Can reverse geocode", func(t *testing.T) {
		var reqOpt ReverseRequestOpts
		reqOpt.Limit = 1

		loc := &base.Location{Longitude: 72.438939, Latitude: 34.074122}

		res, err := geocode.Reverse(loc, &reqOpt)
		if err != nil {
			t.Error(err)
		}

		if res.Type != "FeatureCollection" {
			t.Errorf("Invalid response type: %s", res.Type)
		}

		if len(res.Features) == 0 {
			t.Error("No features returned")
		}
	})

	t.Run("Can use structured input", func(t *testing.T) {
		var reqOpt StructuredInputOpts
		reqOpt.AddressNumber = "1600"
		reqOpt.Street = "Pennsylvania Avenue NW"
		reqOpt.Place = "Washington"
		reqOpt.Region = "DC"
		reqOpt.Country = "US"
		reqOpt.Limit = 1

		res, err := geocode.ForwardStructured(&reqOpt)
		if err != nil {
			t.Error(err)
		}

		if res.Type != "FeatureCollection" {
			t.Errorf("Invalid response type: %s", res.Type)
		}

		if len(res.Features) == 0 {
			t.Error("No features returned")
		}
	})

	t.Run("Can use permanent geocoding", func(t *testing.T) {
		var reqOpt ForwardRequestOpts
		reqOpt.Limit = 1
		reqOpt.Permanent = true

		place := "2 lincoln memorial circle nw"

		res, err := geocode.Forward(place, &reqOpt)
		if err != nil {
			t.Error(err)
		}

		if res.Type != "FeatureCollection" {
			t.Errorf("Invalid response type: %s", res.Type)
		}

		if len(res.Features) == 0 {
			t.Error("No features returned")
		}
	})
}
