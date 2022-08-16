package org.vite.portal.orchestrator.controller;

import javax.annotation.Resource;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import org.vite.portal.orchestrator.beans.ConfigView;
import org.vite.portal.orchestrator.config.FullNodeKafkaProperties;
import org.vite.portal.orchestrator.config.YamlConfig;
import org.vite.portal.orchestrator.config.YamlConfigurator;

@RestController
public class MainController {

  @Autowired
  private YamlConfigurator configurator;

  @Resource
  private FullNodeKafkaProperties kafkaProperties;

  @GetMapping("/")
  public String index() {
    return "vite-portal-orchestrator";
  }

  @GetMapping("/config")
  public ConfigView config() {
    YamlConfig config = configurator.getConfig();
    ConfigView view = new ConfigView();
    view.setName(config.getName());
    view.setEnvironment(config.getEnvironment());
    return view;
  }

}