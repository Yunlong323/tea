package yaml

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"tea/entity"
)

func ParseSettings(reader io.Reader) map[string]interface{} {
	viper.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	err := viper.ReadConfig(reader) // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		log.WithError(err).Panic("fatal error config file")
	}
	return viper.AllSettings()
}

func ParseTemplate(reader io.Reader) *entity.Template {
	settings := ParseSettings(reader)
	if settings["kind"] != "template" {
		log.WithError(fmt.Errorf("get kind %v instead of template", settings["kind"])).Panic("read config error")
	}
	template := entity.TemplateEntity{}
	err := mapstructure.Decode(settings, &template)
	if err != nil {
		log.WithError(err).Panic("failed to parse template")
	}
	return entity.NewTemplate(&template)
}

func ParsePipeline(reader io.Reader) *entity.Pipeline {
	settings := ParseSettings(reader)
	if settings["kind"] != "pipeline" {
		log.WithError(fmt.Errorf("get kind %v instead of pipeline", settings["kind"])).Panic("read config error")
	}
	pipeline := entity.PipelineEntity{}
	err := mapstructure.Decode(settings, &pipeline)
	if err != nil {
		log.WithError(err).Panic("failed to parse pipeline")
	}
	return entity.NewPipeline(&pipeline)
}
