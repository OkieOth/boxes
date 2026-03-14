package svgdrawing

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/okieoth/boxes/pkg/svg"
	"github.com/okieoth/boxes/pkg/types"
)

// Font width factor constants
const (
	MAGIC = 1.1
	// Sans-serif font width factors
	SANSSERIF_NORMAL_WIDTH_FACTOR = 2.6 / 6 * MAGIC
	SANSSERIF_BOLD_WIDTH_FACTOR   = 2.8 / 6 * MAGIC

	// Monospace font width factors
	MONOSPACE_NORMAL_WIDTH_FACTOR = 3.0 / 5 * MAGIC
	MONOSPACE_BOLD_WIDTH_FACTOR   = 3.0 / 5 * MAGIC

	// Serif font width factors
	SERIF_NORMAL_WIDTH_FACTOR = 2.05 / 5 * MAGIC
	SERIF_BOLD_WIDTH_FACTOR   = 2.15 / 5 * MAGIC
)

// Function type for calculating text dimensions
type calcDimensions func(runeCount int, fontSize int, bold bool) (width, height int)

// SvgTextDimensionCalculator calculates text dimensions for SVG rendering
type SvgTextDimensionCalculator struct {
	dimensionCalculators map[string]calcDimensions
}

// NewSvgTextDimensionCalculator creates a new text dimension calculator
func NewSvgTextDimensionCalculator() *SvgTextDimensionCalculator {
	calc := &SvgTextDimensionCalculator{
		dimensionCalculators: make(map[string]calcDimensions),
	}

	calc.dimensionCalculators["serif"] = serifDimensions
	calc.dimensionCalculators["sans-serif"] = sansserifDimensions
	calc.dimensionCalculators["monospace"] = monospaceDimensions

	return calc
}

// Font dimension calculators
func serifDimensions(runeCount int, fontSize int, bold bool) (width, height int) {
	var factor float32
	if bold {
		factor = float32(fontSize*10) * SERIF_BOLD_WIDTH_FACTOR
	} else {
		factor = float32(fontSize*10) * SERIF_NORMAL_WIDTH_FACTOR
	}
	w := int(factor) * runeCount / 10
	return w, fontSize
}

func sansserifDimensions(runeCount int, fontSize int, bold bool) (width, height int) {
	var factor float32
	if bold {
		factor = float32(fontSize*10) * SANSSERIF_BOLD_WIDTH_FACTOR
	} else {
		factor = float32(fontSize*10) * SANSSERIF_NORMAL_WIDTH_FACTOR
	}
	w := int(factor) * runeCount / 10
	return w, fontSize
}

func monospaceDimensions(runeCount int, fontSize int, bold bool) (width, height int) {
	var factor float32
	if bold {
		factor = float32(fontSize*10) * MONOSPACE_BOLD_WIDTH_FACTOR
	} else {
		factor = float32(fontSize*10) * MONOSPACE_NORMAL_WIDTH_FACTOR
	}
	w := int(factor) * runeCount / 10
	return w, fontSize
}

// SplitTxt splits text and calculates dimensions based on font settings
func (s *SvgTextDimensionCalculator) SplitTxt(txt string, format *types.FontDef) (width, height int, lines []types.TextAndDimensions) {
	return s.SplitTxtWithMaxWidth(txt, format, 0)
}

func (s *SvgTextDimensionCalculator) SplitTxtWithMaxWidth(txt string, format *types.FontDef, maxWidth int) (width, height int, lines []types.TextAndDimensions) {
	fontSize := format.Size
	if fontSize == 0 {
		fontSize = 10
	}

	bold := format.Weight != nil && *format.Weight == types.FontDefWeightEnum_bold
	lineHeight := format.LineHeight
	if lineHeight == 0 {
		lineHeight = 1.2
	}

	runeCount := utf8.RuneCountInString(txt)
	fontType := format.Font
	if fontType == "" {
		fontType = "sans-serif"
	}

	// Get appropriate dimension calculator or default to sans-serif
	dimCalc, exists := s.dimensionCalculators[fontType]
	if !exists {
		dimCalc = s.dimensionCalculators["sans-serif"]
	}

	w, h := dimCalc(runeCount, fontSize, bold)

	if maxWidth == 0 {
		maxWidth = format.MaxLenBeforeBreak
	}
	if w < maxWidth {
		lines = []types.TextAndDimensions{{Text: txt, Width: w, Height: h}}
		return w, h, lines
	}
	return splitTxtDimensions(w, runeCount, fontSize, bold, lineHeight, txt, maxWidth, dimCalc)
}

// Dimensions calculates text dimensions
func (s *SvgTextDimensionCalculator) Dimensions(txt string, format *types.FontDef) (width, height int) {
	w, h, lines := s.SplitTxt(txt, format)
	if len(lines) == 1 {
		return w, h //+ types.GlobalPadding // Eiko: removed for tighter layout
	}
	return w, h
}

func (s *SvgTextDimensionCalculator) DimensionsWithMaxWidth(txt string, format *types.FontDef, maxWidth int) (width, height int) {
	w, h, lines := s.SplitTxtWithMaxWidth(txt, format, maxWidth)
	if len(lines) == 1 {
		return w, h //+ types.GlobalPadding // Eiko: removed for tighter layout
	}
	return w, h
}

func splitTxtDimensions(
	originalWidth, runeCount, fontSize int,
	bold bool, lineHeight float32,
	txt string,
	maxLenBeforeBreak int,
	f calcDimensions) (width, height int, lines []types.TextAndDimensions) {

	words := strings.Fields(txt)
	wordsWithWidth := make([]types.TextAndDimensions, 0, len(words))

	// Calculate width for each word
	for _, w := range words {
		rc := utf8.RuneCount([]byte(w))
		width, _ := f(rc, fontSize, bold)

		// Truncate words that are too long
		if width > maxLenBeforeBreak {
			for j := rc - 1; j >= 0; j-- {
				w2, _ := f(j, fontSize, bold)
				if w2 <= maxLenBeforeBreak {
					w = string([]rune(w)[:j])
					width = w2
					break
				}
			}
		}

		wordsWithWidth = append(wordsWithWidth, types.TextAndDimensions{Text: w, Width: width})
	}

	// Build output lines
	curWidth := 0
	var line string
	var maxWidth int
	lines = make([]types.TextAndDimensions, 0)
	lh := int(float32(fontSize) * lineHeight)
	heightSum := 0

	appendLineToOutput := func(txt string, height int) []types.TextAndDimensions {
		lw, _ := f(utf8.RuneCount([]byte(txt)), fontSize, bold)
		if lw > maxWidth {
			maxWidth = lw
		}
		lines = append(lines, types.TextAndDimensions{Text: txt, Width: lw, Height: height})
		heightSum += lh
		return lines
	}

	for _, ww := range wordsWithWidth {
		if curWidth+ww.Width < maxLenBeforeBreak {
			if line != "" {
				line += " "
			}
			line += ww.Text
			curWidth += ww.Width
		} else {
			lines = appendLineToOutput(line, lh)
			line = ww.Text
			curWidth = ww.Width
		}
	}

	if line != "" {
		lines = appendLineToOutput(line, lh)
	}

	if len(lines) > 0 {
		lines[len(lines)-1].Height = fontSize // Eiko: adjustment to remove the lineHeight from the last line of the text
	}
	heightSum = (heightSum - lh) + fontSize
	return maxWidth, heightSum, lines
}

// Drawing represents an SVG drawing
type SvgDrawing struct {
	canvas                 *svg.SVG
	txtDimensionCalculator *SvgTextDimensionCalculator
}

// NewDrawing creates a new SVG drawing
func NewDrawing(w io.Writer) *SvgDrawing {
	return &SvgDrawing{
		canvas:                 svg.New(w),
		txtDimensionCalculator: NewSvgTextDimensionCalculator(),
	}
}

// Start initializes the SVG drawing
func (d *SvgDrawing) Start(title string, height, width int) error {
	args := fmt.Sprintf("viewBox=\"0,0,%d,%d\"", width, height)
	d.canvas.Start(width, height, args)
	return nil
}

func (d *SvgDrawing) InitImages(imageDefs map[string]types.ImageDef) error {
	if len(imageDefs) == 0 {
		return nil
	}
	// Example:
	// <defs>
	//   <image
	//     id="iconA"
	//     width="100"
	//     height="100"
	//     href="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
	//   />
	// </defs>
	d.canvas.Def()
	for key, imgDef := range imageDefs {
		d.canvas.PngWithIdBase64(0, 0, imgDef.Width, imgDef.Height, key, *imgDef.Base64)
	}
	defer d.canvas.DefEnd()
	return nil
}

func (d *SvgDrawing) DrawRaster(width, height, rasterSize int) {
	for x := rasterSize; x < width; x += rasterSize {
		for y := rasterSize; y < height; y += rasterSize {
			// vertical lines
			d.canvas.Line(x, y-1, x, y+1, "stroke:lightgray;stroke-width:0.5")
			// horizontal lines
			d.canvas.Line(x-1, y, x+1, y, "stroke:lightgray;stroke-width:0.5")
		}
	}
}

// textFormat generates CSS formatting for text
func (d *SvgDrawing) textFormat(fontDef *types.FontDef) string {
	var font string

	switch fontDef.Font {
	case "serif":
		font = "Times New Roman, Times, Georgia, serif"
	case "sans-serif":
		font = "Arial, Helvetica, sans-serif"
	case "monospace":
		font = "Courier New, Courier, Lucida Console, monospace"
	default:
		if fontDef.Font != "" {
			font = fontDef.Font
		} else {
			font = "Arial, Helvetica, sans-serif"
		}
	}

	anchor := fmt.Sprintf("%s", fontDef.Anchor)
	if anchor == "right" {
		anchor = "end"
	}
	txtFormat := fmt.Sprintf("text-anchor:%s;font-family:%s;font-size:%dpx;fill:%s",
		anchor, font, fontDef.Size, fontDef.Color)

	if fontDef.Weight != nil && *fontDef.Weight == types.FontDefWeightEnum_bold {
		txtFormat += ";font-weight:bold"
	}

	if fontDef.Type != nil && *fontDef.Type == types.FontDefTypeEnum_italic {
		txtFormat += ";font-style:italic"
	}

	return txtFormat
}

// drawText renders text with appropriate positioning and returns the updated y-position
func (d *SvgDrawing) DrawText(text string, x, currentY, width int, fontDef *types.FontDef) int {
	return d.DrawTextWithId(text, x, currentY, width, fontDef, "")
}

func (d *SvgDrawing) DrawTextWithIdAndClass(text string, x, currentY, width int, fontDef *types.FontDef, id, class string) int {
	if text == "" {
		return currentY
	}

	if fontDef.SpaceTop != 0 {
		currentY += fontDef.SpaceTop
	}

	txtFormat := d.textFormat(fontDef)
	_, _, lines := d.txtDimensionCalculator.SplitTxt(text, fontDef)

	var writeTxt func(int, int, string, ...string)
	if id == "" {
		writeTxt = d.canvas.Text
	} else {
		writeTxt = func(x int, y int, txt string, other ...string) {
			if class != "" {
				d.canvas.TextWithIdAndClass(id, x, y, txt, class, other...)
			} else {
				d.canvas.TextWithId(id, x, y, txt, other...)
			}
		}
	}

	xTxt := x + width/2
	for _, l := range lines {
		yTxt := currentY + fontDef.Size
		writeTxt(xTxt, yTxt, l.Text, txtFormat)
		currentY += l.Height
	}

	if fontDef.SpaceBottom != 0 {
		currentY += fontDef.SpaceBottom
	}

	return currentY
}

func (d *SvgDrawing) DrawTextWithId(text string, x, currentY, width int, fontDef *types.FontDef, id string) int {
	return d.DrawTextWithIdAndClass(text, x, currentY, width, fontDef, id, "")
}

func (d *SvgDrawing) DrawVerticalText(text string, currentX, y, height int, fontDef *types.FontDef) int {
	return d.DrawTextWithId(text, currentX, y, height, fontDef, "")
}

// drawText renders text with appropriate positioning and returns the updated X-position
func (d *SvgDrawing) DrawVerticalTextWithId(text string, currentX, y, height int, fontDef *types.FontDef, id string) int {
	if text == "" {
		return currentX
	}

	if fontDef.SpaceTop != 0 {
		currentX += fontDef.SpaceTop
	}

	txtFormat := d.textFormat(fontDef)
	_, _, lines := d.txtDimensionCalculator.SplitTxt(text, fontDef)

	for _, l := range lines {
		yTxt := y
		xTxt := currentX + fontDef.Size + 2
		// fmt.Sprintf("%s;transform=\"rotate(-90, %d, %d)\"", txtFormat, xTxt, yTxt)
		d.canvas.TextRotated(xTxt, yTxt, l.Text, -90, txtFormat)
		currentX += l.Height
	}

	if fontDef.SpaceBottom != 0 {
		currentX += fontDef.SpaceBottom
	}

	return currentX
}

func (d *SvgDrawing) DrawPng(x, y int, pngId string) error {
	d.canvas.Use(x, y, "#"+pngId)
	return nil
}

const ADDITIONAL_LINK string = "additionalLink"

func (d *SvgDrawing) DrawPngWithAdditionalLink(x, y int, pngId, link string) error {
	d.canvas.UseWithAdditionalDataLink(x, y, "#"+pngId, ADDITIONAL_LINK, link)
	return nil
}

// Draw renders a box with text elements
func (d *SvgDrawing) DrawRectWithText(id, caption, text1, text2 string, x, y, width, height, textYOffset int, format types.RectWithTextFormat) error {
	const onclickClass = "svg-clickable" // "onclick=\"window.shapeClick(event)\""
	if format.Fill != nil || format.Border != nil {
		attr := ""

		if format.Fill != nil {
			attr = fmt.Sprintf("fill: %s", *(*format.Fill).Color)
			if (*format.Fill).Opacity != nil {
				attr += fmt.Sprintf(";opacity: %f", *(*format.Fill).Opacity)
			}
		}

		if format.Border != nil {
			if attr != "" {
				attr += ";"
			}

			w := 1.0
			if (*format.Border).Width != nil {
				w = *(*format.Border).Width
			}

			c := "black"
			if (*format.Border).Color != nil {
				c = *(*format.Border).Color
			}

			attr += fmt.Sprintf("stroke: %s;stroke-Width: %f", c, w)
		}
		if id != "" {
			if format.CornerRadius != nil {
				d.canvas.RoundedRectWithIdAndClass(onclickClass, id, x, y, width, height, *format.CornerRadius, attr)
			} else {
				d.canvas.RectWithIdAndClass(onclickClass, id, x, y, width, height, attr)
			}
		} else {
			if format.CornerRadius != nil {
				d.canvas.Roundrect(x, y, width, height, *format.CornerRadius, *format.CornerRadius, attr)
			} else {
				d.canvas.Rect(x, y, width, height, attr)
			}

		}
	}

	idStr := ""
	if id != "" {
		idStr = id + "_capt"
	}

	if format.VerticalTxt {
		currentX := x + format.Padding
		currentY := y + (height / 2)
		currentX = d.DrawVerticalTextWithId(caption, currentX, currentY, width, &format.FontCaption, idStr)
		currentX = d.DrawVerticalText(text1, currentX, currentY, width, &format.FontText1)
		currentX = d.DrawVerticalText(text2, currentX, currentY, width, &format.FontText2)
	} else {
		currentY := y + format.Padding
		if textYOffset > 0 {
			currentY = y + textYOffset
		}
		currentY = d.DrawTextWithIdAndClass(caption, x, currentY, width, &format.FontCaption, idStr, onclickClass)
		currentY = d.DrawText(text1, x, currentY, width, &format.FontText1)
		currentY = d.DrawText(text2, x, currentY, width, &format.FontText2)
	}
	return nil
}

func lineStyleForType(style types.LineDefStyleEnum) string {
	switch style {
	case types.LineDefStyleEnum_dashed:
		return `stroke-dasharray:10 2;`
	case types.LineDefStyleEnum_dotted:
		return `stroke-dasharray:"1 10;"`
	default:
		return ""
	}
}

func (d *SvgDrawing) DrawLine(x1, y1, x2, y2 int, format types.LineDef) error {
	color := "black"
	if (format.Color != nil) && (*format.Color != "") {
		color = *format.Color
	}
	width := 1.0
	if (format.Width != nil) && (*format.Width != 0) {
		width = *format.Width
	}
	if format.Style != nil {
		d.canvas.Line(x1, y1, x2, y2, fmt.Sprintf("stroke:%s;stroke-width:%f;%s", color, width, lineStyleForType(*format.Style)))
	} else {
		d.canvas.Line(x1, y1, x2, y2, fmt.Sprintf("stroke:%s;stroke-width:%f", color, width))
	}
	return nil
}

func (d *SvgDrawing) DrawLineWithClass(x1, y1, x2, y2 int, format types.LineDef, className string) error {
	color := "black"
	if (format.Color != nil) && (*format.Color != "") {
		color = *format.Color
	}
	width := 1.0
	if (format.Width != nil) && (*format.Width != 0) {
		width = *format.Width
	}
	if format.Style != nil {
		d.canvas.LineWithClass(x1, y1, x2, y2, className, fmt.Sprintf("stroke:%s;stroke-width:%f;%s", color, width, lineStyleForType(*format.Style)))
	} else {
		d.canvas.LineWithClass(x1, y1, x2, y2, className, fmt.Sprintf("stroke:%s;stroke-width:%f", color, width))
	}
	return nil
}

func (d *SvgDrawing) DrawArrow(x, y, angle int, format types.LineDef) error {
	// TODO
	return nil
}

func (d *SvgDrawing) DrawSolidRect(x, y, width, height int, fill *types.FillDef, line *types.LineDef) error {
	attr := ""
	if fill != nil && fill.Color != nil {
		attr = fmt.Sprintf("fill: %s", *fill.Color)
		if fill.Opacity != nil {
			attr += fmt.Sprintf(";opacity: %f", *fill.Opacity)
		}
	}
	if line != nil && line.Color != nil {
		sep := ""
		if attr != "" {
			sep = ";"
		}
		attr = fmt.Sprintf("%s%sstroke: %s", attr, sep, *line.Color)
	}
	d.canvas.Rect(x, y, width, height, attr)
	return nil
}

func (d *SvgDrawing) DrawSolidCircle(x, y, radius int, color string) error {
	// Draw a solid circle at (x, y) with the specified color and radius
	attr := fmt.Sprintf("fill: %s", color)
	d.canvas.Circle(x, y, radius, attr)
	return nil
}

func strOrDefault(s *string, d string) string {
	if s != nil {
		return *s
	}
	return d
}

func floatOrDefault(s *float64, d float64) float64 {
	if s != nil {
		return *s
	}
	return d
}

func (d *SvgDrawing) getCircleStyle(fill *types.FillDef, line *types.LineDef) string {
	fc := strOrDefault(fill.Color, "white")
	fo := floatOrDefault(fill.Opacity, 1.0)
	lc := strOrDefault(line.Color, "black")
	lo := floatOrDefault(line.Opacity, 1.0)
	lw := floatOrDefault(line.Width, 0.5)
	return fmt.Sprintf("fill: %s;fill-opacity: %f; stroke: %s; stroke-width: %f; stroke-opacity: %f;", fc, fo, lc, lw, lo)

}

func (d *SvgDrawing) DrawCircleWithBorder(x, y, radius int, fill *types.FillDef, line *types.LineDef) error {
	attr := d.getCircleStyle(fill, line)
	d.canvas.Circle(x, y, radius, attr)
	return nil
}

func (d *SvgDrawing) DrawCircleWithBorderAndText(text string, x, y, radius int, fill *types.FillDef, line *types.LineDef, font *types.FontDef) error {
	return d.DrawCircleWithBorderTextAndClass(text, x, y, radius, fill, line, font, "")
}

func (d *SvgDrawing) DrawCircleWithBorderTextAndClass(text string, x, y, radius int, fill *types.FillDef, line *types.LineDef, font *types.FontDef, className string) error {
	attr := d.getCircleStyle(fill, line)
	textFormat := d.textFormat(font)
	textFormat += ";dominant-baseline: middle;"
	if className != "" {
		d.canvas.CircleWithClass(x, y, radius, className, attr)
		d.canvas.TextWithClass(x, y, text, className, textFormat)
	} else {
		d.canvas.Circle(x, y, radius, attr)
		d.canvas.Text(x, y, text, textFormat)
	}
	return nil
}

// Done finalizes the SVG document
func (d *SvgDrawing) Done() error {
	d.canvas.End()
	return nil
}
