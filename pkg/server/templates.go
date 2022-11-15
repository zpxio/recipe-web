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
	"fmt"
	"github.com/apex/log"
	"github.com/labstack/echo/v4"
	"github.org/zpxio/recipe-web/pkg/config"
	"html/template"
	"io"
	"path/filepath"
)

type TemplateLibrary struct {
	templates     *template.Template
	IndexTemplate *template.Template
}

func (t *TemplateLibrary) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func ParseTemplates(config config.Config) *TemplateLibrary {
	templatePath := filepath.Clean(filepath.Join(config.BaseDir, config.Content.TemplateDir))
	log.Infof("Loading templates from: %s", templatePath)

	t := &TemplateLibrary{
		templates: template.Must(template.ParseGlob(fmt.Sprintf("%s/*.html", templatePath))),
	}

	return t
}
