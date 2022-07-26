package server

import (
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/dinorain/kalobranded/docs"
)

const (
	maxHeaderBytes = 1 << 20
	readTimeout    = 15 * time.Second
	writeTimeout   = 15 * time.Second
)

func (s *Server) runHttpServer() error {
	s.mapRoutes()

	s.httpS = &http.Server{
		Addr:           s.cfg.Http.Port,
		Handler:        s.mw.RequestLoggerMiddleware(s.mux),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}
	s.logger.Infof("Starting http server on %s", s.cfg.Http.Port)
	return s.httpS.ListenAndServe()
}

func (s *Server) mapRoutes() {
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "API Gateway"
	docs.SwaggerInfo.Description = "API Gateway Kalobranded."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"

	s.mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
}
