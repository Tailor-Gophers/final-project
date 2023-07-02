package models

type FlightClassCancellation struct {
	ID                      uint `gorm:"primaryKey"`
	CancellationConditionID uint
	FlightClassID           uint
	CancellationCondition   CancellationCondition
}
