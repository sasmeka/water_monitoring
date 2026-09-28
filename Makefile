# go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
# ATAU
# Hapus yang lama
# rm $(which migrate)

# Unduh dan pasang versi resmi (sudah include driver PostgreSQL)
# curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz -o migrate.tar.gz
# tar -xzf migrate.tar.gz
# sudo mv migrate /usr/local/bin/

# Cek versi
# migrate -version

# Inisialisasi File
# make migrate-pgdb-crm-init name=t_customer

# -------------------------------------------------
DB_PGDB_CRM_SOURCE="postgres://postgres:r4h4si4@172.16.27.41:5432/d_crm?sslmode=disable&search_path=public"
MIGRATIONS_PGDB_CRM_DIR=./migrations/pgdb/d_crm

migrate-pgdb-crm-init:
	migrate create -dir ${MIGRATIONS_PGDB_CRM_DIR} -ext sql -seq $(name)

DB_DBMS_CRM_SOURCE="sqlserver://sa:SPSj4y4123@172.27.28.20:1433?database=d_crm&connection+timeout=7000"
MIGRATIONS_DBMS_CRM_DIR=./migrations/dbms/d_crm

migrate-dbms-crm-init:
	migrate create -dir ${MIGRATIONS_DBMS_CRM_DIR} -ext sql -seq $(name)

DB_DBMS_TRANS_SOURCE="sqlserver://sa:SPSj4y4123@172.27.28.20:1433?database=d_transaksi&connection+timeout=7000"
MIGRATIONS_DBMS_TRANS_DIR=./migrations/dbms/d_transaksi


migrate-dbms-trans-init:
	migrate create -dir ${MIGRATIONS_DBMS_TRANS_DIR} -ext sql -seq $(name)

DB_DBMS_REPORT_SOURCE="sqlserver://sa:SPSj4y4123@172.27.28.20:1433?database=d_report&connection+timeout=7000"
MIGRATIONS_DBMS_REPORT_DIR=./migrations/dbms/d_report

migrate-dbms-report-init:
	migrate create -dir ${MIGRATIONS_DBMS_REPORT_DIR} -ext sql -seq $(name)

DB_DBMS_WEBIS_SOURCE="sqlserver://sa:SPSj4y4123@172.27.28.20:1433?database=d_webis&connection+timeout=7000"
MIGRATIONS_DBMS_WEBIS_DIR=./migrations/dbms/d_webis

migrate-dbms-webis-init:
	migrate create -dir ${MIGRATIONS_DBMS_WEBIS_DIR} -ext sql -seq $(name)

DB_DBMS_LOG_SOURCE="sqlserver://sa:SPSj4y4123@172.27.28.20:1433?database=d_log&connection+timeout=7000"
MIGRATIONS_DBMS_LOG_DIR=./migrations/dbms/d_log

migrate-dbms-log-init:
	migrate create -dir ${MIGRATIONS_DBMS_LOG_DIR} -ext sql -seq $(name)

DB_DBMS_SFA_SOURCE="sqlserver://sa:SPSj4y4123@172.27.28.20:1433?database=d_sfa&connection+timeout=7000"
MIGRATIONS_DBMS_SFA_DIR=./migrations/dbms/d_sfa

migrate-dbms-sfa-init:
	migrate create -dir ${MIGRATIONS_DBMS_SFA_DIR} -ext sql -seq $(name)

migrate-up:
	migrate -path ${MIGRATIONS_PGDB_CRM_DIR} -database ${DB_PGDB_CRM_SOURCE} -verbose up

migrate-down:
	migrate -path ${MIGRATIONS_PGDB_CRM_DIR} -database ${DB_PGDB_CRM_SOURCE} -verbose down

# migrate-fix:
# 	migrate -path ${MIGRATIONS_PGDB_CRM_DIR} -database ${DB_PGDB_CRM_SOURCE} force 0

# make migrate-up/down/fix