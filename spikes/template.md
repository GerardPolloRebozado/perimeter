# Tech Stack Individual Option Evaluation Tables

Scores come from the 2026-11-11 spikes. A criterion the spike did not run is scored **2**: that means no evidence from this run, not a measured failure. Notes say which is which.

---

## 1. Frontend Technologies

Option scored: **React Native (Expo SDK 58)**, from `spikes/react-native/SPIKE.md`.

| Criterion                               |  Weight  |  Score (1-5)   | Justification / Notes |
| :-------------------------------------- | :------: | :------------: | :-------------------- |
| **Cross-Platform / Single Codebase**    |   20%    |       3        | One `App.js` codebase. Expo web rendered the ping body. The iOS simulator never booted (`simctl` exit 72, no available runtime), and Android was not run. |
| **Native Web/app Compatibility**        |   15%    |       3        | Web worked after adding `react-dom`, `react-native-web`, and `@expo/metro-runtime`. No native screen was shown. |
| **Developer Velocity & Ease of Coding** |   15%    |       4        | `create-expo-app` and a short `App.js` reached a visible `{"status":"ok"}` the same afternoon. The lost time was missing web packages, not the component. |
| **Ecosystem Maturity**                  |   15%    |       4        | The SDK 58 CLI and blank template installed and bundled on Node 26.13. The blank template still did not include web support. |
| **Performance & Speed**                 |   10%    |       3        | Metro bundled 203 modules in about 1.9s. No device startup time or frame-time measurement. |
| **Ease of AI Tools & UI Integration**   |   10%    |       2        | Not exercised. No component library and no model call in this spike. |
| **Accessibility (WCAG AA Compliance)**  |   15%    |       2        | Not exercised. The screen is one `Text` node. No contrast, label, or screen-reader check. |
| **Weighted Total**                      | **100%** | **3.05 / 5.0** | `0.20×3 + 0.15×3 + 0.15×4 + 0.15×4 + 0.10×3 + 0.10×2 + 0.15×2` |

---

## 2. Backend Technologies (Application Server)

Option scored: **NestJS 12**, from `spikes/nestjs/SPIKE.md`.

| Criterion                               |  Weight  |  Score (1-5)   | Justification / Notes |
| :-------------------------------------- | :------: | :------------: | :-------------------- |
| **OpenAPI / Swagger Generation**        |   15%    |       2        | Not exercised. No Swagger module was added, so this run produced no spec. |
| **OAuth 2.0 & OIDC Ecosystem**          |   15%    |       2        | Not exercised. No auth guard, Passport module, or token flow. |
| **Ease of AI Tools & LLM Integration**  |   10%    |       2        | Not exercised. No model client or tool call. |
| **Developer Velocity & Ease of Coding** |   10%    |       4        | `nest new` plus one `GET /ping` handler met the check the same afternoon. The sharp edge was `nodenext`: `./app.module` failed until the import used `.js`. |
| **Async Hook & Background Execution**   |   20%    |       3        | Bootstrap and the request path are async, and the app stayed up under `start:dev`. No queue, cron, or in-process background job was run. The MongoDB check was a separate script. |
| **Ecosystem Maturity**                  |   10%    |       4        | Nest 12.4.2, the CLI, watch mode, and Vitest all ran on Node 26.13. The generated unit test still expected `getHello()` and had to be pointed at `ping()`. |
| **Security**                            |   10%    |       3        | `PORT` comes from the environment and the process throws if it is missing. No `8080` literal in `src/`. `enableCors()` was then opened to `*` so Expo web could read the route. No authentication. |
| **Speed & performance**                 |   10%    |       3        | Watch compile and startup took about a second, and `curl` returned immediately. No load test. |
| **Weighted Total**                      | **100%** | **2.80 / 5.0** | `0.15×2 + 0.15×2 + 0.10×2 + 0.10×4 + 0.20×3 + 0.10×4 + 0.10×3 + 0.10×3` |

---

## 3. Database & Persistence Technologies

Option scored: **MongoDB Community 9.0.8**, from `spikes/mongodb/SPIKE.md`.

The three weights below are what the template lists. They sum to 60%, not 100%. The total is their weighted average, rescaled to /5.0: `(0.20×5 + 0.20×4 + 0.20×2) / 0.60`.

| Criterion                              |  Weight  |  Score (1-5)   | Justification / Notes |
| :------------------------------------- | :------: | :------------: | :-------------------- |
| **Dynamic Workflow / JSON Storage**    |   20%    |       5        | The script inserted `{ testedAt }` into `area.probes` with no schema and read that same value back. |
| **Developer Velocity & Ease of Use**   |   20%    |       4        | The driver script is short and the second run passed immediately. The lost time was Homebrew refusing the untrusted tap, then `mongosh` connecting before `mongod` was listening. |
| **Reliability & Relational Integrity** |   20%    |       2        | One insert and one read. No transaction, crash, or replica test. The spike also showed no relational constraints, which this database does not provide. |
| **Weighted Total**                     | **60%**  | **3.67 / 5.0** | Average of the three scores. The template has no further database criteria to make the weights reach 100%. |

---
