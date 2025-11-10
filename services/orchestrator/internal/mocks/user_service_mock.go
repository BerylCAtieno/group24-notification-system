package mocks

import (
	"time"

	"github.com/BerylCAtieno/group24-notification-system/services/orchestrator/internal/models"
)

type UserServiceMock struct{}

func NewUserServiceMock() *UserServiceMock {
	return &UserServiceMock{}
}

func (m *UserServiceMock) GetPreferences(userID string) (*models.UserPreferences, error) {
	// Return mock user preferences
	return &models.UserPreferences{
		UserID:              userID,
		Email:               "user@example.com",
		Phone:               "+254712345678",
		Timezone:            "Africa/Nairobi",
		Language:            "en",
		NotificationEnabled: true,
		Channels: models.Channels{
			Email: models.EmailChannel{
				Enabled:   true,
				Verified:  true,
				Frequency: "immediate",
				QuietHours: models.QuietHours{
					Enabled:  true,
					Start:    "22:00",
					End:      "07:00",
					Timezone: "Africa/Nairobi",
				},
			},
			Push: models.PushChannel{
				Enabled: true,
				Devices: []models.Device{
					{
						DeviceID: "dev_abc123",
						Platform: "ios",
						Token:    "fcm_token_xyz...",
						LastSeen: time.Now().Add(-1 * time.Hour),
						Active:   true,
					},
					{
						DeviceID: "dev_def456",
						Platform: "android",
						Token:    "fcm_token_abc...",
						LastSeen: time.Now().Add(-24 * time.Hour),
						Active:   true,
					},
				},
				QuietHours: models.QuietHours{
					Enabled: false,
				},
			},
		},
		Preferences: models.NotificationPrefs{
			Marketing:     false,
			Transactional: true,
			Reminders:     true,
			Digest: models.Digest{
				Enabled:   true,
				Frequency: "daily",
				Time:      "09:00",
			},
		},
		UpdatedAt: time.Now(),
	}, nil
}

func (m *UserServiceMock) GetOptOutStatus(userID string) (*models.OptOutStatus, error) {
	return &models.OptOutStatus{
		UserID:   userID,
		OptedOut: false,
		Channels: map[string]bool{
			"email": false,
			"push":  false,
		},
		CheckedAt: time.Now(),
	}, nil
}
