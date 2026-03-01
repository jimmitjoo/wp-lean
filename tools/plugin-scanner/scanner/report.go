package scanner

import (
  "encoding/json"
  "fmt"
  "io"
  "sort"
  "strings"
)

func WriteJSON(w io.Writer, result *Result) error {
  enc := json.NewEncoder(w)
  enc.SetIndent("", "  ")
  return enc.Encode(result)
}

func WriteSummary(w io.Writer, result *Result) {
  fmt.Fprintf(w, "\nPlugin Security Scan Results\n")
  fmt.Fprintf(w, "%s\n\n", strings.Repeat("=", 50))
  fmt.Fprintf(w, "Files scanned: %d\n", result.FilesScanned)
  fmt.Fprintf(w, "Total findings: %d\n\n", len(result.Findings))

  if len(result.Findings) == 0 {
    fmt.Fprintf(w, "No security issues found.\n")
    return
  }

  fmt.Fprintf(w, "  CRITICAL: %d\n", result.CriticalCount)
  fmt.Fprintf(w, "  HIGH:     %d\n", result.HighCount)
  fmt.Fprintf(w, "  MEDIUM:   %d\n\n", result.MediumCount)

  sort.Slice(result.Findings, func(i, j int) bool {
    if result.Findings[i].Severity != result.Findings[j].Severity {
      return result.Findings[i].Severity < result.Findings[j].Severity
    }
    if result.Findings[i].File != result.Findings[j].File {
      return result.Findings[i].File < result.Findings[j].File
    }
    return result.Findings[i].Line < result.Findings[j].Line
  })

  currentSeverity := Severity(-1)
  for _, f := range result.Findings {
    if f.Severity != currentSeverity {
      currentSeverity = f.Severity
      fmt.Fprintf(w, "--- %s ---\n", f.SeverityString())
    }
    fmt.Fprintf(w, "  %s:%d\n", f.File, f.Line)
    fmt.Fprintf(w, "    Pattern: %s\n", f.Pattern)
    fmt.Fprintf(w, "    Code:    %s\n\n", f.Content)
  }

  if result.CriticalCount > 0 {
    fmt.Fprintf(w, "VERDICT: REJECT - %d critical issues found\n", result.CriticalCount)
  } else if result.HighCount > 0 {
    fmt.Fprintf(w, "VERDICT: REVIEW REQUIRED - %d high-severity issues\n", result.HighCount)
  } else {
    fmt.Fprintf(w, "VERDICT: PASS with %d minor notes\n", result.MediumCount)
  }
}
