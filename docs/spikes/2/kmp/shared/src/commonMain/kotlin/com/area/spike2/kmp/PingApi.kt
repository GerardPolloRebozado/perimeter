package com.surfacevolume.spike2.kmp

import io.ktor.client.HttpClient
import io.ktor.client.HttpClientConfig
import io.ktor.client.call.body
import io.ktor.client.engine.HttpClientEngine
import io.ktor.client.plugins.HttpTimeout
import io.ktor.client.plugins.contentnegotiation.ContentNegotiation
import io.ktor.client.request.get
import io.ktor.serialization.kotlinx.json.json
import kotlinx.serialization.json.Json

/** Calls the backend. Shared by Android and web: only the HTTP engine differs (picked automatically). */
class PingApi(
    private val client: HttpClient = createHttpClient(),
    val baseUrl: String = backendBaseUrl(),
) {
    suspend fun ping(): PingResponse = client.get("$baseUrl/ping").body()
}

private val clientConfig: HttpClientConfig<*>.() -> Unit = {
    expectSuccess = true // non-2xx (e.g. 503 when the DB is down) throws instead of parsing
    install(ContentNegotiation) { json(Json { ignoreUnknownKeys = true }) }
    install(HttpTimeout) { requestTimeoutMillis = 5_000 }
}

/** `engine` is only passed by tests (mock engine). Otherwise Ktor uses the platform's engine. */
fun createHttpClient(engine: HttpClientEngine? = null): HttpClient =
    if (engine == null) HttpClient(clientConfig) else HttpClient(engine, clientConfig)
