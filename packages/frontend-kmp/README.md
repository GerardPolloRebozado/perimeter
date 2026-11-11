# AREA Frontend — Kotlin Multiplatform Setup

## What you need to install

### 1. JDK 17 (required)

KMP targets JVM and needs Java 17+.

| OS | Install command |
|---|---|
| macOS | `brew install --cask temurin@17` |
| Linux (deb/rpm) | `sudo apt install openjdk-17-jdk` |
| Windows | Download from [Adoptium](https://adoptium.net/) |

Verify: `java -version` → should show `17.x`

### 2. Android Studio (for Android target)

- Download: https://developer.android.com/studio
- Install with: Android SDK, Android SDK Command-line Tools, Build Tools
- Set environment variable: `ANDROID_HOME` (or `ANDROID_SDK_ROOT`)

Verify: `which android` or `echo $ANDROID_HOME`

### 3. Gradle

The project uses the Gradle Wrapper (`gradlew` / `gradlew.bat`), so you don't need to install Gradle globally.

### 4. Web / Wasm target (optional)

- A modern browser with WasmGC support (Chrome 119+, Firefox 120+, Safari 17.4+)
- Node.js + pnpm (only needed if you run the web dev server)

### 5. iOS target (optional, macOS only)

- Xcode 15+ from the Mac App Store
- For simulator testing: `xcode-select --install`

---

## Project structure

```
packages/frontend-kmp/
├── settings.gradle.kts          # Module declarations
├── build.gradle.kts             # Root build config
├── shared/                      # Shared KMP module (JVM + Android + Wasm)
│   ├── build.gradle.kts
│   └── src/
│       ├── commonMain/          # Code shared by all platforms
│       │   └── kotlin/area/shared/
│       │       ├── ApiClientFactory.kt   # Creates typed API client from OpenAPI
│       │       └── WorkflowService.kt    # Business logic using generated client
│       ├── androidMain/         # Android-specific code
│       ├── wasmJsMain/          # Web/Wasm-specific code
│       └── commonTest/          # Shared tests
└── app/                         # Android app module
    ├── build.gradle.kts
    └── src/main/kotlin/area/app/
```

---

## Key dependencies

| Dependency | Version | Purpose | Docs |
|---|---|---|---|
| Kotlin Multiplatform | 1.9.0 | Cross-platform language | https://kotlinlang.org/docs/multiplatform.html |
| Kotlin Serialization JSON | 1.6.2 | JSON (de)serialization | https://kotlinlang.org/docs/serialization.html |
| Ktor Client Core | 2.0.3 | HTTP client | https://ktor.io/docs/client-create-client.html |
| Ktor Serialization | 2.0.3 | JSON client plugin | https://ktor.io/docs/client-serialization.html |
| OpenAPI Generator | 6.6.0 | Generates Kotlin client from OpenAPI | https://openapi-generator.tech/docs/generators/kotlin |

---

## Running the project

### Shared module (JVM tests only, no device needed)

```bash
cd packages/frontend-kmp
./gradlew :shared:test
```

### Android app

```bash
cd packages/frontend-kmp
./gradlew :app:installDebug   # installs on connected emulator/device
./gradlew :app:uninstallDebug # uninstall
```

### Web / Wasm

```bash
cd packages/frontend-kmp
./gradlew :shared:wasmJsBrowserDevelopmentRun
# Opens http://localhost:8080 (dev server reloads on changes)
```

---

## Type sharing workflow (NestJS ↔ KMP)

1. **Backend** (`packages/backend-nestjs/`): NestJS controllers use `@nestjs/swagger` decorators on DTOs
2. **OpenAPI JSON**: `packages/backend-nestjs/openapi.json` is written by `src/main.ts` on startup
3. **Generation**: `scripts/generate-api-client.sh` runs OpenAPI Generator → produces Kotlin client in `packages/api-client/`
4. **Frontend** (`packages/frontend-kmp/shared/`): imports `@area/api-client` and calls typed methods

```kotlin
// Generated client usage (type-safe)
val api = ApiClientFactory.create()        // DefaultApi from OpenAPI
val workflows: List<WorkflowDto> = api.getWorkflows()  // Kotlin data class from OpenAPI
```

---

## Documentation links

| Topic | Link |
|---|---|
| Kotlin Multiplatform | https://kotlinlang.org/docs/multiplatform.html |
| KMP setup | https://kotlinlang.org/docs/setup-kmm.html |
| Ktor Client | https://ktor.io/docs/client-quick-start.html |
| Kotlin Serialization | https://kotlinlang.org/docs/serialization.html |
| OpenAPI Generator Kotlin | https://openapi-generator.tech/docs/generators/kotlin |
| Android Studio | https://developer.android.com/studio |
| JDK 17 (Adoptium) | https://adoptium.net/ |

---

## Troubleshooting

**`ANDROID_HOME` not found**
```
export ANDROID_HOME="$HOME/Android/Sdk"
export PATH="$ANDROID_HOME/tools:$ANDROID_HOME/platform-tools:$PATH"
```

**Gradle daemon memory issues**
Create `gradle.properties` in `packages/frontend-kmp/`:
```
org.gradle.jvmargs=-Xmx2048m -XX:MaxMetaspaceSize=512m
```

**Wasm build fails**
Make sure your browser supports WasmGC (Chrome 119+, Firefox 120+).

---

## Next steps (to improve)

- [ ] Add iOS target configuration (macOS only)
- [ ] Add Compose Multiplatform UI dependencies
- [ ] Configure API client version pinning in `gradle/libs.versions.toml`
- [ ] Add CI workflow for KMP builds
- [ ] Add shared database/models module
