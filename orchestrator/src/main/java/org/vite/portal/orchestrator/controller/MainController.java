package org.vite.portal.orchestrator.controller;

import javax.annotation.Resource;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import org.vite.portal.orchestrator.config.FullNodeKafkaProperties;
import org.vite.portal.orchestrator.config.YamlConfigurator;

@RestController
public class MainController {

  @Autowired
  private YamlConfigurator configurator;

  @Resource
  private FullNodeKafkaProperties kafkaProperties;

  @GetMapping("/")
  public String index() {
    return "Greetings from Spring Boot!";
  }

  @GetMapping("/config")
  public String config() {
    return configurator.getConfig().getName();
  }

  @GetMapping("/props")
  public FullNodeKafkaProperties props() {
    return kafkaProperties;
  }

}