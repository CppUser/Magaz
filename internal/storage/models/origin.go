package models

// City represents a city.
type City struct {
	ID       uint      `gorm:"primaryKey"`
	Name     string    `gorm:"size:100;not null"`
	Areas    []Area    `gorm:"foreignKey:CityID;constraint:OnDelete:CASCADE;"`
	Products []Product `gorm:"many2many:city_products;constraint:OnDelete:CASCADE;"`
}

// Area represents an area within a city.
type Area struct {
	ID       uint      `gorm:"primaryKey"`
	CityID   uint      `gorm:"index;not null"`
	Name     string    `gorm:"size:100;not null"`
	City     City      `gorm:"foreignKey:CityID;constraint:OnDelete:CASCADE;"`
	Products []Product `gorm:"many2many:area_products;constraint:OnDelete:CASCADE;"`
}
