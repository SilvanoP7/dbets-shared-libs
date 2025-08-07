package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SilvanoP7/dbets-shared-libs/models"
)

// Client handles communication with the events service
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new events service client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateSport creates a new sport in the events service and returns the created sport
func (c *Client) CreateSport(sport models.Sport) (*models.Sport, error) {
	url := fmt.Sprintf("%s/api/v1/sports", c.baseURL)

	jsonData, err := json.Marshal(sport)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sport: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("events service request failed with status: %d", resp.StatusCode)
	}

	// Parse the response to get the created sport
	var response models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("events service returned error: %s", response.Error)
	}

	// Parse the sport from the response
	sportData, err := json.Marshal(response.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sport data: %w", err)
	}

	var createdSport models.Sport
	if err := json.Unmarshal(sportData, &createdSport); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sport: %w", err)
	}

	return &createdSport, nil
}

// CreateEvent creates a new event in the events service and returns the created event
func (c *Client) CreateEvent(event models.CreateEventRequest) (*models.Event, error) {
	url := fmt.Sprintf("%s/api/v1/events", c.baseURL)

	jsonData, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("events service request failed with status: %d", resp.StatusCode)
	}

	// Parse the response to get the created event
	var response models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("events service returned error: %s", response.Error)
	}

	// Parse the event from the response
	eventData, err := json.Marshal(response.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event data: %w", err)
	}

	var createdEvent models.Event
	if err := json.Unmarshal(eventData, &createdEvent); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event: %w", err)
	}

	return &createdEvent, nil
}

// CreateMarket creates a new market in the events service and returns the created market
func (c *Client) CreateMarket(market models.CreateMarketRequest) (*models.Market, error) {
	url := fmt.Sprintf("%s/api/v1/markets", c.baseURL)

	jsonData, err := json.Marshal(market)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal market: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("events service request failed with status: %d", resp.StatusCode)
	}

	// Parse the response to get the created market
	var response models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("events service returned error: %s", response.Error)
	}

	// Parse the market from the response
	marketData, err := json.Marshal(response.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal market data: %w", err)
	}

	var createdMarket models.Market
	if err := json.Unmarshal(marketData, &createdMarket); err != nil {
		return nil, fmt.Errorf("failed to unmarshal market: %w", err)
	}

	return &createdMarket, nil
}

// CreateSelection creates a new selection in the events service and returns the created selection
func (c *Client) CreateSelection(selection models.CreateSelectionRequest) (*models.Selection, error) {
	url := fmt.Sprintf("%s/api/v1/selections", c.baseURL)

	jsonData, err := json.Marshal(selection)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal selection: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("events service request failed with status: %d", resp.StatusCode)
	}

	// Parse the response to get the created selection
	var response models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("events service returned error: %s", response.Error)
	}

	// Parse the selection from the response
	selectionData, err := json.Marshal(response.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal selection data: %w", err)
	}

	var createdSelection models.Selection
	if err := json.Unmarshal(selectionData, &createdSelection); err != nil {
		return nil, fmt.Errorf("failed to unmarshal selection: %w", err)
	}

	return &createdSelection, nil
}

// UpdateOdds updates odds for a specific selection
func (c *Client) UpdateOdds(selectionID string, odds models.UpdateOddsRequest) error {
	url := fmt.Sprintf("%s/api/v1/selections/%s/odds", c.baseURL, selectionID)

	jsonData, err := json.Marshal(odds)
	if err != nil {
		return fmt.Errorf("failed to marshal odds: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("events service request failed with status: %d", resp.StatusCode)
	}

	return nil
}

// GetSportByName retrieves a sport by name from the events service
func (c *Client) GetSportByName(name string) (*models.Sport, error) {
	url := fmt.Sprintf("%s/api/v1/sports?name=%s", c.baseURL, name)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("events service request failed with status: %d", resp.StatusCode)
	}

	var response models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("events service returned error: %s", response.Error)
	}

	// Parse the sports array from the response
	sportData, err := json.Marshal(response.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sport data: %w", err)
	}

	var sport models.Sport
        if err := json.Unmarshal(sportData, &sport); err == nil {
                return &sport, nil
        }

        var sports []models.Sport
	if err := json.Unmarshal(sportData, &sports); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sport data: %w", err)
	}

	// Return the first sport if any found
	if len(sports) == 0 {
		return nil, fmt.Errorf("no sport found with name: %s", name)
	}

	return &sports[0], nil
}

// GetSportByExternalKey retrieves a sport by external key from the events service
func (c *Client) GetSportByExternalKey(externalKey string) (*models.Sport, error) {
	url := fmt.Sprintf("%s/api/v1/sports?external_key=%s", c.baseURL, externalKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("events service request failed with status: %d", resp.StatusCode)
	}

	var response models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("events service returned error: %s", response.Error)
	}

	// Parse the sports array from the response
	sportData, err := json.Marshal(response.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sport data: %w", err)
	}

	var sport models.Sport
        if err := json.Unmarshal(sportData, &sport); err == nil {
                return &sport, nil
        }

        var sports []models.Sport
	if err := json.Unmarshal(sportData, &sports); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sport data: %w", err)
	}

	// Return the first sport if any found
	if len(sports) == 0 {
		return nil, fmt.Errorf("no sport found with external key: %s", externalKey)
	}

	return &sports[0], nil
}

// GetEventByTitle retrieves an event by title from the events service
func (c *Client) GetEventByTitle(title string) (*models.Event, error) {
	url := fmt.Sprintf("%s/api/v1/events?title=%s", c.baseURL, title)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("events service request failed with status: %d", resp.StatusCode)
	}

	var response models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("events service returned error: %s", response.Error)
	}

	// Parse the events array from the response
	eventData, err := json.Marshal(response.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event data: %w", err)
	}

	var events []models.Event
	if err := json.Unmarshal(eventData, &events); err != nil {
		return nil, fmt.Errorf("failed to unmarshal events array: %w", err)
	}

	// Return the first event if any found
	if len(events) == 0 {
		return nil, fmt.Errorf("no event found with title: %s", title)
	}

	return &events[0], nil
}

// GetEventByExternalID retrieves an event by external ID from the events service
func (c *Client) GetEventByExternalID(externalID string) (*models.Event, error) {
	url := fmt.Sprintf("%s/api/v1/events?external_id=%s", c.baseURL, externalID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("events service request failed with status: %d", resp.StatusCode)
	}

	var response models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("events service returned error: %s", response.Error)
	}

	// Parse the events array from the response
	eventData, err := json.Marshal(response.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event data: %w", err)
	}

	var events []models.Event
	if err := json.Unmarshal(eventData, &events); err != nil {
		return nil, fmt.Errorf("failed to unmarshal events array: %w", err)
	}

	// Return the first event if any found
	if len(events) == 0 {
		return nil, fmt.Errorf("no event found with external ID: %s", externalID)
	}

	return &events[0], nil
}

// GetMarketByName retrieves a market by name from the events service
func (c *Client) GetMarketByName(name string, eventID string) (*models.Market, error) {
	url := fmt.Sprintf("%s/api/v1/markets?name=%s&event_id=%s", c.baseURL, name, eventID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("events service request failed with status: %d", resp.StatusCode)
	}

	var response models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("events service returned error: %s", response.Error)
	}

	// Parse the markets array from the response
	marketData, err := json.Marshal(response.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal market data: %w", err)
	}

	var markets []models.Market
	if err := json.Unmarshal(marketData, &markets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal markets array: %w", err)
	}

	// Return the first market if any found
	if len(markets) == 0 {
		return nil, fmt.Errorf("no market found with name: %s", name)
	}

	return &markets[0], nil
}

// GetMarketByExternalKey retrieves a market by external key from the events service
func (c *Client) GetMarketByExternalKey(externalKey string, eventID string) (*models.Market, error) {
	url := fmt.Sprintf("%s/api/v1/markets?external_key=%s&event_id=%s", c.baseURL, externalKey, eventID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("events service request failed with status: %d", resp.StatusCode)
	}

	var response models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("events service returned error: %s", response.Error)
	}

	// Parse the markets array from the response
	marketData, err := json.Marshal(response.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal market data: %w", err)
	}

	var markets []models.Market
	if err := json.Unmarshal(marketData, &markets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal markets array: %w", err)
	}

	// Return the first market if any found
	if len(markets) == 0 {
		return nil, fmt.Errorf("no market found with external key: %s", externalKey)
	}

	return &markets[0], nil
}

// GetSelectionByName retrieves a selection by name from the events service
func (c *Client) GetSelectionByName(name string, marketID string) (*models.Selection, error) {
	url := fmt.Sprintf("%s/api/v1/selections?name=%s&market_id=%s", c.baseURL, name, marketID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("events service request failed with status: %d", resp.StatusCode)
	}

	var response models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("events service returned error: %s", response.Error)
	}

	// Parse the selections array from the response
	selectionData, err := json.Marshal(response.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal selection data: %w", err)
	}

	var selections []models.Selection
	if err := json.Unmarshal(selectionData, &selections); err != nil {
		return nil, fmt.Errorf("failed to unmarshal selections array: %w", err)
	}

	// Return the first selection if any found
	if len(selections) == 0 {
		return nil, fmt.Errorf("no selection found with name: %s", name)
	}

	return &selections[0], nil
}

// GetSelectionByExternalID retrieves a selection by external ID from the events service
func (c *Client) GetSelectionByExternalID(externalID string, marketID string) (*models.Selection, error) {
	url := fmt.Sprintf("%s/api/v1/selections?external_id=%s&market_id=%s", c.baseURL, externalID, marketID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("events service request failed with status: %d", resp.StatusCode)
	}

	var response models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("events service returned error: %s", response.Error)
	}

	// Parse the selections array from the response
	selectionData, err := json.Marshal(response.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal selection data: %w", err)
	}

	var selections []models.Selection
	if err := json.Unmarshal(selectionData, &selections); err != nil {
		return nil, fmt.Errorf("failed to unmarshal selections array: %w", err)
	}

	// Return the first selection if any found
	if len(selections) == 0 {
		return nil, fmt.Errorf("no selection found with external ID: %s", externalID)
	}

	return &selections[0], nil
}
