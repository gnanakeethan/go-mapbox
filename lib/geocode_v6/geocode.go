/**
 * go-mapbox Geocoding Module
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
	"fmt"
	"strconv"
	"strings"

	"github.com/gnanakeethan/go-mapbox/lib/base"
	"github.com/google/go-querystring/query"
)

const (
	apiName    = "search"
	apiVersion = "geocode/v6"
)

// Type defines geocode location response types
type Type string

const (
	// Country level
	Country Type = "country"
	// Region level
	Region Type = "region"
	// Postcode level
	Postcode Type = "postcode"
	// District level
	District Type = "district"
	// Place level
	Place Type = "place"
	// Locality level
	Locality Type = "locality"
	// Neighborhood level
	Neighborhood Type = "neighborhood"
	// Street level (new in v6)
	Street Type = "street"
	// Block level (new in v6, for Japanese addresses)
	Block Type = "block"
	// Address level
	Address Type = "address"
	// Secondary address level (new in v6, US only)
	SecondaryAddress Type = "secondary_address"
)

// Geocode api wrapper instance
type Geocode struct {
	base *base.Base
}

// NewGeocode Create a new Geocode API wrapper
func NewGeocode(base *base.Base) *Geocode {
	return &Geocode{base}
}

// ForwardRequestOpts request options for forward geocoding
type ForwardRequestOpts struct {
	Country      string `url:"country,omitempty"`
	Proximity    string `url:"proximity,omitempty"`
	Types        string `url:"types,omitempty"`
	Autocomplete bool   `url:"autocomplete,omitempty"`
	BBox         string `url:"bbox,omitempty"`
	Limit        uint   `url:"limit,omitempty"`
	Language     string `url:"language,omitempty"`
	Worldview    string `url:"worldview,omitempty"`
	Format       string `url:"format,omitempty"`
	Permanent    bool   `url:"permanent,omitempty"`
}

// StructuredInputOpts options for structured input forward geocoding (new in v6)
type StructuredInputOpts struct {
	AddressLine1  string `url:"address_line1,omitempty"`
	AddressNumber string `url:"address_number,omitempty"`
	Street        string `url:"street,omitempty"`
	Block         string `url:"block,omitempty"`
	Place         string `url:"place,omitempty"`
	Region        string `url:"region,omitempty"`
	Postcode      string `url:"postcode,omitempty"`
	Locality      string `url:"locality,omitempty"`
	Neighborhood  string `url:"neighborhood,omitempty"`
	Country       string `url:"country,omitempty"`
	Autocomplete  bool   `url:"autocomplete,omitempty"`
	BBox          string `url:"bbox,omitempty"`
	Format        string `url:"format,omitempty"`
	Language      string `url:"language,omitempty"`
	Limit         uint   `url:"limit,omitempty"`
	Proximity     string `url:"proximity,omitempty"`
	Types         string `url:"types,omitempty"`
	Worldview     string `url:"worldview,omitempty"`
	Permanent     bool   `url:"permanent,omitempty"`
}

// MatchCode represents the smart address match information (new in v6)
type MatchCode struct {
	AddressNumber string `json:"address_number,omitempty"`
	Street        string `json:"street,omitempty"`
	Postcode      string `json:"postcode,omitempty"`
	Place         string `json:"place,omitempty"`
	Region        string `json:"region,omitempty"`
	Locality      string `json:"locality,omitempty"`
	Country       string `json:"country,omitempty"`
	Confidence    string `json:"confidence,omitempty"`
}

// Coordinates represents the enhanced coordinate information (updated in v6)
type Coordinates struct {
	Longitude      float64 `json:"longitude"`
	Latitude       float64 `json:"latitude"`
	Accuracy       string  `json:"accuracy,omitempty"`
	RoutablePoints []struct {
		Name      string  `json:"name"`
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
	} `json:"routable_points,omitempty"`
}

// ForwardResponse is the response from a forward geocode lookup
type ForwardResponse struct {
	*base.FeatureCollection
}

// ReverseResponse is the response to a reverse geocode request
type ReverseResponse struct {
	*base.FeatureCollection
}

// BatchResponse is the response to a batch geocode request (new in v6)
type BatchResponse struct {
	Batch []base.FeatureCollection `json:"batch"`
}

// Forward geocode lookup using search text
// Finds locations from a place name
func (g *Geocode) Forward(place string, req *ForwardRequestOpts) (*ForwardResponse, error) {
	if req == nil {
		req = &ForwardRequestOpts{}
	}

	v, err := query.Values(req)
	if err != nil {
		return nil, err
	}

	// Add the search query parameter
	v.Set("q", place)

	resp := ForwardResponse{}
	err = g.base.QueryBase(fmt.Sprintf("%s/%s/forward", apiName, apiVersion), &v, &resp)

	return &resp, err
}

// ForwardStructured geocode lookup using structured input (new in v6)
// Provides more accurate results by specifying address components separately
func (g *Geocode) ForwardStructured(req *StructuredInputOpts) (*ForwardResponse, error) {
	v, err := query.Values(req)
	if err != nil {
		return nil, err
	}

	resp := ForwardResponse{}
	err = g.base.QueryBase(fmt.Sprintf("%s/%s/forward", apiName, apiVersion), &v, &resp)

	return &resp, err
}

// ReverseRequestOpts request options for reverse geocoding
type ReverseRequestOpts struct {
	Types     string `url:"types,omitempty"`
	Limit     uint   `url:"limit,omitempty"`
	Country   string `url:"country,omitempty"`
	Language  string `url:"language,omitempty"`
	Worldview string `url:"worldview,omitempty"`
	Permanent bool   `url:"permanent,omitempty"`
}

// Reverse geocode lookup
// Finds place names from a location
func (g *Geocode) Reverse(loc *base.Location, req *ReverseRequestOpts) (*ReverseResponse, error) {
	if req == nil {
		req = &ReverseRequestOpts{}
	}

	v, err := query.Values(req)
	if err != nil {
		return nil, err
	}

	// Add longitude and latitude parameters
	v.Set("longitude", strconv.FormatFloat(loc.Longitude, 'f', -1, 64))
	v.Set("latitude", strconv.FormatFloat(loc.Latitude, 'f', -1, 64))

	resp := ReverseResponse{}
	err = g.base.QueryBase(fmt.Sprintf("%s/%s/reverse", apiName, apiVersion), &v, &resp)

	return &resp, err
}

// BatchQuery represents a single query in a batch request
type BatchQuery struct {
	// For forward geocoding
	Q        string `json:"q,omitempty"`
	Types    string `json:"types,omitempty"`
	BBox     string `json:"bbox,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Country  string `json:"country,omitempty"`
	Language string `json:"language,omitempty"`

	// For structured input forward geocoding
	AddressLine1  string `json:"address_line1,omitempty"`
	AddressNumber string `json:"address_number,omitempty"`
	Street        string `json:"street,omitempty"`
	Block         string `json:"block,omitempty"`
	Place         string `json:"place,omitempty"`
	Region        string `json:"region,omitempty"`
	Postcode      string `json:"postcode,omitempty"`
	Locality      string `json:"locality,omitempty"`
	Neighborhood  string `json:"neighborhood,omitempty"`

	// For reverse geocoding
	Longitude float64 `json:"longitude,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`

	// Common parameters
	Autocomplete bool   `json:"autocomplete,omitempty"`
	Proximity    string `json:"proximity,omitempty"`
	Worldview    string `json:"worldview,omitempty"`
	Format       string `json:"format,omitempty"`
}

// BatchRequestOpts options for batch geocoding requests
type BatchRequestOpts struct {
	Permanent bool `url:"permanent,omitempty"`
}

// Batch geocode lookup (improved in v6)
// Allows up to 1000 forward or reverse geocoding queries in a single request
func (g *Geocode) Batch(queries []BatchQuery, req *BatchRequestOpts) (*BatchResponse, error) {
	if req == nil {
		req = &BatchRequestOpts{}
	}

	// Build query parameters
	queryValues, err := query.Values(req)
	if err != nil {
		return nil, err
	}

	resp := BatchResponse{}
	err = g.base.QueryWithBodyBase(fmt.Sprintf("%s/%s/batch", apiName, apiVersion), &queryValues, queries, &resp)

	return &resp, err
}

// Helper function to convert Type slice to comma-separated string
func typesToString(types []Type) string {
	if len(types) == 0 {
		return ""
	}

	stringTypes := make([]string, len(types))
	for i, t := range types {
		stringTypes[i] = string(t)
	}
	return strings.Join(stringTypes, ",")
}

// Helper function to format proximity coordinates
func proximityToString(proximity []float64) string {
	if len(proximity) != 2 {
		return ""
	}
	return fmt.Sprintf("%f,%f", proximity[0], proximity[1])
}

// Helper function to format bounding box
func bboxToString(bbox base.BoundingBox) string {
	if len(bbox) != 4 {
		return ""
	}
	return fmt.Sprintf("%f,%f,%f,%f", bbox[0], bbox[1], bbox[2], bbox[3])
}

// ForwardWithTypes is a convenience method that accepts Type slice and converts to string
func (g *Geocode) ForwardWithTypes(place string, types []Type, req *ForwardRequestOpts) (*ForwardResponse, error) {
	if req == nil {
		req = &ForwardRequestOpts{}
	}
	if len(types) > 0 {
		req.Types = typesToString(types)
	}
	return g.Forward(place, req)
}

// ForwardWithProximity is a convenience method that accepts proximity coordinates as float slice
func (g *Geocode) ForwardWithProximity(place string, proximity []float64, req *ForwardRequestOpts) (*ForwardResponse, error) {
	if req == nil {
		req = &ForwardRequestOpts{}
	}
	if len(proximity) == 2 {
		req.Proximity = proximityToString(proximity)
	}
	return g.Forward(place, req)
}

// ForwardWithBBox is a convenience method that accepts bounding box and converts to string
func (g *Geocode) ForwardWithBBox(place string, bbox base.BoundingBox, req *ForwardRequestOpts) (*ForwardResponse, error) {
	if req == nil {
		req = &ForwardRequestOpts{}
	}
	if len(bbox) == 4 {
		req.BBox = bboxToString(bbox)
	}
	return g.Forward(place, req)
}

// ReverseWithTypes is a convenience method that accepts Type slice and converts to string
func (g *Geocode) ReverseWithTypes(loc *base.Location, types []Type, req *ReverseRequestOpts) (*ReverseResponse, error) {
	if req == nil {
		req = &ReverseRequestOpts{}
	}
	if len(types) > 0 {
		req.Types = typesToString(types)
	}
	return g.Reverse(loc, req)
}

// ForwardLegacy provides backward compatibility with v5-style API calls
// Deprecated: Use Forward, ForwardWithTypes, or ForwardStructured instead
func (g *Geocode) ForwardLegacy(place string, req *ForwardRequestOpts, permanent ...bool) (*ForwardResponse, error) {
	if req == nil {
		req = &ForwardRequestOpts{}
	}

	// Handle permanent parameter for backward compatibility
	if len(permanent) > 0 && permanent[0] {
		req.Permanent = true
	}

	return g.Forward(place, req)
}

// ReverseLegacy provides backward compatibility with v5-style reverse geocoding
// Deprecated: Use Reverse or ReverseWithTypes instead
func (g *Geocode) ReverseLegacy(loc *base.Location, req *ReverseRequestOpts) (*ReverseResponse, error) {
	return g.Reverse(loc, req)
}
