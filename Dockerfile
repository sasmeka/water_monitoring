FROM registry.tanobel.net/tanobel/public/golang:1.24.1
ENV TZ=Asia/Jakarta
WORKDIR /app/source
COPY . .
RUN go mod download
RUN go build -v -o /app/api-go ./cmd/main.go
WORKDIR /app
COPY cmd/rapidoc /app/rapidoc
RUN rm -rf /app/source
EXPOSE 8080
CMD ["./api-go"]

#
# Database (Transaksi, DBSGPTIS )
# dockebuild -t github.com/sasmeka/water_monitoring .
# docker run --name github.com/sasmeka/water_monitoring-app -e DB_Tr="sqlserver://tis:tanobel@192.168.100.17:1433?database=d_transaksi&connection+timeout=7000" -e JWTEXPIRE=30 -e JWTSECRET=EpanGantengPuol -p 9000:8080 github.com/sasmeka/water_monitoring
