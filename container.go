package lenex

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
)

var zipMagic = []byte{'P', 'K', 0x03, 0x04}

// readXML returns the raw XML bytes from r, transparently handling both
// plain .lef (XML) and .lxf (zip-wrapped XML) inputs.
func readXML(r io.Reader) ([]byte, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}
	if len(buf) >= 4 && bytes.Equal(buf[:4], zipMagic) {
		return extractZipXML(buf)
	}
	return buf, nil
}

func extractZipXML(buf []byte) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
	if err != nil {
		return nil, fmt.Errorf("opening zip: %w", err)
	}
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("opening zip entry %q: %w", f.Name, err)
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("reading zip entry %q: %w", f.Name, err)
		}
		return data, nil
	}
	return nil, fmt.Errorf("zip archive contains no entries")
}
