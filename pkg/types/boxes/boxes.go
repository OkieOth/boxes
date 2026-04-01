package boxes

// Attention, this file is generated. Manual changes get lost with the next
// run of the code generation.
// created by yacg (template: golang_types.mako v1.1.0)

import (
    "github.com/okieoth/boxes/pkg/types"
)


/* Model to describe the input of block diagrams
*/
type Boxes struct {

    // Title of the document
    Title string  `yaml:"title"`

    // format reference used for the title
    TitleFormat *string  `yaml:"titleFormat,omitempty"`

    // allows to include a version for the layout description
    Version *string  `yaml:"version,omitempty"`

    // Legend definition used in this diagram
    Legend *Legend  `yaml:"legend,omitempty"`

    Boxes Layout  `yaml:"boxes"`

    // Map of formats available to be used in the boxes
    Formats map[string]Format  `yaml:"formats,omitempty"`

    // Set of formats that overwrites the style of boxes, if specific conditions are met
    FormatVariations *FormatVariations  `yaml:"formatVariations,omitempty"`

    // optional map of images used in the generated graphic
    Images map[string]types.ImageDef  `yaml:"images,omitempty"`

    // If that is set, then the additional texts are only visible when the box has no visible children
    HideTextsForParents bool  `yaml:"hideTextsForParents"`

    // minimal distance between overlapping lines
    LineDist *int  `yaml:"lineDist,omitempty"`

    // Padding used as default over the whole diagram
    GlobalPadding *int  `yaml:"globalPadding,omitempty"`

    // Minimum margin between boxes
    MinBoxMargin *int  `yaml:"minBoxMargin,omitempty"`

    // Minimum margin between connectors
    MinConnectorMargin *int  `yaml:"minConnectorMargin,omitempty"`

    Overlays []Overlay  `yaml:"overlays,omitempty"`
}

func NewBoxes() *Boxes {
    return &Boxes{
        Legend: NewLegend(),
        Boxes: *NewLayout(),
        Formats: make(map[string]Format, 0),
        FormatVariations: NewFormatVariations(),
        Images: make(map[string]types.ImageDef, 0),
        Overlays: make([]Overlay, 0),
    }
}

func CopyBoxes(src *Boxes) *Boxes {
    if src == nil {
        return nil
    }
    var ret Boxes
    ret.Title = src.Title
    ret.TitleFormat = src.TitleFormat
    ret.Version = src.Version
    ret.Legend = CopyLegend(src.Legend)
    ret.Boxes = *CopyLayout(&src.Boxes)
    ret.Formats = make(map[string]Format, 0)
    for k, v := range src.Formats {
        ret.Formats[k] = v
    }
    ret.FormatVariations = CopyFormatVariations(src.FormatVariations)
    ret.Images = make(map[string]types.ImageDef, 0)
    for k, v := range src.Images {
        ret.Images[k] = v
    }
    ret.HideTextsForParents = src.HideTextsForParents
    ret.LineDist = src.LineDist
    ret.GlobalPadding = src.GlobalPadding
    ret.MinBoxMargin = src.MinBoxMargin
    ret.MinConnectorMargin = src.MinConnectorMargin
    ret.Overlays = make([]Overlay, 0)
    for _, e := range src.Overlays {
        ret.Overlays = append(ret.Overlays, e)
    }

    return &ret
}





/* Definition of the output for the legend
*/
type Legend struct {

    Entries []LegendEntry  `yaml:"entries,omitempty"`

    // format reference used for the legend texts
    Format *string  `yaml:"format,omitempty"`
}

func NewLegend() *Legend {
    return &Legend{
        Entries: make([]LegendEntry, 0),
    }
}

func CopyLegend(src *Legend) *Legend {
    if src == nil {
        return nil
    }
    var ret Legend
    ret.Entries = make([]LegendEntry, 0)
    for _, e := range src.Entries {
        ret.Entries = append(ret.Entries, e)
    }
    ret.Format = src.Format

    return &ret
}





type Layout struct {

    // unique identifier of that entry
    Id string  `yaml:"id"`

    // Some kind of the main text
    Caption string  `yaml:"caption"`

    // First additional text
    Text1 string  `yaml:"text1"`

    // Second additional text
    Text2 string  `yaml:"text2"`

    // additional comment, that can be then included in the created graphic
    Comment *types.Comment  `yaml:"comment,omitempty"`

    // Reference to an image that should be displayed, needs to be declared in the global image section
    Image *string  `yaml:"image,omitempty"`

    // in case the picture is rendered with given expanded IDs, and maxDepth, then if this flag is true, the box is still displayed expanded
    Expand bool  `yaml:"expand"`

    // If set, then the content for 'vertical' attrib is loaded from an external file
    ExtVertical *string  `yaml:"extVertical,omitempty"`

    Vertical []Layout  `yaml:"vertical,omitempty"`

    // If set, then the content for 'horizontal' attrib is loaded from an external file
    ExtHorizontal *string  `yaml:"extHorizontal,omitempty"`

    Horizontal []Layout  `yaml:"horizontal,omitempty"`

    // Tags to annotate the box, tags are used to format and filter
    Tags []string  `yaml:"tags,omitempty"`

    // List of connections to other boxes
    Connections []Connection  `yaml:"connections,omitempty"`

    // reference to the format to use for this box
    Format *string  `yaml:"format,omitempty"`

    // if that is set then connections can run through the box, as long as they don't cross the text
    DontBlockConPaths *bool  `yaml:"dontBlockConPaths,omitempty"`

    // Optional link to a source, related to this element. This can be used for instance for on-click handlers in a UI or simply as documentation.
    DataLink *string  `yaml:"dataLink,omitempty"`

    // is only set by while the layout is processed, don't set it in the definition
    HiddenComments bool  `yaml:"hiddenComments"`
}

func NewLayout() *Layout {
    return &Layout{
        Vertical: make([]Layout, 0),
        Horizontal: make([]Layout, 0),
        Tags: make([]string, 0),
        Connections: make([]Connection, 0),
    }
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
    ret.Image = src.Image
    ret.Expand = src.Expand
    ret.ExtVertical = src.ExtVertical
    ret.Vertical = make([]Layout, 0)
    for _, e := range src.Vertical {
        ret.Vertical = append(ret.Vertical, e)
    }
    ret.ExtHorizontal = src.ExtHorizontal
    ret.Horizontal = make([]Layout, 0)
    for _, e := range src.Horizontal {
        ret.Horizontal = append(ret.Horizontal, e)
    }
    ret.Tags = make([]string, 0)
    for _, e := range src.Tags {
        ret.Tags = append(ret.Tags, e)
    }
    ret.Connections = make([]Connection, 0)
    for _, e := range src.Connections {
        ret.Connections = append(ret.Connections, e)
    }
    ret.Format = src.Format
    ret.DontBlockConPaths = src.DontBlockConPaths
    ret.DataLink = src.DataLink
    ret.HiddenComments = src.HiddenComments

    return &ret
}








type Format struct {

    // sets the width of the object to the width of the parent
    WidthOfParent *bool  `yaml:"widthOfParent,omitempty"`

    // optional fixed width that will be applied on the box
    FixedWidth *int  `yaml:"fixedWidth,omitempty"`

    // optional fixed height that will be applied on the box
    FixedHeight *int  `yaml:"fixedHeight,omitempty"`

    // If true, the text will be displayed vertically
    VerticalTxt *bool  `yaml:"verticalTxt,omitempty"`

    FontCaption *types.FontDef  `yaml:"fontCaption,omitempty"`

    FontText1 *types.FontDef  `yaml:"fontText1,omitempty"`

    FontText2 *types.FontDef  `yaml:"fontText2,omitempty"`

    FontComment *types.FontDef  `yaml:"fontComment,omitempty"`

    FontCommentMarker *types.FontDef  `yaml:"fontCommentMarker,omitempty"`

    Line *types.LineDef  `yaml:"line,omitempty"`

    Fill *types.FillDef  `yaml:"fill,omitempty"`

    // Padding used for this format
    Padding *int  `yaml:"padding,omitempty"`

    // Minimum margin between boxes
    BoxMargin *int  `yaml:"boxMargin,omitempty"`

    // radius of the box corners in pixel
    CornerRadius *int  `yaml:"cornerRadius,omitempty"`
}


func CopyFormat(src *Format) *Format {
    if src == nil {
        return nil
    }
    var ret Format
    ret.WidthOfParent = src.WidthOfParent
    ret.FixedWidth = src.FixedWidth
    ret.FixedHeight = src.FixedHeight
    ret.VerticalTxt = src.VerticalTxt
    ret.FontCaption = types.CopyFontDef(src.FontCaption)
    ret.FontText1 = types.CopyFontDef(src.FontText1)
    ret.FontText2 = types.CopyFontDef(src.FontText2)
    ret.FontComment = types.CopyFontDef(src.FontComment)
    ret.FontCommentMarker = types.CopyFontDef(src.FontCommentMarker)
    ret.Line = types.CopyLineDef(src.Line)
    ret.Fill = types.CopyFillDef(src.Fill)
    ret.Padding = src.Padding
    ret.BoxMargin = src.BoxMargin
    ret.CornerRadius = src.CornerRadius

    return &ret
}





type FormatVariations struct {

    // dictionary with tag values as key, that contains format definitions.
    HasTag map[string]FormatVariation  `yaml:"hasTag,omitempty"`
}

func NewFormatVariations() *FormatVariations {
    return &FormatVariations{
        HasTag: make(map[string]FormatVariation, 0),
    }
}

func CopyFormatVariations(src *FormatVariations) *FormatVariations {
    if src == nil {
        return nil
    }
    var ret FormatVariations
    ret.HasTag = make(map[string]FormatVariation, 0)
    for k, v := range src.HasTag {
        ret.HasTag[k] = v
    }

    return &ret
}








/* Definition of a topic related overlay ... for instance for heatmaps
*/
type Overlay struct {

    // some catchy words to describe the displayed topc
    Caption string  `yaml:"caption"`

    // Optional reference value that defines the reference value for this type of overlay
    RefValue float64  `yaml:"refValue"`

    // radius for having a value of refValue
    RefRadius *float64  `yaml:"refRadius,omitempty"`

    // in case of multiple overlays existing, this allows to define a percentage offset from the center-x of the related layout object
    CenterXOffset float64  `yaml:"centerXOffset"`

    // in case of multiple overlays existing, this allows to define a percentage offset from the center-y of the related layout object
    CenterYOffset float64  `yaml:"centerYOffset"`

    // if set to true, then the value is printed in the overlay
    PrintValue bool  `yaml:"printValue"`

    // dictionary of layout elements, that contain this overlay. The dictionary stores the value for this specific object
    Layouts map[string]float64  `yaml:"layouts,omitempty"`

    // if this is configured the the radius for the layouts is in a percentage of the refValue
    RadiusDefs *OverlayRadiusDef  `yaml:"radiusDefs,omitempty"`

    Formats *OverlayFormatDef  `yaml:"formats,omitempty"`
}

func NewOverlay() *Overlay {
    return &Overlay{
        Layouts: make(map[string]float64, 0),
        Formats: NewOverlayFormatDef(),
    }
}

func CopyOverlay(src *Overlay) *Overlay {
    if src == nil {
        return nil
    }
    var ret Overlay
    ret.Caption = src.Caption
    ret.RefValue = src.RefValue
    ret.RefRadius = src.RefRadius
    ret.CenterXOffset = src.CenterXOffset
    ret.CenterYOffset = src.CenterYOffset
    ret.PrintValue = src.PrintValue
    ret.Layouts = make(map[string]float64, 0)
    for k, v := range src.Layouts {
        ret.Layouts[k] = v
    }
    ret.RadiusDefs = CopyOverlayRadiusDef(src.RadiusDefs)
    ret.Formats = CopyOverlayFormatDef(src.Formats)

    return &ret
}





/* Definition of one legend entry
*/
type LegendEntry struct {

    Text string  `yaml:"text"`

    // this format reference used to identify how the here described object is in the picture displayed.
    Format string  `yaml:"format"`
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








type Connection struct {

    // Caption text of the destination box, can be used as alternative to 'destId'
    Dest string  `yaml:"dest"`

    // box id of the destination
    DestId string  `yaml:"destId"`

    // additional comment, that can be then included in the created graphic
    Comment *types.Comment  `yaml:"comment,omitempty"`

    // Arrow at the source box
    SourceArrow bool  `yaml:"sourceArrow"`

    // Arrow at the destination box
    DestArrow bool  `yaml:"destArrow"`

    // optional format to style the connection
    Format *string  `yaml:"format,omitempty"`

    // is only set by while the layout is processed, don't set it in the definition
    HiddenComments bool  `yaml:"hiddenComments"`

    // optional container to define additional contrains for the specific connection
    ConnRestrictions *types.ConnRestriction  `yaml:"connRestrictions,omitempty"`

    // Tags to annotate the connection, tags are used to format
    Tags []string  `yaml:"tags,omitempty"`

    // optional step where this comment is part of, is filled via processing not by the user
    Step *int  `yaml:"step,omitempty"`
}

func NewConnection() *Connection {
    return &Connection{
        Tags: make([]string, 0),
    }
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
    ret.Format = src.Format
    ret.HiddenComments = src.HiddenComments
    ret.ConnRestrictions = types.CopyConnRestriction(src.ConnRestrictions)
    ret.Tags = make([]string, 0)
    for _, e := range src.Tags {
        ret.Tags = append(ret.Tags, e)
    }
    ret.Step = src.Step

    return &ret
}





/* container to extend the layouts of a given layout element via mixins
*/
type LayoutMixin struct {

    Horizontal []Layout  `yaml:"horizontal,omitempty"`

    Vertical []Layout  `yaml:"vertical,omitempty"`
}

func NewLayoutMixin() *LayoutMixin {
    return &LayoutMixin{
        Horizontal: make([]Layout, 0),
        Vertical: make([]Layout, 0),
    }
}

func CopyLayoutMixin(src *LayoutMixin) *LayoutMixin {
    if src == nil {
        return nil
    }
    var ret LayoutMixin
    ret.Horizontal = make([]Layout, 0)
    for _, e := range src.Horizontal {
        ret.Horizontal = append(ret.Horizontal, e)
    }
    ret.Vertical = make([]Layout, 0)
    for _, e := range src.Vertical {
        ret.Vertical = append(ret.Vertical, e)
    }

    return &ret
}








type FormatVariation struct {

    // fill this in case the format is specified outside
    FormatRef *string  `yaml:"formatRef,omitempty"`

    Format Format  `yaml:"format"`

    // number to define the order if a layout has for instance multiple matching tags
    Priority int  `yaml:"priority"`
}


func CopyFormatVariation(src *FormatVariation) *FormatVariation {
    if src == nil {
        return nil
    }
    var ret FormatVariation
    ret.FormatRef = src.FormatRef
    ret.Format = *CopyFormat(&src.Format)
    ret.Priority = src.Priority

    return &ret
}








/* definition how to calculate radius changes based on a reference value
*/
type OverlayRadiusDef struct {

    // minimal radius to use for the display
    Min float64  `yaml:"min"`

    // maximal radius to use for the display
    Max float64  `yaml:"max"`
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





/* definition what format to use for a specific reference value
*/
type OverlayFormatDef struct {

    // default format to use to display the overlay
    Default string  `yaml:"default"`

    // grations considered for switching for formats to use
    Gradations []OverlayGradation  `yaml:"gradations,omitempty"`
}

func NewOverlayFormatDef() *OverlayFormatDef {
    return &OverlayFormatDef{
        Gradations: make([]OverlayGradation, 0),
    }
}

func CopyOverlayFormatDef(src *OverlayFormatDef) *OverlayFormatDef {
    if src == nil {
        return nil
    }
    var ret OverlayFormatDef
    ret.Default = src.Default
    ret.Gradations = make([]OverlayGradation, 0)
    for _, e := range src.Gradations {
        ret.Gradations = append(ret.Gradations, e)
    }

    return &ret
}





/* gradation entry for switching formats for overlays
*/
type OverlayGradation struct {

    // to what value should the here named format being used
    Limit float64  `yaml:"limit"`

    // name of the defined format to use
    Format string  `yaml:"format"`
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




