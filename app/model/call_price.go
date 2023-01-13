package model

type CallPriceReq struct {
	Uids []int64 `json:"uids" validate:"required"`
}

type CallPriceResp struct {
	PriceCoins map[int64]int64 `json:"price_coins"`
}
