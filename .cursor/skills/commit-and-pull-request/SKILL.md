---
name: commit-and-pull-request
description: >
  Writes branch names, commit messages, and pull request descriptions.
  Branch names stay short and meaningful. Titles start with one relevant
  gitmoji and broadly say what the change does. Bodies stay short, explain
  why, use bullet points where they help, include a decision link only when
  one was used, and note how it was tested. Use when the user asks for a
  branch name, commit message, pull request title, pull request body, or to
  commit or open a PR.
---

# Commit and pull request

Use this format for the branch, the commit message, and the pull request. The pull request title matches the commit title. The pull request body matches the commit body. Do not use a Summary / Test plan template.

## Requirements

- The branch name should be short, meaningful, and describe the change.
- The title should broadly capture what the change is doing, and start with one relevant gitmoji.
- The body of the commit message:
  - Should not be a wall of text.
  - Add bullet points where they make the commit message or the pull request description easier to scan.
  - Should explain why the change is being made. Do not add a "Why" heading.
  - Include any other relevant information about the work. Include a Slack or decision link only when one was actually used.
  - State how the change was tested. If nothing was run, say so. Example: "No runtime tests; this is markdown only."

## Apply

- Branch: a few words, lowercase, hyphens, meaningful to a reviewer. No ticket prefix. No generated suffix.
- Title: one relevant gitmoji, then one imperative line. No ticket key, no period. Use `📝` for docs or skills, `✨` for a feature, `🐛` for a fix, `♻️` for a refactor, `✅` for tests, and `⬆️` for a dependency. Do not add a second gitmoji.
- Stacks: when a change is one part of a stack of pull requests on the same subject, prefix the title with `N/M [Group]`, for example `1/2 [Ingest] ✨ Accept a segment`. The group name is a short label the reader recognises; the number tells reviewers where they are in the stack. A single standalone PR on a subject can still use `[Group]` without the `N/M`.
- Body: a few short paragraphs, with bullet points where a list is clearer than a paragraph. Lead with the reason. Include a Slack or decision link only when one was actually used. One or two sentences on how it was tested. If nothing was run, say so, for example "No runtime tests; this is markdown only."
- Do not add a git trailer. There is no ticket key for this repo.
- Do not add "Made with Cursor". It is not relevant to the work. If a tool appends that line, remove it before the pull request is opened.

```
agent-rules

📝 Add agent rules and skills

Agents need one place that lists the constraints and the skills that add to them, so a later change does not follow another repo by mistake.

No runtime tests; this is markdown only.
```
