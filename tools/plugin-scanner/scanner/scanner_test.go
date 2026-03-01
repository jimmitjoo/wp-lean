package scanner

import (
  "path/filepath"
  "runtime"
  "testing"
)

func testdataDir() string {
  _, filename, _, _ := runtime.Caller(0)
  return filepath.Join(filepath.Dir(filename), "..", "testdata")
}

func TestScanMaliciousFile(t *testing.T) {
  result, err := ScanDirectory(testdataDir())
  if err != nil {
    t.Fatalf("scan failed: %v", err)
  }

  if result.CriticalCount == 0 {
    t.Error("expected critical findings in malicious.php")
  }
  if result.HighCount == 0 {
    t.Error("expected high findings in malicious.php")
  }
  if result.MediumCount == 0 {
    t.Error("expected medium findings in malicious.php")
  }

  t.Logf("Found: %d critical, %d high, %d medium", result.CriticalCount, result.HighCount, result.MediumCount)
}

func TestCleanFileNoFindings(t *testing.T) {
  result, err := ScanDirectory(testdataDir())
  if err != nil {
    t.Fatalf("scan failed: %v", err)
  }

  cleanPath := filepath.Join(testdataDir(), "clean.php")
  for _, f := range result.Findings {
    if f.File == cleanPath {
      t.Errorf("clean.php should have no findings, got: %s at line %d", f.Pattern, f.Line)
    }
  }
}
