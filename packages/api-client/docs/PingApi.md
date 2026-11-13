# PingApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
| ------------- | ------------- | ------------- |
| [**pingControllerPing**](PingApi.md#pingControllerPing) | **GET** /ping | Health check endpoint |


<a id="pingControllerPing"></a>
# **pingControllerPing**
> pingControllerPing()

Health check endpoint

### Example
```kotlin
// Import classes:
//import org.openapitools.client.infrastructure.*
//import surfacevolume.api.models.*

val apiInstance = PingApi()
try {
    apiInstance.pingControllerPing()
} catch (e: ClientException) {
    println("4xx response calling PingApi#pingControllerPing")
    e.printStackTrace()
} catch (e: ServerException) {
    println("5xx response calling PingApi#pingControllerPing")
    e.printStackTrace()
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

