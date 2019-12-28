package presenter

// API 統一回傳結構
type Result struct {
	Data  interface{} `json:"data"`
	Code  StatusCode  `json:"code"` // 回應狀態碼
	Count int         `json:"count"`
}

type StatusCode uint16

// 回應狀態碼
const (
	StatusSuccess = 2001 // 下單成功

	StatusParamError = 4001 // 參數有誤

	StatusSellOut        = 5031 // 目前已經賣完了
	StatusNotAllowToSell = 5032 // 目前不能購買
	StatusOverBuy        = 5033 // 超過目前可購買數
)
