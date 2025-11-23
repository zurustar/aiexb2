# AI Agent Guidelines

This document outlines the rules and guidelines for AI agents working on this project, specifically regarding documentation maintenance and dependency tracking.

## Documentation Rules

### Rule 1: Documentation Consistency
When you modify a file in the `docs/` directory, you **MUST** check its "Depended On By" list (defined in the file header). You are responsible for updating those dependent files to ensure consistency across the documentation set.
*   **Example:** If you add a new requirement to `docs/requirements.md`, you must check `docs/usecases.md` and `docs/ieee830.md` and update them if the new requirement impacts them.

### Rule 2: Dependency Tracking
All markdown files in the `docs/` directory **MUST** include a metadata section at the very top of the file. This section explicitly lists the file's dependencies.

**Format:**
```markdown
<!--
Depends On: [Relative Path to Parent Document]
Depended On By: [Relative Path to Child Document 1], [Relative Path to Child Document 2]
-->
```

*   **Depends On:** The source document(s) that this file is derived from. Changes in the source document usually trigger updates in this file.
*   **Depended On By:** The downstream document(s) that are derived from this file. Changes in this file usually trigger updates in the downstream documents.

### Rule 3: Design Document Code Policy
When working on design documents in the `docs/` directory, you **MUST** follow these guidelines for code inclusion:

**ALLOWED in Design Documents:**
*   Database schema definitions (CREATE TABLE, ALTER TABLE, etc.)
*   Index definitions (CREATE INDEX)
*   Partition definitions (CREATE TABLE ... PARTITION OF)
*   Conceptual algorithms and flow descriptions
*   Mermaid diagrams and flowcharts

**PROHIBITED in Design Documents:**
*   Application source code (Go, JavaScript, Python, etc.)
*   Specific function/method implementations
*   Error handling implementation details
*   Monitoring/metrics collection code
*   Any code that would duplicate what belongs in the actual source code

**Rationale:** Source code itself is the true detailed design document. Design documents should focus on "what to achieve" (concepts), while "how to implement" belongs in the actual source code to avoid duplication and maintain single source of truth.

## Project Structure & Dependencies

The general flow of documentation dependency is as follows:

1.  `docs/requirements.md` (Root)
2.  `docs/usecases.md` (Derives from Requirements)
3.  `docs/ieee830.md` (Derives from Requirements & Use Cases)
4.  `docs/basic_design.md` (Derives from SRS)
5.  `docs/detailed/*.md` (Derives from Basic Design)
