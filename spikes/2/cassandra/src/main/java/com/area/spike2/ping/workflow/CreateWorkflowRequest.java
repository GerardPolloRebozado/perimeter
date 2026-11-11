package com.area.spike2.ping.workflow;

import java.util.Map;
import java.util.UUID;

/** Body of POST /workflows. The configs are free-form JSON objects. */
public record CreateWorkflowRequest(
        UUID userId,
        String name,
        String actionService,
        String actionType,
        Map<String, Object> actionConfig,
        String reactionService,
        String reactionType,
        Map<String, Object> reactionConfig) {
}
