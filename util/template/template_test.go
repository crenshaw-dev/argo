package template

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SimpleValue struct {
	Value string `json:"value,omitempty"`
}

func processTemplate(t *testing.T, tmpl SimpleValue) SimpleValue {
	tmplBytes, err := json.Marshal(tmpl)
	assert.NoError(t, err)
	r, err := Replace(string(tmplBytes), map[string]string{}, true)
	assert.NoError(t, err)
	var newTmpl SimpleValue
	err = json.Unmarshal([]byte(r), &newTmpl)
	assert.NoError(t, err)
	return newTmpl
}

func Test_Template_Replace(t *testing.T) {
	t.Run("ExpressionWithEscapedCharacters", func(t *testing.T) {
		t.Run("SingleQuotes", func(t *testing.T) {
			tmpl := SimpleValue{Value: "{{='test'}}"}
			newTmpl := processTemplate(t, tmpl)
			assert.Equal(t, "test", newTmpl.Value)
		})
		t.Run("DoubleQuotes", func(t *testing.T) {
			tmpl := SimpleValue{Value: `{{="test"}}`}
			newTmpl := processTemplate(t, tmpl)
			assert.Equal(t, "test", newTmpl.Value)
		})
		t.Run("EscapedBackslashInString", func(t *testing.T) {
			tmpl := SimpleValue{Value: `{{='some\\path\\with\\backslashes'}}`}
			newTmpl := processTemplate(t, tmpl)
			assert.Equal(t, `some\path\with\backslashes`, newTmpl.Value)
		})
		t.Run("EscapedNewlineInString", func(t *testing.T) {
			tmpl := SimpleValue{Value: `{{='some\nstring\nwith\nescaped\nnewlines'}}`}
			newTmpl := processTemplate(t, tmpl)
			assert.Equal(t, "some\nstring\nwith\nescaped\nnewlines", newTmpl.Value)
		})
		t.Run("Newline", func(t *testing.T) {
			tmpl := SimpleValue{Value: "{{=1 + \n1}}"}
			newTmpl := processTemplate(t, tmpl)
			assert.Equal(t, "2", newTmpl.Value)
		})
		t.Run("StringAsJson", func(t *testing.T) {
			tmpl := SimpleValue{Value: "{{=toJson('test')}}"}
			newTmpl := processTemplate(t, tmpl)
			assert.Equal(t, `"test"`, newTmpl.Value)
		})
	})
}
