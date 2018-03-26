package watch

import (
	"context"
	"fmt"

	"bitmex.lib/models"
	"github.com/qct/bitmex-go/swagger"
)

func Execusions(a *models.Indicator) {
	apiClient := swagger.NewAPIClient(swagger.NewConfiguration())
	auth := context.WithValue(context.TODO(), swagger.ContextAPIKey, swagger.APIKey{
		Key:    models.API_KEY,
		Secret: models.API_SECRET,
	})

	executionApi := apiClient.ExecutionApi
	// デモ口座へPATH書き換え
	if models.Demo != false {
		apiClient.ChangeBasePath("https://testnet.bitmex.com/api/v1")
	}
	params := map[string]interface{}{
		"filter": `{}`,
		"symbol": "XBTUSD",
		"count":  float32(500),
	}
	e, res, err := executionApi.ExecutionGet(auth, params)
	if err != nil {
		fmt.Errorf("約定情報取得不可", err)
	}

	var (
		buy  int64
		sell int64
	)

	if res.StatusCode == 200 {
		for _, v := range e {
			if v.Side == "Buy" {
				buy += int64(float64(v.Price) * v.SimpleCumQty)
				fmt.Println(v.Timestamp)
			} else if v.Side == "Sell" {
				sell += int64(float64(v.Price) * v.SimpleCumQty)
			}

		}

	} else {
		a.InfoStatus.Status = res.Status
	}
}
