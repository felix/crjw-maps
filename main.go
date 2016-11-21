package main

import (
	"flag"
	"github.com/gocarina/gocsv"
	"github.com/jung-kurt/gofpdf"
	"log"
	"os"
	"time"
)

const (
	pWidth  = 148
	pHeight = 105
)

var c struct {
	in    string
	out   string
	cache string
	cDate time.Time
	debug bool
}

func main() {
	flag.StringVar(&c.in, "in", "", "Input CSV file")
	flag.StringVar(&c.out, "out", "", "Output PDF file")
	flag.StringVar(&c.cache, "cache", "cache", "Cache path")
	flag.BoolVar(&c.debug, "debug", false, "Debug output")
	// flag.StringVar(&c.outPrefix, "outPrefix", "", "Output file prefix")

	flag.Parse()

	if c.in == "" {
		log.Fatal("Missing input file")
	}
	if c.out == "" {
		c.out = "output.pdf"
	}

	if c.debug {
		log.Printf("Reading file %s", c.in)
	}
	f, err := os.Open(c.in)
	if err != nil {
		log.Fatal("Error opening file: %s", err)
	}
	defer f.Close()
	locations := []*Location{}

	err = gocsv.UnmarshalFile(f, &locations)
	if err != nil {
		log.Fatal("Error parsing CSV: ", err)
	}
	fi, err := os.Stat(c.in)
	if err != nil {
		log.Fatal("Error statting CSV: ", err)
	}
	c.cDate = fi.ModTime()

	log.Printf("Processing %d locations", len(locations))

	err = os.MkdirAll(c.cache, os.ModePerm)
	if err != nil {
		log.Fatal("Failed to create cache: ", err)
	}

	// Generate the PDF
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		OrientationStr: "L",
		UnitStr:        "mm",
		Size:           gofpdf.SizeType{Wd: pHeight, Ht: pWidth},
	})
	pdf.AddFont("Norasi", "", "Norasi.json")
	pdf.SetFont("Norasi", "", 16)
	pdf.SetTextColor(0, 0, 0)
	tr := pdf.UnicodeTranslatorFromDescriptor("cp874")

	// Generate the PDF
	for i := range locations {
		loc := locations[i]
		log.Printf("Processing map number %s %s", loc.Number, loc.Name)
		pdf.AddPage()

		imageOpts := gofpdf.ImageOptions{"PNG", true}

		// QR Code
		qrPath := loc.getQRCode()
		// qrInfo := pdf.RegisterImageOptions(qrPath, imageOpts)

		// Map
		smPath := loc.getMap()
		// mapInfo := pdf.RegisterImageOptions(smPath, imageOpts)

		// Name
		if loc.Name != "" {
			pdf.TransformBegin()
			pdf.TransformRotate(90, 5, pHeight-5)
			pdf.SetFontSize(40)
			pdf.Text(5, pHeight+5, loc.Number)
			pdf.SetFontSize(20)
			pdf.Text(5, pHeight+5+7, tr(loc.Name))
			pdf.ImageOptions(qrPath, pHeight-27, pHeight-10, 0, 0, false, imageOpts, 0, "")
			pdf.TransformEnd()
		}

		pdf.ImageOptions(smPath, 25, 5, 0, pHeight-10, false, imageOpts, 0, "")
		if c.debug {
			log.Printf("Added map %s", smPath)
		}
	}

	err = pdf.OutputFileAndClose(c.out)
	if err != nil {
		log.Fatal("Error generating PDF: ", err)
	}
}
