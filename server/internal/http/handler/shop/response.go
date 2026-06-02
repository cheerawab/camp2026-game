package shop

type ItemListResponse struct {
	Items []ShopItemResponse `json:"items"`
}

type ItemDetailResponse struct {
	Item ShopItemResponse `json:"item"`
}

type ShopItemResponse struct {
	ID             string `json:"id" example:"item-crafting-fragment"`
	Name           string `json:"name" example:"合成碎片"`
	Type           string `json:"type" example:"material"`
	Rarity         string `json:"rarity" example:"common"`
	Description    string `json:"description" example:"小石造型合成使用的基礎素材。"`
	PriceOpenPower int    `json:"priceOpenPower" example:"50"`
}

type PurchaseRequest struct {
	ItemID string `json:"itemId" validate:"required" example:"item-crafting-fragment"`
}

type PurchaseResponse struct {
	PurchaseID     string           `json:"purchaseId" example:"purchase_abc123"`
	ItemID         string           `json:"itemId" example:"item-crafting-fragment"`
	Quantity       int              `json:"quantity" example:"1"`
	PriceOpenPower int              `json:"priceOpenPower" example:"50"`
	OpenPower      int              `json:"openPower" example:"1230"`
	Item           ShopItemResponse `json:"item"`
}
