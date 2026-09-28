package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/sasmeka/water_monitoring/config"
	"github.com/sasmeka/water_monitoring/internal/models"
	"github.com/sasmeka/water_monitoring/internal/repositories"
	"github.com/sasmeka/water_monitoring/pkg"

	"github.com/gin-gonic/gin"
)

type Handler_WaterMonitoring struct {
	*repositories.Repo_WaterMonitoring
}

func New_WaterMonitoring(r *repositories.Repo_WaterMonitoring) *Handler_WaterMonitoring {
	return &Handler_WaterMonitoring{r}
}

var locJakarta *time.Location

func getJakartaLocation() *time.Location {
	if locJakarta == nil {
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			loc = time.FixedZone("WIB", 7*3600)
		}
		locJakarta = loc
	}
	return locJakarta
}

// sanitizeWaterData memeriksa dan mengisi nilai default id_device, waktu (Asia/Jakarta), & status jika kosong
func sanitizeWaterData(data *models.BodyWaterMonitoring) {
	loc := getJakartaLocation()
	now := time.Now().In(loc)

	if strings.TrimSpace(data.IDDevice) == "" {
		data.IDDevice = "DEV-001"
	}
	if strings.TrimSpace(data.Waktu) == "" {
		data.Waktu = now.Format("2006-01-02 15:04:05")
	} else {
		// Jika dikirim dengan format DD/MM/YYYY HH:mm:ss atau YYYY-MM-DD HH:mm:ss, set ke timezone Asia/Jakarta
		trimmedWaktu := strings.TrimSpace(data.Waktu)
		if parsedTime, err := time.ParseInLocation("02/01/2006 15:04:05", trimmedWaktu, loc); err == nil {
			data.Waktu = parsedTime.Format("2006-01-02 15:04:05")
		} else if parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", trimmedWaktu, loc); err == nil {
			data.Waktu = parsedTime.Format("2006-01-02 15:04:05")
		}
	}
	if strings.TrimSpace(data.Status) == "" {
		data.Status = "NORMAL"
	}
}

// Create godoc
// @Tags			Water Monitoring
// @Summary		Input Data Water Monitoring
// @Description	Menyimpan data hasil monitoring air (Suhu, pH, TDS, Turbidity, Status, dll)
// @Param			Request	body		models.BodyWaterMonitoring	true	"Data Water Monitoring"
// @Accept			json
// @Produce		json
// @Success		200		{object}	models.ResWaterMonitoring_HTTP200
// @Failure		400		{object}	models.HTTP400
// @Failure		500		{object}	models.HTTP500
// @Router			/water-monitoring [post]
func (h *Handler_WaterMonitoring) Create(ctx *gin.Context) {
	var data models.BodyWaterMonitoring
	if err := ctx.ShouldBind(&data); err != nil {
		pkg.Responses(400, &config.Result{Message: fmt.Sprintf("Format data tidak valid: %s", err.Error())}).Send(ctx)
		return
	}

	sanitizeWaterData(&data)

	res, err := h.Repo_WaterMonitoring.Create(&data)
	if err != nil {
		pkg.Responses(500, &config.Result{Message: fmt.Sprintf("Gagal menyimpan data: %s", err.Error())}).Send(ctx)
		return
	}

	pkg.Responses(200, &config.Result{Data: res, Message: "Data monitoring air berhasil disimpan"}).Send(ctx)
}

// CreateBatch godoc
// @Tags			Water Monitoring
// @Summary		Input Banyak Data Water Monitoring (Batch/Bulk)
// @Description	Menyimpan array data hasil monitoring air sekaligus
// @Param			Request	body		models.BodyWaterMonitoringBatch	true	"Batch Data Water Monitoring"
// @Accept			json
// @Produce		json
// @Success		200		{object}	models.ResWaterMonitoringList_HTTP200
// @Failure		400		{object}	models.HTTP400
// @Failure		500		{object}	models.HTTP500
// @Router			/water-monitoring/batch [post]
func (h *Handler_WaterMonitoring) CreateBatch(ctx *gin.Context) {
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		pkg.Responses(400, &config.Result{Message: "Gagal membaca body request"}).Send(ctx)
		return
	}

	var items []models.BodyWaterMonitoring

	// Coba parsing format {"data": [...]}
	var batchWrapper models.BodyWaterMonitoringBatch
	if err := json.Unmarshal(bodyBytes, &batchWrapper); err == nil && len(batchWrapper.Data) > 0 {
		items = batchWrapper.Data
	} else {
		// Jika bukan objek bertingkat, coba parsing sebagai raw array [...]
		if err := json.Unmarshal(bodyBytes, &items); err != nil {
			pkg.Responses(400, &config.Result{Message: fmt.Sprintf("Format JSON array batch tidak valid: %s", err.Error())}).Send(ctx)
			return
		}
	}

	if len(items) == 0 {
		pkg.Responses(400, &config.Result{Message: "Data array tidak boleh kosong"}).Send(ctx)
		return
	}

	for i := range items {
		sanitizeWaterData(&items[i])
	}

	res, err := h.Repo_WaterMonitoring.CreateBatch(items)
	if err != nil {
		pkg.Responses(500, &config.Result{Message: fmt.Sprintf("Gagal menyimpan batch data: %s", err.Error())}).Send(ctx)
		return
	}

	pkg.Responses(200, &config.Result{Data: res, Message: fmt.Sprintf("Berhasil menyimpan %d data monitoring air", len(res))}).Send(ctx)
}

// GetAll godoc
// @Tags			Water Monitoring
// @Summary		Ambil Riwayat Data Water Monitoring
// @Description	Mendapatkan daftar data hasil monitoring air dengan pagination dan opsional filter id_device
// @Param			id_device	query		string	false	"Filter berdasarkan ID Device (misal: DEV-001)"
// @Param			page		query		int		false	"Halaman (default 1)"
// @Param			limit		query		int		false	"Jumlah data per halaman (default 10)"
// @Accept			json
// @Produce		json
// @Success		200		{object}	models.ResWaterMonitoringList_HTTP200
// @Failure		500		{object}	models.HTTP500
// @Router			/water-monitoring [get]
func (h *Handler_WaterMonitoring) GetAll(ctx *gin.Context) {
	idDevice := ctx.Query("id_device")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	data, total, err := h.Repo_WaterMonitoring.GetAll(idDevice, limit, offset)
	if err != nil {
		pkg.Responses(500, &config.Result{Message: fmt.Sprintf("Gagal mengambil data: %s", err.Error())}).Send(ctx)
		return
	}

	lastPage := (total + limit - 1) / limit
	meta := config.Metas{
		Total_data: total,
		Last_page:  lastPage,
		Prev:       nil,
		Next:       nil,
	}
	if page > 1 {
		meta.Prev = page - 1
	}
	if page < lastPage {
		meta.Next = page + 1
	}

	pkg.Responses(200, &config.Result{
		Data: data,
		Meta: meta,
	}).Send(ctx)
}

// GetByID godoc
// @Tags			Water Monitoring
// @Summary		Ambil Detail Data Water Monitoring by ID
// @Description	Mendapatkan 1 data hasil monitoring air berdasarkan id
// @Param			id	path		int	true	"ID Data Monitoring"
// @Accept			json
// @Produce		json
// @Success		200		{object}	models.ResWaterMonitoring_HTTP200
// @Failure		400		{object}	models.HTTP400
// @Failure		404		{object}	models.HTTP404
// @Failure		500		{object}	models.HTTP500
// @Router			/water-monitoring/{id} [get]
func (h *Handler_WaterMonitoring) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		pkg.Responses(400, &config.Result{Message: "Parameter ID harus berupa angka"}).Send(ctx)
		return
	}

	data, err := h.Repo_WaterMonitoring.GetByID(id)
	if err != nil {
		pkg.Responses(404, &config.Result{Message: err.Error()}).Send(ctx)
		return
	}

	pkg.Responses(200, &config.Result{Data: data}).Send(ctx)
}
