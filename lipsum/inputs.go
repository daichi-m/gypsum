package lipsum

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Input struct encapsulates the input parameters of command line
type Input struct {
	File  string
	Size  int
	Count int

	// internal field
	fileCtr   int
	generator Generator
}

// Export the standard file sizes
const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
	TB = 1024 * GB
)

const fileNamePattern = "lipsum_%04d.txt"
const maxFileSize = 500 * MB
const maxCount = 200

// Validate the input parameters and return error in case of validation failure.
func (inp *Input) Validate() error {
	if inp.Count <= 0 {
		return fmt.Errorf("Count of files cannot be zero or negative")
	}
	if inp.Count > maxCount {
		return fmt.Errorf("Count of files cannot exceed %d", maxCount)
	}

	if inp.Size <= 0 {
		return fmt.Errorf("Size of each file cannot be zero or negative")
	}
	if inp.Size > maxFileSize {
		return fmt.Errorf("Size of each file cannot exceed %d MB", maxFileSize/MB)
	}

	if len(inp.File) == 0 {
		inp.File = "-"
	}
	inp.fileCtr = 1
	var err error
	inp.generator, err = NewGenerator(false)
	if err != nil {
		return err
	}
	return nil
}

// GenerateOut generates the output files for the given Input parameters.
func (inp *Input) GenerateOut() error {
	for ctr := 1; ctr <= inp.Count; ctr++ {
		file, err := inp.createFile()
		if err != nil {
			return err
		}
		c, err := inp.generateContent()
		if err != nil {
			return err
		}
		sz, err := file.WriteString(c)
		if err != nil {
			return err
		}
		err = file.Close()
		if err != nil {
			return err
		}
		log.Printf("%d MB content written in %s file", (sz / MB), file.Name())
	}
	return nil
}

func (inp *Input) createFile() (*os.File, error) {
	if len(inp.File) == 0 || inp.File == "-" {
		return os.Stdout, nil
	}

	path := filepath.Clean(inp.File)
	if inp.Count == 1 {
		return os.Create(path)
	}

	if inp.fileCtr > inp.Count {
		return nil, fmt.Errorf("The file counter has exceeded the count of files, stop generating")
	}

	err := os.MkdirAll(path, os.ModeDir|0755)
	if err != nil {
		return nil, err
	}
	fpath := filepath.Join(path, fmt.Sprintf(fileNamePattern, inp.fileCtr))
	inp.fileCtr++
	return os.Create(fpath)
}

func (inp *Input) generateContent() (string, error) {

	if inp.Size > maxFileSize {
		return "", fmt.Errorf("Max file size %f exceeded, bailing out", float64(maxFileSize)/MB)
	}

	content := new(strings.Builder)
	contSz := 0
	for {
		p := inp.generator.Paragraph()
		_, err := content.WriteString(p)
		if err != nil {
			return "", err
		}
		_, err = content.WriteString(newline)
		if err != nil {
			return "", err
		}
		contSz = contSz + len(p) + 1 // +1 for the newline
		if contSz >= inp.Size {
			return content.String(), nil
		}
	}
}
