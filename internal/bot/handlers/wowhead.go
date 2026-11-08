package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const suggestionsURL = "https://www.wowhead.com/forever/search/suggestions-template"

type item struct {
	Type int    `json:"type"`
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (i item) url() string {
	return "https://www.wowhead.com/forever/item=" + strconv.Itoa(i.ID)
}

func lookupItems(client *http.Client, name string) ([]item, error) {
	u, err := url.Parse(suggestionsURL)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("q", name)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "LootGet/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error from WowHead: %d", resp.StatusCode)
	}

	var result struct {
		Results []item `json:"results"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&result); err != nil {
		return nil, fmt.Errorf("Error from WowHead: %w", err)
	}

	var items []item
	seen := make(map[int]bool)

	for _, i := range result.Results {
		if i.Type == 3 && i.ID > 0 && strings.EqualFold(strings.TrimSpace(i.Name), name) && !seen[i.ID] {
			items = append(items, i)
			seen[i.ID] = true
		}
		if len(items) == 5 {
			break
		}
	}
	return items, nil
}
