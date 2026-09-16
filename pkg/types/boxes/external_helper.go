package boxes

import (
	"fmt"
	"maps"
	"slices"

	"github.com/okieoth/boxes/pkg/types"
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

func (b *Boxes) mixInTagsImplCont(cont []Layout, additional map[string][]string) {
	for i := range len(cont) {
		b.mixInTagsImpl(&cont[i], additional)
	}
}

func (b *Boxes) mixInTagsImpl(l *Layout, additional map[string][]string) {
	if l.Caption != "" {
		if tags, ok := additional[l.Caption]; ok {
			l.Tags = append(l.Tags, tags...)
		} else if tags, ok := additional[l.Id]; ok {
			l.Tags = append(l.Tags, tags...)
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
	if mixin.WrapSubElems == nil {
		// normal mixin of things
		if len(mixin.Horizontal) > 0 {
			// mix in horizontal elements
			b.initIdForMixinsInCase(mixin.Horizontal)
			if mixin.PutAfter != nil {
				for i := range len(l.Horizontal) - 1 {
					e := l.Horizontal[i]
					if e.Caption == *mixin.PutAfter || e.Id == *mixin.PutAfter {
						l.Horizontal = slices.Insert(l.Horizontal, i+1, mixin.Horizontal...)
						return
					}
				}
			} else if mixin.PutBefore != nil {
				for i := range len(l.Horizontal) {
					e := l.Horizontal[i]
					if e.Caption == *mixin.PutBefore || e.Id == *mixin.PutBefore {
						l.Horizontal = slices.Insert(l.Horizontal, i, mixin.Horizontal...)
						return
					}
				}
			}
			l.Horizontal = append(l.Horizontal, mixin.Horizontal...)
		}
		if len(mixin.Vertical) > 0 {
			// mix in vertical elements
			b.initIdForMixinsInCase(mixin.Vertical)
			if mixin.PutAfter != nil {
				for i := range len(l.Vertical) {
					e := l.Vertical[i]
					if e.Caption == *mixin.PutAfter || e.Id == *mixin.PutAfter {
						l.Vertical = slices.Insert(l.Vertical, i+1, mixin.Vertical...)
						return
					}
				}
			} else if mixin.PutBefore != nil {
				for i := range len(l.Vertical) {
					e := l.Vertical[i]
					if e.Caption == *mixin.PutBefore || e.Id == *mixin.PutBefore {
						l.Vertical = slices.Insert(l.Vertical, i, mixin.Vertical...)
						return
					}
				}
			}
			l.Vertical = append(l.Vertical, mixin.Vertical...)
		}
	} else {
		// existing containers (vertical or horizontal) should be wrapped in a dedicated box
		newBox := newLayoutFromWrapSubMixin(mixin.WrapSubElems)
		if newBox != nil {
			// assign the container content to the new Boxes
			if len(l.Horizontal) > 0 {
				newBox.Horizontal = append(newBox.Horizontal, l.Horizontal...)
				l.Horizontal = []Layout{*newBox}
			} else if len(l.Vertical) > 0 {
				newBox.Vertical = append(newBox.Vertical, l.Vertical...)
				l.Vertical = []Layout{*newBox}
			}
		}
	}
}

func newLayoutFromWrapSubMixin(wrapperDef *SubsWrapperObj) *Layout {
	if wrapperDef == nil {
		return nil
	}
	newBox := NewLayout()
	newBox.Expand = true
	if wrapperDef.Caption != nil {
		newBox.Caption = *wrapperDef.Caption
	}
	if wrapperDef.Text1 != nil {
		newBox.Text1 = *wrapperDef.Text1
	}
	if wrapperDef.Text2 != nil {
		newBox.Text2 = *wrapperDef.Text2
	}
	if wrapperDef.Format != nil {
		newBox.Format = wrapperDef.Format
	}
	if len(wrapperDef.Tags) > 0 {
		newBox.Tags = append(newBox.Tags, wrapperDef.Tags...)
	}
	return newBox
}

func (b *Boxes) mixInLayoutsImplCont(cont []Layout, additional *map[string]LayoutMixin) {
	for i := range cont {
		if len(*additional) == 0 {
			return
		}
		b.mixInLayoutsImpl(&cont[i], additional)
	}
}

func (b *Boxes) mixInBoxesImplCont(cont *[]Layout, additional *[]BoxMixin) *[]Layout {
	if len(*additional) == 0 {
		return cont
	}

	offset := 0
	originalLen := len(*cont)

	for contIndex := 0; contIndex < originalLen; contIndex++ {
		i := contIndex + offset
		b.mixInBoxesImpl(&(*cont)[i], additional)
		elem := (*cont)[i]
		var toDelete []int
		for j, mixin := range *additional {
			matches := (elem.Caption != "" && mixin.Reference == elem.Caption) ||
				(elem.Id != "" && mixin.Reference == elem.Id)

			if !matches {
				continue
			}

			if mixin.PutAfter != nil && *mixin.PutAfter {
				cont = insertBoxesMixinInCont(cont, mixin, i, true)
				offset++
				toDelete = append(toDelete, j)
			} else if mixin.PutBefore != nil && *mixin.PutBefore {
				cont = insertBoxesMixinInCont(cont, mixin, i, false)
				offset++
				i = contIndex + offset
				toDelete = append(toDelete, j)
			}
		}
	}
	return cont
}

func insertBoxesMixinInCont(cont *[]Layout, mixin BoxMixin, insertIndex int, putAfter bool) *[]Layout {
	var tmp []Layout
	restIndex := insertIndex
	if insertIndex == 0 {
		if putAfter {
			tmp = slices.Clone((*cont)[0:1])
			restIndex = 1
		} else {
			tmp = make([]Layout, 0)
		}
	} else {
		if putAfter {
			tmp = slices.Clone((*cont)[0 : insertIndex+1])
			restIndex = insertIndex + 1
		} else {
			tmp = slices.Clone((*cont)[0:insertIndex])
		}
	}
	newLayout := NewLayout()
	newLayout.Id = mixin.Id
	newLayout.Caption = mixin.Caption
	newLayout.Text1 = mixin.Text1
	newLayout.Text2 = mixin.Text2
	newLayout.Comment = mixin.Comment
	newLayout.Image = mixin.Image
	newLayout.Expand = mixin.Expand
	newLayout.ExtVertical = mixin.ExtVertical
	newLayout.Vertical = mixin.Vertical
	newLayout.ExtHorizontal = mixin.ExtHorizontal
	newLayout.Horizontal = mixin.Horizontal
	newLayout.Tags = mixin.Tags
	newLayout.Connections = mixin.Connections
	newLayout.Format = mixin.Format
	newLayout.DontBlockConPaths = mixin.DontBlockConPaths
	newLayout.DataLink = mixin.DataLink
	newLayout.HiddenComments = mixin.HiddenComments

	tmp = append(tmp, *newLayout)
	if restIndex < len(*cont) {
		rest := (*cont)[restIndex:]
		tmp = append(tmp, rest...)
	}
	return &tmp
}

func (b *Boxes) mixInBoxesImpl(l *Layout, additional *[]BoxMixin) {
	if len(*additional) == 0 {
		return
	}
	l.Horizontal = *b.mixInBoxesImplCont(&l.Horizontal, additional)
	l.Vertical = *b.mixInBoxesImplCont(&l.Vertical, additional)
}

func (b *Boxes) mixInLayoutsImpl(l *Layout, additional *map[string]LayoutMixin) {
	if len(*additional) == 0 {
		return
	}
	handled := false
	if l.Caption != "" {
		if mixin, ok := (*additional)[l.Caption]; ok {
			b.mixInLayoutNow(l, &mixin)
			delete(*additional, l.Caption)
			handled = true
		}
	}
	if (!handled) && (l.Id != "") {
		if mixin, ok := (*additional)[l.Id]; ok {
			b.mixInLayoutNow(l, &mixin)
			delete(*additional, l.Id)
		}
	}
	b.mixInLayoutsImplCont(l.Horizontal, additional)
	b.mixInLayoutsImplCont(l.Vertical, additional)
}

func (b *Boxes) mixinLegend(legend *Legend) {
	if legend != nil {
		if b.Legend == nil {
			b.Legend = NewLegend()
			b.Legend.Entries = append(b.Legend.Entries, legend.Entries...)
		} else {
			for i := range legend.Entries {
				e := legend.Entries[i]
				if !slices.ContainsFunc(b.Legend.Entries, func(c LegendEntry) bool {
					return c.Text == e.Text && c.Format == e.Format
				}) {
					b.Legend.Entries = append(b.Legend.Entries, e)
				}
			}
		}
	}
}

func (b *Boxes) MixinThings(additional BoxesFileMixings) {
	if additional.Title != nil {
		b.Title += ": " + *additional.Title
		if additional.Version != nil {
			b.Title += fmt.Sprintf(" [%s]", *additional.Version)
		}
	}
	b.mixinLegend(additional.Legend)
	if len(additional.Formats) > 0 {
		if b.Formats == nil {
			b.Formats = make(map[string]Format)
		}
		for k, v := range additional.Formats {
			if existingFormat, ok := b.Formats[k]; ok {
				b.Formats[k] = mergeFormats(existingFormat, v)
			} else {
				b.Formats[k] = v
			}
		}
	}
	b.mixInLayoutsImpl(&b.Boxes, &additional.LayoutMixins)
	b.mixInBoxesImpl(&b.Boxes, &additional.BoxMixins)
	b.mixInConnectionsImpl(&b.Boxes, additional.Connections)
	b.mixInTagsImpl(&b.Boxes, additional.Tags)
	b.mixInCommentsImpl(&b.Boxes, additional.Comments)
	b.Overlays = append(b.Overlays, additional.Overlays...)
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

func mergeFontDef(existingFormat, newFormat *types.FontDef) *types.FontDef {
	if newFormat.Size != 0 {
		existingFormat.Size = newFormat.Size
	}
	if newFormat.Font != "" {
		existingFormat.Font = newFormat.Font
	}

	if newFormat.Type != nil {
		existingFormat.Type = newFormat.Type
	}
	if newFormat.Weight != nil {
		existingFormat.Weight = newFormat.Weight
	}

	if newFormat.LineHeight != 0.0 {
		existingFormat.LineHeight = newFormat.LineHeight
	}

	if newFormat.Color != "" {
		existingFormat.Color = newFormat.Color
	}

	if newFormat.Aligned != nil {
		existingFormat.Aligned = newFormat.Aligned
	}

	if newFormat.SpaceTop != 0 {
		existingFormat.SpaceTop = newFormat.SpaceTop
	}
	if newFormat.SpaceBottom != 0 {
		existingFormat.SpaceBottom = newFormat.SpaceBottom
	}
	if newFormat.MaxLenBeforeBreak != 0 {
		existingFormat.MaxLenBeforeBreak = newFormat.MaxLenBeforeBreak
	}
	return existingFormat
}

func mergeFillDef(existingFormat, newFormat *types.FillDef) *types.FillDef {
	if newFormat.Color != nil {
		existingFormat.Color = newFormat.Color
	}
	if newFormat.Opacity != nil {
		existingFormat.Opacity = newFormat.Opacity
	}
	return existingFormat
}

func mergeLineDef(existingFormat, newFormat *types.LineDef) *types.LineDef {
	if newFormat.Width != nil {
		existingFormat.Width = newFormat.Width
	}
	if newFormat.Style != nil {
		existingFormat.Style = newFormat.Style
	}
	if newFormat.Color != nil {
		existingFormat.Color = newFormat.Color
	}
	if newFormat.Opacity != nil {
		existingFormat.Opacity = newFormat.Opacity
	}
	return existingFormat
}

func mergeFormats(existingFormat, newFormat Format) Format {
	if newFormat.WidthOfParent != nil {
		existingFormat.WidthOfParent = newFormat.WidthOfParent
	}

	if newFormat.FixedWidth != nil {
		existingFormat.FixedWidth = newFormat.FixedWidth
	}
	if newFormat.FixedHeight != nil {
		existingFormat.FixedHeight = newFormat.FixedHeight
	}
	if newFormat.VerticalTxt != nil {
		existingFormat.VerticalTxt = newFormat.VerticalTxt
	}
	if newFormat.FontCaption != nil {
		existingFormat.FontCaption = mergeFontDef(existingFormat.FontCaption, newFormat.FontCaption)
	}
	if newFormat.FontText1 != nil {
		existingFormat.FontText1 = mergeFontDef(existingFormat.FontText1, newFormat.FontText1)
	}
	if newFormat.FontText2 != nil {
		existingFormat.FontText2 = mergeFontDef(existingFormat.FontText2, newFormat.FontText2)
	}
	if newFormat.FontComment != nil {
		existingFormat.FontComment = mergeFontDef(existingFormat.FontComment, newFormat.FontComment)
	}
	if newFormat.FontCommentMarker != nil {
		existingFormat.FontCommentMarker = mergeFontDef(existingFormat.FontCommentMarker, newFormat.FontCommentMarker)
	}
	if newFormat.Line != nil {
		existingFormat.Line = mergeLineDef(existingFormat.Line, newFormat.Line)
	}
	if newFormat.Fill != nil {
		existingFormat.Fill = mergeFillDef(existingFormat.Fill, newFormat.Fill)
	}
	if newFormat.Padding != nil {
		existingFormat.Padding = newFormat.Padding
	}
	if newFormat.BoxMargin != nil {
		existingFormat.BoxMargin = newFormat.BoxMargin
	}
	if newFormat.CornerRadius != nil {
		existingFormat.CornerRadius = newFormat.CornerRadius
	}
	if newFormat.RenderType != nil {
		existingFormat.RenderType = newFormat.RenderType
	}
	return existingFormat
}

func mergeStepLayoutMixins(additional *BoxesFileMixings, step ProcessStep, stepIdx int) {
	if additional.LayoutMixins == nil {
		additional.LayoutMixins = make(map[string]LayoutMixin)
	}
	for k, v := range step.LayoutMixins {
		if existing, ok := additional.LayoutMixins[k]; ok {
			if len(v.Horizontal) > 0 {
				if v.PutAfter != nil {
					for i := range len(existing.Horizontal) - 1 {
						e := existing.Horizontal[i]
						if e.Caption == *v.PutAfter || e.Id == *v.PutAfter {
							existing.Horizontal = slices.Insert(existing.Horizontal, i+1, v.Horizontal...)
							break
						}
					}
				} else if v.PutBefore != nil {
					for i := range len(existing.Horizontal) {
						e := existing.Horizontal[i]
						if e.Caption == *v.PutBefore || e.Id == *v.PutBefore {
							existing.Horizontal = slices.Insert(existing.Horizontal, i, v.Horizontal...)
							break
						}
					}
				} else {
					existing.Horizontal = append(existing.Horizontal, v.Horizontal...)
				}
			}
			if len(v.Vertical) > 0 {
				if v.PutAfter != nil {
					for i := range len(existing.Vertical) - 1 {
						e := existing.Vertical[i]
						if e.Caption == *v.PutAfter || e.Id == *v.PutAfter {
							existing.Vertical = slices.Insert(existing.Vertical, i+1, v.Vertical...)
							break
						}
					}
				} else if v.PutBefore != nil {
					for i := range len(existing.Vertical) {
						e := existing.Vertical[i]
						if e.Caption == *v.PutBefore || e.Id == *v.PutBefore {
							existing.Vertical = slices.Insert(existing.Vertical, i, v.Vertical...)
							break
						}
					}
				} else {
					existing.Vertical = append(existing.Vertical, v.Vertical...)
				}
			}
			if v.WrapSubElems != nil {
				existing.WrapSubElems = v.WrapSubElems
			}
			additional.LayoutMixins[k] = existing
		} else {
			additional.LayoutMixins[k] = v
		}
	}
}

func mergeStepBoxMixins(additional *BoxesFileMixings, step ProcessStep, stepIdx int) {
	if additional.BoxMixins == nil {
		additional.BoxMixins = make([]BoxMixin, 0)
	}
	additional.BoxMixins = append(additional.BoxMixins, step.BoxMixins...)
}

func mergeStepConnections(additional *BoxesFileMixings, step ProcessStep, stepIdx int) {
	if additional.Connections == nil {
		additional.Connections = make(map[string]ConnectionCont)
	}
	for k, v := range step.Connections {
		tagged := ConnectionCont{Connections: make([]Connection, len(v.Connections))}
		for i, c := range v.Connections {
			c.Step = &stepIdx
			tagged.Connections[i] = c
		}
		if existing, ok := additional.Connections[k]; ok {
			existing.Connections = append(existing.Connections, tagged.Connections...)
			additional.Connections[k] = existing
		} else {
			additional.Connections[k] = tagged
		}
	}
}

func mergeStepComments(additional *BoxesFileMixings, step ProcessStep, stepIdx int) {
	if additional.Comments == nil {
		additional.Comments = make(map[string]types.Comment)
	}
	for k, v := range step.Comments {
		v.Step = &stepIdx
		additional.Comments[k] = v
	}
}

func mergeStepFormats(additional *BoxesFileMixings, step ProcessStep, stepIdx int) {
	if additional.Formats == nil {
		additional.Formats = make(map[string]Format)
	}
	for k, v := range step.Formats {
		if existingFormat, ok := additional.Formats[k]; ok {
			additional.Formats[k] = mergeFormats(existingFormat, v)
		} else {
			additional.Formats[k] = v
		}
	}
}

func mergeStepTags(additional *BoxesFileMixings, step ProcessStep) {
	if len(step.Tags) == 0 {
		return
	}
	if additional.Tags == nil {
		additional.Tags = make(map[string][]string, 0)
	}
	for k, v := range step.Tags {
		if existing, ok := additional.Tags[k]; ok {
			existing = append(existing, v...)
			additional.Tags[k] = existing
		} else {
			additional.Tags[k] = v
		}
	}
}

func mergeStepOverlays(additional *BoxesFileMixings, step ProcessStep) {
	additional.Overlays = append(additional.Overlays, step.Overlays...)
}

func mergeStepFormatVariations(additional *BoxesFileMixings, step ProcessStep) {
	if step.FormatVariations == nil || len(step.FormatVariations.HasTag) == 0 {
		return
	}
	if additional.FormatVariations == nil {
		additional.FormatVariations = NewFormatVariations()
	}
	maps.Copy(additional.FormatVariations.HasTag, step.FormatVariations.HasTag)
}

// MixinThingsWithSteps applies the mixin including only the specified workflow steps
// (by index). Root-level connections and comments are always applied as the base layer.
// If activeSteps is empty, only root-level content is applied.
func (b *Boxes) MixinThingsWithSteps(additional BoxesFileMixings, activeSteps []int) {
	if len(additional.Steps) > 0 && len(activeSteps) > 0 {
		for _, idx := range activeSteps {
			if idx < 0 || idx >= len(additional.Steps) {
				continue
			}
			step := additional.Steps[idx]
			mergeStepLayoutMixins(&additional, step, idx)
			mergeStepBoxMixins(&additional, step, idx)
			mergeStepConnections(&additional, step, idx)
			mergeStepComments(&additional, step, idx)
			mergeStepFormats(&additional, step, idx)
			mergeStepTags(&additional, step)
			mergeStepOverlays(&additional, step)
			mergeStepFormatVariations(&additional, step)
		}
	}
	b.MixinThings(additional)
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
