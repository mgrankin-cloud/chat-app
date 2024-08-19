package logic

import (
    "github.com/mgrankin-cloud/user-service/internal/models"
)

func CreateUser(user *models.User) error {
    // написать логику создания пользователя в бд
    return nil
}

func GetUserByID(userID string) (*models.User, error) {
    // написать логику получения пользователя по ID из бд
    return &models.User{
        ID:       userID,
        Username: "example",
        Email:    "example@example.com",
    }, nil
}

func UpdateUser(user *models.User) error {
    // написать логику обновления пользователя в бд
    return nil
}

func DeleteUser(userID string) error {
    // написать логику удаления пользователя из бд
    return nil
}