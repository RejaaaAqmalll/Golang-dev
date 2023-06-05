package auth

import (
	"fmt"
	"nyoba/configg"
	"nyoba/helper"
	"nyoba/models"
)

func insertCode(email string) {

	var user models.Users
	token := models.Users{Code: helper.GenerateCode()}
	err := configg.KoneksiData().Debug().Model(&user).Where("email = ?", email).Update(&token)
	if err != nil {
		panic(err)
	}

	//  Respon akhir
	fmt.Println("Your code hash been send")

}
