package app

import (
	"github.com/osamikoyo/geass-v2/internal/parsers"
	"github.com/osamikoyo/geass-v2/pkg/config"
)

func App() error {
	cfg, err := config.Load("config.yml")
	if err != nil {
		return err
	}

	parser, err := parsers.New(&cfg)
	if err != nil{
		return err
	}

	parser.ParsePage(cfg.StartUrl, cfg.Deep)
	return nil
}