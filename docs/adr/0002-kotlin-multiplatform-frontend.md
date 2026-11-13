# ADR 0002 - Use Kotlin Multiplatform (Compose) for the frontend

* Status: **Accepted** (decided from spike evidence + weighted matrix, see `decision/RECOMMENDATION.md`)
* Date: 2026-11-11
* Deciders: the whole team

## Context

Surface Volume needs a client that runs:

* As a **mobile app** (Android at minimum),
* As a **web app**,
* From **one codebase**, because the a duplicated UI doubles every screen's cost.

The spikes compared **Kotlin Multiplatform (KMP + Compose)**, **Rust + Dioxus** (spike 1) and **React Native with Expo** (spike 5).

## Decision

Use **Kotlin Multiplatform** with **Compose Multiplatform**, targeting Android and web (Wasm). Shared code (UI, state, networking, JSON serialization) lives in a `shared/` module and platform-specific code stays tiny (entry points, backend URL via `expect`/`actual`).

## Consequences

### Positive

* **Single codebase.** The spike's UI, state, networking and JSON were 100% in `commonMain` - cross-platform scored 5/5. Per-platform code was 2 lines of backend URL, 1 dependency, 2 entry points.
* **Type-safe contract with the backend.** `@Serializable` data classes mirror the OpenAPI schema, a changed field is a compile error, not a runtime surprise.
* **Testable logic.** The same business logic runs on the JVM (tested with Ktor's MockEngine fake-backend tests without a server or a device).
* **One language across client and shared logic**, with no JavaScript VM on mobile (the web client talks DOM via WebAssembly).

### Negative

1. **The web bundle is heavy.** The spike measured **~11.5 MB of Wasm** for a single screen, and WasmGC requires a modern browser. First page load is slow, code splitting and lazy loading are mandatory, and mobile web users on poor connections will notice.
2. **Web target is young.** Compose web drew on a canvas in the spike, screen-reader and keyboard behavior was **not verified** and accessibility scored 3–4/5 and must be tested early, if it fails on web, browsers fall back to a server-rendered page.
3. **Gradle tooling is slow and error-prone.** Plugin/target/source-set configuration was the most time-consuming part of the spike, and the first build downloads a lot. We accept a costlier onboarding and CI.
4. **Platforms leak through debugging.** An Android bug was invisible in the web build and vice versa. "Write once" holds for code, not for debugging.
