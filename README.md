# wp-lean

WordPress trimmed for 2030. Based on WordPress 6.7.1.

## What changed

```
                               Original    wp-lean
──────────────────────────────────────────────────
PHP files                         1981       1602
Total lines                     655860     583950
require in wp-settings.php         329        130
```

376 files removed, 72k lines cut, 11% reduction.

### Removed

- **XML-RPC** — xmlrpc.php, IXR library, WP_XMLRPC_Server. REST API exists.
- **Trackbacks & pingbacks** — wp-trackback.php. Dead protocol.
- **Link manager** — wp-links-opml.php, bookmark functions. Removed from UI since 3.5.
- **wp-mail.php** — Post-by-email. Nobody uses it.
- **Classic themes** — twentyten through twentytwentyone. Block themes only.
- **Customizer** — All class-wp-customize-* files. Site Editor replaced it.
- **Deprecated loaders** — deprecated.php, pluggable-deprecated.php, ms-deprecated.php.
- **Snoopy** — HTTP client replaced by WP_Http over a decade ago.

### Added

**Autoloader** — `wp-includes/class-wp-autoloader.php`

Classmap-based autoloader replacing 200 explicit require statements. 292 classes loaded on demand via `spl_autoload_register`. Boot sequence goes from 329 to 130 requires.

**Plugin security sandbox** — `wp-includes/class-wp-plugin-sandbox.php`

Static analysis gate at plugin activation. Hooks into `validate_plugin_requirements`:

- Runs the Go scanner against the plugin directory
- **CRITICAL findings** → activation blocked
- **Undeclared capabilities** → activation blocked
- **HIGH findings** → warning logged, activation allowed

Plugins declare capabilities in their header:

```php
/**
 * Plugin Name: My Plugin
 * Capabilities: database-raw, http-outbound, filesystem-read
 */
```

Available capabilities: `exec`, `eval`, `http-outbound`, `filesystem-read`, `filesystem-write`, `database-raw`, `unserialize`, `obfuscation`.

## Plugin scanner

Go-based static analysis tool in `tools/plugin-scanner/`. 30 security patterns across three severity levels.

```
cd tools/plugin-scanner
go build -o plugin-scanner .

# Scan a plugin
./plugin-scanner /path/to/plugin

# JSON output
./plugin-scanner --json /path/to/plugin

# Verify declared capabilities
./plugin-scanner --json --check-caps "database-raw,http-outbound" /path/to/plugin
```

Exit codes: 0 = pass, 1 = high findings, 2 = critical findings, 3 = capability violation.

### Top 10 results

Scanned against the 10 most installed WordPress plugins:

| Plugin | Installs | CRITICAL | Verdict |
|--------|----------|----------|---------|
| Elementor | 10M | 3 | BLOCKED — Twig eval |
| Yoast SEO | 10M | 0 | PASS |
| Contact Form 7 | 10M | 0 | PASS |
| Classic Editor | 9M | 0 | PASS |
| LiteSpeed Cache | 7M | 1 | BLOCKED — JS minifier |
| WooCommerce | 7M | 0 | PASS |
| Akismet | 6M | 0 | PASS |
| WPForms Lite | 6M | 2 | BLOCKED — HTMLPurifier eval |
| Google Site Kit | 5M | 10 | BLOCKED — phpseclib eval |
| All-in-One Migration | 5M | 0 | PASS |

Every CRITICAL finding is in bundled vendor libraries, not plugin code. See [full results](tools/plugin-scanner/SCAN-RESULTS.md).

## What's preserved

- **Hook system** — `add_action`, `add_filter`, `do_action`, `apply_filters`. The contract plugins depend on.
- **REST API** — Full endpoint system.
- **Block editor** — Gutenberg, block themes, theme.json.
- **Plugin API** — Activation, deactivation, updates. All hooks intact.
- **Multisite** — Full support.
- **All of wp-admin** — Dashboard, post editor, media library, user management.

## Requirements

- PHP 8.0+
- Go 1.22+ (to build the scanner)
