package knowledge

import (
	"fmt"
	"strconv"
	"strings"
)

func clean(source string) string {
	source = strings.Replace(source, "\n", "", -1)
	source = strings.Replace(source, "\t", "", -1)
	return strings.TrimSpace(source)
}

func readUntil(source string, stop string) (string, string, error) {
	token := ""
	source = clean(source)
	for {
		if len(source) < len(stop) {
			return source, source, fmt.Errorf("string [%v] is shorter than stop char %v", source, stop)
		}
		if string(source[0:len(stop)]) == stop {
			return clean(token), clean(source[len(stop):]), nil
		}
		token = token + source[0:1]
		source = source[1:]
	}
}

func getExpressions(snippet string) ([]string, error) {
	var err error
	snippet = clean(snippet)
	exprs := []string{}
	for {
		l0 := len(snippet)
		if l0 == 0 {
			return exprs, nil
		}
		if strings.HasPrefix(snippet, "/*") {
			_, snippet, err = readUntil(snippet, "*/")
			if err != nil {
				return exprs, err
			}
		} else {
			var expr string
			expr, snippet, err = readUntil(snippet, ";")
			if err != nil {
				return exprs, err
			}
			exprs = append(exprs, clean(expr))
		}
		if len(snippet) == l0 {
			return exprs, fmt.Errorf("detected loop while getting expressions")
		}
	}
}

func getBlock(source string) (string, string, error) {
	source = clean(source)
	if strings.HasPrefix(source, "{") {
		source = clean(source[1:])
	} else {
		return "", source, fmt.Errorf("expected {")
	}
	opened := 1
	snippet := ""
	for {
		if len(source) == 0 {
			return snippet, source, fmt.Errorf("expected }")
		}
		if string(source[0]) == "{" {
			opened++
		} else if string(source[0]) == "}" {
			opened--
			if opened == 0 {
				return clean(snippet), clean(source[1:]), nil
			}
		}
		snippet = snippet + string(source[0])
		source = source[1:]
	}
}

func parseEventFunction(expr string) (string, float64, bool, error) {
	if strings.HasPrefix(expr, "event(") && string(expr[len(expr)-1]) == ")" {
		expr = expr[len("event("):]
		expr = expr[:len(expr)-1]
		tokens := strings.Split(expr, ",")
		if len(tokens) != 2 {
			return "", 0, false, fmt.Errorf("malformed arguments of event function: %v", expr)
		}
		w, err := strconv.ParseFloat(clean(tokens[1]), 64)
		if err != nil {
			return "", 0, false, fmt.Errorf("malformed weight of event function: %v", tokens[1])
		}
		return clean(tokens[0]), w, true, nil
	} else {
		return "", 0.0, false, nil
	}
}