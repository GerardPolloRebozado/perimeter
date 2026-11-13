package com.surfacevolume.spike2.kmp

import android.content.pm.PackageManager
import android.os.Build
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.runtime.mutableStateOf

/** Android entry point: shows the shared App(). */
class MainActivity : ComponentActivity() {

    // Android 17+ blocks connections to the local network (our dev backend on 10.0.2.2)
    // unless the user grants this permission at runtime.
    private val localNetworkPermission = "android.permission.ACCESS_LOCAL_NETWORK"
    private val permissionAnswered = mutableStateOf(false)

    private val askPermission =
        registerForActivityResult(ActivityResultContracts.RequestPermission()) {
            // Granted or refused: show the app either way. If refused, the screen shows the connection error.
            permissionAnswered.value = true
        }

    override fun onCreate(savedInstanceState: Bundle?) {
        enableEdgeToEdge()
        super.onCreate(savedInstanceState)

        val needsPermission = Build.VERSION.SDK_INT >= 37 &&
            checkSelfPermission(localNetworkPermission) != PackageManager.PERMISSION_GRANTED
        if (needsPermission) askPermission.launch(localNetworkPermission) else permissionAnswered.value = true

        setContent {
            // Wait for the answer so the first /ping isn't sent before the network is allowed.
            if (permissionAnswered.value) App()
        }
    }
}
