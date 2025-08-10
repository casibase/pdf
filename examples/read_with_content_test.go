package examples

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/casibase/pdf"
)

func TestGetContentText(t *testing.T) {
	text, err := getTextFromPdf("./pdf_test.pdf")
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}

func getTextFromPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	totalPage := r.NumPage()
	var mergedTexts []string
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		var lastTextStyle pdf.Text
		var mergedSentence string

		var texts []pdf.Text
		texts, err = getPageTexts(p)
		if err != nil {
			return "", err
		}
		defer f.Close()

		for _, text := range texts {
			if text.Y == lastTextStyle.Y {
				mergedSentence += text.S
			} else {
				if mergedSentence != "" {
					mergedTexts = append(mergedTexts, mergedSentence)
				}
				lastTextStyle = text
				mergedSentence = text.S
			}
		}

		if mergedSentence != "" {
			mergedTexts = append(mergedTexts, mergedSentence)
		}
	}

	mergedText := strings.Join(mergedTexts, "\n")
	return mergedText, nil
}

func getPageTexts(p pdf.Page) (texts []pdf.Text, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New(fmt.Sprint(x))
			}
		}
	}()

	texts = p.Content().Text
	return
}
