package boxes_test

import (
	"testing"

	"github.com/okieoth/draw.chart.things/pkg/types"
	"github.com/okieoth/draw.chart.things/pkg/types/boxes"
	"github.com/stretchr/testify/require"
)

func TestLoadExternalFormats(t *testing.T) {
	inputFormat := "../../../resources/examples_boxes/ext_formats.yaml"
	inputLayout := "../../../resources/examples_boxes/ext_complex_horizontal_connected_pics.yaml"
	inputLayout2 := "../../../resources/examples_boxes/complex_horizontal_connected_pics.yaml"

	additionalFormats, err := types.LoadInputFromFile[boxes.BoxesFileMixings](inputFormat)
	require.Nil(t, err)
	require.NotNil(t, additionalFormats)

	require.Len(t, additionalFormats.Formats, 1)
	require.Len(t, additionalFormats.Images, 1)

	img, ok := additionalFormats.Images["smilie_01_43"]
	require.True(t, ok)

	b, err := types.LoadInputFromFile[boxes.Boxes](inputLayout)
	require.Nil(t, err)
	require.NotNil(t, b)

	require.NotNil(t, img.Base64)
	require.Nil(t, img.Base64Src)

	b2, err := types.LoadInputFromFile[boxes.Boxes](inputLayout2)
	require.Nil(t, err)
	require.NotNil(t, b2)

	require.NotEqual(t, len(b2.Images), len(b.Formats))

	b.MixinThings(*additionalFormats)

	require.Equal(t, len(b2.Images), len(b.Images))
}

func TestLoadExternalConnections(t *testing.T) {
	input := "../../../resources/examples_boxes/ext_complex_horizontal_connected_pics.yaml"
	inputConnections := "../../../resources/examples_boxes/ext_connections.yaml"

	b, err := types.LoadInputFromFile[boxes.Boxes](input)
	require.Nil(t, err)
	require.NotNil(t, b)

	c, err := types.LoadInputFromFile[boxes.BoxesFileMixings](inputConnections)
	require.Nil(t, err)
	require.NotNil(t, c)

	// r4_1
	require.Len(t, b.Boxes.Horizontal[0].Vertical[0].Connections, 2)
	// r5_2
	require.Len(t, b.Boxes.Horizontal[1].Vertical[1].Connections, 0)

	b.MixinThings(*c)

	// r4_1
	require.Len(t, b.Boxes.Horizontal[0].Vertical[0].Connections, 4)
	// r5_2
	require.Len(t, b.Boxes.Horizontal[1].Vertical[1].Connections, 1)

	cl, ok := b.Formats["connLines"]
	require.True(t, ok)
	require.NotNil(t, cl.Line)
	require.Equal(t, 2.0, *cl.Line.Width)
	require.Equal(t, "pink", *cl.Line.Color)
}

func TestLoadExternalConnections2(t *testing.T) {
	input := "../../../resources/examples_boxes/ext_complex_horizontal_connected_pics.yaml"
	inputConnections := "../../../resources/examples_boxes/ext_connections2.yaml"

	b, err := types.LoadInputFromFile[boxes.Boxes](input)
	require.Nil(t, err)
	require.NotNil(t, b)

	c, err := types.LoadInputFromFile[boxes.BoxesFileMixings](inputConnections)
	require.Nil(t, err)
	require.NotNil(t, c)

	// r4_1
	require.Len(t, b.Boxes.Horizontal[0].Vertical[0].Connections, 2)
	// r5_2
	require.Len(t, b.Boxes.Horizontal[1].Vertical[1].Connections, 0)

	require.Len(t, b.Boxes.Horizontal[1].Vertical[1].Vertical, 0)
	require.Len(t, b.Boxes.Horizontal[1].Vertical[1].Horizontal, 0)
	require.Len(t, b.Boxes.Horizontal[2].Vertical[0].Vertical, 0)
	require.Len(t, b.Boxes.Horizontal[2].Vertical[0].Horizontal, 0)

	b.MixinThings(*c)

	// r4_1
	require.Len(t, b.Boxes.Horizontal[0].Vertical[0].Connections, 4)
	// r5_2
	require.Len(t, b.Boxes.Horizontal[1].Vertical[1].Vertical, 2)
	require.Len(t, b.Boxes.Horizontal[1].Vertical[1].Horizontal, 0)
	require.Len(t, b.Boxes.Horizontal[2].Vertical[0].Vertical, 0)
	require.Len(t, b.Boxes.Horizontal[2].Vertical[0].Horizontal, 2)
}

func TestLoadExternalComments(t *testing.T) {
	input := "../../../resources/examples_boxes/ext_complex_horizontal_connected_pics.yaml"
	inputConnections := "../../../resources/examples_boxes/ext_comments.yaml"

	b, err := types.LoadInputFromFile[boxes.Boxes](input)
	require.Nil(t, err)
	require.NotNil(t, b)

	require.Len(t, b.Legend.Entries, 1, "wrong number of initial legend entries")

	c, err := types.LoadInputFromFile[boxes.BoxesFileMixings](inputConnections)
	require.Nil(t, err)
	require.NotNil(t, c)

	// r5_3
	require.Nil(t, b.Boxes.Horizontal[1].Vertical[0].Comment)
	require.Nil(t, b.Boxes.Horizontal[1].Vertical[1].Comment)
	require.Nil(t, b.Boxes.Horizontal[1].Vertical[2].Comment)

	// Most Left Element
	require.Nil(t, b.Boxes.Horizontal[2].Vertical[0].Comment)
	require.Nil(t, b.Boxes.Horizontal[2].Vertical[1].Comment)
	require.Nil(t, b.Boxes.Horizontal[2].Vertical[2].Comment)

	b.MixinThings(*c)

	require.Len(t, b.Legend.Entries, 2, "wrong number of legend entries")

	// r5_3
	require.Nil(t, b.Boxes.Horizontal[1].Vertical[0].Comment)
	require.Nil(t, b.Boxes.Horizontal[1].Vertical[1].Comment)
	require.NotNil(t, b.Boxes.Horizontal[1].Vertical[2].Comment)
	require.Equal(t, "I am a comment", b.Boxes.Horizontal[1].Vertical[2].Comment.Text)
	require.Nil(t, b.Boxes.Horizontal[1].Vertical[2].Comment.Label)
	require.Nil(t, b.Boxes.Horizontal[1].Vertical[2].Comment.Format)

	// Most Left Element
	require.NotNil(t, b.Boxes.Horizontal[2].Vertical[0].Comment)
	require.Equal(t, "Just another comment", b.Boxes.Horizontal[2].Vertical[0].Comment.Text)
	require.Equal(t, "yyy", *b.Boxes.Horizontal[2].Vertical[0].Comment.Format)
	require.Equal(t, "X", *b.Boxes.Horizontal[2].Vertical[0].Comment.Label)
	require.Nil(t, b.Boxes.Horizontal[2].Vertical[1].Comment)
	require.Nil(t, b.Boxes.Horizontal[2].Vertical[2].Comment)

}
