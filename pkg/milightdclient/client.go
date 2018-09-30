package milightdclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sgrzywna/milightd/internal/app/milightd"
)

// Client represents HTTP client for the milightd daemon.
type Client struct {
	url    string
	client *http.Client
}

// NewClient returns initialized Client object.
func NewClient(url string) *Client {
	return &Client{
		url: url,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// SetLight controls mi-light device through milightd daemon.
func (c *Client) SetLight(l Light) error {
	url := fmt.Sprintf("%s/api/v1/light", c.url)

	var cmd = milightd.Light{
		Color:      l.GetColor(),
		Brightness: l.GetBrightness(),
		Switch:     l.GetSwitch(),
	}

	d, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(d))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("milightd client: unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// GetSequences returns list of defined sequences from milightd daemon.
func (c *Client) GetSequences() error {
	return nil
}

// AddSequence adds sequence through milightd daemon.
func (c *Client) AddSequence() error {
	return nil
}

// GetSequence return sequence definition from milightd daemon.
func (c *Client) GetSequence() error {
	return nil
}

// DeleteSequence deletes sequence through milightd daemon.
func (c *Client) DeleteSequence() error {
	return nil
}

// GetSequenceState returns state of the running sequence from milightd daemon.
func (c *Client) GetSequenceState() error {
	return nil
}

// SetSequenceState control state of the running sequence through milightd daemon.
func (c *Client) SetSequenceState() error {
	return nil
}
