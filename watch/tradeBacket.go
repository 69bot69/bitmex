package watch

import (
	"fmt"
	"time"

	"bitmex.lib/models"
	"github.com/labstack/gommon/log"
	"github.com/qct/bitmex-go/swagger"
)

func GetTradeBacket1mOf4hFor5m(a *models.Indicator) {
	apiClient := swagger.NewAPIClient(swagger.NewConfiguration())
	// デモ口座へPATH書き換え
	if models.Demo != false {
		apiClient.ChangeBasePath("https://testnet.bitmex.com/api/v1")
	}
	tradeApi := apiClient.TradeApi
	now := time.Now().UTC()
	t4 := now.Add(time.Duration(-4) * time.Hour)
	params := map[string]interface{}{
		"binSize":   "5m",
		"symbol":    "XBTUSD",
		"partial":   false,
		"reverse":   false,
		"count":     float32(48),
		"startTime": t4,
	}

	t, res, err := tradeApi.TradeGetBucketed(params)
	if err != nil {
		fmt.Println(res.Status, err)
	}

	if res.StatusCode != 200 {
		fmt.Errorf("%s", res.Status)
	} else {
		// 成功処理
		// // ４時間足最大最小値を取る
		a.InfoStatus.PriceMax4h = 0
		a.InfoStatus.PriceMin4h = a.NowPrice
		var (
			// 出現順番
			apMax int
			apMin int
		)
		for i, v := range t {
			if a.InfoStatus.PriceMax4h < float32(v.High) {
				a.InfoStatus.PriceMax4h = float32(v.High)
				apMax = i
			}
			if float32(v.Low) < a.InfoStatus.PriceMin4h {
				a.InfoStatus.PriceMin4h = float32(v.Low)
				apMin = i
			}
		}

		log.Infof(res.Status)

		if apMin < apMax {
			a.InfoStatus.PriceDiff4h = a.InfoStatus.PriceMax4h - a.InfoStatus.PriceMin4h
		} else if apMax < apMin {
			a.InfoStatus.PriceDiff4h = a.InfoStatus.PriceMin4h - a.InfoStatus.PriceMax4h
		}
	}
}

func GetTradeBacket1mOf2wFor2h(a *models.Indicator) []swagger.TradeBin {
	apiClient := swagger.NewAPIClient(swagger.NewConfiguration())
	// デモ口座へPATH書き換え
	if models.Demo != false {
		apiClient.ChangeBasePath("https://testnet.bitmex.com/api/v1")
	}
	tradeApi := apiClient.TradeApi
	now := time.Now().UTC()
	t4 := now.Add(time.Duration(-336) * time.Hour)
	params := map[string]interface{}{
		"binSize":   "1h",
		"symbol":    "XBTUSD",
		"partial":   false,
		"reverse":   false,
		"count":     float32(500),
		"startTime": t4,
	}

	t, res, err := tradeApi.TradeGetBucketed(params)
	if err != nil {
		fmt.Println(res.Status, err)
	}

	if res.StatusCode != 200 {
		fmt.Errorf("%s", res.Status)
	} else {
		// 成功処理
		log.Infof(res.Status)
		//
		// fmt.Printf("%+v\n", t)

		return t
	}

	return nil
}
