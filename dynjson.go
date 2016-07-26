package dynjson

import (
	"encoding/json"
	"fmt"
	"strings"
)

type dynObject struct {
	dynData map[string]string
}

// New: Creates a new dynObject
func New(jsonFile []byte) *dynObject {
	var jsonfile map[string]interface{}

	err := json.Unmarshal(jsonFile, &jsonfile)
	handleErr(err)

	dynData := &dynObject{dynData: parser(jsonfile, "", 0)}

	return dynData
}

// GetMap: returns the parsed file
func (dynData *dynObject) GetMap() map[string]string {
	return dynData.dynData
}

func parser(jsonData map[string]interface{}, base string, depth int) (retMap map[string]string) {
	var newBase string
	retMap = make(map[string]string)

	for key, value := range jsonData {
		element, isMap := value.(map[string]interface{})
		elem, isSlice := value.([]interface{})
		if isMap {
			if base == "" {
				newBase = fmt.Sprintf("%s", key)
			} else {
				newBase = fmt.Sprintf("%s.%s", base, key)
			}
			retMap = merge(retMap, parser(element, newBase, depth+1))
		} else if isSlice {
			retMap = sliceParser(elem, base, key)
		} else {
			if depth != 0 {
				mapper := fmt.Sprintf("%s.%s", base, key)
				retMap[mapper] = value.(string)
			} else {
				retMap[key] = value.(string)
			}
		}
	}
	return
}

func sliceParser(elem []interface{}, base string, key string) (retMap map[string]string) {
	var mapper string
	length := len(elem)
	retMap = make(map[string]string)

	if base == "" {
		mapper = fmt.Sprintf("dynObject.%s", key)
	} else {
		mapper = fmt.Sprintf("dynObject.%s.%s", base, key)
	}
	valuer := fmt.Sprintf("[")
	for i := 0; i < length-1; i++ {
		valuer = fmt.Sprintf("%s%s, ", valuer, elem[i])
	}
	valuer = fmt.Sprintf("%s%s]", valuer, elem[length-1])
	retMap[mapper] = valuer

	return
}

// SearchKey: searchs a given key in our object
func (dynData *dynObject) SearchKey(regex string) (regexData map[string]string, mapObject map[string]bool) {
	mapObject = make(map[string]bool)
	regexData = make(map[string]string)
	for key, value := range dynData.GetMap() {
		splited := strings.Split(key, "dynObject.")
		var idx = 0
		if splited[0] == "" {
			idx = 1
		}

		if key == regex {
			regexData[key] = value
			break
		} else if strings.Contains(splited[idx], regex) {
			if strings.Contains(key, "dynObject") {
				mapObject[key] = true
			}
			regexData[key] = value
		}
	}
	return
}

// ParseMultiValue: parses a dynObject that contains multiple values
func (dynData *dynObject) ParseMultiValue(regex string) (mapper string, arr []string, notFound bool) {
	notFound = true
	for key, value := range dynData.GetMap() {
		if strings.Contains(key, "dynObject") && strings.Contains(key, regex) {
			mapper = strings.Split(key, "dynObject.")[1]
			trimmed := strings.Trim(strings.Trim(value, "["), "]")
			splited := strings.Split(trimmed, ", ")

			for idx := range splited {
				arr = append(arr, splited[idx])
			}
			notFound = false
		}
	}
	return
}
