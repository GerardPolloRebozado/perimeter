package com.surfacevolume.spike2.ping.workflow;

import java.time.Instant;
import java.util.UUID;
import org.springframework.data.cassandra.core.cql.PrimaryKeyType;
import org.springframework.data.cassandra.core.mapping.CassandraType;
import org.springframework.data.cassandra.core.mapping.Column;
import org.springframework.data.cassandra.core.mapping.PrimaryKeyColumn;
import org.springframework.data.cassandra.core.mapping.Table;

/**
 * Row of surface_volume.workflows_by_user: partitioned by user.
 * Column names are explicit: by default Spring maps "actionService" to "actionservice", not "action_service".
 */
@Table("workflows_by_user")
public record WorkflowByUser(
        @PrimaryKeyColumn(name = "user_id", type = PrimaryKeyType.PARTITIONED) UUID userId,
        @PrimaryKeyColumn(name = "workflow_id", type = PrimaryKeyType.CLUSTERED)
                @CassandraType(type = CassandraType.Name.TIMEUUID) UUID workflowId,
        @Column("name") String name,
        @Column("action_service") String actionService,
        @Column("action_type") String actionType,
        @Column("action_config") String actionConfig,
        @Column("reaction_service") String reactionService,
        @Column("reaction_type") String reactionType,
        @Column("reaction_config") String reactionConfig,
        @Column("active") boolean active,
        @Column("created_at") Instant createdAt) {
}
