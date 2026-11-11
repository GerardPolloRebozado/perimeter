# Tech Stack Individual Option Evaluation Tables

---

## 1. Frontend Technologies

| Criterion                               |  Weight  |  Score (1-5)   | Justification / Notes                                                                               |
| :-------------------------------------- | :------: | :------------: | :-------------------------------------------------------------------------------------------------- |
| **Cross-Platform / Single Codebase**    |   20%    |       5        | React Native with Expo builds iOS, Android, and web from one JavaScript project.                    |
| **Native Web/app Compatibility**        |   15%    |       5        | The same project runs as a native iOS app, a native Android app, and a web app.                     |
| **Developer Velocity & Ease of Coding** |   15%    |       4        | React components are straightforward to write. TypeScript is easy to add and is not effortless.     |
| **Ecosystem Maturity**                  |   15%    |       5        | React Native, Expo, and TypeScript are widely used.                                                 |
| **Performance & Speed**                 |   10%    |       4        | Fast enough for typical screens. Heavy graphics are faster in Swift or Kotlin.                      |
| **Ease of AI Tools & UI Integration**   |   10%    |       4        | React has many ready-made UI libraries, and JavaScript has maintained AI SDKs that are easy to add. |
| **Accessibility (WCAG AA Compliance)**  |   15%    |       4        | React Native provides accessibility props, and web React can implement WCAG patterns.               |
| **Weighted Total**                      | **100%** | **4.50 / 5.0** | `0.20×5 + 0.15×5 + 0.15×4 + 0.15×5 + 0.10×4 + 0.10×4 + 0.15×4`                                      |

---

## 2. Backend Technologies (Application Server)

| Criterion                               |  Weight  |  Score (1-5)   | Justification / Notes                                                                                            |
| :-------------------------------------- | :------: | :------------: | :--------------------------------------------------------------------------------------------------------------- |
| **OpenAPI / Swagger Generation**        |   15%    |       5        | The official `@nestjs/swagger` package builds an OpenAPI document from controller decorators. It is easy to add. |
| **OAuth 2.0 & OIDC Ecosystem**          |   15%    |       4        | NestJS integrates with Passport, which provides OAuth 2.0 and OpenID Connect strategies.                         |
| **Ease of AI Tools & LLM Integration**  |   10%    |       4        | A NestJS service can call LLM APIs with the usual TypeScript SDKs.                                               |
| **Developer Velocity & Ease of Coding** |   10%    |       4        | A route is a TypeScript method in Nest's module layout. TypeScript is easy to use and is not effortless.         |
| **Async Hook & Background Execution**   |   20%    |       5        | Handlers are async. Official packages add cron tasks and job queues.                                             |
| **Ecosystem Maturity**                  |   10%    |       5        | NestJS is an established, widely used Node.js framework.                                                         |
| **Security**                            |   10%    |       4        | Guards, validation pipes, and common Node middleware cover the usual API checks.                                 |
| **Speed & performance**                 |   10%    |       4        | Node.js is enough for a normal API. It is not the fastest server runtime.                                        |
| **Weighted Total**                      | **100%** | **4.45 / 5.0** | `0.15×5 + 0.15×4 + 0.10×4 + 0.10×4 + 0.20×5 + 0.10×5 + 0.10×4 + 0.10×4`                                          |

---

## 3. Database & Persistence Technologies

**MongoDB: taken out directly.** The MongoDB server does not run on iOS. An iPhone app can still reach it through a backend. It is not scored as an on-device database.

| Criterion                              |  Weight  |  Score (1-5)   | Justification / Notes |
| :------------------------------------- | :------: | :------------: | :-------------------- |
| **Dynamic Workflow / JSON Storage**    |   20%    |                |                       |
| **Developer Velocity & Ease of Use**   |   20%    |                |                       |
| **Reliability & Relational Integrity** |   20%    |                |                       |
| **Weighted Total**                     | **100%** | **0.00 / 5.0** |                       |

---

Conclusion:

I worked previously with react and ts. They do work very well together, it is easy to see the structure of a project and changing/adding things is not hard. 

I know very well that ai integration is good. Vivecodabilty is too (I crafted acediploma.com with cursor only). 

Points and justification are mostly done by cursor because I didn't really know the answers to those questions.

Im not sure if react is the best frontend (many ppl like it and many hate it, but I didn't have issues with it), but nodejs looks really good and i want to learn more about it. 