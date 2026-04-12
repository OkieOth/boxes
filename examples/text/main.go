package main

import (
	"os"

	"github.com/okieoth/boxes/pkg/svgdrawing"
	"github.com/okieoth/boxes/pkg/types"
)

func main() {
	f1 := types.FontDef{
		Size:              20,
		Font:              "serif",
		MaxLenBeforeBreak: 200,
		Type:              nil,
		Weight:            nil,
		Color:             "black",
		Anchor:            types.AnchorLeft,
	}
	var w float64 = 0.5
	lc := "pink"
	fl := types.LineDef{
		Color: &lc,
		Width: &w,
	}
	textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()

	outputFile := "temp/45_degree_text.svg"
	output, _ := os.Create(outputFile)
	drawing := svgdrawing.NewDrawing(output)
	drawing.Start("Test Calendar", 1000, 1000)
	testText := "Test Text"
	w1, _ := textDimensionCalulator.Dimensions(testText, &f1)
	x := 100
	y := 500
	drawing.DrawLine(x-5, y, x+5, y, fl)
	drawing.DrawLine(x, y-5, x, y+5, fl)
	drawing.DrawLine(x-5, y-w1, x+5, y-w1, fl)
	drawing.DrawVerticalText(testText, x, y, w1, &f1)

	x = 200
	f1.Anchor = types.AnchorMiddle
	drawing.DrawLine(x-5, y, x+5, y, fl)
	drawing.DrawLine(x, y-5, x, y+5, fl)
	drawing.DrawLine(x-5, y-w1, x+5, y-w1, fl)
	drawing.DrawVerticalText(testText, x, y, w1, &f1)

	x = 300
	f1.Anchor = types.AnchorRight
	drawing.DrawLine(x-5, y, x+5, y, fl)
	drawing.DrawLine(x, y-5, x, y+5, fl)
	drawing.DrawLine(x-5, y-w1, x+5, y-w1, fl)
	drawing.DrawVerticalText(testText, x, y, w1, &f1)

	drawing.Done()
	output.Close()

}
