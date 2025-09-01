package taskmodel

type Task struct {
	// gorm.Model
	// Clustered Index(unique/ not null)
	TaskID int64 `gorm:"primarykey"` // id

	// generate index
	UserID int64 `gorm:"index"` // userid

	Status    int `gorm:"default:0"`
	Title     string
	Content   string `gorm:"type:longtext"`
	StartTime int64
	EndTime   int64
}

// temporarily cover(tablename)
func (*Task) Table() string {
	return "task"
}
