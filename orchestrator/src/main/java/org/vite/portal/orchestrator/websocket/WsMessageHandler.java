package org.vite.portal.orchestrator.websocket;

public interface WsMessageHandler {
  void messageHandler(String msg, NodeWebSocketServer webSocketServer);
}
