# ADR 0001 - Use NestJS for the backend

* Status: **Accepted** (decided from spike evidence + weighted matrix, see `decision/RECOMMENDATION.md`)
* Date: 2026-11-11
* Deciders: the whole team

## Context

AREA needs an application server that:

* Exposes a REST API consumed by a mobile + web client,
* Triggers **actions** when external events happen (webhooks, polling) and runs **reactions** - including asynchronous, long-running work outside the request cycle,
* Talks to third-party services with OAuth 2.0 / OIDC,
* Serves an OpenAPI document so the client can be generated against a stable contract.

The spikes compared **NestJS**, **Java + Spring Boot**, a set of **Go** routers (net/http, chi, echo, fiber, gin, huma - see `spikes/3/evaluation.md`) and **Loco**, the Rust backend framework evaluated in spike 1 (`spikes/1/stack-matrix.md`).

## Decision

Use **NestJS** (TypeScript on Node.js) as the backend framework.

### Alternatives rejected

* **Java + Spring Boot** - capable and mature. Scored lower than NestJS on ecosystem breadth and team velocity in the weighted matrix. No structural blocker, the choice is recorded so it can be revisited deliberately.
* **Go (net/http, chi, echo, fiber, gin, huma)** - strong performance (~4.1/5 in spike 3 matrix) but batteries-not-included: OpenAPI, OAuth, queues and validation must each be assembled from separate libraries. That is integration work.
* **Loco (Rust)** - best raw-merit loser: 5/5 on speed and security (memory safety, auth pre-wired), built-in queues, Rails-like CLI scaffolding, OpenAPI generated from code comments. Lost on ecosystem maturity (3/5): at spike time it was "very new, not even stable, lacks documentation", with breaking changes expected and OAuth a 3/5 because the crates exist but must be wired by hand.

## Consequences

### Positive

* **OpenAPI first-class.** `@nestjs/swagger` builds the OpenAPI document from controller decorators - scored 5/5 in the spike. The KMP client contract comes from it.
* **OAuth/OIDC ecosystem.** Passport strategies cover every provider AREA glues together, scored 4/5 against Spring's 3/5 on this criterion in our matrix.
* **Async hooks and background execution** scored 5/5: handlers are async, and first-party `@nestjs/schedule` (cron) and `@nestjs/bullmq` (queues) cover the reaction pipeline without homegrown infra.
* **Ecosystem maturity** scored 5/5 - an established, widely used framework, so our timeboxed spike produced a working server quickly.
* **Team velocity.** We can recruit knowledge from the largest JS/TS ecosystem, AI tooling knows it deeply, which speeds iteration.

### Negative

1. **Runtime performance ceiling.** Node.js is enough for our API volume but is not the fastest server runtime (spike scored it 4/5). If an action/reaction pipeline becomes a hot path, that work must be moved into queues (`@nestjs/bullmq`) and workers rather than inside request handlers.
2. **Node/JS dependency churn.** `package.json` carries more transitive dependencies than the Spring equivalent, and security updates arrive more often. We pay for this with `npm audit` in CI and by pinning versions via `package-lock.json`.
3. **NestJS magic (decorators, DI) has a learning curve.** Annotations hide a lot, as the Spring spike showed too, we mitigate by keeping the module layout conventional and documented.
