// Attention, this code is generated. Do not modify manually. Changes will
// be overwritten be the next codegen run.
// Generated from BoxesDocument (configs/models/boxes_document.json)



// Internal model to describe the layout if the boxes chart

package boxes
import (
    "github.com/okieoth/boxes/pkg/types"
)




type LayoutElement struct {
    // unique identifier of that entry
Id string `yaml:"id"`
    // Some kind of the main text
Caption string `yaml:"caption"`
    // First additional text
Text1 string `yaml:"text1"`
    // Second additional text
Text2 string `yaml:"text2"`
    // additional comments, that can be then included in the created graphic
Comments []types.Comment `yaml:"comments,omitempty"`
    // true if in the proccessing some comments where truncated with removed child elemens
HiddenComments bool `yaml:"hiddenComments"`
Image *ImageContainer `yaml:"image,omitempty"`
Vertical *LayoutElemContainer `yaml:"vertical,omitempty"`
Horizontal *LayoutElemContainer `yaml:"horizontal,omitempty"`
    // X position of the element
X int `yaml:"x"`
    // Y position of the element
Y int `yaml:"y"`
    // Width of the element
Width int `yaml:"width"`
    // Height of the element
Height int `yaml:"height"`
    // X position of the center of the element
CenterX int `yaml:"centerX"`
    // X position of the center of the element
CenterY int `yaml:"centerY"`
Format *BoxFormat `yaml:"format,omitempty"`
Connections []LayoutElemConnection `yaml:"connections,omitempty"`
    // Y position of the left side of the element to start the connection
LeftYToStart *int `yaml:"leftYToStart,omitempty"`
    // Y position of the right side of the element to start the connection
RightYToStart *int `yaml:"rightYToStart,omitempty"`
    // X position of the top side of the element to start the connection
TopXToStart *int `yaml:"topXToStart,omitempty"`
    // X position of the bottom side of the element to start the connection
BottomXToStart *int `yaml:"bottomXToStart,omitempty"`
    // X position where the text would start
XTextBox *int `yaml:"xTextBox,omitempty"`
    // Y position where the text would start
YTextBox *int `yaml:"yTextBox,omitempty"`
    // Width of the text area
WidthTextBox *int `yaml:"widthTextBox,omitempty"`
    // Height of the text area
HeightTextBox *int `yaml:"heightTextBox,omitempty"`
    // if that is set then connections can run through the box, as long as they don't cross the text
DontBlockConPaths *bool `yaml:"dontBlockConPaths,omitempty"`
    // Tags to annotate the box, tags are used to format and filter
Tags []string `yaml:"tags,omitempty"`
    // Optional link to a source, related to this element. This can be used for instance for on-click handlers in a UI or simply as documentation.
DataLink *string `yaml:"dataLink,omitempty"`
    // aggregated connection point restrictions for this element, derived from all connections to and from it
ConnRestrictions *types.ConnRestrItem `yaml:"connRestrictions,omitempty"`
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

    if src.Comments != nil {
        ret.Comments = make([]types.Comment, len(src.Comments))
        for i, v := range src.Comments {
            ret.Comments[i] = *types.CopyComment(&v)
        }
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

    if src.Connections != nil {
        ret.Connections = make([]LayoutElemConnection, len(src.Connections))
        for i, v := range src.Connections {
            ret.Connections[i] = *CopyLayoutElemConnection(&v)
        }
    }

    if src.LeftYToStart != nil {
        v := *src.LeftYToStart
        ret.LeftYToStart = &v
    }

    if src.RightYToStart != nil {
        v := *src.RightYToStart
        ret.RightYToStart = &v
    }

    if src.TopXToStart != nil {
        v := *src.TopXToStart
        ret.TopXToStart = &v
    }

    if src.BottomXToStart != nil {
        v := *src.BottomXToStart
        ret.BottomXToStart = &v
    }

    if src.XTextBox != nil {
        v := *src.XTextBox
        ret.XTextBox = &v
    }

    if src.YTextBox != nil {
        v := *src.YTextBox
        ret.YTextBox = &v
    }

    if src.WidthTextBox != nil {
        v := *src.WidthTextBox
        ret.WidthTextBox = &v
    }

    if src.HeightTextBox != nil {
        v := *src.HeightTextBox
        ret.HeightTextBox = &v
    }

    if src.DontBlockConPaths != nil {
        v := *src.DontBlockConPaths
        ret.DontBlockConPaths = &v
    }

    if src.Tags != nil {
        ret.Tags = make([]string, len(src.Tags))
        for i, v := range src.Tags {
            ret.Tags[i] = v
        }
    }

    if src.DataLink != nil {
        v := *src.DataLink
        ret.DataLink = &v
    }

    ret.ConnRestrictions = types.CopyConnRestrItem(src.ConnRestrictions)
return &ret
}


func NewLayoutElement() *LayoutElement {
    var ret LayoutElement
    ret.Comments = make([]types.Comment, 0)
    ret.Connections = make([]LayoutElemConnection, 0)
    ret.Tags = make([]string, 0)
    return &ret
}

type LayoutElemContainer struct {
    // X position of the element
X int `yaml:"x"`
    // Y position of the element
Y int `yaml:"y"`
    // Width of the container
Width int `yaml:"width"`
    // Height of the container
Height int `yaml:"height"`
Elems []LayoutElement `yaml:"elems"`
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

    if src.Elems != nil {
        ret.Elems = make([]LayoutElement, len(src.Elems))
        for i, v := range src.Elems {
            ret.Elems[i] = *CopyLayoutElement(&v)
        }
    }
return &ret
}


func NewLayoutElemContainer() *LayoutElemContainer {
    var ret LayoutElemContainer
    ret.Elems = make([]LayoutElement, 0)
    return &ret
}

type LayoutElemConnection struct {
    // box id of the destination
DestId string `yaml:"destId"`
    // additional comment, that can be then included in the created graphic
Comment *types.Comment `yaml:"comment,omitempty"`
    // Arrow at the source box
SourceArrow bool `yaml:"sourceArrow"`
    // Arrow at the destination box
DestArrow bool `yaml:"destArrow"`
Format *types.LineDef `yaml:"format,omitempty"`
    // is only set by while the layout is processed, don't set it in the definition
HiddenComments bool `yaml:"hiddenComments"`
    // optional container to define additional contrains for the specific connection
ConnRestrictions *types.ConnRestriction `yaml:"connRestrictions,omitempty"`
    // Tags to annotate the connection, tags are used to format
Tags []string `yaml:"tags"`
    // optional step where this comment is part of, is filled via processing not by the user
Step *int `yaml:"step,omitempty"`
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


func NewLayoutElemConnection() *LayoutElemConnection {
    var ret LayoutElemConnection
    ret.Tags = make([]string, 0)
    return &ret
}

type ConnectionElem struct {
    // ID of the box where the connector starts
From *string `yaml:"from,omitempty"`
    // ID of the box where the connector ends
To *string `yaml:"to,omitempty"`
    // Arrow at the source box
SourceArrow *bool `yaml:"sourceArrow,omitempty"`
    // Arrow at the destination box
DestArrow *bool `yaml:"destArrow,omitempty"`
Format *types.LineDef `yaml:"format,omitempty"`
    // is only set by while the layout is processed, don't set it in the definition
HiddenComments bool `yaml:"hiddenComments"`
Parts []ConnectionLine `yaml:"parts"`
    // index of this connection, in the boxes_document object
ConnectionIndex int `yaml:"connectionIndex"`
    // optional container to define additional contrains for the specific connection
ConnRestrictions *types.ConnRestriction `yaml:"connRestrictions,omitempty"`
Comment *types.Comment `yaml:"comment,omitempty"`
    // what process step this connection belongs too
Step *int `yaml:"step,omitempty"`
}


func CopyConnectionElem(src *ConnectionElem) *ConnectionElem {
    if src == nil {
        return nil
    }
    var ret ConnectionElem

    if src.From != nil {
        v := *src.From
        ret.From = &v
    }

    if src.To != nil {
        v := *src.To
        ret.To = &v
    }

    if src.SourceArrow != nil {
        v := *src.SourceArrow
        ret.SourceArrow = &v
    }

    if src.DestArrow != nil {
        v := *src.DestArrow
        ret.DestArrow = &v
    }

    ret.Format = types.CopyLineDef(src.Format)

    ret.HiddenComments = src.HiddenComments

    if src.Parts != nil {
        ret.Parts = make([]ConnectionLine, len(src.Parts))
        for i, v := range src.Parts {
            ret.Parts[i] = *CopyConnectionLine(&v)
        }
    }

    ret.ConnectionIndex = src.ConnectionIndex

    ret.ConnRestrictions = types.CopyConnRestriction(src.ConnRestrictions)

    ret.Comment = types.CopyComment(src.Comment)

    if src.Step != nil {
        v := *src.Step
        ret.Step = &v
    }
return &ret
}


func NewConnectionElem() *ConnectionElem {
    var ret ConnectionElem
    ret.Parts = make([]ConnectionLine, 0)
    return &ret
}

type ConnectionLine struct {
StartX int `yaml:"startX"`
StartY int `yaml:"startY"`
EndX int `yaml:"endX"`
EndY int `yaml:"endY"`
    // in case the line is connected to the start layout element, then here is its id
SrcLayoutId *string `yaml:"srcLayoutId,omitempty"`
    // in case the line is connected to a end layout element, then here is its id
DestLayoutId *string `yaml:"destLayoutId,omitempty"`
    // index of the connection, in the boxes_document object, where this line belongs too
ConnectionIndex int `yaml:"connectionIndex"`
    // position of the line in the connections part array
LineIndex int `yaml:"lineIndex"`
Format *types.LineDef `yaml:"format,omitempty"`
IsStart bool `yaml:"isStart"`
IsEnd bool `yaml:"isEnd"`
InverseDirection bool `yaml:"inverseDirection"`
    // to what step belongs this connection line
Step *int `yaml:"step,omitempty"`
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

    if src.SrcLayoutId != nil {
        v := *src.SrcLayoutId
        ret.SrcLayoutId = &v
    }

    if src.DestLayoutId != nil {
        v := *src.DestLayoutId
        ret.DestLayoutId = &v
    }

    ret.ConnectionIndex = src.ConnectionIndex

    ret.LineIndex = src.LineIndex

    ret.Format = types.CopyLineDef(src.Format)

    ret.IsStart = src.IsStart

    ret.IsEnd = src.IsEnd

    ret.InverseDirection = src.InverseDirection

    if src.Step != nil {
        v := *src.Step
        ret.Step = &v
    }
return &ret
}


func NewConnectionLine() *ConnectionLine {
    var ret ConnectionLine
    return &ret
}

type ImageContainer struct {
    // X position of the element
X int `yaml:"x"`
    // Y position of the element
Y int `yaml:"y"`
    // Width of the container
Width int `yaml:"width"`
    // Height of the container
Height int `yaml:"height"`
    // distance top and bottom of the image
MarginTopBottom int `yaml:"marginTopBottom"`
    // distance left and right of the image
MarginLeftRight int `yaml:"marginLeftRight"`
    // reference to the globally declared image
ImgId string `yaml:"imgId"`
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


func NewImageContainer() *ImageContainer {
    var ret ImageContainer
    return &ret
}

type BoxFormat struct {
    // Padding of the box
Padding int `yaml:"padding"`
FontCaption types.FontDef `yaml:"fontCaption"`
FontText1 types.FontDef `yaml:"fontText1"`
FontText2 types.FontDef `yaml:"fontText2"`
FontComment types.FontDef `yaml:"fontComment"`
FontCommentMarker types.FontDef `yaml:"fontCommentMarker"`
Line *types.LineDef `yaml:"line,omitempty"`
    // radius of the box corners in pixel
CornerRadius *int `yaml:"cornerRadius,omitempty"`
Fill *types.FillDef `yaml:"fill,omitempty"`
    // Minimum margin between boxes
MinBoxMargin int `yaml:"minBoxMargin"`
    // sets the width of the object to the width of the parent
WidthOfParent *bool `yaml:"widthOfParent,omitempty"`
    // optional fixed width that will be applied on the box
FixedWidth *int `yaml:"fixedWidth,omitempty"`
    // optional fixed height that will be applied on the box
FixedHeight *int `yaml:"fixedHeight,omitempty"`
    // If true, the text will be displayed vertically
VerticalTxt bool `yaml:"verticalTxt"`
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

    if src.CornerRadius != nil {
        v := *src.CornerRadius
        ret.CornerRadius = &v
    }

    ret.Fill = types.CopyFillDef(src.Fill)

    ret.MinBoxMargin = src.MinBoxMargin

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

    ret.VerticalTxt = src.VerticalTxt
return &ret
}


func NewBoxFormat() *BoxFormat {
    var ret BoxFormat
    return &ret
}

type CommentFormat struct {
FontText types.FontDef `yaml:"fontText"`
FontMarker types.FontDef `yaml:"fontMarker"`
Line types.LineDef `yaml:"line"`
Fill types.FillDef `yaml:"fill"`
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


func NewCommentFormat() *CommentFormat {
    var ret CommentFormat
    return &ret
}

// helper type for calculating the connections between elements
type ConnectionNode struct {
    // X position of the element
X int `yaml:"x"`
    // Y position of the element
Y int `yaml:"y"`
    // in case the edge ends in a layout element, it takes the ID
NodeId *string `yaml:"nodeId,omitempty"`
    // optional box id, only on connection nodes that are the entry points to real box connections
BoxId *string `yaml:"boxId,omitempty"`
Edges []ConnectionEdge `yaml:"edges,omitempty"`
}


func CopyConnectionNode(src *ConnectionNode) *ConnectionNode {
    if src == nil {
        return nil
    }
    var ret ConnectionNode

    ret.X = src.X

    ret.Y = src.Y

    if src.NodeId != nil {
        v := *src.NodeId
        ret.NodeId = &v
    }

    if src.BoxId != nil {
        v := *src.BoxId
        ret.BoxId = &v
    }

    if src.Edges != nil {
        ret.Edges = make([]ConnectionEdge, len(src.Edges))
        for i, v := range src.Edges {
            ret.Edges[i] = *CopyConnectionEdge(&v)
        }
    }
return &ret
}


func NewConnectionNode() *ConnectionNode {
    var ret ConnectionNode
    ret.Edges = make([]ConnectionEdge, 0)
    return &ret
}

// edge type to store edges on a connection node
type ConnectionEdge struct {
    // X position of the element
X int `yaml:"x"`
    // Y position of the element
Y int `yaml:"y"`
    // either the box ID where the edge ends or the ID connection node
DestNodeId *string `yaml:"destNodeId,omitempty"`
    // weight of the connection, based on the distance
Weight int `yaml:"weight"`
}


func CopyConnectionEdge(src *ConnectionEdge) *ConnectionEdge {
    if src == nil {
        return nil
    }
    var ret ConnectionEdge

    ret.X = src.X

    ret.Y = src.Y

    if src.DestNodeId != nil {
        v := *src.DestNodeId
        ret.DestNodeId = &v
    }

    ret.Weight = src.Weight
return &ret
}


func NewConnectionEdge() *ConnectionEdge {
    var ret ConnectionEdge
    return &ret
}

// all parameters to render the comments in the graphic
type CommentContainer struct {
    // text of the comment
Text string `yaml:"text"`
    // displayed in the marker for that note
Label string `yaml:"label"`
    // format name to use to render this note
Format CommentFormat `yaml:"format"`
    // x-coordinate of the marker for that comment
MarkerX int `yaml:"markerX"`
    // x-coordinate of the marker for that comment
MarkerY int `yaml:"markerY"`
    // calculated width of the marker text
MarkerTextWidth int `yaml:"markerTextWidth"`
    // calculated height of the marker text
MarkerTextHeight int `yaml:"markerTextHeight"`
    // calculated width of the comment text
TextWidth int `yaml:"textWidth"`
    // calculated height of the comment text
TextHeight int `yaml:"textHeight"`
    // true if a custom marker is used for this comment
CustomMarker bool `yaml:"customMarker"`
    // in case this comment belongs to a connection, is here the connectionId stored
ConnectionIndex *int `yaml:"connectionIndex,omitempty"`
    // optional step where this comment is part of, is filled via processing not by the user
Step *int `yaml:"step,omitempty"`
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

    if src.ConnectionIndex != nil {
        v := *src.ConnectionIndex
        ret.ConnectionIndex = &v
    }

    if src.Step != nil {
        v := *src.Step
        ret.Step = &v
    }
return &ret
}


func NewCommentContainer() *CommentContainer {
    var ret CommentContainer
    return &ret
}

// Definition of a topic related overlay ... for instance for heatmaps
type DocOverlay struct {
    // some catchy words to describe the displayed topc
Caption string `yaml:"caption"`
    // Optional reference value that defines the reference value for this type of overlay
RefValue float64 `yaml:"refValue"`
    // in case of multiple overlays existing, this allows to define a percentage offset from the center-x of the related layout object
CenterXOffset float64 `yaml:"centerXOffset"`
    // in case of multiple overlays existing, this allows to define a percentage offset from the center-y of the related layout object
CenterYOffset float64 `yaml:"centerYOffset"`
    // if set to true, then the value is printed in the overlay
PrintValue bool `yaml:"printValue"`
    // dictionary of layout elements, that contain this overlay. The dictionary stores the value for this specific object
Layouts map[string]OverlayEntry `yaml:"layouts"`
    // if this is configured the the radius for the layouts is in a percentage of the refValue
RadiusDefs *OverlayRadiusDef `yaml:"radiusDefs,omitempty"`
Formats *OverlayFormatDef `yaml:"formats,omitempty"`
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

    if src.Layouts != nil {
        ret.Layouts = make(map[string]OverlayEntry, len(src.Layouts))
        for k, v := range src.Layouts {
            ret.Layouts[k] = *CopyOverlayEntry(&v)
        }
    }

    ret.RadiusDefs = CopyOverlayRadiusDef(src.RadiusDefs)

    ret.Formats = CopyOverlayFormatDef(src.Formats)
return &ret
}


func NewDocOverlay() *DocOverlay {
    var ret DocOverlay
    ret.Layouts = make(map[string]OverlayEntry, 0)
    return &ret
}

type OverlayEntry struct {
Value float64 `yaml:"value"`
Radius float64 `yaml:"radius"`
X int `yaml:"x"`
Y int `yaml:"y"`
Format BoxFormat `yaml:"format"`
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


func NewOverlayEntry() *OverlayEntry {
    var ret OverlayEntry
    return &ret
}

// Internal model to describe the layout if the boxes chart
type BoxesDocument struct {
    // Title of the document
Title string `yaml:"title"`
    // format reference used for the title
TitleFormat *string `yaml:"titleFormat,omitempty"`
    // Legend definition used in this diagram
Legend *Legend `yaml:"legend,omitempty"`
Boxes LayoutElement `yaml:"boxes"`
    // Height of the document
Height int `yaml:"height"`
    // Width of the document
Width int `yaml:"width"`
Connections []ConnectionElem `yaml:"connections,omitempty"`
ConnectedElems []string `yaml:"connectedElems,omitempty"`
    // minimal distance between overlapping lines
LineDist int `yaml:"lineDist"`
    // Padding used as default over the whole diagram
GlobalPadding int `yaml:"globalPadding"`
    // Minimum margin between boxes
MinBoxMargin int `yaml:"minBoxMargin"`
    // Minimum margin between connectors
MinConnectorMargin int `yaml:"minConnectorMargin"`
    // Map of formats available to be used in the boxes
Formats map[string]BoxFormat `yaml:"formats,omitempty"`
    // Map of images used in the generated graphic
Images map[string]types.ImageDef `yaml:"images,omitempty"`
    // Vertical roads available to connect boxes in a vertical way
VerticalRoads []ConnectionLine `yaml:"verticalRoads,omitempty"`
    // Horizontal roads available to connect boxes in a horizontal way
HorizontalRoads []ConnectionLine `yaml:"horizontalRoads,omitempty"`
    // helper structure to take the node points for the possible connection graph
ConnectionNodes []ConnectionNode `yaml:"connectionNodes,omitempty"`
    // helper structure for resolving the collisions
HorizontalLines []ConnectionLine `yaml:"horizontalLines,omitempty"`
    // helper structure for resolving the collisions
VerticalLines []ConnectionLine `yaml:"verticalLines,omitempty"`
Comments []CommentContainer `yaml:"comments,omitempty"`
Overlays []DocOverlay `yaml:"overlays,omitempty"`
    // is set to true, if there are some comments truncated from the original layout
HasHiddenComments bool `yaml:"hasHiddenComments"`
    // hold the radius of the comment markers to use
CommentMarkerRadius int `yaml:"commentMarkerRadius"`
    // hold the radius of the comment markers to use
LegendEndY int `yaml:"legendEndY"`
    // helper to align tne number of unspecified markers
CommentCurrent int `yaml:"commentCurrent"`
}


func CopyBoxesDocument(src *BoxesDocument) *BoxesDocument {
    if src == nil {
        return nil
    }
    var ret BoxesDocument

    ret.Title = src.Title

    if src.TitleFormat != nil {
        v := *src.TitleFormat
        ret.TitleFormat = &v
    }

    ret.Legend = CopyLegend(src.Legend)

    ret.Boxes = *CopyLayoutElement(&src.Boxes)

    ret.Height = src.Height

    ret.Width = src.Width

    if src.Connections != nil {
        ret.Connections = make([]ConnectionElem, len(src.Connections))
        for i, v := range src.Connections {
            ret.Connections[i] = *CopyConnectionElem(&v)
        }
    }

    if src.ConnectedElems != nil {
        ret.ConnectedElems = make([]string, len(src.ConnectedElems))
        for i, v := range src.ConnectedElems {
            ret.ConnectedElems[i] = v
        }
    }

    ret.LineDist = src.LineDist

    ret.GlobalPadding = src.GlobalPadding

    ret.MinBoxMargin = src.MinBoxMargin

    ret.MinConnectorMargin = src.MinConnectorMargin

    if src.Formats != nil {
        ret.Formats = make(map[string]BoxFormat, len(src.Formats))
        for k, v := range src.Formats {
            ret.Formats[k] = *CopyBoxFormat(&v)
        }
    }

    if src.Images != nil {
        ret.Images = make(map[string]types.ImageDef, len(src.Images))
        for k, v := range src.Images {
            ret.Images[k] = *types.CopyImageDef(&v)
        }
    }

    if src.VerticalRoads != nil {
        ret.VerticalRoads = make([]ConnectionLine, len(src.VerticalRoads))
        for i, v := range src.VerticalRoads {
            ret.VerticalRoads[i] = *CopyConnectionLine(&v)
        }
    }

    if src.HorizontalRoads != nil {
        ret.HorizontalRoads = make([]ConnectionLine, len(src.HorizontalRoads))
        for i, v := range src.HorizontalRoads {
            ret.HorizontalRoads[i] = *CopyConnectionLine(&v)
        }
    }

    if src.ConnectionNodes != nil {
        ret.ConnectionNodes = make([]ConnectionNode, len(src.ConnectionNodes))
        for i, v := range src.ConnectionNodes {
            ret.ConnectionNodes[i] = *CopyConnectionNode(&v)
        }
    }

    if src.HorizontalLines != nil {
        ret.HorizontalLines = make([]ConnectionLine, len(src.HorizontalLines))
        for i, v := range src.HorizontalLines {
            ret.HorizontalLines[i] = *CopyConnectionLine(&v)
        }
    }

    if src.VerticalLines != nil {
        ret.VerticalLines = make([]ConnectionLine, len(src.VerticalLines))
        for i, v := range src.VerticalLines {
            ret.VerticalLines[i] = *CopyConnectionLine(&v)
        }
    }

    if src.Comments != nil {
        ret.Comments = make([]CommentContainer, len(src.Comments))
        for i, v := range src.Comments {
            ret.Comments[i] = *CopyCommentContainer(&v)
        }
    }

    if src.Overlays != nil {
        ret.Overlays = make([]DocOverlay, len(src.Overlays))
        for i, v := range src.Overlays {
            ret.Overlays[i] = *CopyDocOverlay(&v)
        }
    }

    ret.HasHiddenComments = src.HasHiddenComments

    ret.CommentMarkerRadius = src.CommentMarkerRadius

    ret.LegendEndY = src.LegendEndY

    ret.CommentCurrent = src.CommentCurrent
return &ret
}


func NewBoxesDocument() *BoxesDocument {
    var ret BoxesDocument
    ret.Connections = make([]ConnectionElem, 0)
    ret.ConnectedElems = make([]string, 0)
    ret.Formats = make(map[string]BoxFormat, 0)
    ret.Images = make(map[string]types.ImageDef, 0)
    ret.VerticalRoads = make([]ConnectionLine, 0)
    ret.HorizontalRoads = make([]ConnectionLine, 0)
    ret.ConnectionNodes = make([]ConnectionNode, 0)
    ret.HorizontalLines = make([]ConnectionLine, 0)
    ret.VerticalLines = make([]ConnectionLine, 0)
    ret.Comments = make([]CommentContainer, 0)
    ret.Overlays = make([]DocOverlay, 0)
    return &ret
}
