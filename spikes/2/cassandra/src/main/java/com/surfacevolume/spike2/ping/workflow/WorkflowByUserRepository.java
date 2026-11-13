package com.surfacevolume.spike2.ping.workflow;

import java.util.List;
import java.util.UUID;
import org.springframework.data.cassandra.core.mapping.MapId;
import org.springframework.data.cassandra.repository.CassandraRepository;

public interface WorkflowByUserRepository extends CassandraRepository<WorkflowByUser, MapId> {

    /** Spring generates: SELECT * FROM workflows_by_user WHERE user_id = ? */
    List<WorkflowByUser> findByUserId(UUID userId);
}
