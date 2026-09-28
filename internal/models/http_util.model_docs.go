package models

type HTTP200 struct {
	Status int    `json:"status" example:"200"`
	Pesan  string `json:"pesan" example:"OK / Custom Message"`
	Data   string `json:"data" example:"Data (optional)"`
}

type HTTP200_Without_Data struct {
	Status int    `json:"status" example:"200"`
	Pesan  string `json:"pesan" example:"OK / Custom Message"`
}

type HTTP201 struct {
	Status int    `json:"status" example:"201"`
	Pesan  string `json:"pesan" example:"Created"`
}

type HTTP400 struct {
	Status int    `json:"status" example:"400"`
	Pesan  string `json:"pesan" example:"Bad Request / Custom Error"`
}

type HTTP401 struct {
	Status int    `json:"status" example:"401"`
	Pesan  string `json:"pesan" example:"Unauthorized / Custom Error"`
}

type HTTP403 struct {
	Status int    `json:"status" example:"403"`
	Pesan  string `json:"pesan" example:"Forbidden / Custom Error"`
}

type HTTP404 struct {
	Status int    `json:"status" example:"404"`
	Pesan  string `json:"pesan" example:"Not Found / Custom Error"`
}

type HTTP500 struct {
	Status int    `json:"status" example:"500"`
	Pesan  string `json:"pesan" example:"Internal Server Error / Custom Error"`
}

type HTTP501 struct {
	Status int    `json:"status" example:"501"`
	Pesan  string `json:"pesan" example:"Bad Gateway / Custom Error"`
}

type HTTP304 struct {
	Status int    `json:"status" example:"304"`
	Pesan  string `json:"pesan" example:"Not Modified / Custom Error"`
}
