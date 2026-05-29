package rabbitmq

const (

	// User Service Queues
	UserAuthUserEmailVerificationCompletedQueue string = "user.auth.user.email_verification_completed"

	// Notification Service Queues
	NotificationAuthUserEmailVerificationRequestedQueue string = "notification.auth.user.email_verification_requested"
	NotificationAuthUserEmailVerificationCompletedQueue string = "notification.auth.user.email_verification_completed"
	NotificationAuthUserForgotPasswordRequestedQueue    string = "notification.auth.user.forgot_password_requested"
	NotificationAuthUserPasswordChangedQueue            string = "notification.auth.user.password_changed"
	NotificationAuthUserRegisteredQueue                 string = "notification.auth.user.registered"
	NotificationAuthUserEmailVerificationQueue          string = "notification.auth.user.email_verification"

	// Post Service Queues
	PostAuthUserEmailVerificationCompletedQueue string = "post.auth.user.email_verification_completed"
)
