package sio_samples

import (
	"bytes"
	_ "embed"
	"github.com/Arapak/sio-tool/util"
	"testing"
)

func equalSamples(a [][]byte, b [][]byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !bytes.Equal(a[i], b[i]) {
			return false
		}
	}
	return true
}

//go:embed assets/kol.pdf
var kolorowyWazPdf []byte

//go:embed assets/kol1.in
var kolorowyWazInput1 []byte

//go:embed assets/kol1.out
var kolorowyWazOutput1 []byte

var kolorowyWazInput = [][]byte{kolorowyWazInput1}
var kolorowyWazOutput = [][]byte{kolorowyWazOutput1}

func TestKolorowyWaz(t *testing.T) {
	statement, err := util.PdfToTextRaw(kolorowyWazPdf)
	if err != nil {
		t.Errorf("PdfToTextRaw returned an error: %v", err.Error())
		return
	}
	input, output, err := FindSamples(statement, kolorowyWazPdf)
	if err != nil {
		t.Errorf("FindSamples returned an error: %v", err.Error())
	} else {
		if !equalSamples(input, kolorowyWazInput) {
			t.Errorf("Sample inputs don't match")
		}
		if !equalSamples(output, kolorowyWazOutput) {
			t.Errorf("Sample outputs don't match")
		}
	}
}

//go:embed assets/cyk.pdf
var cykPdf []byte

func TestCyk(t *testing.T) {
	statement, err := util.PdfToTextRaw(cykPdf)
	if err != nil {
		t.Errorf("PdfToTextRaw returned an error: %v", err.Error())
		return
	}
	input, output, err := FindSamples(statement, cykPdf)
	if err != nil {
		t.Errorf("FindSamples returned an error: %v", err.Error())
	} else {
		if len(input) != 0 || len(output) != 0 {
			t.Errorf("return some samples for a interactive problem")
		}
	}
}

//go:embed assets/dom.pdf
var dominoPdf []byte

//go:embed assets/dom1.in
var dominoInput1 []byte

//go:embed assets/dom2.in
var dominoInput2 []byte

//go:embed assets/dom1.out
var dominoOutput1 []byte

//go:embed assets/dom2.out
var dominoOutput2 []byte

var dominoInput = [][]byte{dominoInput1, dominoInput2}
var dominoOutput = [][]byte{dominoOutput1, dominoOutput2}

func TestDomino(t *testing.T) {
	statement, err := util.PdfToTextRaw(dominoPdf)
	if err != nil {
		t.Errorf("PdfToTextRaw returned an error: %v", err.Error())
		return
	}
	input, output, err := FindSamples(statement, dominoPdf)
	if err != nil {
		t.Errorf("FindSamples returned an error: %v", err.Error())
	} else {
		if !equalSamples(input, dominoInput) {
			t.Errorf("Sample inputs don't match")
		}
		if !equalSamples(output, dominoOutput) {
			t.Errorf("Sample outputs don't match")
		}
	}
}

//go:embed assets/amm.pdf
var ammPdf []byte

//go:embed assets/amm1.in
var ammInput1 []byte

//go:embed assets/amm2.in
var ammInput2 []byte

//go:embed assets/amm1.out
var ammOutput1 []byte

//go:embed assets/amm2.out
var ammOutput2 []byte

var ammInput = [][]byte{ammInput1, ammInput2}
var ammOutput = [][]byte{ammOutput1, ammOutput2}

func TestAmm(t *testing.T) {
	statement, err := util.PdfToTextRaw(ammPdf)
	if err != nil {
		t.Errorf("PdfToTextRaw returned an error: %v", err.Error())
		return
	}
	input, output, err := FindSamples(statement, ammPdf)
	if err != nil {
		t.Errorf("FindSamples returned an error: %v", err.Error())
	} else {
		if !equalSamples(input, ammInput) {
			t.Errorf("Sample inputs don't match")
		}
		if !equalSamples(output, ammOutput) {
			t.Errorf("Sample outputs don't match")
		}
	}
}
