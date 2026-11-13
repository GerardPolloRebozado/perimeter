package org.surfacevolume.surface

interface Platform {
    val name: String
}

expect fun getPlatform(): Platform