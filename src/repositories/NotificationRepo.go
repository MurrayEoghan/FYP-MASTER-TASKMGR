package repositories

import (
	"log"
	"repo/sqldb"

	"net/http"
	models "repo/models/notificationModels"
)

func GetNotifications(id int, w http.ResponseWriter) []models.Notification {
	var notifications []models.Notification
	rows, err := sqldb.DB3.Query("SELECT * FROM notifications WHERE recipient_id = ? AND viewed = 1", id)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Println(err)

	}
	defer rows.Close()

	for rows.Next() {
		notification := &models.Notification{}
		if err := rows.Scan(&notification.NotificationId, &notification.Type, &notification.Viewed, &notification.InitiatedById, &notification.InitiatedByName, &notification.RecipientName, &notification.RecipientId, &notification.CauseEntity); err != nil {
			log.Println(err)
		}
		notifications = append(notifications, *notification)

	}
	return notifications
}

func DeleteNotifications(id int, w http.ResponseWriter) (int64, error) {
	result, err := sqldb.DB3.Exec(`DELETE FROM notification_service.notifications WHERE recipient_id = ? AND viewed = 1`, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return 0, err
	} else {
		return result.RowsAffected()
	}
}
