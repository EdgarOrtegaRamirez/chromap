// Package palette provides color palette generation and harmony.
package palette

import (
	"math"
	"sort"

	"github.com/EdgarOrtegaRamirez/chromap/internal/color"
)

// Palette represents a set of colors.
type Palette []color.Color

// HarmonyType defines the type of color harmony.
type HarmonyType string

const (
	Complementary   HarmonyType = "complementary"
	Analogous       HarmonyType = "analogous"
	Triadic         HarmonyType = "triadic"
	Tetradic        HarmonyType = "tetradic"
	SplitComplement HarmonyType = "split-complementary"
	Monochromatic   HarmonyType = "monochromatic"
)

// Generate creates a palette using the specified harmony.
func Generate(base color.Color, h HarmonyType, size int) Palette {
	hsl := base.ToHSL()
	p := make(Palette, 0, size)

	switch h {
	case Complementary:
		hue := math.Mod(hsl.H+180, 360)
		p = generateDual(hsl, hue, hsl.S, hsl.L, size)
	case Analogous:
		step := 30.0
		for i := -(size / 2); i <= size/2; i++ {
			if i == 0 {
				p = append(p, base)
			} else {
				hue := math.Mod(hsl.H+step*float64(i), 360)
				if hue < 0 {
					hue += 360
				}
				c := hslToColor(hue, hsl.S, hsl.L, 1)
				p = append(p, c)
			}
		}
	case Triadic:
		h1 := math.Mod(hsl.H+120, 360)
		h2 := math.Mod(hsl.H+240, 360)
		p = generateTriad(base, h1, h2, hsl.S, hsl.L, size)
	case Tetradic:
		h1 := math.Mod(hsl.H+90, 360)
		h2 := math.Mod(hsl.H+180, 360)
		h3 := math.Mod(hsl.H+270, 360)
		p = generateQuad(base, h1, h2, h3, hsl.S, hsl.L, size)
	case SplitComplement:
		h1 := math.Mod(hsl.H+150, 360)
		h2 := math.Mod(hsl.H+210, 360)
		p = generateTriad(base, h1, h2, hsl.S, hsl.L, size)
	case Monochromatic:
		step := 100.0 / float64(size+1)
		for i := 0; i < size; i++ {
			l := math.Max(5, math.Min(95, step*float64(i+1)))
			c := hslToColor(hsl.H, hsl.S, l, 1)
			p = append(p, c)
		}
	}

	// Trim to requested size
	if len(p) > size {
		p = p[:size]
	}
	return p
}

func generateDual(baseHSL color.HSL, hue2 float64, sat, light float64, size int) Palette {
	p := make(Palette, 0, size)
	hues := []float64{baseHSL.H, hue2}
	for i := 0; i < size; i++ {
		idx := i % 2
		var l float64
		switch idx {
		case 0:
			l = math.Max(5, math.Min(95, light+float64(i/2)*10-5))
		default:
			l = math.Max(5, math.Min(95, light+float64(i/2)*10-5))
		}
		c := hslToColor(hues[idx], sat, l, 1)
		p = append(p, c)
	}
	return p
}

func generateTriad(base color.Color, h1, h2 float64, sat, light float64, size int) Palette {
	p := make(Palette, 0, size)
	hues := []float64{base.ToHSL().H, h1, h2}
	for i := 0; i < size; i++ {
		idx := i % 3
		var l float64
		if i < 3 {
			l = light
		} else {
			l = math.Max(5, math.Min(95, light+float64(i/3)*15-5))
		}
		c := hslToColor(hues[idx], sat, l, 1)
		p = append(p, c)
	}
	return p
}

func generateQuad(base color.Color, h1, h2, h3 float64, sat, light float64, size int) Palette {
	p := make(Palette, 0, size)
	hues := []float64{base.ToHSL().H, h1, h2, h3}
	for i := 0; i < size; i++ {
		idx := i % 4
		var l float64
		if i < 4 {
			l = light
		} else {
			l = math.Max(5, math.Min(95, light+float64(i/4)*12-4))
		}
		c := hslToColor(hues[idx], sat, l, 1)
		p = append(p, c)
	}
	return p
}

func hslToColor(h, s, l, a float64) color.Color {
	r, g, b := hslToRGB(h, s/100.0, l/100.0)
	return color.Color{R: r, G: g, B: b, A: a}
}

func hslToRGB(h, s, l float64) (int, int, int) {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}

	var r, g, b float64
	if s < 0.001 {
		r, g, b = l, l, l
	} else {
		hue2rgb := func(p, q, t float64) float64 {
			if t < 0 {
				t += 1
			}
			if t > 1 {
				t -= 1
			}
			if t < 1.0/6.0 {
				return p + (q-p)*6*t
			}
			if t < 1.0/2.0 {
				return q
			}
			if t < 2.0/3.0 {
				return p + (q-p)*(2.0/3.0-t)*6
			}
			return p
		}

		q := 0.0
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q

		r = hue2rgb(p, q, h/360+1.0/3.0)
		g = hue2rgb(p, q, h/360)
		b = hue2rgb(p, q, h/360-1.0/3.0)
	}

	return int(math.Round(r * 255)), int(math.Round(g * 255)), int(math.Round(b * 255))
}

// BestTextColor determines whether black or white text works best on the background color.
func BestTextColor(bg color.Color) color.Color {
	white := color.Color{R: 255, G: 255, B: 255, A: 1}
	black := color.Color{R: 0, G: 0, B: 0, A: 1}
	if bg.Contrast(white) > bg.Contrast(black) {
		return white
	}
	return black
}

// SortByBrightness sorts a palette by brightness (dark to light).
func (p Palette) SortByBrightness() {
	sort.Slice(p, func(i, j int) bool {
		return p[i].Brightness() < p[j].Brightness()
	})
}

// SortByHue sorts a palette by hue.
func (p Palette) SortByHue() {
	sort.Slice(p, func(i, j int) bool {
		return p[i].ToHSL().H < p[j].ToHSL().H
	})
}
