package inputmetrics

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/volvofixthis/go-livecheck/internal/clients"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

func FilterFilesWithRegexp(files []string, r string) ([]string, error) {
	filesF := make([]string, 0, len(files))
	for _, file := range files {
		base := filepath.Base(file)
		found, err := regexp.MatchString(r, base)
		if err != nil {
			color.Red("Regexp is wrong: %s", r)
			return nil, err
		}
		if found {
			filesF = append(filesF, file)
		}
	}
	return filesF, nil
}

func GetMetricsStdinData(d Decoder) (map[string]interface{}, error) {
	data, err := d(os.Stdin)
	if err != nil {
		color.Red("Can't decode json in metrics stdin: %s", err)
		return nil, err
	}
	return data, nil
}

func GetMetricsFileData(path string, d Decoder) (map[string]interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		color.Red("Can't open metrics file: %s, %s", path, err)
		return nil, err
	}
	data, err := d(f)
	if err != nil {
		color.Red("Can't decode json in metrics file: %s, %s", path, err)
		return nil, err
	}
	return data, nil
}

func GetMetricsURLData(url string, d Decoder) (map[string]interface{}, error) {
	httpClient := clients.GetHTTPClient()
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	data, err := d(resp.Body)
	if err != nil {
		color.Red("Can't decode json in metrics url: %s, %s", url, err)
		return nil, err
	}
	return data, nil
}

type Decoder func(r io.Reader) (map[string]interface{}, error)

func JSONDecoder(r io.Reader) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := json.NewDecoder(r).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func YAMLDecoder(r io.Reader) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := yaml.NewDecoder(r).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
