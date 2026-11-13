package com.surfacevolume.spike2.kmp

/**
 * Where the Java backend is, seen from each platform.
 * `expect` = declared here once, `actual` = one implementation per platform (androidMain, wasmJsMain).
 */
expect fun backendBaseUrl(): String
