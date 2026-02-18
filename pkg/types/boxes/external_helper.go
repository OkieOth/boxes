package boxes

import (
	"fmt"
	"maps"
	"slices"

	"github.com/okieoth/draw.chart.things/pkg/types"
)

// used in filter situations in cases where no ID are provided
var GlobalId int

func GetNewId() string {
	GlobalId++
	return fmt.Sprintf("id_xxx_%d", GlobalId)
}

func hasConnectionById(connections []Connection, destId string) bool {
	return slices.ContainsFunc(connections, func(c Connection) bool {
		return c.DestId == destId
	})
}

func hasConnectionByCapt(connections []Connection, caption string) bool {
	return slices.ContainsFunc(connections, func(c Connection) bool {
		return c.Dest == caption
	})
}

func (b *Boxes) mixInConnectionsImplCont(cont []Layout, additional map[string]ConnectionCont) {
	for i := range len(cont) {
		b.mixInConnectionsImpl(&cont[i], additional)
	}
}

func (b *Boxes) mixInConnectionsImpl(l *Layout, additional map[string]ConnectionCont) {
	if l.Id != "" {
		if cc, ok := additional[l.Id]; ok {
			for _, c := range cc.Connections {
				if !hasConnectionById(l.Connections, c.DestId) {
					l.Connections = append(l.Connections, c)
				}
			}
		}
	}
	if l.Caption != "" {
		if cc, ok := additional[l.Caption]; ok {
			for _, c := range cc.Connections {
				if !hasConnectionByCapt(l.Connections, c.Dest) {
					c.DestId = b.FindBoxWithCaption(c.Dest)
					l.Connections = append(l.Connections, c)
				}
			}
		}
	}
	b.mixInConnectionsImplCont(l.Horizontal, additional)
	b.mixInConnectionsImplCont(l.Vertical, additional)
}

func (b *Boxes) mixInTagsImplCont(cont []Layout, additional map[string]Tags) {
	for i := range len(cont) {
		b.mixInTagsImpl(&cont[i], additional)
	}
}

func (b *Boxes) mixInTagsImpl(l *Layout, additional map[string]Tags) {
	if l.Caption != "" {
		if tags, ok := additional[l.Caption]; ok {
			l.Tags = append(l.Tags, tags.Tags...)
		}
	}
	b.mixInTagsImplCont(l.Horizontal, additional)
	b.mixInTagsImplCont(l.Vertical, additional)
}

func (b *Boxes) mixInCommentImplCont(cont []Layout, additional map[string]types.Comment) {
	for i := range len(cont) {
		b.mixInCommentsImpl(&cont[i], additional)
	}
}

func (b *Boxes) mixInCommentsImpl(l *Layout, additional map[string]types.Comment) {
	if l.Id != "" {
		if c, ok := additional[l.Id]; ok {
			l.Comment = &c
		}
		if l.Caption != "" {
			if c, ok := additional[l.Caption]; ok {
				l.Comment = &c
			}
		}
	}
	b.mixInCommentImplCont(l.Horizontal, additional)
	b.mixInCommentImplCont(l.Vertical, additional)
}

func (b *Boxes) initIdForMixinsInCase(mixin []Layout) bool {
	ret := false
	for i := range mixin {
		m := &mixin[i]
		if m.Id == "" {
			ret = true
			m.Id = GetNewId()
		}
	}
	return ret
}

func (b *Boxes) mixInLayoutNow(l *Layout, mixin *LayoutMixin) {
	if mixin == nil {
		return
	}
	if len(mixin.Horizontal) > 0 {
		// mix in horizontal elements
		b.initIdForMixinsInCase(mixin.Horizontal)
		l.Horizontal = append(l.Horizontal, mixin.Horizontal...)
	}
	if len(mixin.Vertical) > 0 {
		// mix in vertical elements
		b.initIdForMixinsInCase(mixin.Vertical)
		l.Vertical = append(l.Vertical, mixin.Vertical...)
	}
}

func (b *Boxes) mixInLayoutsImplCont(cont []Layout, additional *map[string]LayoutMixin) {
	for i := range cont {
		if len(*additional) == 0 {
			return
		}
		b.mixInLayoutsImpl(&cont[i], additional)
	}
}

func (b *Boxes) mixInLayoutsImpl(l *Layout, additional *map[string]LayoutMixin) {
	if len(*additional) == 0 {
		return
	}
	if l.Caption != "" {
		if mixin, ok := (*additional)[l.Caption]; ok {
			b.mixInLayoutNow(l, &mixin)
			delete(*additional, l.Caption)
		}
	} else if l.Id != "" {
		if mixin, ok := (*additional)[l.Id]; ok {
			b.mixInLayoutNow(l, &mixin)
			delete(*additional, l.Id)
		}
	}
	b.mixInLayoutsImplCont(l.Horizontal, additional)
	b.mixInLayoutsImplCont(l.Vertical, additional)
}

func (b *Boxes) MixinThings(additional BoxesFileMixings) {
	if additional.Title != nil {
		b.Title += ": " + *additional.Title
		if additional.Version != nil {
			b.Title += fmt.Sprintf(" [%s]", *additional.Version)
		}
	}
	if additional.Legend != nil {
		if b.Legend == nil {
			b.Legend = NewLegend()
		}
		if len(additional.Legend.Entries) > 0 {
			b.Legend.Entries = append(b.Legend.Entries, additional.Legend.Entries...)
		}
	}
	if len(additional.Formats) > 0 {
		if b.Formats == nil {
			b.Formats = make(map[string]Format, 0)
		}
		maps.Copy(b.Formats, additional.Formats)
	}
	b.mixInLayoutsImpl(&b.Boxes, &additional.LayoutMixins)
	b.mixInConnectionsImpl(&b.Boxes, additional.Connections)
	b.mixInTagsImpl(&b.Boxes, additional.Tags)
	b.mixInCommentsImpl(&b.Boxes, additional.Comments)
	b.Overlays = append(b.Overlays, additional.Overlays...)
	if len(additional.Formats) > 0 {
		if b.Formats == nil {
			b.Formats = make(map[string]Format)
		}
		maps.Copy(b.Formats, additional.Formats)
	}
	if len(additional.Images) > 0 {
		if b.Images == nil {
			b.Images = make(map[string]types.ImageDef)
		}
		maps.Copy(b.Images, additional.Images)
	}
	if additional.FormatVariations != nil {
		if len(additional.FormatVariations.HasTag) > 0 {
			if b.FormatVariations == nil {
				b.FormatVariations = NewFormatVariations()
			}
			maps.Copy(b.FormatVariations.HasTag, additional.FormatVariations.HasTag)
		}
	}
}

func (b *Boxes) findBoxInContWithCaption(cont []Layout, caption string) string {
	if cont == nil {
		return ""
	}
	for i := range len(cont) {
		found := b.findBoxWithCaption(&cont[i], caption)
		if found != "" {
			return found
		}
	}
	return ""
}

func (b *Boxes) findBoxWithCaption(box *Layout, caption string) string {
	if box.Caption == caption {
		return box.Id
	}
	found := b.findBoxInContWithCaption(box.Vertical, caption)
	if found != "" {
		return found
	}
	return b.findBoxInContWithCaption(box.Horizontal, caption)
}

func (b *Boxes) FindBoxWithCaption(caption string) string {
	return b.findBoxWithCaption(&b.Boxes, caption)
}
