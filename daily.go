package reporter

import (
	"fmt"
	"github.com/adminium/logger"
	"github.com/robfig/cron/v3"
	"time"
)

var log = logger.NewLogger("daily")

type Config struct {
	runCron    bool
	cron       *cron.Cron
	expression string
	generator  func() (report string, err error)
	channels   []*Channel
}

type Channel struct {
	provider Provider
	url      string
}

type Option func(*Config)

func WithCron(cron *cron.Cron) Option {
	return func(config *Config) {
		if cron != nil {
			config.cron = cron
			config.runCron = true
		}
	}
}

func WithGenerator(generator func() (report string, err error)) Option {
	return func(config *Config) {
		if generator != nil {
			config.generator = generator
		}
	}
}

func WithExpression(expression string) Option {
	return func(config *Config) {
		if expression != "" {
			config.expression = expression
		}
	}
}

func WithChannel(provider Provider, url string) Option {
	return func(config *Config) {
		config.channels = append(config.channels, &Channel{
			provider: provider,
			url:      url,
		})
	}
}

func nowDate() string {
	return time.Now().Format("2006-01-02")
}

func NewDaily(options ...Option) *Daily {
	conf := &Config{
		runCron:    true,
		cron:       cron.New(),
		expression: "0 23 * * ?",
		generator: func() (string, error) {
			return fmt.Sprintf("report: %s\nhello!", nowDate()), nil
		},
	}
	for _, option := range options {
		option(conf)
	}
	if conf.runCron {
		go conf.cron.Run()
	}
	return &Daily{
		conf: conf,
	}
}

type Daily struct {
	conf *Config
}

func (d *Daily) Start() (err error) {
	_, err = d.conf.cron.AddFunc(d.conf.expression, d.report)
	if err != nil {
		return
	}
	if d.conf.runCron {
		d.conf.cron.Run()
	}
	return
}

func (d *Daily) Stop() {
	if d.conf.runCron {
		d.conf.cron.Stop()
	}
}

func (d *Daily) report() {
	report, err := d.conf.generator()
	if err != nil {
		report = fmt.Sprintf("prepare report %s error: %s", nowDate(), err)
	}
	for _, v := range d.conf.channels {
		err = Report(v.provider, v.url, report)
		if err != nil {
			log.Errorf("report error: %s", err)
		} else {
			log.Infof("report via channel %s success", v.provider)
		}
	}
}
