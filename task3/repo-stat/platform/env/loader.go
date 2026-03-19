package env

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

func MustLoad(path string, cfg interface{}) {
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		log.Fatalf("cannot read config %q: %s", path, err)
	}
}
