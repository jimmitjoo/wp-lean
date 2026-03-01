package scanner

import "regexp"

type Severity int

const (
  Critical Severity = iota
  High
  Medium
)

func (s Severity) String() string {
  switch s {
  case Critical:
    return "CRITICAL"
  case High:
    return "HIGH"
  case Medium:
    return "MEDIUM"
  }
  return "UNKNOWN"
}

type Capability string

const (
  CapExec           Capability = "exec"
  CapEval           Capability = "eval"
  CapHTTPOutbound   Capability = "http-outbound"
  CapFilesystemRead Capability = "filesystem-read"
  CapFilesystemWrite Capability = "filesystem-write"
  CapDatabaseRaw    Capability = "database-raw"
  CapUnserialize    Capability = "unserialize"
  CapObfuscation    Capability = "obfuscation"
  CapNone           Capability = ""
)

type Pattern struct {
  Name       string
  Regex      *regexp.Regexp
  Severity   Severity
  Capability Capability
}

var Patterns = []Pattern{
  {Name: "eval()", Regex: regexp.MustCompile(`\beval\s*\(`), Severity: Critical, Capability: CapEval},
  {Name: "exec()", Regex: regexp.MustCompile(`\bexec\s*\(`), Severity: Critical, Capability: CapExec},
  {Name: "system()", Regex: regexp.MustCompile(`\bsystem\s*\(`), Severity: Critical, Capability: CapExec},
  {Name: "passthru()", Regex: regexp.MustCompile(`\bpassthru\s*\(`), Severity: Critical, Capability: CapExec},
  {Name: "shell_exec()", Regex: regexp.MustCompile(`\bshell_exec\s*\(`), Severity: Critical, Capability: CapExec},
  {Name: "popen()", Regex: regexp.MustCompile(`\bpopen\s*\(`), Severity: Critical, Capability: CapExec},
  {Name: "proc_open()", Regex: regexp.MustCompile(`\bproc_open\s*\(`), Severity: Critical, Capability: CapExec},
  {Name: "base64_decode+exec", Regex: regexp.MustCompile(`base64_decode\s*\(.*\b(eval|exec|include|require|system)\b`), Severity: Critical, Capability: CapObfuscation},
  {Name: "preg_replace /e modifier", Regex: regexp.MustCompile(`preg_replace\s*\(\s*['"][^'"]*\/e['"imsxuU]*\s*,`), Severity: Critical, Capability: CapEval},
  {Name: "create_function()", Regex: regexp.MustCompile(`\bcreate_function\s*\(`), Severity: Critical, Capability: CapEval},
  {Name: "assert() with string eval", Regex: regexp.MustCompile(`\bassert\s*\(\s*['".]`), Severity: Critical, Capability: CapEval},
  {Name: "assert() with user input", Regex: regexp.MustCompile(`\bassert\s*\(\s*\$_(GET|POST|REQUEST|COOKIE)`), Severity: Critical, Capability: CapEval},
  {Name: "file_put_contents to .php", Regex: regexp.MustCompile(`file_put_contents\s*\([^)]*\.php`), Severity: Critical, Capability: CapFilesystemWrite},
  {Name: "include/require with user input", Regex: regexp.MustCompile(`(include|require)(_once)?\s*\(?\s*\$_(GET|POST|REQUEST|COOKIE)`), Severity: Critical, Capability: CapNone},

  {Name: "variable variables ($$)", Regex: regexp.MustCompile(`\$\$[a-zA-Z_]`), Severity: High, Capability: CapNone},
  {Name: "extract()", Regex: regexp.MustCompile(`\bextract\s*\(`), Severity: High, Capability: CapNone},
  {Name: "unserialize()", Regex: regexp.MustCompile(`\bunserialize\s*\(`), Severity: High, Capability: CapUnserialize},
  {Name: "file_get_contents(http)", Regex: regexp.MustCompile(`file_get_contents\s*\(\s*['"]https?://`), Severity: High, Capability: CapHTTPOutbound},
  {Name: "file_get_contents(variable)", Regex: regexp.MustCompile(`file_get_contents\s*\(\s*\$`), Severity: High, Capability: CapFilesystemRead},
  {Name: "curl_exec()", Regex: regexp.MustCompile(`\bcurl_exec\s*\(`), Severity: High, Capability: CapHTTPOutbound},
  {Name: "raw $_GET access", Regex: regexp.MustCompile(`\$_GET\s*\[`), Severity: High, Capability: CapNone},
  {Name: "raw $_POST access", Regex: regexp.MustCompile(`\$_POST\s*\[`), Severity: High, Capability: CapNone},
  {Name: "raw $_REQUEST access", Regex: regexp.MustCompile(`\$_REQUEST\s*\[`), Severity: High, Capability: CapNone},
  {Name: "wpdb->query without prepare", Regex: regexp.MustCompile(`\$wpdb\s*->\s*query\s*\(\s*["'].*\$`), Severity: High, Capability: CapDatabaseRaw},

  {Name: "error_reporting(0)", Regex: regexp.MustCompile(`error_reporting\s*\(\s*0\s*\)`), Severity: Medium, Capability: CapNone},
  {Name: "chmod()", Regex: regexp.MustCompile(`\bchmod\s*\(`), Severity: Medium, Capability: CapFilesystemWrite},
  {Name: "ini_set()", Regex: regexp.MustCompile(`\bini_set\s*\(`), Severity: Medium, Capability: CapNone},
  {Name: "obfuscated hex/chr", Regex: regexp.MustCompile(`(chr\s*\(\s*\d+\s*\)\s*\.?\s*){4,}`), Severity: Medium, Capability: CapObfuscation},
  {Name: "hex2bin()", Regex: regexp.MustCompile(`\bhex2bin\s*\(`), Severity: Medium, Capability: CapObfuscation},
}
