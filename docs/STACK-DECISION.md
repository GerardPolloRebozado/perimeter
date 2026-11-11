# Why we chose this stack - NestJS · Kotlin Multiplatform · PostgreSQL

* Status: **Accepted** (evidence in `spikes/` and `decision/RECOMMENDATION.md`, formal ADRs in `docs/adr/`)
* Date: 2026-11-11

## The stack at a glance

| Layer | Choice | Spike evidence | ADR |
|---|---|---|---|
| Backend | **NestJS** (TypeScript, Node.js) | spike 5 - `spikes/5/area/` | [ADR 0001](adr/0001-nestjs-backend.md) |
| Frontend | **Kotlin Multiplatform + Compose** (Android + web) | spike 2 - `spikes/2/kmp/` | [ADR 0002](adr/0002-kotlin-multiplatform-frontend.md) |
| Database | **PostgreSQL** (JSONB for workflow configs) | spikes 2 & 3 | [ADR 0003](adr/0003-postgresql-database.md) |
| Contingency | MongoDB - **not adopted**, deferred with explicit triggers | spike 5 matrix | [ADR 0004](adr/0004-mongodb-workflow-json-future.md) |

## How the decision was made

Every choice is backed by a **timeboxed spike** we actually ran, an **individual weighted matrix** filled independently (spikes 1,2, 3 and 5), and a **discussion of the divergences** between them.

The criteria that mattered for AREA (an automation platform gluing services together):

1. **Cross-platform from one codebase** - small team, every duplicated screen is doubled work.
2. **Async hooks & background execution** - reactions fire when events happen, not inside a request.
3. **OpenAPI generation** - the client contract must be generated, not hand-maintained.
4. **OAuth 2.0 / OIDC ecosystem** - every third-party integration is OAuth.
5. **Ecosystem maturity** - we have limited novelty budget and a fixed deadline.
6. **Relational integrity** - users own workflows, orphans are a bug, not an edge case.
7. **Dynamic JSON storage** - each workflow's action/reaction config is arbitrary JSON.

## Backend: NestJS

**Why it was chosen:** it scored highest on the criteria that actually drive AREA.

`@nestjs/swagger` generates the OpenAPI document from controller decorators (5/5). Passport covers OAuth for every provider we glue (4/5). Handlers are async with first-party `@nestjs/schedule` and `@nestjs/bullmq` for cron and queues (5/5) - the reaction pipeline needs exactly that.

**What was rejected and why:**

* **Spring Boot** - capable and mature, lost on ecosystem breadth and team velocity in the weighted matrix.
* **Go (net/http, chi, echo, fiber, gin, huma)** - scored well on performance (~4.1/5 in spike 3) but every router meant assembling OpenAPI, OAuth, queues and validation ourselves. Integration work instead of product work.
* **Loco (Rust)** - best raw-merit loser: 5/5 on speed and security, built-in queues, Rails-like CLI scaffolding, OpenAPI generated from code comments. Lost on ecosystem maturity (3/5): at spike time it was "very new, not even stable, lacks documentation", with breaking changes expected and OAuth a 3/5 because the crates must be wired by hand. On a fixed deadline the innovation budget could not be spent on the framework instead of the product.

**Cost we accept:** Node.js is not the fastest runtime. Hot paths go to BullMQ workers instead of request handlers. `pnpm audit` runs in CI, versions are pinned via `package-lock.json`.

## Frontend: Kotlin Multiplatform (Compose)

**Why it was chosen:** one codebase for Android and web (5/5 cross-platform). In the spike, UI, state, networking and JSON were 100% shared (`commonMain`), platform-specific code was a backend URL and two entry points. Shared logic tests on the JVM with Ktor's MockEngine - no device needed. And because the backend is OpenAPI-generated, `@Serializable` data classes give us a compile-time-checked contract.

**What was rejected and why:**

* **Rust + Dioxus** - interesting but niche, the spike's own verdict was "not stable enough, breaking changes ahead, missing integrations" (spike 1, weighted 3.50/5).
* **React Native + Expo** - the closest runner (4.50/5 in spike 5's matrix). We preferred KMP's type-safe shared contract and no-JS-VM mobile runtime, RN's tooling advantage did not outweigh it for a small team already shipping a Kotlin-typed API client and it slower due to JS.

**Cost we accept:** The spike measured an ~11.5 MB Wasm web bundle for one screen.

## Database: PostgreSQL (+ JSONB)

**Why it was chosen:** AREA's data is relational at its corem, users own workflows, workflows reference services, deleting a user must not leave orphans. Postgres gives us foreign keys and transactions **and** JSONB columns for the arbitrary action/reaction configs, queryable via `->>` and GIN indexes. One engine, both worlds.

**What was rejected and why:**

* **Cassandra** - the spike was decisive against it: the same workflow had to live in two tables synced by hand in a logged batch, queries only worked by partition key (`ALLOW FILTERING` errors), and JSON was opaque text. Massive write throughput we do not need.
* **MongoDB (now)** - scored 5/5 on JSON storage but 3/5 on relational integrity, it would add a second engine whose sync with the relational core is our job. Deferred with explicit triggers instead.

**Cost we accept:** JSONB is unvalidated inside the column, so every workflow config is validated at the API boundary with a Zod schema per service.

## The trade-offs we knowingly accept

| Trade-off | How we survive it |
|---|---|
| Node.js performance ceiling | Long reactions run in `@nestjs/bullmq` workers, never inside request handlers |
| ~11.5 MB Wasm web bundle | Code splitting + lazy loading, bundle size tracked in CI, server-rendered fallback |
| Unvalidated JSONB internals | Zod validation per service schema at the API boundary |
| Possible future Mongo migration | Written triggers in ADR 0004, adoption only via a new ADR |

## Where the decision lives

* Spikes and their notes: `spikes/` - spike comparison in `spikes/DISTANCE.md`
* The one-page recommendation: `decision/RECOMMENDATION.md`
* Change procedure: any change to this stack starts a **new ADR** superseding the old one - the old record is never edited, only marked *Superseded by ADR-00XX*.
