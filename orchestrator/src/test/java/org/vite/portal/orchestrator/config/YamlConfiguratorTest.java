package org.vite.portal.orchestrator.config;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;

@SpringBootTest
public class YamlConfiguratorTest {
  @Autowired
  YamlConfigurator configurator;

  @Test
	public void getConfig() throws Exception {
    YamlConfig config = configurator.getConfig();
		assertNotNull(config);
    assertEquals("test-YAML", config.getName());
    assertEquals("testing", config.getEnvironment());
	}
}
