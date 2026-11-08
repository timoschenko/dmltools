# DML Code Formatting Guide

This directory documents the formatting principles and rules enforced by the DML formatter.

---

## Overview

The DML formatter provides an opinionated, deterministic code formatting style
for Device Modeling Language (DML) code.
The goal is to eliminate debates around code style, improve readability,
and ensure consistent layout across device models.

### Core Formatting Principles

- **Indentation**: Standard 4-space indentation for all nested blocks (methods, objects, registers, templates, control flow).
- **Spacing**:
  - Consistent single space around binary, assignment, and logical operators (`+`, `-`, `=`, `==`, etc.).
  - Punctuation clinging: Commas (`,`) and semicolons (`;`) cling directly to the preceding token without leading whitespace (`item,` and `stmt;`).
  - Single space after commas, colons (in types/labels), and around inline comments.
  - Braces: Standard K&R brace placement for blocks and bodies (`{\n ... \n}`), with opening braces preceded by a space.
- **Predictable Line Length**: Respects a configurable line capacity (typically 80 or 100 characters) using an all-or-nothing wrapping strategy.
- **Semantic Integrity**: Formatting never alters code semantics, and comments are preserved with strict attachment to the tokens they annotate.

---

## Detailed Guides

For in-depth explanations and examples of specific formatting behaviors,
refer to the following guides:

- [**Line Wrapping & Wrap Groups**](line_wrapping.md):
  Explains the `black`-style all-or-nothing explosion strategy, wrap groups (expressions, argument lists, signatures, initializers), cascading wrapping levels, atomic non-wrapped constructs, and how to use parentheses `(...)` as an escape hatch.

- [**Comment Processing & Tool Annotations**](comments.md):
  Explains how inline and block comments are formatted, including the "sticky" trailing comment guarantee that prevents line breaks between statements/commas and comments so that static analysis directives (Coverity, linters, formatters) are never broken or misapplied.
