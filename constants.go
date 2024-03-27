package annkesdk

import "fmt"

const (
	timeout      = 5
	randomLenght = 100000000

	loginPath           = "/ISAPI/Security/sessionLogin/capabilities"
	sessionPath         = "/ISAPI/Security/sessionLogin"
	inputChannelsPath   = "/ISAPI/System/Video/inputs/channels"
	motionDetectionPath = "/ISAPI/System/Video/inputs/channels/%d/motionDetection"
	motionSchedule      = "/ISAPI/Event/schedules/motionDetections/VMD_video%d"
	eventTrigger        = "/ISAPI/Event/triggers/VMD-%d"
)

func getMotionDetectionPath(channel int) string {
	return fmt.Sprintf(motionDetectionPath, channel)
}

func getMotionSchedulePath(channel int) string {
	return fmt.Sprintf(motionSchedule, channel)
}

func getEventTriggerPath(channel int) string {
	return fmt.Sprintf(eventTrigger, channel)
}
