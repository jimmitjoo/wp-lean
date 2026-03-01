<?php
class WP_Vendor {

  private static $packages = array();
  private static $loaded = array();
  private static $manifest;
  private static $initialized = false;

  public static function init() {
    if ( self::$initialized ) {
      return;
    }
    self::$initialized = true;

    $manifest_file = ABSPATH . WPINC . '/vendor.json';
    if ( file_exists( $manifest_file ) ) {
      self::$manifest = json_decode( file_get_contents( $manifest_file ), true );
      if ( isset( self::$manifest['packages'] ) ) {
        foreach ( self::$manifest['packages'] as $name => $info ) {
          self::register( $name, $info['version'], ABSPATH . WPINC . '/' . $info['path'] );
        }
      }
    }
  }

  public static function register( $name, $version, $path ) {
    self::$packages[ $name ] = array(
      'version' => $version,
      'path'    => $path,
    );
  }

  public static function require_package( $name, $min_version = null ) {
    if ( isset( self::$loaded[ $name ] ) ) {
      if ( $min_version && version_compare( self::$loaded[ $name ], $min_version, '<' ) ) {
        return new WP_Error(
          'vendor_version_conflict',
          sprintf( 'Package %s loaded at %s, but %s required.', $name, self::$loaded[ $name ], $min_version )
        );
      }
      return true;
    }

    if ( ! isset( self::$packages[ $name ] ) ) {
      return new WP_Error(
        'vendor_not_found',
        sprintf( 'Package %s is not registered. Available: %s', $name, implode( ', ', array_keys( self::$packages ) ) )
      );
    }

    $pkg = self::$packages[ $name ];

    if ( $min_version && version_compare( $pkg['version'], $min_version, '<' ) ) {
      return new WP_Error(
        'vendor_version_mismatch',
        sprintf( 'Package %s is %s, but %s required.', $name, $pkg['version'], $min_version )
      );
    }

    if ( ! is_dir( $pkg['path'] ) ) {
      return new WP_Error(
        'vendor_path_missing',
        sprintf( 'Package %s path does not exist: %s', $name, $pkg['path'] )
      );
    }

    self::$loaded[ $name ] = $pkg['version'];

    do_action( 'wp_vendor_loaded', $name, $pkg['version'] );

    return true;
  }

  public static function is_loaded( $name ) {
    return isset( self::$loaded[ $name ] );
  }

  public static function get_version( $name ) {
    if ( isset( self::$loaded[ $name ] ) ) {
      return self::$loaded[ $name ];
    }
    if ( isset( self::$packages[ $name ] ) ) {
      return self::$packages[ $name ]['version'];
    }
    return false;
  }

  public static function get_path( $name ) {
    if ( isset( self::$packages[ $name ] ) ) {
      return self::$packages[ $name ]['path'];
    }
    return false;
  }

  public static function get_registered() {
    return self::$packages;
  }

  public static function get_loaded() {
    return self::$loaded;
  }

  public static function validate_plugin_vendors( $plugin ) {
    $plugin_file = WP_PLUGIN_DIR . '/' . $plugin;
    if ( ! is_file( $plugin_file ) ) {
      return true;
    }

    $header = get_file_data( $plugin_file, array( 'vendors' => 'Requires Vendors' ) );
    if ( empty( $header['vendors'] ) ) {
      return true;
    }

    $requirements = array_map( 'trim', explode( ',', $header['vendors'] ) );
    $missing = array();

    foreach ( $requirements as $req ) {
      $parts = preg_split( '/\s+/', $req, 2 );
      $name = $parts[0];
      $version = isset( $parts[1] ) ? trim( $parts[1], '><=^ ' ) : null;

      if ( ! isset( self::$packages[ $name ] ) ) {
        $missing[] = $name;
        continue;
      }

      if ( $version && version_compare( self::$packages[ $name ]['version'], $version, '<' ) ) {
        $missing[] = sprintf( '%s (need %s, have %s)', $name, $version, self::$packages[ $name ]['version'] );
      }
    }

    if ( ! empty( $missing ) ) {
      return new WP_Error(
        'plugin_vendor_missing',
        sprintf( 'Plugin requires vendor packages not available: %s', implode( ', ', $missing ) )
      );
    }

    return true;
  }
}
