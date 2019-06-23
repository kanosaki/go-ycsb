package worker

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path"

	"github.com/magiconair/properties"
	"google.golang.org/api/people/v1"
)

type Client struct {
	addr   string
	client *http.Client
}

func NewClient(workerAddr string) *Client {

}

func (c *Client) ListJobs() ([]*Job, error) {
}

func (c *Client) JobStatus(jobId string) (*Job, error) {
}

func (c *Client) StartJob(jobId string, props *properties.Properties, dbName, workload string) error {
	var reqBody bytes.Buffer
	writer := multipart.NewWriter(&reqBody)
	w, err := writer.CreateFormField("properties")
	if err != nil {
		return err
	}
	if _, err := props.Write(w, properties.UTF8); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", path.Join(c.addr, "/job/run", jobId), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		var e errorResponse
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return err
		}
		return errors.New(e.Message)
	}
	return nil
}

func (c *Client) DownloadFile(jobId, key string) (io.ReadCloser, error) {
	req, err := http.NewRequest("POST", path.Join(c.addr, "/job/run", jobId), nil)
	if err != nil {
		return err
	}
}
