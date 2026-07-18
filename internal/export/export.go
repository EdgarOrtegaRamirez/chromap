// Package export provides palette export in various formats.
package export

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/EdgarOrtegaRamirez/chromap/internal/color"
	"github.com/EdgarOrtegaRamirez/chromap/internal/palette"
)

// Format is the output format.
type Format string

const (
	FormatCSS      Format = "css"
	FormatJSON     Format = "json"
	FormatSVG      Format = "svg"
	FormatHTML     Format = "html"
	FormatTailwind Format = "tailwind"
	FormatText     Format = "text"
)

// PaletteData holds palette info for export.
type PaletteData struct {
	Name    string
	Base    color.Color
	Harmony palette.HarmonyType
	Colors  []color.Color
}

// Render converts a palette to the specified format.
func Render(data PaletteData, format Format) (string, error) {
	switch format {
	case FormatCSS:
		return renderCSS(data)
	case FormatJSON:
		return renderJSON(data)
	case FormatSVG:
		return renderSVG(data)
	case FormatHTML:
		return renderHTML(data)
	case FormatTailwind:
		return renderTailwind(data)
	case FormatText:
		return renderText(data), nil
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}

func renderCSS(data PaletteData) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("/* %s — %s palette */\n\n", data.Name, data.Harmony))

	if data.Name == "" {
		data.Name = "palette"
	}

	sb.WriteString(fmt.Sprintf(":root {\n"))
	for i, c := range data.Colors {
		name := fmt.Sprintf("--color-%s-%d", sanitize(data.Name), i+1)
		sb.WriteString(fmt.Sprintf("  %s: %s;\n", name, c.Hex()))
	}
	sb.WriteString("}\n")
	return sb.String(), nil
}

func renderJSON(data PaletteData) (string, error) {
	type colorEntry struct {
		Hex    string `json:"hex"`
		RGB    string `json:"rgb"`
		HSL    string `json:"hsl"`
		Bright string `json:"text_color"`
	}

	type paletteEntry struct {
		Name    string       `json:"name"`
		Base    string       `json:"base"`
		Harmony string       `json:"harmony"`
		Colors  []colorEntry `json:"colors"`
	}

	entries := make([]colorEntry, len(data.Colors))
	for i, c := range data.Colors {
		best := palette.BestTextColor(c)
		hsl := c.ToHSL()
		entries[i] = colorEntry{
			Hex:    c.Hex(),
			RGB:    c.RGB(),
			HSL:    fmt.Sprintf("hsl(%.0f, %.0f%%, %.0f%%)", hsl.H, hsl.S, hsl.L),
			Bright: best.Hex(),
		}
	}

	entry := paletteEntry{
		Name:    data.Name,
		Base:    data.Base.Hex(),
		Harmony: string(data.Harmony),
		Colors:  entries,
	}

	b, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func renderSVG(data PaletteData) (string, error) {
	var sb strings.Builder
	if data.Name == "" {
		data.Name = "palette"
	}

	w := 200
	h := 150
	padding := 10
	radius := 8

	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">\n`,
		w, len(data.Colors)*h, w, len(data.Colors)*h))

	for i, c := range data.Colors {
		y := i * h
		sb.WriteString(fmt.Sprintf(`  <rect x="%d" y="%d" width="%d" height="%d" rx="%d" fill="%s"/>\n`,
			padding, y+padding, w-padding*2, h-padding*2, radius, c.Hex()))

		// Text label
		hsl := c.ToHSL()
		sb.WriteString(fmt.Sprintf(`  <text x="%d" y="%d" font-family="monospace" font-size="11" fill="%s">`,
			w/2, y+40, palette.BestTextColor(c).Hex()))
		sb.WriteString(fmt.Sprintf(`<tspan x="%d" dy="0">#%s</tspan>`, w/2, c.Hex()[1:]))
		sb.WriteString(fmt.Sprintf(`<tspan x="%d" dy="14">HSL(%d,%d%%,%d%%)</tspan>`,
			w/2, int(hsl.H), int(hsl.S), int(hsl.L)))
		sb.WriteString(fmt.Sprintf(`<tspan x="%d" dy="14">RGB(%d,%d,%d)</tspan>`,
			w/2, c.R, c.G, c.B))
		sb.WriteString(`</text>`)
		sb.WriteString("\n")
	}

	sb.WriteString(`</svg>`)
	return sb.String(), nil
}

func renderHTML(data PaletteData) (string, error) {
	var sb strings.Builder
	if data.Name == "" {
		data.Name = "palette"
	}

	sb.WriteString(fmt.Sprintf(`<div class="chromap-palette" data-name="%s" data-harmony="%s">\n`,
		data.Name, data.Harmony))

	sb.WriteString(`<div class="colors" style="display:flex;gap:1rem;flex-wrap:wrap;">`)
	for _, c := range data.Colors {
		text := palette.BestTextColor(c)
		hsl := c.ToHSL()
		sb.WriteString(fmt.Sprintf(`<div class="swatch" style="background-color:%s;min-width:200px;padding:1rem;">`, c.Hex()))
		sb.WriteString(fmt.Sprintf(`<div style="color:%s;font-family:monospace;font-size:12px;">`, text.Hex()))
		sb.WriteString(fmt.Sprintf(`<div style="font-weight:bold;margin-bottom:0.5rem;">%s</div>`, c.Hex()))
		sb.WriteString(fmt.Sprintf(`<div>hsl(%.0f, %.0f%%, %.0f%%)</div>`, hsl.H, hsl.S, hsl.L))
		sb.WriteString(fmt.Sprintf(`<div>RGB(%d,%d,%d)</div>`, c.R, c.G, c.B))
		sb.WriteString(fmt.Sprintf(`<div>Brightness: %.0f</div>`, c.Brightness()))
		sb.WriteString(fmt.Sprintf(`<div>Light/ Dark: %s</div>`, map[bool]string{true: "Light", false: "Dark"}[c.IsLight()]))
		sb.WriteString(`</div></div>`)
	}
	sb.WriteString(`</div></div>`)
	return sb.String(), nil
}

func renderTailwind(data PaletteData) (string, error) {
	var sb strings.Builder
	if data.Name == "" {
		data.Name = "palette"
	}

	sb.WriteString(fmt.Sprintf("/* %s palette — Tailwind CSS config */\n\n", data.Name))
	sb.WriteString("module.exports = {\n")
	sb.WriteString(fmt.Sprintf("  extend: {\n    colors: {\n"))
	for i, c := range data.Colors {
		name := fmt.Sprintf("chromap-%s-%d", sanitize(data.Name), i+1)
		sb.WriteString(fmt.Sprintf("      \"%s\": \"%s\",\n", name, c.Hex()))
	}
	sb.WriteString("    }\n  }\n}\n")
	return sb.String(), nil
}

func renderText(data PaletteData) string {
	var sb strings.Builder
	if data.Name != "" {
		sb.WriteString(fmt.Sprintf("=== %s (%s) ===\n\n", data.Name, data.Harmony))
	}

	for i, c := range data.Colors {
		hsl := c.ToHSL()
		sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, c.Hex()))
		sb.WriteString(fmt.Sprintf("     RGB:  %s\n", c.RGB()))
		sb.WriteString(fmt.Sprintf("     HSL:  hsl(%.0f, %.0f%%, %.0f%%)\n", hsl.H, hsl.S, hsl.L))
		sb.WriteString(fmt.Sprintf("     HSV:  hsv(%.0f, %.0f%%, %.0f%%)\n", c.ToHSV().H, c.ToHSV().S, c.ToHSV().V))
		sb.WriteString(fmt.Sprintf("     Brightness: %.0f\n", c.Brightness()))
		sb.WriteString(fmt.Sprintf("     Type: %s\n", map[bool]string{true: "Light", false: "Dark"}[c.IsLight()]))

		if i < len(data.Colors)-1 {
			ratio := c.Contrast(data.Colors[i+1])
			sb.WriteString(fmt.Sprintf("     → Contrast with next: %.2f:1\n", ratio))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func sanitize(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	// Keep only alphanum and hyphens
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}
