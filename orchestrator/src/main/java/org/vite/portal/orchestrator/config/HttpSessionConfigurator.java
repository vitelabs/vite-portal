package org.vite.portal.orchestrator.config;

import org.apache.catalina.session.StandardSessionFacade;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Profile;
import org.springframework.web.socket.server.standard.ServerEndpointExporter;

import javax.servlet.http.HttpSession;
import javax.websocket.HandshakeResponse;
import javax.websocket.server.HandshakeRequest;
import javax.websocket.server.ServerEndpointConfig;

@Configuration
public class HttpSessionConfigurator extends ServerEndpointConfig.Configurator {
    /*
     * Modifying the handshake means to modify the provided content before the
     * handshake protocol is established
     */
    @Override
    public void modifyHandshake(ServerEndpointConfig sec, HandshakeRequest request, HandshakeResponse response) {
        // If there is no listener, then the HttpSession obtained here is null
        StandardSessionFacade ssf = (StandardSessionFacade) request.getHttpSession();
        if (ssf != null) {
            HttpSession session = (HttpSession) request.getHttpSession();
            sec.getUserProperties().put(HttpSession.class.getName(), session);
        }
        super.modifyHandshake(sec, request, response);
    }

    /**
     * Inject ServerEndpointExporter,
     * This bean will automatically register the Websocket endpoint declared with
     * the @ServerEndpoint annotation
     * 
     * When the active profile is different e.g. test, serverEndpointExporter will
     * be ignored. This is to prevent "javax.websocket.server.ServerContainer not
     * available" error when running tests.
     */
    @Profile({ "dev", "prod" })
    @Bean
    public ServerEndpointExporter serverEndpointExporter() {
        return new ServerEndpointExporter();
    }
}
