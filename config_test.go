package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComments(t *testing.T) {

	const FILE_NAME = "./config.json"
	defer func() {
		require.NoError(t, os.Remove(FILE_NAME))
	}()

	for _, src := range []string{
		`{
	        "param_str1":"text1",
	        "param_str2":"text2",
	        "param_num":50
		}`,
		`{
			// String parameter
	        "param_str1":"text1",
	        "param_str2":"text2", // String parameter

	        "param_num":50,
		}`,
	} {
		require.NoError(t, ioutil.WriteFile("./config.json", []byte(src), os.ModePerm))

		conf, err := NewConfig()
		require.NoError(t, err)
		require.NoError(t, conf.ReadConfig())
		require.Equal(t, "text1", conf.Str("param_str1"))
		require.Equal(t, "text2", conf.Str("param_str2"))
		require.Equal(t, 50, conf.Int("param_num"))
	}
}

func TestRead(t *testing.T) {

	data := `{
        "param_str":"text",
        "param_num":50,
        "param_bool":true,
        "param_array_int":[1,2,3],
        "param_array_str":["1","2","3"],
        "param_map_str":{"k1":"v1","k2":"v2", "k3":"v3"}
}`

	const FILE_NAME = "./config.json"
	require.NoError(t, ioutil.WriteFile(FILE_NAME, []byte(data), os.ModePerm))
	defer func() {
		require.NoError(t, os.Remove(FILE_NAME))
	}()

	conf, err := NewConfig()
	require.NoError(t, err)
	require.NoError(t, conf.ReadConfig())
	require.Equal(t, "text", conf.Str("param_str"))
	require.Equal(t, 50, conf.Int("param_num"))
	require.Equal(t, float64(50), conf.Float64("param_num"))
	require.Equal(t, int64(50), conf.Int64("param_num"))
	require.Equal(t, true, conf.Bool("param_bool"))

	{
		require.Equal(t,
			[]interface{}{float64(1), float64(2), float64(3)},
			conf.Array("param_array_int"))

		require.Equal(t,
			[]string{},
			conf.ArrayStr("param_array_int"))
	}

	{
		require.Equal(t,
			[]interface{}{"1", "2", "3"},
			conf.Array("param_array_str"))

		require.Equal(t,
			[]string{"1", "2", "3"},
			conf.ArrayStr("param_array_str"))
	}

	{
		require.Equal(t,
			map[string]interface{}{
				"k1": "v1",
				"k2": "v2",
				"k3": "v3",
			},
			conf.Map("param_map_str"))

		require.Equal(t,
			map[string]string{
				"k1": "v1",
				"k2": "v2",
				"k3": "v3",
			},
			conf.MapStr("param_map_str"))
	}
}
