// Attention, this code is generated. Do not modify manually. Changes will
// be overwritten be the next codegen run.
// Generated from Shared Types (configs/models/shared.json)



// Types that are shared between different models

package types



// Object to define comments in the diagrams
type Comment struct {
    // text of the comment
Text string `yaml:"text"`
    // optional number or a short text, displayed in the marker of that comment
Label *string `yaml:"label,omitempty"`
    // format to use to render this comment
Format *string `yaml:"format,omitempty"`
    // optional step where this comment is part of, is filled via processing not by the user
Step *int `yaml:"step,omitempty"`
}


func CopyComment(src *Comment) *Comment {
    if src == nil {
        return nil
    }
    var ret Comment

    ret.Text = src.Text

    if src.Label != nil {
        v := *src.Label
        ret.Label = &v
    }

    if src.Format != nil {
        v := *src.Format
        ret.Format = &v
    }

    if src.Step != nil {
        v := *src.Step
        ret.Step = &v
    }
return &ret
}


func NewComment() *Comment {
    var ret Comment
    return &ret
}

// container to define closer how a connection should be routed
type ConnRestriction struct {
    // restricts the start of the connection
Source *ConnRestrItem `yaml:"source,omitempty"`
    // restricts the end of the connection
Dest *ConnRestrItem `yaml:"dest,omitempty"`
}


func CopyConnRestriction(src *ConnRestriction) *ConnRestriction {
    if src == nil {
        return nil
    }
    var ret ConnRestriction

    ret.Source = CopyConnRestrItem(src.Source)

    ret.Dest = CopyConnRestrItem(src.Dest)
return &ret
}


func NewConnRestriction() *ConnRestriction {
    var ret ConnRestriction
    return &ret
}

// possible options to restrict a connection start or end
type ConnRestrItem struct {
    // no connection start or end to the top border of a box
NoTop bool `yaml:"noTop"`
    // no connection start or end to the bottom border of a box
NoBottom bool `yaml:"noBottom"`
    // no connection start or end to the left border of a box
NoLeft bool `yaml:"noLeft"`
    // no connection start or end to the right border of a box
NoRight bool `yaml:"noRight"`
}


func CopyConnRestrItem(src *ConnRestrItem) *ConnRestrItem {
    if src == nil {
        return nil
    }
    var ret ConnRestrItem

    ret.NoTop = src.NoTop

    ret.NoBottom = src.NoBottom

    ret.NoLeft = src.NoLeft

    ret.NoRight = src.NoRight
return &ret
}


func NewConnRestrItem() *ConnRestrItem {
    var ret ConnRestrItem
    return &ret
}
