package content

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

const sitonesFile = "sitones.toml"

var (
	validSitoneTypes = map[string]struct{}{
		"exploration":   {},
		"inspiration":   {},
		"resonance":     {},
		"engineering":   {},
		"entertainment": {},
	}
	validRarities = map[string]struct{}{
		"base":    {},
		"common":  {},
		"rare":    {},
		"limited": {},
	}
)

type Store struct {
	sitones     []Sitone
	sitonesByID map[string]Sitone
}

type Sitone struct {
	ID          string `toml:"id"`
	Name        string `toml:"name"`
	Type        string `toml:"type"`
	Rarity      string `toml:"rarity"`
	Style       string `toml:"style"`
	Description string `toml:"description"`
}

type sitonesDocument struct {
	Sitones []Sitone `toml:"sitones"`
}

func Load(dir string) (*Store, error) {
	resolvedDir, err := resolveDir(dir)
	if err != nil {
		return nil, err
	}

	sitonesPath := filepath.Join(resolvedDir, sitonesFile)
	doc, err := loadTOMLFile[sitonesDocument](sitonesPath)
	if err != nil {
		return nil, err
	}

	sitones, sitonesByID, err := validateSitones(sitonesPath, doc.Sitones)
	if err != nil {
		return nil, err
	}

	return &Store{
		sitones:     sitones,
		sitonesByID: sitonesByID,
	}, nil
}

func (s *Store) ListSitones() []Sitone {
	if s == nil || len(s.sitones) == 0 {
		return nil
	}

	sitones := make([]Sitone, len(s.sitones))
	copy(sitones, s.sitones)
	return sitones
}

func (s *Store) GetSitone(id string) (Sitone, bool) {
	if s == nil {
		return Sitone{}, false
	}

	sitone, ok := s.sitonesByID[id]
	return sitone, ok
}

func resolveDir(dir string) (string, error) {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return "", errors.New("content dir is required")
	}

	for _, candidate := range contentDirCandidates(dir) {
		stat, err := os.Stat(candidate)
		if err == nil {
			if !stat.IsDir() {
				return "", fmt.Errorf("content dir %q is not a directory", candidate)
			}
			return candidate, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("stat content dir %q: %w", candidate, err)
		}
	}

	return "", fmt.Errorf("content dir %q does not exist", dir)
}

func contentDirCandidates(dir string) []string {
	candidates := []string{dir}
	if filepath.IsAbs(dir) {
		return candidates
	}

	if stripped, ok := strings.CutPrefix(dir, "server"+string(filepath.Separator)); ok {
		candidates = append(candidates, stripped)
	} else {
		candidates = append(candidates, filepath.Join("server", dir))
	}
	return candidates
}

func loadTOMLFile[T any](path string) (T, error) {
	var out T

	data, err := os.ReadFile(path)
	if err != nil {
		return out, fmt.Errorf("%s: read: %w", path, err)
	}
	if err := toml.Unmarshal(data, &out); err != nil {
		return out, fmt.Errorf("%s: decode toml: %w", path, err)
	}

	return out, nil
}

func validateSitones(path string, sitones []Sitone) ([]Sitone, map[string]Sitone, error) {
	var errs []error
	seen := make(map[string]struct{}, len(sitones))
	normalized := make([]Sitone, 0, len(sitones))

	for i, sitone := range sitones {
		sitone = normalizeSitone(sitone)
		location := fmt.Sprintf("%s: sitones[%d]", path, i)

		if sitone.ID == "" {
			errs = append(errs, fmt.Errorf("%s.id is required", location))
		} else if _, ok := seen[sitone.ID]; ok {
			errs = append(errs, fmt.Errorf("%s: duplicate sitone id %q", path, sitone.ID))
		} else {
			seen[sitone.ID] = struct{}{}
		}
		if sitone.Name == "" {
			errs = append(errs, fmt.Errorf("%s.name is required", location))
		}
		if sitone.Type == "" {
			errs = append(errs, fmt.Errorf("%s.type is required", location))
		} else if _, ok := validSitoneTypes[sitone.Type]; !ok {
			errs = append(errs, fmt.Errorf("%s.type must be one of %s", location, sortedKeys(validSitoneTypes)))
		}
		if sitone.Rarity == "" {
			errs = append(errs, fmt.Errorf("%s.rarity is required", location))
		} else if _, ok := validRarities[sitone.Rarity]; !ok {
			errs = append(errs, fmt.Errorf("%s.rarity must be one of %s", location, sortedKeys(validRarities)))
		}

		normalized = append(normalized, sitone)
	}

	if err := errors.Join(errs...); err != nil {
		return nil, nil, err
	}

	sort.Slice(normalized, func(i, j int) bool {
		return normalized[i].ID < normalized[j].ID
	})

	byID := make(map[string]Sitone, len(normalized))
	for _, sitone := range normalized {
		byID[sitone.ID] = sitone
	}

	return normalized, byID, nil
}

func normalizeSitone(sitone Sitone) Sitone {
	sitone.ID = strings.TrimSpace(sitone.ID)
	sitone.Name = strings.TrimSpace(sitone.Name)
	sitone.Type = strings.TrimSpace(sitone.Type)
	sitone.Rarity = strings.TrimSpace(sitone.Rarity)
	sitone.Style = strings.TrimSpace(sitone.Style)
	sitone.Description = strings.TrimSpace(sitone.Description)
	return sitone
}

func sortedKeys(values map[string]struct{}) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
