package main

import (
  "flag"
  "fmt"
  "os"
  "strings"

  "plugin-scanner/scanner"
)

func main() {
  jsonOutput := flag.Bool("json", false, "Output JSON to stdout")
  checkCaps := flag.String("check-caps", "", "Comma-separated declared capabilities to verify against")
  flag.Parse()

  args := flag.Args()
  if len(args) < 1 {
    fmt.Fprintf(os.Stderr, "Usage: plugin-scanner [flags] <plugin-directory>\n")
    flag.PrintDefaults()
    os.Exit(1)
  }

  dir := args[0]
  info, err := os.Stat(dir)
  if err != nil || !info.IsDir() {
    fmt.Fprintf(os.Stderr, "Error: %s is not a valid directory\n", dir)
    os.Exit(1)
  }

  result, err := scanner.ScanDirectory(dir)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error scanning: %v\n", err)
    os.Exit(1)
  }

  if *checkCaps != "" {
    declared := parseCaps(*checkCaps)
    violations := scanner.CheckCapabilities(result, declared)
    result.Violations = violations
  }

  if *jsonOutput {
    scanner.WriteJSON(os.Stdout, result)
  } else {
    scanner.WriteSummary(os.Stdout, result)
  }

  if result.CriticalCount > 0 {
    os.Exit(2)
  }
  if result.Violations != nil && len(result.Violations) > 0 {
    os.Exit(3)
  }
  if result.HighCount > 0 {
    os.Exit(1)
  }
}

func parseCaps(s string) []scanner.Capability {
  parts := strings.Split(s, ",")
  caps := make([]scanner.Capability, 0, len(parts))
  for _, p := range parts {
    p = strings.TrimSpace(p)
    if p != "" {
      caps = append(caps, scanner.Capability(p))
    }
  }
  return caps
}
