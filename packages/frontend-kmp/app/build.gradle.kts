plugins {
    kotlin("multiplatform")
    id("com.android.application") version "8.2.0" apply false
    id("org.jetbrains.dokka")
}

android {
    namespace = "area.app"
    compileSdk = 34

    defaultConfig {
        applicationId = "area.app"
        minSdk = 24
        targetSdk = 34
        versionCode = 1
        versionName = "1.0"
    }

    buildTypes {
        release {
            isMinifyEnabled = false
        }
    }

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_17
        targetCompatibility = JavaVersion.VERSION_17
    }
}

kotlin {
    android()
    
    sourceSets {
        val commonMain by getting {
            dependsOn(projects.shared)
        }
        val androidMain by getting {
            dependsOn(commonMain)
        }
    }
}

repositories {
    mavenCentral()
    google()
    mavenCentral()
}

dependencies {
    implementation(project(":shared"))
}
