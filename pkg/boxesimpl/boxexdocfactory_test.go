package boxesimpl_test

import (
	"testing"

	"github.com/okieoth/boxes/pkg/boxesimpl"
	"github.com/okieoth/boxes/pkg/types"
	"github.com/okieoth/boxes/pkg/types/boxes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var TAG1 string = "tag1"

func TestInitLayoutElement(t *testing.T) {
	bf2 := boxes.BoxFormat{
		Padding: 0,
		FontCaption: types.FontDef{
			Size:       12,
			Type:       &boxesimpl.ExpFontDefTypeEnum_normal,
			Weight:     &boxesimpl.ExpFontDefWeightEnum_bold,
			LineHeight: 2.0,
			Color:      "red",
			Aligned:    &boxesimpl.ExpFontDefAlignedEnum_center,
		},
		FontText1: types.FontDef{
			Size:       10,
			Type:       &boxesimpl.ExpFontDefTypeEnum_normal,
			Weight:     &boxesimpl.ExpFontDefWeightEnum_normal,
			LineHeight: 1.5,
			Color:      "black",
			Aligned:    &boxesimpl.ExpFontDefAlignedEnum_left,
		},
		FontText2: types.FontDef{
			Size:       10,
			Type:       &boxesimpl.ExpFontDefTypeEnum_normal,
			Weight:     &boxesimpl.ExpFontDefWeightEnum_normal,
			LineHeight: 1.5,
			Color:      "black",
			Aligned:    &boxesimpl.ExpFontDefAlignedEnum_left,
		},
	}
	tests := []struct {
		name         string
		layout       boxes.Layout
		inputFormats map[string]boxes.BoxFormat
		expected     boxes.LayoutElement
	}{
		{
			name: "Test with empty layout and formats",
			layout: boxes.Layout{
				Id: "test1",
			},
			inputFormats: map[string]boxes.BoxFormat{},
			expected: boxes.LayoutElement{
				Id:     "test1",
				Format: nil,
			},
		},
		{
			name: "Test with layout and formats",
			layout: boxes.Layout{
				Id:      "test2",
				Caption: "Test Caption",
				Text1:   "Test Text1",
				Text2:   "Test Text2",
				Format:  &TAG1,
			},
			inputFormats: map[string]boxes.BoxFormat{
				"tag1": {
					FontCaption: types.FontDef{
						Size:       12,
						Type:       &boxesimpl.ExpFontDefTypeEnum_normal,
						Weight:     &boxesimpl.ExpFontDefWeightEnum_bold,
						LineHeight: 2.0,
						Color:      "red",
						Aligned:    &boxesimpl.ExpFontDefAlignedEnum_center,
					},
				},
			},
			expected: boxes.LayoutElement{
				Id:      "test2",
				Caption: "Test Caption",
				Text1:   "Test Text1",
				Text2:   "Test Text2",
				Format:  &bf2,
			},
		},
	}

	b := boxes.NewBoxes()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := boxes.NewBoxesDocument()
			doc.Formats = tt.inputFormats
			result := boxesimpl.ExpInitLayoutElement(&tt.layout, doc, b)
			if result.Id != tt.expected.Id {
				t.Errorf("expected Id %v, got %v", tt.expected.Id, result.Id)
			}
			if result.Caption != tt.expected.Caption {
				t.Errorf("expected Caption %v, got %v", tt.expected.Caption, result.Caption)
			}
			if result.Text1 != tt.expected.Text1 {
				t.Errorf("expected Text1 %v, got %v", tt.expected.Text1, result.Text1)
			}
			if result.Text2 != tt.expected.Text2 {
				t.Errorf("expected Text2 %v, got %v", tt.expected.Text2, result.Text2)
			}
			if (result.Format == nil && tt.expected.Format != nil) || (result.Format != nil && tt.expected.Format == nil) {
				t.Errorf("expected Format %v, got %v", tt.expected.Format, result.Format)
			}
			if result.Format == nil {
				return
			}
			if result.Format.Padding != tt.expected.Format.Padding {
				t.Errorf("expected Padding %v, got %v", tt.expected.Format.Padding, result.Format.Padding)
			}
			if result.Format.FontCaption.Size != tt.expected.Format.FontCaption.Size {
				t.Errorf("expected FontCaption Size %v, got %v", tt.expected.Format.FontCaption.Size, result.Format.FontCaption.Size)
			}
			if *result.Format.FontCaption.Type != *tt.expected.Format.FontCaption.Type {
				t.Errorf("expected FontCaption Type %v, got %v", tt.expected.Format.FontCaption.Type, result.Format.FontCaption.Type)
			}
			if *result.Format.FontCaption.Weight != *tt.expected.Format.FontCaption.Weight {
				t.Errorf("expected FontCaption Weight %v, got %v", tt.expected.Format.FontCaption.Weight, result.Format.FontCaption.Weight)
			}
			if result.Format.FontCaption.LineHeight != tt.expected.Format.FontCaption.LineHeight {
				t.Errorf("expected FontCaption LineHeight %v, got %v", tt.expected.Format.FontCaption.LineHeight, result.Format.FontCaption.LineHeight)
			}
			if result.Format.FontCaption.Color != tt.expected.Format.FontCaption.Color {
				t.Errorf("expected FontCaption Color %v, got %v", tt.expected.Format.FontCaption.Color, result.Format.FontCaption.Color)
			}
			if *result.Format.FontCaption.Aligned != *tt.expected.Format.FontCaption.Aligned {
				t.Errorf("expected FontCaption Aligned %v, got %v", tt.expected.Format.FontCaption.Aligned, result.Format.FontCaption.Aligned)
			}
		})
	}
}

func TestLoadFromFile(t *testing.T) {
	tests := []struct {
		inputFile string
	}{
		{
			inputFile: "../../resources/examples_boxes/complex_horizontal_connected.yaml",
		},
	}

	for _, test := range tests {
		b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
		require.Nil(t, err)
		require.NotNil(t, b)
		require.Equal(t, 0, len(b.Boxes.Connections))
		require.Equal(t, "r4_1", b.Boxes.Horizontal[0].Vertical[0].Id)
		require.Equal(t, "r5_3", b.Boxes.Horizontal[0].Vertical[0].Connections[0].DestId)
		require.Equal(t, "r5_1", b.Boxes.Horizontal[1].Vertical[0].Id)
		require.Equal(t, 0, len(b.Boxes.Horizontal[1].Vertical[0].Connections))
		require.Equal(t, "r5_2", b.Boxes.Horizontal[1].Vertical[1].Id)
		require.Equal(t, 0, len(b.Boxes.Horizontal[1].Vertical[1].Connections))
		require.Equal(t, "r5_3", b.Boxes.Horizontal[1].Vertical[2].Id)
		require.Equal(t, 0, len(b.Boxes.Horizontal[1].Vertical[2].Connections))

		require.Equal(t, "r6_1", b.Boxes.Horizontal[2].Vertical[0].Id)
		require.Equal(t, 0, len(b.Boxes.Horizontal[2].Vertical[0].Connections))
		require.Equal(t, "r6_2", b.Boxes.Horizontal[2].Vertical[1].Id)
		require.Equal(t, 0, len(b.Boxes.Horizontal[2].Vertical[1].Connections))
		require.Equal(t, "r6_3", b.Boxes.Horizontal[2].Vertical[2].Id)
		require.Equal(t, 0, len(b.Boxes.Horizontal[2].Vertical[2].Connections))
		require.Equal(t, "r6_4", b.Boxes.Horizontal[2].Vertical[3].Id)
		require.Equal(t, 1, len(b.Boxes.Horizontal[2].Vertical[3].Connections))
		require.Equal(t, "r4_1", b.Boxes.Horizontal[2].Vertical[3].Connections[0].DestId)
		require.Equal(t, "r6_5", b.Boxes.Horizontal[2].Vertical[4].Id)
		require.Equal(t, 1, len(b.Boxes.Horizontal[2].Vertical[4].Connections))

	}
}

func TestBuildFromFile(t *testing.T) {
	tests := []struct {
		inputFile string
	}{
		{
			inputFile: "../../resources/examples_boxes/complex_horizontal_connected_pics.yaml",
		},
	}

	for _, test := range tests {
		b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
		require.Nil(t, err)
		require.NotNil(t, b)
		doc, err := boxesimpl.DocumentFromBoxes(b)
		require.Nil(t, err)
		require.NotNil(t, doc)
		assert.Len(t, doc.Images, 3)
		// TODO
	}
}
