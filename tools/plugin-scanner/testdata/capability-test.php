<?php
/**
 * Plugin Name: Capability Test Plugin
 * Capabilities: database-raw, http-outbound
 */

$results = $wpdb->query("SELECT * FROM wp_posts WHERE id = $id");
$response = file_get_contents('https://api.example.com/data');
$ch = curl_init();
curl_exec($ch);
