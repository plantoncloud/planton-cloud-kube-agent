package client

import (
	"encoding/json"
	"fmt"
	"github.com/opencost/opencost/pkg/costmodel"
	"github.com/opencost/opencost/pkg/kubecost"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseUrl string
}

// GetAllocationAggregatedByPodByWindow retrieves data from 'http://localhost:9003/allocation/compute?aggregate=namespace&window=1h&aggregatedMetadata=true'
func (c *Client) GetAllocationAggregatedByPodByWindow(window string) (map[string]*kubecost.AllocationJSON, error) {
	log.Infof("open-cost scrapper triggered at %v for window %v", time.Now(), window)
	url := fmt.Sprintf("%s/allocation/compute?aggregate=namespace&window=%s&aggregatedMetadata=true", c.BaseUrl, window)

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "open-cost request failed using url %s", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(err, "open-cost request failed using url %s. unexpected status code: %d",
			url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read body for open-cost request using url %s", url)
	}

	var costModelResp *costmodel.Response
	if err := json.Unmarshal(body, &costModelResp); err != nil {
		return nil, errors.Wrapf(err, "failed to json unmarshal body for open-cost request using url %s", url)
	}

	jsonData, err := json.Marshal(costModelResp.Data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to json marshal data object for open-cost request using url %s", url)
	}

	listOfWindowObjects := make([]interface{}, 0)
	err = json.Unmarshal(jsonData, &listOfWindowObjects)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to json unmarshal list of window objects for open-cost request using url %s", url)
	}
	//todo: need to create one entry for each time window
	windowCostAllocationEntryMap := make(map[string]*kubecost.AllocationJSON, 0)
	for _, windowObject := range listOfWindowObjects {
		windowCostAllocationEntryMapJson, err := json.Marshal(windowObject)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to json marshal a window-object for open-cost request using url %s", url)
		}
		err = json.Unmarshal(windowCostAllocationEntryMapJson, &windowCostAllocationEntryMap)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to json unmarshal a window-object for open-cost request using url %s", url)
		}
	}
	return windowCostAllocationEntryMap, nil
}
