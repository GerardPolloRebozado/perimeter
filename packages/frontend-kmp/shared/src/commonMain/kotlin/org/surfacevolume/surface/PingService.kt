package org.surfacevolume.surface

import io.ktor.client.*
import io.ktor.client.request.*
import io.ktor.client.statement.*

// Single implementation for ALL platforms (JVM, Android, iOS, JS, Wasm)
// Ktor HttpClient() auto-detects the right engine from the platform classpath:
object PingService {
    private val client = HttpClient {
        expectSuccess = true
    }

    suspend fun ping(): String {
        return client.get("${ApiConfig.baseUrl}/ping").bodyAsText()
    }
}
