package bots

import (
	"errors"
	"fmt"
	"github.com/iyear/E5SubBot/model"
	"github.com/iyear/E5SubBot/util"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
)

// BindUser If Successfully return "",else return error information
func BindUser(m *tb.Message, ClientId, ClientSecret string) error {
	tmp := strings.Split(m.Text, " ")
	if len(tmp) != 2 {
		return errors.New("wrong format")
	}
	code := util.GetURLValue(tmp[0], "code")
	Alias := tmp[1]
	cli := model.NewClient(ClientId, ClientSecret)
	if err := cli.GetTokenWithCode(code); err != nil {
		return err
	}
	bot.Send(m.Chat, "Token获取成功!")

	info, err := cli.GetUserInfo()
	if err != nil {
		return err
	}
	var u = &model.Client{
		TgId: m.Chat.ID,
		//TG的Data传递最高64bytes,一些MsId超过了报错BUTTON_DATA_INVALID (0)，采取md5
		RefreshToken: cli.RefreshToken,
		MsId:         util.Get16MD5Encode(gjson.Get(info, "id").String()),
		Alias:        Alias,
		ClientId:     ClientId,
		ClientSecret: ClientSecret,
		Other:        "",
	}

	//MS User Is Exist
	if MSAppIsExist(u.TgId, u.ClientId) {
		return errors.New("该应用已经绑定过了，无需重复绑定")
	}
	//MS information has gotten
	bot.Send(m.Chat,
		fmt.Sprintf("MS_ID(MD5)： %s\nuserPrincipalName： %s\ndisplayName： %s\n",
			u.MsId,
			gjson.Get(info, "userPrincipalName").String(),
			gjson.Get(info, "displayName").String()),
	)

	if result := model.DB.Create(&u); result.Error != nil {
		return result.Error
	}
	return nil
}

// GetBindNum get bind num
func GetBindNum(TgId int64) int {
	var bindings []*model.Client
	result := model.DB.Where("tg_id = ?", TgId).Find(&bindings)
	return int(result.RowsAffected)
}

// MSAppIsExist return true => exist
func MSAppIsExist(TgId int64, ClientId string) bool {
	result := model.DB.
		Where("tg_id = ? AND client_id = ?", TgId, ClientId).
		First(&model.Client{})
	return util.IF(result.RowsAffected == 0, false, true).(bool)
}
func GetAdmins() []int64 {
	var result []int64
	admins := strings.Split(viper.GetString("admin"), ",")
	for _, v := range admins {
		id, _ := strconv.ParseInt(v, 10, 64)
		result = append(result, id)
	}
	return result
}
