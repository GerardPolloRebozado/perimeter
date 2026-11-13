# ADR 0003 - Use PostgreSQL as the primary database

* Status: **Accepted** (decided from spike evidence + weighted matrix, see `decision/RECOMMENDATION.md`)
* Date: 2026-11-11
* Deciders: the whole team

## Context

Surface Volume stores **users, linked services and workflows**. The data is fundamentally relational:

* A user owns workflows, deleting a user must not leave orphan workflows,
* Workflows reference an action service and a reaction service,
* But each workflow's **action/reaction configuration is arbitrary JSON** defined by the service's parameter schema (repo name, channel, keywords…), which varies per workflow.

The spikes compared **PostgreSQL**, **MongoDB** (spike 5) and **Cassandra** (spike 2, see `spikes/2/cassandra/NOTES.md`).

## Decision

Use **PostgreSQL** for all persistent data. Workflow configurations are stored in **JSONB columns** inside relational tables, not in a separate document database.

## Consequences

### Positive

* **Relational integrity.** Foreign keys and transactions guarantee a deleted user leaves no orphaned workflows. Cassandra's spike showed the opposite: the same workflow in two tables synced by hand ("nothing stops a deleted user from still having workflows").
* **One engine for both worlds.** JSONB holds the dynamic config (scored 5/5 for dynamic workflow storage) **and** is queryable (`->>` / GIN indexes) - Cassandra could only store JSON as opaque text.
* **Mature operations.** Backup, replication and monitoring are well documented, the team's relational knowledge transfers directly.
* **No second database to run** in development or production.

### Negative

1. **JSONB is schemaless inside a schema.** PostgreSQL does not enforce structure inside the JSON column, an invalid config can be written. We must validate workflow configs **at the API boundary** (Zod schema per service, in the NestJS validation pipe) and accept that the database itself trusts us.
2. **We are betting JSONB stays enough.** If workflow-config queries become the product's core access pattern and outgrow JSONB, we pay a later migration. See ADR 0004 for the trigger that would justify MongoDB.
3. **Vertical scaling first.** PostgreSQL scales reads with replicas but writes stay single-primary, if Surface Volume ever becomes write-heavy at scale we can easily migrate to yugabyte db as they are compatible between them
