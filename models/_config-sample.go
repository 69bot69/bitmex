package models

var (
	// 情報取得及び損益チェック
	API_WAIT_CHECK = 5000
	// 発注処理の間隔
	API_WAIT_ORDER = 30000

	// 発注最大枚数
	TRADE_MAX = 20000
)

// 初期設定時は下記項目を全て埋めてください
// その後、当ファイル名先頭の「アンダーバー（_）」を消してプログラムを有効にしてください
const (
	// トレンド転換時送信用
	Gmail       = ""
	GmailPass   = ""
	GmailSendTo = ""

	// BITMEX API Keys
	Name       = "test"
	Demo       = false
	PORT       = ":8080"
	API_KEY    = ""
	API_SECRET = ""

	// // BITMEX API Keys for Demo
	// Name       = "demo"
	// Demo       = true
	// PORT       = ":8080"
	// API_KEY    = ""
	// API_SECRET = ""

	// チャート形状補正（単位: 1000分の1
	REVISION     = 34
	REVISION_ADD = 1.05

	// 約定履歴（各取引所
	// // BitFlyer
	PATH_DOMAIN_BF  = "https://api.bitflyer.jp"
	PATH_EXECUTIONS = "/v1/getexecutions?product_code=FX_BTC_JPY&count=500"
	PATH_TICKER_BF  = "/v1/ticker?product_code=BTC_JPY"
	// // BitFinex
	PATH_EXECUTIONS_fn = "https://api.bitfinex.com/v1/trades/btcusd?limit_trades=500"
	// // Binance
	PATH_EXECUTIONS_bn = "https://api.binance.com/api/v1/trades?symbol=BTCUSDT&limit=500"
)
