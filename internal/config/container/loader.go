package container

import (
	"fmt"
)

func Load() error {
	if err := LoadEnvs(); err != nil {
		return fmt.Errorf("load .env: %w", err)
	}

	if err := LoadLogger(); err != nil {
		return fmt.Errorf("load logger: %w", err)
	}

	return nil
}
