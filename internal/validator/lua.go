package validator

import (
	"fmt"
	"time"

	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

// MapToTable convert map[strign]interface{} into lua table
// Adopted from here https://github.com/yuin/gopher-lua/issues/160
func MapToTable(m map[string]interface{}) *lua.LTable {
	// Main table pointer
	resultTable := &lua.LTable{}

	// Loop map
	for key, element := range m {
		switch element.(type) {
		case float64:
			resultTable.RawSetString(key, lua.LNumber(element.(float64)))
		case int64:
			resultTable.RawSetString(key, lua.LNumber(element.(int64)))
		case string:
			resultTable.RawSetString(key, lua.LString(element.(string)))
		case bool:
			resultTable.RawSetString(key, lua.LBool(element.(bool)))
		case []byte:
			resultTable.RawSetString(key, lua.LString(string(element.([]byte))))
		case map[string]interface{}:
			// Get table from map
			tble := MapToTable(element.(map[string]interface{}))
			resultTable.RawSetString(key, tble)
		case time.Time:
			resultTable.RawSetString(key, lua.LNumber(element.(time.Time).Unix()))
		case []map[string]interface{}:
			// Create slice table
			sliceTable := &lua.LTable{}
			// Loop element
			for _, s := range element.([]map[string]interface{}) {
				// Get table from map
				tble := MapToTable(s)
				sliceTable.Append(tble)
			}
			// Set slice table
			resultTable.RawSetString(key, sliceTable)
		case []interface{}:
			// Create slice table
			sliceTable := &lua.LTable{}
			// Loop interface slice
			for _, s := range element.([]interface{}) {
				// Switch interface type
				switch s.(type) {
				case map[string]interface{}:
					// Convert map to table
					t := MapToTable(s.(map[string]interface{}))
					// Append result
					sliceTable.Append(t)
				case float64:
					// Append result as number
					sliceTable.Append(lua.LNumber(s.(float64)))
				case string:
					// Append result as string
					sliceTable.Append(lua.LString(s.(string)))
				case bool:
					// Append result as bool
					sliceTable.Append(lua.LBool(s.(bool)))
				}
			}
			// Append to main table
			resultTable.RawSetString(key, sliceTable)
		}
	}

	return resultTable
}

type Result struct {
	value bool
}

func (r *Result) Get() bool {
	return r.value
}

func (r *Result) Set(value bool) {
	r.value = value
}

type LuaValidator struct {
	config       *config.ValidatorConfig
	luaState     *lua.LState
	result       *Result
	preparedRule string
}

type Helper struct {
}

func (h *Helper) UnixTime() float64 {
	return float64(time.Now().UTC().UnixNano()) / float64(time.Second)
}

func (h *Helper) Duration(value string) float64 {
	duration, err := time.ParseDuration(value)
	if err != nil {
		panic(fmt.Sprintf("wrong duration %s", err))
	}
	return float64(duration) / float64(time.Second)
}

func (v *LuaValidator) Exec(data map[string]interface{}) (bool, error) {
	v.luaState.SetGlobal("data", MapToTable(data))
	if err := v.luaState.DoString(v.preparedRule); err != nil {
		return false, err
	}
	return v.result.Get(), nil
}

func (v *LuaValidator) Name() string {
	return v.config.Name
}

func (v *LuaValidator) Title() string {
	return v.config.Title
}

func (v *LuaValidator) IsMajor() bool {
	return v.config.Major
}

func NewLuaValidator(vc *config.ValidatorConfig) (*LuaValidator, error) {
	l := lua.NewState()

	r := &Result{
		value: false,
	}
	l.SetGlobal("result", luar.New(l, r))
	h := &Helper{}
	l.SetGlobal("helper", luar.New(l, h))

	v := &LuaValidator{config: vc, luaState: l, result: r, preparedRule: vc.Rule}

	if v.config.Type == config.LuaEngine {
		v.preparedRule = fmt.Sprintf("result:Set(%s)", v.preparedRule)
	}

	if _, err := l.LoadString(v.preparedRule); err != nil {
		return nil, err
	}

	return v, nil
}
