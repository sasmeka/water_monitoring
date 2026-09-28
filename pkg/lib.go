package pkg

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sasmeka/water_monitoring/config"

	"github.com/jmoiron/sqlx"
	"github.com/mitchellh/mapstructure"
	"github.com/valyala/fasthttp"
	"github.com/zsefvlol/timezonemapper"
)

type genNomor struct {
	Branch   string `json:"branch" db:"branch"`
	Document string `json:"document" db:"document"`
	Year     string `json:"year" db:"year"`
	Periode  string `json:"periode" db:"periode"`
	Prefix   string `json:"prefix" db:"prefix"`
	Sufix    string `json:"sufix" db:"sufix"`
	Eom      string `json:"eom" db:"eom"`
	Length   string `json:"length" db:"length"`
	Reset    string `json:"reset" db:"reset"`
}

func StructtoString(data interface{}) (string, error) {
	jsonValue, err := json.Marshal(data)
	if err != nil {
		return "", errors.New("gagal generate struct to string")
	}
	return string(jsonValue), nil
}

func HttpClientSGP(method, endpoint, body string) (*config.ResponseHttpClientSGP, error) {
	url := os.Getenv("URLSGP")
	if os.Getenv("GIN_MODE") == "release" {
		url += "/antrian"
	} else {
		url += "/antrian_dev"
	}

	result := config.ResponseHttpClientSGP{}
	client := &http.Client{}
	req, err := http.NewRequest(method, url+endpoint, nil)
	if err != nil {
		return nil, errors.New("error new request: " + err.Error())
	}
	req.Header.Set("tanobel-api-key", "PrzNvotiZkwc6nWwTHSj3E4hevUJpqMw")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Body = io.NopCloser(bytes.NewBuffer([]byte(body)))

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.New("error new do req: " + err.Error())
	}
	defer res.Body.Close()

	resbody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("error read all: " + err.Error())
	}
	response := make(map[string]interface{}, 0)
	err_unmarshal := json.Unmarshal([]byte(resbody), &response)
	if err_unmarshal != nil {
		return nil, errors.New("error unmarshal: " + err_unmarshal.Error())
	}
	mapstructure.Decode(response, &result)
	if result.Status != true {
		if result.Msg != nil {
			return &result, errors.New(result.Msg.(string))
		} else {
			return &result, errors.New("Request to SGP gagal: " + response["message"].(string))
		}
	}
	return &result, nil
}

func Strtoint(data string) int {
	data = strings.Split(data, ".")[0]
	num, _ := strconv.Atoi(data)
	return num
}
func Inttostr(data int) string {
	str := strconv.Itoa(data)
	return str
}

func StrToFloat64(data string) float64 {
	data = strings.Split(data, ".")[0]
	num, _ := strconv.ParseFloat(strings.ReplaceAll(data, ",", ""), 64)
	return num
}
func Float64toStr(data float64) string {
	str := strconv.FormatFloat(data, 'f', 2, 64)
	return str
}

func TransformDate(data string) string {
	if data != "" {
		t, _ := time.Parse(time.RFC3339, data)
		return t.Format("2006-01-02 15:04:05")
	} else {
		return ""
	}
}

func InsertValueNull(data *string) *string {
	var temp = ""
	if data == nil {
		return &temp
	} else {
		return data
	}
}

func GenNomorV2(db *sqlx.Tx, branch, eom, document, year, periode, prefix, sufix, length, reset string) (string, error) {
	database := db
	resultquery := []genNomor{}
	data := genNomor{
		Branch:   branch,
		Document: document,
		Year:     year,
		Periode:  periode,
		Prefix:   prefix,
		Sufix:    sufix,
		Eom:      eom,
		Length:   length,
		Reset:    reset,
	}
	var (
		nameDB       string
		cekcountdata string
		counterr     error
	)

	err_dbname := database.Get(&nameDB, database.Rebind("SELECT DB_NAME()"))
	if err_dbname != nil {
		database.Rollback()
		return "", err_dbname
	}

	// ---- insert new t nomor
	if reset == "0" {
		get_err := database.Select(&resultquery, database.Rebind(`select top 1 cBranch branch, cEOM eom, cDocument document, TRIM(cYear) year, TRIM(cPeriode) periode, TRIM(cPrefix) prefix, TRIM(cSufix) sufix, nlength length, cReset reset from t_nomor WITH(NOLOCK) WHERE cbranch=? and cDocument=? and cReset='0' order by cYear DESC, cPeriode DESC`), data.Branch, data.Document, data.Reset)
		if get_err != nil {
			database.Rollback()
			return "", get_err
		}
		if len(resultquery) == 0 {
			_, err_insert := database.Exec(database.Rebind(`insert into t_nomor WITH(XLOCK) values (?,?,?,?,?,?,?,0,?,?)`), data.Branch, data.Eom, data.Document, data.Year, data.Periode, data.Prefix, data.Sufix, data.Length, data.Reset)
			if err_insert != nil {
				database.Rollback()
				return "", errors.New(fmt.Sprintf("Gagal insert t_nomor baru dokumen '%s': ", document) + err_insert.Error())
			}
		} else {
			if resultquery[0].Year == data.Year && resultquery[0].Periode == data.Periode {
				data.Branch = resultquery[0].Branch
				data.Document = resultquery[0].Document
				data.Year = resultquery[0].Year
				data.Periode = resultquery[0].Periode
			} else {
				var ndocno int
				get_err := database.Get(&ndocno, database.Rebind(`select top 1 ndocno from t_nomor WITH(NOLOCK) WHERE cbranch=? and cDocument=? and cReset='0' order by cYear DESC, cPeriode DESC`), data.Branch, data.Document, data.Reset)
				if get_err != nil {
					database.Rollback()
					return "", get_err
				}
				_, err_insert := database.Exec(database.Rebind(`insert into t_nomor WITH(XLOCK) values (?,?,?,?,?,?,?,?,?,?)`), data.Branch, data.Eom, data.Document, data.Year, data.Periode, data.Prefix, data.Sufix, ndocno, data.Length, data.Reset)
				if err_insert != nil {
					database.Rollback()
					return "", errors.New(fmt.Sprintf("Gagal insert t_nomor baru dokumen '%s': ", document) + err_insert.Error())
				}
			}
		}
	} else {
		counterr = database.Get(&cekcountdata, database.Rebind(`select count(*) from t_nomor WITH(NOLOCK) WHERE cbranch=? and cDocument=? and cYear=? and cPeriode=?`), data.Branch, data.Document, data.Year, data.Periode)
		if counterr != nil {
			database.Rollback()
			return "", counterr
		}
		if cekcountdata == "0" {
			_, err_insert := database.Exec(database.Rebind(`insert into t_nomor WITH(XLOCK) values (?,?,?,?,?,?,?,0,?,?)`), data.Branch, data.Eom, data.Document, data.Year, data.Periode, data.Prefix, data.Sufix, data.Length, data.Reset)
			if err_insert != nil {
				database.Rollback()
				return "", errors.New(fmt.Sprintf("Gagal insert t_nomor baru dokumen '%s': ", document) + err_insert.Error())
			}
		}
	}

	// ----
	var result *string
	err := database.Get(&result, database.Rebind(`EXEC [pr_nomor] ?, ?, ?, ?`), data.Branch, data.Document, data.Year, data.Periode)
	if err != nil {
		database.Rollback()
		return "", err
	}
	if result == nil || *result == "" {
		database.Rollback()
		return "", errors.New("Generate nomor gagal (NULL) - Hubungi IT")
	}

	return strings.TrimSpace(*result), nil
}
func UpdateNomorV2(db *sqlx.Tx, branch, eom, document, year, periode, prefix, sufix, length, reset string) error {
	database := db
	data := genNomor{
		Branch:   branch,
		Document: document,
		Year:     year,
		Periode:  periode,
		Prefix:   prefix,
		Sufix:    sufix,
		Eom:      eom,
		Length:   length,
		Reset:    reset,
	}
	var (
		nameDB string
	)

	err_dbname := database.Get(&nameDB, database.Rebind("SELECT DB_NAME()"))
	if err_dbname != nil {
		return err_dbname
	}
	if strings.Contains(nameDB, "b2b") {
		_, err := database.NamedExec(`update t_nomor WITH(UPDLOCK) set ndocno=ndocno+1 where csuppcode=:branch and cDocument=:document and cYear=:year and cPeriode=:periode`, data)
		if err != nil {
			return errors.New("Gagal update t_nomor: " + err.Error())
		}
	} else {
		_, err := database.NamedExec(`update t_nomor WITH(UPDLOCK) set ndocno=ndocno+1 where cBranch=:branch and cDocument=:document and cYear=:year and cPeriode=:periode`, data)
		if err != nil {
			return errors.New("Gagal update t_nomor: " + err.Error())
		}
	}
	return nil
}

func GenNomor(db *sqlx.Tx, branch, eom, document, year, periode, prefix, sufix, length, reset string) (string, error) {
	database := db
	resultquery := []genNomor{}
	data := genNomor{
		Branch:   branch,
		Document: document,
		Year:     year,
		Periode:  periode,
		Prefix:   prefix,
		Sufix:    sufix,
		Eom:      eom,
		Length:   length,
		Reset:    reset,
	}
	var (
		nameDB       string
		cekcountdata string
		counterr     error
	)

	err_dbname := database.Get(&nameDB, database.Rebind("SELECT DB_NAME()"))
	if err_dbname != nil {
		database.Rollback()
		return "", err_dbname
	}

	// ---- insert new t nomor
	if reset == "0" {
		get_err := database.Select(&resultquery, database.Rebind(`select top 1 cBranch branch, cEOM eom, cDocument document, TRIM(cYear) year, TRIM(cPeriode) periode, TRIM(cPrefix) prefix, TRIM(cSufix) sufix, nlength length, cReset reset from d_transaksi.dbo.t_nomor WITH (NOLOCK) WHERE cbranch=? and cDocument=? and cReset='0' order by cYear DESC, cPeriode DESC`), data.Branch, data.Document, data.Reset)
		if get_err != nil {
			database.Rollback()
			return "", get_err
		}
		if len(resultquery) == 0 {
			_, err_insert := database.Exec(database.Rebind(`insert into d_transaksi.dbo.t_nomor values (?,?,?,?,?,?,?,0,?,?)`), data.Branch, data.Eom, data.Document, data.Year, data.Periode, data.Prefix, data.Sufix, data.Length, data.Reset)
			if err_insert != nil {
				database.Rollback()
				return "", errors.New(fmt.Sprintf("Gagal insert t_nomor baru dokumen '%s': ", document) + err_insert.Error())
			}
		} else {
			if resultquery[0].Year == data.Year && resultquery[0].Periode == data.Periode {
				data.Branch = resultquery[0].Branch
				data.Document = resultquery[0].Document
				data.Year = resultquery[0].Year
				data.Periode = resultquery[0].Periode
			} else {
				var ndocno int
				get_err := database.Get(&ndocno, database.Rebind(`select top 1 ndocno from d_transaksi.dbo.t_nomor WITH (NOLOCK) WHERE cbranch=? and cDocument=? and cReset='0' order by cYear DESC, cPeriode DESC`), data.Branch, data.Document, data.Reset)
				if get_err != nil {
					database.Rollback()
					return "", get_err
				}
				_, err_insert := database.Exec(database.Rebind(`insert into d_transaksi.dbo.t_nomor values (?,?,?,?,?,?,?,?,?,?)`), data.Branch, data.Eom, data.Document, data.Year, data.Periode, data.Prefix, data.Sufix, ndocno, data.Length, data.Reset)
				if err_insert != nil {
					database.Rollback()
					return "", errors.New(fmt.Sprintf("Gagal insert t_nomor baru dokumen '%s': ", document) + err_insert.Error())
				}
			}
		}
	} else {
		counterr = database.Get(&cekcountdata, database.Rebind(`select count(*) from d_transaksi.dbo.t_nomor WITH (NOLOCK) WHERE cBranch=? and cDocument=? and cYear=? and cPeriode=?`), data.Branch, data.Document, data.Year, data.Periode)
		if counterr != nil {
			database.Rollback()
			return "", counterr
		}
		if cekcountdata == "0" {
			_, err_insert := database.Exec(database.Rebind(`insert into d_transaksi.dbo.t_nomor values (?,?,?,?,?,?,?,0,?,?)`), data.Branch, data.Eom, data.Document, data.Year, data.Periode, data.Prefix, data.Sufix, data.Length, data.Reset)
			if err_insert != nil {
				database.Rollback()
				return "", errors.New(fmt.Sprintf("Gagal insert t_nomor baru dokumen '%s': ", document) + err_insert.Error())
			}
		}
	}

	// ----
	var result *string
	err := database.Get(&result, database.Rebind(`EXEC [pr_nomor] ?, ?, ?, ?`), data.Branch, data.Document, data.Year, data.Periode)
	if err != nil {
		database.Rollback()
		return "", err
	}
	if result == nil || *result == "" {
		database.Rollback()
		return "", errors.New("Generate nomor gagal (NULL) - Hubungi IT")
	} else {
		if strings.Contains(nameDB, "b2b") {
			_, err := database.NamedExec(`update t_nomor WITH (UPDLOCK) set ndocno=ndocno+1 where csuppcode=:branch and cDocument=:document and cYear=:year and cPeriode=:periode`, data)
			if err != nil {
				database.Rollback()
				return "", errors.New("Gagal Query Update: " + err.Error())
			}
		} else {
			_, err := database.NamedExec(`update d_transaksi.dbo.t_nomor WITH (UPDLOCK) set ndocno=ndocno+1 where cBranch=:branch and cDocument=:document and cYear=:year and cPeriode=:periode`, data)
			if err != nil {
				database.Rollback()
				return "", errors.New("Gagal Query Update: " + err.Error())
			}
		}
	}
	return strings.TrimSpace(*result), nil
}

func GenNomor2(db *sqlx.Tx, branch, eom, document, year, periode, prefix, sufix, length, reset string) (string, error) {
	database := db
	resultquery := []genNomor{}
	data := genNomor{
		Branch:   branch,
		Document: document,
		Year:     year,
		Periode:  periode,
		Prefix:   prefix,
		Sufix:    sufix,
		Eom:      eom,
		Length:   length,
		Reset:    reset,
	}
	var (
		nameDB       string
		cekcountdata string
		counterr     error
	)

	err_dbname := database.Get(&nameDB, database.Rebind("SELECT DB_NAME()"))
	if err_dbname != nil {
		database.Rollback()
		return "", err_dbname
	}

	// ---- insert new t nomor
	if reset == "0" {
		get_err := database.Select(&resultquery, database.Rebind(`select top 1 cBranch branch, cEOM eom, cDocument document, TRIM(cYear) year, TRIM(cPeriode) periode, TRIM(cPrefix) prefix, TRIM(cSufix) sufix, nlength length, cReset reset from d_transaksi.dbo.t_nomor WITH (NOLOCK) WHERE cbranch=? and cDocument=? and cReset='0' order by cYear DESC, cPeriode DESC`), data.Branch, data.Document, data.Reset)
		if get_err != nil {
			database.Rollback()
			return "", get_err
		}
		if len(resultquery) == 0 {
			_, err_insert := database.Exec(database.Rebind(`insert into d_transaksi.dbo.t_nomor values (?,?,?,?,?,?,?,0,?,?)`), data.Branch, data.Eom, data.Document, data.Year, data.Periode, data.Prefix, data.Sufix, data.Length, data.Reset)
			if err_insert != nil {
				database.Rollback()
				return "", errors.New(fmt.Sprintf("Gagal insert t_nomor baru dokumen '%s': ", document) + err_insert.Error())
			}
		} else {
			if resultquery[0].Year == data.Year && resultquery[0].Periode == data.Periode {
				data.Branch = resultquery[0].Branch
				data.Document = resultquery[0].Document
				data.Year = resultquery[0].Year
				data.Periode = resultquery[0].Periode
			} else {
				var ndocno int
				get_err := database.Get(&ndocno, database.Rebind(`select top 1 ndocno from d_transaksi.dbo.t_nomor WITH (NOLOCK) WHERE cbranch=? and cDocument=? and cReset='0' order by cYear DESC, cPeriode DESC`), data.Branch, data.Document, data.Reset)
				if get_err != nil {
					database.Rollback()
					return "", get_err
				}
				_, err_insert := database.Exec(database.Rebind(`insert into d_transaksi.dbo.t_nomor values (?,?,?,?,?,?,?,?,?,?)`), data.Branch, data.Eom, data.Document, data.Year, data.Periode, data.Prefix, data.Sufix, ndocno, data.Length, data.Reset)
				if err_insert != nil {
					database.Rollback()
					return "", errors.New(fmt.Sprintf("Gagal insert t_nomor baru dokumen '%s': ", document) + err_insert.Error())
				}
			}
		}
	} else {
		counterr = database.Get(&cekcountdata, database.Rebind(`select count(*) from d_transaksi.dbo.t_nomor WITH (NOLOCK) WHERE cBranch=? and cDocument=? and cYear=? and cPeriode=?`), data.Branch, data.Document, data.Year, data.Periode)
		if counterr != nil {
			database.Rollback()
			return "", counterr
		}
		if cekcountdata == "0" {
			_, err_insert := database.Exec(database.Rebind(`insert into d_transaksi.dbo.t_nomor values (?,?,?,?,?,?,?,0,?,?)`), data.Branch, data.Eom, data.Document, data.Year, data.Periode, data.Prefix, data.Sufix, data.Length, data.Reset)
			if err_insert != nil {
				database.Rollback()
				return "", errors.New(fmt.Sprintf("Gagal insert t_nomor baru dokumen '%s': ", document) + err_insert.Error())
			}
		}
	}

	// ----
	var result *string
	err := database.Get(&result, database.Rebind(`EXEC [pr_nomor] ?, ?, ?, ?`), data.Branch, data.Document, data.Year, data.Periode)
	if err != nil {
		database.Rollback()
		return "", err
	}
	return strings.TrimSpace(*result), nil
}

func FetchAPI(url string) ([]byte, error) {
	// Create a new fasthttp client
	client := &fasthttp.Client{}

	// Create a request object
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	// Set the request URL
	req.SetRequestURI(url)

	// Create a response object
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// Send the request and get the response
	if err := client.Do(req, resp); err != nil {
		return nil, err
	}

	// Return the response body
	return resp.Body(), nil
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GetDatetimeByTimezone(lat, long float64) (string, error) {
	// tz := latlong.LookupZoneName(lat, long)
	tz := timezonemapper.LatLngToTimezoneString(lat, long)
	if tz == "Indian/Cocos" {
		tz = "Asia/Jakarta"
	}
	// Get the current time in the timezone
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return err.Error(), err
	}

	t := time.Now().In(loc).Format("2006-01-02 15:04:05")

	return t, nil
}

func GetTimeByTimezone(lat, long float64) (string, error) {
	// tz := latlong.LookupZoneName(lat, long)
	tz := timezonemapper.LatLngToTimezoneString(lat, long)

	// Get the current time in the timezone
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return err.Error(), err
	}

	t := time.Now().In(loc).Format("15:04:05")

	return t, nil
}

func GetTimezoneByLongLat(lat, long float64) (int, string) {
	// Use the timezone function from the latlong library to get the timezone
	// tz := latlong.LookupZoneName()
	tz := timezonemapper.LatLngToTimezoneString(lat, long)

	// TIME ZONE CODE
	// WIB : 0
	// WITA : 1
	// WIT : 2
	var timezoneCode int
	var timeZone string
	switch value := tz; value {
	case "Asia/Jakarta":
		timezoneCode = 0
		timeZone = "WIB"
	case "Asia/Makassar":
		timezoneCode = 1
		timeZone = "WITA"
	case "Asia/Jayapura":
		timezoneCode = 2
		timeZone = "WIT"
	}

	return timezoneCode, timeZone
}

func DateTimeFormat(param string) string {
	curTime := time.Now()
	var varturn string
	switch param {
	case "year":
		varturn = curTime.Format("2006")
	case "month":
		varturn = curTime.Format("01")
	case "day":
		varturn = curTime.Format("02")
	case "hour":
		varturn = curTime.Format("15")
	case "minute":
		varturn = curTime.Format("04")
	case "second":
		varturn = curTime.Format("05")
	case "his":
		varturn = curTime.Format("15:04:05")
	case "datetime":
		varturn = curTime.Format("2006-01-02 15:04:05")
	case "date-":
		varturn = curTime.Format("2006-01-02")
	case "date":
		varturn = curTime.Format("20060102")
	}

	return varturn
}

func EuclideanDistance(faceDB, faceInput []float64) (float64, error) {
	if faceDB == nil || faceInput == nil {
		return 0.0, errors.New("null argument")
	}

	if len(faceDB) != len(faceInput) {
		return 0.0, errors.New("input lists must have the same length")
	}

	sum := 0.0
	for i := 0; i < len(faceDB); i++ {
		sum += math.Pow(faceDB[i]-faceInput[i], 2)
	}

	return math.Sqrt(sum), nil
}

func JsonStringtoArray(data string) ([]float64, error) {
	var facePointInput []float64
	if err := json.Unmarshal([]byte(data), &facePointInput); err != nil {
		return nil, err
	}
	return facePointInput, nil
}

func SendPayment(jsonMap any, endpoint string) (map[string]interface{}, error) {
	json_string, err := json.Marshal(jsonMap)
	if err != nil {
		return nil, errors.New("Error Marshal req to Payment, " + err.Error())
	}
	resp, err := http.Post(os.Getenv("SERV_API_PAYMENT")+endpoint, "application/json", io.NopCloser(bytes.NewBuffer(json_string)))
	if err != nil {
		return nil, errors.New("Error HTTP Req POST, " + err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Error Reading Response Body, " + err.Error())
	}

	var result map[string]interface{}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, errors.New("Error Unmarshal Response Body, " + err.Error())
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(result["pesan"].(string))
	}

	if result["data"] == nil {
		return nil, nil
	}
	return result["data"].(map[string]interface{}), nil
}

func SendNotif(phone, pesan string) error {
	if phone == "" {
		return fmt.Errorf("phone tidak ditemukan")
	}
	pesan = fmt.Sprintf("%s\n\n_Pesan Otomatis dari PT Sentralsari Primasentosa_", pesan)
	jsonPayload := map[string]string{
		"phone": phone,
		"body":  pesan,
	}

	payload, err := json.Marshal(jsonPayload)
	if err != nil {
		log.Println("error marshal,", err)
		return fmt.Errorf("error marshal, %s", err.Error())
	}

	url := fmt.Sprintf("%s/chat/send/text?token=%s", os.Getenv("SERV_WUZZAPI"), os.Getenv("SERV_WUZZAPI_TOKEN"))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("error hit api notif wa,", err)
		return fmt.Errorf("error hit api notif wa, %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("status code not 200 but", resp.StatusCode)
		return fmt.Errorf("status code not 200")
	}

	return nil
}

func Logout(token string) {
	rdb := Redis("REDISDB")
	ctxx := context.Background()

	arrDel := []string{
		fmt.Sprintf("%s", token),
	}

	deleteKeys := make(map[string]struct{})

	for _, key := range arrDel {
		// wildcard
		if strings.Contains(key, "*") {
			iter := rdb.Scan(ctxx, 0, key, 0).Iterator()
			for iter.Next(ctxx) {
				deleteKeys[iter.Val()] = struct{}{}
			}
			continue
		}

		// key notif username -> value adalah nama key lain
		if strings.Contains(key, "NOTIF-SETUP-USERNAME") {
			val, _ := rdb.Get(ctxx, key).Result()
			if val != "" {
				deleteKeys[val] = struct{}{}
			}
		}

		// key utama
		deleteKeys[key] = struct{}{}
	}

	// convert ke slice
	keys := make([]string, 0, len(deleteKeys))
	for k := range deleteKeys {
		keys = append(keys, k)
	}

	// sekali pipeline
	if len(keys) > 0 {
		pipe := rdb.Pipeline()
		pipe.Del(ctxx, keys...)
		_, _ = pipe.Exec(ctxx)
	}
}
