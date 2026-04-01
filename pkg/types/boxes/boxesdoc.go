package boxes

// Attention, this file is generated. Manual changes get lost with the next
// run of the code generation.
// created by yacg (template: golang_types.mako v1.1.0)

import (
    "github.com/okieoth/boxes/pkg/types"
)


/* Internal model to describe the layout if the boxes chart
*/
type BoxesDocument struct {

    // Title of the document
    Title string  `yaml:"title"`

    // format reference used for the title
    TitleFormat *string  `yaml:"titleFormat,omitempty"`

    // Legend definition used in this diagram
    Legend *Legend  `yaml:"legend,omitempty"`

    Boxes LayoutElement  `yaml:"boxes"`

    // Height of the document
    Height int  `yaml:"height"`

    // Width of the document
    Width int  `yaml:"width"`

    Connections []ConnectionElem  `yaml:"connections,omitempty"`

    ConnectedElems []string  `yaml:"connectedElems,omitempty"`

    // minimal distance between overlapping lines
    LineDist int  `yaml:"lineDist"`

    // Padding used as default over the whole diagram
    GlobalPadding int  `yaml:"globalPadding"`

    // Minimum margin between boxes
    MinBoxMargin int  `yaml:"minBoxMargin"`

    // Minimum margin between connectors
    MinConnectorMargin int  `yaml:"minConnectorMargin"`

    // Map of formats available to be used in the boxes
    Formats map[string]BoxFormat  `yaml:"formats,omitempty"`

    // Map of images used in the generated graphic
    Images map[string]types.ImageDef  `yaml:"images,omitempty"`

    // Vertical roads available to connect boxes in a vertical way
    VerticalRoads []ConnectionLine  `yaml:"verticalRoads,omitempty"`

    // Horizontal roads available to connect boxes in a horizontal way
    HorizontalRoads []ConnectionLine  `yaml:"horizontalRoads,omitempty"`

    // helper structure to take the node points for the possible connection graph
    ConnectionNodes []ConnectionNode  `yaml:"connectionNodes,omitempty"`

    // helper structure for resolving the collisions
    HorizontalLines []ConnectionLine  `yaml:"horizontalLines,omitempty"`

    // helper structure for resolving the collisions
    VerticalLines []ConnectionLine  `yaml:"verticalLines,omitempty"`

    Comments []CommentContainer  `yaml:"comments,omitempty"`

    Overlays []DocOverlay  `yaml:"overlays,omitempty"`

    // is set to true, if there are some comments truncated from the original layout
    HasHiddenComments bool  `yaml:"hasHiddenComments"`

    // hold the radius of the comment markers to use
    CommentMarkerRadius int  `yaml:"commentMarkerRadius"`

    // hold the radius of the comment markers to use
    LegendEndY int  `yaml:"legendEndY"`

    // helper to align tne number of unspecified markers
    CommentCurrent int  `yaml:"commentCurrent"`
}

func NewBoxesDocument() *BoxesDocument {
    return &BoxesDocument{
        Legend: NewLegend(),
        Boxes: *NewLayoutElement(),
        Connections: make([]ConnectionElem, 0),
        ConnectedElems: make([]string, 0),
        Formats: make(map[string]BoxFormat, 0),
        Images: make(map[string]types.ImageDef, 0),
        VerticalRoads: make([]ConnectionLine, 0),
        HorizontalRoads: make([]ConnectionLine, 0),
        ConnectionNodes: make([]ConnectionNode, 0),
        HorizontalLines: make([]ConnectionLine, 0),
        VerticalLines: make([]ConnectionLine, 0),
        Comments: make([]CommentContainer, 0),
        Overlays: make([]DocOverlay, 0),
    }
}

func CopyBoxesDocument(src *BoxesDocument) *BoxesDocument {
    if src == nil {
        return nil
    }
    var ret BoxesDocument
    ret.Title = src.Title
    ret.TitleFormat = src.TitleFormat
    ret.Legend = CopyLegend(src.Legend)
    ret.Boxes = *CopyLayoutElement(&src.Boxes)
    ret.Height = src.Height
    ret.Width = src.Width
    ret.Connections = make([]ConnectionElem, 0)
    for _, e := range src.Connections {
        ret.Connections = append(ret.Connections, e)
    }
    ret.ConnectedElems = make([]string, 0)
    for _, e := range src.ConnectedElems {
        ret.ConnectedElems = append(ret.ConnectedElems, e)
    }
    ret.LineDist = src.LineDist
    ret.GlobalPadding = src.GlobalPadding
    ret.MinBoxMargin = src.MinBoxMargin
    ret.MinConnectorMargin = src.MinConnectorMargin
    ret.Formats = make(map[string]BoxFormat, 0)
    for k, v := range src.Formats {
        ret.Formats[k] = v
    }
    ret.Images = make(map[string]types.ImageDef, 0)
    for k, v := range src.Images {
        ret.Images[k] = v
    }
    ret.VerticalRoads = make([]ConnectionLine, 0)
    for _, e := range src.VerticalRoads {
        ret.VerticalRoads = append(ret.VerticalRoads, e)
    }
    ret.HorizontalRoads = make([]ConnectionLine, 0)
    for _, e := range src.HorizontalRoads {
        ret.HorizontalRoads = append(ret.HorizontalRoads, e)
    }
    ret.ConnectionNodes = make([]ConnectionNode, 0)
    for _, e := range src.ConnectionNodes {
        ret.ConnectionNodes = append(ret.ConnectionNodes, e)
    }
    ret.HorizontalLines = make([]ConnectionLine, 0)
    for _, e := range src.HorizontalLines {
        ret.HorizontalLines = append(ret.HorizontalLines, e)
    }
    ret.VerticalLines = make([]ConnectionLine, 0)
    for _, e := range src.VerticalLines {
        ret.VerticalLines = append(ret.VerticalLines, e)
    }
    ret.Comments = make([]CommentContainer, 0)
    for _, e := range src.Comments {
        ret.Comments = append(ret.Comments, e)
    }
    ret.Overlays = make([]DocOverlay, 0)
    for _, e := range src.Overlays {
        ret.Overlays = append(ret.Overlays, e)
    }
    ret.HasHiddenComments = src.HasHiddenComments
    ret.CommentMarkerRadius = src.CommentMarkerRadius
    ret.LegendEndY = src.LegendEndY
    ret.CommentCurrent = src.CommentCurrent

    return &ret
}





type LayoutElement struct {

    // unique identifier of that entry
    Id string  `yaml:"id"`

    // Some kind of the main text
    Caption string  `yaml:"caption"`

    // First additional text
    Text1 string  `yaml:"text1"`

    // Second additional text
    Text2 string  `yaml:"text2"`

    // additional comments, that can be then included in the created graphic
    Comments []types.Comment  `yaml:"comments,omitempty"`

    // true if in the proccessing some comments where truncated with removed child elemens
    HiddenComments bool  `yaml:"hiddenComments"`

    Image *ImageContainer  `yaml:"image,omitempty"`

    Vertical *LayoutElemContainer  `yaml:"vertical,omitempty"`

    Horizontal *LayoutElemContainer  `yaml:"horizontal,omitempty"`

    // X position of the element
    X int  `yaml:"x"`

    // Y position of the element
    Y int  `yaml:"y"`

    // Width of the element
    Width int  `yaml:"width"`

    // Height of the element
    Height int  `yaml:"height"`

    // X position of the center of the element
    CenterX int  `yaml:"centerX"`

    // X position of the center of the element
    CenterY int  `yaml:"centerY"`

    Format *BoxFormat  `yaml:"format,omitempty"`

    Connections []LayoutElemConnection  `yaml:"connections,omitempty"`

    // Y position of the left side of the element to start the connection
    LeftYToStart *int  `yaml:"leftYToStart,omitempty"`

    // Y position of the right side of the element to start the connection
    RightYToStart *int  `yaml:"rightYToStart,omitempty"`

    // X position of the top side of the element to start the connection
    TopXToStart *int  `yaml:"topXToStart,omitempty"`

    // X position of the bottom side of the element to start the connection
    BottomXToStart *int  `yaml:"bottomXToStart,omitempty"`

    // X position where the text would start
    XTextBox *int  `yaml:"xTextBox,omitempty"`

    // Y position where the text would start
    YTextBox *int  `yaml:"yTextBox,omitempty"`

    // Width of the text area
    WidthTextBox *int  `yaml:"widthTextBox,omitempty"`

    // Height of the text area
    HeightTextBox *int  `yaml:"heightTextBox,omitempty"`

    // if that is set then connections can run through the box, as long as they don't cross the text
    DontBlockConPaths *bool  `yaml:"dontBlockConPaths,omitempty"`

    // Tags to annotate the box, tags are used to format and filter
    Tags []string  `yaml:"tags,omitempty"`

    // Optional link to a source, related to this element. This can be used for instance for on-click handlers in a UI or simply as documentation.
    DataLink *string  `yaml:"dataLink,omitempty"`

    // aggregated connection point restrictions for this element, derived from all connections to and from it
    ConnRestrictions *types.ConnRestrItem  `yaml:"connRestrictions,omitempty"`
}

func NewLayoutElement() *LayoutElement {
    return &LayoutElement{
        Comments: make([]types.Comment, 0),
        Vertical: NewLayoutElemContainer(),
        Horizontal: NewLayoutElemContainer(),
        Connections: make([]LayoutElemConnection, 0),
        Tags: make([]string, 0),
    }
}

func CopyLayoutElement(src *LayoutElement) *LayoutElement {
    if src == nil {
        return nil
    }
    var ret LayoutElement
    ret.Id = src.Id
    ret.Caption = src.Caption
    ret.Text1 = src.Text1
    ret.Text2 = src.Text2
    ret.Comments = make([]types.Comment, 0)
    for _, e := range src.Comments {
        ret.Comments = append(ret.Comments, e)
    }
    ret.HiddenComments = src.HiddenComments
    ret.Image = CopyImageContainer(src.Image)
    ret.Vertical = CopyLayoutElemContainer(src.Vertical)
    ret.Horizontal = CopyLayoutElemContainer(src.Horizontal)
    ret.X = src.X
    ret.Y = src.Y
    ret.Width = src.Width
    ret.Height = src.Height
    ret.CenterX = src.CenterX
    ret.CenterY = src.CenterY
    ret.Format = CopyBoxFormat(src.Format)
    ret.Connections = make([]LayoutElemConnection, 0)
    for _, e := range src.Connections {
        ret.Connections = append(ret.Connections, e)
    }
    ret.LeftYToStart = src.LeftYToStart
    ret.RightYToStart = src.RightYToStart
    ret.TopXToStart = src.TopXToStart
    ret.BottomXToStart = src.BottomXToStart
    ret.XTextBox = src.XTextBox
    ret.YTextBox = src.YTextBox
    ret.WidthTextBox = src.WidthTextBox
    ret.HeightTextBox = src.HeightTextBox
    ret.DontBlockConPaths = src.DontBlockConPaths
    ret.Tags = make([]string, 0)
    for _, e := range src.Tags {
        ret.Tags = append(ret.Tags, e)
    }
    ret.DataLink = src.DataLink
    ret.ConnRestrictions = types.CopyConnRestrItem(src.ConnRestrictions)

    return &ret
}





type ConnectionElem struct {

    // ID of the box where the connector starts
    From *string  `yaml:"from,omitempty"`

    // ID of the box where the connector ends
    To *string  `yaml:"to,omitempty"`

    // Arrow at the source box
    SourceArrow *bool  `yaml:"sourceArrow,omitempty"`

    // Arrow at the destination box
    DestArrow *bool  `yaml:"destArrow,omitempty"`

    Format *types.LineDef  `yaml:"format,omitempty"`

    // is only set by while the layout is processed, don't set it in the definition
    HiddenComments bool  `yaml:"hiddenComments"`

    Parts []ConnectionLine  `yaml:"parts,omitempty"`

    // index of this connection, in the boxes_document object
    ConnectionIndex int  `yaml:"connectionIndex"`

    // optional container to define additional contrains for the specific connection
    ConnRestrictions *types.ConnRestriction  `yaml:"connRestrictions,omitempty"`

    Comment *types.Comment  `yaml:"comment,omitempty"`

    // what process step this connection belongs too
    Step *int  `yaml:"step,omitempty"`
}

func NewConnectionElem() *ConnectionElem {
    return &ConnectionElem{
        Parts: make([]ConnectionLine, 0),
    }
}

func CopyConnectionElem(src *ConnectionElem) *ConnectionElem {
    if src == nil {
        return nil
    }
    var ret ConnectionElem
    ret.From = src.From
    ret.To = src.To
    ret.SourceArrow = src.SourceArrow
    ret.DestArrow = src.DestArrow
    ret.Format = types.CopyLineDef(src.Format)
    ret.HiddenComments = src.HiddenComments
    ret.Parts = make([]ConnectionLine, 0)
    for _, e := range src.Parts {
        ret.Parts = append(ret.Parts, e)
    }
    ret.ConnectionIndex = src.ConnectionIndex
    ret.ConnRestrictions = types.CopyConnRestriction(src.ConnRestrictions)
    ret.Comment = types.CopyComment(src.Comment)
    ret.Step = src.Step

    return &ret
}








type BoxFormat struct {

    // Padding of the box
    Padding int  `yaml:"padding"`

    FontCaption types.FontDef  `yaml:"fontCaption"`

    FontText1 types.FontDef  `yaml:"fontText1"`

    FontText2 types.FontDef  `yaml:"fontText2"`

    FontComment types.FontDef  `yaml:"fontComment"`

    FontCommentMarker types.FontDef  `yaml:"fontCommentMarker"`

    Line *types.LineDef  `yaml:"line,omitempty"`

    // radius of the box corners in pixel
    CornerRadius *int  `yaml:"cornerRadius,omitempty"`

    Fill *types.FillDef  `yaml:"fill,omitempty"`

    // Minimum margin between boxes
    MinBoxMargin int  `yaml:"minBoxMargin"`

    // sets the width of the object to the width of the parent
    WidthOfParent *bool  `yaml:"widthOfParent,omitempty"`

    // optional fixed width that will be applied on the box
    FixedWidth *int  `yaml:"fixedWidth,omitempty"`

    // optional fixed height that will be applied on the box
    FixedHeight *int  `yaml:"fixedHeight,omitempty"`

    // If true, the text will be displayed vertically
    VerticalTxt bool  `yaml:"verticalTxt"`
}


func CopyBoxFormat(src *BoxFormat) *BoxFormat {
    if src == nil {
        return nil
    }
    var ret BoxFormat
    ret.Padding = src.Padding
    ret.FontCaption = *types.CopyFontDef(&src.FontCaption)
    ret.FontText1 = *types.CopyFontDef(&src.FontText1)
    ret.FontText2 = *types.CopyFontDef(&src.FontText2)
    ret.FontComment = *types.CopyFontDef(&src.FontComment)
    ret.FontCommentMarker = *types.CopyFontDef(&src.FontCommentMarker)
    ret.Line = types.CopyLineDef(src.Line)
    ret.CornerRadius = src.CornerRadius
    ret.Fill = types.CopyFillDef(src.Fill)
    ret.MinBoxMargin = src.MinBoxMargin
    ret.WidthOfParent = src.WidthOfParent
    ret.FixedWidth = src.FixedWidth
    ret.FixedHeight = src.FixedHeight
    ret.VerticalTxt = src.VerticalTxt

    return &ret
}








type ConnectionLine struct {

    StartX int  `yaml:"startX"`

    StartY int  `yaml:"startY"`

    EndX int  `yaml:"endX"`

    EndY int  `yaml:"endY"`

    // in case the line is connected to the start layout element, then here is its id
    SrcLayoutId *string  `yaml:"srcLayoutId,omitempty"`

    // in case the line is connected to a end layout element, then here is its id
    DestLayoutId *string  `yaml:"destLayoutId,omitempty"`

    // index of the connection, in the boxes_document object, where this line belongs too
    ConnectionIndex int  `yaml:"connectionIndex"`

    // position of the line in the connections part array
    LineIndex int  `yaml:"lineIndex"`

    Format *types.LineDef  `yaml:"format,omitempty"`

    IsStart bool  `yaml:"isStart"`

    IsEnd bool  `yaml:"isEnd"`

    InverseDirection bool  `yaml:"inverseDirection"`

    // to what step belongs this connection line
    Step *int  `yaml:"step,omitempty"`
}


func CopyConnectionLine(src *ConnectionLine) *ConnectionLine {
    if src == nil {
        return nil
    }
    var ret ConnectionLine
    ret.StartX = src.StartX
    ret.StartY = src.StartY
    ret.EndX = src.EndX
    ret.EndY = src.EndY
    ret.SrcLayoutId = src.SrcLayoutId
    ret.DestLayoutId = src.DestLayoutId
    ret.ConnectionIndex = src.ConnectionIndex
    ret.LineIndex = src.LineIndex
    ret.Format = types.CopyLineDef(src.Format)
    ret.IsStart = src.IsStart
    ret.IsEnd = src.IsEnd
    ret.InverseDirection = src.InverseDirection
    ret.Step = src.Step

    return &ret
}





/* helper type for calculating the connections between elements
*/
type ConnectionNode struct {

    // X position of the element
    X int  `yaml:"x"`

    // Y position of the element
    Y int  `yaml:"y"`

    // in case the edge ends in a layout element, it takes the ID
    NodeId *string  `yaml:"nodeId,omitempty"`

    // optional box id, only on connection nodes that are the entry points to real box connections
    BoxId *string  `yaml:"boxId,omitempty"`

    Edges []ConnectionEdge  `yaml:"edges,omitempty"`
}

func NewConnectionNode() *ConnectionNode {
    return &ConnectionNode{
        Edges: make([]ConnectionEdge, 0),
    }
}

func CopyConnectionNode(src *ConnectionNode) *ConnectionNode {
    if src == nil {
        return nil
    }
    var ret ConnectionNode
    ret.X = src.X
    ret.Y = src.Y
    ret.NodeId = src.NodeId
    ret.BoxId = src.BoxId
    ret.Edges = make([]ConnectionEdge, 0)
    for _, e := range src.Edges {
        ret.Edges = append(ret.Edges, e)
    }

    return &ret
}





/* all parameters to render the comments in the graphic
*/
type CommentContainer struct {

    // text of the comment
    Text string  `yaml:"text"`

    // displayed in the marker for that note
    Label string  `yaml:"label"`

    // format name to use to render this note
    Format CommentFormat  `yaml:"format"`

    // x-coordinate of the marker for that comment
    MarkerX int  `yaml:"markerX"`

    // x-coordinate of the marker for that comment
    MarkerY int  `yaml:"markerY"`

    // calculated width of the marker text
    MarkerTextWidth int  `yaml:"markerTextWidth"`

    // calculated height of the marker text
    MarkerTextHeight int  `yaml:"markerTextHeight"`

    // calculated width of the comment text
    TextWidth int  `yaml:"textWidth"`

    // calculated height of the comment text
    TextHeight int  `yaml:"textHeight"`

    // true if a custom marker is used for this comment
    CustomMarker bool  `yaml:"customMarker"`

    // in case this comment belongs to a connection, is here the connectionId stored
    ConnectionIndex *int  `yaml:"connectionIndex,omitempty"`

    // optional step where this comment is part of, is filled via processing not by the user
    Step *int  `yaml:"step,omitempty"`
}


func CopyCommentContainer(src *CommentContainer) *CommentContainer {
    if src == nil {
        return nil
    }
    var ret CommentContainer
    ret.Text = src.Text
    ret.Label = src.Label
    ret.Format = *CopyCommentFormat(&src.Format)
    ret.MarkerX = src.MarkerX
    ret.MarkerY = src.MarkerY
    ret.MarkerTextWidth = src.MarkerTextWidth
    ret.MarkerTextHeight = src.MarkerTextHeight
    ret.TextWidth = src.TextWidth
    ret.TextHeight = src.TextHeight
    ret.CustomMarker = src.CustomMarker
    ret.ConnectionIndex = src.ConnectionIndex
    ret.Step = src.Step

    return &ret
}





/* Definition of a topic related overlay ... for instance for heatmaps
*/
type DocOverlay struct {

    // some catchy words to describe the displayed topc
    Caption string  `yaml:"caption"`

    // Optional reference value that defines the reference value for this type of overlay
    RefValue float64  `yaml:"refValue"`

    // in case of multiple overlays existing, this allows to define a percentage offset from the center-x of the related layout object
    CenterXOffset float64  `yaml:"centerXOffset"`

    // in case of multiple overlays existing, this allows to define a percentage offset from the center-y of the related layout object
    CenterYOffset float64  `yaml:"centerYOffset"`

    // if set to true, then the value is printed in the overlay
    PrintValue bool  `yaml:"printValue"`

    // dictionary of layout elements, that contain this overlay. The dictionary stores the value for this specific object
    Layouts map[string]OverlayEntry  `yaml:"layouts,omitempty"`

    // if this is configured the the radius for the layouts is in a percentage of the refValue
    RadiusDefs *OverlayRadiusDef  `yaml:"radiusDefs,omitempty"`

    Formats *OverlayFormatDef  `yaml:"formats,omitempty"`
}

func NewDocOverlay() *DocOverlay {
    return &DocOverlay{
        Layouts: make(map[string]OverlayEntry, 0),
        Formats: NewOverlayFormatDef(),
    }
}

func CopyDocOverlay(src *DocOverlay) *DocOverlay {
    if src == nil {
        return nil
    }
    var ret DocOverlay
    ret.Caption = src.Caption
    ret.RefValue = src.RefValue
    ret.CenterXOffset = src.CenterXOffset
    ret.CenterYOffset = src.CenterYOffset
    ret.PrintValue = src.PrintValue
    ret.Layouts = make(map[string]OverlayEntry, 0)
    for k, v := range src.Layouts {
        ret.Layouts[k] = v
    }
    ret.RadiusDefs = CopyOverlayRadiusDef(src.RadiusDefs)
    ret.Formats = CopyOverlayFormatDef(src.Formats)

    return &ret
}





type ImageContainer struct {

    // X position of the element
    X int  `yaml:"x"`

    // Y position of the element
    Y int  `yaml:"y"`

    // Width of the container
    Width int  `yaml:"width"`

    // Height of the container
    Height int  `yaml:"height"`

    // distance top and bottom of the image
    MarginTopBottom int  `yaml:"marginTopBottom"`

    // distance left and right of the image
    MarginLeftRight int  `yaml:"marginLeftRight"`

    // reference to the globally declared image
    ImgId string  `yaml:"imgId"`
}


func CopyImageContainer(src *ImageContainer) *ImageContainer {
    if src == nil {
        return nil
    }
    var ret ImageContainer
    ret.X = src.X
    ret.Y = src.Y
    ret.Width = src.Width
    ret.Height = src.Height
    ret.MarginTopBottom = src.MarginTopBottom
    ret.MarginLeftRight = src.MarginLeftRight
    ret.ImgId = src.ImgId

    return &ret
}





type LayoutElemContainer struct {

    // X position of the element
    X int  `yaml:"x"`

    // Y position of the element
    Y int  `yaml:"y"`

    // Width of the container
    Width int  `yaml:"width"`

    // Height of the container
    Height int  `yaml:"height"`

    Elems []LayoutElement  `yaml:"elems,omitempty"`
}

func NewLayoutElemContainer() *LayoutElemContainer {
    return &LayoutElemContainer{
        Elems: make([]LayoutElement, 0),
    }
}

func CopyLayoutElemContainer(src *LayoutElemContainer) *LayoutElemContainer {
    if src == nil {
        return nil
    }
    var ret LayoutElemContainer
    ret.X = src.X
    ret.Y = src.Y
    ret.Width = src.Width
    ret.Height = src.Height
    ret.Elems = make([]LayoutElement, 0)
    for _, e := range src.Elems {
        ret.Elems = append(ret.Elems, e)
    }

    return &ret
}





type LayoutElemConnection struct {

    // box id of the destination
    DestId string  `yaml:"destId"`

    // additional comment, that can be then included in the created graphic
    Comment *types.Comment  `yaml:"comment,omitempty"`

    // Arrow at the source box
    SourceArrow bool  `yaml:"sourceArrow"`

    // Arrow at the destination box
    DestArrow bool  `yaml:"destArrow"`

    Format *types.LineDef  `yaml:"format,omitempty"`

    // is only set by while the layout is processed, don't set it in the definition
    HiddenComments bool  `yaml:"hiddenComments"`

    // optional container to define additional contrains for the specific connection
    ConnRestrictions *types.ConnRestriction  `yaml:"connRestrictions,omitempty"`

    // Tags to annotate the connection, tags are used to format
    Tags []string  `yaml:"tags,omitempty"`

    // optional step where this comment is part of, is filled via processing not by the user
    Step *int  `yaml:"step,omitempty"`
}

func NewLayoutElemConnection() *LayoutElemConnection {
    return &LayoutElemConnection{
        Tags: make([]string, 0),
    }
}

func CopyLayoutElemConnection(src *LayoutElemConnection) *LayoutElemConnection {
    if src == nil {
        return nil
    }
    var ret LayoutElemConnection
    ret.DestId = src.DestId
    ret.Comment = types.CopyComment(src.Comment)
    ret.SourceArrow = src.SourceArrow
    ret.DestArrow = src.DestArrow
    ret.Format = types.CopyLineDef(src.Format)
    ret.HiddenComments = src.HiddenComments
    ret.ConnRestrictions = types.CopyConnRestriction(src.ConnRestrictions)
    ret.Tags = make([]string, 0)
    for _, e := range src.Tags {
        ret.Tags = append(ret.Tags, e)
    }
    ret.Step = src.Step

    return &ret
}





type CommentFormat struct {

    FontText types.FontDef  `yaml:"fontText"`

    FontMarker types.FontDef  `yaml:"fontMarker"`

    Line types.LineDef  `yaml:"line"`

    Fill types.FillDef  `yaml:"fill"`
}


func CopyCommentFormat(src *CommentFormat) *CommentFormat {
    if src == nil {
        return nil
    }
    var ret CommentFormat
    ret.FontText = *types.CopyFontDef(&src.FontText)
    ret.FontMarker = *types.CopyFontDef(&src.FontMarker)
    ret.Line = *types.CopyLineDef(&src.Line)
    ret.Fill = *types.CopyFillDef(&src.Fill)

    return &ret
}





/* edge type to store edges on a connection node
*/
type ConnectionEdge struct {

    // X position of the element
    X int  `yaml:"x"`

    // Y position of the element
    Y int  `yaml:"y"`

    // either the box ID where the edge ends or the ID connection node
    DestNodeId *string  `yaml:"destNodeId,omitempty"`

    // weight of the connection, based on the distance
    Weight int  `yaml:"weight"`
}


func CopyConnectionEdge(src *ConnectionEdge) *ConnectionEdge {
    if src == nil {
        return nil
    }
    var ret ConnectionEdge
    ret.X = src.X
    ret.Y = src.Y
    ret.DestNodeId = src.DestNodeId
    ret.Weight = src.Weight

    return &ret
}








type OverlayEntry struct {

    Value float64  `yaml:"value"`

    Radius float64  `yaml:"radius"`

    X int  `yaml:"x"`

    Y int  `yaml:"y"`

    Format BoxFormat  `yaml:"format"`
}


func CopyOverlayEntry(src *OverlayEntry) *OverlayEntry {
    if src == nil {
        return nil
    }
    var ret OverlayEntry
    ret.Value = src.Value
    ret.Radius = src.Radius
    ret.X = src.X
    ret.Y = src.Y
    ret.Format = *CopyBoxFormat(&src.Format)

    return &ret
}




