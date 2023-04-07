package lib

import (
	"github.com/flosch/pongo2/v6"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

// Generate applies a set of user defined options (ie: the 'context') to a set of template
// files stored in srcPath, and produces a complete project in the targetPath with the
// user defined parameters applied throughout
func Generate(srcPath string, targetPath string, context map[string]any) error {
	// loop through all files
	var err = filepath.WalkDir(srcPath, func(path string, info os.DirEntry, err error) error {
		// If walk encountered an error attempting to enumerate the file system object
		// we are processing, it tells us here. For now we just assume we can not proceed
		// if we hit this condition.
		// TODO: Consider how best to handle error conditions
		//		https://pkg.go.dev/io/fs#WalkDirFunc
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return errors.Wrap(err, "Failed to parse relative Path name")
		}
		// Skip processing the root dir
		if relPath == "." {
			return nil
		}
		// Skip Rejigger manifest file
		if relPath == ".rejig.yml" {
			return nil
		}

		// apply template to the Path being processed
		tpl, err := pongo2.FromString(relPath)
		if err != nil {
			return errors.Wrap(err, "Failed to load Path as template")
		}
		newDirName, err := tpl.Execute(context)
		if err != nil {
			return errors.Wrap(err, "Failed to apply template to Path")
		}
		newOutputPath := filepath.Join(targetPath, newDirName)

		temp, err := info.Info()
		if err != nil {
			return errors.Wrap(err, "Failed to load file system info for input file")
		}
		originalMode := temp.Mode()

		// Generate output content
		if info.IsDir() {
			// Make sure to preserve the file mode
			err = os.MkdirAll(newOutputPath, originalMode)
			if err != nil {
				return errors.Wrap(err, "Failed to create output folder")
			}
		} else {
			// Apply template to the file contents
			var data []byte
			data, err = os.ReadFile(path)
			if err != nil {
				return errors.Wrap(err, "Failed to read source file")
			}
			tpl, err = pongo2.FromString(string(data))
			if err != nil {
				return errors.Wrap(err, "Failed to load template source")
			}

			var newData string
			newData, err = tpl.Execute(context)
			if err != nil {
				return errors.Wrap(err, "Failed to apply template")
			}

			// Create a new output file with the processed content
			// making sure to preserve the file mode in the process
			err = os.WriteFile(newOutputPath, []byte(newData), originalMode)
			if err != nil {
				return errors.Wrap(err, "Failed to create output file")
			}
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "Failed to generate project")
	}
	return nil
}
