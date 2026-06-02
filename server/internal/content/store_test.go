package content

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadSitones(t *testing.T) {
	dir := writeContent(t, `
[[sitones]]
id = "sitone-engineering"
name = "工程型小石"
type = "engineering"
rarity = "base"
style = "default"
description = "完成技術任務、分享解法或協助除錯。"

[[sitones]]
id = "sitone-exploration"
name = "探索型小石"
type = "exploration"
rarity = "base"
style = "default"
description = "逛攤位、問問題、參與社群事件。"
`)

	store, err := Load(dir)
	if err != nil {
		t.Fatalf("load content: %v", err)
	}

	sitones := store.ListSitones()
	if len(sitones) != 2 {
		t.Fatalf("expected 2 sitones, got %d", len(sitones))
	}
	if sitones[0].ID != "sitone-engineering" || sitones[1].ID != "sitone-exploration" {
		t.Fatalf("expected sitones sorted by id, got %#v", sitones)
	}

	sitone, ok := store.GetSitone("sitone-engineering")
	if !ok {
		t.Fatal("expected sitone-engineering to exist")
	}
	if sitone.Name != "工程型小石" {
		t.Fatalf("expected engineering sitone name, got %q", sitone.Name)
	}
	if _, ok := store.GetSitone("missing"); ok {
		t.Fatal("expected missing sitone not to exist")
	}

	sitones[0].ID = "mutated"
	sitones = store.ListSitones()
	if sitones[0].ID != "sitone-engineering" {
		t.Fatalf("expected ListSitones to return a copy, got %q", sitones[0].ID)
	}
}

func TestLoadRejectsDuplicateSitoneID(t *testing.T) {
	dir := writeContent(t, `
[[sitones]]
id = "sitone-engineering"
name = "Engineering"
type = "engineering"
rarity = "base"

[[sitones]]
id = "sitone-engineering"
name = "Engineering Again"
type = "engineering"
rarity = "base"
`)

	_, err := Load(dir)
	if err == nil {
		t.Fatal("expected duplicate id error")
	}
	if !strings.Contains(err.Error(), `duplicate sitone id "sitone-engineering"`) {
		t.Fatalf("expected duplicate id error, got %v", err)
	}
}

func TestLoadRejectsMissingRequiredSitoneFields(t *testing.T) {
	dir := writeContent(t, `
[[sitones]]
id = ""
name = ""
type = ""
rarity = ""
`)

	_, err := Load(dir)
	if err == nil {
		t.Fatal("expected required field error")
	}
	for _, want := range []string{
		"sitones[0].id is required",
		"sitones[0].name is required",
		"sitones[0].type is required",
		"sitones[0].rarity is required",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("expected error to contain %q, got %v", want, err)
		}
	}
}

func TestLoadRejectsInvalidSitoneEnums(t *testing.T) {
	dir := writeContent(t, `
[[sitones]]
id = "sitone-unknown"
name = "Unknown"
type = "unknown"
rarity = "mythic"
`)

	_, err := Load(dir)
	if err == nil {
		t.Fatal("expected enum error")
	}
	for _, want := range []string{
		"sitones[0].type must be one of",
		"sitones[0].rarity must be one of",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("expected error to contain %q, got %v", want, err)
		}
	}
}

func TestLoadResolvesServerContentFallback(t *testing.T) {
	root := t.TempDir()
	serverDir := filepath.Join(root, "server")
	contentDir := filepath.Join(serverDir, "content")
	if err := os.MkdirAll(contentDir, 0o755); err != nil {
		t.Fatalf("mkdir content dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(contentDir, sitonesFile), []byte(validSitonesTOML()), 0o644); err != nil {
		t.Fatalf("write sitones: %v", err)
	}

	oldCWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}
	if err := os.Chdir(serverDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldCWD); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})

	store, err := Load("server/content")
	if err != nil {
		t.Fatalf("load content through fallback: %v", err)
	}
	if len(store.ListSitones()) != 1 {
		t.Fatalf("expected 1 sitone, got %d", len(store.ListSitones()))
	}
}

func writeContent(t *testing.T, sitones string) string {
	t.Helper()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, sitonesFile), []byte(strings.TrimSpace(sitones)), 0o644); err != nil {
		t.Fatalf("write sitones: %v", err)
	}
	return dir
}

func validSitonesTOML() string {
	return strings.TrimSpace(`
[[sitones]]
id = "sitone-engineering"
name = "Engineering"
type = "engineering"
rarity = "base"
`)
}
