package com.area.spike2.ping;

import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.content;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.header;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

import com.datastax.oss.driver.api.core.CqlSession;
import com.datastax.oss.driver.api.core.NoNodeAvailableException;
import com.datastax.oss.driver.api.core.cql.ResultSet;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.webmvc.test.autoconfigure.WebMvcTest;
import org.springframework.context.annotation.Import;
import org.springframework.http.MediaType;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;

/** No real Cassandra here: the session is a mock, so we can simulate "up" and "down". */
@WebMvcTest(PingController.class)
@Import(WebConfig.class)
class PingControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockitoBean
    private CqlSession session;

    @Test
    void pingIsOkWhenCassandraAnswers() throws Exception {
        when(session.execute(anyString())).thenReturn(mock(ResultSet.class));

        mockMvc.perform(get("/ping"))
                .andExpect(status().isOk())
                .andExpect(content().contentTypeCompatibleWith(MediaType.APPLICATION_JSON))
                .andExpect(content().string("{\"status\":\"ok\"}"));
    }

    @Test
    void pingIs503WhenCassandraIsDown() throws Exception {
        when(session.execute(anyString())).thenThrow(new NoNodeAvailableException());

        mockMvc.perform(get("/ping"))
                .andExpect(status().isServiceUnavailable())
                .andExpect(content().string("{\"status\":\"error\"}"));
    }

    @Test
    void pingAllowsWebClientOrigin() throws Exception {
        when(session.execute(anyString())).thenReturn(mock(ResultSet.class));

        mockMvc.perform(get("/ping").header("Origin", "http://localhost:8081"))
                .andExpect(status().isOk())
                .andExpect(header().string("Access-Control-Allow-Origin", "http://localhost:8081"));
    }
}
