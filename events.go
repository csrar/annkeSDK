package annkesdk

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/csrar/annkeSDK/models"
)

func (c Connector) GetChannels() (models.VideoInputChannelList, error) {
	inputChannels := models.VideoInputChannelList{}
	err := c.makeGetRequest(inputChannelsPath, &inputChannels)
	return inputChannels, err
}

func (c Connector) GetMotionDetection(channel int) (models.MotionDetection, error) {
	motionDetection := models.MotionDetection{}
	err := c.makeGetRequest(getMotionDetectionPath(channel), &motionDetection)
	return motionDetection, err
}

func (c Connector) GetMotionSchedule(channel int) (models.MotionSchedule, error) {
	motionSchedule := models.MotionSchedule{}
	err := c.makeGetRequest(getMotionSchedulePath(channel), &motionSchedule)
	return motionSchedule, err
}

func (c Connector) GetEventTrigger(channel int) (models.EventTrigger, error) {
	motionTrigger := models.EventTrigger{}
	err := c.makeGetRequest(getEventTriggerPath(channel), &motionTrigger)
	return motionTrigger, err
}

func (c Connector) UpdateMotionDetection(channel int, motion models.MotionDetection) error {
	return c.makeUpdateRequest(getMotionDetectionPath(channel), motion)
}

func (c Connector) UpdateMotionSchedule(channel int, motionSchelude models.MotionSchedule) error {
	return c.makeUpdateRequest(getMotionSchedulePath(channel), motionSchelude)
}

func (c Connector) UpdateEventTrigger(channel int, eventTrigger models.EventTrigger) error {
	return c.makeUpdateRequest(getEventTriggerPath(channel), eventTrigger)
}

func (c Connector) makeUpdateRequest(path string, body interface{}) error {
	url := fmt.Sprintf("%s://%s%s", c.getProtocol(), c.Host, path)
	data, err := xml.Marshal(body)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(data)
	req, err := http.NewRequest("PUT", url, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	response, err := c.client.Do(req)
	if err != nil {
		return err
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return NewAnnkeError(response.StatusCode, string(responseBody), req.RequestURI)
	}
	return nil
}

func (c Connector) makeGetRequest(path string, data interface{}) error {
	url := fmt.Sprintf("%s://%s%s", c.getProtocol(), c.Host, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return NewAnnkeError(resp.StatusCode, string(body), req.RequestURI)
	}
	err = xml.Unmarshal(body, data)

	if err != nil {
		return fmt.Errorf("Error unmarshaling %s response %w", req.RequestURI, err)
	}
	return nil
}
