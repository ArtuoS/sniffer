package downloader

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type KaggleDownloader struct {
}

func NewKaggleDownloader() *KaggleDownloader {
	return &KaggleDownloader{}
}

func (d *KaggleDownloader) FetchCSV(ctx context.Context, datasetURL string) (*bytes.Reader, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, datasetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", os.Getenv("KAGGLE_KEY"))
	req.Header.Set("User-Agent", "sniffer-ingest/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download dataset: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}

	for _, f := range zipReader.File {
		name := f.Name
		if strings.HasSuffix(name, "fra_perfumes.csv") {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("open csv in zip: %w", err)
			}
			defer rc.Close() //nolint:errcheck

			csvBytes, err := io.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("read csv from zip: %w", err)
			}
			return bytes.NewReader(csvBytes), nil
		}
	}

	return nil, fmt.Errorf("fra_perfumes.csv not found in zip archive")
}
