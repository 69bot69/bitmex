package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	BTC_SATOSHI = 100000000 // 1BTC = 100000000 satoshi
)

// save 取引決断になる分析情報
type LineSingle struct {
	Av        float64
	Timestamp time.Time
}

type Indicator struct {
	Setting struct {
		// 証拠金活用割合
		Collateral int
		Leverage   float64

		// 強制損益確定割合
		Check_PL float32

		// 損益確定幅
		Profit int
		Loss   int

		// 板監視範囲
		BoardRange float32

		// 待機注文の絶対値（単位金額は簡素化する）
		OrderAb int64
		// 出来高の絶対値（単位金額は簡素化する）
		VolumeAb int64

		// 取引枚数
		Size int
	}

	InfoStatus struct {
		StartCol        int // 開始時の証拠金
		NowCol          int // 現在の証拠金
		ChangeCol       int // 証拠金変化額
		AvailableMargin int // 残証拠金

		Waiting  int  // 証拠金変動での待機時間
		Progress bool // 強制決済執行中

		// 建玉情報
		Leverage          float32 //レバレッジ設定
		PositionsSymbol   string  // 建玉種類
		PositionsSide     string  // 建玉の売買方向
		PositionsSize     int     // 建玉の枚数
		PositionsAvg      float32 // 建玉の平均取得値
		PositionsCost     float32 // 建玉の証拠金背景
		PositionsCostPcnt float64 // 証拠金使用率
		PositionsPL       float32 // 建玉合計の損益（Profit OR Loss）
		PositionsPLPcnt   float64 // 建玉合計の証拠金全体のROE

		// 取引所
		Health string // 取引所健康状態
		Status string // ボットの判断状態

		// 市場
		VolatilitySize int // 直近の出来高（合算し小数点を捨てる）

		// チャート
		// // 直近板の最大最小値
		PriceMax float32
		PriceMin float32
		// // ４時間足で最大最小値
		PriceMax4h  float32
		PriceMin4h  float32
		PriceDiff4h float32

		Mean       float32 // 直近平均値
		Dispersion float32 // 分散
		Deviation  float32 // 標準偏差

		Yen float32 // 現時点の日本円価格
	}

	StrongSide   int // 初期値: 0 / 買い: 1 / 売り: -1
	VolumeLastID int
	VolumeLength int
	BuyAppear    int
	SellAppear   int
	BuySide      int
	SellSide     int
	ChangeAmount int
	BuyVolume    int
	SellVolume   int
	Exchanges    [3]bool
	ChangePrice  int
	ChangeRatio  int

	// 待機注文分析
	OrderAll       int64   // 待機注文全額
	OrderAllRatio  float32 // 絶対値と比較した割合
	OrderBuy       int64   // 買い待機注文
	OrderSell      int64   // 売り待機注文
	OrderBuyRatio  int     // 買い待機注文比率
	OrderSellRatio int     // 売り待機注文比率

	// 出来高分析
	VolumeAll       int64   // 出来高全額
	VolumeAllRatio  float32 // 絶対値と比較した割合
	VolumeBuy       int64   // 買い約定出来高
	VolumeSell      int64   // 売り約定出来高
	VolumeBuyRatio  int     // 買い約定出来高比率
	VolumeSellRatio int     // 売り約定出来高比率

	// Before
	PastPrice float32
	PastRatio float32

	// After
	NowPrice float32
	NowRatio float32

	// 直近最小比率と直近最大比率
	SmallRatio int
	BigRatio   int
	DiffRatio  int
	// 75以上の変化が起きるとtrue
	ChangeTheDiff bool

	// 判断
	TrendStrong    bool
	TrendLength    int
	TrendRatioLong string
	TrendRatio     string
	// // 全期間
	TrendAll     int
	TrendAll_b   int
	TrendAll_s   int
	TrendAllDiff int
	// // 長期トレンド
	TrendLong     int
	TrendLong_b   int
	TrendLong_s   int
	TrendLongDiff int
	// // 中期トレンド
	Trend     int
	Trend_b   int
	Trend_s   int
	TrendDiff int
	// // 短期トレンド
	TrendMin     int
	TrendMin_b   int
	TrendMin_s   int
	TrendMinDiff int
	// 超短期トレンド
	TrendMinMin     int
	TrendMinMin_b   int
	TrendMinMin_s   int
	TrendMinMinDiff int
	// // 総合判断
	TrendJudgment int
	Decision      int // BUY / SELL
	DecisionMin   int // BUY / SELL
	DecisionSMA   int // math2's

	// Demo用損益
	Demo struct {
		StartClose bool
		Side       int
		PLMax      float32
		PL         float32
		// 損益
		ProfitLoss float32
	}

	// Time
	ExchangesAt string // 取得時最終約定時間
	CreatedAt   string
	SideHistory []int
}

// 日本円の取得
func (i *Indicator) GetCurrentJPY() {
	var body TickerInfo_bf

	method := "GET"
	url := PATH_DOMAIN_BF + PATH_TICKER_BF

	// BitFlyer
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("Content-Type", "application/json")

	// Request本体
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("アクセスエラー: Bitflyer - %s", err)
		i.InfoStatus.Yen = 0
		return
	}
	defer res.Body.Close()

	bytes, _ := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(bytes, &body); err != nil {
		fmt.Printf("アンマーシャルエラー")
		i.InfoStatus.Yen = 0
		return
	}
	i.InfoStatus.Yen = float32(body.Ltp)
	return
}

// 日本円に変換
func (i Indicator) JPY(satoshi float32) float32 {
	return satoshi / BTC_SATOSHI * i.InfoStatus.Yen
}

// 約定履歴
// // 約定履歴 BitFlyer
type ApiExecution_bf struct {
	ID                        int     `json:"id"`
	Child_order_id            string  `json:"child_order_id"`
	Side                      string  `json:"side"`
	Price                     float32 `json:"price"`
	Size                      float32 `json:"size"`
	Commission                float32 `json:"commission"`
	Exec_date                 string  `json:"exec_date"`
	Child_order_acceptance_id string  `json:"child_order_acceptance_id"`
}

// 約定履歴 BitFinex
type ApiExecution_fn struct {
	ID        int    `json:"tid"`
	Side      string `json:"type"` // buy / sell
	Price     string `json:"price"`
	Size      string `json:"amount"`
	Timestamp int    `json:"timestamp"`
}

// 約定履歴 Binance
type ApiExecution_bn struct {
	ID        int    `json:"id"`
	Side      bool   `json:"isBuyerMaker"`
	Price     string `json:"price"`
	Size      string `json:"qty"`
	Timestamp int    `json:"time"`
}

type WalletSummary struct {
	Account       int64  `json:"account"`
	Currency      string `json:"currency"`
	TransactType  string `json:"transactType"`
	Symbol        string `json:"symbol"`
	Amount        int64  `json:"amount"`
	PendingDebit  int64  `json:"pendingDebit"`
	RealisedPnl   int64  `json:"realisedPnl"`
	WalletBalance int64  `json:"walletBalance"`
	UnrealisedPnl int64  `json:"unrealisedPnl"`
	MarginBalance int64  `json:"marginBalance"`
}

// 現在価格 bitFlyer
type TickerInfo_bf struct {
	ProductCode     string  `json:"product_code"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}
