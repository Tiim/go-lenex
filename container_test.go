package lenex

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

const testXML = `<?xml version="1.0"?><LENEX version="3.0"></LENEX>`

func TestReadXMLPlain(t *testing.T) {
	got, err := readXML(strings.NewReader(testXML))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != testXML {
		t.Errorf("plain: got %q, want %q", got, testXML)
	}
}

func TestReadXMLZip(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create("meet.lef")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte(testXML)); err != nil {
		t.Fatal(err)
	}
	zw.Close()

	got, err := readXML(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != testXML {
		t.Errorf("zip: got %q, want %q", got, testXML)
	}
}
