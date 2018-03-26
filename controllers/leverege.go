package controllers

import (
	"context"
	"fmt"
	"log"

	"bitmex.lib/models"
	"github.com/qct/bitmex-go/swagger"
)

func ChangeLeverage(a *models.Indicator) {
	symbol := "XBTUSD"
	leverage := a.Setting.Leverage
	apiClient := swagger.NewAPIClient(swagger.NewConfiguration())
	// デモ口座へPATH書き換え
	if models.Demo != false {
		apiClient.ChangeBasePath("https://testnet.bitmex.com/api/v1")
	}
	auth := context.WithValue(context.TODO(), swagger.ContextAPIKey, swagger.APIKey{
		Key:    models.API_KEY,
		Secret: models.API_SECRET,
	})

	positionsApi := apiClient.PositionApi
	l, res, err := positionsApi.PositionUpdateLeverage(auth, symbol, leverage)
	if err != nil {
		log.Println("error: ", err)
	}

	if res.StatusCode != 200 {
		fmt.Printf("%s エラー", res.Status)
	} else {
		fmt.Printf("レバレッジ設定: %f\n", l.Leverage)
		a.InfoStatus.Leverage = float32(l.Leverage)
	}
}
