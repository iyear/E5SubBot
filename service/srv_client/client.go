package srv_client

import (
	"github.com/iyear/E5SubBot/db"
	"github.com/iyear/E5SubBot/model"
)

func Add(c *model.Client) error {
	return db.DB.Create(c).Error
}

func Update(c *model.Client) error {
	return db.DB.Save(c).Error
}

func Del(id int) error {
	return db.DB.Where("id = ?", id).Delete(&model.Client{}).Error
}

func GetAllClients() []*model.Client {
	var clients []*model.Client
	db.DB.Find(&clients)
	return clients
}

func GetClients(uid int64) []*model.Client {
	var clients []*model.Client
	db.DB.Where("tg_id = ?", uid).Find(&clients)
	return clients
}

func GetClient(id int) (*model.Client, error) {
	var client model.Client
	err := db.DB.Where("id = ?", id).First(&client).Error
	return &client, err
}

func IsExist(tgID int64, clientID string) bool {
	return !(db.DB.
		Where("tg_id = ? AND client_id = ?", tgID, clientID).
		First(&model.Client{}).RowsAffected == 0)
}
