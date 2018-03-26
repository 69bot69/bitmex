package writer

import (
	"bitmex.lib/models"
)

var (
	// インジケータ用保存データ配列
	DATA []models.Indicator
)

func ChangePointer(a models.Indicator) {
	DATA = append(DATA, a)
	if len(DATA) != 0 && len(DATA)%1440 == 0 {
		// トレンド素材情報を消す
		for i, _ := range DATA {
			DATA[i].SideHistory = []int{}
		}

		WriteJson(DATA)
		DATA = []models.Indicator{}
	}
}
