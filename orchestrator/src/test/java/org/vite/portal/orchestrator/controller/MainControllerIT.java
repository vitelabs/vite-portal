package org.vite.portal.orchestrator.controller;

import org.junit.jupiter.api.Test;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.http.ResponseEntity;
import org.vite.portal.orchestrator.beans.ConfigView;

import static org.assertj.core.api.Assertions.assertThat;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
public class MainControllerIT {

	@Autowired
	private TestRestTemplate template;

	@Test
	public void getIndex() throws Exception {
		ResponseEntity<String> response = template.getForEntity("/", String.class);
		assertThat(response.getBody()).isEqualTo("vite-portal-orchestrator");
	}

	@Test
	public void getConfig() throws Exception {
		ResponseEntity<ConfigView> response = template.getForEntity("/config", ConfigView.class);
		assertThat(response.getBody().getName()).isEqualTo("test-YAML");
		assertThat(response.getBody().getEnvironment()).isEqualTo("testing");
	}
}