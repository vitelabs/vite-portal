package org.vite.portal.orchestrator;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.ConfigurableApplicationContext;
import org.vite.portal.orchestrator.utils.ApplicationContextUtils;

@SpringBootApplication
public class OrchestratorApplication {
	public static final Logger MSG_LOGGER = LoggerFactory.getLogger("MSG");

	public static void main(String[] args) {
		ConfigurableApplicationContext ctx = SpringApplication.run(OrchestratorApplication.class, args);
		ApplicationContextUtils.setApplicationContext(ctx);
	}

}
