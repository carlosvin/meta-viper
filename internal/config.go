package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Cfg interface {
	Load() error
}

type cfgLoader struct {
	cfg interface{}
}

// New The the configuration from files, env and flags
func New(cfg interface{}, args []string) Cfg {
	a := &appConfig{
		v:   viper.New(),
		cfg: cfg,
		t:   reflect.ValueOf(cfg).Elem()}
	a.initEnv()
	a.initFlags(args)
	a.initFiles()
	a.loadValues()
	return a
}

// Load The the configuration from files, env and flags
func (a *appConfig) Load() error {
	a.loadValues()
	return nil
}

// appConfig Application configuration
type appConfig struct {
	v   *viper.Viper
	cfg interface{}
	t   reflect.Value
}

func (a *appConfig) initEnv() {
	// Env config
	a.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	a.v.AutomaticEnv()
}

func (a *appConfig) initFlags(args []string) {
	// Flags config
	flagSet := pflag.NewFlagSet("flagsConfig", pflag.ContinueOnError)
	flagSet.String("config", "prod", "Configuration name")
	for i := 0; i < a.t.NumField(); i++ {
		a.initFlag(a.t.Field(i), a.t.Type().Field(i), flagSet)
	}
	flagSet.Parse(args)
	a.v.BindPFlags(flagSet)
}

func (a *appConfig) initFlag(v reflect.Value, t reflect.StructField, flagSet *pflag.FlagSet) {
	name, desc := t.Tag.Get("cfg_name"), t.Tag.Get("cfg_desc")
	switch t.Type.Kind() {
	case reflect.String:
		flagSet.String(name, v.String(), desc)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		flagSet.Int64(name, v.Int(), desc)
	case reflect.Float32, reflect.Float64:
		flagSet.Float64(name, v.Float(), desc)
	// TODO check for nested struct and call recursively
	default:
		panic("Unexpected type " + t.Name)
	}
}

func (a *appConfig) initFiles() {
	// Config files
	configName := a.v.GetString("config")
	a.v.SetConfigName(configName)
	a.v.AddConfigPath("configs")
	a.v.AddConfigPath("../configs")
	a.v.WatchConfig()
	a.v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	err := a.v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func (a *appConfig) loadValues() {
	for i := 0; i < a.t.NumField(); i++ {
		a.loadValue(a.t.Field(i), a.t.Type().Field(i))
	}
}

func (a *appConfig) loadValue(v reflect.Value, t reflect.StructField) {
	name := t.Tag.Get("cfg_name")
	switch t.Type.Kind() {
	case reflect.String:
		v.SetString(a.v.GetString(name))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(a.v.GetInt64(name))
	case reflect.Float32, reflect.Float64:
		v.SetFloat(a.v.GetFloat64(name))
	// TODO check for nested struct and call recursively
	default:
		panic("Unexpected type " + t.Name)
	}
}
