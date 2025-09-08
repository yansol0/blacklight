package tester

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/yansol0/blacklight/parser"
	"github.com/yansol0/blacklight/utils"
)

type Results struct {
	Unauth         [][2]string
	Auth           [][2]string
	Bypass         [][2]string
	IDORCandidates []string
	BypassHits     []string
}

var bypassHeaders = []map[string]string{
	{"X-Forwarded-For": "127.0.0.1"},
	{"X-Originating-IP": "127.0.0.1"},
	{"X-Client-IP": "127.0.0.1"},
	{"X-Remote-IP": "127.0.0.1"},
	{"Forwarded": "for=127.0.0.1"},
	{"Authorization": "Bearer faketoken"},
	{"Cookie": "session=fakecookie"},
}

func RunTests(endpoints []parser.Endpoint, token string, cookie string) Results {
	results := Results{
		Unauth:         make([][2]string, 0),
		Auth:           make([][2]string, 0),
		Bypass:         make([][2]string, 0),
		IDORCandidates: make([]string, 0),
		BypassHits:     make([]string, 0),
	}

	utils.LogInfo("Starting probes against API...")
	utils.LogInfo("Total endpoints discovered: " + fmt.Sprint(len(endpoints)))

	client := &http.Client{Timeout: 5 * time.Second}

	var headersAuth map[string]string
	if token != "" {
		headersAuth = map[string]string{"Authorization": "Bearer " + token}
		utils.LogInfo("Using JWT-based authentication")
	} else if cookie != "" {
		headersAuth = map[string]string{"Cookie": cookie}
		utils.LogInfo("Using Cookie-based authentication")
	}

	for i, ep := range endpoints {
		progress := fmt.Sprintf("[%d/%d]", i+1, len(endpoints))
		utils.LogInfo(progress + " Testing endpoint: " + ep.Method + " " + ep.URL)

		unauthStatus := doRequest(client, ep.Method, ep.URL, nil)
		utils.LogInfo("  [Unauth] → " + unauthStatus)
		results.Unauth = append(results.Unauth, [2]string{ep.URL, unauthStatus})

		authStatus := doRequest(client, ep.Method, ep.URL, headersAuth)
		utils.LogInfo("  [Auth]   → " + authStatus)
		results.Auth = append(results.Auth, [2]string{ep.URL, authStatus})

		for _, hset := range bypassHeaders {
			status := doRequest(client, ep.Method, ep.URL, hset)
			key := ""
			for k := range hset {
				key = k
			}
			utils.LogInfo("  [Bypass " + key + "] → " + status)
			results.Bypass = append(results.Bypass, [2]string{ep.URL + " (" + key + ")", status})

			if status != unauthStatus {
				hit := fmt.Sprintf("%s %s (%s) → baseline=%s, bypass=%s",
					ep.Method, ep.URL, key, unauthStatus, status)
				results.BypassHits = append(results.BypassHits, hit)
			}
		}

		if containsIDORHint(ep.Path) {
			utils.LogWarn("Potential IDOR candidate: " + ep.URL)
			results.IDORCandidates = append(results.IDORCandidates, ep.URL)
		}
	}

	utils.LogSuccess("Finished probing all endpoints")

	if len(results.BypassHits) > 0 {
		utils.LogCritical("====================================")
		for _, hit := range results.BypassHits {
			utils.LogCritical("AUTH BYPASS FOUND - " + hit)
		}
		utils.LogCritical("====================================")
	} else {
		utils.LogInfo("No auth bypasses detected")
	}

	return results
}

func doRequest(client *http.Client, method, url string, headers map[string]string) string {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "ERR"
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "ERR"
	}
	defer resp.Body.Close()
	return fmt.Sprintf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
}

func containsIDORHint(path string) bool {
	lower := strings.ToLower(path)

	if strings.Contains(lower, "{") && strings.Contains(lower, "}") {
		return true
	}

	idorKeywords := []string{"user", "account", "project", "org", "team", "profile"}
	for _, kw := range idorKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}

	return false
}
