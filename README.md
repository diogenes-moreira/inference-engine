<p align="center">
  <img src="logo.png" alt="Inference Engine" width="300">
</p>

# Inference Engine

A forward-chaining inference engine in Go with a structured 6-step deterministic pipeline. It transforms input signals into validated, explainable outputs with confidence levels, reasoning traces, and risk analysis.

## Quick Start

```bash
# Build
go build ./...

# Run tests
go test ./...

# Run the demo web UI
go run ./cmd/demo
# Open http://localhost:8080
```

## Architecture

### Core Engine

The core engine manages facts, derives new facts via rules, tracks certainty/probability, and detects contradictions.

- **Fact** — Atomic unit of knowledge with ID, Value, Source, and DerivedFrom tracking
- **Rule** — An [Expr language](https://github.com/expr-lang/expr) expression evaluated against facts. `WeightedRule` pairs a rule with a probability weight
- **Inference** — Weighted rules that, when satisfied, produce a new fact. Supports dynamic ID/Value via Expr expressions
- **Conclusion** — Asserts whether expected facts hold; `Certainty()` gives partial-match confidence
- **Contradiction** — Declares mutually exclusive facts with automatic detection and resolution
- **KnowledgeBase** — Central orchestrator: `AddFact()` -> `Infer()` -> `ResolveContradictions()`

### 6-Step Pipeline

The pipeline wraps the core engine in a structured process:

```
Input Facts
    |
    v
1. Domain Detection    -> Select domain-specific scoring weights
2. Intent Classification -> Classify what the user wants (query/decision/action/analysis)
3. Entity Extraction    -> Extract structured entities, add as facts to KB
4. Constraint Check     -> Identify hard/soft constraints
5. Knowledge Application -> Run inference engine (Infer + ResolveContradictions)
6. Risk Analysis        -> Evaluate risk rules + check contradictions/low certainty
    |
    v
Structured Output (Result + Reasoning + Confidence + Follow-up)
```

### Pipeline Types

- **ConfidenceLevel** (`confidence.go`) — High/Medium/Low classification from certainty scores
- **Domain** (`domain.go`) — Business domain detection via keyword signals (finance, ecommerce, infrastructure, etc.)
- **Intent** (`intent.go`) — User intent classification using Expr rules (query, decision, analysis, action)
- **Entity** (`entity.go`) — Structured entity extraction from facts using Expr rules
- **Constraint** (`constraint.go`) — Hard/soft constraints with Expr-based evaluation
- **Risk** (`risk.go`) — Risk analysis including explicit rules, contradiction detection, and low-certainty warnings
- **Solution** (`solution.go`) — Multi-dimensional solution scoring and ranking
- **PipelineResult** (`output.go`) — Structured output with result, reasoning, confidence, and follow-up signals

## Usage

### Standalone Knowledge Base

```go
kb := &inference.KnowledgeBase{
    Facts: map[string]inference.Fact{},
    Inferences: []inference.Inference{
        {
            Description: "Can vote",
            Rules: []inference.WeightedRule{
                {Rule: inference.Rule{Expression: "age >= 18"}, Weight: 1},
            },
            FactID:           "can_vote",
            FactValue:        "age >= 18",
            IsValeCalculated: true,
        },
    },
}
kb.Start()
kb.AddFact(inference.Fact{ID: "age", Value: 20})
// kb.Facts["can_vote"].Value == true
```

### Full Pipeline

```go
config := inference.PipelineConfig{
    KnowledgeBase:    kb,
    DomainDetector:   &inference.DomainDetector{Signals: signals},
    IntentClassifier: &inference.IntentClassifier{Rules: intentRules},
    EntityExtractor:  &inference.EntityExtractor{Rules: extractionRules},
    RiskAnalyzer:     &inference.RiskAnalyzer{Risks: risks},
}

pipeline := inference.NewPipeline(config)
result, err := pipeline.Run(inputFacts)
// result.Result, result.Confidence, result.Reasoning, result.Risks, etc.
```

### Loading from JSON

```go
config, err := inference.LoadPipelineConfig("examples/triage/definition.json")
pipeline := inference.NewPipeline(*config)
config.KnowledgeBase.Start()
result, err := pipeline.Run(inputFacts)
```

## Examples

### Hospital Triage (`examples/triage/`)

Classifies emergency patients by vital signs into red/yellow/green urgency levels. Demonstrates entity extraction (fever, tachycardia detection), risk rules (critical temperature, low SpO2), and constraint validation.

### Pizza Gamification (`examples/gamification/`)

Loyalty rewards system: accumulated pizza sales trigger discount tiers (25% at 3 sales, 50% at 5, free at 10). Demonstrates accumulative facts, overwrite inferences, and ordered rule execution.

## Web Demo

```bash
go run ./cmd/demo
```

Opens a web UI at `http://localhost:8080` with both examples. Select an example, add input facts (or use presets), and run the pipeline to see structured output with confidence, reasoning, and risk analysis.

## Project Structure

```
fact.go, rule.go, inference.go       # Core engine types
conclution.go, contradiction.go      # Assertions and conflict resolution
knowledgebase.go, utils.go           # Orchestration and expression evaluation
confidence.go, domain.go, intent.go  # Pipeline step types
entity.go, constraint.go, risk.go    # Pipeline step types
solution.go, output.go               # Scoring and structured output
pipeline.go, pipeline_config.go      # Pipeline orchestrator and config I/O
cmd/demo/                            # Web UI demo server
examples/triage/                     # Hospital triage example
examples/gamification/               # Pizza loyalty example
```
