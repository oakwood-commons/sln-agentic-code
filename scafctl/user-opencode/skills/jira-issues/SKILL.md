---
name: jira-issues
description: "Use when the user asks to create a Jira ticket/story, read a Jira key (PROJ-123), or look at Jira tickets to build a plan or Herdr DAG via the Atlassian/Jira MCP server."
---

# Jira Tickets & Stories (`jira-issues`)

Use the configured Jira / Atlassian MCP tools whenever the user says *"create a
Jira ticket"*, references a Jira key (`PROJ-123`), or asks to *"look at my Jira
tickets"* for planning (`planner`) or parallel orchestration (`work-herd`).

If no Jira MCP server is connected yet in the user's harness config, point them
to the remote OAuth 2.1 endpoint (`https://mcp.atlassian.com/v2/mcp`) rather
than guessing CLI commands.

## 1. Creating a Jira Ticket (After `create-issue` Draft Approval)

1. Draft the content with `create-issue` first (Plain English summary, behavior,
   and acceptance criteria at the top; grounded codebase context at the bottom)
   and confirm the draft with the human before creating the ticket.
2. Determine the target **Project Key** (for example `PROJ` from the user's
   prompt, `AGENTS.md`, or by listing accessible Jira projects via MCP) and
   issue type (`Story`, `Task`, or `Bug`).
3. Create the ticket via the Jira MCP create-issue tool and return the new key
   and browse URL.

## 2. Reading & Querying Jira Tickets (For `planner` and `work-herd`)

- **Single ticket pickup (`PROJ-123`):** Fetch the ticket summary, description,
  acceptance criteria, status, and comments via the Jira MCP get-issue tool,
  then validate every assumption against live code (`planner`).
- **Shortlisting Jira tickets for a Herdr plan (`work-herd`):** When the user
  asks to *"look at my Jira tickets to make a Herdr plan"*, query via JQL
  scoped to their project, assignee, or active sprint rather than pulling an
  unbounded backlog:
  - Current sprint assigned to user:
    `project = PROJ AND assignee = currentUser() AND sprint in openSprints() AND statusCategory != Done ORDER BY Rank ASC`
  - Unassigned / team sprint or label filter:
    `project = PROJ AND sprint in openSprints() AND statusCategory != Done ORDER BY priority DESC`
  Present the 4 to 8 candidate tickets and build the `work-herd` dependency and
  hotspot DAG using the Jira keys (`PROJ-101`, `PROJ-102`).

## 3. Branch & PR Linking Conventions

- **Branch name:** `<type>/<KEY>-<slug>` preserving the uppercase Jira key (for
  example, `feat/PROJ-123-expiring-links`).
- **Pull request title / footer:** Reference `<KEY>` in the PR body (`Refs
  PROJ-123`) so Jira's GitHub integration links the branch and PR automatically.
