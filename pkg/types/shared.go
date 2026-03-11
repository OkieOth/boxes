package types

// Attention, this file is generated. Manual changes get lost with the next
// run of the code generation.
// created by yacg (template: golang_types.mako v1.1.0)

import (
)


/* Object to define comments in the diagrams
*/
type Comment struct {

    // text of the comment
    Text string  `yaml:"text"`

    // optional number or a short text, displayed in the marker of that comment
    Label *string  `yaml:"label,omitempty"`

    // format to use to render this comment
    Format *string  `yaml:"format,omitempty"`
}


func CopyComment(src *Comment) *Comment {
    if src == nil {
        return nil
    }
    var ret Comment
    ret.Text = src.Text
    ret.Label = src.Label
    ret.Format = src.Format

    return &ret
}





/* container to define closer how a connection should be routed
*/
type ConnRestriction struct {

    // restricts the start of the connection
    Source *ConnRestItem  `yaml:"source,omitempty"`

    // restricts the end of the connection
    Dest *ConnRestItem  `yaml:"dest,omitempty"`
}


func CopyConnRestriction(src *ConnRestriction) *ConnRestriction {
    if src == nil {
        return nil
    }
    var ret ConnRestriction
    ret.Source = CopyConnRestItem(src.Source)
    ret.Dest = CopyConnRestItem(src.Dest)

    return &ret
}





type ConnRestItem struct {
}


func CopyConnRestItem(src *ConnRestItem) *ConnRestItem {
    if src == nil {
        return nil
    }
    var ret ConnRestItem

    return &ret
}





/* possible options to restrict a connection start or end
*/
type ConnRestrItem struct {

    // no connection start or end to the top border of a box
    NoTop bool  `yaml:"noTop"`

    // no connection start or end to the bottom border of a box
    NoBottom bool  `yaml:"noBottom"`

    // no connection start or end to the left border of a box
    NoLeft bool  `yaml:"noLeft"`

    // no connection start or end to the right border of a box
    NoRight bool  `yaml:"noRight"`
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




