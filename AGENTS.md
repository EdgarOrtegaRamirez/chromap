# chromap — AI Agent Guide

## Project Overview
chromap is a CLI tool for color palette generation, conversion, validation, and analysis.

## Build & Test
```bash
go build ./cmd/chromap
go test ./...
go vet ./...
```

## Running
```bash
go run ./cmd/chromap parse "#ff6b35"
go run ./cmd/chromap harmony "#4a90d9"
go run ./cmd/chromap contrast "#ffffff" "#333333"
```

## Structure
- `cmd/chromap/main.go` — CLI entry point
- `internal/color/color.go` — Color parsing, conversion, luminance, contrast
- `internal/palette/palette.go` — Palette harmony generation
- `internal/export/export.go` — Palette export formats
- `tests/` — Test suite

## Key Design Decisions
- No external dependencies (stdlib only)
- Colors use RGBA values (0-255 for RGB, 0-1 for alpha)
- HSL uses degrees for hue (0-360), percentages for saturation/lightness
- WCAG 2.1 relative luminance formula for contrast
- Hex output is always uppercase (#RRGGBB format)