# Setup and spike tests — NestJS, React Native, MongoDB

This branch is for testing three candidates:

| Layer | Technology | What “it works” means |
| --- | --- | --- |
| Server | NestJS | `GET /ping` on port 8080 returns exactly `{"status":"ok"}` |
| Database | MongoDB | The NestJS app can write one document and read it back |
| Client | React Native (Expo) | The app shows the body returned by `/ping` |

The AREA bootstrap (`area-bootstrap`, v1.1.1) asks for **timeboxed spikes**, not finished features. A spike buys information. When the timer rings, you stop, even if it is not working.

The brief asks for three **backend** technologies run against the same target. This guide covers the three layers above. NestJS is the HTTP server. MongoDB and React Native are tested against that same server so the three runs stay comparable. NestJS runs on Node.js. Node.js is the runtime, not the candidate being scored.

Do the NestJS spike first. The other two need the server.

---

## What you must produce

Each spike is 40 to 90 minutes, run on its own. Before you start a spike, write down the timebox length and why you picked it. Do not extend the timebox because you like the tool.

When the timer rings, stop and immediately write the note for that technology. Commit the note. A decision you cannot show is a decision you did not make.

| File | When | What it must contain |
| --- | --- | --- |
| `spikes/nestjs/SPIKE.md` | End of the NestJS timebox | Where you were when time ran out, what ate the most minutes, what you would need to go further |
| `spikes/mongodb/SPIKE.md` | End of the MongoDB timebox | Same three points |
| `spikes/react-native/SPIKE.md` | End of the React Native timebox | Same three points |
| `spikes/DISTANCE.md` | After all three | A table comparing the three runs. You define what “distance travelled” means, and you justify that definition |

Suggested columns for `spikes/DISTANCE.md`:

| Technology | Timebox | Reached the check? | Minutes lost to setup | Blocker when the timer rang | Distance travelled |
| --- | --- | --- | --- | --- | --- |
| NestJS | | | | | |
| MongoDB | | | | | |
| React Native | | | | | |

“Distance travelled” is yours to define. Write the definition above the table. Example you can adopt or replace: *how far the run got along setup → process running → official check passed → you could explain the failure without reading the docs again.*

After the spikes, fill `spikes/template.md` (the weighted matrix) from what you actually saw. Scores filled in before you run anything are invented.

---

## Machine already in place

Checked on this machine:

- Node.js **26.13.0** (`/opt/homebrew/bin/node`), the runtime NestJS needs
- npm **12.2.2**
- NestJS CLI **12.0.12** (`nest`)

Confirm before you start:

```bash
node -v
npm -v
nest --version
```

---

## 1. NestJS server

Target from the bootstrap, all three of these:

- The process listens on port **8080**.
- `GET /ping` answers with the exact body `{"status":"ok"}`.
- The port comes from the environment variable `PORT`. There must be no `8080` literal in the source. If `PORT` is missing, the process exits instead of picking a default.

### Create the app

From the repo root:

```bash
mkdir -p spikes/nestjs
cd spikes/nestjs
nest new server --package-manager npm --skip-git --strict
cd server
```

If the CLI still asks a question, accept npm and strict TypeScript.

### Read the port from the environment

Replace `server/src/main.ts` with:

```typescript
import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module.js';

async function bootstrap() {
  const port = process.env.PORT;
  if (!port) {
    throw new Error('PORT environment variable is required');
  }

  const app = await NestFactory.create(AppModule);
  app.enableCors();
  await app.listen(port);
}

bootstrap();
```

### Answer `GET /ping`

In `server/src/app.controller.ts`, replace the controller with:

```typescript
import { Controller, Get } from '@nestjs/common';

@Controller()
export class AppController {
  @Get('ping')
  ping(): { status: string } {
    return { status: 'ok' };
  }
}
```

You can leave `AppService` unused, or delete it and drop it from `app.module.ts`. The ping route must not depend on it.

### Run and check

```bash
PORT=8080 npm run start:dev
```

In another terminal, from anywhere:

```bash
curl -sS http://localhost:8080/ping
```

The body must be exactly:

```text
{"status":"ok"}
```

Then confirm the port is not hardcoded:

```bash
grep -R --include='*.ts' -n "8080" src || echo "no hardcoded port"
PORT= npm run start
```

The second command must fail with `PORT environment variable is required`. Stop that process. Keep the first server running for the MongoDB and React Native checks, or start it again with `PORT=8080`.

### Pass / fail

Pass only if all of these are true:

- `curl` prints `{"status":"ok"}` and nothing else.
- Nothing in `src/` contains `8080`.
- Starting without `PORT` throws.

Write `spikes/nestjs/SPIKE.md` as soon as the timer rings, whether or not it passed.

---

## 2. MongoDB

The official `/ping` body stays `{"status":"ok"}`. Do not fold the database into that response. Prove MongoDB with a separate check.

### Install and start

Docker, if you already use it:

```bash
docker run -d --name area-mongo -p 27017:27017 mongo:7
```

Homebrew otherwise. Current Homebrew refuses this tap until you trust the formulae:

```bash
brew tap mongodb/brew
brew trust --formula mongodb/brew/mongodb-community
brew trust --formula mongodb/brew/mongodb-database-tools
brew trust --formula mongodb/brew/mongosh
brew install mongodb-community
brew services start mongodb-community
```

Check the server itself:

```bash
mongosh --eval 'db.runCommand({ ping: 1 })'
```

You want `ok: 1`.

### Prove NestJS can use it

Still inside `spikes/nestjs/server`:

```bash
npm install mongodb
```

Save this as `spikes/nestjs/server/scripts/mongo-check.mjs`:

```javascript
import { MongoClient } from 'mongodb';

const uri = process.env.MONGODB_URI;
if (!uri) {
  throw new Error('MONGODB_URI environment variable is required');
}

const client = new MongoClient(uri);

try {
  await client.connect();
  const db = client.db('area');
  const probes = db.collection('probes');
  const testedAt = new Date().toISOString();

  await probes.insertOne({ testedAt });
  const found = await probes.findOne({ testedAt });

  if (!found) {
    throw new Error('inserted document was not read back');
  }

  console.log(JSON.stringify({ status: 'ok', testedAt: found.testedAt }));
} finally {
  await client.close();
}
```

Run it:

```bash
MONGODB_URI='mongodb://127.0.0.1:27017/area' node scripts/mongo-check.mjs
```

Pass only if it prints a line like:

```text
{"status":"ok","testedAt":"2026-11-11T12:00:00.000Z"}
```

The `testedAt` value will differ. `status` must be `ok`, and that timestamp must be a document you just wrote.

Then confirm `/ping` is unchanged:

```bash
curl -sS http://localhost:8080/ping
```

Still exactly `{"status":"ok"}`.

Write `spikes/mongodb/SPIKE.md` when the timer rings.

To stop MongoDB later:

```bash
docker stop area-mongo
```

or, with Homebrew:

```bash
brew services stop mongodb-community
```

---

## 3. React Native

Use Expo. A bare React Native project usually spends the whole timebox on native toolchains (Xcode, Android Studio, CocoaPods).

### Create the app

From the repo root, in a new terminal:

```bash
mkdir -p spikes/react-native
cd spikes/react-native
npx create-expo-app@latest mobile
cd mobile
npx expo start
```

Use the iOS simulator if you are on this Mac. In the simulator, `localhost` is your machine, so the app can call `http://localhost:8080/ping`. This Mac currently has Xcode installed, but `xcode-select` points at the Command Line Tools and no iOS simulator runtime is available, so `simctl` cannot boot a phone. Expo web at `http://localhost:8081` is the check that can actually run here. The blank template has no web packages, so install them first:

```bash
npx expo install react-dom react-native-web @expo/metro-runtime
```

A browser page on port 8081 cannot read port 8080 unless the NestJS app calls `app.enableCors()`. The iOS simulator does not need that. The snippet in section 1 already includes it so both clients can run.

A physical phone cannot use `localhost`. It must use your computer’s LAN address, for example `http://192.168.1.20:8080/ping`, and the phone and the computer must be on the same network. An Android emulator uses `http://10.0.2.2:8080/ping` instead of `localhost`.

### Show the ping result

Replace `App.js` or `App.tsx` (whichever the generator created) with a screen that loads `/ping` on startup and renders the raw body. Example for `App.tsx`:

```tsx
import { useEffect, useState } from 'react';
import { SafeAreaView, Text } from 'react-native';

const PING_URL = process.env.EXPO_PUBLIC_PING_URL;

export default function App() {
  const [body, setBody] = useState('loading');

  useEffect(() => {
    if (!PING_URL) {
      setBody('EXPO_PUBLIC_PING_URL is required');
      return;
    }

    fetch(PING_URL)
      .then((response) => response.text())
      .then((text) => setBody(text))
      .catch((error: unknown) => {
        setBody(error instanceof Error ? error.message : 'request failed');
      });
  }, []);

  return (
    <SafeAreaView>
      <Text>{body}</Text>
    </SafeAreaView>
  );
}
```

Start it with the URL in the environment, not hardcoded:

```bash
EXPO_PUBLIC_PING_URL='http://localhost:8080/ping' npx expo start
```

The NestJS server from section 1 must be running with `PORT=8080`.

### Pass / fail

Pass only if the screen shows:

```text
{"status":"ok"}
```

If it shows `Network request failed`, you are calling the wrong host for that simulator or device. Fix the URL and note how long that took. That cost belongs in `spikes/react-native/SPIKE.md`.

---

## Order to run the session

1. Write the three timebox lengths and the reason for each length.
2. Start the NestJS timer. Reach the `curl` check. Stop when it rings. Write `spikes/nestjs/SPIKE.md`.
3. Start MongoDB. Start the MongoDB timer. Reach `mongo-check.mjs`. Stop. Write `spikes/mongodb/SPIKE.md`.
4. Start the React Native timer. Reach the screen that shows the ping body. Stop. Write `spikes/react-native/SPIKE.md`.
5. Write `spikes/DISTANCE.md`, including your definition of distance.
6. Only then score `spikes/template.md`.

Do not keep coding after the timer to “just finish”. That extra time is the bias the exercise is there to expose. If you finish early, stop anyway and write what you would need in order to trust the technology in the real project.

---

## SPIKE.md skeleton

Copy this into each `SPIKE.md` and fill it when the timer rings.

```markdown
# Spike — <NestJS | MongoDB | React Native>

- Timebox: <minutes>
- Why this length:
- Started:
- Stopped:

## Where I was when time ran out

## What consumed the most minutes

## What I would need to go further

## Check

- [ ] Official check passed
- Command I ran:
- Output I got:
```
