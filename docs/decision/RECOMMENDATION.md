# Recommendation

**Stack:** NestJS, Kotlin Multiplatform + Compose (Android and web), PostgreSQL. Workflow configs sit in JSONB. MongoDB is deferred (ADR 0004).

Two criteria decided it. The team is small, and the product is connecting services.

**One codebase for Android and web.** In the KMP spike the screen, the state and the HTTP client lived in `commonMain`. Both apps were that screen. Platform code was the backend URL and two entry points.

**Background work.** Reactions run outside the HTTP request. NestJS scored 5/5 on this in the spike 5 matrix: async handlers, `@nestjs/schedule` and `@nestjs/bullmq`. Weighted total 4.45/5.

Postgres holds both shapes of data we have. A user owns workflows (foreign keys, transactions). Each workflow's config is JSON that changes per service (JSONB, same database). Cassandra stored one workflow in two tables and rejected queries off the partition key. MongoDB scored 5/5 on JSON and 3/5 on integrity, and would be a second database to sync.

The trade-off is performance. NestJS scored speed 4/5. The KMP web build was ~11.5 MB of Wasm for one screen. Long reactions will run in BullMQ workers, and the web app will be code-split. The spike never started a worker, and CI does not measure the Wasm size yet.

The disagreement was the client. React Native scored 4.50/5 in spike 5, and React is the one we have already used. KMP's web target is young, draws on a canvas, and 11.5 MB is a heavy first load. We kept KMP: one shared screen, and Kotlin generated from the NestJS OpenAPI spec, so a changed backend field fails the build. ADR 0002 records that cost. A different client would be a new ADR.
