package gomergy

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/asdfgugus/godiffy/pkg/godiffy"
)

func MergeToPath(diff godiffy.Diff, path string) error {
	if _, err := os.ReadDir(path); err != nil {
		return fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	for _, file := range diff.Files {
		switch file.Status {
		case godiffy.FileStatusDeleted:
			err := handleDeletedFile(file, path)
			if err != nil {
				return fmt.Errorf("failed to handle deleted file %s: %w", file.NewPath, err)
			}
		case godiffy.FileStatusNew:
			err := handleNewFile(file, path)
			if err != nil {
				return fmt.Errorf("failed to handle new file %s: %w", file.NewPath, err)
			}
		case godiffy.FileStatusModified:
			err := handleModifiedFile(file, path)
			if err != nil {
				return fmt.Errorf("failed to handle modified file %s: %w", file.NewPath, err)
			}
		}
	}
	return nil
}

func handleDeletedFile(file *godiffy.FileDiff, path string) error {
	if _, err := os.Open(filepath.Join(path, file.NewPath)); os.IsNotExist(err) {
		return nil
	}

	err := os.Remove(filepath.Join(path, file.NewPath))
	if err != nil {
		return fmt.Errorf("failed to remove file %s: %w", file.NewPath, err)
	}

	return nil
}

func handleNewFile(file *godiffy.FileDiff, path string) error {
	err := os.MkdirAll(filepath.Dir(filepath.Join(path, file.NewPath)), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", filepath.Dir(filepath.Join(path, file.NewPath)), err)
	}

	content := ""
	for _, hunk := range file.Hunks {
		for _, line := range hunk.Lines {
			content += line.Content
		}
	}
	fileMode, err := strconv.ParseInt(file.NewMode, 8, 0)
	if err != nil {
		return fmt.Errorf("failed to convert file mode %s: %w", file.NewMode, err)
	}
	err = os.WriteFile(filepath.Join(path, file.NewPath), []byte(content), os.FileMode(fileMode))
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", file.NewPath, err)
	}

	return nil
}

func handleModifiedFile(file *godiffy.FileDiff, path string) error {
	err := os.MkdirAll(filepath.Dir(filepath.Join(path, file.NewPath)), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", filepath.Dir(filepath.Join(path, file.NewPath)), err)
	}

	content := ""
	for _, hunk := range file.Hunks {
		for _, line := range hunk.Lines {
			if line.Type == godiffy.HunkLineAdded || line.Type == godiffy.HunkLineContext {
				content += line.Content
			}
		}
	}
	fileMode, err := strconv.ParseInt(file.NewMode, 8, 0)
	if err != nil {
		return fmt.Errorf("failed to convert file mode %s: %w", file.NewMode, err)
	}
	err = os.WriteFile(filepath.Join(path, file.NewPath), []byte(content), os.FileMode(fileMode))
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", file.NewPath, err)
	}

	return nil
}
