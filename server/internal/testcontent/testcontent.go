package testcontent

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/sitcon-tw/camp2026-game/internal/content"
)

const quizQuestionsCSV = `id,prompt,choice_a,choice_b,choice_c,choice_d,correct_choice,explanation
quiz-001,Question 1,A,B,C,D,A,Explanation 1
quiz-002,Question 2,A,B,C,D,B,Explanation 2
quiz-003,Question 3,A,B,C,D,C,Explanation 3
quiz-004,Question 4,A,B,C,D,D,Explanation 4
quiz-005,Question 5,A,B,C,D,A,Explanation 5
quiz-006,Question 6,A,B,C,D,B,Explanation 6
quiz-007,Question 7,A,B,C,D,C,Explanation 7
quiz-008,Question 8,A,B,C,D,D,Explanation 8
quiz-009,Question 9,A,B,C,D,A,Explanation 9
quiz-010,Question 10,A,B,C,D,B,Explanation 10
`

func Load(t testing.TB) *content.Store {
	t.Helper()

	store, err := content.Load(Dir(t))
	if err != nil {
		t.Fatalf("load test content: %v", err)
	}
	return store
}

func Dir(t testing.TB) string {
	t.Helper()

	dir := t.TempDir()
	sourceDir := sourceContentDir(t)

	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		t.Fatalf("read source content dir: %v", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == "quiz_questions.csv" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(sourceDir, entry.Name()))
		if err != nil {
			t.Fatalf("read source content %s: %v", entry.Name(), err)
		}
		if err := os.WriteFile(filepath.Join(dir, entry.Name()), data, 0o644); err != nil {
			t.Fatalf("write test content %s: %v", entry.Name(), err)
		}
	}

	if err := os.WriteFile(filepath.Join(dir, "quiz_questions.csv"), []byte(quizQuestionsCSV), 0o644); err != nil {
		t.Fatalf("write test quiz questions: %v", err)
	}
	return dir
}

func sourceContentDir(t testing.TB) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve testcontent path")
	}

	dir := filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "content"))
	stat, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("stat source content dir %s: %v", dir, err)
	}
	if !stat.IsDir() {
		t.Fatalf("%s is not a directory", dir)
	}
	if _, err := os.Stat(filepath.Join(dir, "sitones.toml")); err != nil {
		t.Fatalf("source content dir %s is missing sitones.toml: %v", dir, err)
	}
	return dir
}
