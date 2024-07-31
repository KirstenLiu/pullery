package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	jsonInputFilename = "input.json"
)

var log = logrus.New()

func logInit() {
	formatter := &logrus.TextFormatter{}
	formatter.DisableQuote = true

	formatter.DisableTimestamp = false
	log.SetFormatter(formatter)

	//Please set to Warn when running
	log.Level = logrus.WarnLevel
	//log.Level = logrus.InfoLevel
	//log.Level = logrus.DebugLevel
}

func main() {
	logInit()
	jsonInputFile, err := os.Open(jsonInputFilename)
	if err != nil {
		log.Fatalf("Can't open json input file %v: %v\n", jsonInputFilename, err)
	}
	defer jsonInputFile.Close()

	jsonBytes, err := io.ReadAll(jsonInputFile)
	if err != nil {
		log.Fatalf("Can't read all from jsoninputFile %v: %v\n", jsonInputFilename, err)
	}

	var jsonInputMap map[string]map[string]any
	json.Unmarshal(jsonBytes, &jsonInputMap)

	log.Debugln(jsonInputMap)

	jsonResultMap := make(map[string]any)

	for key, value := range jsonInputMap {
		if key == "" {
			continue
		} else {
			key = strings.TrimSpace(key)
		}

		for dataType, v := range value {
			//In the transformation instruction, **N** denotes the value's data type, and the sanitize of trailing and leading whitespace is only defined to "value"
			//So whitespace in dataType is not processed and considered illegal.
			//For example, "null_1": { "NULL ": "true"} should be considered illegal.
			//But in the sample output it is included, so I process the trailing zero for data type as well.
			dataType := strings.TrimSpace(dataType)

			switch dataType {
			case "S":
				matched, str := processString(v)
				if matched {
					jsonResultMap[key] = str
				}
			case "N":
				matched, num := processNumber(v)
				if matched {
					jsonResultMap[key] = num
				}
			case "BOOL":
				matched, b := processBoolean(v)
				if matched {
					jsonResultMap[key] = b
				}
			case "NULL":
				matched, n := processNull(v)
				if matched {
					jsonResultMap[key] = n
				}
			case "L":
				switch l := v.(type) {
				case []any:
					matched, result := processList(l)
					if matched {
						jsonResultMap[key] = result
					}
				}
			case "M":
				switch m := v.(type) {
				case map[string]any:
					matched, result := processMap(m)
					if matched {
						jsonResultMap[key] = result
					}
				}
			}
		}

	}
	//Map is always sorted in fmt print.
	log.Debugln(jsonResultMap)
	mapJson, err := json.Marshal(jsonResultMap)
	if err != nil {
		log.Fatalf("Error occurs while marshalling map %v to json: %v\n", jsonResultMap, err)
	}
	fmt.Println(string(mapJson))
}

func processNumber(value any) (bool, any) {
	//Don't consider 8. as legal, it should be 8.0 for legal numeric
	numStr := strings.TrimSpace(value.(string))
	match, err := regexp.MatchString(`[+-]?\d+(.\d+)?$`, numStr)
	if err != nil {
		log.Fatalf("Error occurs while doing regex match for number %v: %v", numStr, err)
	}
	if !match {
		return false, nil
	} else {
		//Consider it is float64 when there is ".": according to instruction "be transformed to the relevant `Numeric` data type."
		//so 8.0 will be considered required for float64.

		//Leading zero is trimmed automatically by Go Atoi and ParseFloat.
		if strings.Contains(numStr, ".") {
			vFloat, _ := strconv.ParseFloat(numStr, 64)
			return true, vFloat
		} else {
			vInt, _ := strconv.Atoi(numStr)
			return true, vInt
		}
	}
}

func processString(value any) (bool, any) {
	str := strings.TrimSpace(value.(string))
	if len(str) == 0 {
		return false, ""
	} else {
		match, err := regexp.MatchString(`^\d{4}-\d{2}-\d{2}T(\d{2}:\d{2}:\d{2}(\.\d+)?)(Z|[\+-]\d{2}:\d{2})$`, str)
		if err != nil {
			log.Fatalf("Error occurs while doing regex match for string %v: %v", str, err)
		}
		if match {
			t, err := time.Parse(time.RFC3339, str)
			if err != nil {
				log.Fatalf("Error occurs while paring time str %v: %v\n", str, err)
			}
			return true, t.Unix()
		} else {
			return true, str
		}
	}
}

func processBoolean(value any) (bool, bool) {
	b := strings.TrimSpace(value.(string))

	if b == "1" || b == "t" || b == "T" || b == "TRUE" || b == "true" || b == "True" {
		return true, true
	} else if b == "0" || b == "f" || b == "F" || b == "FALSE" || b == "false" || b == "False" {
		return true, false
	} else {
		return false, false
	}
}

func processNull(value any) (bool, string) {
	b := strings.TrimSpace(value.(string))

	if b == "1" || b == "t" || b == "T" || b == "TRUE" || b == "true" || b == "True" {
		return true, "null"
	} else if b == "0" || b == "f" || b == "F" || b == "FALSE" || b == "false" || b == "False" {
		return false, ""
	} else {
		return false, ""
	}
}

func processList(values []any) (bool, []any) {
	var final []any
	for _, item := range values {
		switch item := item.(type) {
		case map[string]any:
			for dataType, v := range item {
				dataType := strings.TrimSpace(dataType)

				switch dataType {
				case "S":
					matched, str := processString(v)
					if matched {
						final = append(final, str)
					}
				case "N":
					matched, num := processNumber(v)
					if matched {
						final = append(final, num)
					}
				case "BOOL":
					matched, b := processBoolean(v)
					if matched {
						final = append(final, b)
					}
				}
			}
		}
	}
	if len(final) > 0 {
		return true, final
	} else {
		return false, []any{}
	}
}

func processMap(values map[string]any) (bool, map[string]any) {
	final := make(map[string]any)

	for k, value := range values {
		switch value.(type) {
		case map[string]any:
			for dataType, v := range value.(map[string]any) {
				dataType := strings.TrimSpace(dataType)

				switch dataType {
				case "S":
					matched, str := processString(v)
					if matched {
						final[k] = str
					}
				case "N":
					matched, num := processNumber(v)
					if matched {
						final[k] = num
					}
				case "BOOL":
					matched, b := processBoolean(v)
					if matched {
						final[k] = b
					}
				case "NULL":
					matched, n := processNull(v)
					if matched {
						final[k] = n
					}
				case "L":
					switch l := v.(type) {
					case []any:
						matched, result := processList(l)
						if matched {
							final[k] = result
						}
					}
				case "M":
					switch m := v.(type) {
					case map[string]any:
						matched, result := processMap(m)
						if matched {
							final[k] = result
						}
					}
				}
			}
		}
	}
	if len(final) > 0 {
		return true, final
	} else {
		return false, final
	}
}
