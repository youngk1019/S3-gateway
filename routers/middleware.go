package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"s3-gateway/command/vars"
	"s3-gateway/jwt"
	"s3-gateway/log"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := jwt.FetchJWTToken(c.Request)
		claims, err := jwt.ParseJWT(token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set(vars.UIDKey, claims[vars.UIDKey])
		c.Set(vars.UNAMEKey, claims[vars.UNAMEKey])
		c.Set(vars.LoginUserKey, claims[vars.LoginUserKey])
	}
}

func Logger(reqs ...string) gin.HandlerFunc {
	req := ""
	for _, u := range reqs {
		req += u
	}

	return func(c *gin.Context) {
		c.Set(vars.UUIDKey, uuid.New().String())
		startTime := time.Now()
		log.Infow("request start",
			vars.UUIDKey, c.Value(vars.UUIDKey),
			requestUrlKey, c.Request.URL,
			requestMethod, c.Request.Method,
			requestStartTimeKey, startTime.Format(timeFormat),
			vars.UIDKey, c.Value(vars.UIDKey),
			vars.UNAMEKey, c.Value(vars.UNAMEKey),
			vars.LoginUserKey, c.Value(vars.LoginUserKey),
		)
		c.Next()
		endTime := time.Now()
		log.Infow("request finish",
			vars.UUIDKey, c.Value(vars.UUIDKey),
			requestFinishTimeKey, endTime.Format(timeFormat),
			timeElapsedKey, endTime.Sub(startTime).Seconds(),
			requestMethod, c.Request.Method,
			responseHTTPCodeKey, c.Writer.Status(),
			vars.UIDKey, c.Value(vars.UIDKey),
			vars.UNAMEKey, c.Value(vars.UNAMEKey),
			vars.LoginUserKey, c.Value(vars.LoginUserKey),
		)
	}
}
