package boxes

// Attention, this file is generated. Manual changes get lost with the next
// run of the code generation.
// created by yacg (template: golang_types.mako v1.1.0)

import (
    "github.com/okieoth/boxes/pkg/types"
)


/* Model to inject additional things in a boxes layout definition
*/
type BoxesFileMixings struct {

    // optional title, that's appended to the original layout title
    Title *string  `yaml:"title,omitempty"`

    // allows to include a version for the layout description
    Version *string  `yaml:"version,omitempty"`

    // Legend definition used in this diagram
    Legend *Legend  `yaml:"legend,omitempty"`

    // dictionary for layout mixins. key of the dictionary is the caption of the box that will take the additional content
    LayoutMixins map[string]LayoutMixin  `yaml:"layoutMixins,omitempty"`

    // dictionary of connection objects
    Connections map[string]ConnectionCont  `yaml:"connections,omitempty"`

    Formats map[string]Format  `yaml:"formats,omitempty"`

    // Set of formats that overwrites the style of boxes, if specific conditions are met
    FormatVariations *FormatVariations  `yaml:"formatVariations,omitempty"`

    // dictionary of comment objects, this comment will applied on layout objects and replace existing comments there
    Comments map[string]types.Comment  `yaml:"comments,omitempty"`

    // dictionary of tag array, the additional tags will be applied on the existing layout and can be used for instance to define display formats
    Tags map[string]Tags  `yaml:"tags,omitempty"`

    // optional map of images used in the generated graphic
    Images map[string]types.ImageDef  `yaml:"images,omitempty"`

    Overlays []Overlay  `yaml:"overlays,omitempty"`

    // additional container to allow step separation in workflows
    Steps []ProcessStep  `yaml:"steps,omitempty"`
}

func NewBoxesFileMixings() *BoxesFileMixings {
    return &BoxesFileMixings{
        Legend: NewLegend(),
        LayoutMixins: make(map[string]LayoutMixin, 0),
        Connections: make(map[string]ConnectionCont, 0),
        Formats: make(map[string]Format, 0),
        FormatVariations: NewFormatVariations(),
        Comments: make(map[string]types.Comment, 0),
        Tags: make(map[string]Tags, 0),
        Images: make(map[string]types.ImageDef, 0),
        Overlays: make([]Overlay, 0),
        Steps: make([]ProcessStep, 0),
    }
}

func CopyBoxesFileMixings(src *BoxesFileMixings) *BoxesFileMixings {
    if src == nil {
        return nil
    }
    var ret BoxesFileMixings
    ret.Title = src.Title
    ret.Version = src.Version
    ret.Legend = CopyLegend(src.Legend)
    ret.LayoutMixins = make(map[string]LayoutMixin, 0)
    for k, v := range src.LayoutMixins {
        ret.LayoutMixins[k] = v
    }
    ret.Connections = make(map[string]ConnectionCont, 0)
    for k, v := range src.Connections {
        ret.Connections[k] = v
    }
    ret.Formats = make(map[string]Format, 0)
    for k, v := range src.Formats {
        ret.Formats[k] = v
    }
    ret.FormatVariations = CopyFormatVariations(src.FormatVariations)
    ret.Comments = make(map[string]types.Comment, 0)
    for k, v := range src.Comments {
        ret.Comments[k] = v
    }
    ret.Tags = make(map[string]Tags, 0)
    for k, v := range src.Tags {
        ret.Tags[k] = v
    }
    ret.Images = make(map[string]types.ImageDef, 0)
    for k, v := range src.Images {
        ret.Images[k] = v
    }
    ret.Overlays = make([]Overlay, 0)
    for _, e := range src.Overlays {
        ret.Overlays = append(ret.Overlays, e)
    }
    ret.Steps = make([]ProcessStep, 0)
    for _, e := range src.Steps {
        ret.Steps = append(ret.Steps, e)
    }

    return &ret
}











type ConnectionCont struct {

    Connections []Connection  `yaml:"connections,omitempty"`
}

func NewConnectionCont() *ConnectionCont {
    return &ConnectionCont{
        Connections: make([]Connection, 0),
    }
}

func CopyConnectionCont(src *ConnectionCont) *ConnectionCont {
    if src == nil {
        return nil
    }
    var ret ConnectionCont
    ret.Connections = make([]Connection, 0)
    for _, e := range src.Connections {
        ret.Connections = append(ret.Connections, e)
    }

    return &ret
}














type Tags struct {

    Tags []string  `yaml:"tags,omitempty"`
}

func NewTags() *Tags {
    return &Tags{
        Tags: make([]string, 0),
    }
}

func CopyTags(src *Tags) *Tags {
    if src == nil {
        return nil
    }
    var ret Tags
    ret.Tags = make([]string, 0)
    for _, e := range src.Tags {
        ret.Tags = append(ret.Tags, e)
    }

    return &ret
}








type ProcessStep struct {

    // title, that's appended to the original layout title
    Caption string  `yaml:"caption"`

    // dictionary of connection objects
    Connections map[string]ConnectionCont  `yaml:"connections,omitempty"`

    // dictionary of comment objects, this comment will applied on layout objects and replace existing comments there
    Comments map[string]types.Comment  `yaml:"comments,omitempty"`

    // dictionary of tag array, the additional tags will be applied on the existing layout and can be used for instance to define display formats
    Tags map[string]Tags  `yaml:"tags,omitempty"`

    Overlays []Overlay  `yaml:"overlays,omitempty"`

    // Set of formats that overwrites the style of boxes, if specific conditions are met
    FormatVariations *FormatVariations  `yaml:"formatVariations,omitempty"`
}

func NewProcessStep() *ProcessStep {
    return &ProcessStep{
        Connections: make(map[string]ConnectionCont, 0),
        Comments: make(map[string]types.Comment, 0),
        Tags: make(map[string]Tags, 0),
        Overlays: make([]Overlay, 0),
        FormatVariations: NewFormatVariations(),
    }
}

func CopyProcessStep(src *ProcessStep) *ProcessStep {
    if src == nil {
        return nil
    }
    var ret ProcessStep
    ret.Caption = src.Caption
    ret.Connections = make(map[string]ConnectionCont, 0)
    for k, v := range src.Connections {
        ret.Connections[k] = v
    }
    ret.Comments = make(map[string]types.Comment, 0)
    for k, v := range src.Comments {
        ret.Comments[k] = v
    }
    ret.Tags = make(map[string]Tags, 0)
    for k, v := range src.Tags {
        ret.Tags[k] = v
    }
    ret.Overlays = make([]Overlay, 0)
    for _, e := range src.Overlays {
        ret.Overlays = append(ret.Overlays, e)
    }
    ret.FormatVariations = CopyFormatVariations(src.FormatVariations)

    return &ret
}













