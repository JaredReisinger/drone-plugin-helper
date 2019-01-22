// Package cmd creates command-lines based on tagged structure fields
package cmd

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	tagName      = "cmd"
	termReString = "[[:upper:]][^[:upper:]]*"
)

var (
	matchRe = regexp.MustCompile(fmt.Sprintf("^(?:%s)+$", termReString))
	termRe  = regexp.MustCompile(termReString)
)

// Create generates the command line
// TODO: pass the order? Is that defined by the struct?
func Create(cfg interface{}) (params []string, err error) {
	s, _ := indirect(reflect.ValueOf(cfg))
	params, err = createStructFlags(s)
	return
}

func createStructFlags(s reflect.Value) (params []string, err error) {
	params = make([]string, 0)
	for i := 0; i < s.NumField(); i++ {
		field, _ := indirect(s.Field(i))
		if field.Kind() == reflect.Struct {
			var innerParams []string
			innerParams, err = createStructFlags(field)
			if err != nil {
				return
			}
			// if len(innerParams) > 0 {
			params = append(params, innerParams...)
			// }
			continue
		}

		err = addFieldFlag(&params, s.Type().Field(i), field)
		if err != nil {
			return
		}
	}
	return
}

// indirect returns the underlying field (so long as the pointer isn't null)
func indirect(field reflect.Value) (fieldOut reflect.Value, hadPtr bool) {
	fieldOut = field
	for fieldOut.Kind() == reflect.Ptr {
		hadPtr = true
		fieldOut = fieldOut.Elem()
	}
	return
}

func addFieldFlag(line *[]string, sf reflect.StructField, val reflect.Value) (err error) {
	info, err := infoFromField(sf)
	if err != nil {
		return
	}
	// log.Printf("using info %+v", info)
	if info.omit {
		return
	}

	field, hadPtr := indirect(val)
	kind := field.Kind()
	// log.Printf("adding flag for %v...", kind)
	switch kind {

	case reflect.Bool:
		// no check for hadPtr?
		if val.Bool() {
			*line = append(*line, info.flag)
		} else if info.boolNo {
			negatedFlag, ok := negatedBool(info.flag)
			if !ok {
				err = fmt.Errorf("unable to negate boolean flag %q", info.flag)
				return
			}
			*line = append(*line, negatedFlag)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if hadPtr || val.Int() != 0 {
			if !info.positional {
				*line = append(*line, info.flag)
			}
			*line = append(*line, strconv.FormatInt(val.Int(), 10))
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if hadPtr || val.Uint() != 0 {
			if !info.positional {
				*line = append(*line, info.flag)
			}
			*line = append(*line, strconv.FormatUint(val.Uint(), 10))
		}

	case reflect.String:
		if hadPtr || val.String() != "" {
			if !info.positional {
				*line = append(*line, info.flag)
			}
			*line = append(*line, val.String())
		}

	default:
		err = fmt.Errorf("unsupported parameter type for %q: %q", sf.Name, kind)
		return
	}

	return
}

type tagInfo struct {
	flag       string
	omit       bool
	positional bool
	boolNo     bool
}

func infoFromField(sf reflect.StructField) (info tagInfo, err error) {
	tag := sf.Tag.Get(tagName)

	if tag != "" {
		// log.Printf("got flag tag %q for %q...", tag, sf.Name)
		info, err = parseTagInfo(tag)
		// log.Printf("got info/err: %+v %+v", info, err)
		// } else {
		// 	log.Printf("no flag tag for %q... ???", sf.Name)
	}

	if info.flag == "" {
		param, ok := fieldToParamName(sf.Name)
		if !ok {
			err = fmt.Errorf("unable to create param from %q", sf.Name)
			return
		}
		info.flag = fmt.Sprintf("--%s", param)
	}

	// info, err = parseTagInfo(tag)
	return
}

// parseTagInfo doesn't care about the field type or name... it simply parses
// all of the structured information from the tag value.
func parseTagInfo(tag string) (info tagInfo, err error) {
	tagParts := strings.Split(tag, ",")
	// log.Printf("got tag parts: (%d) %+v", len(tagParts), tagParts)

	for i, part := range tagParts {
		part = strings.TrimSpace(part)

		if i == 0 {
			info.flag = part
			continue
		}

		switch part {
		case "omit":
			info.omit = true
		case "no":
			info.boolNo = true
		case "positional":
			info.positional = true
		default:
			err = fmt.Errorf("unknown cmd tag option: %q", part)
			return
		}
	}

	return
}

func fieldToParamName(name string) (string, bool) {
	if !matchRe.MatchString(name) {
		return "", false
	}

	terms := termRe.FindAllString(name, -1)

	var b strings.Builder
	for i, term := range terms {
		if i > 0 {
			b.WriteString("-")
		}
		b.WriteString(strings.ToLower(term))
	}
	// TODO: separate and hyphenate on capital letter boundaries
	return b.String(), true
}

func negatedBool(flag string) (string, bool) {
	if !strings.HasPrefix(flag, "--") {
		return "", false
	}

	return fmt.Sprintf("--no-%s", strings.TrimPrefix(flag, "--")), true
}
