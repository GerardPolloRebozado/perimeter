plugins {
    kotlin("multiplatform")
    id("org.jetbrains.dokka")
}

repositories {
    mavenCentral()
}

kotlin {
    jvm()
    iosX64()
    iosArm64()
    iosSimulatorArm64()
    wasmJs()
    
    sourceSets {
        val commonMain by getting {
            dependencies {
                // Import the generated API client
                implementation(project(":api-client"))
                
                // Ktor HTTP client for calling NestJS backend
                implementation("io.ktor:ktor-client-core:2.0.3")
                implementation("io.ktor:ktor-client-content-negotiation:2.0.3")
                implementation("io.ktor:ktor-client-serialization:2.0.3")
                implementation("io.ktor:ktor-serialization-kotlinx-json:2.0.3")
                
                // Kotlinx serialization for JSON
                implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.6.2")
            }
        }
        val androidMain by getting {
            dependsOn(commonMain)
        }
        val iosX64Main by getting {
            dependsOn(commonMain)
        }
        val iosArm64Main by getting {
            dependsOn(commonMain)
        }
        val iosSimulatorArm64Main by getting {
            dependsOn(commonMain)
        }
        val wasmJsMain by getting {
            dependsOn(commonMain)
        }
        
        val commonTest by getting {
            dependencies {
                implementation(kotlin("test"))
            }
        }
    }
}

publishing {
    publications {
        create<MavenPublication>("shared") {
            from(components["jvm"])
        }
    }
}
