package repositories

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/sasmeka/water_monitoring/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repo_WaterMonitoring struct {
	pgdb *sqlx.DB
}

func New_WaterMonitoring(db *sqlx.DB) *Repo_WaterMonitoring {
	repo := &Repo_WaterMonitoring{pgdb: db}
	if err := repo.AutoMigrate(); err != nil {
		log.Printf("[AutoMigrate] Gagal migrasi tabel water_monitoring: %v\n", err)
	}
	return repo
}

// AutoMigrate membuat tabel water_monitoring secara otomatis jika belum ada
func (r *Repo_WaterMonitoring) AutoMigrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS water_monitoring (
		id SERIAL PRIMARY KEY,
		id_device VARCHAR(255) NOT NULL DEFAULT 'DEV-001',
		waktu TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		suhu_air NUMERIC(10,2) NOT NULL DEFAULT 0,
		ph NUMERIC(10,2) NOT NULL DEFAULT 0,
		tds NUMERIC(10,2) NOT NULL DEFAULT 0,
		turbidity NUMERIC(10,2) NOT NULL DEFAULT 0,
		status VARCHAR(50) NOT NULL DEFAULT 'NORMAL',
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

	-- Tambahkan kolom jika tabel lama sudah ada
	ALTER TABLE water_monitoring ADD COLUMN IF NOT EXISTS id_device VARCHAR(255) NOT NULL DEFAULT 'DEV-001';
	ALTER TABLE water_monitoring ADD COLUMN IF NOT EXISTS waktu TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;

	-- Konversi tipe data waktu ke TIMESTAMPTZ jika sebelumnya tipe lain
	DO $$
	BEGIN
		BEGIN
			ALTER TABLE water_monitoring ALTER COLUMN waktu TYPE TIMESTAMPTZ USING waktu::timestamptz;
		EXCEPTION
			WHEN OTHERS THEN NULL;
		END;
		BEGIN
			ALTER TABLE water_monitoring ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at::timestamptz;
		EXCEPTION
			WHEN OTHERS THEN NULL;
		END;
	END $$;

	CREATE INDEX IF NOT EXISTS idx_water_monitoring_id_device ON water_monitoring(id_device);
	CREATE INDEX IF NOT EXISTS idx_water_monitoring_waktu ON water_monitoring(waktu);
	CREATE INDEX IF NOT EXISTS idx_water_monitoring_created_at ON water_monitoring(created_at);
	`
	_, err := r.pgdb.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

var locJakartaRepo *time.Location

func getJakartaLocationRepo() *time.Location {
	if locJakartaRepo == nil {
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			loc = time.FixedZone("WIB", 7*3600)
		}
		locJakartaRepo = loc
	}
	return locJakartaRepo
}

// parseWaktuHelper mengonversi string ke time.Time dengan timezone Asia/Jakarta
func parseWaktuHelper(waktuStr string) time.Time {
	loc := getJakartaLocationRepo()
	waktuStr = strings.TrimSpace(waktuStr)
	if waktuStr == "" {
		return time.Now().In(loc)
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", waktuStr, loc); err == nil {
		return t.In(loc)
	}
	if t, err := time.ParseInLocation("02/01/2006 15:04:05", waktuStr, loc); err == nil {
		return t.In(loc)
	}
	if t, err := time.Parse(time.RFC3339, waktuStr); err == nil {
		return t.In(loc)
	}
	return time.Now().In(loc)
}

// Create menyisipkan single data water monitoring
func (r *Repo_WaterMonitoring) Create(item *models.BodyWaterMonitoring) (*models.WaterMonitoring, error) {
	query := `
	INSERT INTO water_monitoring (id_device, waktu, suhu_air, ph, tds, turbidity, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, id_device, waktu, suhu_air, ph, tds, turbidity, status, created_at
	`
	var result models.WaterMonitoring
	parsedWaktu := parseWaktuHelper(item.Waktu)

	err := r.pgdb.QueryRowx(query,
		item.IDDevice,
		parsedWaktu,
		item.SuhuAir,
		item.PH,
		item.TDS,
		item.Turbidity,
		item.Status,
	).StructScan(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateBatch menyisipkan banyak data water monitoring sekaligus
func (r *Repo_WaterMonitoring) CreateBatch(items []models.BodyWaterMonitoring) ([]models.WaterMonitoring, error) {
	if len(items) == 0 {
		return nil, errors.New("data batch tidak boleh kosong")
	}

	tx, err := r.pgdb.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO water_monitoring (id_device, waktu, suhu_air, ph, tds, turbidity, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, id_device, waktu, suhu_air, ph, tds, turbidity, status, created_at
	`

	results := make([]models.WaterMonitoring, 0, len(items))
	for _, item := range items {
		var row models.WaterMonitoring
		parsedWaktu := parseWaktuHelper(item.Waktu)

		err := tx.QueryRowx(query,
			item.IDDevice,
			parsedWaktu,
			item.SuhuAir,
			item.PH,
			item.TDS,
			item.Turbidity,
			item.Status,
		).StructScan(&row)
		if err != nil {
			return nil, err
		}
		results = append(results, row)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return results, nil
}

// GetAll mengambil list data water monitoring dengan pagination & opsional filter id_device
func (r *Repo_WaterMonitoring) GetAll(idDevice string, limit, offset int) ([]models.WaterMonitoring, int, error) {
	var total int
	var countQuery string
	var selectQuery string
	var err error

	results := []models.WaterMonitoring{}

	if strings.TrimSpace(idDevice) != "" {
		countQuery = `SELECT COUNT(1) FROM water_monitoring WHERE id_device = $1`
		err = r.pgdb.Get(&total, countQuery, idDevice)
		if err != nil {
			return nil, 0, err
		}

		selectQuery = `
		SELECT id, id_device, waktu, suhu_air, ph, tds, turbidity, status, created_at
		FROM water_monitoring
		WHERE id_device = $1
		ORDER BY id DESC
		LIMIT $2 OFFSET $3
		`
		err = r.pgdb.Select(&results, selectQuery, idDevice, limit, offset)
	} else {
		countQuery = `SELECT COUNT(1) FROM water_monitoring`
		err = r.pgdb.Get(&total, countQuery)
		if err != nil {
			return nil, 0, err
		}

		selectQuery = `
		SELECT id, id_device, waktu, suhu_air, ph, tds, turbidity, status, created_at
		FROM water_monitoring
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
		`
		err = r.pgdb.Select(&results, selectQuery, limit, offset)
	}

	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// GetByID mengambil detail data water monitoring berdasarkan ID
func (r *Repo_WaterMonitoring) GetByID(id int) (*models.WaterMonitoring, error) {
	var result models.WaterMonitoring
	query := `
	SELECT id, id_device, waktu, suhu_air, ph, tds, turbidity, status, created_at
	FROM water_monitoring
	WHERE id = $1
	`
	err := r.pgdb.Get(&result, query, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, errors.New("data tidak ditemukan")
		}
		return nil, err
	}

	return &result, nil
}
