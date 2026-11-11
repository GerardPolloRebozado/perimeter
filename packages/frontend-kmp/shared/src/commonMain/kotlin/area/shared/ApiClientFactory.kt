package area.shared

import area.api.client.DefaultApi
import area.api.client.ApiClient

// Factory for creating typed API client from generated OpenAPI client
object ApiClientFactory {
    
    private const val BACKEND_URL = "http://localhost:8080"
    
    fun create(): DefaultApi {
        val client = ApiClient(BACKEND_URL)
        return DefaultApi(client)
    }
    
    // Example type-safe usage:
    // val api = ApiClientFactory.create()
    // val workflows = api.getWorkflows()  // Returns typed List<WorkflowDto>
}
