package chart

import (
	"github.com/go-test/deep"
	"os"
	"path"
	"swinch/domain/datastore"
	_ "swinch/testing"
	"testing"
)

type renderTest struct {
	chartPath            string
	valuesFile           string
	outputPath           string
	fullRender           bool
	excludeDefaultValues bool
	control              string
}

var simpleRender = renderTest{
	"test/charts/test_template",
	"",
	"test_template_simple",
	false,
	false,
	"test/manifests/test_template_simple/pipeline.yaml",
}

var optionsRender = renderTest{
	"test/charts/test_template",
	"test/values/test_template_options.yaml",
	"test_template_options",
	false,
	true,
	"test/manifests/test_template_options/pipeline.yaml",
}

var fullRender = renderTest{
	"test/charts/test_template",
	"test/values/test_template_full_render.yaml",
	"test_template_full_render",
	true,
	false,
	"test/manifests/test_template_full_render/pipeline.yaml",
}

func TestTemplateChart(t *testing.T) {
	r := renderTest{}
	r.runRenderTest(simpleRender, t)
	r.runRenderTest(optionsRender, t)
	r.runRenderTest(fullRender, t)
}

func (r renderTest) runRenderTest(test renderTest, t *testing.T){
	t.Run(test.outputPath, func(t *testing.T) {
		control, render := r.renderer(test)
		if len(control) == 0 || len(render) == 0 {
			t.Error("Failed to load test params.")
		}
		if diff := deep.Equal(control, render); diff != nil {
			t.Error(diff)
		}
	})
}

func (r renderTest) renderer(test renderTest) ([]byte, []byte){
	r = test
	tp := Template{}
	d := datastore.Datastore{}

	outputPath := d.CreateTmpFolder()
	defer os.RemoveAll(outputPath)

	r.outputPath = path.Join(outputPath + r.outputPath)
	tp.TemplateChart(
		r.chartPath,
		r.valuesFile,
		r.outputPath,
		r.fullRender,
		r.excludeDefaultValues)

	control := d.ReadFile(r.control)
	render := d.ReadFile(path.Join(r.outputPath + "/pipeline.yaml"))

	return control, render
}
