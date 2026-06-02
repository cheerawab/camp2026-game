package content

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

const sitonesFile = "sitones.toml"

var validSitoneTypes = map[string]struct{}{
	"exploration":   {},
	"inspiration":   {},
	"resonance":     {},
	"engineering":   {},
	"entertainment": {},
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
