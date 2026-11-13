package com.surfacevolume.spike2.ping.workflow;

import java.util.List;
import java.util.UUID;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.server.ResponseStatusException;

@RestController
public class WorkflowController {

    private final WorkflowService service;

    public WorkflowController(WorkflowService service) {
        this.service = service;
    }

    @PostMapping("/workflows")
    @ResponseStatus(HttpStatus.CREATED)
    public WorkflowDto create(@RequestBody CreateWorkflowRequest req) {
        if (req.userId() == null || isBlank(req.actionService()) || isBlank(req.reactionService())) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST,
                    "userId, actionService and reactionService are required");
        }
        return service.create(req);
    }

    @GetMapping("/users/{userId}/workflows")
    public List<WorkflowDto> byUser(@PathVariable UUID userId) {
        return service.findByUser(userId);
    }

    @GetMapping("/services/{service}/workflows")
    public List<WorkflowDto> byActionService(@PathVariable String service) {
        return this.service.findByActionService(service);
    }

    private static boolean isBlank(String s) {
        return s == null || s.isBlank();
    }
}
