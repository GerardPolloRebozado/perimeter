package area.shared

import area.api.client.DefaultApi
import area.api.client.models.WorkflowDto  // Generated model from OpenAPI

data class WorkflowViewModel(
    val id: String,
    val name: String,
    val actionService: String,
    val reactionService: String,
    val config: Map<String, Any>
)

object WorkflowService {
    
    private val api: DefaultApi = ApiClientFactory.create()
    
    suspend fun getWorkflows(): List<WorkflowViewModel> {
        // Call generated client - type-safe
        val dtos: List<WorkflowDto> = api.getWorkflows()
        
        // Map generated DTO to view model
        return dtos.map { dto ->
            WorkflowViewModel(
                id = dto.id?.toString() ?: "",
                name = dto.name ?: "",
                actionService = dto.actionService ?: "",
                reactionService = dto.reactionService ?: "",
                config = dto.actionConfig?.entries?.associate { it.key to it.value } ?: emptyMap()
            )
        }
    }
    
    suspend fun createWorkflow(name: String, actionService: String, actionConfig: Map<String, Any>) {
        // Type-safe creation through generated client
        api.createWorkflow(
            // Generated request types
        )
    }
}
