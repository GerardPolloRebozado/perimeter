package com.area.spike2.ping;

import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class PingController {

    @GetMapping(path = "/ping", produces = MediaType.APPLICATION_JSON_VALUE)
    public PingResponse ping() {
        return PingResponse.ok();
    }
}
