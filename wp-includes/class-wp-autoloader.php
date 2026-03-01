<?php
class WP_Autoloader {
  private static $classmap = [];
  private static $registered = false;

  public static function register() {
    if (self::$registered) {
      return;
    }
    self::$registered = true;
    self::$classmap = require __DIR__ . '/classmap.php';
    spl_autoload_register([self::class, 'load']);
  }

  public static function load($class) {
    if (isset(self::$classmap[$class])) {
      require self::$classmap[$class];
      return true;
    }
    return false;
  }
}
