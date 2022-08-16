package org.vite.portal.orchestrator.websocket;

public interface WSMessageHandler {
  void messageHandler(String msg, NodeWebSocketServer webSocketServer);
}
