package main

import (
	"strings"
)

type authorData struct {
	Description string `json:"name"`
	AuthorURL   string `json:"href"`
}

// searchResource loads the cached data and searches for the author.
func (website resource) searchResource(query query, cacheDir string) resource {
	if website.Complex {
		return website.searchComplexResource(query, cacheDir)
	}

	cacheFileName := cacheDir + "/" + website.Name + ".json"
	data := []authorData{}

	update, err := fileNotExist(cacheFileName)
	if update {
		data, err = website.updateCache(cacheDir, cacheFileName)
	}
	if err != nil {
		website.Error = err
		return website
	}

	if !update {
		err = loadFileJSON(cacheFileName, &data)
		if err != nil {
			website.Error = err
			return website
		}
	}

	results := []authorData{}
	for _, a := range data {
		if website.match(a.Description, query) {
			results = append(results, a)
		}
	}
	website.Results = results
	return website
}

func (website resource) match(authorDesc string, query query) bool {
	if !strings.Contains(authorDesc, query.LastName) {
		return false
	}

	if website.FirstName && !strings.Contains(authorDesc, query.FirstName) {
		return false
	}

	return !website.Year || strings.Contains(authorDesc, query.Year)
}

// updateCache carries out an http get request and saves the response body
// into a file
func (website resource) updateCache(cacheDir string, cacheFileName string) ([]authorData, error) {
	fullURL := website.BaseURL + website.QueryURL
	body, err := getResource(fullURL)
	if err != nil {
		return []authorData{}, err
	}

	data, err := website.readResource(body)
	if err != nil {
		return []authorData{}, err
	}
	filteredData := website.dedupe(data)

	return filteredData, writeFileJSON(cacheFileName, filteredData)
}
