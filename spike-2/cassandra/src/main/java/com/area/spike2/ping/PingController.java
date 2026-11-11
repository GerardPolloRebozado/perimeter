package com.area.spike2.ping;

import com.datastax.oss.driver.api.core.CqlSession;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

/** GET /ping answers "ok" only if Cassandra actually answers a query. */
@RestController
public class PingController {

    private static final Logger log = LoggerFactory.getLogger(PingController.class);

    // Cheapest possible query: reads a built-in system table, needs no schema of ours.
    private static final String HEALTH_QUERY = "SELECT release_version FROM system.local";

    private final CqlSession session;

    public PingController(CqlSession session) {
        this.session = session;
    }

    @GetMapping(path = "/ping", produces = MediaType.APPLICATION_JSON_VALUE)
    public ResponseEntity<PingResponse> ping() {
        try {
            session.execute(HEALTH_QUERY);
            return ResponseEntity.ok(PingResponse.ok());
        } catch (RuntimeException e) {
            log.warn("Cassandra health query failed: {}", e.toString());
            return ResponseEntity.status(HttpStatus.SERVICE_UNAVAILABLE).body(PingResponse.error());
        }
    }
}
