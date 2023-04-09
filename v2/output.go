package testmtx

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

// TODO
// - 出力するディレクトリ構成等を考慮し、適切にファイルを作成する
// - ファイル自体はformatterに任せる

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
		dir := filepath.Join(baseDir, string(propName))
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

				fp := filepath.Join(dir, fmt.Sprintf("%s_%s.%s", sheet.Name, cn, o.formatter.Extension()))

				file, err := os.Create(fp)
				if err != nil {
					return fmt.Errorf("fail to create file: %w", err)
				}
				defer file.Close()

				o.formatter.Write(file, col, cn, 0)

				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			return err
		}
	}

	return nil
}
