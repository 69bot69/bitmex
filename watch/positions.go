package watch

import (
	"context"
	"fmt"

	"bitmex.lib/models"
	"github.com/qct/bitmex-go/swagger"
)

func GetPositions(a *models.Indicator) {
	apiClient := swagger.NewAPIClient(swagger.NewConfiguration())
	// デモ口座へPATH書き換え
	if models.Demo != false {
		apiClient.ChangeBasePath("https://testnet.bitmex.com/api/v1")
	}
	auth := context.WithValue(context.TODO(), swagger.ContextAPIKey, swagger.APIKey{
		Key:    models.API_KEY,
		Secret: models.API_SECRET,
	})

	positionApi := apiClient.PositionApi
	params := map[string]interface{}{
		"filter":  `{"symbol": "XBTUSD"}`,
		"columns": "",
		"count":   float32(10),
	}
	p, res, err := positionApi.PositionGet(auth, params)
	if err != nil {
		fmt.Println(err)
	}

	if res.StatusCode != 200 {
		// return fmt.Errorf(res.Status)
	}

	// // todo: 未完
	// for i,v := range p {
	// 	if i == 0 {
	// 		// 建玉種類
	// 		a.InfoStatus.PosisionsSymbol = v.Symbol
	// 		// 建玉の売買方向
	//
	// 	}
	//
	//
	//
	// 	a.InfoStatus.PosisionsAvg += int(v.)
	// 	a.InfoStatus.PosisionsAvg += int(v.AvgEntryPrice)
	//
	// }

	if len(p) != 0 {
		a.InfoStatus.PositionsSymbol = p[0].Symbol
		a.InfoStatus.PositionsAvg = float32(p[0].AvgEntryPrice)
		var (
			size float32
			// size_b float32
			// size_s float32
			cost     float32
			pl       float32
			roe_pl_i int
			roe_pl   float64
		)

		for _, v := range p {
			size += v.CurrentQty
			// size_b += v.ExecBuyQty
			// size_s -= v.ExecSellQty
			// cost += v.ExecBuyCost
			// cost += v.ExecSellCost

			cost += v.MaintMargin
			pl += v.UnrealisedGrossPnl
			roe_pl_i++
			roe_pl += v.UnrealisedPnlPcnt
			fmt.Printf("%f - %f", v.ExecBuyQty, v.ExecBuyCost)
		}

		// 売買建玉方向
		// size = size_b - size_s
		if 0 < size {
			a.InfoStatus.PositionsSide = "Buy"
		} else if size < 0 {
			a.InfoStatus.PositionsSide = "Sell"
		} else {
			a.InfoStatus.PositionsSide = ""
		}

		a.InfoStatus.PositionsSize = int(size)
		// 証拠金背景（取得コスト）
		a.InfoStatus.PositionsCost = float32(cost)
		a.InfoStatus.PositionsPL = float32(pl)
		a.InfoStatus.PositionsPLPcnt = roe_pl / float64(roe_pl_i)
	}
}
