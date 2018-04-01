package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"

	c "bitmex.lib/controllers"
	ev "bitmex.lib/evaluations"
	"bitmex.lib/models"
	"bitmex.lib/order"
	w "bitmex.lib/watch"
	"github.com/qct/bitmex-go/swagger"
)

// BitMEX API Rules
// - 5分あたり300リクエストに制限
// - UNIXタイムスタンプ
// - 1契約当たり最大200発注
// - 1契約当たり最大10ストップ注文
// - 1契約当たり最大10件の注文
// - 0.0025 XBT未満の注文は自動的に隠れ注文になり、隠れ注文は注文書には表示されず、常に入金手数料がかかる。
//
//
//

var (
	a     models.Indicator
	chart []swagger.TradeBin
)

func init() {
	// 初期設定
	a.Setting.Collateral = 400
	a.Setting.Check_PL = 0.25
	a.Setting.Profit = 10100
	a.Setting.Loss = 9915
	a.Setting.BoardRange = 0.03
	a.Setting.OrderAb = 100000000000
	a.Setting.VolumeAb = 1000000
	a.Setting.Size = 5000
	a.Setting.Leverage = 25

	// 初期証拠金取得
	// w.GetWallet(&a)
	w.GetWalletMargin(&a)

	// レバレッジ変更
	c.ChangeLeverage(&a)

	// 現建玉照会
	w.GetPositions(&a)

	// 取引所全約定履歴取得
	w.OrderBookL2(&a)

	// 4時間分の5分足を取得
	// w.GetTradeBacket1mOf4hFor5m(&a)
	// 336時間分の2時間足を取得
	chart = w.GetTradeBacket1mOf2wFor2h(&a)
	ev.CalcIndcatorsAveLine(&a, chart)

	// // テスト
	// if err := orders.PostOrderLimitSingle(a); err != nil {
	// 	fmt.Println("発注エラー出た！", err)
	// }
}

func main() {
	// 乱数作成
	rand.Seed(time.Now().UnixNano())
	// 並列ルーティン初期設定
	// // 板情報・損切りチェック
	tick_Check := time.NewTicker(time.Duration(models.API_WAIT_CHECK) * time.Millisecond)
	// // 板履歴
	tick_15s := time.NewTicker(time.Duration(15) * time.Second)
	// // 設定調整 && メモリCLEAR
	// // && ノーポジ注文チェック
	tick_60s := time.NewTicker(time.Duration(60) * time.Second)
	// // 市場時間別対応
	tick_5m := time.NewTicker(time.Duration(5) * time.Minute)

	// 並列ルーティン
	go func() {
		for {

			select {
			case <-tick_Check.C:
				// 取引所全約定履歴取得
				w.GetApiExecutionsExchanges(&a)

			case <-tick_15s.C:
				// 取引所全約定履歴取得
				w.OrderBookL2(&a)
				// インジケータ作成し、取引判断
				// ここでのインジケータを元に数回の発注処理を行う
				// 取引決断: a.decisionの発行

				// // 注意: 下記関数は.gitignoreしています。
				// // BitFlyerで使っていたEMAと複雑なロジックです。
				// ev.CalcIndcators(&a)

			case <-tick_60s.C:
				// 現建玉照会
				w.GetPositions(&a)
				// // 現アカウント状況取得
				// w.GetWallet(&a)
				w.GetWalletMargin(&a)
				// レバレッジ変更
				c.ChangeLeverage(&a)

			case <-tick_5m.C:
				// // 4時間分の5分足を取得
				// w.GetTradeBacket1mOf4hFor5m(&a)
				// 336時間分の2時間足を取得
				chart = w.GetTradeBacket1mOf2wFor2h(&a)
				ev.CalcIndcatorsAveLine(&a, chart)

			}
		}
	}()

	// 上記までをTicker → go func for Selectで並列的に取得及び監視
	//
	//
	//
	//
	//
	// 下記はAPI_WAITで待機以外は常に一定にFor Loop
	// 情報はリアルタイムで変化する(参照渡しではなくポインタ渡し)

	go func() {
		for {
			flow(&a)

			// 発注処理の間隔
			time.Sleep(time.Duration(models.API_WAIT_ORDER) * time.Millisecond)
		}
	}()

	// サーバー起動
	// 顧客への表示出力
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(
			w,
			html,
			a.Setting.Collateral,
			int(a.Setting.Check_PL*100),
			a.Setting.Profit, a.Setting.Loss,
			int(a.Setting.BoardRange*100),
			a.Setting.Size,

			a.InfoStatus.StartCol, a.InfoStatus.NowCol,
			a.InfoStatus.AvailableMargin,
			a.InfoStatus.ChangeCol,

			int(a.InfoStatus.Leverage),
			a.InfoStatus.PositionsSymbol,
			a.InfoStatus.PositionsSide,
			a.InfoStatus.PositionsSize,
			a.InfoStatus.PositionsAvg,
			int(a.InfoStatus.PositionsCost),
			float32(a.InfoStatus.PositionsCostPcnt*100),
			int(a.InfoStatus.PositionsPL),
			float32(a.InfoStatus.PositionsPLPcnt*100),

			a.InfoStatus.Health,
			a.DecisionSMA,
			a.InfoStatus.Status,

			a.NowPrice,
			a.PastRatio, a.NowRatio,
			a.NowRatio-a.PastRatio,
			a.InfoStatus.PriceMax, a.InfoStatus.PriceMin,
			a.InfoStatus.PriceMax4h, a.InfoStatus.PriceMin4h,
			a.InfoStatus.PriceDiff4h,

			a.OrderAll, a.OrderAllRatio,
			a.OrderBuyRatio, a.OrderSellRatio,

			a.VolumeAll, a.VolumeAllRatio,
			a.VolumeBuyRatio, a.VolumeSellRatio,
			a.BigRatio, a.SmallRatio,
			a.DiffRatio,

			a.Decision,
			a.TrendLength,
			a.TrendAllDiff, a.TrendLongDiff, a.TrendDiff, a.TrendMinDiff, a.TrendMinMinDiff,
			a.TrendAll, a.TrendLong, a.Trend, a.TrendMin, a.TrendMinMin,
		)
	})
	http.ListenAndServe(models.PORT, nil)
}

// 発注処理を行う
func flow(a *models.Indicator) {
	// チェック: 証拠金状態
	// チェック: 取引所状態（遅延・Status）

	// if a.InfoStatus.PositionsSize < models.TRADE_MAX {
	// 	if err := orders.PostOrderLimitSingle(*a); err != nil {
	// 		fmt.Println("発注エラー出た！", err)
	// 	}
	// } else {
	// 	a.InfoStatus.Status = fmt.Sprintf("否決 - 最大取引枚数を超え、%d枚保持なので発注見送り", a.InfoStatus.PositionsSize)
	// }

	if a.DecisionSMA != 0 {
		// 残建玉成行き決済
		// 保有建玉があれば
		if a.InfoStatus.PositionsSize != 0 {
			if err := orders.PostOrderSMA_Done(*a); err != nil {
				fmt.Println("決済注文発注エラー: ", err)
			}
		}

		// 新規注文
		if err := orders.PostOrderSMA_Market(*a); err != nil {
			fmt.Println("新規注文発注エラー: ", err)
		}

		// 取引決断をリセット
		a.Decision = 0
		a.DecisionMin = 0
		a.DecisionSMA = 0

	} else {
		a.InfoStatus.Status = fmt.Sprint("否決 - SMAクロストレンド検出なし")

		// 利益計算後、決済チェック
		if a.InfoStatus.PositionsSize != 0 {
			var (
				// ポジション取得時から1割
				p1 float32 = a.InfoStatus.PositionsAvg * 0.1
				// ポジション取得時と現在価格の差額
				p2 float32 = float32(math.Abs(float64(a.NowPrice - a.InfoStatus.PositionsAvg)))
			)

			// 差額が1割超えなら決済
			if p1 < p2 {
				if err := orders.PostOrderSMA_Done(*a); err != nil {
					fmt.Println("決済注文発注エラー: ", err)
				}
			} else {
				fmt.Println("建玉利益チェック - 設定利益未満")
			}
		} else {
			fmt.Println("建玉利益チェック - 建玉0")
		}
	}
}

//
//
//
//
//
//
//
//
//
//
//
// 顧客表記用HTML
var html string = `<!DOCTYPE html>
<html lang="ja">
<head>
  <meta charset="UTF-8" />
  <meta data-n-head="true" name="robots" content="noindex, nofollow"/>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>自動取引ボットプログラム</title>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.6.2/css/bulma.min.css">
</head>
<body>
  <section class="section">
    <div class="container">
      <div class="content">
        <dl>
          <!--
          <dt>動的セッティング</dt>
          <dd>証拠金活用割合: %d％</dd>
          <dd>強制損益確定割合: %d％</dd>
          <dd>
            利確/損切目安幅: %d / %d
          </dd>
          <dd>板監視範囲: %d％</dd>
          <dd>取引枚数: %d枚</dd>
           -->

          <dt>ステータス</dt>
          <dd>初期/現在証拠金: %d / %d</dd>
          <dd>使用可能証拠金: %d</dd>
          <dd>証拠金変動額: %d</dd>

          <dt>建玉状況</dt>
          <dd>レバレッジ: %d</dd>
          <dd>建玉種類: %s</dd>
          <dd>売買方向: %s</dd>
          <dd>建玉枚数: %d</dd>
          <dd>平均取得価額: %f</dd>
          <dd>使用証拠金: %d XBt</dd>
          <dd>証拠金使用率: %f ％</dd>
          <dd>総建玉損益合計: %d XBt</dd>
          <dd>総資産ROE: %f ％</dd>

          <dt>現情報</dt>
          <dd>取引所状態: %s</dd>

          <dt>アルゴリズム</dt>
          <dd><strong>SMA9/16クロスシグナル: %d</strong></dd>
          <dd>アルゴリズムは「<strong>%s</strong>」と判断しています。</dd>

          <dt>チャート関連</dt>
          <dd>現在値: %f</dd>

          <!--
          <dd>現在買売強弱指数: %f → <strong>%f</strong></dd>
          <dd>買売強弱指数差: <strong>%f</strong></dd>
          <dd>直近最大/最小値: %f / %f</dd>


          <dd>4h最大/最小値: %f / %f</dd>
          <dd>4h変化額: %f</dd>
          <br>
          <dd>待機注文/基本値比較: %d / %f</dd>
          <dd>待機注文買売割合: %d:%d</dd>


          <dd>瞬間出来高/基本値比較: %d / %f</dd>
          <dd>約定出来高買売比率: %d:%d</dd>
          <dd>板比率: 最大%d / 最小%d</dd>
          <dd>板比率最大変動: %d</dd>


          <dt>アルゴリズム取引判断: %d</dt>
          <dd>判断要素数: %d</dd>
          <dd>トレンド要素: 全期間 %d, 長期 %d, 中期 %d, 短期 %d, 超短期 %d</dd>
          <dd>トレンド: 全期間 %d, 長期 %d, 中期 %d, 短期 %d, 超短期 %d</dd>
         -->
        </dl>
      </div>
    </div>
  </section>
</body>
</html>
`
