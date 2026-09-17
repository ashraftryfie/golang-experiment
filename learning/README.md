# Memory-First Knowledge System (`learning/`)

This directory is designed for **high-efficiency revision and lifelong retention**. Six months from now, you will not have to re-watch tutorials or re-read long blogs to remember how Go works. Everything here is organized for instantaneous conceptual recall.

---

## Structure

```text
learning/
├── README.md               # Overview of the revision and retention system
├── progress.yaml           # Machine-readable progress tracker with verified statuses
├── roadmap.md              # Detailed stage-by-stage learning pathway (Stages 00-21)
├── glossary.md             # Core Go terminology with "Remember It Like This" memory anchors
├── mistakes.md             # Common Go traps, anti-patterns, and debugging lessons
├── cheat-sheets/
│   ├── go-syntax.md        # Rapid syntax reference
│   └── go-commands.md      # Toolchain commands reference
└── notes/                  # Freeform conceptual deep-dives and runtime mechanics
```

---

## 5-Level Mastery Scale

Every stage and topic in `progress.yaml` moves through five verified states:

1. **`NOT_STARTED`**: No work done yet.
2. **`IN_PROGRESS`**: Actively studying examples and writing exercises.
3. **`PRACTICED`**: All starter exercises completed and passing tests (`go test`).
4. **`MASTERED`**: Practical challenge built, self-review questions passed, code reviewed by **Go Reviewer**, zero warnings from `go vet` and `go test -race`.
5. **`NEEDS_REVIEW`**: Spaced-repetition trigger flagged if a concept has not been touched in >30 days or if a related mistake occurred.

---

## The "Remember It Like This" Rule

Every difficult concept in this repository has a short, high-impact mental anchor designed to stick in long-term memory. Whenever you learn a tricky concept (e.g. slice re-allocation, value vs pointer receivers, channel deadlocks), distill it into a *"Remember It Like This"* block and document it here.
