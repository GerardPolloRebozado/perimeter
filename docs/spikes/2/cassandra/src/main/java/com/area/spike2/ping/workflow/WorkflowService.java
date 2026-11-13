package com.surfacevolume.spike2.ping.workflow;

import com.datastax.oss.driver.api.core.uuid.Uuids;
import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import org.springframework.data.cassandra.core.CassandraTemplate;
import org.springframework.stereotype.Service;
import tools.jackson.core.type.TypeReference;
import tools.jackson.databind.ObjectMapper;

@Service
public class WorkflowService {

    private static final TypeReference<Map<String, Object>> JSON_OBJECT = new TypeReference<>() {};

    private final CassandraTemplate template;
    private final WorkflowByUserRepository byUser;
    private final WorkflowByActionServiceRepository byActionService;
    private final ObjectMapper json;

    public WorkflowService(CassandraTemplate template,
                           WorkflowByUserRepository byUser,
                           WorkflowByActionServiceRepository byActionService,
                           ObjectMapper json) {
        this.template = template;
        this.byUser = byUser;
        this.byActionService = byActionService;
        this.json = json;
    }

    /**
     * Writes the workflow to BOTH tables in one logged batch, so the two copies stay in sync.
     * This duplication is the price of Cassandra's one-table-per-query modelling.
     */
    public WorkflowDto create(CreateWorkflowRequest req) {
        UUID workflowId = Uuids.timeBased();
        // Cassandra timestamps have millisecond precision: truncate so what we return matches what we read back.
        Instant now = Instant.now().truncatedTo(ChronoUnit.MILLIS);
        String actionConfig = json.writeValueAsString(orEmpty(req.actionConfig()));
        String reactionConfig = json.writeValueAsString(orEmpty(req.reactionConfig()));

        var userRow = new WorkflowByUser(req.userId(), workflowId, req.name(),
                req.actionService(), req.actionType(), actionConfig,
                req.reactionService(), req.reactionType(), reactionConfig, true, now);
        var serviceRow = new WorkflowByActionService(req.actionService(), workflowId, req.userId(), req.name(),
                req.actionType(), actionConfig,
                req.reactionService(), req.reactionType(), reactionConfig, true, now);

        template.batchOps().insert(userRow).insert(serviceRow).execute();
        return toDto(userRow);
    }

    public List<WorkflowDto> findByUser(UUID userId) {
        return byUser.findByUserId(userId).stream().map(this::toDto).toList();
    }

    public List<WorkflowDto> findByActionService(String actionService) {
        return byActionService.findByActionService(actionService).stream().map(this::toDto).toList();
    }

    private WorkflowDto toDto(WorkflowByUser r) {
        return new WorkflowDto(r.workflowId(), r.userId(), r.name(),
                r.actionService(), r.actionType(), parse(r.actionConfig()),
                r.reactionService(), r.reactionType(), parse(r.reactionConfig()),
                r.active(), r.createdAt());
    }

    private WorkflowDto toDto(WorkflowByActionService r) {
        return new WorkflowDto(r.workflowId(), r.userId(), r.name(),
                r.actionService(), r.actionType(), parse(r.actionConfig()),
                r.reactionService(), r.reactionType(), parse(r.reactionConfig()),
                r.active(), r.createdAt());
    }

    private Map<String, Object> parse(String jsonText) {
        return jsonText == null ? Map.of() : json.readValue(jsonText, JSON_OBJECT);
    }

    private static Map<String, Object> orEmpty(Map<String, Object> map) {
        return map == null ? Map.of() : map;
    }
}
