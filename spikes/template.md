# Tech Stack Individual Option Evaluation Tables

---

## 1. Frontend Technologies: Kotlin Multiplatform + Compose

| Criterion                               |  Weight  |  Score (1-5)   | Justification / Notes |
| :-------------------------------------- | :------: | :------------: | :-------------------- |
| **Cross-Platform / Single Codebase**    |   20%    |       4        | Almost all the code is shared between web and Android. Each platform still had its own bugs to fix. |
| **Native Web/app Compatibility**        |   15%    |       3        | Feels native on Android. On the web it's drawn like an image instead of a normal web page. |
| **Developer Velocity & Ease of Coding** |   15%    |       3        | Writing the screen was quick. Project setup and build times were slow. |
| **Ecosystem Maturity**                  |   15%    |       3        | Mature on Android. The web version is still new, with few examples. |
| **Performance & Speed**                 |   10%    |       3        | Smooth on Android. The web version is heavy to load for a single screen. |
| **Ease of AI Tools & UI Integration**   |   10%    |       3        | AI tools know Kotlin and Android well. They know much less about the web version. |
| **Accessibility (WCAG AA Compliance)**  |   15%    |       2        | Works on Android. On the web, screen reader support is uncertain and wasn't tested. |
| **Weighted Total**                      | **100%** | **3.05 / 5.0** | Good on Android; the web version is the weak side. |

<sub>Calculation: 0.20×4 + 0.15×3 + 0.15×3 + 0.15×3 + 0.10×3 + 0.10×3 + 0.15×2 = 3.05</sub>

---

## 2. Backend Technologies (Application Server): Java + Spring Boot

| Criterion                               |  Weight  |  Score (1-5)   | Justification / Notes |
| :-------------------------------------- | :------: | :------------: | :-------------------- |
| **OpenAPI / Swagger Generation**        |   15%    |       4        | API documentation can be generated with one add-on. Not tested. |
| **OAuth 2.0 & OIDC Ecosystem**          |   15%    |       5        | Login with Google, GitHub, etc. is built in. It's hard to configure. Not tested. |
| **Ease of AI Tools & LLM Integration**  |   10%    |       4        | Good AI libraries, and AI tools know Spring well. Many examples are for the older version. |
| **Developer Velocity & Ease of Coding** |   10%    |       4        | The target took very little code and a couple of minutes. |
| **Async Hook & Background Execution**   |   20%    |       4        | Running tasks in the background is built in. Not tested. |
| **Ecosystem Maturity**                  |   10%    |       5        | Huge community, and an official add-on exists for almost everything. Adding the database was easy. |
| **Security**                            |   10%    |       4        | Mature and widely used security tools. They're easy to misconfigure. |
| **Speed & performance**                 |   10%    |       4        | Starts in a few seconds. Uses more memory than lighter options, which is fine for AREA. |
| **Weighted Total**                      | **100%** | **4.25 / 5.0** | Strongest of the three, though most of what AREA needs wasn't tested yet. |

<sub>Calculation: 0.15×4 + 0.15×5 + 0.10×4 + 0.10×4 + 0.20×4 + 0.10×5 + 0.10×4 + 0.10×4 = 4.25</sub>

---

## 3. Database & Persistence Technologies: Cassandra


| Criterion                              |  Weight  |  Score (1-5)   | Justification / Notes |
| :------------------------------------- | :------: | :------------: | :-------------------- |
| **Dynamic Workflow / JSON Storage**    |   20%    |       2        | Workflow settings can be saved and loaded. You can't search inside them. |
| **Developer Velocity & Ease of Use**   |   20%    |       2        | Every new kind of search needs a new table planned in advance. It's also slow to start and uses a lot of memory. |
| **Reliability & Relational Integrity** |   20%    |       2        | It recovers well when it goes down and comes back. Nothing keeps related data consistent; our code has to. |
| **Weighted Total**                     | **60%** | **2.00 / 5.0** | Poor fit: it's built for huge scale AREA doesn't need, and gives up the searching and consistency AREA does need. |

<sub>Calculation: (0.20×2 + 0.20×2 + 0.20×2) / 0.60 = 1.20 / 0.60 = 2.00</sub>

---
