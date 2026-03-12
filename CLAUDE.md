# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Generated code

The project uses codegen approaches to create some files out of JSON schema data models.
Attention, this files are never touched manually. They are ajusted via editing the related
JSON schemas in `configs/models` and are generated via `make generate-all`

* `pkg/types/boxes/boxes.go` - input definition model
* `pkg/types/boxes/externals.go` - additional input definition models
* `pkg/types/boxes/boxesdoc.go` - domain model, created out of the input model
* `pkg/types/boxes/formats.go` - model for format definitions, used in input- and domain models
* `pkg/types/boxes/shared.go` - model for some shared types, used in input- and domain models


## Commands

```bash
# Build CLI tool
make build               # outputs to build/draw

# Run tests
make test                # go test --cover ./...
go test ./pkg/...        # test specific packages
go test ./pkg/boxesimpl/... -run TestName  # run a single test

# Build WebAssembly for the UI
make build-wasm

# Run code generation (requires Docker)
make generate-all

# Run the UI locally via Docker
make run-ui-docker       # serves on port 8081
```

## Architecture

This is a **diagram-as-code** tool that renders hierarchical box diagrams from YAML specifications. It targets both a CLI tool and a browser-based UI (via WebAssembly).

### Model-Driven Development

The project uses YACG to generate Go type definitions from JSON schemas. The workflow is:
1. Edit schemas in `configs/models/` (e.g. `boxes.json`, `formats.json`, `shared.json`)
2. Run `make generate-all` to regenerate types in `pkg/types/`
3. Hand-written implementation files are kept separate from generated files

Generated files: `pkg/types/formats.go`, `pkg/types/shared.go`, `pkg/types/boxes/boxes.go`, `pkg/types/boxes/boxesdoc.go`
Custom extensions: `pkg/types/boxes/boxestype_impl.go`, `pkg/types/boxes/layout.go`, `pkg/types/boxes/connect.go`, etc.

See `ImplementationSteps.md` for the full model-driven workflow.

### Data Flow

```
YAML input → load (pkg/types/load.go)
           → internal document (boxesdocfactory.go)
           → layout positioning (pkg/types/boxes/layout.go)
           → connection routing via Dijkstra (pkg/types/boxes/dijkstra.go, connect.go)
           → SVG rendering (pkg/svgdrawing/, pkg/svg/)
           → output SVG file or WASM return value
```

### Key Packages

- **`cmd/`** — Cobra CLI entry point; subcommands: `boxes`, `randomize`, `truncate`, `version`
- **`pkg/boxesimpl/`** — Orchestrates the full pipeline: load → layout → connect → draw
- **`pkg/types/`** — All data types; `boxes/` subdirectory has all box-diagram-specific types including layout, connection routing, roads, and comments
- **`pkg/svgdrawing/`** — SVG canvas and drawing primitives (text, shapes, colors)
- **`wasm/main.go`** — WebAssembly entry point; exports `createSvg`, `createSvgExt`, `createSvgForConnected` to JavaScript
- **`ui/html/`** — Vanilla JS + HTMX frontend; calls WASM functions for real-time rendering

### Input Model Concepts

- **Boxes** — Hierarchical containers arranged in rows/columns; can be nested
- **Formats** — Style definitions (fonts, colors, fills, line styles) referenced by name
- **Format variations** — Conditional styling based on tags
- **Connections/Roads** — Paths between boxes, routed to avoid overlaps
- **External mixins** — Inject connections or box definitions from separate YAML files
- **Comments** — Multiple comments per box, aggregated and displayed on the diagram
- **Overlays** — Heatmaps and visual layers on top of the diagram

### Known Issues

See `known_issues.md`. Notably: connections with captions only work for mixins.
