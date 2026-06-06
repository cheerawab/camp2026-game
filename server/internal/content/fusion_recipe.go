package content

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

const (
	fusionRecipesFile = "fusion_recipes.toml"
	FusionKindItem    = "item"
	FusionKindSitone  = "sitone"
)

var validFusionKinds = map[string]struct{}{
	FusionKindItem:   {},
	FusionKindSitone: {},
}

type FusionRecipe struct {
	ID          string            `toml:"id"`
	Name        string            `toml:"name"`
	Description string            `toml:"description"`
	Enabled     bool              `toml:"enabled"`
	Inputs      []FusionComponent `toml:"inputs"`
	Outputs     []FusionComponent `toml:"outputs"`
}

type FusionComponent struct {
	Kind     string `toml:"kind"`
	ID       string `toml:"id"`
	Quantity int    `toml:"quantity"`
}

type fusionRecipesDocument struct {
	FusionRecipes []FusionRecipe `toml:"fusion_recipes"`
}

func (s *Store) ListFusionRecipes() []FusionRecipe {
	if s == nil || len(s.fusionRecipes) == 0 {
		return nil
	}

	recipes := make([]FusionRecipe, len(s.fusionRecipes))
	copy(recipes, s.fusionRecipes)
	for i := range recipes {
		recipes[i].Inputs = copyFusionComponents(recipes[i].Inputs)
		recipes[i].Outputs = copyFusionComponents(recipes[i].Outputs)
	}
	return recipes
}

func (s *Store) GetFusionRecipe(id string) (FusionRecipe, bool) {
	if s == nil {
		return FusionRecipe{}, false
	}

	recipe, ok := s.fusionRecipesByID[id]
	if !ok {
		return FusionRecipe{}, false
	}
	recipe.Inputs = copyFusionComponents(recipe.Inputs)
	recipe.Outputs = copyFusionComponents(recipe.Outputs)
	return recipe, true
}

func validateFusionRecipes(
	path string,
	recipes []FusionRecipe,
	sitonesByID map[string]Sitone,
	itemsByID map[string]Item,
) ([]FusionRecipe, map[string]FusionRecipe, error) {
	var errs []error
	seen := make(map[string]struct{}, len(recipes))
	normalized := make([]FusionRecipe, 0, len(recipes))

	for i, recipe := range recipes {
		recipe = normalizeFusionRecipe(recipe)
		location := fmt.Sprintf("%s: fusion_recipes[%d]", path, i)

		if recipe.ID == "" {
			errs = append(errs, fmt.Errorf("%s.id is required", location))
		} else if _, ok := seen[recipe.ID]; ok {
			errs = append(errs, fmt.Errorf("%s: duplicate fusion recipe id %q", path, recipe.ID))
		} else {
			seen[recipe.ID] = struct{}{}
		}
		if recipe.Name == "" {
			errs = append(errs, fmt.Errorf("%s.name is required", location))
		}
		if len(recipe.Inputs) == 0 {
			errs = append(errs, fmt.Errorf("%s.inputs is required", location))
		}
		if len(recipe.Outputs) == 0 {
			errs = append(errs, fmt.Errorf("%s.outputs is required", location))
		}
		errs = append(errs, validateFusionComponents(location+".inputs", recipe.Inputs, sitonesByID, itemsByID)...)
		errs = append(errs, validateFusionComponents(location+".outputs", recipe.Outputs, sitonesByID, itemsByID)...)

		normalized = append(normalized, recipe)
	}

	if err := errors.Join(errs...); err != nil {
		return nil, nil, err
	}

	sort.Slice(normalized, func(i, j int) bool {
		return normalized[i].ID < normalized[j].ID
	})

	byID := make(map[string]FusionRecipe, len(normalized))
	for _, recipe := range normalized {
		byID[recipe.ID] = recipe
	}
	return normalized, byID, nil
}

func validateFusionComponents(
	location string,
	components []FusionComponent,
	sitonesByID map[string]Sitone,
	itemsByID map[string]Item,
) []error {
	var errs []error
	for i, component := range components {
		componentLocation := fmt.Sprintf("%s[%d]", location, i)
		if component.Kind == "" {
			errs = append(errs, fmt.Errorf("%s.kind is required", componentLocation))
		} else if _, ok := validFusionKinds[component.Kind]; !ok {
			errs = append(errs, fmt.Errorf("%s.kind must be one of %s", componentLocation, sortedKeys(validFusionKinds)))
		}
		if component.ID == "" {
			errs = append(errs, fmt.Errorf("%s.id is required", componentLocation))
		}
		if component.Quantity <= 0 {
			errs = append(errs, fmt.Errorf("%s.quantity must be greater than 0", componentLocation))
		}
		if component.Kind == FusionKindSitone && component.ID != "" {
			if _, ok := sitonesByID[component.ID]; !ok {
				errs = append(errs, fmt.Errorf("%s.id references unknown sitone %q", componentLocation, component.ID))
			}
		}
		if component.Kind == FusionKindItem && component.ID != "" {
			if _, ok := itemsByID[component.ID]; !ok {
				errs = append(errs, fmt.Errorf("%s.id references unknown item %q", componentLocation, component.ID))
			}
		}
	}
	return errs
}

func normalizeFusionRecipe(recipe FusionRecipe) FusionRecipe {
	recipe.ID = strings.TrimSpace(recipe.ID)
	recipe.Name = strings.TrimSpace(recipe.Name)
	recipe.Description = strings.TrimSpace(recipe.Description)
	for i := range recipe.Inputs {
		recipe.Inputs[i] = normalizeFusionComponent(recipe.Inputs[i])
	}
	for i := range recipe.Outputs {
		recipe.Outputs[i] = normalizeFusionComponent(recipe.Outputs[i])
	}
	return recipe
}

func normalizeFusionComponent(component FusionComponent) FusionComponent {
	component.Kind = strings.TrimSpace(component.Kind)
	component.ID = strings.TrimSpace(component.ID)
	return component
}

func copyFusionComponents(components []FusionComponent) []FusionComponent {
	if len(components) == 0 {
		return nil
	}
	out := make([]FusionComponent, len(components))
	copy(out, components)
	return out
}
