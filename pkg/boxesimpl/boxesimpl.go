package boxesimpl

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/okieoth/draw.chart.things/pkg/svgdrawing"
	"github.com/okieoth/draw.chart.things/pkg/types"
	"github.com/okieoth/draw.chart.things/pkg/types/boxes"
)

func LoadBoxesFromFile(inputFile string) (*boxes.Boxes, error) {
	layout, err := types.LoadInputFromFile[boxes.Boxes](inputFile)
	if err != nil {
		return nil, fmt.Errorf("error while loading boxes from file: %w", err)
	}
	err = replaceAlternativePaths(&layout.Boxes, inputFile)
	if err != nil {
		return nil, fmt.Errorf("error while replacing relative paths: %w", err)
	}
	if layout != nil && layout.Version != nil && layout.Title != "" {
		layout.Title += fmt.Sprintf(" [%s]", *layout.Version)
	}
	return layout, nil
}

func DrawBoxesFromFile(inputFile, outputFile string) error {
	layout, err := LoadBoxesFromFile(inputFile)
	if err != nil {
		return fmt.Errorf("error while loading input: %v", err)
	}
	textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()

	doc, err := InitialLayoutBoxes(layout, textDimensionCalulator)
	if err != nil {
		return err
	}

	doc.ConnectBoxes()
	doc.AdjustDocHeightTLegend(textDimensionCalulator)
	doc.IncludeComments(textDimensionCalulator)
	doc.IncludeOverlays(textDimensionCalulator)
	output, err := os.Create(outputFile)
	svgdrawing := svgdrawing.NewDrawing(output)
	svgdrawing.Start(doc.Title, doc.Height, doc.Width)
	svgdrawing.InitImages(doc.Images)
	svgdrawing.DrawRaster(doc.Width, doc.Height, types.RasterSize)
	doc.DrawBoxes(svgdrawing)
	doc.DrawConnections(svgdrawing)
	doc.DrawOverlays(svgdrawing, textDimensionCalulator)
	doc.DrawTitle(svgdrawing, textDimensionCalulator)
	doc.DrawLegend(svgdrawing, textDimensionCalulator)
	doc.DrawComments(svgdrawing, textDimensionCalulator)
	svgdrawing.Done()
	output.Close()
	return nil
}

type UIReturn struct {
	SVG      string
	ErrorMsg string
}

func isRelatedToId(b boxes.Layout, filter []string) bool {
	if slices.Contains(filter, b.Id) {
		return true
	}
	for _, e := range b.Horizontal {
		if isRelatedToId(e, filter) {
			return true
		}
	}
	for _, e := range b.Vertical {
		if isRelatedToId(e, filter) {
			return true
		}
	}
	return false
}

func collectTruncatedFromCont(b []boxes.Layout, newId string, truncated *map[string]TruncatedInfo) {
	for _, l := range b {
		collectTruncated(l, newId, truncated)
	}
}

func collectTruncated(b boxes.Layout, newId string, truncated *map[string]TruncatedInfo) {
	if b.Id != "" {
		(*truncated)[b.Id] = TruncatedInfo{
			truncated: b,
			newId:     newId,
		}
	}
	collectTruncatedFromCont(b.Horizontal, newId, truncated)
	collectTruncatedFromCont(b.Vertical, newId, truncated)
}

type TruncatedInfo struct {
	truncated boxes.Layout
	newId     string
}

func addTruncated(src map[string]TruncatedInfo, dest *map[string]TruncatedInfo) {
	for k, v := range src {
		(*dest)[k] = v
	}
}

func truncBoxes(b boxes.Layout, currentDepth, maxDepth int, expanded, blacklisted []string) (boxes.Layout, map[string]TruncatedInfo) {
	truncatedBoxes := make(map[string]TruncatedInfo, 0)
	if (currentDepth >= maxDepth) && (!b.Expand) && (!isRelatedToId(b, expanded)) {
		// possible removed connections to this object
		if b.Id == "" {
			b.Id = boxes.GetNewId()
		}
		collectTruncatedFromCont(b.Horizontal, b.Id, &truncatedBoxes)
		collectTruncatedFromCont(b.Vertical, b.Id, &truncatedBoxes)
		b.Horizontal = make([]boxes.Layout, 0)
		b.Vertical = make([]boxes.Layout, 0)
	} else {
		if len(b.Horizontal) > 0 {
			cont := make([]boxes.Layout, 0)
			for i := range len(b.Horizontal) {
				id := b.Horizontal[i].Id
				if blacklisted != nil && id != "" && slices.Contains(blacklisted, id) {
					// possible removed connections to this object
					if b.Id == "" {
						b.Id = boxes.GetNewId()
					}
					collectTruncated(b.Horizontal[i], b.Id, &truncatedBoxes)
					continue
				}
				l, trunc := truncBoxes(b.Horizontal[i], currentDepth+1, maxDepth, expanded, blacklisted)
				addTruncated(trunc, &truncatedBoxes)
				cont = append(cont, l)
			}
			b.Horizontal = cont
		}
		if len(b.Vertical) > 0 {
			cont := make([]boxes.Layout, 0)
			for i := range len(b.Vertical) {
				id := b.Vertical[i].Id
				if blacklisted != nil && id != "" && slices.Contains(blacklisted, id) {
					// possible removed connections to this object
					if b.Id == "" {
						b.Id = boxes.GetNewId()
					}
					collectTruncated(b.Vertical[i], b.Id, &truncatedBoxes)
					continue
				}
				l, trunc := truncBoxes(b.Vertical[i], currentDepth+1, maxDepth, expanded, blacklisted)
				addTruncated(trunc, &truncatedBoxes)
				cont = append(cont, l)
			}
			b.Vertical = cont
		}
	}
	return b, truncatedBoxes
}

func connectionExistsByDestId(connections []boxes.Connection, destId string) bool {
	return slices.ContainsFunc(connections, func(c boxes.Connection) bool {
		return c.DestId == destId
	})
}

func copyTruncatedComments(layout *boxes.Layout, truncatedObjects map[string]TruncatedInfo) {
	// find comments from truncated objects to the annotate current object
	found := false
	for _, v := range truncatedObjects {
		if v.newId == layout.Id {
			if v.truncated.Comment != nil && (v.truncated.Comment.Text != "" || v.truncated.Comment.Label != nil) {
				found = true
				break
			}
		}
	}
	if found {
		layout.HiddenComments = true
	}
}

func copyTruncatedTags(layout *boxes.Layout, truncatedObjects map[string]TruncatedInfo) {
	// find comments from truncated objects to the annotate current object
	for _, v := range truncatedObjects {
		if v.newId == layout.Id {
			if len(v.truncated.Tags) > 0 {
				for i := range v.truncated.Tags {
					tag := v.truncated.Tags[i]
					if strings.HasPrefix(tag, "_") {
						continue
					}
					if !slices.Contains(layout.Tags, tag) {
						layout.Tags = append(layout.Tags, tag)
					}
				}
			}
		}
	}
}

func copyTruncatedConnections(layout *boxes.Layout, truncatedObjects map[string]TruncatedInfo) {
	// copy all needed connections from truncated objects to the current object
	// Seems to have no impact
	// for i := range layout.Connections {
	// 	c := layout.Connections[i]
	// 	if truncated, ok := truncatedObjects[c.DestId]; ok {
	// 		c.Dest = fmt.Sprintf("CHANGED-1: %s", c.DestId)
	// 		c.DestId = truncated.newId
	// 	}
	// }
	for _, v := range truncatedObjects {
		if v.newId == layout.Id {
			// copy all truncated connections
			for _, c := range v.truncated.Connections {
				if c.DestId == layout.Id {
					continue
				}
				// destId is also replaced
				destIdToUse := c.DestId
				if obj, ok := truncatedObjects[c.DestId]; ok {
					destIdToUse = obj.newId
				}
				if !connectionExistsByDestId(layout.Connections, destIdToUse) {
					if destIdToUse != layout.Id {
						c.Dest = fmt.Sprintf("CHANGED-2: %s", c.DestId)
						c.DestId = destIdToUse
						(*layout).Connections = append(layout.Connections, c)
					}
				} else if c.Comment != nil {
					for i := range layout.Connections {
						con := &layout.Connections[i]
						if con.DestId == destIdToUse {
							con.HiddenComments = true
						}
					}
				}
			}
		}
	}
}

func adjustDestIdInRespectOfTruncated(layout *boxes.Layout, truncatedObjects map[string]TruncatedInfo) {
	normalizedConnections := make([]boxes.Connection, 0)
	for _, c := range layout.Connections {
		if trunc, found := truncatedObjects[c.DestId]; found {
			c.DestId = trunc.newId
		}
		if !connectionExistsByDestId(normalizedConnections, c.DestId) {
			normalizedConnections = append(normalizedConnections, c)
		}
	}
	layout.Connections = normalizedConnections
}

func adjustTruncated(layout *boxes.Layout, truncatedObjects map[string]TruncatedInfo) {
	copyTruncatedConnections(layout, truncatedObjects)
	copyTruncatedComments(layout, truncatedObjects)
	copyTruncatedTags(layout, truncatedObjects)
	adjustDestIdInRespectOfTruncated(layout, truncatedObjects)
	if len(layout.Horizontal) > 0 {
		adjustTruncatedForCont(&layout.Horizontal, truncatedObjects)
	}
	if len(layout.Vertical) > 0 {
		adjustTruncatedForCont(&layout.Vertical, truncatedObjects)
	}
}

func adjustTruncatedForCont(cont *[]boxes.Layout, truncatedObjects map[string]TruncatedInfo) {
	if cont != nil {
		for i := range len(*cont) {
			adjustTruncated(&(*cont)[i], truncatedObjects)
		}
	}
}

func FilterBoxes(layout boxes.Boxes, defaultDepth int, expanded, blacklisted []string) boxes.Boxes {
	filteredBoxes := boxes.CopyBoxes(&layout)
	b, truncatedObjects := truncBoxes(layout.Boxes, 0, defaultDepth, expanded, blacklisted)
	adjustTruncated(&b, truncatedObjects)
	filteredBoxes.Boxes = b
	adjustOverlaysWithTruncated(filteredBoxes, truncatedObjects)
	return *filteredBoxes
}

func adjustOverlaysWithTruncated(filteredBoxes *boxes.Boxes, truncatedObjects map[string]TruncatedInfo) {
	for i := range filteredBoxes.Overlays {
		o := &filteredBoxes.Overlays[i]
		newOverlays := make(map[string]float64, 0)
		for k, v := range o.Layouts {
			foundId := getTruncatedObject(k, truncatedObjects)
			if foundId == "" {
				// the layout elem still exist ...
				if existing, ok := newOverlays[k]; ok {
					v += existing
				}
				newOverlays[k] = v
			} else {
				if existing, ok := newOverlays[foundId]; ok {
					v += existing
				}
				newOverlays[foundId] = v
			}
		}
		o.Layouts = newOverlays
	}
}

func getTruncatedObject(text string, truncatedObjects map[string]TruncatedInfo) string {
	for _, v := range truncatedObjects {
		if v.truncated.Id == text || v.truncated.Caption == text {
			return v.newId
		}
	}
	return ""
}

func DrawBoxesFilteredExt(
	layout boxes.Boxes,
	mixins []boxes.BoxesFileMixings,
	defaultDepth int,
	expanded, blacklisted []string,
	debug bool) UIReturn {
	for _, m := range mixins {
		layout.MixinThings(m)
	}
	return DrawBoxesFiltered(layout, defaultDepth, expanded, blacklisted, debug)
}

func searchForExpanded(layout *boxes.Layout, expanded *[]string) {
	if layout.Id != "" && len(layout.Connections) > 0 {
		if !slices.Contains(*expanded, layout.Id) {
			*expanded = append(*expanded, layout.Id)
		}
		for i := range layout.Connections {
			c := layout.Connections[i]
			if c.DestId != "" && (!slices.Contains(*expanded, c.DestId)) {
				*expanded = append(*expanded, c.DestId)
			}
		}
	}
	searchForExpandedCont(layout.Horizontal, expanded)
	searchForExpandedCont(layout.Vertical, expanded)
}

func searchForExpandedCont(cont []boxes.Layout, expanded *[]string) {
	for i := range cont {
		searchForExpanded(&cont[i], expanded)
	}
}

func searchForBlacklisted(layout *boxes.Layout, expanded []string, blacklisted *[]string) {
	if !isRelatedToId(*layout, expanded) {
		*blacklisted = append(*blacklisted, layout.Id)
	} else {
		searchForBlacklistedCont(layout.Horizontal, expanded, blacklisted)
		searchForBlacklistedCont(layout.Vertical, expanded, blacklisted)
	}
}

func searchForBlacklistedCont(cont []boxes.Layout, expanded []string, blacklisted *[]string) {
	for i := range cont {
		searchForBlacklisted(&cont[i], expanded, blacklisted)
	}
}

func DrawBoxesRelatedToConnections(
	layout boxes.Boxes,
	mixins []boxes.BoxesFileMixings,
	debug bool) UIReturn {
	for _, m := range mixins {
		layout.MixinThings(m)
	}
	// find all IDs that have connections
	expanded := make([]string, 0)
	blacklisted := make([]string, 0)
	searchForExpanded(&layout.Boxes, &expanded)
	searchForBlacklisted(&layout.Boxes, expanded, &blacklisted)
	return DrawBoxesFiltered(layout, 2, expanded, blacklisted, debug)
}

func DrawBoxesFiltered(layout boxes.Boxes, defaultDepth int, expanded, blacklisted []string, debug bool) UIReturn {
	return DrawBoxesFilteredComments(layout, defaultDepth, expanded, blacklisted, false, debug)
}

func DrawBoxesFilteredComments(layout boxes.Boxes, defaultDepth int, expanded, blacklisted []string, hideComments, debug bool) UIReturn {
	textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()
	filteredLayout := FilterBoxes(layout, defaultDepth, expanded, blacklisted)
	doc, err := InitialLayoutBoxes(&filteredLayout, textDimensionCalulator)
	if err != nil {
		return UIReturn{ErrorMsg: fmt.Sprintf("error while initialy layout: %v", err)}
	}
	doc.ConnectBoxes()
	doc.AdjustDocHeightTLegend(textDimensionCalulator)
	if !hideComments {
		doc.IncludeComments(textDimensionCalulator)
	}
	doc.IncludeOverlays(textDimensionCalulator)
	var svgBuilder strings.Builder
	svgdrawing := svgdrawing.NewDrawing(&svgBuilder)
	svgdrawing.Start(doc.Title, doc.Height, doc.Width)
	svgdrawing.InitImages(doc.Images)
	if debug {
		svgdrawing.DrawRaster(doc.Width, doc.Height, types.RasterSize)
	}
	doc.DrawBoxes(svgdrawing)
	if debug {
		doc.InitStartPositions() // not needed to be called separately
		doc.InitRoads()          // not needed to be called separately
		doc.DrawRoads(svgdrawing)
		doc.DrawStartPositions(svgdrawing)
		doc.DrawConnectionNodes(svgdrawing)
		doc.Boxes.DrawTextBoxes(svgdrawing)
	}
	doc.DrawConnections(svgdrawing)
	doc.DrawOverlays(svgdrawing, textDimensionCalulator)
	doc.DrawTitle(svgdrawing, textDimensionCalulator)
	doc.DrawLegend(svgdrawing, textDimensionCalulator)
	doc.DrawComments(svgdrawing, textDimensionCalulator)
	svgdrawing.Done()
	return UIReturn{SVG: svgBuilder.String()}
}

func initContainerFromAlternativePath(extPathToLoad, origInputFile string) ([]boxes.Layout, error) {
	var pathToLoad string
	if path.IsAbs(extPathToLoad) {
		pathToLoad = extPathToLoad
	} else {
		// relative path, so we need to resolve it against the original input file
		origDir := filepath.Dir(origInputFile)
		pathToLoad = filepath.Join(origDir, extPathToLoad)
	}
	layoutArray, err := types.LoadInputFromFile[[]boxes.Layout](pathToLoad)
	if err != nil {
		return make([]boxes.Layout, 0), err
	}
	return *layoutArray, nil
}

func replaceAlternativePaths(l *boxes.Layout, inputFile string) error {
	if l == nil {
		return nil
	}
	if l.ExtHorizontal != nil {
		if cont, err := initContainerFromAlternativePath(*l.ExtHorizontal, inputFile); err != nil {
			return fmt.Errorf("error while initialize horizontal container from relative path: %w", err)
		} else {
			l.ExtHorizontal = nil // clear the alternative path
			l.Horizontal = cont
		}
	}
	if l.ExtVertical != nil {
		if cont, err := initContainerFromAlternativePath(*l.ExtVertical, inputFile); err != nil {
			return fmt.Errorf("error while initialize vertical container from relative path: %w", err)
		} else {
			l.ExtVertical = nil // clear the alternative path
			l.Vertical = cont
		}
	}
	if err := replaceAlternativePathsForContainer(l.Horizontal, inputFile); err != nil {
		return fmt.Errorf("error while replacing alternative paths for horizontal container: %w", err)
	}
	if err := replaceAlternativePathsForContainer(l.Vertical, inputFile); err != nil {
		return fmt.Errorf("error while replacing alternative paths for vertical container: %w", err)
	}
	return nil
}

func replaceAlternativePathsForContainer(cont []boxes.Layout, inputFile string) error {
	if cont == nil {
		return nil
	}
	for i, _ := range cont {
		if err := replaceAlternativePaths(&cont[i], inputFile); err != nil {
			return fmt.Errorf("error while replacing alternative paths for container: %w", err)
		}
	}
	return nil
}

func InitialLayoutBoxes(b *boxes.Boxes, c types.TextDimensionCalculator) (*boxes.BoxesDocument, error) {
	doc, err := DocumentFromBoxes(b)
	if err != nil {
		return nil, err
	}

	doc.Boxes.X = doc.GlobalPadding
	doc.Boxes.Y = doc.GlobalPadding

	doc.Boxes.InitDimensions(c)
	doc.Width = doc.Boxes.Width + doc.GlobalPadding*2
	doc.Height = doc.Boxes.Height + doc.GlobalPadding*2
	doc.Boxes.Center()
	doc.Height = doc.AdjustDocHeight(&doc.Boxes, 0) + (2 * doc.GlobalPadding)

	if doc.Title != "" {
		format2Use := doc.GetTitleFormat()
		w, h := c.Dimensions(doc.Title, &format2Use)
		doc.Height += h + (2 * doc.GlobalPadding)
		if w > doc.Width {
			doc.Width = w + (2 * doc.GlobalPadding)
		}
	}

	return doc, nil
}

func LayoutBoxesWithFilter(b *boxes.Boxes, c types.TextDimensionCalculator, defaultDepth int) (*boxes.BoxesDocument, error) {
	doc, err := DocumentFromBoxes(b)
	if err != nil {
		return nil, err
	}
	doc.Boxes.X = doc.GlobalPadding
	doc.Boxes.Y = doc.GlobalPadding
	doc.Boxes.InitDimensions(c)
	doc.Width = doc.Boxes.Width + doc.GlobalPadding*2
	doc.Height = doc.Boxes.Height + doc.GlobalPadding*2
	doc.Boxes.Center()
	doc.Height = doc.AdjustDocHeight(&doc.Boxes, 0) + (2 * doc.GlobalPadding)

	if doc.Title != "" {
		format2Use := doc.GetTitleFormat()
		w, h := c.Dimensions(doc.Title, &format2Use)
		doc.Height += h + (2 * doc.GlobalPadding)
		if w > doc.Width {
			doc.Width = w + (2 * doc.GlobalPadding)
		}
	}
	return doc, nil
}
