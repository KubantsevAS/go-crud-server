package stat

import (
	"demo/go-server/pkg/db"
	"time"

	"gorm.io/datatypes"
)

type StatRepository struct {
	DataBase *db.Db
}

func NewStatRepository(database *db.Db) *StatRepository {
	return &StatRepository{
		DataBase: database,
	}
}

func (repo *StatRepository) AddClick(linkId uint) {
	var stat Stat
	currentDate := datatypes.Date(time.Now())
	repo.DataBase.DB.Find(&stat, "link_id = ? and date = ?", linkId, currentDate)

	if stat.ID == 0 {
		newStat := Stat{
			LinkId: linkId,
			Clicks: 1,
			Date:   currentDate,
		}
		repo.DataBase.DB.Create(&newStat)
	} else {
		stat.Clicks += 1
		repo.DataBase.DB.Save(&stat)
	}
}
