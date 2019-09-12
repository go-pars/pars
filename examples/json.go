package examples

import (
	"io"

	"github.com/ktnyt/pars"
)

// JSON parser parts.
var (
	Value  pars.Parser
	Null   = pars.String("null").Bind(nil)
	True   = pars.String("true").Bind(true)
	False  = pars.String("false").Bind(false)
	String = pars.Quoted('"')
	Number = pars.Number.Map(pars.ParseFloat(64))
	Array  = pars.Phrase('[', pars.Cut, pars.Sep(&Value, ','), ']').Map(func(result *pars.Result) error {
		c := result.Children[2].Children
		v := make([]interface{}, len(c))
		for i := range c {
			v[i] = c[i].Value
		}
		result.Value = v
		result.Children = nil
		return nil
	})
	Object = pars.Phrase('{', pars.Cut, pars.Sep(pars.Phrase(String, ":", &Value), ','), '}').Map(func(result *pars.Result) error {
		c := result.Children[2].Children
		v := make(map[string]interface{})
		for i := range c {
			v[c[i].Children[0].Value.(string)] = c[i].Children[2].Value
		}
		result.Value = v
		result.Children = nil
		return nil
	})
)

func init() {
	Value = pars.Any(Null, True, False, String, Number, Array, Object)
}

// Unmarshal json string into an interface{}.
func Unmarshal(r io.Reader) (interface{}, error) {
	return pars.Apply(Value, pars.NewState(r))
}
