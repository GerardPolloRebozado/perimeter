# Surface Volume — Technology Stack Evaluation

## 1. Frontend Technologies

| Criterion                           | Weight   | Score (1-5)    | Justification          |
|-------------------------------------|----------|----------------|------------------------|
| Cross-Platform / Single Codebase    | 20%      | 4              | Web runs everywhere    |
| Native Web/app Compatibility        | 15%      | 4              | Browser plus mobile    |
| Developer Velocity & Ease of Coding | 15%      | 4              | Svelte runes simple    |
| Ecosystem Maturity                  | 15%      | 3              | Smaller than React/Vue |
| Performance & Speed                 | 10%      | 5              | Compiled fast          |
| Ease of AI Tools & UI Integration   | 10%      | 3              | Fewer AI snippets      |
| Accessibility (WCAG AA Compliance)  | 15%      | 4              | Bits UI headless       |
| **Weighted Total**                  | **100%** | **~3.8 / 5.0** |                        |

## 2. Backend Technologies (Application Server)

| Criterion                           | Weight   | Score (1-5)    | Justification             |
|-------------------------------------|----------|----------------|---------------------------|
| OpenAPI / Swagger Generation        | 15%      | 5              | oapi-codegen first-class  |
| OAuth 2.0 & OIDC Ecosystem          | 15%      | 4              | go-oidc JWKS verify       |
| Ease of AI Tools & LLM Integration  | 10%      | 3              | LLMs are stupid           |
| Developer Velocity & Ease of Coding | 10%      | 3              | net/http verbose explicit |
| Async Hook & Background Execution   | 20%      | 4              | Goroutines Asynq queue    |
| Ecosystem Maturity                  | 10%      | 5              | Go libs production-grade  |
| Security                            | 10%      | 4              | Simple middleware audit   |
| Speed & performance                 | 10%      | 5              | Goroutines cheap pgx fast |
| **Weighted Total**                  | **100%** | **~4.1 / 5.0** |                           |

## 3. Database & Persistence Technologies

| Criterion                          | Weight   | Score (1-5)     | Justification            |
|------------------------------------|----------|-----------------|--------------------------|
| Dynamic Workflow / JSON Storage    | 20%      | 3               | Postgres JSONB flexible  |
| Developer Velocity & Ease of Use   | 20%      | 4               | pgx sqlc goose           |
| Reliability & Relational Integrity | 20%      | 4               | ACID transactions proven |
| **Weighted Total**                 | **100%** | **~3.75 / 5.0** |                          |
