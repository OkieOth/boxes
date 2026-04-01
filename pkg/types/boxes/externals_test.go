package boxes_test

import (
	"testing"

	"github.com/okieoth/boxes/pkg/types"
	"github.com/okieoth/boxes/pkg/types/boxes"
	"github.com/stretchr/testify/require"
)

func buildBoxesWithIds() boxes.Boxes {
	return boxes.Boxes{
		Boxes: boxes.Layout{
			Id: "root",
			Horizontal: []boxes.Layout{
				{Id: "box_a", Caption: "Box A"},
				{Id: "box_b", Caption: "Box B"},
				{Id: "box_c", Caption: "Box C"},
			},
		},
	}
}

func buildMixinWithSteps() boxes.BoxesFileMixings {
	destC := "box_c"
	commentText1 := "step 1 comment"
	commentText2 := "step 2 comment"
	return boxes.BoxesFileMixings{
		Steps: []boxes.ProcessStep{
			{
				Caption: "Step 1",
				Connections: map[string]boxes.ConnectionCont{
					"box_a": {Connections: []boxes.Connection{{DestId: destC}}},
				},
				Comments: map[string]types.Comment{
					"box_b": {Text: commentText1},
				},
			},
			{
				Caption: "Step 2",
				Connections: map[string]boxes.ConnectionCont{
					"box_b": {Connections: []boxes.Connection{{DestId: destC}}},
				},
				Comments: map[string]types.Comment{
					"box_b": {Text: commentText2},
				},
			},
		},
	}
}

func TestMixinThingsWithSteps_NoActiveSteps(t *testing.T) {
	b := buildBoxesWithIds()
	m := buildMixinWithSteps()

	b.MixinThingsWithSteps(m, []int{})

	require.Len(t, b.Boxes.Horizontal[0].Connections, 0, "box_a should have no connections")
	require.Len(t, b.Boxes.Horizontal[1].Connections, 0, "box_b should have no connections")
	require.Nil(t, b.Boxes.Horizontal[1].Comment, "box_b should have no comment")
}

func TestMixinThingsWithSteps_Step0Active(t *testing.T) {
	b := buildBoxesWithIds()
	m := buildMixinWithSteps()

	b.MixinThingsWithSteps(m, []int{0})

	require.Len(t, b.Boxes.Horizontal[0].Connections, 1, "box_a should have 1 connection from step 0")
	require.Equal(t, "box_c", b.Boxes.Horizontal[0].Connections[0].DestId)
	require.Len(t, b.Boxes.Horizontal[1].Connections, 0, "box_b should have no connections")
	require.NotNil(t, b.Boxes.Horizontal[1].Comment, "box_b should have comment from step 0")
	require.Equal(t, "step 1 comment", b.Boxes.Horizontal[1].Comment.Text)
}

func TestMixinThingsWithSteps_Step1Active(t *testing.T) {
	b := buildBoxesWithIds()
	m := buildMixinWithSteps()

	b.MixinThingsWithSteps(m, []int{1})

	require.Len(t, b.Boxes.Horizontal[0].Connections, 0, "box_a should have no connections")
	require.Len(t, b.Boxes.Horizontal[1].Connections, 1, "box_b should have 1 connection from step 1")
	require.Equal(t, "box_c", b.Boxes.Horizontal[1].Connections[0].DestId)
	require.NotNil(t, b.Boxes.Horizontal[1].Comment, "box_b should have comment from step 1")
	require.Equal(t, "step 2 comment", b.Boxes.Horizontal[1].Comment.Text)
}

func TestMixinThingsWithSteps_BothStepsActive(t *testing.T) {
	b := buildBoxesWithIds()
	m := buildMixinWithSteps()

	b.MixinThingsWithSteps(m, []int{0, 1})

	// connections from both steps
	require.Len(t, b.Boxes.Horizontal[0].Connections, 1, "box_a should have 1 connection from step 0")
	require.Len(t, b.Boxes.Horizontal[1].Connections, 1, "box_b should have 1 connection from step 1")
	// step 1 comment overwrites step 0 comment (last writer wins)
	require.NotNil(t, b.Boxes.Horizontal[1].Comment)
	require.Equal(t, "step 2 comment", b.Boxes.Horizontal[1].Comment.Text)
}

func TestMixinThingsWithSteps_OutOfBoundsIndex(t *testing.T) {
	b := buildBoxesWithIds()
	m := buildMixinWithSteps()

	b.MixinThingsWithSteps(m, []int{0, 99})

	// only step 0 applied, out-of-bounds index silently ignored
	require.Len(t, b.Boxes.Horizontal[0].Connections, 1)
	require.Len(t, b.Boxes.Horizontal[1].Connections, 0)
}

func TestMixinThingsWithSteps_RootConnectionsAlwaysApplied(t *testing.T) {
	b := buildBoxesWithIds()
	m := boxes.BoxesFileMixings{
		Connections: map[string]boxes.ConnectionCont{
			"box_a": {Connections: []boxes.Connection{{DestId: "box_b"}}},
		},
		Steps: []boxes.ProcessStep{
			{
				Caption: "Step 1",
				Connections: map[string]boxes.ConnectionCont{
					"box_b": {Connections: []boxes.Connection{{DestId: "box_c"}}},
				},
			},
		},
	}

	// activate no steps — root connection should still be applied
	b.MixinThingsWithSteps(m, []int{})

	require.Len(t, b.Boxes.Horizontal[0].Connections, 1, "root connection on box_a always applied")
	require.Len(t, b.Boxes.Horizontal[1].Connections, 0, "step 0 connection on box_b not applied")
}

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
