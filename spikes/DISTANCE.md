# Distance

Distance travelled means how far a run got along this path:

1. The process is installed and starts.
2. The official check passed.
3. The failure, if there was one, can be explained without opening the docs again.

A run that only installed the tool travelled less than a run whose check passed. A passed check with an unexplained failure along the way travelled less than a passed check whose blocker is written down. The scores below use that definition. No timebox length was locked before these runs. The bootstrap band for each spike is 40 to 90 minutes.

| Technology | Timebox | Reached the check? | Minutes lost to setup | Blocker when the timer rang | Distance travelled |
| --- | --- | --- | --- | --- | --- |
| NestJS | not chosen beforehand | yes | the missing `.js` import under `nodenext`, and restarting the dev server | none. `curl` returned `{"status":"ok"}` | Server running, check passed, and the import failure is explained |
| MongoDB | not chosen beforehand | yes | Homebrew refusing the untrusted `mongodb/brew` formulae | none. The write/read script returned `status: ok` | Database running, check passed, and the trust plus early-connection failure are explained |
| React Native | not chosen beforehand | yes on Expo web. iOS simulator did not run | blank template had no web packages, and no usable iOS simulator runtime | none on web. The simulator never started | Web check passed after CORS. The simulator gap is explained and still open |

The weighted matrix is filled in `spikes/template.md` from these runs. Criteria the spike did not exercise are scored 2, and the note says so.
