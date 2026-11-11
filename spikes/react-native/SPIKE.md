# Spike — React Native

- Timebox: not written down before the run. The bootstrap sizes a spike at 40 to 90 minutes.
- Why this length: no length was chosen ahead of time. The run was the client check from `spikes/SETUP.md`.
- Started: 2026-11-11 14:28, at `npx create-expo-app@latest mobile --template blank`.
- Stopped: 2026-11-11 14:31, when the Expo web page text was `{"status":"ok"}`.

## Where I was when time ran out

The app is `spikes/react-native/mobile`, Expo SDK 58, blank template. `App.js` loads `EXPO_PUBLIC_PING_URL` and renders the response body.

The iOS simulator did not run. `xcode-select` points at `/Library/Developer/CommandLineTools`, `xcrun simctl` exits 72, and the only simulator runtimes Xcode lists are unavailable. Expo itself logged `Unable to run simctl`.

Expo web is running at `http://localhost:8081` with `EXPO_PUBLIC_PING_URL=http://localhost:8080/ping`. The page text is `{"status":"ok"}`.

## What consumed the most minutes

Two gaps in the written steps, not the component:

- The blank template does not include web support. `react-dom`, `react-native-web`, and `@expo/metro-runtime` had to be installed before `expo start --web` could render.
- The first page text was `Failed to fetch`. NestJS was not sending `Access-Control-Allow-Origin`, so the browser on port 8081 could not read port 8080. `app.enableCors()` fixed that. A React Native simulator would not have hit this, because that fetch is not a browser request.

## What I would need to go further

- Point `xcode-select` at `/Applications/Xcode.app` and install a current iOS simulator runtime, then run `npm run ios` and confirm the same body on a phone screen.
- Restrict CORS to the origins the project actually uses.
- On a physical device, replace `localhost` with the computer's LAN address.

## Check

- [x] Official check passed on Expo web
- [ ] iOS simulator check not run
- Command I ran: `EXPO_PUBLIC_PING_URL='http://localhost:8080/ping' npx expo start --web --port 8081`, then opened `http://localhost:8081`
- Output I got: page text `{"status":"ok"}` after CORS was enabled. Before that, the page text was `Failed to fetch`.
