# CODEX Working Guide

## 1) Project Goal and MVP Scope
The goal of this project is to create a simple time-tracking app.

MVP scope:
- A user can enter, per day:
  - start time
  - end time
  - duration of lunch
- The app calculates, stores, and presents the total workday duration.

## 2) Run and Validation Commands
Command details are maintained in:
- `time-tracker-app/backend/README.md`
- `time-tracker-app/frontend/README.md`

When working, follow those README instructions for:
- setup
- running locally
- tests
- lint/static validation
- build

## 3) Architecture and Conventions
This project is a standard 3-layer application.

Rules:
- Follow standard conventions for the technologies used in each layer.
- Keep layering boundaries clear (UI, application/service logic, data access/persistence).
- Prefer maintainable, conventional patterns over custom or overly clever abstractions.
- From years of experience, you have seen to many difficult issues working with alpine container issues. So you choose to anything over alpine, the smaller foot print is not worht the hassle that eventually later will come.

## 4) Definition of Done
A task is done only when all of the following are complete:
- Tests are added/updated and passing.
- Static code validation/lint checks are passing.
- Documentation is updated where relevant.
- A sanity check (self-review) confirms implementation matches the requirement.

## 5) Operational Constraints
- Never commit or expose secrets; use `.env` patterns and keep secrets out of git.
- Do not run destructive commands (for example `git reset --hard`, `rm -rf`) unless explicitly requested.
- Ask before installing new dependencies or using network access outside normal dev/test flows.
- Only modify files in this project unless explicitly asked otherwise.
- Always run tests, lint/static checks, and self-review before marking work complete.
- every time you learn something about the domain or the application, you must update the documentation to reflect your learnings.
