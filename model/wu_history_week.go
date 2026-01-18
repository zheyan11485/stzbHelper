package model

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

// WuHistoryWeek 存储每周武勋历史数据
type WuHistoryWeek struct {
	Id          int       `json:"id" gorm:"column:id;primaryKey"`
	GroupId     string    `json:"group_id" gorm:"column:group_id"`           // 分组ID
	GroupName   string    `json:"group_name" gorm:"column:group_name"`       // 分组名称
	MemberCount int       `json:"member_count" gorm:"column:member_count"`   // 成员数量
	TotalWu     int       `json:"total_wu" gorm:"column:total_wu"`           // 总武勋
	AverageWu   int       `json:"average_wu" gorm:"column:average_wu"`       // 平均武勋
	ZeroWuCount int       `json:"zero_wu_count" gorm:"column:zero_wu_count"` // 0武勋人数
	RecordDate  time.Time `json:"record_date" gorm:"column:record_date"`     // 记录日期
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`       // 创建时间
}

func (WuHistoryWeek) TableName() string {
	return "wu_history_week"
}

// SaveGroupWuHistory 保存分组武勋历史数据
func SaveGroupWuHistory() error {
	// 查询当前各分组的统计数据
	type GroupWuStats struct {
		Group       string `json:"group"`
		MemberCount int    `json:"member_count"`
		TotalWu     int    `json:"total_wu"`
		AverageWu   int    `json:"average_wu"`
		ZeroWuCount int    `json:"zero_wu_count"`
	}

	var stats []GroupWuStats

	subQuery := Conn.Model(&TeamUser{}).
		Select("`group`, COUNT(*) as zero_wu_count").
		Where("wu = 0").
		Group("`group`")

	err := Conn.Model(&TeamUser{}).
		Select("`team_user`.`group`, COUNT(*) as member_count, SUM(wu) as total_wu, ROUND(AVG(wu)) as average_wu, IFNULL(sub.zero_wu_count, 0) as zero_wu_count").
		Joins("LEFT JOIN (?) as sub ON sub.`group` = `team_user`.`group`", subQuery).
		Group("`team_user`.`group`").
		Scan(&stats).Error

	if err != nil {
		return err
	}

	currentTime := time.Now()
	// 只保留日期部分，确保时间部分为00:00:00
	currentDate := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	currentDateStr := currentDate.Format("2006-01-02")

	// 使用事务确保数据一致性
	return Conn.Transaction(func(tx *gorm.DB) error {
		for _, stat := range stats {
			var existingHistory WuHistoryWeek
			// 使用更明确的日期比较，确保只匹配同一天的记录
			result := tx.Where("strftime('%Y-%m-%d', record_date) = ? AND group_name = ?", currentDateStr, stat.Group).First(&existingHistory)

			if result.Error != nil {
				// 检查是否是因为记录未找到
				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					// 如果没有找到当天该分组的记录，则创建新记录
					history := WuHistoryWeek{
						GroupName:   stat.Group,
						GroupId:     stat.Group,
						MemberCount: stat.MemberCount,
						TotalWu:     stat.TotalWu,
						AverageWu:   stat.AverageWu,
						ZeroWuCount: stat.ZeroWuCount,
						RecordDate:  currentDate, // 使用规范化后的日期
						CreatedAt:   currentTime,
					}
					if err := tx.Create(&history).Error; err != nil {
						return err
					}
				} else {
					// 其他数据库错误
					return result.Error
				}
			} else {
				// 如果找到当天该分组的记录，则更新它
				existingHistory.MemberCount = stat.MemberCount
				existingHistory.TotalWu = stat.TotalWu
				existingHistory.AverageWu = stat.AverageWu
				existingHistory.ZeroWuCount = stat.ZeroWuCount
				existingHistory.RecordDate = currentDate // 使用规范化后的日期
				if err := tx.Save(&existingHistory).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
