import org.jetbrains.kotlin.gradle.ExperimentalWasmDsl
import org.jetbrains.kotlin.gradle.dsl.JvmTarget

plugins {
    alias(libs.plugins.kotlinMultiplatform)
    alias(libs.plugins.androidMultiplatformLibrary)
    alias(libs.plugins.composeMultiplatform)
    alias(libs.plugins.composeCompiler)
}

kotlin {
    listOf(
        iosArm64(),
        iosSimulatorArm64()
    ).forEach { iosTarget ->
        iosTarget.binaries.framework {
            baseName = "Shared"
            isStatic = true
        }
    }
    jvm()
    js {
        browser()
        binaries.executable()
    }
    @OptIn(ExperimentalWasmDsl::class)
    wasmJs {
        browser()
        binaries.executable()
    }
    android {
       namespace = "org.surfacevolume.surface.shared"
       compileSdk = libs.versions.android.compileSdk.get().toInt()
       minSdk = libs.versions.android.minSdk.get().toInt()
       compilerOptions {
           jvmTarget = JvmTarget.JVM_11
       }
       androidResources {
           enable = true
       }
       withHostTest {
           isIncludeAndroidResources = true
       }
       withDeviceTestBuilder {
           sourceSetTreeName = "test"
       }.configure {
           instrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
       }
    }
    sourceSets {
        commonMain.dependencies {
            implementation(libs.compose.runtime)
            implementation(libs.compose.foundation)
            implementation(libs.compose.material3)
            implementation(libs.compose.ui)
            implementation(libs.compose.components.resources)
            implementation(libs.compose.uiToolingPreview)
            implementation(libs.androidx.lifecycle.viewmodelCompose)
            implementation(libs.androidx.lifecycle.runtimeCompose)

            // Ktor core - multiplatform artifact, works on all targets
            implementation("io.ktor:ktor-client-core:3.6.1")
            implementation("io.ktor:ktor-client-content-negotiation:3.6.1")
            implementation("io.ktor:ktor-client-serialization:3.6.1")
            implementation("io.ktor:ktor-serialization-kotlinx-json:3.6.1")
            implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.7.0")
        }
        commonTest.dependencies {
            implementation(libs.kotlin.test)
        }
        androidMain.dependencies {
            implementation(libs.compose.uiToolingPreview)
            implementation(libs.compose.uiTooling)
            // OkHttp engine for Android
            implementation("io.ktor:ktor-client-okhttp:3.6.1")
        }
        jvmMain.dependencies {
            // OkHttp engine for JVM desktop
            implementation("io.ktor:ktor-client-okhttp:3.6.1")
        }
        iosArm64Main.dependencies {
            // Darwin engine for iOS device
            implementation("io.ktor:ktor-client-darwin:3.6.1")
        }
        iosSimulatorArm64Main.dependencies {
            // Darwin engine for iOS simulator
            implementation("io.ktor:ktor-client-darwin:3.6.1")
        }
        jsMain.dependencies {
            implementation(libs.wrappers.browser)
            // JS engine for web
            implementation("io.ktor:ktor-client-js:3.6.1")
        }
        wasmJsMain.dependencies {
            implementation(libs.wrappers.browser)
            // JS engine for Wasm (uses browser fetch)
            implementation("io.ktor:ktor-client-js:3.6.1")
        }
        getByName("androidDeviceTest").dependencies {
            implementation(libs.androidx.testRunner)
        }
    }
}

dependencies {
    androidRuntimeClasspath(libs.compose.uiTooling)
}
