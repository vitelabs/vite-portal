package org.vite.portal.orchestrator.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

@Component
public class FullNodeKafkaProperties {

    @Value("${kafka.servers}")
    private String servers;

    @Value("${kafka.awardTopic}")
    private String awardTopic;

    @Value("${kafka.sslLocation}")
    private String sslLocation;

    public String getServers() {
        return servers;
    }

    public String getSslLocation() {
        return sslLocation;
    }

    public void setSslLocation(String sslLocation) {
        this.sslLocation = sslLocation;
    }

    public void setServers(String servers) {
        this.servers = servers;
    }

    public String getAwardTopic() {
        return awardTopic;
    }

    public void setAwardTopic(String awardTopic) {
        this.awardTopic = awardTopic;
    }
}
