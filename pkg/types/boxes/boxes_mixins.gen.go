// Attention, this code is generated. Do not modify manually. Changes will
// be overwritten be the next codegen run.
// Generated from Boxes file mixings v0.4.0 (configs/models/boxes_external_mixins.json)



// Model to inject additional things in a boxes layout definition

package boxes
import (
    "github.com/okieoth/boxes/pkg/types"
)




type ProcessStep struct {
    // title, that's appended to the original layout title
Caption string `yaml:"caption"`
    // dictionary of connection objects
Connections map[string]ConnectionCont `yaml:"connections,omitempty"`
    // dictionary of comment objects, this comment will applied on layout objects and replace existing comments there
Comments map[string]types.Comment `yaml:"comments,omitempty"`
    // dictionary of tag array, the additional tags will be applied on the existing layout and can be used for instance to define display formats
Tags map[string]Tags `yaml:"tags,omitempty"`
Overlays []Overlay `yaml:"overlays,omitempty"`
    // Set of formats that overwrites the style of boxes, if specific conditions are met
FormatVariations *FormatVariations `yaml:"formatVariations,omitempty"`
}


func CopyProcessStep(src *ProcessStep) *ProcessStep {
    if src == nil {
        return nil
    }
    var ret ProcessStep

    ret.Caption = src.Caption

    if src.Connections != nil {
        ret.Connections = make(map[string]ConnectionCont, len(src.Connections))
        for k, v := range src.Connections {
            ret.Connections[k] = *CopyConnectionCont(&v)
        }
    }

    if src.Comments != nil {
        ret.Comments = make(map[string]types.Comment, len(src.Comments))
        for k, v := range src.Comments {
            ret.Comments[k] = *types.CopyComment(&v)
        }
    }

    if src.Tags != nil {
        ret.Tags = make(map[string]Tags, len(src.Tags))
        for k, v := range src.Tags {
            ret.Tags[k] = *CopyTags(&v)
        }
    }

    if src.Overlays != nil {
        ret.Overlays = make([]Overlay, len(src.Overlays))
        for i, v := range src.Overlays {
            ret.Overlays[i] = *CopyOverlay(&v)
        }
    }

    ret.FormatVariations = CopyFormatVariations(src.FormatVariations)
return &ret
}


func NewProcessStep() *ProcessStep {
    var ret ProcessStep
    ret.Connections = make(map[string]ConnectionCont, 0)
    ret.Comments = make(map[string]types.Comment, 0)
    ret.Tags = make(map[string]Tags, 0)
    ret.Overlays = make([]Overlay, 0)
    return &ret
}

type ConnectionCont struct {
Connections []Connection `yaml:"connections,omitempty"`
}


func CopyConnectionCont(src *ConnectionCont) *ConnectionCont {
    if src == nil {
        return nil
    }
    var ret ConnectionCont

    if src.Connections != nil {
        ret.Connections = make([]Connection, len(src.Connections))
        for i, v := range src.Connections {
            ret.Connections[i] = *CopyConnection(&v)
        }
    }
return &ret
}


func NewConnectionCont() *ConnectionCont {
    var ret ConnectionCont
    ret.Connections = make([]Connection, 0)
    return &ret
}

type Tags struct {
Tags []string `yaml:"tags,omitempty"`
}


func CopyTags(src *Tags) *Tags {
    if src == nil {
        return nil
    }
    var ret Tags

    if src.Tags != nil {
        ret.Tags = make([]string, len(src.Tags))
        for i, v := range src.Tags {
            ret.Tags[i] = v
        }
    }
return &ret
}


func NewTags() *Tags {
    var ret Tags
    ret.Tags = make([]string, 0)
    return &ret
}

// Model to inject additional things in a boxes layout definition
type BoxesFileMixings struct {
    // optional title, that's appended to the original layout title
Title *string `yaml:"title,omitempty"`
    // allows to include a version for the layout description
Version *string `yaml:"version,omitempty"`
    // Legend definition used in this diagram
Legend *Legend `yaml:"legend,omitempty"`
    // dictionary for layout mixins. key of the dictionary is the caption of the box that will take the additional content
LayoutMixins map[string]LayoutMixin `yaml:"layoutMixins,omitempty"`
    // dictionary of connection objects
Connections map[string]ConnectionCont `yaml:"connections,omitempty"`
Formats map[string]Format `yaml:"formats,omitempty"`
    // Set of formats that overwrites the style of boxes, if specific conditions are met
FormatVariations *FormatVariations `yaml:"formatVariations,omitempty"`
    // dictionary of comment objects, this comment will applied on layout objects and replace existing comments there
Comments map[string]types.Comment `yaml:"comments,omitempty"`
    // dictionary of tag array, the additional tags will be applied on the existing layout and can be used for instance to define display formats
Tags map[string]Tags `yaml:"tags,omitempty"`
    // optional map of images used in the generated graphic
Images map[string]types.ImageDef `yaml:"images,omitempty"`
Overlays []Overlay `yaml:"overlays,omitempty"`
    // additional container to allow step separation in workflows
Steps []ProcessStep `yaml:"steps,omitempty"`
}


func CopyBoxesFileMixings(src *BoxesFileMixings) *BoxesFileMixings {
    if src == nil {
        return nil
    }
    var ret BoxesFileMixings

    if src.Title != nil {
        v := *src.Title
        ret.Title = &v
    }

    if src.Version != nil {
        v := *src.Version
        ret.Version = &v
    }

    ret.Legend = CopyLegend(src.Legend)

    if src.LayoutMixins != nil {
        ret.LayoutMixins = make(map[string]LayoutMixin, len(src.LayoutMixins))
        for k, v := range src.LayoutMixins {
            ret.LayoutMixins[k] = *CopyLayoutMixin(&v)
        }
    }

    if src.Connections != nil {
        ret.Connections = make(map[string]ConnectionCont, len(src.Connections))
        for k, v := range src.Connections {
            ret.Connections[k] = *CopyConnectionCont(&v)
        }
    }

    if src.Formats != nil {
        ret.Formats = make(map[string]Format, len(src.Formats))
        for k, v := range src.Formats {
            ret.Formats[k] = *CopyFormat(&v)
        }
    }

    ret.FormatVariations = CopyFormatVariations(src.FormatVariations)

    if src.Comments != nil {
        ret.Comments = make(map[string]types.Comment, len(src.Comments))
        for k, v := range src.Comments {
            ret.Comments[k] = *types.CopyComment(&v)
        }
    }

    if src.Tags != nil {
        ret.Tags = make(map[string]Tags, len(src.Tags))
        for k, v := range src.Tags {
            ret.Tags[k] = *CopyTags(&v)
        }
    }

    if src.Images != nil {
        ret.Images = make(map[string]types.ImageDef, len(src.Images))
        for k, v := range src.Images {
            ret.Images[k] = *types.CopyImageDef(&v)
        }
    }

    if src.Overlays != nil {
        ret.Overlays = make([]Overlay, len(src.Overlays))
        for i, v := range src.Overlays {
            ret.Overlays[i] = *CopyOverlay(&v)
        }
    }

    if src.Steps != nil {
        ret.Steps = make([]ProcessStep, len(src.Steps))
        for i, v := range src.Steps {
            ret.Steps[i] = *CopyProcessStep(&v)
        }
    }
return &ret
}


func NewBoxesFileMixings() *BoxesFileMixings {
    var ret BoxesFileMixings
    ret.LayoutMixins = make(map[string]LayoutMixin, 0)
    ret.Connections = make(map[string]ConnectionCont, 0)
    ret.Formats = make(map[string]Format, 0)
    ret.Comments = make(map[string]types.Comment, 0)
    ret.Tags = make(map[string]Tags, 0)
    ret.Images = make(map[string]types.ImageDef, 0)
    ret.Overlays = make([]Overlay, 0)
    ret.Steps = make([]ProcessStep, 0)
    return &ret
}
