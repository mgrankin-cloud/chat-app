package validator

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
)

func ValidateField(field, fieldName string, validator func(string) (bool, error)) error {
	if field == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("%s is required", fieldName))
	}
	valid, err := validator(field)
	if err != nil {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("invalid %s", fieldName))
	}
	if !valid {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("invalid %s format", fieldName))
	}
	return nil
}

func IsValidEmail(email string) (bool, error) {
	if len(email) < 8 {
		return false, status.Error(codes.InvalidArgument, "email is too short, minimum 8 characters required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return false, status.Error(codes.InvalidArgument, "invalid email format")
	}

	return true, nil
}

func IsValidPassword(password string) (bool, error) {
	if len(password) < 8 {
		return false, status.Error(codes.InvalidArgument, "password is too short, minimum 8 characters required")
	}

	passwordRegex := regexp.MustCompile(`^[a-zA-Z0-9]{6,}$`)
	if !passwordRegex.MatchString(password) {
		return false, status.Error(codes.InvalidArgument, "invalid password format")
	}
	return true, nil
}

func IsValidUsername(username string) (bool, error) {
	if len(username) < 8 {
		return false, status.Error(codes.InvalidArgument, "username is too short, minimum 8 characters")
	}

	userRegex := regexp.MustCompile(`^[a-zA-Z0-9]{6,}$`)
	if !userRegex.MatchString(username) {
		return false, status.Error(codes.InvalidArgument, "invalid username format")
	}
	return true, nil
}

func IsValidPhoneNumber(phoneNumber string) (bool, error) {
	if len(phoneNumber) < 7 {
		return false, status.Error(codes.InvalidArgument, "phone number is too short")
	}

	return true, nil
}
