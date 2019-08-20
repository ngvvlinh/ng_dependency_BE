package l

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const DefaultRoute = "/==/logging"
const MaxVerbosity = 9
const EnvKey = "ETOP_LOG"
const deprecatedEnvKey = "ETOP_LOG_DEBUG"

// verbosity from 1 to 9
type WildcardPatterns [MaxVerbosity][]*regexp.Regexp

var (
	errEnablerNotFound = errors.New("enabler not found")
	errLevelNil        = errors.New("must specify a logging level")
	errInvalidLevel    = errors.New("invalid level")

	enablers    = make(map[string]zap.AtomicLevel)
	envPatterns WildcardPatterns
)

func RegisterHTTPHandler(mux *http.ServeMux) {
	mux.HandleFunc(DefaultRoute, ServeHTTP)
}

// ServeHTTP supports logging level with an HTTP request.
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	type payload struct {
		Name    string `json:"name"`
		Pattern string `json:"pattern,omitempty"`
		Level   string `json:"level"`
	}
	enc := json.NewEncoder(w)
	badRequest := func(msg string) {
		w.WriteHeader(http.StatusBadRequest)
		_ = enc.Encode(errorResponse{Error: msg})
	}
	encodeEnablers := func(enablers map[string]zap.AtomicLevel) {
		var payloads []payload
		for k, e := range enablers {
			lvl := capitalLevel(e.Level())
			payloads = append(payloads, payload{
				Name:  k,
				Level: lvl,
			})
		}
		sort.Slice(payloads, func(i, j int) bool {
			return payloads[i].Name < payloads[j].Name
		})
		_ = enc.Encode(payloads)
	}

	switch r.Method {
	case "GET":
		encodeEnablers(enablers)

	case "POST":
		var req payload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			badRequest(err.Error())
			return
		}

		if (req.Name == "") == (req.Pattern == "") {
			badRequest("must provide name or pattern (but not both)")
			return
		}

		if req.Name != "" {
			if req.Level == "" {
				badRequest(errLevelNil.Error())
				return
			}
			level, ok := unmarshalLevel(req.Level)
			if !ok {
				badRequest(errInvalidLevel.Error())
				return
			}

			enabler, ok := enablers[req.Name]
			if !ok {
				badRequest(errEnablerNotFound.Error())
				return
			}
			enabler.SetLevel(level)

		} else {
			if req.Level != "" {
				badRequest("pattern must not be used with level")
				return
			}

			patterns, errPattern, err := parseWildcardPatterns(req.Pattern)
			if err != nil {
				badRequest(fmt.Sprintf("bad pattern (%v): %v", errPattern, err))
				return
			}
			for name, enabler := range enablers {
				setLogLevelFromPatterns(patterns, name, enabler)
			}
		}
		encodeEnablers(enablers)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = enc.Encode(errorResponse{
			Error: "only GET and POST are supported",
		})
	}
}

func init() {
	err := zap.RegisterEncoder(ConsoleEncoderName,
		func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
			return NewConsoleEncoder(cfg), nil
		})
	if err != nil {
		panic(err)
	}

	ll = New()
	xl = New(zap.AddCallerSkip(1))

	envLog := os.Getenv(EnvKey)
	if envLog == "" {
		envLog = os.Getenv(deprecatedEnvKey)
	}
	if envLog == "" {
		return
	}

	var errPattern string
	envPatterns, errPattern, err = parseWildcardPatterns(envLog)
	if err != nil {
		ll.Fatal(fmt.Sprintf("Unable to parse %v", EnvKey), String("invalid", errPattern), Error(err))
	}
	ll.Info("Enable debug log", String(EnvKey, envLog))
}

func parseWildcardPatterns(input string) (result WildcardPatterns, errPattern string, err error) {
	patterns := strings.Split(input, ",")
	for _, p := range patterns {
		lvl, r, err := parsePattern(p)
		if err != nil {
			return WildcardPatterns{}, p, err
		}
		if lvl < 0 {
			lvl = -lvl
		}
		if lvl > MaxVerbosity {
			return WildcardPatterns{}, p, fmt.Errorf("verbosity must be from 1 to %v", MaxVerbosity)
		}
		result[lvl] = append(result[lvl], r)
	}
	return
}

func parsePattern(p string) (lvl int, pattern *regexp.Regexp, err error) {
	parts := strings.Split(p, ":")
	switch len(parts) {
	case 1:
		lvl = 1 // default to DebugLevel

	case 2:
		lvl, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, nil, err
		}
		p = parts[0]

	default:
		return 0, nil, fmt.Errorf("must be format PATTERN:LEVEL")
	}

	p = strings.Replace(strings.Trim(p, " "), "*", ".*", -1)
	pattern, err = regexp.Compile(p)
	return
}

func setLogLevelFromPatterns(wildcardPatterns WildcardPatterns, name string, enabler zap.AtomicLevel) {
	for lvl, patterns := range wildcardPatterns {
		for _, pattern := range patterns {
			if pattern.MatchString(name) {
				enabler.SetLevel(zapcore.Level(-lvl))
			}
		}
	}
}
