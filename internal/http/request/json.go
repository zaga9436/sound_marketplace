package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
)

func DecodeJSON(r *http.Request, target interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		return fmt.Errorf("empty body")
	}

	normalized := body
	if !utf8.Valid(body) {
		decoded, err := charmap.Windows1251.NewDecoder().Bytes(body)
		if err != nil {
			return fmt.Errorf("invalid request encoding: %w", err)
		}
		normalized = decoded
	}

	decoder := json.NewDecoder(bytes.NewReader(normalized))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}

	return nil
}
