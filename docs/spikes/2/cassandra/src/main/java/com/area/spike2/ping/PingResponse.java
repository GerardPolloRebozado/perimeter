package com.area.spike2.ping;

/** Body of GET /ping. Serialized by Jackson as {"status":"..."}. */
public record PingResponse(String status) {

    public static PingResponse ok() {
        return new PingResponse("ok");
    }

    public static PingResponse error() {
        return new PingResponse("error");
    }
}
