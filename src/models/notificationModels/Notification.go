package notificationModels

type Notification struct {
	NotificationId  int    `json:"id"`
	Type            int    `json:"notification_type"`
	Viewed          bool   `json:"viewed"`
	InitiatedById   int    `json:"initiated_by_id"`
	InitiatedByName string `json:"initiated_by_name"`
	RecipientId     int    `json:"recipient_id"`
	RecipientName   string `json:"recipient_name"`
	CauseEntity     int    `json:"cause_entity"`
}
