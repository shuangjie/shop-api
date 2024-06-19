package forms

type ShopCartForm struct {
	GoodsId int32 `json:"goods" binding:"required"`
	Nums    int32 `json:"nums" binding:"required,min=1"`
}

type ShopCartUpdateForm struct {
	Nums    int32 `json:"nums" binding:"required,min=1"`
	Checked *bool `json:"checked"` // *bool主要是为了有三种状态：选中、未选中、不传，如果是bool，就只有选中和未选中两种状态，不便于表单判断
}
