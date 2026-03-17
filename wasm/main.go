//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"syscall/js"

	"github.com/okieoth/boxes/pkg/boxesimpl"
	"github.com/okieoth/boxes/pkg/types"
	"github.com/okieoth/boxes/pkg/types/boxes"

	"gopkg.in/yaml.v3"
	y "gopkg.in/yaml.v3"
)

const unknownSvg string = `<svg xmlns="http://www.w3.org/2000/svg" width="800" height="400" viewBox="0 0 800 400">
	<rect x="150" y="100" width="500" height="200" fill="#eeeeee" stroke="#555555" stroke-width="2" />
	<text x="400" y="200" font-family="sans-serif" font-size="20" fill="#333333" text-anchor="middle" dominant-baseline="middle">unknown content</text>
</svg>
`

// getSVG implements the requested signature
func createSvg(boxesYaml string, defaultDepth int, expanded, blacklisted []string, hideComments, debug bool) string {
	fmt.Printf("createSvg: hideComments=%v, debug:=%v\n", hideComments, debug)
	var boxes boxes.Boxes
	if err := y.Unmarshal([]byte(boxesYaml), &boxes); err != nil {
		return unknownSvg
	}
	if boxes.Version != nil {
		boxes.Title += fmt.Sprintf(" [%s]", *boxes.Version)
	}
	ret := boxesimpl.DrawBoxesFilteredComments(boxes, defaultDepth, expanded, blacklisted, hideComments, debug)
	if ret.ErrorMsg != "" {
		fmt.Println(ret.ErrorMsg)
		return unknownSvg
	}
	return ret.SVG
}

func createSvgExt(boxesYaml string, mixins []string, defaultDepth int, expanded, blacklisted []string, hideComments, debug bool) string {
	var b boxes.Boxes
	fmt.Printf("createSvgExt: hideComments=%v, debug=%v\n", hideComments, debug)

	if err := y.Unmarshal([]byte(boxesYaml), &b); err != nil {
		fmt.Printf("error while unmarshalling boxes layout: %v", err)
		return unknownSvg
	}
	if b.Version != nil {
		b.Title += fmt.Sprintf(" [%s]", *b.Version)
	}
	for i, c := range mixins {
		var m boxes.BoxesFileMixings
		if err := y.Unmarshal([]byte(c), &m); err != nil {
			fmt.Printf("error while unmarshalling external connections (%d): %v", i, err)
			return unknownSvg
		}
		b.MixinThings(m)
	}

	ret := boxesimpl.DrawBoxesFilteredComments(b, defaultDepth, expanded, blacklisted, hideComments, debug)
	if ret.ErrorMsg != "" {
		fmt.Println(ret.ErrorMsg)
		return unknownSvg
	}
	return ret.SVG
}

func getArrayFromJsValue(args []js.Value, index int) ([]string, error) {
	jsArray := args[index]
	if jsArray.Type() != js.TypeObject {
		return []string{}, errors.New("")
	}
	length := jsArray.Length()
	ret := make([]string, 0, length)
	for i := range length {
		val := jsArray.Index(i)
		if val.Type() == js.TypeString {
			ret = append(ret, val.String())
		}
	}
	return ret, nil
}

// JS wrapper: exposes getSvg to JavaScript
func createSvgWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 6 {
		return "error: expected (string, number, string[], string[], bool, bool)"
	}
	input := args[0].String()
	depth := args[1].Int()
	hideComments := args[4].Bool()
	debug := args[5].Bool()
	expanded, err := getArrayFromJsValue(args, 2)
	if err != nil {
		return "error: expanded must be an array"
	}
	blacklisted, err := getArrayFromJsValue(args, 3)
	if err != nil {
		return "error: blacklisted must be an array"
	}
	return createSvg(input, depth, expanded, blacklisted, hideComments, debug)
}

func createSvgExtWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 7 {
		return "error: expected (string, string[], number, string[], string[], bool, bool)"
	}
	input := args[0].String()
	mixins, err := getArrayFromJsValue(args, 1)
	if err != nil {
		return "error: mixins need to be an array"
	}
	expanded, err := getArrayFromJsValue(args, 3)
	if err != nil {
		return "error: expanded must be an array"
	}
	blacklisted, err := getArrayFromJsValue(args, 4)
	if err != nil {
		return "error: blacklisted must be an array"
	}
	depth := args[2].Int()
	hideComments := args[5].Bool()
	debug := args[6].Bool()
	return createSvgExt(input, mixins, depth, expanded, blacklisted, hideComments, debug)
}

type NewReturn struct {
	SVG         string
	Expanded    []string
	Blacklisted []string
}

func returnErrorSvg() NewReturn {
	return NewReturn{
		SVG:         unknownSvg,
		Expanded:    []string{},
		Blacklisted: []string{},
	}
}

func createSvgForConnected(boxesYaml string, mixins []string, debug bool) NewReturn {
	var b boxes.Boxes
	fmt.Printf("createSvgForConnected: debug=%v\n", debug)

	if err := y.Unmarshal([]byte(boxesYaml), &b); err != nil {
		fmt.Printf("error while unmarshalling boxes layout: %v", err)
		return returnErrorSvg()
	}
	if b.Version != nil {
		b.Title += fmt.Sprintf(" [%s]", *b.Version)
	}
	boxMixins := make([]boxes.BoxesFileMixings, 0)
	for i, c := range mixins {
		var m boxes.BoxesFileMixings
		if err := y.Unmarshal([]byte(c), &m); err != nil {
			fmt.Printf("error while unmarshalling external connections (%d): %v\n", i, err)
		}
		boxMixins = append(boxMixins, m)
	}

	ret := boxesimpl.DrawBoxesRelatedToConnections(b, boxMixins, debug)
	if ret.ErrorMsg != "" {
		fmt.Println(ret.ErrorMsg)
		return returnErrorSvg()
	}
	return NewReturn{
		SVG:         ret.SVG,
		Expanded:    ret.Expanded,
		Blacklisted: ret.Blacklisted,
	}
}
func sliceToJSArray(items []string) js.Value {
	arr := js.Global().Get("Array").New()
	for _, v := range items {
		arr.Call("push", v)
	}
	return arr
}

func createSvgForConnectedWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return "error: expected (string, string[], bool)"
	}
	input := args[0].String()
	mixins, err := getArrayFromJsValue(args, 1)
	if err != nil {
		return "error: mixins need to be an array"
	}
	debug := args[2].Bool()
	result := createSvgForConnected(input, mixins, debug)

	// Create JS object
	obj := js.Global().Get("Object").New()

	obj.Set("svg", result.SVG)
	obj.Set("expanded", sliceToJSArray(result.Expanded))
	obj.Set("blacklisted", sliceToJSArray(result.Blacklisted))

	return obj
}

type searchItem struct {
	Id      string `json:"id"`
	Caption string `json:"caption"`
}

func collectSearchItems(layout boxes.Layout, items *[]searchItem) {
	if layout.Id != "" && layout.Caption != "" {
		*items = append(*items, searchItem{Id: layout.Id, Caption: layout.Caption})
	}
	for _, child := range layout.Vertical {
		collectSearchItems(child, items)
	}
	for _, child := range layout.Horizontal {
		collectSearchItems(child, items)
	}
}

func getSearchableItems(boxesYaml string, mixins []string) string {
	var b boxes.Boxes
	if err := y.Unmarshal([]byte(boxesYaml), &b); err != nil {
		fmt.Printf("getSearchableItems: error unmarshalling input: %v\n", err)
		return "[]"
	}
	for i, c := range mixins {
		var m boxes.BoxesFileMixings
		if err := y.Unmarshal([]byte(c), &m); err != nil {
			fmt.Printf("getSearchableItems: error unmarshalling mixin (%d): %v\n", i, err)
		}
		b.MixinThings(m)
	}
	items := make([]searchItem, 0)
	collectSearchItems(b.Boxes, &items)
	jsonBytes, err := json.Marshal(items)
	if err != nil {
		return "[]"
	}
	return string(jsonBytes)
}

func getSearchableItemsWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return "[]"
	}
	input := args[0].String()
	mixins, err := getArrayFromJsValue(args, 1)
	if err != nil {
		return "[]"
	}
	return getSearchableItems(input, mixins)
}

// Function returns a mixin to highlight the the selected items
func getSearchMixin(selectedIds []string) string {
	if len(selectedIds) == 0 {
		return ""
	}
	const searchResultColor = "#f50057"
	mixin := boxes.NewBoxesFileMixings()
	format := boxes.Format{
		Line: types.InitLineDef2(searchResultColor, 1),
		Fill: types.InitFillDef2(searchResultColor, 1),
	}
	mixin.Legend = boxes.NewLegend()
	mixin.Legend.Entries = append(mixin.Legend.Entries, boxes.LegendEntry{
		Text:   "Search Result",
		Format: "searchResult",
	})
	mixin.Formats["searchResult"] = format
	overlay := boxes.NewOverlay()
	overlay.Caption = "Search Result"
	overlay.CenterXOffset = -1
	overlay.CenterYOffset = -1
	overlay.Formats = boxes.NewOverlayFormatDef()
	overlay.Formats.Default = "searchResult"
	for _, id := range selectedIds {
		overlay.Layouts[id] = 5
	}
	mixin.Overlays = append(mixin.Overlays, *overlay)
	bytes, err := yaml.Marshal(mixin)
	if err != nil {
		fmt.Println("error while marshal search mixin:", *mixin)
		return ""

	}
	return string(bytes)
}

func getSearchMixinWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return ""
	}
	selectedIds, err := getArrayFromJsValue(args, 0)
	if err != nil {
		return ""
	}
	return getSearchMixin(selectedIds)
}

func main() {
	// Expose the function to JS as `getSvg`
	js.Global().Set("createSvg", js.FuncOf(createSvgWrapper))
	js.Global().Set("createSvgExt", js.FuncOf(createSvgExtWrapper))
	js.Global().Set("createSvgForConnected", js.FuncOf(createSvgForConnectedWrapper))
	js.Global().Set("getSearchableItems", js.FuncOf(getSearchableItemsWrapper))
	js.Global().Set("getSearchMixin", js.FuncOf(getSearchMixinWrapper))

	// Keep Go running
	select {}
}
