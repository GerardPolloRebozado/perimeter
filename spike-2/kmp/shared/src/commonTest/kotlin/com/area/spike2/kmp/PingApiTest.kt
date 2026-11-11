package com.area.spike2.kmp

import io.ktor.client.engine.mock.MockEngine
import io.ktor.client.engine.mock.respond
import io.ktor.http.HttpHeaders
import io.ktor.http.HttpStatusCode
import io.ktor.http.headersOf
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertIs
import kotlinx.coroutines.test.runTest

/** Fake backend (MockEngine): checks how the shared code reacts to each answer. Runs without any server. */
class PingApiTest {

    private fun apiAnswering(status: HttpStatusCode, body: String) = PingApi(
        client = createHttpClient(MockEngine { request ->
            assertEquals("/ping", request.url.encodedPath)
            respond(body, status, headersOf(HttpHeaders.ContentType, "application/json"))
        }),
        baseUrl = "http://backend.test",
    )

    @Test
    fun okBodyGivesSuccess() = runTest {
        val state = loadPingState(apiAnswering(HttpStatusCode.OK, """{"status":"ok"}"""))
        assertEquals(PingUiState.Success("ok"), state)
    }

    @Test
    fun backendDownGivesError() = runTest {
        val state = loadPingState(apiAnswering(HttpStatusCode.ServiceUnavailable, """{"status":"error"}"""))
        assertIs<PingUiState.Error>(state)
    }

    @Test
    fun connectionFailureGivesError() = runTest {
        // Same failure type the browser engine raises when nothing listens on the port.
        val api = PingApi(
            client = createHttpClient(MockEngine { throw Error("Fail to fetch") }),
            baseUrl = "http://backend.test",
        )
        assertEquals(PingUiState.Error("Fail to fetch"), loadPingState(api))
    }

    @Test
    fun malformedBodyGivesError() = runTest {
        val state = loadPingState(apiAnswering(HttpStatusCode.OK, """not json"""))
        assertIs<PingUiState.Error>(state)
    }
}
