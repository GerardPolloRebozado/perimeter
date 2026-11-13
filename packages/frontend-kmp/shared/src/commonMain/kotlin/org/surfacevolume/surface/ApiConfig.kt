package org.surfacevolume.surface

// Backend URL configuration - platform-specific implementation
expect object ApiConfig {
    val baseUrl: String
}
