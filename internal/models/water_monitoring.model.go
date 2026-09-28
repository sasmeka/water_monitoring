package models

import "time"

type WaterMonitoring struct {
	ID        int       `json:"id" db:"id"`
	IDDevice  string    `json:"id_device" db:"id_device" example:"DEV-001"`
	Waktu     time.Time `json:"waktu" db:"waktu" example:"2026-09-15T07:14:50Z"`
	SuhuAir   float64   `json:"suhu_air" db:"suhu_air" example:"26.75"`
	PH        float64   `json:"ph" db:"ph" example:"7.51"`
	TDS       float64   `json:"tds" db:"tds" example:"55.21"`
	Turbidity float64   `json:"turbidity" db:"turbidity" example:"2049.89"`
	Status    string    `json:"status" db:"status" example:"TIDAK_NORMAL"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type BodyWaterMonitoring struct {
	IDDevice  string  `json:"id_device" form:"id_device" valid:"-" example:"DEV-001"`
	Waktu     string  `json:"waktu" form:"waktu" valid:"-" example:"2026-09-15 07:14:50"`
	SuhuAir   float64 `json:"suhu_air" form:"suhu_air" valid:"required~suhu_air is required" example:"26.75"`
	PH        float64 `json:"ph" form:"ph" valid:"required~ph is required" example:"7.51"`
	TDS       float64 `json:"tds" form:"tds" valid:"required~tds is required" example:"55.21"`
	Turbidity float64 `json:"turbidity" form:"turbidity" valid:"required~turbidity is required" example:"2049.89"`
	Status    string  `json:"status" form:"status" valid:"-" example:"TIDAK_NORMAL"`
}

type BodyWaterMonitoringBatch struct {
	Data []BodyWaterMonitoring `json:"data" form:"data" valid:"required~data is required"`
}

type ResWaterMonitoring_HTTP200 struct {
	Status int             `json:"status" example:"200"`
	Pesan  string          `json:"pesan" example:"OK"`
	Data   WaterMonitoring `json:"data"`
}

type ResWaterMonitoringList_HTTP200 struct {
	Status int               `json:"status" example:"200"`
	Pesan  string            `json:"pesan" example:"OK"`
	Data   []WaterMonitoring `json:"data"`
}
