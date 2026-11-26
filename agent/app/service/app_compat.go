package service

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

)

// MergeAppEnvironment flattens install parameters into environment variables.
func MergeAppEnvironment(params map[string]interface{}, envParams map[string]string) {
	for k, v := range params {
		switch t := v.(type) {
		case string:
			envParams[k] = t
		case float64:
			envParams[k] = strconv.FormatFloat(t, 'f', -1, 32)
		case uint:
			envParams[k] = strconv.Itoa(int(t))
		case int:
			envParams[k] = strconv.Itoa(t)
		case []interface{}:
			strArray := make([]string, len(t))
			for i := range t {
				strArray[i] = strings.ToLower(fmt.Sprintf("%v", t[i]))
			}
			envParams[k] = strings.Join(strArray, ",")
		case map[string]interface{}:
			MergeAppEnvironment(t, envParams)
		}
	}
}

func appDownloadURL(rawURL, version string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	query := parsedURL.Query()
	query.Set("v", version)
	parsedURL.RawQuery = query.Encode()
	return parsedURL.String(), nil
}

func currentAppDownloadURL(rawURL string) (string, error) {
	storeVersion, err := settingRepo.GetValueByKey("AppStoreLastModified")
	if err != nil {
		return "", err
	}
	return appDownloadURL(rawURL, storeVersion)
}
