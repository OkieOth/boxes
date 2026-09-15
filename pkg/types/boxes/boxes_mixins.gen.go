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
    // dictionary for layout mixins. key of the dictionary is the caption of the box that will take the additional content
LayoutMixins map[string]LayoutMixin `yaml:"layoutMixins,omitempty"`
    // array for box mixins. It's useful in cases where the dictionary reference to boxes to mix in is limiting
BoxMixins []BoxMixin `yaml:"boxMixins,omitempty"`
    // dictionary of connection objects
Connections map[string]ConnectionCont `yaml:"connections,omitempty"`
    // dictionary of comment objects, this comment will applied on layout objects and replace existing comments there
Comments map[string]types.Comment `yaml:"comments,omitempty"`
    // dictionary of tag array, the additional tags will be applied on the existing layout and can be used for instance to define display formats
Tags map[string][]string `yaml:"tags,omitempty"`
Overlays []Overlay `yaml:"overlays,omitempty"`
    // Be careful it can mess up the stype of the picture! This allows to overwrite global formats in an available step. It can be helpful in cases this step wants to highlight things, that are visual pressed back by existing formats
Formats map[string]Format `yaml:"formats,omitempty"`
    // Set of formats that overwrites the style of boxes, if specific conditions are met
FormatVariations *FormatVariations `yaml:"formatVariations,omitempty"`
}


func CopyProcessStep(src *ProcessStep) *ProcessStep {
    if src == nil {
        return nil
    }
    var ret ProcessStep

    ret.Caption = src.Caption

    if src.LayoutMixins != nil {
        ret.LayoutMixins = make(map[string]LayoutMixin, len(src.LayoutMixins))
        for k, v := range src.LayoutMixins {
            ret.LayoutMixins[k] = *CopyLayoutMixin(&v)
        }
    }

    if src.BoxMixins != nil {
        ret.BoxMixins = make([]BoxMixin, len(src.BoxMixins))
        for i, v := range src.BoxMixins {
            ret.BoxMixins[i] = *CopyBoxMixin(&v)
        }
    }

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
        ret.Tags = make(map[string][]string, len(src.Tags))
        for k, v := range src.Tags {
            ret.Tags[k] = v
        }
    }

    if src.Overlays != nil {
        ret.Overlays = make([]Overlay, len(src.Overlays))
        for i, v := range src.Overlays {
            ret.Overlays[i] = *CopyOverlay(&v)
        }
    }

    if src.Formats != nil {
        ret.Formats = make(map[string]Format, len(src.Formats))
        for k, v := range src.Formats {
            ret.Formats[k] = *CopyFormat(&v)
        }
    }

    ret.FormatVariations = CopyFormatVariations(src.FormatVariations)
return &ret
}


func NewProcessStep() *ProcessStep {
    var ret ProcessStep
    ret.LayoutMixins = make(map[string]LayoutMixin, 0)
    ret.BoxMixins = make([]BoxMixin, 0)
    ret.Connections = make(map[string]ConnectionCont, 0)
    ret.Comments = make(map[string]types.Comment, 0)
    ret.Tags = make(map[string][]string, 0)
    ret.Overlays = make([]Overlay, 0)
    ret.Formats = make(map[string]Format, 0)
    return &ret
}

// definition of a box to be mixed in
type BoxMixin struct {
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
    // either ID or caption of the box where this mixin is placed before or after
Reference string `yaml:"reference"`
    // triggers the new mixin elemente to put either left or above the object with the given ID or caption
PutBefore *string `yaml:"putBefore,omitempty"`
    // triggers the new mixin elemente to put either right or below the object with the given ID or caption
PutAfter *string `yaml:"putAfter,omitempty"`
}


func CopyBoxMixin(src *BoxMixin) *BoxMixin {
    if src == nil {
        return nil
    }
    var ret BoxMixin

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

    ret.Reference = src.Reference

    if src.PutBefore != nil {
        v := *src.PutBefore
        ret.PutBefore = &v
    }

    if src.PutAfter != nil {
        v := *src.PutAfter
        ret.PutAfter = &v
    }
return &ret
}


func NewBoxMixin() *BoxMixin {
    var ret BoxMixin
    ret.Vertical = make([]Layout, 0)
    ret.Horizontal = make([]Layout, 0)
    ret.Tags = make([]string, 0)
    ret.Connections = make([]Connection, 0)
    return &ret
}

type BoxMixinBase struct {
    // either ID or caption of the box where this mixin is placed before or after
Reference string `yaml:"reference"`
    // triggers the new mixin elemente to put either left or above the object with the given ID or caption
PutBefore *string `yaml:"putBefore,omitempty"`
    // triggers the new mixin elemente to put either right or below the object with the given ID or caption
PutAfter *string `yaml:"putAfter,omitempty"`
}


func CopyBoxMixinBase(src *BoxMixinBase) *BoxMixinBase {
    if src == nil {
        return nil
    }
    var ret BoxMixinBase

    ret.Reference = src.Reference

    if src.PutBefore != nil {
        v := *src.PutBefore
        ret.PutBefore = &v
    }

    if src.PutAfter != nil {
        v := *src.PutAfter
        ret.PutAfter = &v
    }
return &ret
}


func NewBoxMixinBase() *BoxMixinBase {
    var ret BoxMixinBase
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
    // array for box mixins. It's useful in cases where the dictionary reference to boxes to mix in is limiting
BoxMixins []BoxMixin `yaml:"boxMixins,omitempty"`
    // dictionary of connection objects
Connections map[string]ConnectionCont `yaml:"connections,omitempty"`
Formats map[string]Format `yaml:"formats,omitempty"`
    // Set of formats that overwrites the style of boxes, if specific conditions are met
FormatVariations *FormatVariations `yaml:"formatVariations,omitempty"`
    // dictionary of comment objects, this comment will applied on layout objects and replace existing comments there
Comments map[string]types.Comment `yaml:"comments,omitempty"`
    // dictionary of tag array, the additional tags will be applied on the existing layout and can be used for instance to define display formats
Tags map[string][]string `yaml:"tags,omitempty"`
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

    if src.BoxMixins != nil {
        ret.BoxMixins = make([]BoxMixin, len(src.BoxMixins))
        for i, v := range src.BoxMixins {
            ret.BoxMixins[i] = *CopyBoxMixin(&v)
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
        ret.Tags = make(map[string][]string, len(src.Tags))
        for k, v := range src.Tags {
            ret.Tags[k] = v
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
    ret.BoxMixins = make([]BoxMixin, 0)
    ret.Connections = make(map[string]ConnectionCont, 0)
    ret.Formats = make(map[string]Format, 0)
    ret.Comments = make(map[string]types.Comment, 0)
    ret.Tags = make(map[string][]string, 0)
    ret.Images = make(map[string]types.ImageDef, 0)
    ret.Overlays = make([]Overlay, 0)
    ret.Steps = make([]ProcessStep, 0)
    return &ret
}
