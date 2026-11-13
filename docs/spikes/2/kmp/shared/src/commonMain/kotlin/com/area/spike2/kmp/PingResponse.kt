package com.surfacevolume.spike2.kmp

import kotlinx.serialization.Serializable

/** Body of the backend's GET /ping: {"status":"ok"}. */
@Serializable
data class PingResponse(val status: String)
