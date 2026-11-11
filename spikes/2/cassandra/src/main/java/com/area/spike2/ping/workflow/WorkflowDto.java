package com.area.spike2.ping.workflow;

import java.time.Instant;
import java.util.Map;
import java.util.UUID;

/** What the API returns for a workflow, whichever table it was read from. */
public record WorkflowDto(
        UUID workflowId,
        UUID userId,
        String name,
        String actionService,
        String actionType,
        Map<String, Object> actionConfig,
        String reactionService,
        String reactionType,
        Map<String, Object> reactionConfig,
        boolean active,
        Instant createdAt) {
}
