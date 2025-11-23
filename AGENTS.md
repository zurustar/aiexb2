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

## Task Management Rules

### Rule 4: Task Status Tracking
All implementation work **MUST** be tracked in `docs/tasks.md`. You are responsible for updating task status before and after each work session.

**Task Status Indicators:**
- `[ ]` - Not started
- `[/]` - In progress (include assignee name)
- `[x]` - Completed
- `[R]` - Under review / Review requested
- `⚠️` - Requires individual review (critical tasks)

**Mandatory Update Points:**

1. **Before Starting Work:**
   - Change task status from `[ ]` to `[/]`
   - Add assignee name (e.g., `[/] AI Assistant - 2025-11-23`)
   - Update the "最終更新" (Last Updated) timestamp at the top of the file

2. **After Completing Work:**
   - Change task status from `[/]` to `[x]` or `[R]` (if review needed)
   - Update the "最終更新" timestamp
   - If the task has `⚠️` marker, change status to `[R]` for human review

3. **When Requesting Review:**
   - Change task status to `[R]`
   - Use `notify_user` tool to request review
   - Include the task number in the review request

**Example Workflow:**
```markdown
# Before starting task 2.1
- [ ] 2.1 ユーザードメイン

# When starting work
- [/] 2.1 ユーザードメイン (AI Assistant - 2025-11-23)

# After completing implementation
- [x] 2.1 ユーザードメイン

# For critical tasks requiring review
- [R] 2.5 予約ドメイン ⚠️ (AI Assistant - 2025-11-23)
```

**Important Notes:**
- **NEVER skip task updates** - This is critical for project continuity
- Update tasks.md **in the same commit** as the implementation
- If you're unsure which task you're working on, check `docs/implementation_plan.md`
- Always check dependencies before starting a task
- If a task is blocked by dependencies, note it in the task file

### Rule 5: Phase Completion Checklist
When completing a Phase:
1. Ensure all tasks in the Phase are marked `[x]` or `[R]`
2. Update the Phase checkpoint task status
3. Request human review via `notify_user` with paths to all completed files
4. Do not proceed to the next Phase until the current Phase is reviewed and approved
