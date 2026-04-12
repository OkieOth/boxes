// Attention, this code is generated. Do not modify manually. Changes will
// be overwritten be the next codegen run.
// Generated from Formats (configs/models/formats.json)



// Some types to define base formats

package types



type Type string

const (
    TypeNormal Type = "normal"
    TypeItalic Type = "italic"
)

type Weight string

const (
    WeightNormal Weight = "normal"
    WeightBold Weight = "bold"
)

type Aligned string

const (
    AlignedLeft Aligned = "left"
    AlignedCenter Aligned = "center"
    AlignedRight Aligned = "right"
)

type Anchor string

const (
    AnchorMiddle Anchor = "middle"
    AnchorLeft Anchor = "left"
    AnchorRight Anchor = "right"
)

// Defines the font a text
type FontDef struct {
Size int `yaml:"size"`
Font string `yaml:"font"`
Type *Type `yaml:"type,omitempty"`
Weight *Weight `yaml:"weight,omitempty"`
    // Line height of the box
LineHeight float32 `yaml:"lineHeight"`
Color string `yaml:"color"`
Aligned *Aligned `yaml:"aligned,omitempty"`
SpaceTop int `yaml:"spaceTop"`
SpaceBottom int `yaml:"spaceBottom"`
    // Maximum length of the text before it breaks
MaxLenBeforeBreak int `yaml:"maxLenBeforeBreak"`
Anchor Anchor `yaml:"anchor"`
}


func CopyFontDef(src *FontDef) *FontDef {
    if src == nil {
        return nil
    }
    var ret FontDef

    ret.Size = src.Size

    ret.Font = src.Font

    if src.Type != nil {
        v := *src.Type
        ret.Type = &v
    }

    if src.Weight != nil {
        v := *src.Weight
        ret.Weight = &v
    }

    ret.LineHeight = src.LineHeight

    ret.Color = src.Color

    if src.Aligned != nil {
        v := *src.Aligned
        ret.Aligned = &v
    }

    ret.SpaceTop = src.SpaceTop

    ret.SpaceBottom = src.SpaceBottom

    ret.MaxLenBeforeBreak = src.MaxLenBeforeBreak

    ret.Anchor = src.Anchor
return &ret
}


func NewFontDef() *FontDef {
    var ret FontDef
    return &ret
}

type Style string

const (
    StyleSolid Style = "solid"
    StyleDotted Style = "dotted"
    StyleDashed Style = "dashed"
)

// Defines how the border of the box looks like
type LineDef struct {
Width *float64 `yaml:"width,omitempty"`
Style *Style `yaml:"style,omitempty"`
Color *string `yaml:"color,omitempty"`
Opacity *float64 `yaml:"opacity,omitempty"`
}


func CopyLineDef(src *LineDef) *LineDef {
    if src == nil {
        return nil
    }
    var ret LineDef

    if src.Width != nil {
        v := *src.Width
        ret.Width = &v
    }

    if src.Style != nil {
        v := *src.Style
        ret.Style = &v
    }

    if src.Color != nil {
        v := *src.Color
        ret.Color = &v
    }

    if src.Opacity != nil {
        v := *src.Opacity
        ret.Opacity = &v
    }
return &ret
}


func NewLineDef() *LineDef {
    var ret LineDef
    return &ret
}

// Defines the fill of the box
type FillDef struct {
Color *string `yaml:"color,omitempty"`
Opacity *float64 `yaml:"opacity,omitempty"`
}


func CopyFillDef(src *FillDef) *FillDef {
    if src == nil {
        return nil
    }
    var ret FillDef

    if src.Color != nil {
        v := *src.Color
        ret.Color = &v
    }

    if src.Opacity != nil {
        v := *src.Opacity
        ret.Opacity = &v
    }
return &ret
}


func NewFillDef() *FillDef {
    var ret FillDef
    return &ret
}

// parameters of an image to be displayed in the SVG
type ImageDef struct {
    // some words to explain what this image is about
Description *string `yaml:"description,omitempty"`
    // with of the displayed image
Width int `yaml:"width"`
    // height of the displayed image
Height int `yaml:"height"`
    // base64 string of the image to use
Base64 *string `yaml:"base64,omitempty"`
    // distance top and bottom of the image
MarginTopBottom *int `yaml:"marginTopBottom,omitempty"`
    // distance left and right of the image
MarginLeftRight *int `yaml:"marginLeftRight,omitempty"`
    // file path to a text file that contains the base64 of the png
Base64Src *string `yaml:"base64Src,omitempty"`
}


func CopyImageDef(src *ImageDef) *ImageDef {
    if src == nil {
        return nil
    }
    var ret ImageDef

    if src.Description != nil {
        v := *src.Description
        ret.Description = &v
    }

    ret.Width = src.Width

    ret.Height = src.Height

    if src.Base64 != nil {
        v := *src.Base64
        ret.Base64 = &v
    }

    if src.MarginTopBottom != nil {
        v := *src.MarginTopBottom
        ret.MarginTopBottom = &v
    }

    if src.MarginLeftRight != nil {
        v := *src.MarginLeftRight
        ret.MarginLeftRight = &v
    }

    if src.Base64Src != nil {
        v := *src.Base64Src
        ret.Base64Src = &v
    }
return &ret
}


func NewImageDef() *ImageDef {
    var ret ImageDef
    return &ret
}
