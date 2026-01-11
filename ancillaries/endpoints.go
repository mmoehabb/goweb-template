package ancillaries

import (
	"log"
	"os"
	"path"
	"strings"
)

func GetEndpoint(dir string) []string {
	var endpoints []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Println("WARNING:", dir,"directory not found!")
	}

	for _, entry := range entries {
		if (strings.HasSuffix(entry.Name(), "_templ.go")) {
			var endpoint = strings.TrimSuffix(entry.Name(), "_templ.go")
			if strings.Compare(endpoint, "index") == 0 {
				endpoints = append(endpoints, "/")
			}
			endpoints = append(endpoints, endpoint)
			continue
		}
		if (entry.IsDir()) {
			eps := GetEndpoint(path.Join(dir, entry.Name()))
			var inner []string
			for _, ep := range eps {
				inner = append(inner, path.Join(entry.Name(), ep))
			}
			endpoints = append(endpoints, inner...)
		}
	}

	return endpoints
}
