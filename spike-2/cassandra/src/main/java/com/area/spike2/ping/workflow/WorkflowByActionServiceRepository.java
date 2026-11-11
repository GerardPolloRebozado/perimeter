package com.area.spike2.ping.workflow;

import java.util.List;
import org.springframework.data.cassandra.core.mapping.MapId;
import org.springframework.data.cassandra.repository.CassandraRepository;

public interface WorkflowByActionServiceRepository extends CassandraRepository<WorkflowByActionService, MapId> {

    /** Spring generates: SELECT * FROM workflows_by_action_service WHERE action_service = ? */
    List<WorkflowByActionService> findByActionService(String actionService);
}
