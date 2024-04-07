package abstract

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Snow-Gremlin/goToolbox/utils"
)

type (
	jsonData any
	jsonMap  map[string]jsonData
	jsonList []jsonData
)

func (j *jsonList) add(d jsonData) *jsonList {
	*j = append(*j, d)
	return j
}

func (j jsonMap) add(name string, d jsonData) jsonMap {
	j[name] = d
	return j
}

func isEmpty(d jsonData) bool {
	if utils.IsZero(d) {
		return true
	}
	len, _ := utils.Length(d)
	return len == 0
}

func (j jsonMap) addNotEmpty(name string, d jsonData) jsonMap {
	if !isEmpty(d) {
		j.add(name, d)
	}
	return j
}

func (j jsonMap) append(name string, d jsonData) jsonMap {
	if v, has := j[name]; has {
		if s, ok := v.(jsonList); ok {
			j[name] = append(s, d)
		} else {
			j[name] = jsonList{s, d}
		}
	} else {
		j[name] = jsonList{d}
	}
	return j
}

func jsonMarshal(minimize bool, data jsonData) ([]byte, error) {
	if minimize {
		return json.Marshal(data)
	}
	return json.MarshalIndent(data, ``, `  `)
}

func writeJson(path string, minimize bool, data jsonData) error {
	b, err := jsonMarshal(minimize, data)
	if err != nil {
		return err
	}

	if len(path) > 0 {
		return os.WriteFile(path, b, 0666)
	}

	_, err = fmt.Println(string(b))
	return err
}
