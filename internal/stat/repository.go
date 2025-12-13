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

func (repo *StatRepository) GetAll(by string, from, to time.Time) []GetStatResponse {
	var stats []GetStatResponse
	var selectQuery string

	switch by {
	case GroupByMonth:
		selectQuery = "to_char(date, 'YYYY-MM') as period, sum(clicks)"
	default:
		selectQuery = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	}

	repo.DataBase.DB.Table("stats").
		Select(selectQuery).
		Where("date BETWEEN ? AND ?", from, to).
		Group("period").
		Order("period").
		Scan(&stats)

	return stats
}
