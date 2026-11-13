# Spike comparison - distance travelled

Three timeboxed spikes, same target: a server listening on port 8080, answering `GET /ping` with exactly `{"status":"ok"}`, port read from an environment variable.

**Our definition of "distance travelled":** how far each spike got *beyond* the bare ping target - because any framework can serve `/ping`. Distance = ping target + a workflow CRUD slice + a persistence layer, measured to the point the spike was demonstrably working (curl-able, or on-device for the frontend spike). A spike that reached less but taught us more still counts - that is recorded in the "what stopped us" column.

## The three runs

| | **NestJS** (`spikes/5/surface-volume/`) | **Kotlin Multiplatform** (`spikes/2/kmp/`) | **Cassandra + Spring** (`spikes/2/cassandra/`) |
|---|---|---|---|
| Layer probed | Backend | Frontend (Android + web) | Persistence |
| Ping target (`:8080/ping`, env port) | ✅ reached | ✅ reached (web + Android calling `:8080`) | ✅ reached - with a **503** variant when the DB is down |
| Went beyond ping | Full NestJS service scaffold (`nest-cli`, vitest configs, scripts) | Real UI screens, Ktor client, `expect/actual` backend URLs, 4 MockEngine tests | Full workflow schema (2 tables), logged-batch dual writes, REST CRUD for workflows, 3 tests |
| Timebox used mostly on | Project structure + OpenAPI/service conventions | Gradle setup: plugins, targets, source sets, version catalog | Cassandra container startup (~75 s) + per-query table modelling |
| What stopped us / what we'd need | Nothing structural - untested: real OAuth flows, BullMQ workers | ~11.5 MB Wasm bundle, untested web accessibility (canvas + screen readers) | `ALLOW FILTERING` errors for non-key queries, JSON as opaque text, ~1.7 GB RAM |
| Verdict fed into | **Adopted** (ADR 0001) | **Adopted** (ADR 0002) | **Rejected** - informed ADR 0003 (Postgres) & 0004 (deferred Mongo) |

## What the exercise actually taught us

* **The bias check worked.** The spike we liked least on paper (Cassandra) was the one that produced the most decision-relevant evidence - it is why PostgreSQL won on *integrity*, not on vibes, and why MongoDB became a documented contingency instead of a hunch.
* **"Distance" is layer-relative.** The frontend spike travelled the least in *features* and the most in *risks retired* (Wasm bundle size, platform-leak debugging, Gradle cost).
* **Per-query table modelling is the Cassandra cost**, not its syntax - that finding, not a benchmark, killed it for Surface Volume.