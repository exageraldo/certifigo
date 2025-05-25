package certifigo

import (
	"bytes"
	"fmt"
	"path"
	"strconv"
	"strings"
	tt "text/template"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/pelletier/go-toml/v2/unstable"
)

type WxHSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`

	raw string // "WxH"
}

// func (s *WxHSize) UnmarshalJSON(p []byte) error {
// 	return nil
// }

func (s *WxHSize) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%dx%d"`, s.Width, s.Height)), nil
}

func (s *WxHSize) UnmarshalTOML(value *unstable.Node) error {
	pattern := `^(?P<width>\d+)[xX](?P<height>\d+)$`
	matches, err := FindNamedMatches(string(value.Data), pattern)
	if err != nil {
		return fmt.Errorf("error parsing WxHSize: %v", err)
	}

	width, err := strconv.Atoi(matches["width"])
	if err != nil {
		return fmt.Errorf("error converting width to float64: %v", err)
	}

	height, err := strconv.Atoi(matches["height"])
	if err != nil {
		return fmt.Errorf("error converting height to float64: %v", err)
	}

	s.Width = width
	s.Height = height
	s.raw = string(value.Data)
	return nil
}

type HexColor struct {
	R, G, B, A uint8

	raw string // "#RRGGBB" or "#RRGGBB[AAA%]"
}

func (h *HexColor) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, h.raw)), nil
}

func (h *HexColor) UnmarshalTOML(value *unstable.Node) error {
	pattern := `^#(?P<r>[0-9A-Fa-f]{2})(?P<g>[0-9A-Fa-f]{2})(?P<b>[0-9A-Fa-f]{2})(?:\[(?P<a>\d{1,3})%\])?$`
	matches, err := FindNamedMatches(string(value.Data), pattern)
	if err != nil {
		return fmt.Errorf("error parsing HexColor: %v", err)
	}

	r, err := strconv.ParseUint(matches["r"], 16, 8)
	if err != nil {
		return fmt.Errorf("error converting red component to uint8: %v", err)
	}

	g, err := strconv.ParseUint(matches["g"], 16, 8)
	if err != nil {
		return fmt.Errorf("error converting green component to uint8: %v", err)
	}

	b, err := strconv.ParseUint(matches["b"], 16, 8)
	if err != nil {
		return fmt.Errorf("error converting blue component to uint8: %v", err)
	}

	h.R = uint8(r)
	h.G = uint8(g)
	h.B = uint8(b)

	var alphaValue int
	if alpha, ok := matches["a"]; ok && alpha != "" {
		alphaValue, err = strconv.Atoi(alpha)
		if err != nil {
			return fmt.Errorf("error converting alpha percentage to int: %v", err)
		}
		if alphaValue < 0 || alphaValue > 100 {
			return fmt.Errorf("alpha percentage must be between 0 and 100: %v", alphaValue)
		}
		a, err := strconv.ParseUint(alpha, 10, 8)
		if err != nil {
			return fmt.Errorf("error converting alpha component to int8: %v", err)
		}
		h.A = uint8(a * 255 / 100)
	} else {
		alphaValue = 100
		h.A = 255 // default alpha value
	}

	h.raw = fmt.Sprintf("#%02X%02X%02X[%d%%]", h.R, h.G, h.B, alphaValue)
	return nil
}

type StringDate string // "DD/MM/YYYY"

func (s *StringDate) ParseDate(date string) (*time.Time, error) {
	const layout = "02/01/2006" // DD/MM/YYYY format
	parsedDate, err := time.Parse(layout, date)
	if err != nil {
		return nil, fmt.Errorf("error parsing StringDate: %v", err)
	}
	return &parsedDate, nil
}

func (s *StringDate) UnmarshalTOML(value *unstable.Node) error {
	_, err := s.ParseDate(string(value.Data))
	if err != nil {
		return err
	}
	*s = StringDate(value.Data)
	return nil
}

func ParseTOMLFile(fileContent []byte, v any) error {
	decoder := toml.NewDecoder(bytes.NewReader(fileContent))
	decoder.EnableUnmarshalerInterface()

	if err := decoder.Decode(v); err != nil {
		return err
	}

	return nil
}

func ParseTOMLTemplate(filePath string, v any, data map[string]any) error {
	t, err := tt.New(path.Base(filePath)).ParseFiles(filePath)
	if err != nil {
		return err
	}
	var buff strings.Builder
	if err := t.Execute(&buff, data); err != nil {
		return err
	}
	return ParseTOMLFile([]byte(buff.String()), v)
}

func ParseTOMLTemplateFS(filePath string, v any, data map[string]any) error {
	t, err := tt.New(path.Base(filePath)).ParseFS(assetsDir, filePath)
	if err != nil {
		return err
	}
	buff := new(bytes.Buffer)
	if err := t.Execute(buff, data); err != nil {
		return err
	}

	return ParseTOMLFile(buff.Bytes(), v)
}
