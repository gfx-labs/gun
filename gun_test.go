package gun_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gfx-labs/gun"
)

// clearEnv unsets all env vars that could interfere with tests.
func clearEnv(keys ...string) {
	for _, k := range keys {
		os.Unsetenv(k)
	}
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(orig) })
}

func TestDefaults(t *testing.T) {
	clearEnv("HOST", "PORT", "DB_URL")

	tmp := t.TempDir()
	chdir(t, tmp)

	var cfg struct {
		Host string `default:"localhost"`
		Port int    `default:"8080"`
	}
	gun.Load(&cfg)

	if cfg.Host != "localhost" {
		t.Errorf("expected Host=localhost, got %q", cfg.Host)
	}
	if cfg.Port != 8080 {
		t.Errorf("expected Port=8080, got %d", cfg.Port)
	}
}

func TestEnvVars(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	t.Setenv("HOST", "envhost")
	t.Setenv("PORT", "9090")
	t.Setenv("DB_URL", "postgres://localhost/mydb")

	var cfg struct {
		Host  string
		Port  int
		DbUrl string
	}
	gun.Load(&cfg)

	if cfg.Host != "envhost" {
		t.Errorf("expected Host=envhost, got %q", cfg.Host)
	}
	if cfg.Port != 9090 {
		t.Errorf("expected Port=9090, got %d", cfg.Port)
	}
	if cfg.DbUrl != "postgres://localhost/mydb" {
		t.Errorf("expected DbUrl=postgres://localhost/mydb, got %q", cfg.DbUrl)
	}
}

func TestYAMLFile(t *testing.T) {
	clearEnv("HOST", "PORT", "ITEMS", "NESTED_VALUE")

	tmp := t.TempDir()
	chdir(t, tmp)

	yamlContent := `Host: yamlhost
Port: 3000
Items: alpha,bravo,charlie
Nested:
  Value: deep
`
	if err := os.WriteFile(filepath.Join(tmp, "config.yml"), []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	var cfg struct {
		Host   string
		Port   int
		Items  []string
		Nested struct {
			Value string
		}
	}
	gun.Load(&cfg)

	if cfg.Host != "yamlhost" {
		t.Errorf("expected Host=yamlhost, got %q", cfg.Host)
	}
	if cfg.Port != 3000 {
		t.Errorf("expected Port=3000, got %d", cfg.Port)
	}
	if len(cfg.Items) != 3 || cfg.Items[0] != "alpha" || cfg.Items[1] != "bravo" || cfg.Items[2] != "charlie" {
		t.Errorf("expected Items=[alpha bravo charlie], got %v", cfg.Items)
	}
	if cfg.Nested.Value != "deep" {
		t.Errorf("expected Nested.Value=deep, got %q", cfg.Nested.Value)
	}
}

func TestYAMLWithTags(t *testing.T) {
	clearEnv("FIELD_ONE", "FIELD_2")

	tmp := t.TempDir()
	chdir(t, tmp)

	yamlContent := `field_one: tagged_value
Field2: 42
`
	if err := os.WriteFile(filepath.Join(tmp, "config.yaml"), []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	var cfg struct {
		Field1 string `yaml:"field_one"`
		Field2 int32
	}
	gun.Load(&cfg)

	if cfg.Field1 != "tagged_value" {
		t.Errorf("expected Field1=tagged_value, got %q", cfg.Field1)
	}
	if cfg.Field2 != 42 {
		t.Errorf("expected Field2=42, got %d", cfg.Field2)
	}
}

func TestJSONFile(t *testing.T) {
	clearEnv("HOST", "PORT", "ENABLED")

	tmp := t.TempDir()
	chdir(t, tmp)

	jsonContent := `{
  "Host": "jsonhost",
  "Port": 4000,
  "Enabled": true
}`
	if err := os.WriteFile(filepath.Join(tmp, "config.json"), []byte(jsonContent), 0644); err != nil {
		t.Fatal(err)
	}

	var cfg struct {
		Host    string
		Port    int
		Enabled bool
	}
	gun.Load(&cfg)

	if cfg.Host != "jsonhost" {
		t.Errorf("expected Host=jsonhost, got %q", cfg.Host)
	}
	if cfg.Port != 4000 {
		t.Errorf("expected Port=4000, got %d", cfg.Port)
	}
	if !cfg.Enabled {
		t.Errorf("expected Enabled=true, got false")
	}
}

func TestEnvOverridesFile(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	yamlContent := `Host: fromfile
Port: 3000
`
	if err := os.WriteFile(filepath.Join(tmp, "config.yml"), []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("HOST", "fromenv")

	var cfg struct {
		Host string
		Port int
	}
	gun.Load(&cfg)

	if cfg.Host != "fromenv" {
		t.Errorf("expected Host=fromenv (env override), got %q", cfg.Host)
	}
	if cfg.Port != 3000 {
		t.Errorf("expected Port=3000 (from file), got %d", cfg.Port)
	}
}

func TestEnvOverridesDefault(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	t.Setenv("HOST", "fromenv")

	var cfg struct {
		Host string `default:"defaulthost"`
	}
	gun.Load(&cfg)

	if cfg.Host != "fromenv" {
		t.Errorf("expected Host=fromenv, got %q", cfg.Host)
	}
}

func TestPrefix(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	t.Setenv("MYAPP_HOST", "prefixed")
	t.Setenv("MYAPP_PORT", "7070")
	// unprefixed should not be picked up
	t.Setenv("HOST", "nope")

	var cfg struct {
		Host string
		Port int
	}
	gun.LoadPrefix(&cfg, "MYAPP")

	if cfg.Host != "prefixed" {
		t.Errorf("expected Host=prefixed, got %q", cfg.Host)
	}
	if cfg.Port != 7070 {
		t.Errorf("expected Port=7070, got %d", cfg.Port)
	}
}

func TestPrefixWithYAMLFile(t *testing.T) {
	clearEnv("HOST", "PORT", "MYAPP_HOST", "MYAPP_PORT")

	tmp := t.TempDir()
	chdir(t, tmp)

	// With prefix "MYAPP", gun looks for MYAPP.yml instead of config.yml
	yamlContent := `Host: fromprefixfile
Port: 5555
`
	if err := os.WriteFile(filepath.Join(tmp, "MYAPP.yml"), []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	var cfg struct {
		Host string
		Port int
	}
	gun.LoadPrefix(&cfg, "MYAPP")

	if cfg.Host != "fromprefixfile" {
		t.Errorf("expected Host=fromprefixfile, got %q", cfg.Host)
	}
	if cfg.Port != 5555 {
		t.Errorf("expected Port=5555, got %d", cfg.Port)
	}
}

func TestCommaSlice(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	t.Setenv("TAGS", "one,two,three")

	var cfg struct {
		Tags []string
	}
	gun.Load(&cfg)

	if len(cfg.Tags) != 3 {
		t.Fatalf("expected 3 tags, got %d: %v", len(cfg.Tags), cfg.Tags)
	}
	expected := []string{"one", "two", "three"}
	for i, v := range expected {
		if cfg.Tags[i] != v {
			t.Errorf("Tags[%d]: expected %q, got %q", i, v, cfg.Tags[i])
		}
	}
}

func TestYAMLMergesWithEnv(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	yamlContent := `Host: filehost
Port: 1111
Debug: false
`
	if err := os.WriteFile(filepath.Join(tmp, "config.yml"), []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Override only some fields via env
	t.Setenv("PORT", "2222")
	t.Setenv("DEBUG", "true")

	var cfg struct {
		Host  string
		Port  int
		Debug bool
	}
	gun.Load(&cfg)

	if cfg.Host != "filehost" {
		t.Errorf("expected Host=filehost (from yaml), got %q", cfg.Host)
	}
	if cfg.Port != 2222 {
		t.Errorf("expected Port=2222 (env override), got %d", cfg.Port)
	}
	if !cfg.Debug {
		t.Errorf("expected Debug=true (env override), got false")
	}
}

func TestNestedStructFromYAML(t *testing.T) {
	clearEnv("DATABASE_HOST", "DATABASE_PORT", "DATABASE_NAME")

	tmp := t.TempDir()
	chdir(t, tmp)

	yamlContent := `Database:
  Host: dbhost
  Port: 5432
  Name: mydb
`
	if err := os.WriteFile(filepath.Join(tmp, "config.yml"), []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	var cfg struct {
		Database struct {
			Host string
			Port int
			Name string
		}
	}
	gun.Load(&cfg)

	if cfg.Database.Host != "dbhost" {
		t.Errorf("expected Database.Host=dbhost, got %q", cfg.Database.Host)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("expected Database.Port=5432, got %d", cfg.Database.Port)
	}
	if cfg.Database.Name != "mydb" {
		t.Errorf("expected Database.Name=mydb, got %q", cfg.Database.Name)
	}
}

func TestNestedStructFromEnv(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	t.Setenv("DATABASE_HOST", "envdbhost")
	t.Setenv("DATABASE_PORT", "3306")

	var cfg struct {
		Database struct {
			Host string
			Port int
		}
	}
	gun.Load(&cfg)

	if cfg.Database.Host != "envdbhost" {
		t.Errorf("expected Database.Host=envdbhost, got %q", cfg.Database.Host)
	}
	if cfg.Database.Port != 3306 {
		t.Errorf("expected Database.Port=3306, got %d", cfg.Database.Port)
	}
}

func TestEmptyConfig(t *testing.T) {
	clearEnv("HOST", "PORT")

	tmp := t.TempDir()
	chdir(t, tmp)

	var cfg struct {
		Host string
		Port int
	}
	gun.Load(&cfg)

	if cfg.Host != "" {
		t.Errorf("expected Host empty, got %q", cfg.Host)
	}
	if cfg.Port != 0 {
		t.Errorf("expected Port=0, got %d", cfg.Port)
	}
}
