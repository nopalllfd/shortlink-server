package qr

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"golang.org/x/image/draw"
)

type CustomQRConfig struct {
	Size       int
	Foreground string
	Background string
	Style      string
	LogoURL    string
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

func GeneratePNG(content string, size int) ([]byte, error) {
	return GenerateCustom(content, CustomQRConfig{
		Size: size,
	})
}

func resizeImage(src image.Image, targetWidth, targetHeight int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}

func GenerateCustom(content string, cfg CustomQRConfig) ([]byte, error) {
	// 1. Create the QR code instance
	qrc, err := qrcode.New(content)
	if err != nil {
		return nil, fmt.Errorf("failed to create qrcode instance: %w", err)
	}

	// 2. Fetch and decode logo image if provided
	var logoImage image.Image
	if cfg.LogoURL != "" {
		resp, err := http.Get(cfg.LogoURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch logo URL: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to fetch logo URL, status code: %d", resp.StatusCode)
		}
		img, _, errDec := image.Decode(resp.Body)
		if errDec != nil {
			return nil, fmt.Errorf("failed to decode logo image: %w", errDec)
		}
		logoImage = img
	}

	// 3. Prepare writer options
	var options []standard.ImageOption

	// Block width based on size
	blockWidth := uint8(10)
	if cfg.Size > 0 {
		bw := cfg.Size / 30
		if bw < 5 {
			bw = 5
		}
		if bw > 100 {
			bw = 100
		}
		blockWidth = uint8(bw)
	}
	options = append(options, standard.WithQRWidth(blockWidth))

	// Colors
	if cfg.Foreground != "" {
		options = append(options, standard.WithFgColorRGBHex(cfg.Foreground))
	} else {
		options = append(options, standard.WithFgColorRGBHex("#000000"))
	}

	if cfg.Background != "" {
		options = append(options, standard.WithBgColorRGBHex(cfg.Background))
	} else {
		options = append(options, standard.WithBgColorRGBHex("#ffffff"))
	}

	// Style
	if cfg.Style == "circle" {
		options = append(options, standard.WithCircleShape())
	}

	// Logo
	if logoImage != nil {
		// Scale logo to ~16% (1/6) of the QR code image size
		logoSize := 80
		if cfg.Size > 0 {
			logoSize = cfg.Size / 6
		}
		resizedLogo := resizeImage(logoImage, logoSize, logoSize)
		options = append(options, standard.WithLogoImage(resizedLogo))
	}

	// 4. Draw to memory buffer
	buf := new(bytes.Buffer)
	w := standard.NewWithWriter(nopCloser{buf}, options...)

	if err := qrc.Save(w); err != nil {
		return nil, fmt.Errorf("failed to save qrcode to writer: %w", err)
	}

	return buf.Bytes(), nil
}
