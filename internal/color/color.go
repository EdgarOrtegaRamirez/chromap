// Package color provides color parsing, conversion, and manipulation.
package color

import (
	"fmt"
	"html"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Color represents a color in multiple formats.
type Color struct {
	R, G, B int     // 0-255
	A       float64 // 0-1 (alpha)
}

// HSL represents a color in HSL format.
type HSL struct {
	H, S, L float64 // H: 0-360, S: 0-100, L: 0-100
}

// HSV represents a color in HSV format.
type HSV struct {
	H, S, V float64 // H: 0-360, S: 0-100, V: 0-100
}

// HTML named colors map.
var namedColors = map[string]Color{
	"aliceblue":            {240, 248, 255, 1},
	"antiquewhite":         {250, 235, 215, 1},
	"aqua":                 {0, 255, 255, 1},
	"aquamarine":           {127, 255, 212, 1},
	"azure":                {240, 255, 255, 1},
	"beige":                {245, 245, 220, 1},
	"bisque":               {255, 228, 196, 1},
	"black":                {0, 0, 0, 1},
	"blanchedalmond":       {255, 235, 205, 1},
	"blue":                 {0, 0, 255, 1},
	"blueviolet":           {138, 43, 226, 1},
	"brown":                {165, 42, 42, 1},
	"burlywood":            {222, 184, 135, 1},
	"cadetblue":            {95, 158, 160, 1},
	"chartreuse":           {127, 255, 0, 1},
	"chocolate":            {210, 105, 30, 1},
	"coral":                {255, 127, 80, 1},
	"cornflowerblue":       {100, 149, 237, 1},
	"cornsilk":             {255, 248, 220, 1},
	"crimson":              {220, 20, 60, 1},
	"cyan":                 {0, 255, 255, 1},
	"darkblue":             {0, 0, 139, 1},
	"darkcyan":             {0, 139, 139, 1},
	"darkgoldenrod":        {184, 134, 11, 1},
	"darkgray":             {169, 169, 169, 1},
	"darkgreen":            {0, 100, 0, 1},
	"darkgrey":             {169, 169, 169, 1},
	"darkkhaki":            {189, 183, 107, 1},
	"darkmagenta":          {139, 0, 139, 1},
	"darkolivegreen":       {85, 107, 47, 1},
	"darkorange":           {255, 140, 0, 1},
	"darkorchid":           {153, 50, 204, 1},
	"darkred":              {139, 0, 0, 1},
	"darksalmon":           {233, 150, 122, 1},
	"darkseagreen":         {143, 188, 143, 1},
	"darkslateblue":        {72, 61, 139, 1},
	"darkslategray":        {47, 79, 79, 1},
	"darkslategrey":        {47, 79, 79, 1},
	"darkturquoise":        {0, 206, 209, 1},
	"darkviolet":           {148, 0, 211, 1},
	"deeppink":             {255, 20, 147, 1},
	"deepskyblue":          {0, 191, 255, 1},
	"dimgray":              {105, 105, 105, 1},
	"dimgrey":              {105, 105, 105, 1},
	"dodgerblue":           {30, 144, 255, 1},
	"firebrick":            {178, 34, 34, 1},
	"floralwhite":          {255, 250, 240, 1},
	"forestgreen":          {34, 139, 34, 1},
	"fuchsia":              {255, 0, 255, 1},
	"gainsboro":            {220, 220, 220, 1},
	"ghostwhite":           {248, 248, 255, 1},
	"gold":                 {255, 215, 0, 1},
	"goldenrod":            {218, 165, 32, 1},
	"gray":                 {128, 128, 128, 1},
	"green":                {0, 128, 0, 1},
	"greenyellow":          {173, 255, 47, 1},
	"grey":                 {128, 128, 128, 1},
	"honeydew":             {240, 255, 240, 1},
	"hotpink":              {255, 105, 180, 1},
	"indianred":            {205, 92, 92, 1},
	"indigo":               {75, 0, 130, 1},
	"ivory":                {255, 255, 240, 1},
	"khaki":                {240, 230, 140, 1},
	"lavender":             {230, 230, 250, 1},
	"lavenderblush":        {255, 240, 245, 1},
	"lawngreen":            {124, 252, 0, 1},
	"lemonchiffon":         {255, 250, 205, 1},
	"lightblue":            {173, 216, 230, 1},
	"lightcoral":           {240, 128, 128, 1},
	"lightcyan":            {224, 255, 255, 1},
	"lightgoldenrodyellow": {250, 250, 210, 1},
	"lightgray":            {211, 211, 211, 1},
	"lightgreen":           {144, 238, 144, 1},
	"lightgrey":            {211, 211, 211, 1},
	"lightpink":            {255, 182, 193, 1},
	"lightsalmon":          {255, 160, 122, 1},
	"lightseagreen":        {32, 178, 170, 1},
	"lightskyblue":         {135, 206, 250, 1},
	"lightslategray":       {119, 136, 153, 1},
	"lightslategrey":       {119, 136, 153, 1},
	"lightsteelblue":       {176, 196, 222, 1},
	"lightyellow":          {255, 255, 224, 1},
	"lime":                 {0, 255, 0, 1},
	"limegreen":            {50, 205, 50, 1},
	"linen":                {250, 240, 230, 1},
	"magenta":              {255, 0, 255, 1},
	"maroon":               {128, 0, 0, 1},
	"mediumaquamarine":     {102, 205, 170, 1},
	"mediumblue":           {0, 0, 205, 1},
	"mediumorchid":         {186, 85, 211, 1},
	"mediumpurple":         {147, 112, 219, 1},
	"mediumseagreen":       {60, 179, 113, 1},
	"mediumslateblue":      {123, 104, 238, 1},
	"mediumspringgreen":    {0, 250, 154, 1},
	"mediumturquoise":      {72, 209, 204, 1},
	"mediumvioletred":      {199, 21, 133, 1},
	"midnightblue":         {25, 25, 112, 1},
	"mintcream":            {245, 255, 250, 1},
	"mistyrose":            {255, 228, 225, 1},
	"moccasin":             {255, 228, 181, 1},
	"navajowhite":          {255, 222, 173, 1},
	"navy":                 {0, 0, 128, 1},
	"oldlace":              {253, 245, 230, 1},
	"olive":                {128, 128, 0, 1},
	"olivedrab":            {107, 142, 35, 1},
	"orange":               {255, 165, 0, 1},
	"orangered":            {255, 69, 0, 1},
	"orchid":               {218, 112, 214, 1},
	"palegoldenrod":        {238, 232, 170, 1},
	"palegreen":            {152, 251, 152, 1},
	"paleturquoise":        {175, 238, 238, 1},
	"palevioletred":        {219, 112, 147, 1},
	"papayawhip":           {255, 239, 213, 1},
	"peachpuff":            {255, 218, 185, 1},
	"peru":                 {205, 133, 63, 1},
	"pink":                 {255, 192, 203, 1},
	"plum":                 {221, 160, 221, 1},
	"powderblue":           {176, 224, 230, 1},
	"purple":               {128, 0, 128, 1},
	"rebeccapurple":        {102, 51, 153, 1},
	"red":                  {255, 0, 0, 1},
	"rosybrown":            {188, 143, 143, 1},
	"royalblue":            {65, 105, 225, 1},
	"saddlebrown":          {139, 69, 19, 1},
	"salmon":               {250, 128, 114, 1},
	"sandybrown":           {244, 164, 96, 1},
	"seagreen":             {46, 139, 87, 1},
	"seashell":             {255, 245, 238, 1},
	"sienna":               {160, 82, 45, 1},
	"silver":               {192, 192, 192, 1},
	"skyblue":              {135, 206, 235, 1},
	"slateblue":            {106, 90, 205, 1},
	"slategray":            {112, 128, 144, 1},
	"slategrey":            {112, 128, 144, 1},
	"snow":                 {255, 250, 250, 1},
	"springgreen":          {0, 255, 127, 1},
	"steelblue":            {70, 130, 180, 1},
	"tan":                  {210, 180, 140, 1},
	"teal":                 {0, 128, 128, 1},
	"thistle":              {216, 191, 216, 1},
	"tomato":               {255, 99, 71, 1},
	"turquoise":            {64, 224, 208, 1},
	"violet":               {238, 130, 238, 1},
	"wheat":                {245, 222, 179, 1},
	"white":                {255, 255, 255, 1},
	"whitesmoke":           {245, 245, 245, 1},
	"yellow":               {255, 255, 0, 1},
	"yellowgreen":          {154, 205, 50, 1},
}

var hexRegex = regexp.MustCompile(`^#?([a-fA-F0-9]{3,8})$`)

// Parse parses a color string in hex (#RGB, #RRGGBB, #RRGGBBAA),
// RGB, HSL, or named color format.
func Parse(s string) (Color, error) {
	s = strings.TrimSpace(strings.ToLower(s))

	// Try named color
	if c, ok := namedColors[s]; ok {
		return c, nil
	}

	// Try hex
	if m := hexRegex.FindStringSubmatch(s); m != nil {
		return parseHex(m[1])
	}

	// Try rgb(r,g,b)
	if strings.HasPrefix(s, "rgb") {
		return parseRGB(s)
	}

	// Try hsl(h,s%,l%)
	if strings.HasPrefix(s, "hsl") {
		return parseHSL(s)
	}

	return Color{}, fmt.Errorf("unable to parse color: %s", s)
}

func parseHex(hex string) (Color, error) {
	hex = strings.TrimPrefix(hex, "#")

	// Expand 3-digit hex to 6-digit
	if len(hex) == 3 {
		hex = string(hex[0]) + string(hex[0]) +
			string(hex[1]) + string(hex[1]) +
			string(hex[2]) + string(hex[2])
	}

	switch len(hex) {
	case 6:
		r, err := strconv.ParseInt(hex[0:2], 16, 64)
		if err != nil {
			return Color{}, fmt.Errorf("invalid hex color: %s", hex)
		}
		g, err := strconv.ParseInt(hex[2:4], 16, 64)
		if err != nil {
			return Color{}, fmt.Errorf("invalid hex color: %s", hex)
		}
		b, err := strconv.ParseInt(hex[4:6], 16, 64)
		if err != nil {
			return Color{}, fmt.Errorf("invalid hex color: %s", hex)
		}
		return Color{int(r), int(g), int(b), 1}, nil
	case 8:
		r, err := strconv.ParseInt(hex[0:2], 16, 64)
		if err != nil {
			return Color{}, fmt.Errorf("invalid hex color: %s", hex)
		}
		g, err := strconv.ParseInt(hex[2:4], 16, 64)
		if err != nil {
			return Color{}, fmt.Errorf("invalid hex color: %s", hex)
		}
		b, err := strconv.ParseInt(hex[4:6], 16, 64)
		if err != nil {
			return Color{}, fmt.Errorf("invalid hex color: %s", hex)
		}
		a, err := strconv.ParseInt(hex[6:8], 16, 64)
		if err != nil {
			return Color{}, fmt.Errorf("invalid hex color: %s", hex)
		}
		return Color{int(r), int(g), int(b), float64(a) / 255.0}, nil
	default:
		return Color{}, fmt.Errorf("invalid hex color length: %d", len(hex))
	}
}

func parseRGB(s string) (Color, error) {
	s = strings.TrimSuffix(strings.TrimPrefix(s, "rgb("), ")")
	parts := strings.Split(s, ",")
	if len(parts) < 3 || len(parts) > 4 {
		return Color{}, fmt.Errorf("invalid RGB format: %s", s)
	}

	r, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return Color{}, fmt.Errorf("invalid RGB color: %s", s)
	}
	g, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return Color{}, fmt.Errorf("invalid RGB color: %s", s)
	}
	b, err := strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil {
		return Color{}, fmt.Errorf("invalid RGB color: %s", s)
	}

	a := 1.0
	if len(parts) == 4 {
		a, err = strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
		if err != nil || a < 0 || a > 1 {
			return Color{}, fmt.Errorf("invalid alpha in RGB: %s", s)
		}
	}

	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return Color{}, fmt.Errorf("RGB values out of range: %s", s)
	}

	return Color{r, g, b, a}, nil
}

func parseHSL(s string) (Color, error) {
	s = strings.TrimSuffix(strings.TrimPrefix(s, "hsl("), ")")
	parts := strings.Split(s, ",")
	if len(parts) < 3 || len(parts) > 4 {
		return Color{}, fmt.Errorf("invalid HSL format: %s", s)
	}

	h, err := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(parts[0]), "%"), 64)
	if err != nil {
		return Color{}, fmt.Errorf("invalid HSL color: %s", s)
	}
	sVal, err := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(parts[1]), "%"), 64)
	if err != nil {
		return Color{}, fmt.Errorf("invalid HSL color: %s", s)
	}
	l, err := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(parts[2]), "%"), 64)
	if err != nil {
		return Color{}, fmt.Errorf("invalid HSL color: %s", s)
	}

	sVal /= 100.0
	l /= 100.0

	c, a, x := hslToRGB(h, sVal, l)
	alpha := 1.0
	if len(parts) == 4 {
		alpha, err = strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
		if err != nil || alpha < 0 || alpha > 1 {
			return Color{}, fmt.Errorf("invalid alpha in HSL: %s", s)
		}
	}

	return Color{c, a, x, alpha}, nil
}

func hslToRGB(h, s, l float64) (int, int, int) {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}

	var r, g, b float64
	if s == 0 {
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

	return int(round(r * 255)), int(round(g * 255)), int(round(b * 255))
}

// Round helper
func round(f float64) float64 {
	if f < 0 {
		return math.Ceil(f - 0.5)
	}
	return math.Floor(f + 0.5)
}

// ToHSL converts RGB to HSL.
func (c Color) ToHSL() HSL {
	r, g, b := float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	l := (max + min) / 2.0

	if max == min {
		return HSL{0, 0, l * 100}
	}

	d := max - min
	var s float64
	if l > 0.5 {
		s = d / (2.0 - max - min)
	} else {
		s = d / (max + min)
	}

	var h float64
	switch max {
	case r:
		h = (g - b) / d
		if g < b {
			h += 6
		}
	case g:
		h = (b-r)/d + 2
	case b:
		h = (r-g)/d + 4
	}

	return HSL{h * 60, s * 100, l * 100}
}

// ToHSV converts RGB to HSV.
func (c Color) ToHSV() HSV {
	r, g, b := float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	v := max

	if max == 0 {
		return HSV{0, 0, 0}
	}

	s := (max - min) / max
	h := 0.0

	if max == min {
		h = 0
	} else {
		d := max - min
		switch max {
		case r:
			h = (g - b) / d
			if g < b {
				h += 6
			}
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}
		h *= 60
	}

	return HSV{h, s * 100, v * 100}
}

// Hex returns the hex string representation.
func (c Color) Hex() string {
	if c.A >= 1.0 {
		return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
	}
	return fmt.Sprintf("#%02X%02X%02X%02X", c.R, c.G, c.B, int(round(c.A*255)))
}

// RGB returns the RGB string.
func (c Color) RGB() string {
	if c.A >= 1.0 {
		return fmt.Sprintf("rgb(%d, %d, %d)", c.R, c.G, c.B)
	}
	return fmt.Sprintf("rgba(%d, %d, %d, %.2f)", c.R, c.G, c.B, c.A)
}

// String returns the hex representation.
func (c Color) String() string {
	return c.Hex()
}

// RelativeLuminance computes the relative luminance per WCAG 2.1.
func (c Color) RelativeLuminance() float64 {
	var linearize func(float64) float64
	linearize = func(c float64) float64 {
		if c <= 0.04045 {
			return c / 12.92
		}
		return math.Pow((c+0.055)/1.055, 2.4)
	}

	r := linearize(float64(c.R) / 255.0)
	g := linearize(float64(c.G) / 255.0)
	b := linearize(float64(c.B) / 255.0)

	return 0.2126*r + 0.7152*g + 0.0722*b
}

// ContrastRatio returns the WCAG contrast ratio between two colors.
func ContrastRatio(c1, c2 float64) float64 {
	l1, l2 := c1, c2
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}

// Contrast returns the contrast ratio between two colors.
func (c Color) Contrast(other Color) float64 {
	return ContrastRatio(c.RelativeLuminance(), other.RelativeLuminance())
}

// IsLight returns true if the color is considered light.
func (c Color) IsLight() bool {
	hsl := c.ToHSL()
	return hsl.L > 50
}

// IsDark returns true if the color is considered dark.
func (c Color) IsDark() bool {
	return !c.IsLight()
}

// Brightness returns the perceived brightness (0-100).
func (c Color) Brightness() float64 {
	return math.Round((0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)) / 2.55)
}

// ParseNamed returns a named color by name (case-insensitive).
func ParseNamed(name string) (Color, error) {
	key := strings.ToLower(strings.TrimSpace(html.UnescapeString(name)))
	if c, ok := namedColors[key]; ok {
		return c, nil
	}
	return Color{}, fmt.Errorf("unknown named color: %s", name)
}

// NamedColors returns a sorted list of all available named color keys.
func NamedColors() []string {
	names := make([]string, 0, len(namedColors))
	for name := range namedColors {
		names = append(names, name)
	}
	return names
}
