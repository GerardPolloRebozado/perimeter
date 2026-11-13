package com.surfacevolume.spike2.kmp

import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.runtime.Composable
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier

/** Shared root of the app: the exact same code runs on Android and in the browser. */
@Composable
fun App() {
    MaterialTheme {
        Surface(Modifier.fillMaxSize()) {
            val api = remember { PingApi() }
            PingScreen(api)
        }
    }
}
