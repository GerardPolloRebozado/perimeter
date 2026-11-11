# Spike — NestJS

- Timebox: not written down before the run. The bootstrap sizes a spike at 40 to 90 minutes.
- Why this length: no length was chosen ahead of time. The run was the server check from `spikes/SETUP.md`.
- Started: 2026-11-11, early afternoon, when `nest new server` was created.
- Stopped: 2026-11-11 14:32, after the checks below passed.

## Where I was when time ran out

The NestJS app in `spikes/nestjs/server` is listening with `PORT=8080`. `GET /ping` returns the required body. The port is read from the environment. Starting with `PORT` empty throws before the process listens. The generated unit test now calls `ping()` and passes.

`app.enableCors()` was added after the React Native web check. Without it the browser on port 8081 cannot read port 8080. `curl` still returns the same body.

## What consumed the most minutes

NestJS 12 compiles with `moduleResolution: nodenext`. The first `main.ts` imported `./app.module` with no extension, and `nest start --watch` failed with `TS2307: Cannot find module './app.module'`. The generated files already use a `.js` extension on local imports. The fix was `import { AppModule } from './app.module.js'`.

The dev server had also been stopped in the same terminal that was used for Homebrew, so it had to be started again before `curl` could pass.

## What I would need to go further

- Put MongoDB behind a Nest provider instead of a standalone script.
- Decide which origins `enableCors()` should allow before this is more than a spike.
- `AppService` is still registered and unused.

## Check

- [x] Official check passed
- Command I ran: `curl -sS http://localhost:8080/ping`, then `grep -R --include='*.ts' -n "8080" src`, then `PORT= npm run start`
- Output I got: `{"status":"ok"}`, then `no hardcoded port`, then `Error: PORT environment variable is required` and exit code 1. `npm test` — 1 file, 1 test, passed.
