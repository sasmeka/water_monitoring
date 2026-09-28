package pkg

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Postgres_Database(db string) *sqlx.DB {
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")

	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable port=%s", host, user, password, db, port)

	database, err := sqlx.Connect("postgres", config)
	if err != nil {
		log.Fatal(err)
	}
	database.SetMaxIdleConns(5)
	database.SetMaxOpenConns(50)
	database.SetConnMaxIdleTime(3 * time.Minute)

	return database

}

var (
	clients = map[string]*redis.Client{}
	mu      sync.Mutex
)

func Redis(db string) *redis.Client {
	mu.Lock()
	defer mu.Unlock()

	// kalau sudah ada, reuse
	if client, ok := clients[db]; ok {
		return client
	}

	defaultPoolSize := 100
	defaultMinIdleConns := 20

	if poolSizeEnv := os.Getenv("REDISDBPOOL"); poolSizeEnv != "" {
		if poolSize, err := strconv.Atoi(poolSizeEnv); err == nil {
			defaultPoolSize = poolSize
		}
	}

	if minIdleConnsEnv := os.Getenv("REDISDBMINIDLE"); minIdleConnsEnv != "" {
		if minIdleConns, err := strconv.Atoi(minIdleConnsEnv); err == nil {
			defaultMinIdleConns = minIdleConns
		}
	}

	// buat sekali saja
	client := redis.NewClient(&redis.Options{
		Addr:         os.Getenv(db),
		Username:     os.Getenv("REDISUSER"),
		Password:     os.Getenv("REDISPASSWORD"),
		PoolSize:     defaultPoolSize,
		MinIdleConns: defaultMinIdleConns,
	})

	clients[db] = client
	return client
}

var (
	es   *elasticsearch.Client
	once sync.Once
)

func Elastic() *elasticsearch.Client {
	once.Do(func() {
		client, err := elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{os.Getenv("ELASTIC_HOST")},
			Username:  os.Getenv("ELASTIC_USER"),
			Password:  os.Getenv("ELASTIC_PASS"),
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		})

		if err != nil {
			log.Fatal(err)
		}

		es = client
	})

	return es
}
