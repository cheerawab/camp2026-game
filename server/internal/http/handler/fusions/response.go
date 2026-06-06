package fusions

type RecipeListResponse struct {
	Recipes []FusionRecipeResponse `json:"recipes"`
}

type FusionRecipeResponse struct {
	ID          string                    `json:"id" example:"fusion-engineering-route-frame"`
	Name        string                    `json:"name" example:"工程路線展示框"`
	Description string                    `json:"description" example:"把工程小石與合成碎片組合成基地展示用的收藏外框。"`
	Enabled     bool                      `json:"enabled" example:"true"`
	Available   bool                      `json:"available" example:"true"`
	Inputs      []FusionComponentResponse `json:"inputs"`
	Outputs     []FusionComponentResponse `json:"outputs"`
}

type FusionComponentResponse struct {
	Kind     string `json:"kind" example:"item"`
	ID       string `json:"id" example:"item-crafting-fragment"`
	Name     string `json:"name" example:"合成碎片"`
	Type     string `json:"type,omitempty" example:"material"`
	Rarity   string `json:"rarity,omitempty" example:"common"`
	Quantity int    `json:"quantity" example:"3"`
}

type CreateRequest struct {
	RecipeID string `json:"recipeId" validate:"required" example:"fusion-engineering-route-frame"`
}

type CreateResponse struct {
	FusionID string               `json:"fusionId" example:"fusion_abc123"`
	Recipe   FusionRecipeResponse `json:"recipe"`
}
