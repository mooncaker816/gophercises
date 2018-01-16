package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if newurl, ok := pathsToUrls[req.URL.Path]; ok {
			http.Redirect(w, req, newurl, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, req)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	urlmap, err := parseYml(yml)
	if err != nil {
		return nil, err
	}
	return MapHandler(urlmap, fallback), nil
}

func JsonHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	urlmap, err := parseJson(jsonBytes)
	if err != nil {
		return nil, err
	}
	return MapHandler(urlmap, fallback), nil
}

func parseYml(yml []byte) (map[string]string, error) {
	var urlmap = make(map[string]string)
	var urls []urlpath
	err := yaml.Unmarshal(yml, &urls)
	if err != nil {
		return nil, err
	}
	for _, v := range urls {
		urlmap[v.Path] = v.Url
	}
	return urlmap, nil
}

func parseJson(jsonBytes []byte) (map[string]string, error) {
	var urlmap = make(map[string]string)
	var urls []urlpath
	err := json.Unmarshal(jsonBytes, &urls)
	if err != nil {
		return nil, err
	}
	for _, v := range urls {
		urlmap[v.Path] = v.Url
	}
	return urlmap, nil
}

type urlpath struct {
	Path string `json:"path" ymal:"path"`
	Url  string `json:"url" ymal:"url"`
}
