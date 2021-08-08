// Package enzap implements the function to create *zap.Logger.
package enzap

import (
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	EnvCaller      = "ZAP_CALLER"
	EnvDColor      = "ZAP_DCOLOR"
	EnvDevelopment = "ZAP_DEVELOPMENT"
	EnvLevel       = "ZAP_LEVEL"
	EnvStackTrace  = "ZAP_STACK_TRACE"
	EnvTimeLayout  = "ZAP_TIME_LAYOUT"
)

const (
	TimeLayoutDevelopment = `02T15:04:05.000000Z0700`
	TimeLayoutProduction  = time.RFC3339Nano
)

// ReplaceGlobals calls zap.ReplaceGlobals(enzap.New()).
func ReplaceGlobals() func() {
	return zap.ReplaceGlobals(New())
}

// New returns *zap.Logger.
func New() *zap.Logger {
	return NewConfig().Build()
}

// NewConfig returns *Config.
func NewConfig() *Config {
	development := lookupBool(EnvDevelopment, false)

	config := &Config{
		Caller:      true,
		DColor:      true,
		Development: development,
		StackTrace:  zapcore.ErrorLevel,
	}

	if development {
		config.Level = zapcore.DebugLevel
		config.TimeLayout = TimeLayoutDevelopment
	} else {
		config.Level = zapcore.InfoLevel
		config.TimeLayout = TimeLayoutProduction
	}

	config.Caller = lookupBool(EnvCaller, config.Caller)
	config.DColor = lookupBool(EnvDColor, config.DColor)
	config.Level = lookupLevel(EnvLevel, config.Level)
	config.StackTrace = lookupLevel(EnvStackTrace, config.StackTrace)
	config.TimeLayout = lookupString(EnvTimeLayout, config.TimeLayout)

	return config
}

func lookupBool(key string, defaultValue bool) bool {
	if v, found := os.LookupEnv(key); found {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return defaultValue
}

func lookupString(key string, defaultValue string) string {
	if v, found := os.LookupEnv(key); found {
		return v
	}
	return defaultValue
}

func lookupLevel(key string, defaultValue zapcore.Level) zapcore.Level {
	if v, found := os.LookupEnv(key); found {
		_ = defaultValue.Set(v)
	}
	return defaultValue
}

type Config struct {
	Caller      bool
	DColor      bool
	Development bool
	Level       zapcore.Level
	StackTrace  zapcore.Level
	TimeLayout  string
}

func (config *Config) Build() *zap.Logger {
	encoder := config.newEncoder()

	stderr := zapcore.Lock(os.Stderr)
	stderrLevel := config.stderrLevelEnabler()

	stdout := zapcore.Lock(os.Stdout)
	stdoutLevel := config.stdoutLevelEnabler()

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, stderr, stderrLevel),
		zapcore.NewCore(encoder, stdout, stdoutLevel),
	)

	opts := []zap.Option{
		zap.WithCaller(config.Caller),
		zap.AddStacktrace(config.StackTrace),
	}
	if config.Development {
		opts = append(opts, zap.Development())
	}

	return zap.New(core, opts...)
}

func (config *Config) newEncoder() zapcore.Encoder {
	if config.Development {
		zapcfg := zap.NewDevelopmentEncoderConfig()
		zapcfg.EncodeTime = zapcore.TimeEncoderOfLayout(config.TimeLayout)
		if config.DColor {
			zapcfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		} else {
			zapcfg.EncodeLevel = zapcore.CapitalLevelEncoder
		}
		return zapcore.NewConsoleEncoder(zapcfg)
	}

	zapcfg := zap.NewProductionEncoderConfig()
	zapcfg.EncodeTime = zapcore.TimeEncoderOfLayout(config.TimeLayout)
	zapcfg.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(zapcfg)
}

func (config *Config) stderrLevelEnabler() zapcore.LevelEnabler {
	level := config.Level
	if level < zapcore.ErrorLevel {
		level = zapcore.ErrorLevel
	}
	return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool { return lvl >= level })
}

func (config *Config) stdoutLevelEnabler() zapcore.LevelEnabler {
	level := config.Level
	switch {
	case level <= zapcore.DebugLevel:
		return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.ErrorLevel
		})

	case level >= zapcore.ErrorLevel:
		return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return false
		})

	default:
		return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return level <= lvl && lvl < zapcore.ErrorLevel
		})
	}
}
