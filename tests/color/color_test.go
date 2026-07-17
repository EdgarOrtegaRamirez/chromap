package color_test

import (
	"strings"
	"testing"

	"github.com/EdgarOrtegaRamirez/chromap/internal/color"
)

func TestParseHex(t *testing.T) {
	tests := []struct {
		input  string
		wantR  int
		wantG  int
		wantB  int
	}{
		{"#ff6b35", 255, 107, 53},
		{"#FF6B35", 255, 107, 53},
		{"#f63", 255, 102, 51},
		{"ff6b35", 255, 107, 53},
		{"black", 0, 0, 0},
		{"white", 255, 255, 255},
		{"rebeccapurple", 102, 51, 153},
		{"red", 255, 0, 0},
		{"blue", 0, 0, 255},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			c, err := color.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", tt.input, err)
			}
			if c.R != tt.wantR || c.G != tt.wantG || c.B != tt.wantB {
				t.Errorf("Parse(%q) = {R:%d G:%d B:%d}, want {R:%d G:%d B:%d}",
					tt.input, c.R, c.G, c.B, tt.wantR, tt.wantG, tt.wantB)
			}
		})
	}
}

func TestParseInvalid(t *testing.T) {
	_, err := color.Parse("notacolor")
	if err == nil {
		t.Error("Expected error for invalid color, got nil")
	}
}

func TestHexRoundTrip(t *testing.T) {
	tests := []string{"#ff0000", "#00ff00", "#0000ff", "#123456", "#0abc00", "#ff6b35"}
	for _, hex := range tests {
		c, err := color.Parse(hex)
		if err != nil {
			t.Fatalf("Parse(%q) error: %v", hex, err)
		}
		if got := c.Hex(); got != strings.ToUpper(hex) {
			t.Errorf("Round trip %s -> %s", hex, got)
		}
	}
}

func TestRGB(t *testing.T) {
	c := color.Color{R: 255, G: 128, B: 64, A: 1}
	expected := "rgb(255, 128, 64)"
	if got := c.RGB(); got != expected {
		t.Errorf("RGB() = %q, want %q", got, expected)
	}
}

func TestHSL(t *testing.T) {
	c := color.Color{R: 100, G: 150, B: 200, A: 1}
	hsl := c.ToHSL()
	if hsl.H < 0 || hsl.H > 360 {
		t.Errorf("HSL H out of range: %f", hsl.H)
	}
	if hsl.S < 0 || hsl.S > 100 {
		t.Errorf("HSL S out of range: %f", hsl.S)
	}
	if hsl.L < 0 || hsl.L > 100 {
		t.Errorf("HSL L out of range: %f", hsl.L)
	}
}

func TestContrastRatio(t *testing.T) {
	white := color.Color{R: 255, G: 255, B: 255, A: 1}
	black := color.Color{R: 0, G: 0, B: 0, A: 1}
	ratio := white.Contrast(black)
	if ratio < 20 || ratio > 22 {
		t.Errorf("White-black contrast ratio = %.2f, expected ~21", ratio)
	}
}

func TestIsLight(t *testing.T) {
	white := color.Color{R: 255, G: 255, B: 255, A: 1}
	if !white.IsLight() {
		t.Error("White should be light")
	}
	black := color.Color{R: 0, G: 0, B: 0, A: 1}
	if black.IsLight() {
		t.Error("Black should not be light")
	}
}

func TestIsDark(t *testing.T) {
	black := color.Color{R: 0, G: 0, B: 0, A: 1}
	if !black.IsDark() {
		t.Error("Black should be dark")
	}
	white := color.Color{R: 255, G: 255, B: 255, A: 1}
	if white.IsDark() {
		t.Error("White should not be dark")
	}
}

func TestBrightness(t *testing.T) {
	white := color.Color{R: 255, G: 255, B: 255, A: 1}
	if got := white.Brightness(); got != 100 {
		t.Errorf("White brightness = %.0f, expected 100", got)
	}
	black := color.Color{R: 0, G: 0, B: 0, A: 1}
	if got := black.Brightness(); got != 0 {
		t.Errorf("Black brightness = %.0f, expected 0", got)
	}
}

func TestRelativeLuminance(t *testing.T) {
	white := color.Color{R: 255, G: 255, B: 255, A: 1}
	if lum := white.RelativeLuminance(); lum < 0.99 || lum > 1.01 {
		t.Errorf("White luminance = %.4f, expected ~1.0", lum)
	}
	black := color.Color{R: 0, G: 0, B: 0, A: 1}
	if lum := black.RelativeLuminance(); lum > 0.01 {
		t.Errorf("Black luminance = %.4f, expected ~0", lum)
	}
}

func TestNamedColors(t *testing.T) {
	names := color.NamedColors()
	if len(names) == 0 {
		t.Fatal("Expected named colors, got empty list")
	}
	found := false
	for _, n := range names {
		if n == "rebeccapurple" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected rebeccapurple in named colors")
	}
}

func TestParseHSL(t *testing.T) {
	c, err := color.Parse("hsl(0, 100%, 50%)")
	if err != nil {
		t.Fatalf("Parse HSL error: %v", err)
	}
	if c.R != 255 || c.G != 0 || c.B != 0 {
		t.Errorf("hsl(0,100%%,50%%) = {R:%d G:%d B:%d}, want {R:255 G:0 B:0}", c.R, c.G, c.B)
	}
}

func TestParseRGB(t *testing.T) {
	c, err := color.Parse("rgb(100, 150, 200)")
	if err != nil {
		t.Fatalf("Parse RGB error: %v", err)
	}
	if c.R != 100 || c.G != 150 || c.B != 200 {
		t.Errorf("rgb(100,150,200) = {R:%d G:%d B:%d}, want {R:100 G:150 B:200}", c.R, c.G, c.B)
	}
}

func TestColorString(t *testing.T) {
	c := color.Color{R: 100, G: 150, B: 200, A: 1}
	if got := c.String(); got != "#6496C8" {
		t.Errorf("Color.String() = %q, want \"#6496C8\"", got)
	}
}