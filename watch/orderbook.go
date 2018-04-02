package watch

import (
	"fmt"
	"log"

	"bitmex.pub/models"
	"github.com/qct/bitmex-go/restful"
	"github.com/qct/bitmex-go/swagger"
)

var (
	d int
)

func OrderBookL2(a *models.Indicator) {
	// 初回実行ならdebt = 0
	if d == 0 {
		d = 1
	}

	apiClient := swagger.NewAPIClient(swagger.NewConfiguration())
	// デモ口座へPATH書き換え
	if models.Demo != false {
		apiClient.ChangeBasePath("https://testnet.bitmex.com/api/v1")
	}
	orderBookApi := restful.NewOrderBookApi(apiClient.OrderBookApi)
	o, err := orderBookApi.OrderBookGetL2("XBTUSD", float32(d))
	if err != nil {
		log.Println("error wihle get orderbook: ", err)
		// 情報取得及び損益チェック
		models.API_WAIT_CHECK = 10000
		fmt.Println(models.API_WAIT_CHECK)
	} else {
		// 情報取得及び損益チェック
		models.API_WAIT_CHECK = 5000
		fmt.Println(models.API_WAIT_CHECK)
	}

	var (
		buy  int64
		sell int64
	)
	// 現価格からの差額
	if o != nil {
		l := len(o.AskList) + len(o.BidList)
		if l != 0 {
			a.NowPrice = float32((o.AskList[0].Price + o.BidList[0].Price) / 2)
		}
		d = int(a.NowPrice * a.Setting.BoardRange)

		for _, v := range o.BidList {
			// 取得価格範囲以内の出来高を取得
			// 少額指値は判断より除外する
			if 100000 < v.Size {
				buy += int64(float32(v.Price) * v.Size)
			}
		}

		for _, v := range o.AskList {
			if 100000 < v.Size {
				sell += int64(float32(v.Price) * v.Size)
			}
		}

		a.InfoStatus.PriceMax = float32(o.AskList[0].Price)
		a.InfoStatus.PriceMin = float32(o.BidList[0].Price)

		a.OrderAll = sell + buy
		a.OrderAllRatio = float32(a.OrderAll) / float32(a.Setting.OrderAb)

		// 待機注文価格
		a.OrderBuy = buy
		a.OrderSell = sell
		// 比率
		if buy != 0 && sell != 0 {
			a.OrderBuyRatio = int(buy * 100 / a.OrderAll)
			a.OrderSellRatio = int(sell * 100 / a.OrderAll)
		}
	}

}
