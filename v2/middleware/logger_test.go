package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/masonhubco/rebar/v2/middleware"
	"github.com/masonhubco/rebar/v2/mocks"
)

func Test_Logger(t *testing.T) {
	okHandler := func(c *gin.Context) {
		c.String(http.StatusOK, "200 OK")
	}
	unauthorizedHandler := func(c *gin.Context) {
		c.AbortWithError(http.StatusUnauthorized, errors.New("401 Unauthorized"))
	}

	tests := []struct {
		name         string
		givenReqPath string
		mockLogger   func(logger *mocks.Logger)
	}{
		{
			name:         "request ok path and log an info entry",
			givenReqPath: "/ok",
			mockLogger: func(logger *mocks.Logger) {
				logger.EXPECT().With(gomock.Any())
				logger.EXPECT().Info("[rebar] /ok", gomock.Any())
			},
		},
		{
			name:         "request unauthorized path and log an error entry",
			givenReqPath: "/unauthorized",
			mockLogger: func(logger *mocks.Logger) {
				logger.EXPECT().With(gomock.Any())
				logger.EXPECT().Error("[rebar] /unauthorized", gomock.Any())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := mocks.NewLogger(ctrl)
			tc.mockLogger(logger)

			router := gin.New()
			router.Use(middleware.Logger(logger))
			router.GET("/ok", okHandler)
			router.GET("/unauthorized", unauthorizedHandler)

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.givenReqPath, rr.Body)
			router.ServeHTTP(rr, req)
		})
	}
}

func Test_LoggerWithConfig(t *testing.T) {
	okHandler := func(c *gin.Context) {
		c.String(http.StatusOK, "200 OK")
	}
	unauthorizedHandler := func(c *gin.Context) {
		c.AbortWithError(http.StatusUnauthorized, errors.New("401 Unauthorized"))
	}

	tests := []struct {
		name         string
		givenConf    middleware.LoggerConfig
		givenReqPath string
		mockLogger   func(logger *mocks.Logger)
	}{
		{
			name: "request a skipped path",
			givenConf: middleware.LoggerConfig{
				SkipPaths: []string{"/ok"},
			},
			givenReqPath: "/ok",
			mockLogger: func(logger *mocks.Logger) {
				logger.EXPECT().With(gomock.Any())
				logger.EXPECT().
					Info(gomock.Any(), gomock.Any()).
					Times(0) // no expected call to Info
				logger.EXPECT().
					Error(gomock.Any(), gomock.Any()).
					Times(0) // no expected call to Error
			},
		},
		{
			name:         "request ok path and log an info entry",
			givenConf:    middleware.LoggerConfig{},
			givenReqPath: "/ok",
			mockLogger: func(logger *mocks.Logger) {
				logger.EXPECT().With(gomock.Any())
				logger.EXPECT().
					Info("[rebar] /ok", gomock.Any()).
					Times(1)
				logger.EXPECT().
					Error(gomock.Any(), gomock.Any()).
					Times(0) // no expected call to Error
			},
		},
		{
			name:         "request unauthorized path and log an error entry",
			givenConf:    middleware.LoggerConfig{},
			givenReqPath: "/unauthorized",
			mockLogger: func(logger *mocks.Logger) {
				logger.EXPECT().With(gomock.Any())
				logger.EXPECT().
					Error("[rebar] /unauthorized", gomock.Any()).
					Times(1)
				logger.EXPECT().
					Info(gomock.Any(), gomock.Any()).
					Times(0) // no expected call to Info
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := mocks.NewLogger(ctrl)
			tc.mockLogger(logger)

			router := gin.New()
			router.Use(middleware.LoggerWithConfig(logger, tc.givenConf))
			router.GET("/ok", okHandler)
			router.GET("/unauthorized", unauthorizedHandler)

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.givenReqPath, rr.Body)
			router.ServeHTTP(rr, req)
		})
	}
}
