<?php
class WP_Input {

  private static $get;
  private static $post;
  private static $request;
  private static $server;

  public static function init() {
    self::$get = $_GET;
    self::$post = $_POST;
    self::$request = $_REQUEST;
    self::$server = $_SERVER;
  }

  public static function get_string( $key, $default = '' ) {
    return self::sanitize_string( self::$get, $key, $default );
  }

  public static function get_int( $key, $default = 0 ) {
    return self::sanitize_int( self::$get, $key, $default );
  }

  public static function get_array( $key, $default = array() ) {
    return self::sanitize_array( self::$get, $key, $default );
  }

  public static function get_raw( $key, $default = '' ) {
    return isset( self::$get[ $key ] ) ? self::$get[ $key ] : $default;
  }

  public static function post_string( $key, $default = '' ) {
    return self::sanitize_string( self::$post, $key, $default );
  }

  public static function post_int( $key, $default = 0 ) {
    return self::sanitize_int( self::$post, $key, $default );
  }

  public static function post_array( $key, $default = array() ) {
    return self::sanitize_array( self::$post, $key, $default );
  }

  public static function post_raw( $key, $default = '' ) {
    return isset( self::$post[ $key ] ) ? self::$post[ $key ] : $default;
  }

  public static function request_string( $key, $default = '' ) {
    return self::sanitize_string( self::$request, $key, $default );
  }

  public static function request_int( $key, $default = 0 ) {
    return self::sanitize_int( self::$request, $key, $default );
  }

  public static function request_array( $key, $default = array() ) {
    return self::sanitize_array( self::$request, $key, $default );
  }

  public static function request_raw( $key, $default = '' ) {
    return isset( self::$request[ $key ] ) ? self::$request[ $key ] : $default;
  }

  public static function server( $key, $default = '' ) {
    $value = isset( self::$server[ $key ] ) ? self::$server[ $key ] : $default;
    return sanitize_text_field( wp_unslash( $value ) );
  }

  public static function has_get( $key ) {
    return isset( self::$get[ $key ] );
  }

  public static function has_post( $key ) {
    return isset( self::$post[ $key ] );
  }

  public static function has_request( $key ) {
    return isset( self::$request[ $key ] );
  }

  public static function all_post() {
    return array_map( 'sanitize_text_field', wp_unslash( self::$post ) );
  }

  public static function all_get() {
    return array_map( 'sanitize_text_field', wp_unslash( self::$get ) );
  }

  private static function sanitize_string( $source, $key, $default ) {
    if ( ! isset( $source[ $key ] ) ) {
      return $default;
    }
    return sanitize_text_field( wp_unslash( $source[ $key ] ) );
  }

  private static function sanitize_int( $source, $key, $default ) {
    if ( ! isset( $source[ $key ] ) ) {
      return $default;
    }
    return intval( $source[ $key ] );
  }

  private static function sanitize_array( $source, $key, $default ) {
    if ( ! isset( $source[ $key ] ) || ! is_array( $source[ $key ] ) ) {
      return $default;
    }
    return array_map( 'sanitize_text_field', wp_unslash( $source[ $key ] ) );
  }
}
