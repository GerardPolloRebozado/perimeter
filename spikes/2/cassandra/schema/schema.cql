-- Spike 2 schema. Safe to run several times (IF NOT EXISTS everywhere).
CREATE KEYSPACE IF NOT EXISTS area
  WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

-- Cassandra is modelled query-first: one table per query, data duplicated between them.

-- Query 1: "all workflows of a user" (the user's dashboard).
CREATE TABLE IF NOT EXISTS area.workflows_by_user (
  user_id          uuid,
  workflow_id      timeuuid,
  name             text,
  action_service   text,
  action_type      text,
  action_config    text,      -- JSON, stored as text (Cassandra has no document type)
  reaction_service text,
  reaction_type    text,
  reaction_config  text,      -- JSON
  active           boolean,
  created_at       timestamp,
  PRIMARY KEY ((user_id), workflow_id)
) WITH CLUSTERING ORDER BY (workflow_id DESC);

-- Query 2: "all workflows triggered by service X" (what AREA's trigger poller needs).
-- Same columns, different partition key.
CREATE TABLE IF NOT EXISTS area.workflows_by_action_service (
  action_service   text,
  workflow_id      timeuuid,
  user_id          uuid,
  name             text,
  action_type      text,
  action_config    text,
  reaction_service text,
  reaction_type    text,
  reaction_config  text,
  active           boolean,
  created_at       timestamp,
  PRIMARY KEY ((action_service), workflow_id)
) WITH CLUSTERING ORDER BY (workflow_id DESC);
