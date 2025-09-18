package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gnanakeethan/go-mapbox/lib/base"
	"github.com/gnanakeethan/go-mapbox/lib/geocode"
)

func main() {
	// Get your Mapbox access token from environment variable
	token := os.Getenv("MAPBOX_TOKEN")
	if token == "" {
		log.Fatal("MAPBOX_TOKEN environment variable is required")
	}

	// Create base client
	b, err := base.NewBase(token)
	if err != nil {
		log.Fatal("Failed to create base client:", err)
	}

	// Create geocoding client
	gc := geocode.NewGeocode(b)

	fmt.Println("=== Mapbox Geocoding API v6 Examples ===\n")

	// Example 1: Basic Forward Geocoding
	fmt.Println("1. Basic Forward Geocoding")
	basicForwardExample(gc)

	// Example 2: Forward Geocoding with Options
	fmt.Println("\n2. Forward Geocoding with Options")
	forwardWithOptionsExample(gc)

	// Example 3: Structured Input Forward Geocoding
	fmt.Println("\n3. Structured Input Forward Geocoding")
	structuredInputExample(gc)

	// Example 4: Reverse Geocoding
	fmt.Println("\n4. Reverse Geocoding")
	reverseGeocodingExample(gc)

	// Example 5: Permanent Geocoding
	fmt.Println("\n5. Permanent Geocoding")
	permanentGeocodingExample(gc)

	// Example 6: Batch Geocoding
	fmt.Println("\n6. Batch Geocoding")
	batchGeocodingExample(gc)

	// Example 7: Using Helper Methods
	fmt.Println("\n7. Using Helper Methods")
	helperMethodsExample(gc)
}

func basicForwardExample(gc *geocode.Geocode) {
	// Simple forward geocoding
	result, err := gc.Forward("2 Lincoln Memorial Circle NW", nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(result.Features) > 0 {
		feature := result.Features[0]
		coords := feature.Geometry.Coordinates
		fmt.Printf("Found: %s\n", feature.Properties["name"])
		fmt.Printf("Coordinates: [%f, %f]\n", coords[0], coords[1])
		fmt.Printf("Formatted Address: %s\n", feature.Properties["place_formatted"])
	}
}

func forwardWithOptionsExample(gc *geocode.Geocode) {
	// Forward geocoding with various options
	opts := &geocode.ForwardRequestOpts{
		Country:      "us",
		Limit:        5,
		Types:        "address",
		Autocomplete: false,
		Language:     "en",
		Worldview:    "us",
	}

	result, err := gc.Forward("1600 Pennsylvania Avenue", opts)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Found %d results:\n", len(result.Features))
	for i, feature := range result.Features {
		fmt.Printf("  %d. %s - %s\n", i+1,
			feature.Properties["name"],
			feature.Properties["place_formatted"])
	}
}

func structuredInputExample(gc *geocode.Geocode) {
	// Using structured input for more precise geocoding
	opts := &geocode.StructuredInputOpts{
		AddressNumber: "1600",
		Street:        "Pennsylvania Avenue NW",
		Place:         "Washington",
		Region:        "DC",
		Country:       "US",
		Autocomplete:  false,
	}

	result, err := gc.ForwardStructured(opts)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(result.Features) > 0 {
		feature := result.Features[0]
		coords := feature.Geometry.Coordinates
		fmt.Printf("Structured Input Result:\n")
		fmt.Printf("  Name: %s\n", feature.Properties["name"])
		fmt.Printf("  Coordinates: [%f, %f]\n", coords[0], coords[1])

		// Check for match code (Smart Address Match)
		if matchCode, ok := feature.Properties["match_code"].(map[string]interface{}); ok {
			if confidence, exists := matchCode["confidence"]; exists {
				fmt.Printf("  Match Confidence: %s\n", confidence)
			}
		}
	}
}

func reverseGeocodingExample(gc *geocode.Geocode) {
	// Reverse geocoding - coordinates to address
	location := &base.Location{
		Longitude: -77.036556,
		Latitude:  38.897708,
	}

	opts := &geocode.ReverseRequestOpts{
		Types:   "address",
		Limit:   1,
		Country: "us",
	}

	result, err := gc.Reverse(location, opts)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(result.Features) > 0 {
		feature := result.Features[0]
		fmt.Printf("Reverse geocoded address: %s\n", feature.Properties["name"])
		fmt.Printf("Full address: %s\n", feature.Properties["place_formatted"])
		fmt.Printf("Feature type: %s\n", feature.Properties["feature_type"])
	}
}

func permanentGeocodingExample(gc *geocode.Geocode) {
	// Using permanent geocoding (results can be stored)
	opts := &geocode.ForwardRequestOpts{
		Permanent: true,
		Limit:     1,
		Country:   "us",
	}

	result, err := gc.Forward("Empire State Building", opts)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(result.Features) > 0 {
		feature := result.Features[0]
		coords := feature.Geometry.Coordinates
		fmt.Printf("Permanent geocoding result:\n")
		fmt.Printf("  %s\n", feature.Properties["name"])
		fmt.Printf("  Coordinates: [%f, %f]\n", coords[0], coords[1])
		fmt.Printf("  Note: This result can be stored permanently\n")
	}
}

func batchGeocodingExample(gc *geocode.Geocode) {
	// Batch geocoding - multiple queries in one request
	queries := []geocode.BatchQuery{
		{
			Q:       "New York City",
			Types:   "place",
			Country: "us",
			Limit:   1,
		},
		{
			Longitude: -73.986136,
			Latitude:  40.748895,
			Types:     "address",
		},
		{
			AddressNumber: "123",
			Street:        "Main Street",
			Place:         "Boston",
			Region:        "MA",
			Country:       "us",
		},
	}

	opts := &geocode.BatchRequestOpts{
		Permanent: false,
	}

	result, err := gc.Batch(queries, opts)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Batch geocoding completed - %d queries processed:\n", len(result.Batch))
	for i, batch := range result.Batch {
		fmt.Printf("  Query %d: %d results\n", i+1, len(batch.Features))
		if len(batch.Features) > 0 {
			feature := batch.Features[0]
			fmt.Printf("    -> %s\n", feature.Properties["name"])
		}
	}
}

func helperMethodsExample(gc *geocode.Geocode) {
	// Using helper methods for easier API usage

	// 1. Forward geocoding with types using helper
	types := []geocode.Type{geocode.Address, geocode.Place}
	result1, err := gc.ForwardWithTypes("Central Park", types, nil)
	if err != nil {
		fmt.Printf("Error with types: %v\n", err)
	} else if len(result1.Features) > 0 {
		fmt.Printf("With types filter: %s\n", result1.Features[0].Properties["name"])
	}

	// 2. Forward geocoding with proximity using helper
	proximity := []float64{-73.968285, 40.785091} // Central Park coordinates
	result2, err := gc.ForwardWithProximity("Starbucks", proximity, &geocode.ForwardRequestOpts{Limit: 1})
	if err != nil {
		fmt.Printf("Error with proximity: %v\n", err)
	} else if len(result2.Features) > 0 {
		fmt.Printf("Near Central Park: %s\n", result2.Features[0].Properties["name"])
	}

	// 3. Reverse geocoding with types using helper
	location := &base.Location{Longitude: -73.968285, Latitude: 40.785091}
	reverseTypes := []geocode.Type{geocode.Address}
	result3, err := gc.ReverseWithTypes(location, reverseTypes, nil)
	if err != nil {
		fmt.Printf("Error with reverse types: %v\n", err)
	} else if len(result3.Features) > 0 {
		fmt.Printf("Reverse with types: %s\n", result3.Features[0].Properties["name"])
	}

	// 4. Legacy method (for backward compatibility)
	result4, err := gc.ForwardLegacy("Washington DC", nil, true) // permanent = true
	if err != nil {
		fmt.Printf("Error with legacy method: %v\n", err)
	} else if len(result4.Features) > 0 {
		fmt.Printf("Legacy method: %s\n", result4.Features[0].Properties["name"])
	}
}
