<?php
class WP_Plugin_Sandbox {

  private static $scanner_binary;
  private static $initialized = false;

  public static function init() {
    if (self::$initialized) {
      return;
    }
    self::$initialized = true;
    self::$scanner_binary = ABSPATH . WPINC . '/plugin-scanner';

    add_filter('validate_plugin_requirements', array(self::class, 'scan_plugin'), 10, 2);
  }

  public static function scan_plugin($result, $plugin) {
    if (is_wp_error($result)) {
      return $result;
    }

    $plugin_dir = WP_PLUGIN_DIR . '/' . dirname($plugin);
    if (!is_dir($plugin_dir) || dirname($plugin) === '.') {
      $plugin_dir = WP_PLUGIN_DIR . '/' . $plugin;
      if (!is_file($plugin_dir)) {
        return $result;
      }
      $plugin_dir = dirname($plugin_dir);
    }

    $scan_result = self::run_scanner($plugin_dir, $plugin);
    if ($scan_result === null) {
      return $result;
    }

    if ($scan_result['critical_count'] > 0) {
      return new WP_Error(
        'plugin_security_critical',
        sprintf(
          'Plugin blocked: %d critical security issues found. Run "plugin-scanner %s" for details.',
          $scan_result['critical_count'],
          $plugin_dir
        ),
        array('scan_result' => $scan_result)
      );
    }

    if (!empty($scan_result['violations'])) {
      $undeclared = array();
      foreach ($scan_result['violations'] as $v) {
        $undeclared[] = $v['capability'];
      }
      return new WP_Error(
        'plugin_capability_violation',
        sprintf(
          'Plugin uses undeclared capabilities: %s. Add "Capabilities: %s" to plugin header.',
          implode(', ', $undeclared),
          implode(', ', array_merge(
            self::get_declared_capabilities($plugin),
            $undeclared
          ))
        ),
        array('scan_result' => $scan_result)
      );
    }

    if ($scan_result['high_count'] > 0) {
      do_action('plugin_security_warning', $plugin, $scan_result);
    }

    return $result;
  }

  private static function run_scanner($plugin_dir, $plugin) {
    if (!is_executable(self::$scanner_binary)) {
      return null;
    }

    $args = array(
      '--json',
    );

    $declared_caps = self::get_declared_capabilities($plugin);
    if (!empty($declared_caps)) {
      $args[] = '--check-caps';
      $args[] = implode(',', $declared_caps);
    }

    $args[] = $plugin_dir;

    $cmd = escapeshellarg(self::$scanner_binary);
    foreach ($args as $arg) {
      $cmd .= ' ' . escapeshellarg($arg);
    }

    $output = array();
    $exit_code = 0;
    exec($cmd . ' 2>/dev/null', $output, $exit_code);

    $json = implode("\n", $output);
    $data = json_decode($json, true);

    if ($data === null) {
      return null;
    }

    return $data;
  }

  public static function get_declared_capabilities($plugin) {
    $plugin_file = WP_PLUGIN_DIR . '/' . $plugin;
    if (!is_file($plugin_file)) {
      return array();
    }

    $header = get_file_data($plugin_file, array('capabilities' => 'Capabilities'));
    if (empty($header['capabilities'])) {
      return array();
    }

    $caps = array_map('trim', explode(',', $header['capabilities']));
    return array_filter($caps);
  }

  public static function scan_all_active() {
    $active = get_option('active_plugins', array());
    $results = array();
    foreach ($active as $plugin) {
      $plugin_dir = WP_PLUGIN_DIR . '/' . dirname($plugin);
      if (is_dir($plugin_dir) && dirname($plugin) !== '.') {
        $scan = self::run_scanner($plugin_dir, $plugin);
        if ($scan !== null) {
          $results[$plugin] = $scan;
        }
      }
    }
    return $results;
  }
}
