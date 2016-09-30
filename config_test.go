package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestRead(t *testing.T) {

	data := `{
        "param_str":"text",
        "param_num":50,
        "param_bool":true,
        "param_array_int":[1,2,3],
        "param_array_str":["1","2","3"],
        "param_map_str":{"k1":"v1","k2":"v2", "k3":"v3"}
}`

	ioutil.WriteFile("./config.json", []byte(data), os.ModePerm)

	conf, err := NewConfig()
	if err != nil {
		t.Error(err.Error())
	}

	err = conf.ReadConfig()

	if err != nil {
		t.Error(err.Error())
	}

	{
		val := conf.Str("param_str")
		etalon := "text"
		if val != etalon {
			t.Errorf("'%s' != '%s'", val, etalon)
		}
	}

	{
		val := conf.Int("param_num")
		etalon := 50
		if val != etalon {
			t.Errorf("%#v != %v", val, etalon)
		}
	}

	{
		val := conf.Float64("param_num")
		etalon := 50.0
		if val != etalon {
			t.Errorf("%#v != %v", val, etalon)
		}
	}

	{
		val := conf.Int64("param_num")
		etalon := int64(50)
		if val != etalon {
			t.Errorf("%#v != %v", val, etalon)
		}
	}

	{
		val := conf.Bool("param_bool")
		etalon := true
		if val != etalon {
			t.Errorf("%#v != %v", val, etalon)
		}
	}

	{
		val := conf.Array("param_array_int")
		etalon := []interface{}{1, 2, 3}
		if len(val) != len(etalon) {
			t.Errorf("%#v != %#v", val, etalon)
		}

		for i, vEtalon := range etalon {
			vRetval := val[i]
			if vEtalon.(int) != int(vRetval.(float64)) {
				t.Errorf("%#v != %#v", vEtalon, vRetval)
			}
		}
	}

	{
		val := conf.Array("param_array_str")
		etalon := []interface{}{"1", "2", "3"}
		if len(val) != len(etalon) {
			t.Errorf("%#v != %#v", val, etalon)
		}

		for i, vEtalon := range etalon {
			vRetval := val[i]
			if vEtalon.(string) != vRetval.(string) {
				t.Errorf("%#v != %#v", vEtalon, vRetval)
			}
		}
	}

	{
		val := conf.Map("param_map_str")
		etalon := map[string]interface{}{"k1": "v1", "k2": "v2", "k3": "v3"}
		if len(val) != len(etalon) {
			t.Errorf("%#v != %#v", val, etalon)
		}

		for kEtalon, vEtalon := range etalon {
			vRetval := val[kEtalon]
			if vEtalon.(string) != vRetval.(string) {
				t.Errorf("%s: %#v != %#v", kEtalon, vEtalon, vRetval)
			}
		}
	}

	{
		val := conf.MapStr("param_map_str")
		etalon := map[string]string{"k1": "v1", "k2": "v2", "k3": "v3"}
		if len(val) != len(etalon) {
			t.Errorf("%#v != %#v", val, etalon)
		}

		for kEtalon, vEtalon := range etalon {
			vRetval := val[kEtalon]
			if vEtalon != vRetval {
				t.Errorf("%s: %#v != %#v", kEtalon, vEtalon, vRetval)
			}
		}
	}
}
