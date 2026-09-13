package models

type Theatre struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`

	Screens []Screen `gorm:"foreignKey:TheatreID" json:"-"`
}
