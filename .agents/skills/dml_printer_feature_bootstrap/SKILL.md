---
name: dml-printer-feature-bootstrap
description: A workflow for scaffolding new DML printer features from grammar rules, including tests and abstraction skeletons.
---

# DML Printer Feature Bootstrap

This skill outlines a systematic, scalable approach to adding new formatting structures to the `dmlint` printer package after a grammar update in `dml.g4`. It enforces test-driven development, strict mapping to ANTLR contexts, and proactive error suppression.

## 1. Plan & Research
- **Analyze Grammar:** Inspect `pkg/dmlex/dml.g4` to identify the new rules to be implemented. Trace down dependencies and nested structures (e.g., initializers, statements, expressions).
  - User may provide some pseudo g4 code with annotation:
    - \n - means force wrap
    - \r - optional wrap
  - Check potential optimization of the `pkg/dmlex/dml.g4` file (Group similar rules). Any changes to the file has to be appoved by user.
- **Propose Implementation:** Draft a plan that clarifies any formatting ambiguity with the user. Specifically ask about:
  - Line wrapping strategies (e.g. keeping elements on a single line if they fit post-wrap).
  - Trailing comma or specific grammar symbol handling (e.g., `...`).
  - Compiler workarounds that require skipping tests via `task006Skip: true`.

## 2. Test Implementation (TDD Phase)
- **Target Location:** `test/printer/` (e.g., `dml_initializer_test.go`).
- **Required Coverage:**
  - Basic usage (single item, many items).
  - Exhaustive line wrapping logic respecting `lineCap = 40` (items exceeding cap, multiple long items, long prefixes forcing inner wraps).
  - Nested combinations (e.g., designated inside single, deeply nested structures).
  - Comment alignment for both `//` and `/* */`, placing them immediately before brackets, after brackets, and inline.
  - Stripping of unsupported syntax tokens (like trailing commas) enforced by adding `task006Skip: true` to the respective test cases.
- **Verification:** Run `go test -c ./test/printer` using a tester subagent to ensure there are no compilation errors in the newly created test structs. All new tests must initially call `t.Skip("pending printer implementation...")` to avoid breaking CI.

## 3. Abstraction Skeleton Creation
- **Target Location:** `pkg/printer/` (e.g., `dml_initializer.go`).
- **Structure:** For each relevant grammar rule, define a print method on the `Printer` struct matching the ANTLR generated context interfaces from `dmlex`.
  - Naming Convention: CamelCase `print<RuleName>` (e.g. `printSingleInitializer`).
  - Argument: `node dmlex.I<RuleName_Capitalized>Context`.
  - Prefix with a tag: `// g4-tag: <token_name>`.
- **Implementation Stub:** All new methods should return `p.printUnimplemented(node)`.
- **Unused Function Suppression:** To satisfy linters and Go compilers before full integration, link the abstraction methods locally. For instance, have the parent rule's print method dummy-assign to the child rule's print method (`_ = p.printChildRule`), or integrate a placeholder into existing handlers (like `printParamStatement`).

## 4. Final Review
- **Walkthrough:** Summarize the additions explicitly in the `walkthrough.md` artifact, categorized by test suites and printer stubs created.
