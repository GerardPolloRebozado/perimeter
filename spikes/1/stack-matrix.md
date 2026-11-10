# Tech Stack Individual Option Evaluation Tables

---

## 1. Frontend Technologies

| Criterion                               |  Weight  |  Score (1-5)   | Justification / Notes                                                                                                                 |
| :-------------------------------------- | :------: | :------------: | :------------------------------------------------------------------------------------------------------------------------------------ |
| **Cross-Platform / Single Codebase**    |   20%    |       4        | It has very nice support for multiplatform but some native things like notifications depend on platform specific code                 |
| **Native Web/app Compatibility**        |   15%    |       4        | On web it uses the dom and talks to it using webassembly, on phones it uses webview but does not need a javascript vm improving speed |
| **Developer Velocity & Ease of Coding** |   15%    |       3        | It's using rust so its harder due to compiler                                                                                         |
| **Ecosystem Maturity**                  |   15%    |       2        | Its kind of new and has many experimental features                                                                                    |
| **Performance & Speed**                 |   10%    |       4        | Faster than javascript full-stack frameworks due to native code but a little bit slower than native mobile apps                       |
| **Ease of AI Tools & UI Integration**   |   10%    |       3        | Comes with skills and other resources for AI agents but since its very new ai agents don't know that much about it                    |
| **Accessibility (WCAG AA Compliance)**  |   15%    |       3        | Comes with WAI-ARIA compliance and unstyles component swiht accesibility features and i18n                                            |
| **Weighted Total**                      | **100%** | **3.50 / 5.0** | 66%                                                                                                                                   |

---

## 2. Backend Technologies (Application Server)

| Criterion                               |  Weight  |  Score (1-5)   | Justification / Notes                                                                                                                             |
| :-------------------------------------- | :------: | :------------: | :------------------------------------------------------------------------------------------------------------------------------------------------ |
| **OpenAPI / Swagger Generation**        |   15%    |       5        | The swagger is done automatically from comments of the code                                                                                       |
| **OAuth 2.0 & OIDC Ecosystem**          |   15%    |       3        | It has crates like oauth2 that make all the work for you but the user is the one that connects it to the framework                                |
| **Ease of AI Tools & LLM Integration**  |   10%    |       3        | It is a very new framework that is not even stable and lacks documentation                                                                        |
| **Developer Velocity & Ease of Coding** |   10%    |       4        | It's using rust so its harder due to compiler but has a cli that does a lot of boilerplate for you                                                |
| **Async Hook & Background Execution**   |   20%    |       4        | It has build in queues and async functions that symplify long running taks but there is not scheduler or task deduplication                       |
| **Ecosystem Maturity**                  |   10%    |       3        | Very nbew framework but rust crates is mature enough and it uses a mature low level backend framework called axum that has a lot of crates for it |
| **Security**                            |   10%    |       5        | It has rust memory safety and auth is already implemented so the user has less room to screw up auth                                              |
| **Speed & performance**                 |   10%    |       5        | Rust is a very fast language as there is no garbage collector                                                                                     |
| **Weighted Total**                      | **100%** | **0.00 / 5.0** |                                                                                                                                                   |

---

## 3. Database & Persistence Technologies

| Criterion                              |  Weight  |  Score (1-5)   | Justification / Notes                               |
| :------------------------------------- | :------: | :------------: | :-------------------------------------------------- |
| **Dynamic Workflow / JSON Storage**    |   20%    |       5        | It has a lot of types including json for storing it |
| **Developer Velocity & Ease of Use**   |   20%    |       5        | The most used databsae a lot of docs and resources  |
| **Reliability & Relational Integrity** |   20%    |       5        | It has tools to backup and do replicas              |
| **Weighted Total**                     | **100%** | **0.00 / 5.0** |                                                     |

---
