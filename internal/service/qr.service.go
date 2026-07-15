package service

import (
	"context"
	"fmt"

	"github.com/nopalllfd/shortlink-server/internal/dto"
	"github.com/nopalllfd/shortlink-server/internal/storage"
	"github.com/nopalllfd/shortlink-server/pkg/qr"
)

type QRService struct {
	storage storage.ObjectStorage
}

func NewQRService(storage storage.ObjectStorage) *QRService {
	return &QRService{
		storage: storage,
	}
}

func (s *QRService) Generate(
	ctx context.Context,
	linkID int,
	shortURL string,
	cfg *dto.QRConfig,
) (string, error) {

	var customCfg qr.CustomQRConfig
	if cfg != nil {
		customCfg = qr.CustomQRConfig{
			Size:       cfg.Size,
			Foreground: cfg.Foreground,
			Background: cfg.Background,
			Style:      cfg.Style,
			LogoURL:    cfg.LogoURL,
		}
	} else {
		customCfg = qr.CustomQRConfig{
			Size: 512,
		}
	}

	imageBytes, err := qr.GenerateCustom(shortURL, customCfg)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("qr/%d.png", linkID)

	url, err := s.storage.Upload(
		ctx,
		key,
		imageBytes,
		"image/png",
	)
	if err != nil {
		return "", err
	}

	return url, nil
}
