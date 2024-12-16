package tests

/*import (
	"Messenger-android/messenger/auth-service/sso/internal/lib/otp"
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSendOTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailSender := NewMockEmailSender(ctrl)
	auth := NewAuth(mockEmailSender, 5*time.Minute)

	email := "test@example.com"
	otpNew := otp.GenerateOTP()

	mockEmailSender.EXPECT().SendEmail(email, "Email Verification OTP", "Your OTP is: "+otpNew).Return(nil)

	err := auth.SendOTP(context.Background(), email)
	assert.NoError(t, err)

	storedOTP, exists := otp.StoreOTP[email]
	assert.True(t, exists)
	assert.Equal(t, otpNew, storedOTP)
}
*/
