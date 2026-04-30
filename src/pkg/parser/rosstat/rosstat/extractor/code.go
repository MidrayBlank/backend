package extractor

import (
	"backend/src/pkg/parser/rosstat/rosstat/dao"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type CodeExtractor struct {
}

func NewCodeExtractor() *CodeExtractor {
	return &CodeExtractor{}
}

func (ext *CodeExtractor) Extract(ctx context.Context, filePath string) ([]*dao.CodeExtracted, error) {
	fileContentByte, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	content := string(fileContentByte)

	codes, err := ext.extractCodes(content)

	if err != nil {
		return nil, err
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	names, err := ext.extractNames(content)

	if err != nil {
		return nil, err
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if len(names) != len(codes) {
		return nil, fmt.Errorf("error different count: names %d, codes %d", len(names), len(codes))
	}

	result := make([]*dao.CodeExtracted, len(names))
	for i := range names {
		result[i] = &dao.CodeExtracted{
			Name: names[i],
			Code: codes[i],
		}
	}

	return result, nil
}

func (_ *CodeExtractor) extractCodes(html string) ([]int, error) {
	re := regexp.MustCompile(`p_oktmo\[(\d+)\]="(\d+)"`)
	matches := re.FindAllStringSubmatch(html, -1)
	if matches == nil {
		log.Printf("CODE_EXTRACTOR_ERROR p_oktmo not found: html=%s\n", html)
		return nil, fmt.Errorf("error p_oktmo not found")
	}

	maxIdx := -1
	for _, m := range matches {
		idx, _ := strconv.Atoi(m[1])
		if idx > maxIdx {
			maxIdx = idx
		}
	}

	codes := make([]int, maxIdx+1)
	for _, m := range matches {
		idx, _ := strconv.Atoi(m[1])
		code, err := strconv.Atoi(m[2])

		if err != nil {
			log.Printf("CODE_EXTRACTOR_ERROR convert code to int: html=%s\n", html)
			return nil, fmt.Errorf("error converting code to int: %w", err)
		}

		codes[idx] = code
	}
	return codes, nil
}

func (_ *CodeExtractor) extractNames(html string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var names []string
	doc.Find("select[name='oktmo'] option").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Text())
		if name != "" {
			names = append(names, name)
		}
	})

	if len(names) == 0 {
		log.Printf("CODE_EXTRACTOR_ERROR <option> not found: html=%s\n", html)
		return nil, fmt.Errorf("error <option> not found")
	}
	return names, nil
}
