package palette_test

import (
	"fmt"
	"testing"

	"github.com/EdgarOrtegaRamirez/chromap/internal/color"
	"github.com/EdgarOrtegaRamirez/chromap/internal/palette"
)

func TestDebugComplementary(t *testing.T) {
	base, err := color.Parse("#4a90d9")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Base: %s HSL: %+v", base.Hex(), base.ToHSL())

	p := palette.Generate(base, palette.Complementary, 5)
	t.Logf("Palette size: %d", len(p))
	for i, c := range p {
		t.Logf("[%d] %s HSL: %+v", i, c.Hex(), c.ToHSL())
	}

	// Print raw output
	var output string
	for _, c := range p {
		output += fmt.Sprintf("%s\n", c.Hex())
	}
	t.Logf("Output:\n%s", output)
}