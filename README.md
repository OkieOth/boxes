[![CI](https://github.com/OkieOth/draw.chart.things/actions/workflows/test.yml/badge.svg?branch=main&event=push)](https://github.com/OkieOth/draw.chart.things/actions/workflows/test.yml)

# boxes

> Diagrams as code. Boxes, connections, layouts — all from YAML.

A Go tool and library for rendering **hierarchical box diagrams** as SVG from plain YAML specifications. No GUI drag-and-drop nonsense. You describe the structure, the tool figures out the layout and routes the connections. Ship it as a CLI binary or run it in the browser via WebAssembly.

---

## What it does

You write YAML. You get SVG. Between those two things, the tool:

- **Auto-lays out** nested boxes arranged in rows and columns
- **Routes connections** between boxes using Dijkstra's algorithm, avoiding collisions
- **Applies styles** via named format definitions (fonts, colors, fill, line styles, arrows)
- **Truncates deep hierarchies** at a configurable depth, with optional expansion of specific nodes
- **Merges mixins** — inject connections, formats, or extra boxes from separate YAML files
- **Renders overlays** — heatmaps and visual layers on top of the diagram
- **Aggregates comments** per box and renders them as a side panel
- **Exports SVG** to file or stdout, or serves live previews in the browser

---

## Core concepts

### Boxes

A `boxes` tree of nested containers, arranged either `horizontal`ly or `vertical`ly. Each box has an `id`, optional `caption`, `text1`/`text2` fields, a `format` name, and optional `connections`.

```yaml
title: My System
boxes:
  vertical:
    - id: frontend
      caption: Frontend
      format: b_ui
      horizontal:
        - id: spa
          caption: Single Page App
          format: s_ui
        - id: mobile
          caption: Mobile Client
          format: s_ui
      connections:
        - destId: api_gateway
          format: http
    - id: api_gateway
      caption: API Gateway
      format: b_infra
```

### Formats

Named style definitions. Define once, reference everywhere. Supports fonts, colors, fill, borders, line styles, and arrow types. Can be conditionally applied via tags.

```yaml
formats:
  b_ui:
    fill:
      color: "#dbeafe"
    border:
      color: "#3b82f6"
      width: 2
  http:
    line:
      color: "#6366f1"
      width: 2
    destArrow: true
```

### Mixins

External YAML files that inject connections, formats, or additional layout nodes into a base diagram. Useful for layering concerns — e.g., a base architecture diagram + a separate alarm-flow mixin.

```yaml
# ext_alarms.yaml
connections:
  some_service:
    connections:
      - dest: broker
        format: amqp
formats:
  amqp:
    line:
      color: orange
```

### Depth & filtering

Large diagrams can be truncated at a given depth. Specific subtrees can be pinned open (`expand: true` in YAML or `--expand id` on the CLI). Nodes can be hidden entirely with `--blacklisted`.

---

## Architecture

```
YAML input
  └─► load          (pkg/types/load.go)
  └─► internal doc  (boxesdocfactory.go)
  └─► layout        (pkg/types/boxes/layout.go)
  └─► routing       (pkg/types/boxes/dijkstra.go, connect.go)
  └─► SVG render    (pkg/svgdrawing/, pkg/svg/)
  └─► output SVG
```

The type system is **model-driven**: Go structs in `pkg/types/` are generated from JSON schemas in `configs/models/` via [YACG](https://github.com/OkieOth/yacg). Hand-written implementation lives in separate files, never touched by codegen.

---

## CLI usage

```bash
# Build
make build          # outputs to build/draw

# Draw a diagram
./build/draw boxes \
  --from my_diagram.yaml \
  --output diagram.svg

# With mixins, depth control, and expansion
./build/draw boxes \
  --from architecture.yaml \
  --mixin ext_connections.yaml \
  --mixin ext_formats.yaml \
  --depth 3 \
  --expand id_some_deep_node \
  --blacklisted id_legacy_subsystem \
  --output architecture.svg

# Debug mode: renders routing roads and connection nodes
./build/draw boxes --from diagram.yaml --debug --output debug.svg

# Output to stdout (pipe into your toolchain)
./build/draw boxes --from diagram.yaml | rsvg-convert -f pdf > diagram.pdf
```

### All flags for `boxes`

| Flag | Short | Default | Description |
|---|---|---|---|
| `--from` | `-f` | *(required)* | Path to the input YAML file |
| `--output` | `-o` | stdout | Output SVG file path |
| `--mixin` | `-m` | — | Mixin YAML file path (repeatable) |
| `--depth` | `-d` | `2` | Default truncation depth |
| `--expand` | `-e` | — | Box ID to force-expand (repeatable) |
| `--blacklisted` | `-b` | — | Box ID to hide entirely (repeatable) |
| `--debug` | | `false` | Render routing roads and graph nodes |

---

## Browser UI

A live preview UI runs in the browser via WebAssembly. Write YAML in the left panel, get SVG on the right, instantly.

```bash
make run-ui-docker   # serves on http://localhost:8081
```

Features:
- Real-time SVG rendering as you type
- Select and stack multiple mixin files via a combobox
- Toggle individual mixins on/off without removing them
- Click boxes in the SVG to pin/expand them in a filtered view
- Blacklist boxes to remove them from the current view
- Undo support for box selection changes
- Minimap for large diagrams
- Comment legend panel

> **Chrome rendering note:** If the SVG appears blurry after zooming, press `Ctrl++` then `Ctrl+-` to reset the raster.

---

## Development

```bash
make test            # go test --cover ./...
make build-wasm      # build the .wasm for the UI
make generate-all    # regenerate types from JSON schemas (requires Docker)
```

Generated files (do not edit manually):
- `pkg/types/boxes/boxes.go`
- `pkg/types/boxes/externals.go`
- `pkg/types/boxes/boxesdoc.go`
- `pkg/types/boxes/formats.go`
- `pkg/types/boxes/shared.go`

Edit the schemas in `configs/models/` instead, then run `make generate-all`.

---

## Miscellaneous

```bash
# Encode an image as base64 for embedding in YAML
base64 -w 0 my_icon.png
```
