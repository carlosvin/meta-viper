package config

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Load The the configuration from files, env and flags
func Load(cfg interface{}) {
	a := &appConfig{
		v:   viper.New(),
		cfg: cfg,
		t:   reflect.ValueOf(cfg).Elem()}
	a.initEnv()
	a.initFlags()
	a.initFiles()
	a.loadValues()
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
	log.Println(a.v.GetString("HOST"))
}

func (a *appConfig) initFlags() {
	// Flags config
	pflag.String("config", "prod", "Configuration name")
	for i := 0; i < a.t.NumField(); i++ {
		a.initFlag(a.t.Field(i), a.t.Type().Field(i))
	}
	pflag.Parse()
	a.v.BindPFlags(pflag.CommandLine)
}

func (a *appConfig) initFlag(v reflect.Value, t reflect.StructField) {
	name, desc := t.Tag.Get("cfg_name"), t.Tag.Get("cfg_desc")
	switch t.Type.Kind() {
	case reflect.String:
		pflag.String(name, v.String(), desc)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		pflag.Int64(name, v.Int(), desc)
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
		/*
			v := a.t.Field(i)
			name := a.t.Type().Field(i).Tag.Get("cfg_name")
			newVal := a.v.Get(name)
			if newVal == nil {
				log.Printf("No config found for '%s'\n", name)
			} else {
				v.Set(reflect.ValueOf(newVal))
			}*/
	}
}

func (a *appConfig) loadValue(v reflect.Value, t reflect.StructField) {
	name := t.Tag.Get("cfg_name")
	switch t.Type.Kind() {
	case reflect.String:
		v.SetString(a.v.GetString(name))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(a.v.GetInt64(name))
	// TODO check for nested struct and call recursively
	default:
		panic("Unexpected type " + t.Name)
	}
}
