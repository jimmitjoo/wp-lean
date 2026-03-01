<?php
class WP_Lazy_Loader {

  private static $deferred = array();
  private static $loaded = array();
  private static $hooks = array();

  public static function register( $file, $load_on = null ) {
    self::$deferred[ $file ] = true;
    if ( $load_on ) {
      if ( ! isset( self::$hooks[ $load_on ] ) ) {
        self::$hooks[ $load_on ] = array();
        add_action( $load_on, array( self::class, 'fire_hook' ), 0 );
      }
      self::$hooks[ $load_on ][] = $file;
    }
  }

  public static function load( $file ) {
    if ( isset( self::$loaded[ $file ] ) ) {
      return;
    }
    self::$loaded[ $file ] = true;
    unset( self::$deferred[ $file ] );
    require ABSPATH . WPINC . '/' . $file;
  }

  public static function fire_hook() {
    $hook = current_action();
    if ( ! isset( self::$hooks[ $hook ] ) ) {
      return;
    }
    foreach ( self::$hooks[ $hook ] as $file ) {
      self::load( $file );
    }
    unset( self::$hooks[ $hook ] );
  }

  public static function load_all_deferred() {
    foreach ( array_keys( self::$deferred ) as $file ) {
      self::load( $file );
    }
  }

  public static function is_loaded( $file ) {
    return isset( self::$loaded[ $file ] );
  }

  public static function get_deferred() {
    return array_keys( self::$deferred );
  }

  public static function get_loaded() {
    return array_keys( self::$loaded );
  }
}
