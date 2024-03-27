package models

import "encoding/xml"

type VideoInputChannelList struct {
	XMLName           xml.Name `xml:"VideoInputChannelList"`
	Text              string   `xml:",chardata"`
	Version           string   `xml:"version,attr"`
	VideoInputChannel []struct {
		Text              string `xml:",chardata"`
		Version           string `xml:"version,attr"`
		ID                string `xml:"id"`
		InputPort         string `xml:"inputPort"`
		VideoInputEnabled string `xml:"videoInputEnabled"`
		Name              string `xml:"name"`
		VideoFormat       string `xml:"videoFormat"`
		ResDesc           string `xml:"resDesc"`
	} `xml:"VideoInputChannel"`
}

type MotionDetection struct {
	XMLName          xml.Name `xml:"MotionDetection"`
	Version          string   `xml:"version,attr"`
	Xmlns            string   `xml:"xmlns,attr"`
	Enabled          string   `xml:"enabled"`
	EnableHighlight  string   `xml:"enableHighlight"`
	SamplingInterval string   `xml:"samplingInterval"`
	StartTriggerTime string   `xml:"startTriggerTime"`
	EndTriggerTime   string   `xml:"endTriggerTime"`
	RegionType       string   `xml:"regionType"`
	Grid             struct {
		RowGranularity    string `xml:"rowGranularity"`
		ColumnGranularity string `xml:"columnGranularity"`
	} `xml:"Grid"`
	MotionDetectionLayout struct {
		SensitivityLevel string `xml:"sensitivityLevel"`
		Layout           struct {
			GridMap    string `xml:"gridMap"`
			RegionList struct {
				Size   string `xml:"size,attr"`
				Region struct {
					ID                    string `xml:"id"`
					Xmlns                 string `xml:"xmlns,attr"`
					RegionCoordinatesList struct {
						Size              string `xml:"size,attr"`
						RegionCoordinates []struct {
							PositionX string `xml:"positionX"`
							PositionY string `xml:"positionY"`
						} `xml:"RegionCoordinates"`
					} `xml:"RegionCoordinatesList"`
				} `xml:"Region"`
			} `xml:"RegionList"`
		} `xml:"layout"`
		TargetType string `xml:"targetType"`
	} `xml:"MotionDetectionLayout"`
}

type MotionSchedule struct {
	XMLName             xml.Name `xml:"Schedule"`
	ID                  string   `xml:"id"`
	EventType           string   `xml:"eventType"`
	VideoInputChannelID string   `xml:"videoInputChannelID"`
	TimeBlockList       struct {
		Size      string `xml:"size,attr"`
		TimeBlock []struct {
			DayOfWeek string `xml:"dayOfWeek"`
			TimeRange struct {
				BeginTime string `xml:"beginTime"`
				EndTime   string `xml:"endTime"`
			} `xml:"TimeRange"`
		} `xml:"TimeBlock"`
	} `xml:"TimeBlockList"`
	HolidayBlockList string `xml:"HolidayBlockList"`
}

type EventTrigger struct {
	XMLName                      xml.Name `xml:"EventTrigger"`
	ID                           string   `xml:"id"`
	EventType                    string   `xml:"eventType"`
	VideoInputChannelID          string   `xml:"videoInputChannelID"`
	EventTriggerNotificationList struct {
		EventTriggerNotification []struct {
			ID                 string `xml:"id"`
			NotificationMethod string `xml:"notificationMethod"`
			VideoInputID       string `xml:"videoInputID"`
		} `xml:"EventTriggerNotification"`
	} `xml:"EventTriggerNotificationList"`
}
