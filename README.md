# chromap

A CLI tool for color palette generation, conversion, validation, and analysis. Supports hex, RGB, HSL, HSV, and 148 named colors. Generates color harmonies and checks WCAG contrast ratios.

## Features

- **Color Parsing** — Parse hex (#RGB, #RRGGBB, #RRGGBBAA), RGB, HSL, and 148 named colors
- **Conversion** — Convert between all color formats instantly
- **Color Harmonies** — Generate complementary, analogous, triadic, tetradic, split-complementary, and monochromatic palettes
- **Contrast Checking** — WCAG 2.1 compliant contrast ratio calculations with AA/AAA pass/fail
- **Palette Export** — Export palettes as CSS variables, JSON, SVG, HTML, Tailwind CSS, or plain text

## Install

```bash
# Clone and build
git clone https://github.com/EdgarOrtegaRamirez/chromap.git
cd chromap
go build -o chromap ./cmd/chromap

# Install globally
go install github.com/EdgarOrtegaRamirez/chromap/cmd/chromap@latest
```

## Usage

### Parse a color

```bash
chromap parse "#ff6b35"
chromap parse "rebeccapurple"
chromap parse "rgb(255, 107, 53)"
chromap parse "hsl(19, 100%, 61%)"
```

Output:
```
Hex:     #FF6B35
RGB:     rgb(255, 107, 53)
HSL:     hsl(19, 100%, 61%)
HSV:     hsv(19, 79%, 100%)
Brightness: 71
Light/Dark: Light
Luminance: 0.2500
```

### Convert a color

```bash
chromap convert "#4a90d9"
```

### Generate color harmonies

```bash
# Default: complementary palette (5 colors)
chromap harmony "#4a90d9"

# With harmony type and size
chromap harmony "#4a90d9" analogous 7
chromap harmony "#ff6b35" triadic 6
chromap harmony "#4a90d9" monochromatic 5
```

Available harmonies: `complementary`, `analogous`, `triadic`, `tetradic`, `split-complementary`, `monochromatic`

### Check contrast

```bash
chromap contrast "#ffffff" "#333333"
chromap contrast "#ff6b35" "#000000"
```

Output:
```
Contrast ratio: 15.23:1
WCAG AA (normal text): Pass
WCAG AA (large text):  Pass
WCAG AAA (normal text): Pass
WCAG AAA (large text):  Pass
```

### List named colors

```bash
chromap colors
```

Lists all 148 available named CSS colors with their hex values.

## API Usage (Go)

```go
import "github.com/EdgarOrtegaRamirez/chromap/internal/color"
import "github.com/EdgarOrtegaRamirez/chromap/internal/palette"

// Parse a color
c, err := color.Parse("#ff6b35")

// Convert formats
fmt.Println(c.Hex())    // "#FF6B35"
fmt.Println(c.RGB())    // "rgb(255, 107, 53)"
hsl := c.ToHSL()        // HSL{H: 19, S: 100, L: 61}

// Check contrast
contrast := c.Contrast(otherColor)

// Generate a palette
p := palette.Generate(c, palette.Complementary, 5)

// Best text color on background
textColor := palette.BestTextColor(bgColor)
```

## License

MIT