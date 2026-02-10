# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands

```bash
go build ./...                     # Build all packages
go test ./...                      # Run all tests (including examples)
go test -v                         # Run root package tests with verbose output
go test -run TestName              # Run a single test by name
go test -v ./examples/gamification # Run gamification example tests
go run ./cmd/demo                  # Run demo web server (http://localhost:8080)
```

No Makefile, linter, or CI pipeline is configured. Standard Go toolchain only.

## Architecture

This is a **forward-chaining inference engine** in Go wrapped in a **6-step deterministic pipeline** that transforms input signals into validated, explainable outputs.

### Core Types (all in package `inference`)

- **Fact** (`fact.go`) — Atomic unit of knowledge with ID, Value, Source, and DerivedFrom tracking. The `Accumulative` flag controls whether derived facts survive re-evaluation. `Source` tracks provenance (user, extracted, inferred).
- **Rule** (`rule.go`) — A string expression evaluated against the current fact set using the [Expr language](https://github.com/expr-lang/expr). `WeightedRule` pairs a rule with a probability weight.
- **Inference** (`inference.go`) — A set of weighted rules (conditions) that, when satisfied, produce a new fact. Supports dynamic ID and Value via Expr expressions. Inferences are applied in `Order`.
- **Conclusion** (`conclution.go`) — Asserts whether a set of expected facts hold; `Certainty()` gives partial-match confidence.
- **Contradiction** (`contradiction.go`) — Declares mutually exclusive facts. `Detect()` finds conflicts; `Resolve()` removes conflicting facts (cascading through `RemoveDerivedFrom`).
- **KnowledgeBase** (`knowledgebase.go`) — Central orchestrator. Key flow: `AddFact()` -> `Infer()` -> `ResolveContradictions()`. `GetMissingFactIDs()` returns fact IDs needed by pending inferences.

### Pipeline Types

- **ConfidenceLevel** (`confidence.go`) — High (>=0.8) / Medium (>=0.5) / Low (<0.5) from certainty scores.
- **Domain** (`domain.go`) — Domain detection via keyword signals: finance, ecommerce, infrastructure, data, aiml, general.
- **Intent** (`intent.go`) — Intent classification using Expr rules: query, decision, analysis, action.
- **Entity** (`entity.go`) — Entity extraction from facts using Expr rules, producing new facts for the KB.
- **Constraint** (`constraint.go`) — Hard/soft constraints with Expr-based evaluation and identification.
- **Risk** (`risk.go`) — Risk analysis: explicit rules + contradiction detection + low-certainty warnings.
- **Solution** (`solution.go`) — Multi-dimensional scoring (BusinessImpact, Complexity, Risk, TimeToValue) and ranking.
- **PipelineResult** (`output.go`) — Structured output: Result, Reasoning, Confidence, FollowUp, Risks.
- **Pipeline** (`pipeline.go`) — 6-step orchestrator: Domain -> Intent -> Extract -> Constrain -> Infer -> Risk.
- **PipelineConfig** (`pipeline_config.go`) — JSON serialization/deserialization for pipeline configs.

### Expression Evaluation

Rule conditions and dynamic values are Expr strings compiled and evaluated at runtime (`utils.go`). `extractFacts()` walks the Expr AST to discover which facts a rule depends on.

### Persistence

`KnowledgeBase.Dump()` serializes to JSON. `LoadPipelineConfig()` / `SavePipelineConfig()` handle full pipeline configs. See `examples/triage/definition.json` for the pipeline config schema.

### Sessions

`KnowledgeBase.Start()` increments a running counter and clears facts while preserving structure.

### Web Demo

`cmd/demo/main.go` — Go web server using stdlib `net/http` with embedded static files. REST endpoints: `/api/examples`, `/api/pipeline/load`, `/api/pipeline/run`, `/api/pipeline/pending`, `/api/pipeline/reset`.

### File Structure

```
fact.go, rule.go, inference.go, conclution.go, contradiction.go  # Core engine
knowledgebase.go, utils.go                                        # Orchestration
confidence.go, domain.go, intent.go, entity.go                   # Pipeline steps
constraint.go, risk.go, solution.go, output.go                   # Pipeline steps
pipeline.go, pipeline_config.go                                   # Pipeline orchestrator
cmd/demo/main.go, cmd/demo/static/index.html                     # Web demo
examples/triage/definition.json                                   # Hospital triage config
examples/gamification/definition.json, pipeline.json              # Gamification configs
```
