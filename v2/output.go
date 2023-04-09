package testmtx

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

type Outputter interface {
	Output(baseDir string, sheet *Sheet) error
}

type oneColumnOneCaseOutputter struct {
	formatter Formatter
}

func NewOneColumnOneCaseOutputter(formatter Formatter) Outputter {
	return &oneColumnOneCaseOutputter{
		formatter: formatter,
	}
}

func (o *oneColumnOneCaseOutputter) Output(baseDir string, sheet *Sheet) error {

	if sheet == nil {
		return errors.New("sheet is nil")
	}

	for propName, col := range sheet.Collections {
		eg := &errgroup.Group{}
		for _, cn := range sheet.ColumnNames {
			cn := cn
			eg.Go(func() (er error) {
				defer func() {
					if e := recover(); e != nil {
						er = fmt.Errorf("panic recovered in goroutine: %v", e)
					}
				}()

				buf := &bytes.Buffer{}
				o.formatter.Write(buf, col, cn, 0)
				if buf.String() == "" {
					return nil
				}

				dir := filepath.Join(baseDir, sheet.Name, string(cn))
				if err := os.MkdirAll(dir, 0777); err != nil {
					return fmt.Errorf("fail to create directory: %w", err)
				}

				fp := filepath.Join(dir, fmt.Sprintf("%s.%s", propName, o.formatter.Extension()))
				file, err := os.Create(fp)
				if err != nil {
					return fmt.Errorf("fail to create file: %w", err)
				}
				defer file.Close()

				fmt.Fprint(file, buf.String())

				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			return err
		}
	}

	return nil
}

type oneSheetOneCaseOutputter struct {
	formatter Formatter
}

func NewOneSheetOneCaseOutputter(formatter Formatter) Outputter {
	return &oneSheetOneCaseOutputter{
		formatter: formatter,
	}
}

func (o *oneSheetOneCaseOutputter) Output(baseDir string, sheet *Sheet) error {

	if sheet == nil {
		return errors.New("sheet is nil")
	}

	for propName, col := range sheet.Collections {
		dir := filepath.Join(baseDir, sheet.Name, string(propName))
		if err := os.MkdirAll(dir, 0777); err != nil {
			return fmt.Errorf("fail to create directory: %w", err)
		}
		eg := &errgroup.Group{}
		for _, cn := range sheet.ColumnNames {
			cn := cn
			eg.Go(func() (er error) {
				defer func() {
					if e := recover(); e != nil {
						er = fmt.Errorf("panic recovered in goroutine: %v", e)
					}
				}()

				buf := &bytes.Buffer{}
				o.formatter.Write(buf, col, cn, 0)
				if buf.String() == "" {
					return nil
				}

				fp := filepath.Join(dir, fmt.Sprintf("%s.%s", cn, o.formatter.Extension()))
				file, err := os.Create(fp)
				if err != nil {
					return fmt.Errorf("fail to create file: %w", err)
				}
				defer file.Close()

				fmt.Fprint(file, buf.String())

				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			return err
		}
	}

	return nil
}
