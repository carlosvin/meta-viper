package config

import (
	"fmt"
	"math"
	"os"
	"testing"

	config "github.com/carlosvin/meta-viper"
	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	A   string   `cfg_name:"a"`
	B   int      `cfg_name:"b" cfg_desc:"B is a flag for someting"`
	B64 int64    `cfg_name:"b64" cfg_desc:"B 64 value for something"`
	C   float64  `cfg_name:"c" cfg_desc:"C is a config value for something"`
	D   []string `cfg_name:"d"`
	E   []int    `cfg_name:"e" cfg_desc:"Description of list of integers"`
	F   bool     `cfg_name:"f" cfg_desc:"a boolean"`
}

func TestEnv(t *testing.T) {

	tc := &TestConfig{
		A: "default",
		B: 123,
		C: 3.1415926,
		D: []string{"a", "b"},
		E: []int{-1, 1}}
	cfg, err := config.New(tc, []string{})

	assert.NoError(t, err)

	os.Setenv("B", "999")
	os.Setenv("E", "1 2 3")

	assert.NoError(t, cfg.Reload())

	assert.Equal(t, "default", tc.A)
	assert.Equal(t, 999, tc.B)
	assert.Equal(t, int64(0), tc.B64)
	assert.Equal(t, 3.1415926, tc.C)
	assert.Equal(t, []string{"a", "b"}, tc.D)
	assert.Equal(t, []int{1, 2, 3}, tc.E)
	assert.Equal(t, false, tc.F)

	os.Setenv("A", "new value")
	os.Setenv("B64", fmt.Sprintf("%d", math.MaxInt64))
	os.Setenv("C", "0.99999")
	os.Setenv("D", "1 2 3")
	os.Setenv("F", "1")

	assert.NoError(t, cfg.Reload())
	assert.Equal(t, "new value", tc.A)
	assert.Equal(t, 999, tc.B)
	assert.Equal(t, int64(math.MaxInt64), tc.B64)
	assert.Equal(t, 0.99999, tc.C)
	assert.Equal(t, []string{"1", "2", "3"}, tc.D)
	assert.Equal(t, []int{1, 2, 3}, tc.E)
	assert.Equal(t, true, tc.F)
	cleanupEnv()
}

func TestFlags(t *testing.T) {
	args := []string{
		"test-flagss",
		"--a=imaflag",
		fmt.Sprintf("--b64=%d", math.MinInt64),
		"--d=1,2,3",
		"--e=1,2,3",
	}

	tc := &TestConfig{
		A:   "default",
		B:   123,
		B64: 1,
		C:   1.1111,
		D:   []string{"a", "b"},
		E:   []int{9, 9, 9}}
	_, err := config.New(tc, args)
	assert.NoError(t, err)

	assert.Equal(t, "imaflag", tc.A)
	assert.Equal(t, 123, tc.B)
	assert.Equal(t, int64(math.MinInt64), tc.B64)
	assert.Equal(t, 1.1111, tc.C)
	assert.Equal(t, []string{"1", "2", "3"}, tc.D)
	assert.Equal(t, []int{1, 2, 3}, tc.E)
}

func TestFiles(t *testing.T) {

	tc := &TestConfig{A: "default", B: 123, B64: 1, C: 1.1111, D: []string{"a", "b"}}
	_, err := config.New(tc, []string{"--config=test"})
	assert.NoError(t, err)

	assert.Equal(t, "from json", tc.A)
	assert.Equal(t, -2, tc.B)
	assert.Equal(t, int64(1000000000000000), tc.B64)
	assert.Equal(t, 0.001, tc.C)
	assert.Equal(t, []string{"10", "20", "30"}, tc.D)
	assert.Equal(t, []int{10, 20, 30}, tc.E)
}

func cleanupEnv() {
	os.Setenv("A", "")
	os.Setenv("B64", "")
	os.Setenv("B", "")
	os.Setenv("C", "")
	os.Setenv("D", "")
	os.Setenv("E", "")
}

type invalidNestedConfig struct {
	A      string
	Nested TestConfig `cfg_name:"nested" cfg_desc:"nested config does not work"`
}

func TestInvalidNestedFiles(t *testing.T) {
	tc := &invalidNestedConfig{A: "default"}
	_, err := config.New(tc, []string{".", "--config=test"})
	assert.Error(t, err)
}

func TestInvalidConfigFile(t *testing.T) {
	tc := &TestConfig{A: "default", B: 123, B64: 1, C: 1.1111, D: []string{"a", "b"}}
	_, err := config.New(tc, []string{"--config-dirs=./invalid", "--config=test"})
	assert.Error(t, err)
}

type invalidConfigSliceType struct {
	A       string
	Invalid []float32 `cfg_name:"invalid_float" cfg_desc:"float slice"`
}

func TestInvalidSliceType(t *testing.T) {
	tc := &invalidConfigSliceType{
		A:       "default",
		Invalid: []float32{1.1, 2.2}}
	_, err := config.New(tc, []string{
		"test",
		"--config=test"})
	assert.Error(t, err)
}
