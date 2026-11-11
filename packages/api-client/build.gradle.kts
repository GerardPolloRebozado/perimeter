plugins {
    `kotlin`
    `maven-publish`
}

repositories {
    mavenCentral()
}

dependencies {
    // OpenAPI Generator will add necessary dependencies here
    // When generated, the build.gradle.kts will be replaced/merged with generated content
}

kotlin {
    jvm()
    
    sourceSets {
        val commonMain by getting {
            dependencies {
                implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.6.2")
                implementation("io.ktor:ktor-client-core:2.0.3")
                implementation("io.ktor:ktor-client-content-negotiation:2.0.3")
                implementation("io.ktor:ktor-client-serialization:2.0.3")
                implementation("io.ktor:ktor-serialization-kotlinx-json:2.0.3")
            }
        }
    }
}

publishing {
    publications {
        create<MavenPublication>("areaApi") {
            groupId = "com.area"
            artifactId = "area-api-client"
            version = "1.0.0"
        }
    }
}
