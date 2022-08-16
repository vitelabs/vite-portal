package org.vite.portal.orchestrator.utils;

import org.springframework.context.ApplicationContext;

public class ApplicationContextUtils {

  private static ApplicationContext applicationContext;

  public static void setApplicationContext(ApplicationContext applicationContext) {
    ApplicationContextUtils.applicationContext = applicationContext;
  }

  public static Object getBean(Class<?> tClass) {
    return applicationContext.getBean(tClass);
  }
}
