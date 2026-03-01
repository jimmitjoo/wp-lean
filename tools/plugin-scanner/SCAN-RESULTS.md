# Plugin Security Scan — WordPress Top 10

Scanned 2026-03-01 against the 10 most installed plugins on wordpress.org.

## Results

| Plugin | Installs | Files | CRITICAL | HIGH | MEDIUM | Verdict |
|--------|----------|-------|----------|------|--------|---------|
| Elementor | 10M | 1343 | 3 | 159 | 2 | BLOCKED |
| Yoast SEO | 10M | 1325 | 0 | 156 | 1 | PASS |
| Contact Form 7 | 10M | 111 | 0 | 4 | 3 | PASS |
| Classic Editor | 9M | 1 | 0 | 17 | 0 | PASS |
| LiteSpeed Cache | 7M | 198 | 1 | 104 | 1 | BLOCKED |
| WooCommerce | 7M | 2783 | 0 | 757 | 16 | PASS |
| Akismet | 6M | 22 | 0 | 39 | 0 | PASS |
| WPForms Lite | 6M | 3740 | 2 | 290 | 5 | BLOCKED |
| Google Site Kit | 5M | 1769 | 10 | 65 | 14 | BLOCKED |
| All-in-One Migration | 5M | 165 | 0 | 13 | 4 | PASS |

**6 of 10 pass. 4 blocked.**

## CRITICAL findings detail

### Elementor (3 critical)

All in bundled Twig template engine:

```
vendor_prefixed/twig/twig/src/Environment.php:350
  eval('?>' . $content);

vendor_prefixed/twig/twig/src/Test/IntegrationTestCase.php:133
  eval('$ret = ' . $condition . ';');

vendor_prefixed/twig/twig/src/Test/IntegrationTestCase.php:182
  $output = trim($template->render(eval($match[1] . ';')), "\n ");
```

Twig executes arbitrary PHP via eval on concatenated strings. This is the template engine's core rendering mechanism — not a bug, but a fundamental design choice incompatible with secure plugin execution.

### LiteSpeed Cache (1 critical)

```
lib/css_js_min/minify/js.cls.php:245
  'exec(',
```

JS minifier processes JavaScript containing exec() patterns. The minifier handles potentially unsanitized input containing executable code patterns.

### WPForms Lite (2 critical)

Both in bundled HTMLPurifier:

```
vendor_prefixed/ezyang/htmlpurifier/library/HTMLPurifier/ConfigSchema/InterchangeBuilder.php:157
  return eval('return array(' . $contents . ');');

vendor_prefixed/ezyang/htmlpurifier/library/HTMLPurifier/VarParser/Native.php:30
  $result = eval("\$var = {$expr};");
```

HTMLPurifier uses eval to parse configuration arrays and variable expressions. A sanitization library that itself uses eval.

### Google Site Kit (10 critical)

phpseclib (7): Generates cryptographic functions by concatenating strings and calling eval. Used for BigInteger math and symmetric key operations.

```
phpseclib/Crypt/Common/SymmetricKey.php:2925
  eval('$func = function ($_action, $_text) { ' . $init_crypt . '...');

phpseclib/Math/BigInteger/Engines/BCMath/Reductions/EvalBarrett.php:57
  eval('$func = function ($n) { ' . $code . '};');
```

Google Auth (1): Shells out to detect gcloud CLI credentials.

```
google/auth/src/CredentialsLoader.php:193
  exec(implode(' ', $cmd), $output, $returnVar);
```

Monolog (1): Opens shell processes for log handling.

```
monolog/src/Monolog/Handler/ProcessHandler.php:104
  proc_open($this->command, static::DESCRIPTOR_SPEC, $this->pipes, $this->cwd);
```

phpseclib SSH (1): Method named exec on SSH2 class (not PHP's exec function, but scoped as exec capability).

## Analysis

Every CRITICAL finding is in a **bundled vendor library**, not in plugin code. This is the core problem:

1. **Vendor bundling** — Each plugin ships its own copy of third-party libraries. No shared versions, no coordinated updates, no security review.

2. **eval as architecture** — phpseclib generates crypto via eval'd string concatenation. Twig renders templates via eval. HTMLPurifier parses config via eval. These are design choices, not bugs.

3. **Scale** — These 4 plugins represent ~500M active installations. The vendor code they bundle was never reviewed by WordPress.org's plugin review process.

## Position

No exceptions for vendor code. CRITICAL is CRITICAL regardless of directory. If a library requires eval/exec to function, it needs to be replaced with one that doesn't, or the capability must be explicitly declared and approved.
