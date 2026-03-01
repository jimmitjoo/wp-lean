package scanner

import (
  "bufio"
  "os"
  "path/filepath"
  "strings"
)

type Finding struct {
  File       string     `json:"file"`
  Line       int        `json:"line"`
  Severity   Severity   `json:"severity"`
  Pattern    string     `json:"pattern"`
  Content    string     `json:"content"`
  Capability Capability `json:"capability,omitempty"`
}

func (f Finding) SeverityString() string {
  return f.Severity.String()
}

type Violation struct {
  Capability Capability `json:"capability"`
  Findings   []Finding  `json:"findings"`
}

type Result struct {
  Findings      []Finding    `json:"findings"`
  Capabilities  []Capability `json:"capabilities"`
  Violations    []Violation  `json:"violations,omitempty"`
  CriticalCount int          `json:"critical_count"`
  HighCount     int          `json:"high_count"`
  MediumCount   int          `json:"medium_count"`
  FilesScanned  int          `json:"files_scanned"`
}

func ScanDirectory(dir string) (*Result, error) {
  result := &Result{}
  err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return nil
    }
    if info.IsDir() || !strings.HasSuffix(path, ".php") {
      return nil
    }
    findings, err := scanFile(path)
    if err != nil {
      return nil
    }
    result.FilesScanned++
    result.Findings = append(result.Findings, findings...)
    return nil
  })
  if err != nil {
    return nil, err
  }
  capSet := map[Capability]bool{}
  for _, f := range result.Findings {
    switch f.Severity {
    case Critical:
      result.CriticalCount++
    case High:
      result.HighCount++
    case Medium:
      result.MediumCount++
    }
    if f.Capability != "" {
      capSet[f.Capability] = true
    }
  }
  for cap := range capSet {
    result.Capabilities = append(result.Capabilities, cap)
  }
  return result, nil
}

func scanFile(path string) ([]Finding, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var findings []Finding
  s := bufio.NewScanner(file)
  buf := make([]byte, 0, 1024*1024)
  s.Buffer(buf, 1024*1024)
  lineNum := 0

  for s.Scan() {
    lineNum++
    line := s.Text()
    trimmed := strings.TrimSpace(line)

    if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "*") || strings.HasPrefix(trimmed, "/*") {
      continue
    }

    for _, p := range Patterns {
      if p.Regex.MatchString(line) {
        if isSanitized(p.Name, line) {
          continue
        }
        content := trimmed
        if len(content) > 200 {
          content = content[:200] + "..."
        }
        findings = append(findings, Finding{
          File:       path,
          Line:       lineNum,
          Severity:   p.Severity,
          Pattern:    p.Name,
          Content:    content,
          Capability: p.Capability,
        })
      }
    }
  }
  return findings, s.Err()
}

var sanitizers = []string{
  "esc_html", "esc_attr", "esc_url", "esc_sql",
  "sanitize_text_field", "sanitize_email", "sanitize_url",
  "sanitize_key", "sanitize_file_name", "sanitize_title",
  "wp_kses", "wp_unslash", "intval(", "absint(",
  "wp_verify_nonce",
}

var superglobalPatterns = []string{
  "raw $_GET", "raw $_POST", "raw $_REQUEST", "raw $_SERVER",
}

func CheckCapabilities(result *Result, declared []Capability) []Violation {
  declaredSet := map[Capability]bool{}
  for _, c := range declared {
    declaredSet[c] = true
  }
  violationMap := map[Capability][]Finding{}
  for _, f := range result.Findings {
    if f.Capability == "" {
      continue
    }
    if !declaredSet[f.Capability] {
      violationMap[f.Capability] = append(violationMap[f.Capability], f)
    }
  }
  violations := make([]Violation, 0, len(violationMap))
  for cap, findings := range violationMap {
    violations = append(violations, Violation{Capability: cap, Findings: findings})
  }
  return violations
}

func isSanitized(patternName string, line string) bool {
  isSuperglobal := false
  for _, sp := range superglobalPatterns {
    if patternName == sp+" access" {
      isSuperglobal = true
      break
    }
  }
  if !isSuperglobal {
    if patternName == "wpdb->query without prepare" {
      return strings.Contains(line, "->prepare(")
    }
    if patternName == "exec()" {
      return strings.Contains(line, "->exec(")
    }
    return false
  }
  lower := strings.ToLower(line)
  for _, s := range sanitizers {
    if strings.Contains(lower, strings.ToLower(s)) {
      return true
    }
  }
  return false
}
