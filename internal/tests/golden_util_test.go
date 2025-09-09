package tests

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"gotest.tools/v3/golden"
)

type TestCase struct {
	Name    string
	Engine  string
	Package string
}

const YAML_TEMPLATE = `
version: "2"
plugins:
  - name: php
    wasm:
      url: file://%s
sql:
- schema: schema.sql
  queries: query.sql
  engine: %s
  codegen:
  - out: generated
    plugin: php
    options:
      package: "%s"
`

func runGoldenTest(t *testing.T, tc TestCase) {
	t.Helper()

	ensureWASMPlugin(t)

	tempDir := t.TempDir()
	copyTestFiles(t, tc.Name, tempDir)
	createSQLCConfig(t, tempDir, tc)
	runSQLCGenerate(t, tempDir)
	compareWithGolden(t, tc, tempDir)
}

func ensureWASMPlugin(t *testing.T) {
	t.Helper()

	// We are in internal/tests, so go up two levels to the project root
	_, b, _, _ := runtime.Caller(0)
	projectRoot := path.Join(filepath.Dir(b), "..", "..")

	// Build the WASM plugin
	wasmPath := filepath.Join(projectRoot, "bin", "sqlc-gen-php.wasm")
	cmd := exec.Command("go", "build", "-o", wasmPath, "main.go")
	cmd.Dir = filepath.Join(projectRoot, "plugin")
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build WASM plugin: %v", err)
	}
}

func copyTestFiles(t *testing.T, testName, destDir string) {
	t.Helper()

	srcDir := filepath.Join("testdata", testName)
	copyFile(t, filepath.Join(srcDir, "schema.sql"), filepath.Join(destDir, "schema.sql"))
	copyFile(t, filepath.Join(srcDir, "query.sql"), filepath.Join(destDir, "query.sql"))
}

func copyFile(t *testing.T, src, dest string) {
	t.Helper()

	srcFile, err := os.Open(src)
	if err != nil {
		t.Fatalf("Failed to open source file %s: %v", src, err)
	}

	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		t.Fatalf("Failed to create dest file %s: %v", dest, err)
	}

	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		t.Fatalf("Failed to copy file: %v", err)
	}
}

func createSQLCConfig(t *testing.T, dir string, tc TestCase) {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	projectRoot := wd
	if strings.HasSuffix(wd, "internal/tests") {
		projectRoot = filepath.Join(wd, "..", "..")
	}

	projectRoot, err = filepath.Abs(projectRoot)
	if err != nil {
		t.Fatalf("Failed to get absolute project root: %v", err)
	}

	wasmPath := filepath.Join(projectRoot, "bin", "sqlc-gen-php.wasm")
	config := fmt.Sprintf(
		YAML_TEMPLATE,
		wasmPath,
		tc.Engine,
		strings.ReplaceAll(tc.Package, `\`, `\\`),
	)

	configPath := filepath.Join(dir, "sqlc.yaml")
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatalf("Failed to write sqlc config: %v", err)
	}
}

func runSQLCGenerate(t *testing.T, dir string) {
	t.Helper()
	cmd := exec.Command("sqlc", "generate")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("sqlc generate failed: %v\nOutput: %s", err, string(output))
	}
}

func compareWithGolden(t *testing.T, tc TestCase, tempDir string) {
	t.Helper()

	generatedDir := filepath.Join(tempDir, "generated")

	err := filepath.WalkDir(generatedDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(generatedDir, path)
		if err != nil {
			return err
		}

		goldenName := filepath.Join(tc.Name, "expected", relPath)
		golden.Assert(t, string(content), goldenName)
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to compare with golden files: %v", err)
	}
}
