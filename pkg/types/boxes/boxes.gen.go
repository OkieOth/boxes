// Attention, this code is generated. Do not modify manually. Changes will
// be overwritten be the next codegen run.
// Generated from Boxes (configs/models/boxes.json)



// Model to describe the input of block diagrams

package boxes
import (
    "github.com/okieoth/boxes/pkg/types"
)




// Definition of the output for the legend
type Legend struct {
Entries []LegendEntry `yaml:"entries,omitempty"`
    // format reference used for the legend texts
Format *string `yaml:"format,omitempty"`
}


func CopyLegend(src *Legend) *Legend {
    if src == nil {
        return nil
    }
    var ret Legend

    if src.Entries != nil {
        ret.Entries = make([]LegendEntry, len(src.Entries))
        for i, v := range src.Entries {
            ret.Entries[i] = *CopyLegendEntry(&v)
        }
    }

    if src.Format != nil {
        v := *src.Format
        ret.Format = &v
    }
return &ret
}


func NewLegend() *Legend {
    var ret Legend
    ret.Entries = make([]LegendEntry, 0)
    return &ret
}

// Definition of one legend entry
type LegendEntry struct {
Text string `yaml:"text"`
    // this format reference used to identify how the here described object is in the picture displayed.
Format string `yaml:"format"`
}


func CopyLegendEntry(src *LegendEntry) *LegendEntry {
    if src == nil {
        return nil
    }
    var ret LegendEntry

    ret.Text = src.Text

    ret.Format = src.Format
return &ret
}


func NewLegendEntry() *LegendEntry {
    var ret LegendEntry
    return &ret
}

type Layout struct {
    // unique identifier of that entry
Id string `yaml:"id"`
    // Some kind of the main text
Caption string `yaml:"caption"`
    // First additional text
Text1 string `yaml:"text1"`
    // Second additional text
Text2 string `yaml:"text2"`
    // additional comment, that can be then included in the created graphic
Comment *types.Comment `yaml:"comment,omitempty"`
    // Reference to an image that should be displayed, needs to be declared in the global image section
Image *string `yaml:"image,omitempty"`
    // in case the picture is rendered with given expanded IDs, and maxDepth, then if this flag is true, the box is still displayed expanded
Expand bool `yaml:"expand"`
    // If set, then the content for 'vertical' attrib is loaded from an external file
ExtVertical *string `yaml:"extVertical,omitempty"`
Vertical []Layout `yaml:"vertical,omitempty"`
    // If set, then the content for 'horizontal' attrib is loaded from an external file
ExtHorizontal *string `yaml:"extHorizontal,omitempty"`
Horizontal []Layout `yaml:"horizontal,omitempty"`
    // Tags to annotate the box, tags are used to format and filter
Tags []string `yaml:"tags,omitempty"`
    // List of connections to other boxes
Connections []Connection `yaml:"connections,omitempty"`
    // reference to the format to use for this box
Format *string `yaml:"format,omitempty"`
    // if that is set then connections can run through the box, as long as they don't cross the text
DontBlockConPaths *bool `yaml:"dontBlockConPaths,omitempty"`
    // Optional link to a source, related to this element. This can be used for instance for on-click handlers in a UI or simply as documentation.
DataLink *string `yaml:"dataLink,omitempty"`
    // is only set by while the layout is processed, don't set it in the definition
HiddenComments bool `yaml:"hiddenComments"`
}


func CopyLayout(src *Layout) *Layout {
    if src == nil {
        return nil
    }
    var ret Layout

    ret.Id = src.Id

    ret.Caption = src.Caption

    ret.Text1 = src.Text1

    ret.Text2 = src.Text2

    ret.Comment = types.CopyComment(src.Comment)

    if src.Image != nil {
        v := *src.Image
        ret.Image = &v
    }

    ret.Expand = src.Expand

    if src.ExtVertical != nil {
        v := *src.ExtVertical
        ret.ExtVertical = &v
    }

    if src.Vertical != nil {
        ret.Vertical = make([]Layout, len(src.Vertical))
        for i, v := range src.Vertical {
            ret.Vertical[i] = *CopyLayout(&v)
        }
    }

    if src.ExtHorizontal != nil {
        v := *src.ExtHorizontal
        ret.ExtHorizontal = &v
    }

    if src.Horizontal != nil {
        ret.Horizontal = make([]Layout, len(src.Horizontal))
        for i, v := range src.Horizontal {
            ret.Horizontal[i] = *CopyLayout(&v)
        }
    }

    if src.Tags != nil {
        ret.Tags = make([]string, len(src.Tags))
        for i, v := range src.Tags {
            ret.Tags[i] = v
        }
    }

    if src.Connections != nil {
        ret.Connections = make([]Connection, len(src.Connections))
        for i, v := range src.Connections {
            ret.Connections[i] = *CopyConnection(&v)
        }
    }

    if src.Format != nil {
        v := *src.Format
        ret.Format = &v
    }

    if src.DontBlockConPaths != nil {
        v := *src.DontBlockConPaths
        ret.DontBlockConPaths = &v
    }

    if src.DataLink != nil {
        v := *src.DataLink
        ret.DataLink = &v
    }

    ret.HiddenComments = src.HiddenComments
return &ret
}


func NewLayout() *Layout {
    var ret Layout
    ret.Vertical = make([]Layout, 0)
    ret.Horizontal = make([]Layout, 0)
    ret.Tags = make([]string, 0)
    ret.Connections = make([]Connection, 0)
    return &ret
}

// container to extend the layouts of a given layout element via mixins
type LayoutMixin struct {
Horizontal []Layout `yaml:"horizontal,omitempty"`
Vertical []Layout `yaml:"vertical,omitempty"`
}


func CopyLayoutMixin(src *LayoutMixin) *LayoutMixin {
    if src == nil {
        return nil
    }
    var ret LayoutMixin

    if src.Horizontal != nil {
        ret.Horizontal = make([]Layout, len(src.Horizontal))
        for i, v := range src.Horizontal {
            ret.Horizontal[i] = *CopyLayout(&v)
        }
    }

    if src.Vertical != nil {
        ret.Vertical = make([]Layout, len(src.Vertical))
        for i, v := range src.Vertical {
            ret.Vertical[i] = *CopyLayout(&v)
        }
    }
return &ret
}


func NewLayoutMixin() *LayoutMixin {
    var ret LayoutMixin
    ret.Horizontal = make([]Layout, 0)
    ret.Vertical = make([]Layout, 0)
    return &ret
}

type Connection struct {
    // Caption text of the destination box, can be used as alternative to 'destId'
Dest string `yaml:"dest"`
    // box id of the destination
DestId string `yaml:"destId"`
    // additional comment, that can be then included in the created graphic
Comment *types.Comment `yaml:"comment,omitempty"`
    // Arrow at the source box
SourceArrow bool `yaml:"sourceArrow"`
    // Arrow at the destination box
DestArrow bool `yaml:"destArrow"`
    // optional format to style the connection
Format *string `yaml:"format,omitempty"`
    // is only set by while the layout is processed, don't set it in the definition
HiddenComments bool `yaml:"hiddenComments"`
    // optional container to define additional contrains for the specific connection
ConnRestrictions *types.ConnRestriction `yaml:"connRestrictions,omitempty"`
    // Tags to annotate the connection, tags are used to format
Tags []string `yaml:"tags"`
    // optional step where this comment is part of, is filled via processing not by the user
Step *int `yaml:"step,omitempty"`
}


func CopyConnection(src *Connection) *Connection {
    if src == nil {
        return nil
    }
    var ret Connection

    ret.Dest = src.Dest

    ret.DestId = src.DestId

    ret.Comment = types.CopyComment(src.Comment)

    ret.SourceArrow = src.SourceArrow

    ret.DestArrow = src.DestArrow

    if src.Format != nil {
        v := *src.Format
        ret.Format = &v
    }

    ret.HiddenComments = src.HiddenComments

    ret.ConnRestrictions = types.CopyConnRestriction(src.ConnRestrictions)

    if src.Tags != nil {
        ret.Tags = make([]string, len(src.Tags))
        for i, v := range src.Tags {
            ret.Tags[i] = v
        }
    }

    if src.Step != nil {
        v := *src.Step
        ret.Step = &v
    }
return &ret
}


func NewConnection() *Connection {
    var ret Connection
    ret.Tags = make([]string, 0)
    return &ret
}

type Format struct {
    // sets the width of the object to the width of the parent
WidthOfParent *bool `yaml:"widthOfParent,omitempty"`
    // optional fixed width that will be applied on the box
FixedWidth *int `yaml:"fixedWidth,omitempty"`
    // optional fixed height that will be applied on the box
FixedHeight *int `yaml:"fixedHeight,omitempty"`
    // If true, the text will be displayed vertically
VerticalTxt *bool `yaml:"verticalTxt,omitempty"`
FontCaption *types.FontDef `yaml:"fontCaption,omitempty"`
FontText1 *types.FontDef `yaml:"fontText1,omitempty"`
FontText2 *types.FontDef `yaml:"fontText2,omitempty"`
FontComment *types.FontDef `yaml:"fontComment,omitempty"`
FontCommentMarker *types.FontDef `yaml:"fontCommentMarker,omitempty"`
Line *types.LineDef `yaml:"line,omitempty"`
Fill *types.FillDef `yaml:"fill,omitempty"`
    // Padding used for this format
Padding *int `yaml:"padding,omitempty"`
    // Minimum margin between boxes
BoxMargin *int `yaml:"boxMargin,omitempty"`
    // radius of the box corners in pixel
CornerRadius *int `yaml:"cornerRadius,omitempty"`
}


func CopyFormat(src *Format) *Format {
    if src == nil {
        return nil
    }
    var ret Format

    if src.WidthOfParent != nil {
        v := *src.WidthOfParent
        ret.WidthOfParent = &v
    }

    if src.FixedWidth != nil {
        v := *src.FixedWidth
        ret.FixedWidth = &v
    }

    if src.FixedHeight != nil {
        v := *src.FixedHeight
        ret.FixedHeight = &v
    }

    if src.VerticalTxt != nil {
        v := *src.VerticalTxt
        ret.VerticalTxt = &v
    }

    ret.FontCaption = types.CopyFontDef(src.FontCaption)

    ret.FontText1 = types.CopyFontDef(src.FontText1)

    ret.FontText2 = types.CopyFontDef(src.FontText2)

    ret.FontComment = types.CopyFontDef(src.FontComment)

    ret.FontCommentMarker = types.CopyFontDef(src.FontCommentMarker)

    ret.Line = types.CopyLineDef(src.Line)

    ret.Fill = types.CopyFillDef(src.Fill)

    if src.Padding != nil {
        v := *src.Padding
        ret.Padding = &v
    }

    if src.BoxMargin != nil {
        v := *src.BoxMargin
        ret.BoxMargin = &v
    }

    if src.CornerRadius != nil {
        v := *src.CornerRadius
        ret.CornerRadius = &v
    }
return &ret
}


func NewFormat() *Format {
    var ret Format
    return &ret
}

type FormatVariations struct {
    // dictionary with tag values as key, that contains format definitions.
HasTag map[string]FormatVariation `yaml:"hasTag,omitempty"`
}


func CopyFormatVariations(src *FormatVariations) *FormatVariations {
    if src == nil {
        return nil
    }
    var ret FormatVariations

    if src.HasTag != nil {
        ret.HasTag = make(map[string]FormatVariation, len(src.HasTag))
        for k, v := range src.HasTag {
            ret.HasTag[k] = *CopyFormatVariation(&v)
        }
    }
return &ret
}


func NewFormatVariations() *FormatVariations {
    var ret FormatVariations
    ret.HasTag = make(map[string]FormatVariation, 0)
    return &ret
}

type FormatVariation struct {
    // fill this in case the format is specified outside
FormatRef *string `yaml:"formatRef,omitempty"`
Format Format `yaml:"format"`
    // number to define the order if a layout has for instance multiple matching tags
Priority int `yaml:"priority"`
}


func CopyFormatVariation(src *FormatVariation) *FormatVariation {
    if src == nil {
        return nil
    }
    var ret FormatVariation

    if src.FormatRef != nil {
        v := *src.FormatRef
        ret.FormatRef = &v
    }

    ret.Format = *CopyFormat(&src.Format)

    ret.Priority = src.Priority
return &ret
}


func NewFormatVariation() *FormatVariation {
    var ret FormatVariation
    return &ret
}

// Definition of a topic related overlay ... for instance for heatmaps
type Overlay struct {
    // some catchy words to describe the displayed topc
Caption string `yaml:"caption"`
    // Optional reference value that defines the reference value for this type of overlay
RefValue float64 `yaml:"refValue"`
    // radius for having a value of refValue
RefRadius *float64 `yaml:"refRadius,omitempty"`
    // in case of multiple overlays existing, this allows to define a percentage offset from the center-x of the related layout object
CenterXOffset float64 `yaml:"centerXOffset"`
    // in case of multiple overlays existing, this allows to define a percentage offset from the center-y of the related layout object
CenterYOffset float64 `yaml:"centerYOffset"`
    // if set to true, then the value is printed in the overlay
PrintValue bool `yaml:"printValue"`
    // dictionary of layout elements, that contain this overlay. The dictionary stores the value for this specific object
Layouts map[string]float64 `yaml:"layouts"`
    // if this is configured the the radius for the layouts is in a percentage of the refValue
RadiusDefs *OverlayRadiusDef `yaml:"radiusDefs,omitempty"`
Formats *OverlayFormatDef `yaml:"formats,omitempty"`
}


func CopyOverlay(src *Overlay) *Overlay {
    if src == nil {
        return nil
    }
    var ret Overlay

    ret.Caption = src.Caption

    ret.RefValue = src.RefValue

    if src.RefRadius != nil {
        v := *src.RefRadius
        ret.RefRadius = &v
    }

    ret.CenterXOffset = src.CenterXOffset

    ret.CenterYOffset = src.CenterYOffset

    ret.PrintValue = src.PrintValue

    if src.Layouts != nil {
        ret.Layouts = make(map[string]float64, len(src.Layouts))
        for k, v := range src.Layouts {
            ret.Layouts[k] = v
        }
    }

    ret.RadiusDefs = CopyOverlayRadiusDef(src.RadiusDefs)

    ret.Formats = CopyOverlayFormatDef(src.Formats)
return &ret
}


func NewOverlay() *Overlay {
    var ret Overlay
    ret.Layouts = make(map[string]float64, 0)
    return &ret
}

// definition how to calculate radius changes based on a reference value
type OverlayRadiusDef struct {
    // minimal radius to use for the display
Min float64 `yaml:"min"`
    // maximal radius to use for the display
Max float64 `yaml:"max"`
}


func CopyOverlayRadiusDef(src *OverlayRadiusDef) *OverlayRadiusDef {
    if src == nil {
        return nil
    }
    var ret OverlayRadiusDef

    ret.Min = src.Min

    ret.Max = src.Max
return &ret
}


func NewOverlayRadiusDef() *OverlayRadiusDef {
    var ret OverlayRadiusDef
    return &ret
}

// definition what format to use for a specific reference value
type OverlayFormatDef struct {
    // default format to use to display the overlay
Default string `yaml:"default"`
    // grations considered for switching for formats to use
Gradations []OverlayGradation `yaml:"gradations,omitempty"`
}


func CopyOverlayFormatDef(src *OverlayFormatDef) *OverlayFormatDef {
    if src == nil {
        return nil
    }
    var ret OverlayFormatDef

    ret.Default = src.Default

    if src.Gradations != nil {
        ret.Gradations = make([]OverlayGradation, len(src.Gradations))
        for i, v := range src.Gradations {
            ret.Gradations[i] = *CopyOverlayGradation(&v)
        }
    }
return &ret
}


func NewOverlayFormatDef() *OverlayFormatDef {
    var ret OverlayFormatDef
    ret.Gradations = make([]OverlayGradation, 0)
    return &ret
}

// gradation entry for switching formats for overlays
type OverlayGradation struct {
    // to what value should the here named format being used
Limit float64 `yaml:"limit"`
    // name of the defined format to use
Format string `yaml:"format"`
}


func CopyOverlayGradation(src *OverlayGradation) *OverlayGradation {
    if src == nil {
        return nil
    }
    var ret OverlayGradation

    ret.Limit = src.Limit

    ret.Format = src.Format
return &ret
}


func NewOverlayGradation() *OverlayGradation {
    var ret OverlayGradation
    return &ret
}

// Model to describe the input of block diagrams
type Boxes struct {
    // Title of the document
Title string `yaml:"title"`
    // format reference used for the title
TitleFormat *string `yaml:"titleFormat,omitempty"`
    // allows to include a version for the layout description
Version *string `yaml:"version,omitempty"`
    // Legend definition used in this diagram
Legend *Legend `yaml:"legend,omitempty"`
Boxes Layout `yaml:"boxes"`
    // Map of formats available to be used in the boxes
Formats map[string]Format `yaml:"formats,omitempty"`
    // Set of formats that overwrites the style of boxes, if specific conditions are met
FormatVariations *FormatVariations `yaml:"formatVariations,omitempty"`
    // optional map of images used in the generated graphic
Images map[string]types.ImageDef `yaml:"images,omitempty"`
    // If that is set, then the additional texts are only visible when the box has no visible children
HideTextsForParents bool `yaml:"hideTextsForParents"`
    // minimal distance between overlapping lines
LineDist *int `yaml:"lineDist,omitempty"`
    // Padding used as default over the whole diagram
GlobalPadding *int `yaml:"globalPadding,omitempty"`
    // Minimum margin between boxes
MinBoxMargin *int `yaml:"minBoxMargin,omitempty"`
    // Minimum margin between connectors
MinConnectorMargin *int `yaml:"minConnectorMargin,omitempty"`
Overlays []Overlay `yaml:"overlays,omitempty"`
}


func CopyBoxes(src *Boxes) *Boxes {
    if src == nil {
        return nil
    }
    var ret Boxes

    ret.Title = src.Title

    if src.TitleFormat != nil {
        v := *src.TitleFormat
        ret.TitleFormat = &v
    }

    if src.Version != nil {
        v := *src.Version
        ret.Version = &v
    }

    ret.Legend = CopyLegend(src.Legend)

    ret.Boxes = *CopyLayout(&src.Boxes)

    if src.Formats != nil {
        ret.Formats = make(map[string]Format, len(src.Formats))
        for k, v := range src.Formats {
            ret.Formats[k] = *CopyFormat(&v)
        }
    }

    ret.FormatVariations = CopyFormatVariations(src.FormatVariations)

    if src.Images != nil {
        ret.Images = make(map[string]types.ImageDef, len(src.Images))
        for k, v := range src.Images {
            ret.Images[k] = *types.CopyImageDef(&v)
        }
    }

    ret.HideTextsForParents = src.HideTextsForParents

    if src.LineDist != nil {
        v := *src.LineDist
        ret.LineDist = &v
    }

    if src.GlobalPadding != nil {
        v := *src.GlobalPadding
        ret.GlobalPadding = &v
    }

    if src.MinBoxMargin != nil {
        v := *src.MinBoxMargin
        ret.MinBoxMargin = &v
    }

    if src.MinConnectorMargin != nil {
        v := *src.MinConnectorMargin
        ret.MinConnectorMargin = &v
    }

    if src.Overlays != nil {
        ret.Overlays = make([]Overlay, len(src.Overlays))
        for i, v := range src.Overlays {
            ret.Overlays[i] = *CopyOverlay(&v)
        }
    }
return &ret
}


func NewBoxes() *Boxes {
    var ret Boxes
    ret.Formats = make(map[string]Format, 0)
    ret.Images = make(map[string]types.ImageDef, 0)
    ret.Overlays = make([]Overlay, 0)
    return &ret
}
