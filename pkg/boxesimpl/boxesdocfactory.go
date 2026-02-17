package boxesimpl

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/okieoth/draw.chart.things/pkg/types"
	"github.com/okieoth/draw.chart.things/pkg/types/boxes"
)

func initLayoutElemContainer(l []boxes.Layout, doc *boxes.BoxesDocument, b *boxes.Boxes) *boxes.LayoutElemContainer {
	if len(l) == 0 {
		return nil
	}
	var ret boxes.LayoutElemContainer
	ret.Elems = make([]boxes.LayoutElement, 0)
	for _, elem := range l {
		ret.Elems = append(ret.Elems, initLayoutElement(&elem, doc, b))
	}
	return &ret
}

func initConnections(l []boxes.Connection, inputFormats map[string]boxes.BoxFormat) []boxes.LayoutElemConnection {
	ret := make([]boxes.LayoutElemConnection, 0)
	defaultConnectionFormat, ok := inputFormats["defaultConnection"]
	for _, elem := range l {
		var conn boxes.LayoutElemConnection
		conn.DestId = elem.DestId
		conn.Comment = elem.Comment
		conn.SourceArrow = elem.SourceArrow
		conn.DestArrow = elem.DestArrow
		conn.HiddenComments = elem.HiddenComments
		conn.Tags = elem.Tags
		if elem.Format != nil {
			if formatInst, ok := inputFormats[*elem.Format]; ok {
				conn.Format = formatInst.Line
			}
		} else if ok {
			conn.Format = defaultConnectionFormat.Line
		}
		ret = append(ret, conn)
	}
	return ret
}

func initFontFormatAnchorInCase(f *types.FontDef) {
	if f != nil {
		if f.Anchor == "" {
			f.Anchor = types.FontDefAnchorEnum_middle
		}
	}
}

func initFontFormatAnchorLeftInCase(f *types.FontDef) {
	if f != nil {
		if f.Anchor == "" {
			f.Anchor = types.FontDefAnchorEnum_left
		}
	}
}

func initBoxFormat(f *boxes.Format) boxes.BoxFormat {
	var border *types.LineDef
	var fill *types.FillDef

	var fontCaption, fontText1, fontText2, fontComment, fontCommentMarker *types.FontDef
	var verticalTxt bool
	padding := types.GlobalPadding
	boxMargin := types.GlobalMinBoxMargin
	var fixedHeight, fixedWidth, cornerRadius *int
	var widthOfParent *bool
	if f != nil {
		initFontFormatAnchorInCase(f.FontCaption)
		initFontFormatAnchorInCase(f.FontText1)
		initFontFormatAnchorInCase(f.FontText2)
		initFontFormatAnchorLeftInCase(f.FontComment)
		initFontFormatAnchorInCase(f.FontCommentMarker)
		fontCaption = f.FontCaption
		fontText1 = f.FontText1
		fontText2 = f.FontText2
		border = f.Line
		fill = f.Fill
		if f.Padding != nil {
			padding = *f.Padding
		} else {
			f.Padding = &padding
		}
		if f.BoxMargin != nil {
			boxMargin = *f.BoxMargin
		} else {
			f.BoxMargin = &boxMargin
		}
		if f.VerticalTxt != nil {
			verticalTxt = *f.VerticalTxt
		}
		fixedHeight = f.FixedHeight
		fixedWidth = f.FixedWidth
		cornerRadius = f.CornerRadius
		widthOfParent = f.WidthOfParent
	}

	return boxes.BoxFormat{
		Padding:           padding,
		MinBoxMargin:      boxMargin,
		FontCaption:       types.InitFontDef(fontCaption, "sans-serif", 10, true, false, 0),
		FontText1:         types.InitFontDef(fontText1, "serif", 8, false, false, 10),
		FontText2:         types.InitFontDef(fontText2, "monospace", 8, false, true, 10),
		FontComment:       types.InitFontDef(fontComment, "serif", 10, false, false, 10),
		FontCommentMarker: types.InitFontDef(fontCommentMarker, "serif", 8, false, false, 10),
		Line:              border,
		Fill:              fill,
		FixedWidth:        fixedWidth,
		FixedHeight:       fixedHeight,
		WidthOfParent:     widthOfParent,
		VerticalTxt:       verticalTxt,
		CornerRadius:      cornerRadius,
	}
}

func adjustBoxFormat(f *boxes.BoxFormat, adjustment *boxes.Format) boxes.BoxFormat {
	var border *types.LineDef
	var fill *types.FillDef

	var fontCaption, fontText1, fontText2, fontComment, fontCommentMarker *types.FontDef
	var verticalTxt bool
	padding := types.GlobalPadding
	boxMargin := types.GlobalMinBoxMargin
	var fixedHeight, fixedWidth, cornerRadius *int
	var widthOfParent *bool
	if f != nil {
		fontCaption = &f.FontCaption
		fontText1 = &f.FontText1
		fontText2 = &f.FontText2
		border = types.CopyLineDef(f.Line)
		if f.Fill != nil {
			fill = types.CopyFillDef(f.Fill)
		}
		if f.Padding > 0 {
			padding = f.Padding
		}
		if f.MinBoxMargin > 0 {
			boxMargin = f.MinBoxMargin
		}
		fixedHeight = f.FixedHeight
		fixedWidth = f.FixedWidth
		widthOfParent = f.WidthOfParent
		cornerRadius = f.CornerRadius
	}

	ret := boxes.BoxFormat{
		Padding:           padding,
		MinBoxMargin:      boxMargin,
		FontCaption:       types.InitFontDef(fontCaption, "sans-serif", 10, true, false, 0),
		FontText1:         types.InitFontDef(fontText1, "serif", 8, false, false, 10),
		FontText2:         types.InitFontDef(fontText2, "monospace", 8, false, true, 10),
		FontComment:       types.InitFontDef(fontComment, "serif", 10, false, false, 10),
		FontCommentMarker: types.InitFontDef(fontCommentMarker, "serif", 8, false, false, 10),
		Line:              border,
		Fill:              fill,
		FixedWidth:        fixedWidth,
		FixedHeight:       fixedHeight,
		WidthOfParent:     widthOfParent,
		VerticalTxt:       verticalTxt,
		CornerRadius:      cornerRadius,
	}

	if adjustment != nil {
		if adjustment.Fill != nil {
			if adjustment.Fill.Color != nil {
				ret.Fill.Color = adjustment.Fill.Color
			}
			if adjustment.Fill.Opacity != nil {
				ret.Fill.Opacity = adjustment.Fill.Opacity
			}
		}
		if adjustment.Line != nil {
			if adjustment.Line.Color != nil {
				ret.Line.Color = adjustment.Line.Color
			}
			if adjustment.Line.Opacity != nil {
				ret.Line.Opacity = adjustment.Line.Opacity
			}
			if adjustment.Line.Width != nil {
				ret.Line.Width = adjustment.Line.Width
			}
			if adjustment.Line.Style != nil {
				ret.Line.Style = adjustment.Line.Style
			}
		}
	}

	return ret
}

func getDefaultCommentFormat() boxes.BoxFormat {
	return boxes.BoxFormat{
		FontComment:       types.InitFontDef(nil, "serif", 10, false, false, 10),
		FontCommentMarker: types.InitFontDef(nil, "serif", 10, false, false, 10),
		Line:              types.InitLineDef(nil),
		Fill:              types.InitFillDef(nil),
	}
}

func getDefaultFormat() boxes.BoxFormat {
	return boxes.BoxFormat{
		Padding:           types.GlobalPadding,
		MinBoxMargin:      types.GlobalMinBoxMargin,
		FontCaption:       types.InitFontDef(nil, "sans-serif", 10, true, false, 0),
		FontText1:         types.InitFontDef(nil, "serif", 8, false, false, 10),
		FontText2:         types.InitFontDef(nil, "monospace", 8, false, true, 10),
		FontComment:       types.InitFontDef(nil, "serif", 10, false, false, 10),
		FontCommentMarker: types.InitFontDef(nil, "serif", 10, false, false, 10),
		Line:              types.InitLineDef(nil),
		Fill:              nil,
	}
}

func initFormats(inputFormat map[string]boxes.Format) map[string]boxes.BoxFormat {
	var res = make(map[string]boxes.BoxFormat)
	for key, elem := range inputFormat {
		res[key] = initBoxFormat(&elem)
	}
	if _, hasDefault := res["default"]; !hasDefault {
		res["default"] = getDefaultFormat()
	}
	if _, hasDefault := res["defaultConnection"]; !hasDefault {
		res["defaultConnection"] = getDefaultFormat()
	}
	if _, hasDefault := res["defaultComment"]; !hasDefault {
		res["defaultComment"] = getDefaultCommentFormat()
	}
	return res
}

func initImage(l *boxes.Layout, definedImages map[string]types.ImageDef) *boxes.ImageContainer {
	if l.Image == nil {
		return nil
	}
	imageDef, ok := definedImages[*l.Image]
	if !ok {
		return nil
	}
	marginTopBottom := 5
	if imageDef.MarginTopBottom != nil {
		marginTopBottom = *imageDef.MarginTopBottom
	}
	marginLeftRight := 2
	if imageDef.MarginLeftRight != nil {
		marginLeftRight = *imageDef.MarginLeftRight
	}
	img := boxes.ImageContainer{
		Width:           imageDef.Width,
		Height:          imageDef.Height,
		ImgId:           *l.Image,
		MarginTopBottom: marginTopBottom,
		MarginLeftRight: marginLeftRight,
	}
	return &img
}

// func initLayoutElement(l *boxes.Layout, inputFormats map[string]boxes.BoxFormat, connectedIds *[]string, hideTextsForParents bool, definedImages map[string]types.ImageDef) boxes.LayoutElement {
func initLayoutElement(l *boxes.Layout, doc *boxes.BoxesDocument, b *boxes.Boxes) boxes.LayoutElement {
	var f *boxes.BoxFormat
	// for _, tag := range l.Tags {
	// 	if val, ok := inputFormats[tag]; ok {
	// 		f = &val
	// 		break
	// 	}
	// }
	if l.Format != nil {
		if val, ok := doc.Formats[*l.Format]; ok {
			f = &val
		}
	}
	if len(l.Connections) > 0 && l.Id != "" {
		if !slices.Contains(doc.ConnectedElems, l.Id) {
			doc.ConnectedElems = append(doc.ConnectedElems, l.Id)
		}
		for _, c := range l.Connections {
			if !slices.Contains(doc.ConnectedElems, c.DestId) {
				doc.ConnectedElems = append(doc.ConnectedElems, c.DestId)
			}
		}
	}
	if (f == nil) && (l.Caption != "" || l.Text1 != "" || l.Text2 != "") {
		formatKey := "default"
		if l.Format != nil {
			formatKey = *l.Format
		}
		for key, format := range doc.Formats {
			if key == formatKey {
				f = &format
				break
			}
		}
		if f == nil {
			d := getDefaultFormat()
			f = &d
		}
	}
	if l.Id == "" {
		if l.Caption != "" {
			l.Id = strings.ToLower(l.Caption)
		}
	}
	var text1, text2 string
	if !b.HideTextsForParents || (len(l.Vertical) == 0 && len(l.Horizontal) == 0) {
		text1 = l.Text1
		text2 = l.Text2
	}
	comments := make([]types.Comment, 0)
	if l.Comment != nil {
		comments = append(comments, *l.Comment)
	}
	return boxes.LayoutElement{
		Id:                l.Id,
		Caption:           l.Caption,
		Text1:             text1,
		Text2:             text2,
		Comments:          comments,
		HiddenComments:    l.HiddenComments,
		Image:             initImage(l, b.Images),
		Vertical:          initLayoutElemContainer(l.Vertical, doc, b),
		Horizontal:        initLayoutElemContainer(l.Horizontal, doc, b),
		Format:            adjustFormatBasedOnVariations(l, b, f, doc),
		DontBlockConPaths: l.DontBlockConPaths,
		DataLink:          l.DataLink,
		Connections:       initConnections(l.Connections, doc.Formats),
	}
}

func adjustFormatBasedOnVariations(l *boxes.Layout, b *boxes.Boxes, f *boxes.BoxFormat, doc *boxes.BoxesDocument) *boxes.BoxFormat {
	if b.FormatVariations != nil && b.FormatVariations.HasTag != nil {
		// mix in format adjustment
		var maxPriority *int
		var format2Use *boxes.Format
		for _, t := range l.Tags {
			if adjustment, ok := b.FormatVariations.HasTag[t]; ok {
				// needed adjustment found
				if maxPriority == nil || adjustment.Priority > *maxPriority {
					maxPriority = &adjustment.Priority
					if adjustment.FormatRef != nil {
						if f, ok := doc.Formats[*adjustment.FormatRef]; ok {
							bf := boxes.Format{
								Line:        f.Line,
								Fill:        f.Fill,
								FontCaption: &f.FontCaption,
							}
							format2Use = &bf
						}
					} else {
						format2Use = &adjustment.Format
					}
				}
			}
		}
		if format2Use != nil {
			newFormat := adjustBoxFormat(f, format2Use)
			return &newFormat
		}
	}
	return f
}

func initExternalImages(doc *boxes.BoxesDocument) error {
	for key, img := range doc.Images {
		if img.Base64 == nil && img.Base64Src == nil {
			return fmt.Errorf("Missing 'base64' or 'base64Src' attribute for imageDef id=%s", key)
		}
		if img.Base64Src != nil {
			// load the base64 string from the file given by the attrib
			bytes, err := os.ReadFile(*img.Base64Src)
			if err != nil {
				return fmt.Errorf("Error while reading content of 'base64Src' (%s) for imageDef id=%s", *img.Base64Src, key)
			}
			base64Str := string(bytes)
			img.Base64 = &base64Str
			doc.Images[key] = img
		}
	}
	return nil
}

func DocumentFromBoxes(b *boxes.Boxes) (*boxes.BoxesDocument, error) {
	doc := boxes.NewBoxesDocument()
	doc.Title = b.Title
	doc.TitleFormat = b.TitleFormat
	doc.Legend = b.Legend
	doc.Formats = initFormats(b.Formats)
	doc.Images = b.Images
	if b.MinBoxMargin != nil {
		doc.MinBoxMargin = *b.MinBoxMargin
	}
	if b.MinConnectorMargin != nil {
		doc.MinConnectorMargin = *b.MinConnectorMargin
	}
	if b.GlobalPadding != nil {
		doc.GlobalPadding = *b.GlobalPadding
	}
	if b.LineDist != nil {
		doc.LineDist = *b.LineDist
	}
	if err := initExternalImages(doc); err != nil {
		return nil, err
	}

	if err := initOverlays(b, doc); err != nil {
		return nil, err
	}

	doc.Boxes = initLayoutElement(&b.Boxes, doc, b)
	if doc.MinBoxMargin == 0 {
		doc.MinBoxMargin = types.GlobalMinBoxMargin
	}
	if doc.MinConnectorMargin == 0 {
		doc.MinConnectorMargin = types.GlobalMinConnectorMargin
	}
	if doc.GlobalPadding == 0 {
		doc.GlobalPadding = types.GlobalPadding
	}
	if doc.LineDist == 0 {
		doc.LineDist = types.LineDist
	}
	return doc, nil
}

func initOverlays(b *boxes.Boxes, doc *boxes.BoxesDocument) error {
	if len(b.Overlays) > 0 {
		for i := range b.Overlays {
			o := b.Overlays[i]
			newOverlay := boxes.NewDocOverlay()
			newOverlay.Caption = o.Caption
			newOverlay.RefValue = o.RefValue
			newOverlay.CenterXOffset = o.CenterXOffset
			newOverlay.CenterYOffset = o.CenterYOffset
			if o.Formats != nil {
				newOverlay.Formats = o.Formats
				slices.SortFunc(newOverlay.Formats.Gradations, func(c1, c2 boxes.OverlayGradation) int {
					if c1.Limit == c2.Limit {
						return 0
					} else if c1.Limit < c2.Limit {
						return -1
					} else {
						return 1
					}
				})
			} else {
				newOverlay.Formats = boxes.NewOverlayFormatDef()
			}
			max := float64(doc.Width / 4)
			min := float64(0)
			if o.RadiusDefs != nil && o.RadiusDefs.Max > 0 {
				max = o.RadiusDefs.Max
			}
			if o.RadiusDefs != nil && o.RadiusDefs.Min > 0 {
				min = o.RadiusDefs.Min
			}
			newOverlay.RadiusDefs = &boxes.OverlayRadiusDef{
				Min: min,
				Max: max, // max radius is when not set one quarter of the document width
			}
			for k, v := range o.Layouts {
				newOverlay.Layouts[k] = boxes.OverlayEntry{
					Value: v,
				}
			}
			doc.Overlays = append(doc.Overlays, *newOverlay)
		}
	}
	return nil
}
