package initialisers

import "jwt/models"

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}
