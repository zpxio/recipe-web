/*
 * Copyright 2022 zpxio (Jeff Sharpe)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.org/zpxio/recipe-web/pkg/config"
	"github.org/zpxio/recipe-web/pkg/server/middleware"
	"github.org/zpxio/recipe-web/pkg/server/page"
	"net/http"
	"time"
)

type Server struct {
	eInst  *echo.Echo
	Config config.Config
}

func CreateServer(config config.Config) *Server {
	// Setup
	e := echo.New()
	// Slightly more casual logging
	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level}")
	}
	e.Logger.SetLevel(99)

	s := &Server{
		eInst:  e,
		Config: config,
	}

	s.initGlobalMiddleware()
	s.initStandardEndpoints()
	s.initDiagnosticEndpoints()

	return s
}

func (s *Server) Logger() echo.Logger {
	return s.eInst.Logger
}

func (s *Server) Start() error {
	go func() {
		addrString := fmt.Sprintf("%s:%d", s.Config.Server.LocalIP, s.Config.Server.Port)

		if err := s.eInst.Start(addrString); err != nil && err != http.ErrServerClosed {
			s.Logger().Fatal("shutting down the server")
		}
	}()

	return nil
}

func (s *Server) Shutdown(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := s.eInst.Shutdown(ctx); err != nil {
		s.Logger().Fatal(err)
	}
}

func (s *Server) initGlobalMiddleware() {
	s.eInst.Pre(middleware.RequestTime)
	s.eInst.Pre(middleware.RequestID)
	s.eInst.Pre(middleware.RequestLogger)
}

func (s *Server) initStandardEndpoints() {
	s.eInst.GET("/", page.Index)
}

func (s *Server) initDiagnosticEndpoints() {
	s.eInst.GET("/ping", PingHandler)
}

func (s *Server) SetLogLevel(l log.Lvl) {
	s.eInst.Logger.SetLevel(l)
}
