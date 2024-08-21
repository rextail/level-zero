package models

type Order struct {
	OrderUID          string   `json:"order_uid" validate:"required"`
	TrackNumber       string   `json:"track_number" validate:"required"`
	Entry             string   `json:"entry" validate:"required"`
	Delivery          Delivery `json:"delivery" validate:"required"`
	Payment           Payment  `json:"payment" validate:"required"`
	Items             []Item   `json:"items" validate:"required,dive,required"`
	Locale            string   `json:"locale" validate:"required"`
	InternalSignature string   `json:"internal_signature"`
	CustomerID        string   `json:"customer_id" validate:"required"`
	DeliveryService   string   `json:"delivery_service" validate:"required"`
	ShardKey          string   `json:"shardkey" validate:"required"`
	SMID              int      `json:"sm_id" validate:"required"`
	DateCreated       string   `json:"date_created" validate:"required,validateDate"`
	OofShard          string   `json:"oof_shard" validate:"required"`
}

type Delivery struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Zip     string `json:"zip" validate:"required"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

type Payment struct {
	Transaction  string  `json:"transaction" validate:"required"`
	RequestID    string  `json:"request_id"`
	Currency     string  `json:"currency" validate:"required"`
	Provider     string  `json:"provider" validate:"required"`
	Amount       float64 `json:"amount" validate:"required,gt=0"`
	PaymentDT    int64   `json:"payment_dt" validate:"required"`
	Bank         string  `json:"bank" validate:"required"`
	DeliveryCost float64 `json:"delivery_cost" validate:"gte=0"`
	GoodsTotal   float64 `json:"goods_total" validate:"gt=0"`
	CustomFee    *int    `json:"custom_fee" validate:"required,gte=0"`
}

type Item struct {
	ChrtID      int     `json:"chrt_id" validate:"required"`
	TrackNumber string  `json:"track_number" validate:"required"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	RID         string  `json:"rid" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Sale        int     `json:"sale" validate:"required"`
	Size        string  `json:"size" validate:"required"`
	TotalPrice  float64 `json:"total_price" validate:"required,gt=0"`
	NmID        int     `json:"nm_id" validate:"required"`
	Brand       string  `json:"brand" validate:"required"`
	Status      int     `json:"status" validate:"required"`
}
