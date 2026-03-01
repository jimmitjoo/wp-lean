package scanner

import (
  "testing"
)

func TestDetectedCapabilities(t *testing.T) {
  result, err := ScanDirectory(testdataDir())
  if err != nil {
    t.Fatalf("scan failed: %v", err)
  }

  if len(result.Capabilities) == 0 {
    t.Fatal("expected detected capabilities, got none")
  }

  capSet := map[Capability]bool{}
  for _, c := range result.Capabilities {
    capSet[c] = true
  }

  expected := []Capability{CapExec, CapEval, CapHTTPOutbound}
  for _, e := range expected {
    if !capSet[e] {
      t.Errorf("expected capability %s to be detected", e)
    }
  }

  t.Logf("Detected capabilities: %v", result.Capabilities)
}

func TestCheckCapabilitiesAllDeclared(t *testing.T) {
  result, err := ScanDirectory(testdataDir())
  if err != nil {
    t.Fatalf("scan failed: %v", err)
  }

  violations := CheckCapabilities(result, result.Capabilities)
  if len(violations) != 0 {
    t.Errorf("expected no violations when all capabilities declared, got %d", len(violations))
    for _, v := range violations {
      t.Logf("  undeclared: %s (%d findings)", v.Capability, len(v.Findings))
    }
  }
}

func TestCheckCapabilitiesUndeclared(t *testing.T) {
  result, err := ScanDirectory(testdataDir())
  if err != nil {
    t.Fatalf("scan failed: %v", err)
  }

  violations := CheckCapabilities(result, []Capability{CapHTTPOutbound})

  if len(violations) == 0 {
    t.Fatal("expected violations for undeclared capabilities")
  }

  violatedCaps := map[Capability]bool{}
  for _, v := range violations {
    violatedCaps[v.Capability] = true
  }

  if violatedCaps[CapHTTPOutbound] {
    t.Error("http-outbound was declared but still reported as violation")
  }
  if !violatedCaps[CapExec] {
    t.Error("exec should be a violation (not declared)")
  }
  if !violatedCaps[CapEval] {
    t.Error("eval should be a violation (not declared)")
  }

  t.Logf("Violations: %d undeclared capabilities", len(violations))
}
