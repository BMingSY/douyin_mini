package initial

import "tiktok/utils"

func InitTokenBucket() {

	utils.TB = &utils.TokenBucket{}
	utils.TB.Set(500, 3000)
}
