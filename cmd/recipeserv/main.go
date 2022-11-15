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

package main

import (
	"fmt"
	"github.com/apex/log"
	"github.com/mkideal/cli"
	config2 "github.org/zpxio/recipe-web/pkg/config"
	"github.org/zpxio/recipe-web/pkg/server"
	"os"
	"os/signal"
	"time"
)

type argT struct {
	cli.Helper
	ConfigPath    string `cli:"config" usage:"Server configuration file"`
	BaseDirectory string `cli:"d" usage:"Base server directory"`

	//Options

}

func main() {

	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		log.Infof("Starting up.")

		ctx.JSONln(ctx.Argv())
		argv := ctx.Argv().(*argT)

		// Load config
		log.Infof("Loading configuration file: %s", argv.ConfigPath)
		config, configErr := config2.Load(argv.ConfigPath)
		if configErr != nil {
			return fmt.Errorf("could not load config: %s", configErr)
		}

		// Setup
		srv := server.CreateServer(config)

		// Start server
		err := srv.Start()
		if err != nil {
			log.Fatalf("failed to start server: %s", err)
		}

		// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
		// Use a buffered channel to avoid missing signals as recommended for signal.Notify
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		// Execute a timed shutdown of the server
		srv.Shutdown(10 * time.Second)

		return nil
	}))
}
