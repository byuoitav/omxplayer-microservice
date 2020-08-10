package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
)

type controlConfig struct {
	Streams []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"streams"`
}

type controlTemplateData struct {
	Streams []controlTemplateDataStream
}

type controlTemplateDataStream struct {
	Name     string
	URL      string
	Selected bool
}

func (h *Handlers) ControlPageHandler(path string) echo.HandlerFunc {
	tmpl := template.Must(template.ParseFiles(path))

	return echo.HandlerFunc(func(c echo.Context) error {
		data, err := ioutil.ReadFile(h.ControlConfigPath)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("unable to read config file: %s", err))
		}

		var config controlConfig
		if err := json.Unmarshal(data, &config); err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("unable to parse config file: %s", err))
		}

		var currStreamURL string
		// TODO add back in
		//conn, err := helpers.ConnectToDbus()
		//if err == nil {
		//	currStreamURL, _ = helpers.GetStream(conn)
		//}

		var tmplData controlTemplateData
		for _, stream := range config.Streams {
			tmplData.Streams = append(tmplData.Streams, controlTemplateDataStream{
				Name:     stream.Name,
				URL:      stream.URL,
				Selected: currStreamURL == stream.URL,
			})
		}

		if err := tmpl.Execute(c.Response().Writer, tmplData); err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("unable to execute template: %s", err))
		}

		return nil
	})
}
