package org.area.perimeter

// JVM/Desktop: read from system property or environment variable
actual object ApiConfig {
    actual val baseUrl: String
        get() = System.getProperty("AREA_API_BASE_URL")
            ?: System.getenv("AREA_API_BASE_URL")
            ?: "http://localhost:8080"
}
