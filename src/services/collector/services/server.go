package services

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func (app *App) SetupRouter() *gin.Engine {
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
		rd := BuildErrorResponse(http.StatusBadRequest, "error", "missing id parameter", "id parameter is required", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	resp, code, err := app.getConnToken(id)
	if err != nil {
		rd := BuildErrorResponse(code, "error", "failed to generate connection token", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	rd := BuildSuccessResponse(code, "connection token generated successfully", resp)
	c.JSON(code, rd)
}

func (app *App) subTokenHandler(c *gin.Context) {
	var req ChannelSubTokenReq

	id := c.Param("id")
	if id == "" {
		rd := BuildErrorResponse(http.StatusBadRequest, "error", "missing id parameter", "id parameter is required", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		rd := BuildErrorResponse(http.StatusBadRequest, "error", "invalid request body", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	resp, code, err := app.getSubToken(id, req)
	if err != nil {
		rd := BuildErrorResponse(code, "error", "failed to generate subscription token", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	rd := BuildSuccessResponse(code, "subscription token generated successfully", resp)
	c.JSON(code, rd)
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
