package org.vite.portal.orchestrator.config;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

@Component
public class YamlConfigurator {
  @Autowired
  private YamlConfig config;

  public YamlConfig getConfig() {
    return config;
  }
}
