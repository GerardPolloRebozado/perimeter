package com.surfacevolume.spike2.kmp

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.safeContentPadding
import androidx.compose.foundation.layout.size
import androidx.compose.material3.Button
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableIntStateOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.semantics.LiveRegionMode
import androidx.compose.ui.semantics.heading
import androidx.compose.ui.semantics.liveRegion
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.unit.dp

@Composable
fun PingScreen(api: PingApi) {
    var state by remember { mutableStateOf<PingUiState>(PingUiState.Loading) }
    var attempt by remember { mutableIntStateOf(0) }

    // Runs on first display and again every time `attempt` changes (the button).
    LaunchedEffect(attempt) {
        state = PingUiState.Loading
        state = loadPingState(api)
    }

    Column(
        modifier = Modifier.fillMaxSize().safeContentPadding().padding(24.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(16.dp, Alignment.CenterVertically),
    ) {
        Text(
            "Surface Volume backend status",
            style = MaterialTheme.typography.headlineSmall,
            modifier = Modifier.semantics { heading() },
        )

        // liveRegion: screen readers announce the new text when the status changes.
        Row(
            modifier = Modifier.semantics { liveRegion = LiveRegionMode.Polite },
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(8.dp),
        ) {
            when (val s = state) {
                PingUiState.Loading -> {
                    CircularProgressIndicator(Modifier.size(20.dp), strokeWidth = 2.dp)
                    Text("Checking…")
                }
                is PingUiState.Success -> Text(
                    "Status: ${s.status}",
                    style = MaterialTheme.typography.titleLarge,
                    color = MaterialTheme.colorScheme.primary,
                )
                is PingUiState.Error -> Text(
                    "Backend unreachable: ${s.message}",
                    color = MaterialTheme.colorScheme.error,
                )
            }
        }

        Button(onClick = { attempt++ }, enabled = state != PingUiState.Loading) {
            Text("Check again")
        }

        Text("Backend: ${api.baseUrl}", style = MaterialTheme.typography.bodySmall)
    }
}
