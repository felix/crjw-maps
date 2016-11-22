package main

import (
	"fmt"
	"github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	qrcode "github.com/skip2/go-qrcode"
	"image/color"
	"log"
	"os"
	"path/filepath"
)

type Location struct {
	Number    string  `csv:"number"`
	Name      string  `csv:"name"`
	AltName   string  `csv:"altname"`
	Latitude  float64 `csv:"latitude"`
	Longitude float64 `csv:"longitude"`
	Comments  string  `csv:"comments"`
}

func (l *Location) getQRCode() string {
	path := filepath.Join(c.cache, fmt.Sprintf("%s-qrc.png", l.Number))
	old := true
	fi, err := os.Stat(path)
	if err == nil && fi.ModTime().After(c.cDate) {
		old = false
		if c.debug {
			log.Printf("Using existing QR code %s", path)
		}
	}
	if old || c.force {
		if c.debug {
			log.Printf("Creating new QR code %s", path)
		}
		url := fmt.Sprintf("geo:%f,%f", l.Latitude, l.Longitude)
		qrcode, err := qrcode.New(url, qrcode.High)
		if err != nil {
			log.Fatal("Failed to create QR code: %s", err)
		}
		err = qrcode.WriteFile(110, path)
		if err != nil {
			log.Fatal("Failed to write QR code: %s", err)
		}
	}
	return path
}

func (l *Location) getMap() string {
	// Flag
	old := true
	path := filepath.Join(c.cache, fmt.Sprintf("%s-map.png", l.Number))
	fi, err := os.Stat(path)
	if err == nil && fi.ModTime().After(c.cDate) {
		old = false
		if c.debug {
			log.Printf("Using existing map file %s", path)
		}
	}

	if old {
		if c.debug {
			log.Printf("Creating new map file %s", path)
		}
		ctx := sm.NewContext()
		// ctx.SetSize(640, 513)
		ctx.SetSize(960, 770)
		pos := s2.LatLngFromDegrees(l.Latitude, l.Longitude)
		kh := s2.LatLngFromDegrees(19.89830, 99.81805)
		ctx.AddMarker(sm.NewMarker(pos, color.RGBA{0xff, 0, 0, 0xff}, 32.0))
		ctx.AddMarker(sm.NewMarker(kh, color.RGBA{0, 0, 0xff, 0xff}, 32.0))
		// ctx.SetCenter(s2.LatLngFromDegrees(19.89830, 99.81805))
		// ctx.SetCenter(pos)
		// ctx.SetTileProvider(sm.NewTileProviderOpenTopoMap())
		// ctx.SetTileProvider(sm.NewTileProviderThunderforestOutdoors())
		ctx.SetTileProvider(sm.NewTileProviderThunderforestLandscape())

		if c.debug {
			log.Printf("Rendering map for %s", l.Number)
		}
		img, err := ctx.Render()
		if err != nil {
			log.Fatal("Failed to render map: ", err)
		}
		if c.debug {
			log.Printf("Saving map for %s", l.Number)
		}
		err = gg.SavePNG(path, img)
		if err != nil {
			log.Fatal("Failed to create image: ", err)
		}
	}
	return path
}
