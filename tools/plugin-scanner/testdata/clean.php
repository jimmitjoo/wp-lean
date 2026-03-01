<?php
function hello() {
  echo "Hello, World!";
}

$id = intval($_GET['id']);
$name = sanitize_text_field($_POST['name']);
$results = $wpdb->get_results($wpdb->prepare("SELECT * FROM wp_posts WHERE id = %d", $id));
add_action('init', 'hello');
