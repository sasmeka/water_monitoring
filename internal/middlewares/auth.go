package middleware

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sasmeka/water_monitoring/config"
	"github.com/sasmeka/water_monitoring/pkg"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// func AuthJwt(role ...string) gin.HandlerFunc {
func AuthJwt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// var valid bool
		var header string

		if header = ctx.GetHeader("Authorization"); header == "" {
			pkg.Responses(401, &config.Result{Message: "Please login"}).Send(ctx)
			return
		}

		if !strings.Contains(header, "Bearer") {
			// pkg.Responses(401, &config.Result{Message: "Invalid header value"}).Send(ctx)
			// return
			header = fmt.Sprintf(`Bearer %s`, header)
		}

		tokens := strings.Replace(header, "Bearer ", "", -1)
		check, err := pkg.VerifyToken(tokens)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			pkg.Responses(401, &config.Result{Message: err.Error()}).Send(ctx)
			return
		}

		var (
			rdb  = pkg.Redis("REDISDB")
			ctxx = context.Background()
			key  = tokens
		)

		// Pipeline: ambil token saja
		pipe := rdb.Pipeline()
		tokenCmd := pipe.HGet(ctxx, key, "token")
		_, err = pipe.Exec(ctxx)
		if err != nil && err != redis.Nil {
			pkg.Logout(check.Username)
			pkg.Responses(401, &config.Result{Message: "logout"}).Send(ctx)
			return
		}
		token, err := tokenCmd.Result()

		// Jika key / field token tidak ada
		if err == redis.Nil || token == "" {
			pkg.Logout(check.Username)
			pkg.Responses(401, &config.Result{Message: "logout"}).Send(ctx)
			return
		}

		// Jika token beda
		if header != fmt.Sprintf("Bearer %s", token) {
			pkg.Logout(check.Username)
			pkg.Responses(401, &config.Result{Message: "logout"}).Send(ctx)
			return
		}

		// Jika valid baru refresh expired
		loginTime, _ := strconv.Atoi(os.Getenv("LOGIN_TIME"))
		err = rdb.Expire(ctxx, key, time.Duration(loginTime)*time.Minute).Err()
		if err != nil {
			pkg.Logout(check.Username)
			pkg.Responses(401, &config.Result{Message: "logout"}).Send(ctx)
			return
		}

		ctx.Set("username", check.Username)
		ctx.Set("platform", check.Platform)
		ctx.Set("token", tokens)
		ctx.Set("role_id", check.RoleID)
		ctx.Set("id", check.ID)
		ctx.Next()

	}
}
