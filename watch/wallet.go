package watch

import (
	"context"

	"bitmex.pub/models"
	"github.com/labstack/gommon/log"
	"github.com/qct/bitmex-go/swagger"
)

func GetWallet(a *models.Indicator) {
	apiClient := swagger.NewAPIClient(swagger.NewConfiguration())
	// デモ口座へPATH書き換え
	if models.Demo != false {
		apiClient.ChangeBasePath("https://testnet.bitmex.com/api/v1")
	}
	auth := context.WithValue(context.TODO(), swagger.ContextAPIKey, swagger.APIKey{
		Key:    models.API_KEY,
		Secret: models.API_SECRET,
	})

	userApi := apiClient.UserApi
	params := map[string]interface{}{
		"currency": "XBt",
	}
	w, res, err := userApi.UserGetWalletHistory(auth, params)
	if err != nil {
		log.Errorf("error: %s", err)
	}

	if res.StatusCode != 200 {
		log.Errorf("Get wallet summary: %s", res.Status)
	} else {
		if len(w) != 0 {
			// ウォレット情報取得
			// // 現在証拠金
			if a.InfoStatus.StartCol == 0 {
				a.InfoStatus.StartCol = int(w[len(w)-1].Amount)
			}

			var change float32
			for _, v := range w {
				change += v.Amount
			}
			// change
			a.InfoStatus.NowCol = int(w[len(w)-1].Amount) + int(change)

			// // 証拠金変化額 a.InfoStatus.NowCol - a.InfoStatus.StartCol
			a.InfoStatus.ChangeCol = int(change)
		}
	}
}

func GetWalletMargin(a *models.Indicator) {
	apiClient := swagger.NewAPIClient(swagger.NewConfiguration())
	// デモ口座へPATH書き換え
	if models.Demo != false {
		apiClient.ChangeBasePath("https://testnet.bitmex.com/api/v1")
	}
	auth := context.WithValue(context.TODO(), swagger.ContextAPIKey, swagger.APIKey{
		Key:    models.API_KEY,
		Secret: models.API_SECRET,
	})

	userApi := apiClient.UserApi
	params := map[string]interface{}{
		"currency": "XBt",
	}
	w, res, err := userApi.UserGetMargin(auth, params)
	if err != nil {
		log.Errorf("error: %s", err)
	}

	if res.StatusCode != 200 {
		log.Errorf("Get wallet summary: %s", res.Status)
	} else {
		// ウォレット情報取得
		// // 現在証拠金
		if a.InfoStatus.StartCol == 0 {
			a.InfoStatus.StartCol = int(w.Amount)
		}
		// change
		a.InfoStatus.NowCol = int(w.WalletBalance)
		a.InfoStatus.PositionsCostPcnt = w.MarginUsedPcnt
		a.InfoStatus.AvailableMargin = int(w.AvailableMargin)
	}
}
