package com.surfacevolume.spike2.kmp

import kotlin.coroutines.cancellation.CancellationException

/** Everything the screen can show. */
sealed interface PingUiState {
    data object Loading : PingUiState
    data class Success(val status: String) : PingUiState
    data class Error(val message: String) : PingUiState
}

/** Calls the API and turns the result (or the failure) into a screen state. Kept out of the UI so it can be tested. */
suspend fun loadPingState(api: PingApi): PingUiState =
    try {
        PingUiState.Success(api.ping().status)
    } catch (e: CancellationException) {
        throw e // the screen went away: don't turn that into an error message
    } catch (e: Throwable) {
        // Throwable, not Exception: in the browser, Ktor reports "can't connect" as a kotlin.Error
        // ("Fail to fetch"), which is NOT an Exception. Catching only Exception left the screen stuck on "Checking…".
        PingUiState.Error(e.message ?: e::class.simpleName ?: "Unknown error")
    }
