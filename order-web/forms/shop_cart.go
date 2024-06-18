package forms

type ShopCartForm struct {
	GoodsId int32 `json:"goods" binding:"required"`
	Nums    int32 `json:"nums" binding:"required,min=1"`
}
