package utils

import (
	"bufio"
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

func LoadEnv(filepath string) error {
	file, err := os.Open(filepath)

	if err != nil {
		log.Errorf("unable to open file %s: %s\n", filepath, err.Error())
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Errorf("error reading %s: %s", filepath, err.Error())
		return err
	}

	for i, line := range lines {
		if len(line) < 1 {
			log.Debugf("line %d is empty", i)
			continue
		}

		if strings.HasPrefix(line, "#") {
			log.Debugf("line %d is commented", i)
			continue
		}

		parts := strings.SplitN(line, "=", 2)

		if len(parts) != 2 {
			log.Debugf("line %d missing val", i)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		os.Setenv(key, value)
		log.Debugf("set key %s", key)
	}

	return nil
}
