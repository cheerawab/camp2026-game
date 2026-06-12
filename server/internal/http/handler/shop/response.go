package shop

type ItemListResponse struct {
	Items []ShopItemResponse `json:"items"`
}

type ItemDetailResponse struct {
	Item ShopItemResponse `json:"item"`
}

type ShopItemResponse struct {
	ID             string `json:"id" example:"item_adventure_backpack"`
	Name           string `json:"name" example:"冒險背包"`
	Type           string `json:"type" example:"material"`
	Rarity         string `json:"rarity" example:"common"`
	Description    string `json:"description" example:"冒險背包，可用於小石合成。"`
	PriceOpenPower int    `json:"priceOpenPower" example:"50"`
	Redeemed       bool   `json:"redeemed" example:"false"`
}

type PurchaseRequest struct {
	ItemID string `json:"itemId" validate:"required" example:"item_adventure_backpack"`
}

type PurchaseResponse struct {
	PurchaseID     string           `json:"purchaseId" example:"purchase_abc123"`
	ItemID         string           `json:"itemId" example:"item_adventure_backpack"`
	Quantity       int              `json:"quantity" example:"1"`
	PriceOpenPower int              `json:"priceOpenPower" example:"50"`
	OpenPower      int              `json:"openPower" example:"1230"`
	Item           ShopItemResponse `json:"item"`
}
