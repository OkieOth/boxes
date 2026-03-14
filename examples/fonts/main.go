package main

import (
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/okieoth/boxes/pkg/svg"
	"github.com/okieoth/boxes/pkg/svgdrawing"
	"github.com/okieoth/boxes/pkg/types"
)

func main() {
	outputFile := "temp/ExampleTexts.svg"
	output, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer output.Close()
	canvas := svg.New(output)
	canvas.Start(1000, 3000, "")
	canvas.Rect(0, 0, 1000, 3000, "fill:white")
	for i := 0; i < 100; i++ {
		canvas.Line(i*10, 5, i*10, 2995, "stroke:#adadad;stroke-width:1")
		txt := fmt.Sprintf("%d", i)
		canvas.Text(i*10+1, 5, txt, "text-anchor:start;font-size:5px;fill:#adadad")
	}
	basicTxt := "Hello, I am a text, that's used to calculate text bounding boxes"
	runeCount := utf8.RuneCountInString(basicTxt)

	styleInfo := "text-anchor:start;font-size:8px;fill:black"

	dc := svgdrawing.NewSvgTextDimensionCalculator()
	wb := types.FontDefWeightEnum_bold
	si := types.FontDefTypeEnum_italic

	for i := 0; i <= 10; i++ {
		y := 20 * (i + 1)
		fontSize := 10 + i
		style := fmt.Sprintf("text-anchor:start;font-family:Times New Roman, Times, Georgia, serif;font-size:%dpx;fill:black", fontSize)
		info := fmt.Sprintf("(font-size: %d, chars: %d", fontSize, runeCount)
		canvas.Text(10, y, basicTxt, style)
		// factor := float32(fontSize*10) * 1.95 / 5
		// canvas.Text(15+(int(factor)*runeCount/10), y, info, styleInfo)
		w, _ := dc.Dimensions(basicTxt, &types.FontDef{Size: fontSize, Font: "serif"})
		canvas.Text(w, y, info, styleInfo)
	}
	for i := 0; i <= 10; i++ {
		y := 250 + (20 * (i + 1))
		fontSize := 10 + i
		style := fmt.Sprintf("text-anchor:start;font-style: italic; font-family:Times New Roman, Times, Georgia, serif;font-size:%dpx;fill:black", fontSize)
		info := fmt.Sprintf("(font-size: %d, chars: %d", fontSize, runeCount)
		canvas.Text(10, y, basicTxt, style)
		// factor := float32(fontSize*10) * 1.95 / 5
		// canvas.Text(15+(int(factor)*runeCount/10), y, info, styleInfo)
		w, _ := dc.Dimensions(basicTxt, &types.FontDef{Size: fontSize, Font: "serif", Type: &si})
		canvas.Text(w, y, info, styleInfo)
	}
	for i := 0; i <= 10; i++ {
		y := 500 + (20 * (i + 1))
		fontSize := 10 + i
		style := fmt.Sprintf("text-anchor:start;font-weight: bold; font-family:Times New Roman, Times, Georgia, serif;font-size:%dpx;fill:black", fontSize)
		info := fmt.Sprintf("(font-size: %d, chars: %d", fontSize, runeCount)
		canvas.Text(10, y, basicTxt, style)
		// factor := float32(fontSize*10) * 2.1 / 5
		// canvas.Text(15+(int(factor)*runeCount/10), y, info, styleInfo)
		w, _ := dc.Dimensions(basicTxt, &types.FontDef{Size: fontSize, Font: "serif", Weight: &wb})
		canvas.Text(w, y, info, styleInfo)
	}
	for i := 0; i <= 10; i++ {
		y := 750 + (20 * (i + 1))
		fontSize := 10 + i
		style := fmt.Sprintf("text-anchor:start;font-weight: bold; font-style: italic; font-family:Times New Roman, Times, Georgia, serif;font-size:%dpx;fill:black", fontSize)
		info := fmt.Sprintf("(font-size: %d, chars: %d", fontSize, runeCount)
		canvas.Text(10, y, basicTxt, style)
		// factor := float32(fontSize*10) * 2.1 / 5
		// canvas.Text(15+(int(factor)*runeCount/10), y, info, styleInfo)
		w, _ := dc.Dimensions(basicTxt, &types.FontDef{Size: fontSize, Font: "serif", Weight: &wb, Type: &si})
		canvas.Text(w, y, info, styleInfo)
	}
	canvas.Line(0, 995, 1000, 995, "stroke:#adadad;stroke-width:1")
	for i := 0; i <= 10; i++ {
		y := 1000 + (20 * (i + 1))
		fontSize := 10 + i
		style := fmt.Sprintf("text-anchor:start;font-family:Arial, Helvetica, sans-serif;font-size:%dpx;fill:black", fontSize)
		info := fmt.Sprintf("(font-size: %d, chars: %d", fontSize, runeCount)
		canvas.Text(10, y, basicTxt, style)
		// factor := float32(fontSize*10) * 2.6 / 6
		// canvas.Text(15+(int(factor)*runeCount/10), y, info, styleInfo)
		w, _ := dc.Dimensions(basicTxt, &types.FontDef{Size: fontSize, Font: "sans-serif"})
		canvas.Text(15+w, y, info, styleInfo)
	}
	for i := 0; i <= 10; i++ {
		y := 1250 + (20 * (i + 1))
		fontSize := 10 + i
		style := fmt.Sprintf("text-anchor:start;font-style: italic; font-family:Arial, Helvetica, sans-serif;font-size:%dpx;fill:black", fontSize)
		info := fmt.Sprintf("(font-size: %d, chars: %d", fontSize, runeCount)
		canvas.Text(10, y, basicTxt, style)
		// factor := float32(fontSize*10) * 2.6 / 6
		// canvas.Text(15+(int(factor)*runeCount/10), y, info, styleInfo)
		w, _ := dc.Dimensions(basicTxt, &types.FontDef{Size: fontSize, Font: "sans-serif", Type: &si})
		canvas.Text(15+w, y, info, styleInfo)
	}
	for i := 0; i <= 10; i++ {
		y := 1500 + (20 * (i + 1))
		fontSize := 10 + i
		style := fmt.Sprintf("text-anchor:start;font-weight: bold; font-family:Arial, Helvetica, sans-serif;font-size:%dpx;fill:black", fontSize)
		info := fmt.Sprintf("(font-size: %d, chars: %d", fontSize, runeCount)
		canvas.Text(10, y, basicTxt, style)
		// factor := float32(fontSize*10) * 2.8 / 6
		// canvas.Text(15+(int(factor)*runeCount/10), y, info, styleInfo)
		w, _ := dc.Dimensions(basicTxt, &types.FontDef{Size: fontSize, Font: "sans-serif", Weight: &wb})
		canvas.Text(15+w, y, info, styleInfo)
	}
	for i := 0; i <= 10; i++ {
		y := 1750 + (20 * (i + 1))
		fontSize := 10 + i
		style := fmt.Sprintf("text-anchor:start;font-style: italic; font-weight: bold; font-family:Arial, Helvetica, sans-serif;font-size:%dpx;fill:black", fontSize)
		info := fmt.Sprintf("(font-size: %d, chars: %d", fontSize, runeCount)
		canvas.Text(10, y, basicTxt, style)
		// factor := float32(fontSize*10) * 2.8 / 6
		// canvas.Text(15+(int(factor)*runeCount/10), y, info, styleInfo)
		w, _ := dc.Dimensions(basicTxt, &types.FontDef{Size: fontSize, Font: "sans-serif", Weight: &wb, Type: &si})
		canvas.Text(15+w, y, info, styleInfo)
	}

	canvas.Line(0, 1995, 1000, 1995, "stroke:#adadad;stroke-width:1")
	for i := 0; i <= 10; i++ {
		y := 2000 + (20 * (i + 1))
		fontSize := 10 + i
		style := fmt.Sprintf("text-anchor:start;font-family:Courier New, Courier, Lucida Console, monospace;font-size:%dpx;fill:black", fontSize)
		info := fmt.Sprintf("(font-size: %d, chars: %d", fontSize, runeCount)
		canvas.Text(10, y, basicTxt, style)
		// factor := float32(fontSize*10) * 3 / 5
		// canvas.Text(15+(int(factor)*runeCount/10), y, info, styleInfo)
		w, _ := dc.Dimensions(basicTxt, &types.FontDef{Size: fontSize, Font: "monospace"})
		canvas.Text(15+w, y, info, styleInfo)
	}

	for i := 0; i <= 10; i++ {
		y := 2250 + (20 * (i + 1))
		fontSize := 10 + i
		style := fmt.Sprintf("text-anchor:start;font-style: italic; font-family:Courier New, Courier, Lucida Console, monospace;font-size:%dpx;fill:black", fontSize)
		info := fmt.Sprintf("(font-size: %d, chars: %d", fontSize, runeCount)
		canvas.Text(10, y, basicTxt, style)
		// factor := float32(fontSize*10) * 3 / 5
		// canvas.Text(15+(int(factor)*runeCount/10), y, info, styleInfo)
		w, _ := dc.Dimensions(basicTxt, &types.FontDef{Size: fontSize, Font: "monospace", Type: &si})
		canvas.Text(15+w, y, info, styleInfo)
	}
	for i := 0; i <= 10; i++ {
		y := 2500 + (20 * (i + 1))
		fontSize := 10 + i
		style := fmt.Sprintf("text-anchor:start;font-weight: bold; font-family:Courier New, Courier, Lucida Console, monospace;font-size:%dpx;fill:black", fontSize)
		info := fmt.Sprintf("(font-size: %d, chars: %d", fontSize, runeCount)
		canvas.Text(10, y, basicTxt, style)
		// factor := float32(fontSize*10) * 3 / 5
		// canvas.Text(15+(int(factor)*runeCount/10), y, info, styleInfo)
		w, _ := dc.Dimensions(basicTxt, &types.FontDef{Size: fontSize, Font: "monospace", Weight: &wb})
		canvas.Text(15+w, y, info, styleInfo)
	}
	for i := 0; i <= 10; i++ {
		y := 2750 + (20 * (i + 1))
		fontSize := 10 + i
		style := fmt.Sprintf("text-anchor:start;font-style: italic; font-weight: bold; font-family:Courier New, Courier, Lucida Console, monospace;font-size:%dpx;fill:black", fontSize)
		info := fmt.Sprintf("(font-size: %d, chars: %d", fontSize, runeCount)
		canvas.Text(10, y, basicTxt, style)
		// factor := float32(fontSize*10) * 3 / 5
		// canvas.Text(15+(int(factor)*runeCount/10), y, info, styleInfo)
		w, _ := dc.Dimensions(basicTxt, &types.FontDef{Size: fontSize, Font: "monospace", Weight: &wb, Type: &si})
		canvas.Text(15+w, y, info, styleInfo)
	}

	canvas.End()
}
