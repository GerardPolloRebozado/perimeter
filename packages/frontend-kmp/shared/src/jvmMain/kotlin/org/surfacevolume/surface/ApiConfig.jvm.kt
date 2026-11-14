package org.surfacevolume.surface

// JVM/Desktop: read from system property or environment variable
actual object ApiConfig {
    actual val baseUrl: String
        get() = System.getProperty("SURFACE_VOLUME_API_BASE_URL")
            ?: System.getenv("SURFACE_VOLUME_API_BASE_URL")
            ?: "http://localhost:8080"
}
