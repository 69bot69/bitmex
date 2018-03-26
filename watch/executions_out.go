package watch

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"

	"bitmex.lib/models"
	"github.com/labstack/gommon/log"
)

// 最終約定時間を取得するため
func GetApiExecutionsExchanges(me *models.Indicator) {
	me.Exchanges = [3]bool{}
	// 約定履歴をリクエスト
	method := "GET"
	// url := config.DOMAIN + config.PATH_EXECUTIONS + "&after=" + strconv.Itoa(me.VolumeLastID)
	url := models.PATH_DOMAIN_BF + models.PATH_EXECUTIONS
	url_fn := models.PATH_EXECUTIONS_fn
	url_bn := models.PATH_EXECUTIONS_bn

	// BitFlyer
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("Content-Type", "application/json")

	// BitFinex
	req_fn, err := http.NewRequest(method, url_fn, nil)
	req_fn.Header.Set("Content-Type", "application/json")
	// Binance
	req_bn, err := http.NewRequest(method, url_bn, nil)
	req_bn.Header.Set("Content-Type", "application/json")

	// Request本体
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Errorf("アクセスエラー: Bitflyer - %s", err)
	}
	defer res.Body.Close()

	var (
		a []models.ApiExecution_bf
		b []models.ApiExecution_fn
		c []models.ApiExecution_bn
	)
	bytes, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(bytes, &a); err != nil {
		log.Error("アンマーシャルエラー")
		me.Exchanges[0] = false
	}

	// 取引所の状態からエラーを返す
	if res.StatusCode != 200 {
		log.Errorf("%s", res.Status)
	} else {
		var (
			index                 int
			all, buy, sell, vsize float32
		)
		if len(a) != 0 {
			me.Exchanges[0] = true
			// 約定の買いと売りの強度を返す
			for _, v := range a {
				all += v.Price
				// 小さい指値を捨てる
				if 0.49 < v.Size {
					if v.Side == "BUY" {
						buy += me.NowPrice * v.Size
					}
					if v.Side == "SELL" {
						sell += me.NowPrice * v.Size
					}
				}
				// 出来高
				vsize += v.Size

				index++
			}
		}

		// Bifinex
		// API制限回避
		var res_fn *http.Response
		if me.TrendLength%5 == 0 {
			client_fn := &http.Client{}
			r, err := client_fn.Do(req_fn)
			if err != nil {
				log.Errorf("アクセスエラー: BitFinex - %s", err)
			}
			res_fn = r
			defer r.Body.Close()
			defer res_fn.Body.Close()

			if res_fn.StatusCode == 200 {
				bytes_fn, _ := ioutil.ReadAll(res_fn.Body)
				if err := json.Unmarshal(bytes_fn, &b); err != nil {
					log.Errorf("アンマーシャルエラー: BitFinexE - %s", res_fn.Status)
					me.Exchanges[1] = false
				}
			}

			if len(b) != 0 {
				for i, v := range b {
					if i == 0 {
						me.Exchanges[1] = true
					}
					var size float64
					s, _ := strconv.ParseFloat(v.Size, 32)
					if 0.49 < s {
						if v.Side == "buy" {
							size += s
							buy += me.NowPrice * float32(s)
						}
						if v.Side == "sell" {
							size += s
							sell += me.NowPrice * float32(s)
						}
					}
					// 出来高
					vsize += float32(size)
				}
			}
		} else {
			me.Exchanges[1] = false
		}

		// Binance
		client_bn := &http.Client{}
		res_bn, err := client_bn.Do(req_bn)
		if err != nil {
			log.Errorf("アクセスエラー: Binance - %s", err)
		}
		defer res_bn.Body.Close()
		if res_bn.StatusCode == 200 {
			bytes_bn, _ := ioutil.ReadAll(res_bn.Body)
			if err := json.Unmarshal(bytes_bn, &c); err != nil {
				log.Errorf("アンマーシャルエラー: BinanceE - %s", res_bn.Status)
				me.Exchanges[2] = false
			}
		}
		if len(c) != 0 {
			var (
				bb, ss float32
			)
			for i, v := range c {
				if i == 0 {
					me.Exchanges[2] = true
				}
				var size float64
				s, _ := strconv.ParseFloat(v.Size, 32)
				if 0.49 < s {
					if v.Side != false {
						size += s
						buy += me.NowPrice * float32(s)
						bb += me.NowPrice * float32(s)
					}
					if v.Side != true {
						size += s
						sell += me.NowPrice * float32(s)
						ss += me.NowPrice * float32(s)
					}
				}
				// 出来高
				vsize += float32(size)
			}
		}

		me.InfoStatus.VolatilitySize = int(vsize)

		// 標準偏差を求める
		var f1 float64

		for _, v := range a {
			f1 += math.Pow(float64(v.Price-float32(me.NowPrice)), 2)
		}
		// 分散
		me.InfoStatus.Dispersion = float32(f1) / float32(index-1)
		// 標準偏差
		me.InfoStatus.Deviation = float32(math.Sqrt(float64(me.InfoStatus.Dispersion)))

		me.VolumeLength = len(a)
		if me.VolumeLength != 0 {
			me.VolumeLastID = a[0].ID
		}

		me.VolumeAll = int64(sell + buy)
		me.VolumeAllRatio = float32(me.VolumeAll) / float32(me.Setting.VolumeAb)

		// 約定金額
		me.VolumeBuy = int64(buy)
		me.VolumeSell = int64(sell)
		// 比率
		if me.VolumeBuy != 0 && me.VolumeSell != 0 {
			me.VolumeBuyRatio = int(int64(buy) * 100 / me.VolumeAll)
			me.VolumeSellRatio = int(int64(sell) * 100 / me.VolumeAll)
		}

		// 最終取引情報を返す
		if len(a) != 0 {
			me.ExchangesAt = a[0].Exec_date
		}

		log.Infof("%s", res.Status)

	}
}
