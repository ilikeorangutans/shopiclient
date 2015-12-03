package shopify

import (
	"time"
)

type Timestamps struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommonFields struct {
	ID     int64 `json:"id"`
	ShopID int64 `json:"shop_id"`
	Timestamps
}
