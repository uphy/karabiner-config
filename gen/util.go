package gen

import (
	"encoding/json"
	"io"
	"log"
)

type JSONObject map[string]interface{}
type JSONArray []interface{}

func ToObject(v interface{}) JSONObject {
	b, err := json.Marshal(v)
	if err != nil {
		log.Fatal("failed on toObject:", err)
	}
	var obj JSONObject
	if err := json.Unmarshal(b, &obj); err != nil {
		log.Fatal("failed to unmarshal in toObject:", err)
	}
	return obj
}

func (o JSONObject) InjectMembersTo(dest JSONObject) {
	for name, value := range o {
		dest[name] = value
	}
}

func (o JSONObject) Write(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "   ")
	return encoder.Encode(o)
}
