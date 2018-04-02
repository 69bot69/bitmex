package orders

import (
	"context"
	"fmt"
	"math"

	"bitmex.pub/models"
	"github.com/qct/bitmex-go/swagger"
)

func round(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor(f*shift+.5) / shift
}

func PostOrderLimitSingle(a models.Indicator) error {
	apiClient := swagger.NewAPIClient(swagger.NewConfiguration())
	// デモ口座へPATH書き換え
	if models.Demo != false {
		apiClient.ChangeBasePath("https://testnet.bitmex.com/api/v1")
	}
	auth := context.WithValue(context.TODO(), swagger.ContextAPIKey, swagger.APIKey{
		Key:    models.API_KEY,
		Secret: models.API_SECRET,
	})

	oApi := apiClient.OrderApi

	var (
		side string
		size float32
	)

	// 売買方向記述
	if a.Decision == 1 || a.DecisionMin == 1 {
		side = "Buy"
	} else if a.Decision == -1 || a.DecisionMin == -1 {
		side = "Sell"
	}
	if a.Decision == 0 {
		// ポジションがあり、かつ、短期売買ならば
		// 2倍で反対売買取引
		if a.DecisionMin != 0 {
			if a.InfoStatus.PositionsSize != 0 && a.InfoStatus.PositionsSide != side {
				size = float32(a.InfoStatus.PositionsSize)
			} else {
				size = float32(a.Setting.Size)
			}
		}
	} else {
		// a.Decisionがある場合は、大きく反対売買

		// 1. 全てのポジションを整理して
		if a.InfoStatus.PositionsSide != side {
			size = float32(math.Abs(float64(a.InfoStatus.PositionsSize)) + 10000)
		} else {
			size = float32(10000 - math.Abs(float64(a.InfoStatus.PositionsSize)))
		}
		// 2. 証拠金最適枚数で取引開始
		// 3. 一定以上損益が発生すれば、適宜指値発注
		// 4.
	}

	params := map[string]interface{}{
		"symbol": "XBTUSD",
		"side":   side,
		// 枚数
		//"simpleOrderQty": ,
		"orderQty": size,
		"price":    round(float64(a.NowPrice), 0),
		"ordType":  "Limit",
	}
	o, res, err := oApi.OrderNew(auth, "XBTUSD", params)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf(res.Status)
	} else {
		fmt.Println(o.OrderID)
		return nil
	}

	return fmt.Errorf(res.Status)
}

// DecicionSMA用発注
func PostOrderSMA_Market(a models.Indicator) error {
	apiClient := swagger.NewAPIClient(swagger.NewConfiguration())
	// デモ口座へPATH書き換え
	if models.Demo != false {
		apiClient.ChangeBasePath("https://testnet.bitmex.com/api/v1")
	}
	auth := context.WithValue(context.TODO(), swagger.ContextAPIKey, swagger.APIKey{
		Key:    models.API_KEY,
		Secret: models.API_SECRET,
	})

	oApi := apiClient.OrderApi

	var (
		side string
		size float64
	)

	// // 活用可能な金額を算出
	// size = float32(a.InfoStatus.NowCol / 20000000)
	// テスト用サイズ
	size = float64(a.InfoStatus.NowCol) / 20000000.0
	fmt.Println(a.InfoStatus.NowCol, size)
	// 売買方向記述
	if a.DecisionSMA == 1 {
		side = "Buy"
	} else if a.DecisionSMA == -1 {
		side = "Sell"
	}

	params := map[string]interface{}{
		"symbol": "XBTUSD",
		"side":   side,
		// Btcの枚数 float64
		"simpleOrderQty": size,
		// XBtの枚数
		// "orderQty": size,
		// "price":   round(float64(a.NowPrice), 0),
		"ordType": "Market",
	}
	o, res, err := oApi.OrderNew(auth, "XBTUSD", params)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf(res.Status)
	} else {
		fmt.Println(o.OrderID)
		return nil
	}

	return fmt.Errorf(res.Status)
}

// トレンド検出後
// ポジションがある場合
func PostOrderSMA_Done(a models.Indicator) error {
	apiClient := swagger.NewAPIClient(swagger.NewConfiguration())
	// デモ口座へPATH書き換え
	if models.Demo != false {
		apiClient.ChangeBasePath("https://testnet.bitmex.com/api/v1")
	}
	auth := context.WithValue(context.TODO(), swagger.ContextAPIKey, swagger.APIKey{
		Key:    models.API_KEY,
		Secret: models.API_SECRET,
	})

	oApi := apiClient.OrderApi

	var (
		side string
		size float32
	)

	// 活用可能な金額を算出
	size = float32(a.InfoStatus.PositionsSize)
	// 売買方向記述
	if a.InfoStatus.PositionsSide == "Buy" {
		side = "Sell"
	} else if a.InfoStatus.PositionsSide == "Sell" {
		side = "Buy"
	}

	params := map[string]interface{}{
		"symbol": "XBTUSD",
		"side":   side,
		// // Btcの枚数
		// "simpleOrderQty": size,
		// XBtの枚数 float32
		"orderQty": size,
		// "price":    round(float64(a.NowPrice), 0),
		"ordType": "Market",
	}
	o, res, err := oApi.OrderNew(auth, "XBTUSD", params)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf(res.Status)
	} else {
		fmt.Println(o.OrderID)
		return nil
	}

	return fmt.Errorf(res.Status)
}
