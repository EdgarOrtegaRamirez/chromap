package palette_test

import (
	"testing"

	"github.com/EdgarOrtegaRamirez/chromap/internal/color"
	"github.com/EdgarOrtegaRamirez/chromap/internal/palette"
)

func TestGenerateComplementary(t *testing.T) {
	base := color.Color{R: 74, G: 144, B: 217, A: 1} // blue
	p := palette.Generate(base, palette.Complementary, 5)
	if len(p) != 5 {
		t.Fatalf("Expected 5 colors, got %d", len(p))
	}
}

func TestGenerateAnalogous(t *testing.T) {
	base := color.Color{R: 74, G: 144, B: 217, A: 1}
	p := palette.Generate(base, palette.Analogous, 7)
	if len(p) != 7 {
		t.Fatalf("Expected 7 colors, got %d", len(p))
	}
}

func TestGenerateMonochromatic(t *testing.T) {
	base := color.Color{R: 74, G: 144, B: 217, A: 1}
	p := palette.Generate(base, palette.Monochromatic, 5)
	if len(p) != 5 {
		t.Fatalf("Expected 5 colors, got %d", len(p))
	}
}

func TestGenerateTriadic(t *testing.T) {
	base := color.Color{R: 255, G: 107, B: 53, A: 1}
	p := palette.Generate(base, palette.Triadic, 6)
	if len(p) != 6 {
		t.Fatalf("Expected 6 colors, got %d", len(p))
	}
}

func TestGenerateTetradic(t *testing.T) {
	base := color.Color{R: 255, G: 107, B: 53, A: 1}
	p := palette.Generate(base, palette.Tetradic, 8)
	if len(p) != 8 {
		t.Fatalf("Expected 8 colors, got %d", len(p))
	}
}

func TestBestTextColor(t *testing.T) {
	dark := color.Color{R: 30, G: 30, B: 30, A: 1}
	text := palette.BestTextColor(dark)
	if text.R != 255 || text.G != 255 || text.B != 255 {
		t.Errorf("Best text on dark should be white, got %s", text.Hex())
	}

	light := color.Color{R: 255, G: 255, B: 255, A: 1}
	text = palette.BestTextColor(light)
	if text.R != 0 || text.G != 0 || text.B != 0 {
		t.Errorf("Best text on light should be black, got %s", text.Hex())
	}
}

func TestSortByBrightness(t *testing.T) {
	p := palette.Palette{
		color.Color{R: 255, G: 255, B: 255, A: 1},
		color.Color{R: 0, G: 0, B: 0, A: 1},
		color.Color{R: 128, G: 128, B: 128, A: 1},
	}
	p.SortByBrightness()
	if p[0].R != 0 || p[0].G != 0 || p[0].B != 0 {
		t.Error("First after sort should be black (darkest)")
	}
	if p[len(p)-1].R != 255 || p[len(p)-1].G != 255 || p[len(p)-1].B != 255 {
		t.Error("Last after sort should be white (lightest)")
	}
}

func TestPaletteTrimming(t *testing.T) {
	base := color.Color{R: 74, G: 144, B: 217, A: 1}
	// Request fewer than generated
	p := palette.Generate(base, palette.Complementary, 3)
	if len(p) > 3 {
		t.Errorf("Expected max 3 colors, got %d", len(p))
	}
}
