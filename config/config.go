package config

import (
	"mime/multipart"

	"github.com/gin-contrib/cors"
)

type Metas struct {
	Next       interface{} `json:"next"`
	Prev       interface{} `json:"prev"`
	Last_page  interface{} `json:"last_page"`
	Total_data interface{} `json:"total_data"`
}

type Result struct {
	Data             interface{} `json:"data"`
	Data1            interface{} `json:"data1"`
	Data2            interface{} `json:"data2"`
	Data3            interface{} `json:"data3"`
	Data_Now         interface{} `json:"data_now"`
	Data_Best_Seller interface{} `json:"data_best_seller"`
	Data_Sum_Month   interface{} `json:"data_sum_month"`
	Registered       interface{} `json:"registered"`
	Notregistered    interface{} `json:"notregistered"`
	Data_bonus       interface{} `json:"data_bonus"`
	Data_harga       interface{} `json:"data_harga"`
	Data_diskon      interface{} `json:"data_diskon"`
	Data_barang      interface{} `json:"data_barang"`
	Data_stok        interface{} `json:"data_stok"`
	Data_ppn         interface{} `json:"data_ppn"`
	DataMst          interface{} `json:"dataMst"`
	DataDtl          interface{} `json:"dataDtl"`
	Data_mst         interface{} `json:"data_mst"`
	Data_dtl         interface{} `json:"data_dtl"`
	DataCust         interface{} `json:"dataCust"`
	DataSales        interface{} `json:"dataSales"`
	Data_dtl1        interface{} `json:"data_dtl1"`
	Data_dtl2        interface{} `json:"data_dtl2"`
	Climitair        interface{} `json:"climitair"`
	Climitberas      interface{} `json:"climitberas"`
	Meta             interface{} `json:"meta"`
	Message          interface{} `json:"pesan"`
	Status           interface{} `json:"status"`
	Token            interface{} `json:"token"`
	Step             interface{} `json:"step"`
}

var CorsConfig = cors.Config{
	AllowOrigins:     []string{"*"},
	AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "HEAD", "OPTIONS"},
	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	ExposeHeaders:    []string{"Content-Length"},
	AllowCredentials: true,
}

type ResponseHttpClientSGP struct {
	Status  interface{} `json:"status,omitempty" form:"status"`
	Msg     interface{} `json:"msg,omitempty" form:"msg"`
	Result  interface{} `json:"result,omitempty" form:"result"`
	Message interface{} `json:"message,omitempty" form:"message"`
}

type File struct {
	Source    *multipart.FileHeader
	Name      string
	Extention string
	Size      int64 //satuannya bytes
	Mime      string
}
