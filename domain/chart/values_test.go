package chart

import (
	"github.com/go-test/deep"
	_ "swinch/testing"
	"testing"
)

func TestLoadValuesFile(t *testing.T) {
	v := Values{}
	result := v.loadValuesFile("test/charts/test_values", "test/values/values1.yaml,test/values/values2.yaml", false)
	values := Values{
		Values: map[interface{}]interface{}{
			"test": map[string]interface{}{
				"default_values": true, "success": true, "values_1": true, "values_2": true, "list": []interface{}{2, 3, 4},
			},
		},
	}
	if diff := deep.Equal(result, values); diff != nil {
		t.Error(diff)
	}
}
