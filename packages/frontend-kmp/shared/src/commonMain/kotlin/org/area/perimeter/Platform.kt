package org.area.perimeter

interface Platform {
    val name: String
}

expect fun getPlatform(): Platform