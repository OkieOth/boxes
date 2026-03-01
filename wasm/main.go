//go:build js && wasm
// +build js,wasm

package main

import (
	"errors"
	"fmt"
	"syscall/js"

	"github.com/okieoth/draw.chart.things/pkg/boxesimpl"
	"github.com/okieoth/draw.chart.things/pkg/types/boxes"

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

func main() {
	// Expose the function to JS as `getSvg`
	js.Global().Set("createSvg", js.FuncOf(createSvgWrapper))
	js.Global().Set("createSvgExt", js.FuncOf(createSvgExtWrapper))
	js.Global().Set("createSvgForConnected", js.FuncOf(createSvgForConnectedWrapper))

	// Keep Go running
	select {}
}
