# Mapbox Geocoding API v6 Migration Guide

This document outlines the changes and migration steps for upgrading from Mapbox Geocoding API v5 to v6.

## Overview

The Mapbox Geocoding API v6 brings significant improvements and changes:

- **Separate endpoints** for forward and reverse geocoding
- **Unified permanent/temporary** handling via parameters
- **New feature types** including `street`, `block`, and `secondary_address`
- **Structured Input** for more accurate geocoding
- **Smart Address Match** for better result validation
- **Improved Batch geocoding** supporting up to 1000 queries
- **Enhanced Japanese address support**
- **Removed POI support** (use Search Box API instead)

## Breaking Changes

### 1. API Endpoints

**v5:**
```
https://api.mapbox.com/geocoding/v5/mapbox.places/{query}.json
https://api.mapbox.com/geocoding/v5/mapbox.places-permanent/{query}.json
```

**v6:**
```
https://api.mapbox.com/search/geocode/v6/forward?q={query}
https://api.mapbox.com/search/geocode/v6/reverse?longitude={lon}&latitude={lat}
https://api.mapbox.com/search/geocode/v6/batch
```

### 2. Permanent vs Temporary Geocoding

**v5:** Different endpoints for permanent and temporary
**v6:** Single endpoint with `permanent` parameter

### 3. Feature Types

**Removed:**
- `poi` (use Search Box API instead)

**Added:**
- `street`
- `block` (for Japanese addresses)
- `secondary_address` (US only)

### 4. Response Structure Changes

- Removed `query` field from responses
- Added `match_code` object for address features
- Enhanced `coordinates` object with accuracy and routable points
- Updated context structure

## Migration Steps

### 1. Update Method Signatures

**v5 Forward Geocoding:**
```go
func (g *Geocode) Forward(place string, req *ForwardRequestOpts, permanent ...bool) (*ForwardResponse, error)
```

**v6 Forward Geocoding:**
```go
func (g *Geocode) Forward(place string, req *ForwardRequestOpts) (*ForwardResponse, error)
```

### 2. Handle Permanent Geocoding

**v5:**
```go
// Permanent geocoding
result, err := geocode.Forward("New York", &opts, true)
```

**v6:**
```go
opts := &ForwardRequestOpts{
    Permanent: true,
}
result, err := geocode.Forward("New York", opts)
```

### 3. Update Type Handling

**v5:**
```go
opts := &ForwardRequestOpts{
    Types: []Type{Address, Place},
}
```

**v6:**
```go
opts := &ForwardRequestOpts{
    Types: "address,place", // String format
}

// Or use helper method
result, err := geocode.ForwardWithTypes("New York", []Type{Address, Place}, opts)
```

### 4. Update Proximity Handling

**v5:**
```go
opts := &ForwardRequestOpts{
    Proximity: []float64{-74.0, 40.7},
}
```

**v6:**
```go
opts := &ForwardRequestOpts{
    Proximity: "-74.0,40.7", // String format
}

// Or use helper method
result, err := geocode.ForwardWithProximity("New York", []float64{-74.0, 40.7}, opts)
```

### 5. Update Reverse Geocoding

**v5:**
```go
opts := &ReverseRequestOpts{
    Types: []Type{Address},
}
```

**v6:**
```go
opts := &ReverseRequestOpts{
    Types: "address", // String format
}

// Or use helper method
result, err := geocode.ReverseWithTypes(location, []Type{Address}, opts)
```

## New Features

### 1. Structured Input (New in v6)

More accurate geocoding by specifying address components:

```go
opts := &StructuredInputOpts{
    AddressNumber: "1600",
    Street:        "Pennsylvania Avenue NW",
    Place:         "Washington",
    Region:        "DC",
    Country:       "US",
}
result, err := geocode.ForwardStructured(opts)
```

### 2. Smart Address Match

Access match quality information for address results:

```go
for _, feature := range result.Features {
    if feature.Properties.MatchCode != nil {
        confidence := feature.Properties.MatchCode.Confidence
        streetMatch := feature.Properties.MatchCode.Street
    }
}
```

### 3. Enhanced Batch Geocoding

Support for mixed query types and up to 1000 queries:

```go
queries := []BatchQuery{
    {Q: "New York", Types: "place"},
    {Longitude: -74.0, Latitude: 40.7, Types: "address"},
    {
        AddressNumber: "123",
        Street:        "Main St",
        Place:         "Anytown",
        Country:       "US",
    },
}

opts := &BatchRequestOpts{Permanent: true}
result, err := geocode.Batch(queries, opts)
```

### 4. New Optional Parameters

- `language`: Set response language
- `worldview`: Choose geographic worldview
- `format`: Response format (default: geojson)

```go
opts := &ForwardRequestOpts{
    Language:  "es",
    Worldview: "us",
    Format:    "geojson",
}
```

## Backward Compatibility

For easier migration, legacy methods are provided:

```go
// Deprecated but still available
result, err := geocode.ForwardLegacy("New York", opts, true)
result, err := geocode.ReverseLegacy(location, opts)
```

## Testing Your Migration

1. Update your import statements if needed
2. Replace permanent parameter usage with `Permanent` field
3. Convert slice parameters to string format or use helper methods
4. Update response handling to remove dependency on `query` field
5. Test with your existing queries to ensure results are still accurate

## Error Handling

Error handling remains largely the same, but v6 may return different error messages. Update your error handling to be more generic or check the new error format.

## Performance Considerations

- v6 may have different rate limits
- Batch geocoding is more efficient for multiple queries
- Structured input can provide more accurate results for well-formatted data

## Resources

- [Official v6 Documentation](https://docs.mapbox.com/api/search/geocoding/)
- [Migration from v5 Guide](https://docs.mapbox.com/api/search/geocoding/#migrating-from-geocoding-v5)
- [Search Box API](https://docs.mapbox.com/api/search/search-box/) (for POI searches)