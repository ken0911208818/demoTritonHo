package httputil

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func Bind(r *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}

	return nil
}

func BindForUpdate(r *http.Request, obj interface{}) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	inputBytes := buf.Bytes()
	//json Unmarshal json 轉換乘 struct
	//json Marshal  struct 轉換成 json
	if err := json.Unmarshal(inputBytes, obj); err != nil {
		return err
	}
	return nil
}
