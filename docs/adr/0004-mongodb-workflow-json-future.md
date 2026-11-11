# ADR 0004 - Do NOT adopt MongoDB now, keep it as a contingency for workflow JSON

* Status: **Proposed** (deliberately deferred - revisit only if the trigger below fires)
* Date: 2026-11-11
* Deciders: the whole team

## Context

Workflow configurations are arbitrary JSON. MongoDB scored 5/5 on dynamic JSON storage (spike 5) versus PostgreSQL's JSONB, and the idea of a dedicated document store for workflows keeps coming up. The Cassandra spike, however, already demonstrated the cost of a non-relational store for this data: no integrity, sync-it-yourself multi-table writes, and JSON that could not be queried.

## Decision

**Do not introduce MongoDB now.** PostgreSQL with JSONB is the single store (ADR 0003). MongoDB is kept as a **named contingency**, adopted only if **all three** triggers are met:

1. Product queries need to **search inside workflow configurations** (e.g. "find every workflow reacting to repo X with keyword Y") at a frequency where JSONB/GIN performance is measurably insufficient,
2. Workflow JSON volume or change rate makes single-primary PostgreSQL a measured bottleneck, not a hypothetical one,
3. The relational core (users, services, ownership) stays in PostgreSQL - MongoDB would only hold workflow documents, with the workflow ID mirroring a Postgres row.

## Consequences

### Positive

* Zero operational cost today: one database to run, back up and reason about.
* The decision is reversible cheaply **now**, and expensive later - which is exactly why it is written down with explicit triggers instead of "we'll see".
* We avoid the Cassandra spike's trap on day one: a second database means the workflow's existence is recorded in two engines whose sync is our job.

### Negative

1. **We accept a possible future migration.** If the triggers fire, moving workflow configs from JSONB to MongoDB is a data migration plus a driver/service rewrite - real weeks of work at the worst time (mid-project).
2. **Two engines, if adopted, means eventual consistency between them** (workflow deleted in Postgres, document still in MongoDB until cleanup). Any adoption must ship the deletion path first.
3. **The team must learn a second database** (operations, backup, indexing) if this ever activates - novelty budget spent on a component, not a feature.
