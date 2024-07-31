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
	jsonInputFilename = "input2.json"
)

var log = logrus.New()

func logInit() {
	formatter := &logrus.TextFormatter{}
	formatter.DisableQuote = true

	formatter.DisableTimestamp = false
	log.SetFormatter(formatter)

	//TODO:: TEST:: need to set to Info when running
	//log.Level = logrus.WarnLevel
	//log.Level = logrus.InfoLevel
	log.Level = logrus.DebugLevel
}

func main() {
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

	log.Println(jsonInputMap)

	jsonResultMap := make(map[string]any)

	for key, value := range jsonInputMap {
		if key == "" {
			continue
		} else {
			key = strings.TrimSpace(key)
		}

		//keyParts := strings.Split(key, "_")
		//dataType := keyParts[0]

		for dataType, v := range value {
			switch dataType {
			case "S":
				//fmt.Println(key, value)

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

			case "null":
				fmt.Println(key, value)
			case "list":
				fmt.Println(key, value)
			case "map":
				fmt.Println(key, value)
			}
		}

	}
	fmt.Println(jsonResultMap)
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
