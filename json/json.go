package examples

import (
	"io"

	"github.com/ktnyt/pars"
)

func arrayMap(result *pars.Result) error {
	children := result.Children[1].Children
	v := make([]interface{}, len(children))
	for i, child := range children {
		v[i] = child.Value
	}
	result.SetValue(v)
	return nil
}

func objectMap(result *pars.Result) error {
	children := result.Children[1].Children
	v := make(map[string]interface{})
	for _, child := range children {
		name, value := child.Children[0].Value.(string), child.Children[2].Value
		v[name] = value
	}
	result.SetValue(v)
	return nil
}

// JSON parser parts.
var (
	Value  pars.Parser
	Null   = pars.String("null").Bind(nil)
	True   = pars.String("true").Bind(true)
	False  = pars.String("false").Bind(false)
	Number = pars.Number
	String = pars.Quoted('"').Map(pars.ToString)
	Array  = pars.Seq('[', pars.Delim(&Value, ','), ']').Map(arrayMap)
	prop   = pars.Seq(String, ':', &Value)
	Object = pars.Seq('{', pars.Delim(prop, ','), '}').Map(objectMap)
)

func init() {
	Value = pars.Any(Null, True, False, String, Number, Array, Object)
}

// Unmarshal json string into an interface{}.
func Unmarshal(r io.Reader) (interface{}, error) {
	state := pars.NewState(pars.NewReader(r))
	result, err := pars.Exact(Value).Parse(state)
	return result.Value, err
}
