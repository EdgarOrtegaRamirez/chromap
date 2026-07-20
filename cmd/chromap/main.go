// Package main provides the chromap CLI.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/EdgarOrtegaRamirez/chromap/internal/color"
	"github.com/EdgarOrtegaRamirez/chromap/internal/palette"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "parse":
		runParse(args)
	case "convert":
		runConvert(args)
	case "harmony":
		runHarmony(args)
	case "contrast":
		runContrast(args)
	case "colors":
		runListColors()
	case "version", "--version", "-v":
		fmt.Println("chromap", version)
	case "help", "--help", "-h", "":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`chromap — Color Palette Manager v1.0.0

Usage:
  chromap <command> [arguments]

Commands:
  parse <color>       Parse a color (hex, rgb, hsl, named)
  convert <color>     Convert color to all formats
  harmony <color>     Generate color harmonies
  contrast <c1> <c2>  Calculate contrast ratio between two colors
  colors              List available named colors
  version             Show version

Examples:
  chromap parse "#ff6b35"
  chromap parse "rebeccapurple"
  chromap parse "rgb(255, 107, 53)"
  chromap parse "hsl(19, 100%, 61%)"
  chromap convert "#ff6b35"
  chromap harmony "#4a90d9"
  chromap harmony "#4a90d9" analogous 7
  chromap contrast "#ffffff" "#333333"
  chromap harmony "#ff6b35" monochromatic 5 | chromap export css
`)
}

func runParse(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: chromap parse <color>\n")
		os.Exit(1)
	}

	c, err := color.Parse(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	hsl := c.ToHSL()
	hsv := c.ToHSV()
	fmt.Printf("Hex:     %s\n", c.Hex())
	fmt.Printf("RGB:     %s\n", c.RGB())
	fmt.Printf("HSL:     hsl(%.0f, %.0f%%, %.0f%%)\n", hsl.H, hsl.S, hsl.L)
	fmt.Printf("HSV:     hsv(%.0f, %.0f%%, %.0f%%)\n", hsv.H, hsv.S, hsv.V)
	fmt.Printf("Brightness: %.0f\n", c.Brightness())
	fmt.Printf("Light/Dark: %s\n", map[bool]string{true: "Light", false: "Dark"}[c.IsLight()])
	fmt.Printf("Luminance: %.4f\n", c.RelativeLuminance())
}

func runConvert(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: chromap convert <color>\n")
		os.Exit(1)
	}

	c, err := color.Parse(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== All formats ===")
	fmt.Printf("Hex:     %s\n", c.Hex())
	fmt.Printf("RGB:     %s\n", c.RGB())

	hsl := c.ToHSL()
	hsv := c.ToHSV()
	fmt.Printf("HSL:     hsl(%.0f, %.0f%%, %.0f%%)\n", hsl.H, hsl.S, hsl.L)
	fmt.Printf("HSV:     hsv(%.0f, %.0f%%, %.0f%%)\n", hsv.H, hsv.S, hsv.V)
	fmt.Printf("Luminance: %.4f\n", c.RelativeLuminance())
	fmt.Printf("Brightness: %.0f\n", c.Brightness())

	// Best text color
	text := palette.BestTextColor(c)
	fmt.Printf("Best text: %s\n", text.Hex())
}

func runHarmony(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: chromap harmony <color> [harmony] [size]\n")
		fmt.Println("\nHarmonies: complementary, analogous, triadic, tetradic, split-complementary, monochromatic")
		os.Exit(1)
	}

	baseStr := args[0]
	base, err := color.Parse(baseStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing base color: %v\n", err)
		os.Exit(1)
	}

	harm := palette.Complementary
	size := 5

	if len(args) > 1 {
		harm = palette.HarmonyType(strings.ToLower(args[1]))
	}
	switch harm {
	case "complementary", "analogous", "triadic", "tetradic",
		"split-complementary", "split-complement", "monochromatic":
	default:
		fmt.Fprintf(os.Stderr, "Unknown harmony: %s\n", harm)
		os.Exit(1)
	}

	if len(args) > 2 {
		fmt.Sscan(args[2], &size)
	}

	p := palette.Generate(base, harm, size)

	// Check if outputting to pipe (export mode)
	if len(os.Args) > 2 && os.Args[2] == "export" {
		// Will be handled by export command
		return
	}

	for _, c := range p {
		fmt.Printf("%s  %s\n", c.Hex(), map[bool]string{true: "Light", false: "Dark"}[c.IsLight()])
	}
}

func runContrast(args []string) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: chromap contrast <color1> <color2>\n")
		os.Exit(1)
	}

	c1, err := color.Parse(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing color1: %v\n", err)
		os.Exit(1)
	}

	c2, err := color.Parse(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing color2: %v\n", err)
		os.Exit(1)
	}

	ratio := c1.Contrast(c2)
	wcagAA, wcagAAA := "Fail", "Fail"
	if ratio >= 4.5 {
		wcagAA = "Pass"
	}
	if ratio >= 7.0 {
		wcagAAA = "Pass"
	}

	fmt.Printf("Contrast ratio: %.2f:1\n", ratio)
	fmt.Printf("WCAG AA (normal text): %s\n", wcagAA)
	fmt.Printf("WCAG AA (large text):  %s\n", map[bool]string{true: "Pass", false: "Fail"}[ratio >= 3.0])
	fmt.Printf("WCAG AAA (normal text): %s\n", wcagAAA)
	fmt.Printf("WCAG AAA (large text):  %s\n", map[bool]string{true: "Pass", false: "Fail"}[ratio >= 4.5])

	// Preview
	fmt.Printf("\n  %s  On %s: %s\n", c1.Hex(), c2.Hex(),
		map[bool]string{true: "Readable", false: "Not readable"}[ratio >= 4.5])
}

func runListColors() {
	names := color.NamedColors()
	fmt.Println("Available named colors:")
	for _, name := range names {
		c, _ := color.ParseNamed(name)
		fmt.Printf("  %-22s %s\n", name, c.Hex())
	}
}
