package properties

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"unicode"
)

var propertyTypes = map[string]string{
	"DIFFICULTY":    "difficulty",
	"GAMEMODE":      "gamemode",
	"MAX_PLAYERS":   "max-players",
	"SERVER_PORT":   "server-port",
	"VIEW_DISTANCE": "view-distance",
	"WHITE_LIST":    "white-list",
}

type Properties struct {
	path   string
	values []string
}

func New(path string) *Properties {
	return &Properties{
		path:   path,
		values: make([]string, 0, 10),
	}
}

func (p *Properties) LoadFromEnv() error {
	for k := range propertyTypes {
		value := os.Getenv(k)
		slog.Debug("reading prop from env:", k, value)
		if value == "" {
			continue
		}

		err := p.Add(fmt.Sprintf("%s=%s", k, value))
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Properties) Add(prop string) error {
	parsedProp, err := parseProp(prop)
	if err != nil {
		return err
	}

	p.values = append(p.values, parsedProp)
	return nil
}

func (p *Properties) Write() error {
	file, err := os.Create(p.path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, v := range p.values {
		_, err := writer.WriteString(v + "\n")
		if err != nil {
			return err
		}
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func parseProp(prop string) (string, error) {
	vals := strings.Split(prop, "=")

	if len(vals) != 2 {
		return "", fmt.Errorf("unable to parse invalid property: %s", prop)
	}

	key := propertyTypes[vals[0]]
	if key == "" {
		return "", fmt.Errorf("invalid or unsupported property: %s", prop)
	}

	return fmt.Sprintf("%s=%s", key, parse(vals[1])), nil
}

func parse(key string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case '_':
			return '-'
		default:
			return unicode.ToLower(r)
		}
	}, key)
}
