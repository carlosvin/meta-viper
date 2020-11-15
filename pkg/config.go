package config

import (
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var cfgDirNames = [...]string{"config", "configs", "cfg"}

const flagName = "config"
const flagDirs = "config-dirs"

// Cfg Configuration loader
type Cfg interface {
	Load() error
}

type cfgLoader struct {
	cfg interface{}
}

// New The the configuration from files, env and flags
func New(cfg interface{}, args []string) Cfg {
	a := &appConfig{
		v:          viper.New(),
		cfg:        cfg,
		t:          reflect.ValueOf(cfg).Elem(),
		searchDirs: getSearchDirs(args[0]),
	}
	a.initEnv()
	a.initFlags(args)
	a.initFiles()
	a.loadValues()
	return a
}

func getBaseSearchDirs(program string) []string {
	dir, err := filepath.Abs(filepath.Dir(program))
	if err != nil {
		log.Printf("can't find the directory of %s, using current working path", program)
		return []string{"."}
	}
	return []string{dir, "."}
}

func getSearchDirs(program string) []string {
	dirs := make([]string, 0)
	for _, b := range getBaseSearchDirs(program) {
		dirs = append(dirs, b)
		for _, s := range cfgDirNames {
			dirs = append(dirs, filepath.Join(b, s))
		}
	}
	return dirs
}

// Load The the configuration from files, env and flags
func (a *appConfig) Load() error {
	a.loadValues()
	return nil
}

// appConfig Application configuration
type appConfig struct {
	v          *viper.Viper
	cfg        interface{}
	t          reflect.Value
	searchDirs []string
}

func (a *appConfig) initEnv() {
	// Env config
	a.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	a.v.AutomaticEnv()
}

func (a *appConfig) initFlags(args []string) {
	// Flags config
	flagSet := pflag.NewFlagSet("flagsConfig", pflag.ExitOnError)
	flagSet.String(flagName, "", "Configuration name")
	flagSet.StringSlice(flagDirs, a.searchDirs, "Configuration directories search paths")
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
	case reflect.Bool:
		flagSet.Bool(name, v.Bool(), desc)
	case reflect.Slice:
		a.initFlagSlice(v, t, flagSet, name, desc)
	default:
		panicType(t.Type)
	}
}

func (a *appConfig) initFlagSlice(v reflect.Value, t reflect.StructField, flagSet *pflag.FlagSet, name, desc string) {
	switch t.Type.Elem().Kind() {
	case reflect.String:
		slice, ok := v.Interface().([]string)
		if !ok {
			panicType(t.Type)
		}
		flagSet.StringSlice(name, slice, desc)
	case reflect.Int:
		slice, ok := v.Interface().([]int)
		if !ok {
			panicType(t.Type)
		}
		flagSet.IntSlice(name, slice, desc)
	default:
		panicType(t.Type)
	}
}

func panicType(t reflect.Type) {
	panic(fmt.Sprintf("Unexpected type %d %v", t.Kind(), t.Name()))
}

func (a *appConfig) initFiles() {
	// Config files
	configName := a.v.GetString("config")
	if configName == "" {
		log.Println("No configuration name has been specified, so no configuration file will be loaded. Using flags and environment variables.")
		return
	}
	a.v.SetConfigName(configName)
	searchDirs := a.v.GetStringSlice(flagDirs)
	for _, d := range searchDirs {
		a.v.AddConfigPath(d)
	}
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
	case reflect.Bool:
		v.SetBool(a.v.GetBool(name))
	case reflect.Slice:
		a.loadValueSlice(v, t, name)
	default:
		panicType(t.Type)
		// TODO check for nested struct and call recursively
	}
}

func (a *appConfig) loadValueSlice(v reflect.Value, t reflect.StructField, name string) {
	switch t.Type.Elem().Kind() {
	case reflect.String:
		v.Set(reflect.ValueOf(a.v.GetStringSlice(name)))
	case reflect.Int:
		v.Set(reflect.ValueOf(a.getIntSlice(name)))
	default:
		panicType(t.Type)
	}
}

func (a *appConfig) getIntSlice(name string) []int {
	ints := a.v.GetIntSlice(name)
	// for some reason when it is an env, the GetIntSlice doesn't return the slice
	if len(ints) == 0 {
		strs := a.v.GetStringSlice(name)
		ints = make([]int, len(strs))
		for i := range strs {
			ints[i], _ = strconv.Atoi(strs[i])
		}
	}
	return ints
}
