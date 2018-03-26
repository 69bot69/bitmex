package evaluations

import (
	"fmt"
	"log"
	"time"

	"bitmex.lib/models"
	"github.com/qct/bitmex-go/swagger"
	"togashi.dev/util"
)

var (
	change string
	now    string
	past   string
	allPL  float64
)

func CalcIndcatorsAveLine(a *models.Indicator, chart []swagger.TradeBin) {
	start := time.Now()

	var (
		// 2h足をつくる
		chart2h []models.LineSingle
		// 移動平均線
		ma9  []models.LineSingle
		ma16 []models.LineSingle

		// 平均足を取るための過去足配列
		tc []models.LineSingle
		// 合計値
		sum float64
		// 平均値
		av float64

		// 直近検出
	)

	for _, v := range chart {
		if v.Timestamp.Hour()%2 == 0 {
			var this models.LineSingle
			this.Av = v.Close

			jst := time.FixedZone("Asia/Tokyo", 9*60*60)

			nowJST := v.Timestamp.In(jst)
			t := nowJST.Add(time.Duration(-4) * time.Hour)

			this.Timestamp = t
			chart2h = append(chart2h, this)
		}
	}

	for i, v := range chart2h {
		if 15 < i {
			// 9*2h足平均線を作る
			tc = chart2h[i-9 : i]
			sum = 0
			for _, val := range tc {
				sum += val.Av
			}
			av = sum / float64(len(tc))
			var this models.LineSingle
			this.Av = av
			this.Timestamp = v.Timestamp
			ma9 = append(ma9, this)

			// 16*2h足平均線を作る
			tc = chart2h[i-16 : i]
			sum = 0
			for _, val := range tc {
				sum += val.Av
			}
			av = sum / float64(len(tc))
			var _this models.LineSingle
			_this.Av = av
			_this.Timestamp = v.Timestamp
			ma16 = append(ma16, _this)
		}
	}

	// クロス判断
	for i, v := range ma16 {
		// ゴールデン・デッドクロス
		if v.Av < ma9[i].Av {
			now = "買い"
		} else if ma9[i].Av < v.Av {
			now = "売り"
		}

		// 転換点による取引方向判断
		if now != past {
			fmt.Printf("CROSS & CHANGE!! %s\n", v.Timestamp)
			change = fmt.Sprintf("CROSS & CHANGE!! %s", v.Timestamp)
			if now == "買い" {
				a.DecisionSMA = 1
			} else if now == "売り" {
				a.DecisionSMA = -1
			}
		} else {
			change = ""
			a.DecisionSMA = 0
		}

		a.InfoStatus.Status = fmt.Sprintf("%s %s", now, change)

		past = now
	}

	// 転換点を経て通して悪い時
	// - ポジションを持っているが逆ポジ転換点ではな時
	// 現在ポジション方向性
	var which int
	if 0 < a.InfoStatus.PositionsSize {
		which = 1
	} else if a.InfoStatus.PositionsSize < 0 {
		which = -1
	}

	// 建玉有り
	// 建玉ないなら、Decisionを通す
	if which != 0 {
		// かつ、建玉とトレンドが合致したら、変更なし
		if a.DecisionSMA == which {
			a.DecisionSMA = 0
		}
	}

	// トレンド転換報告
	if a.DecisionSMA != 0 {
		var email util.SendMailForError
		m := util.Mail{
			// 通知送信用メアドから送信する
			From:     models.Gmail,
			Username: models.Gmail,
			Password: models.GmailPass,
			To:       models.GmailSendTo,
			Sub:      fmt.Sprintf("トレンド転換! %s", models.Name),
			Msg:      email,
		}

		// gmailからメールを送っちゃうよ
		if err := util.GmailSend(m); err != nil {
			log.Fatalf("メール送信失敗: %s", err)
		}
	}

	// // 保存用データ
	// writer.ChangePointer(*a)
	t := time.Now()
	fmt.Println("Calc2実行時間: ", t.Sub(start))
}
