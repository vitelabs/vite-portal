package org.vite.portal.orchestrator.websocket;

import com.alibaba.fastjson.JSONObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;
import org.vite.portal.orchestrator.beans.NodeInfo;
import org.vite.portal.orchestrator.config.HttpSessionConfigurator;
import org.vite.portal.orchestrator.utils.ApplicationContextUtils;

import javax.servlet.http.HttpSession;
import javax.websocket.EndpointConfig;
import javax.websocket.OnClose;
import javax.websocket.OnError;
import javax.websocket.OnMessage;
import javax.websocket.OnOpen;
import javax.websocket.Session;
import javax.websocket.server.PathParam;
import javax.websocket.server.ServerEndpoint;

import java.io.IOException;
import java.nio.ByteBuffer;
import java.util.concurrent.CopyOnWriteArraySet;
import java.util.concurrent.Future;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicLong;

import static org.vite.portal.orchestrator.OrchestratorApplication.MSG_LOGGER;

@ServerEndpoint(value = "/ws/gvite/{sid}", configurator = HttpSessionConfigurator.class)
@Component
public class NodeWebSocketServer {
  private static final String OS_COMAND = "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"" + "dashboard_osInfo"
      + "\",\"params\":[\"request1.3.0\"]}";
  public static final String RUNTIME_COMMAND = "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"" + "dashboard_runtimeInfo"
      + "\",\"params\":[\"request1\"]}";
  private static final String PROCESS_COMMAND = "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\""
      + "dashboard_processInfo" + "\",\"params\":[\"request1\"]}";

  public static AtomicLong gloableHeight = new AtomicLong(0);

  public static AtomicLong snapshotLastestTime = new AtomicLong(0);

  private static Logger log = LoggerFactory.getLogger(NodeWebSocketServer.class);

  private static AtomicInteger onlineCount = new AtomicInteger(0);
  private static CopyOnWriteArraySet<NodeWebSocketServer> webSocketSet = new CopyOnWriteArraySet<NodeWebSocketServer>();

  private Session session;
  private String sid = ""; // The unique identifier of the node network, the format is
                           // 'vite_network_id@node_id'
  private String remoteIp;
  private String version;
  private int peerNum;
  private JSONObject snapshot;
  private NodeInfo nodeInfo;
  private WSMessageHandler wsMessageHandler;

  public static CopyOnWriteArraySet<NodeWebSocketServer> getWebSockets() {
    return webSocketSet;
  }

  public void ping() {
    try {
      if (session != null) {
        if (session.isOpen()) {
          session.getAsyncRemote().sendPing(ByteBuffer.wrap(new byte[] {}));
        }
      }
    } catch (Exception e) {
      log.error("send ping fail, sid:{}", sid, e);
    }
  }

  private String ipTrans(String ip) {
    if (ip == null) {
      return null;
    }

    String result = ip;
    String[] strs = ip.split("\\.");
    if (strs.length == 4) {
      result = strs[0] + "." + strs[1] + "." + "*" + "." + "*";
    }
    return result;
  }

  /**
   * The method called when the connection is established successfully
   */
  @OnOpen
  public void onOpen(Session session, @PathParam("sid") String sid, EndpointConfig config) {
    // if (ips == null) {
    // ips = WSMessageHandlerImpl.getDBByName("ip");
    // }
    wsMessageHandler = (WSMessageHandler) ApplicationContextUtils.getBean(WSMessageHandlerImpl.class);
    HttpSession httpSession = (HttpSession) config.getUserProperties().get(HttpSession.class.getName());
    if (httpSession.getAttribute("ClientIP") != null) {
      this.remoteIp = (String) httpSession.getAttribute("ClientIP");
      // log.info("{} getMaxTextMessageBufferSize : {}", remoteIp,
      // session.getMaxTextMessageBufferSize());
      int maxSize = 200 * 1024;
      // Maximum length of incoming binary messages that can be buffered
      session.setMaxBinaryMessageBufferSize(maxSize);
      // Maximum length of incoming text messages that can be buffered
      session.setMaxTextMessageBufferSize(maxSize);
      // log.info("{} SETMaxTextMessageBufferSize : {}", remoteIp,
      // session.getMaxTextMessageBufferSize());
    }

    this.session = session;
    addOnlineCount(); // Add 1 to the online count
    MSG_LOGGER.info("nodeJoin {} {} {}", sid, this.remoteIp, getOnlineCount());
    this.sid = sid;

    this.nodeInfo = new NodeInfo();
    this.nodeInfo.setSid(sid);
    this.nodeInfo.setIp(ipTrans(remoteIp));

    try {
      sendBasicMessage(OS_COMAND);
      sendBasicMessage(PROCESS_COMMAND);
      long t = System.currentTimeMillis();
      sendBasicMessage(NodeWebSocketServer.RUNTIME_COMMAND.replace("\"id\":1", "\"id\":" + t));
    } catch (Exception e) {
      // TODO: If you can't get the basic information, just turn off the websocket
      MSG_LOGGER.error("sendBase fail {} {}", this.sid, this.remoteIp, e);
    }

    // JSONObject ipJson = null;
    // if (ips.containsKey(this.remoteIp)) {
    // ipJson = ips.get(this.remoteIp);
    // } else {
    // try {
    // ipJson = HttpRequestUtils.httpPost("https://iplocation.com?ip=" +
    // this.remoteIp, null);
    // if (ipJson != null) {
    // ips.put(this.remoteIp, ipJson);
    // }
    // } catch (Exception e) {
    // log.error("get ip info fail", e);
    // }
    // }

    // if (ipJson != null) {
    // this.nodeInfo.setIpInfo(JSONObject.toJavaObject(ipJson, IPInfo.class));
    // }

    webSocketSet.add(this); // Add to set
  }

  /**
   * Method to be called on connection close
   */
  @OnClose
  public void onClose() {
    webSocketSet.remove(this); // Remove from set
    subOnlineCount(); // Online count minus 1
    MSG_LOGGER.info("nodeLeave {} {} {}", sid, this.remoteIp, getOnlineCount());
  }

  /**
   * Method to be called when a client message is received
   * 
   * @param message Message sent by client
   */
  @OnMessage
  public void onMessage(String message, Session session) {
    wsMessageHandler.messageHandler(message, this);
  }

  /**
   * @param session
   * @param error
   */
  @OnError
  public void onError(Session session, Throwable error) {
    MSG_LOGGER.error("wsError {} {}", sid, remoteIp, error);
    // TODO: Why delete it here? Will it close immediately after the error?
    // error.printStackTrace();
    // webSocketSet.remove(this);
    // subOnlineCount(); // Online count minus 1
    // log.info("error leave {}", remoteIp);
  }

  /**
   * Implement active server push
   * 
   * @return
   */
  public synchronized Future<Void> sendMessage(String message) throws IOException {
    return this.session.getAsyncRemote().sendText(message);
  }

  /**
   * Implement active server push
   */
  public synchronized void sendBasicMessage(String message) throws IOException {
    this.session.getBasicRemote().sendText(message);
  }

  /**
   * Bulk custom message
   */
  public static void sendInfo(String message, @PathParam("sid") String sid) throws IOException {
    for (NodeWebSocketServer item : webSocketSet) {
      try {
        // Here you can set only push to this sid, if it is null, push all
        if (sid == null || item.sid.equals(sid)) {
          item.sendMessage(message);
        }
      } catch (IOException e) {
        // TODO: Failed to send close the socket
        log.error("send fail {} {}", sid, item.remoteIp, e);
        continue;
      }
    }
  }

  public static int getOnlineCount() {
    return onlineCount.get();
  }

  public static void addOnlineCount() {
    NodeWebSocketServer.onlineCount.addAndGet(1);
  }

  public static void subOnlineCount() {
    NodeWebSocketServer.onlineCount.decrementAndGet();
  }

  public Session getSession() {
    return session;
  }

  public void setSession(Session session) {
    this.session = session;
  }

  public String getSid() {
    return sid;
  }

  public void setSid(String sid) {
    this.sid = sid;
  }

  public String getRemoteIp() {
    return remoteIp;
  }

  public void setRemoteIp(String remoteIp) {
    this.remoteIp = remoteIp;
  }

  public String getVersion() {
    return version;
  }

  public void setVersion(String version) {
    this.version = version;
  }

  public int getPeerNum() {
    return peerNum;
  }

  public void setPeerNum(int peerNum) {
    this.peerNum = peerNum;
  }

  public JSONObject getSnapshot() {
    return snapshot;
  }

  public void setSnapshot(JSONObject snapshot) {
    this.snapshot = snapshot;
  }

  public NodeInfo getNodeInfo() {
    return nodeInfo;
  }

  public void setNodeInfo(NodeInfo nodeInfo) {
    this.nodeInfo = nodeInfo;
  }
}
