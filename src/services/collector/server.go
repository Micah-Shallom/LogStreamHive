package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)


func (app *App) setupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(corsMiddleware())

	router.GET("/conn/:id", app.connTokenHandler)
	router.POST("/sub/:id", app.subTokenHandler)

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func (app *App) connTokenHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id parameter"})
		return
	}

	resp, code, err := app.getConnToken(id)
	if err != nil {
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(code, resp)
}

func (app *App) subTokenHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id parameter"})
		return
	}

	var req ChannelSubTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, code, err := app.getSubToken(id, req)
	if err != nil {
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(code, resp)
}

type ChannelSubTokenReq struct {
	Token   string `json:"token"`
	Channel string `json:"channel" binding:"required"`
}

func (app *App) getConnToken(id string) (gin.H, int, error) {
	claims := jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(app.Config.Centrifugo.Secret))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return gin.H{"token": signed}, http.StatusOK, nil
}

func (app *App) getSubToken(id string, req ChannelSubTokenReq) (gin.H, int, error) {
	if _, err := uuid.Parse(req.Channel); err != nil {
		return nil, http.StatusBadRequest, errors.New("invalid channel id format")
	}

	// TODO: Add channel ownership/access validation here
	// if !app.userHasAccessToChannel(id, req.Channel) {
	//     return nil, http.StatusForbidden, errors.New("access denied")
	// }

	claims := jwt.MapClaims{
		"sub":     id,
		"channel": req.Channel,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(app.Config.Centrifugo.Secret))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return gin.H{"token": signed}, http.StatusOK, nil
}

