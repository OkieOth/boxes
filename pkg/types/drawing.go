package types

// textAndDimensions represents text and its calculated dimensions
type TextAndDimensions struct {
	Text   string
	Width  int
	Height int
}

type TextDimensionCalculator interface {
	Dimensions(txt string, format *FontDef) (width, height int)
	SplitTxt(txt string, format *FontDef) (width, height int, lines []TextAndDimensions)
	DimensionsWithMaxWidth(txt string, format *FontDef, maxWidth int) (width, height int)
	SplitTxtWithMaxWidth(txt string, format *FontDef, maxWidth int) (width, height int, lines []TextAndDimensions)
}

type RectWithTextFormat struct {

	// Padding of the box
	Padding int `yaml:"padding"`

	FontCaption FontDef `yaml:"fontCaption"`

	FontText1 FontDef `yaml:"fontText1"`

	FontText2 FontDef `yaml:"fontText2"`

	Border *LineDef `yaml:"border,omitempty"`

	Fill *FillDef `yaml:"fill,omitempty"`

	// Minimum margin between boxes
	MinBoxMargin int `yaml:"minBoxMargin"`

	// If true, the text will be displayed vertically
	VerticalTxt bool `yaml:"verticalTxt"`

	CornerRadius *int `yaml:"cornerRadius"`
}

type Drawing interface {
	Start(title string, height, width int) error
	DrawRectWithText(id, caption, text1, text2 string, x, y, width, height, textYOffset int, format RectWithTextFormat, isLeaf bool) error
	DrawPng(x, y int, pngId string) error
	DrawPngWithAdditionalLink(x, y int, pngId, link string) error
	DrawLine(x1, y1, x2, y2 int, format LineDef) error
	DrawLineWithClass(x1, y1, x2, y2 int, format LineDef, className string) error
	DrawArrow(x, y, angle int, format LineDef) error
	DrawSolidRect(x, y, width, height int, fill *FillDef, line *LineDef) error
	DrawSolidCircle(x, y, radius int, color string) error
	DrawCircleWithBorder(x, y, radius int, fill *FillDef, line *LineDef) error
	DrawCircleWithBorderAndText(text string, x, y, radius int, fill *FillDef, line *LineDef, font *FontDef) error
	DrawCircleWithBorderTextAndClass(text string, x, y, radius int, fill *FillDef, line *LineDef, font *FontDef, className string) error
	DrawText(text string, x, y, width int, fontDef *FontDef) int
	DrawVerticalText(text string, currentX, y, height int, fontDef *FontDef) int
	Done() error
}
