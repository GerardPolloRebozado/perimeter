# Spike — MongoDB

- Timebox: not written down before the run. The bootstrap sizes a spike at 40 to 90 minutes.
- Why this length: no length was chosen ahead of time. The run was the database check from `spikes/SETUP.md`.
- Started: 2026-11-11 14:26, at `brew tap mongodb/brew`.
- Stopped: 2026-11-11 14:28, when the write/read script printed `status: ok`.

## Where I was when time ran out

MongoDB Community **9.0.8** is installed from `mongodb/brew` and the Homebrew service is running. `mongod` listens on `127.0.0.1:27017`. `mongosh` reports `{ ok: 1 }`.

`spikes/nestjs/server/scripts/mongo-check.mjs` inserted one document into `area.probes` and read that same `testedAt` back. The NestJS `/ping` body was still `{"status":"ok"}` afterward.

## What consumed the most minutes

Homebrew refused the tap. The first install stopped on untrusted `mongodb/brew/mongodb-community`. After that formula was trusted, the install stopped again on untrusted `mongodb-database-tools`. Trusting these three formulae was enough:

- `mongodb/brew/mongodb-community`
- `mongodb/brew/mongodb-database-tools`
- `mongodb/brew/mongosh`

The first `mongosh` then failed with `MongoNetworkError: connect ECONNREFUSED 127.0.0.1:27017` because it ran before `mongod` had opened the port. The service was already marked running. A second ping, a few seconds later, returned `ok: 1`.

The Node driver and the script itself were not the slow part.

## What I would need to go further

- Call MongoDB from a Nest injectable, with `MONGODB_URI` read the same way `PORT` is read.
- Keep `/ping` independent of the database, which this run already does.
- Homebrew printed that Xcode 26.2 is older than the version it wants. The bottled install still finished.

## Check

- [x] Official check passed
- Command I ran: `mongosh --eval 'db.runCommand({ ping: 1 })'`, then `MONGODB_URI='mongodb://127.0.0.1:27017/area' node scripts/mongo-check.mjs`, then `curl -sS http://localhost:8080/ping`
- Output I got: `{ ok: 1 }`, then `{"status":"ok","testedAt":"2026-11-11T13:28:26.076Z"}`, then `{"status":"ok"}`
