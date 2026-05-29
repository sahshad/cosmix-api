package rabbitmq

const (

	// Auth Routing Keys
	AuthUserEmailVerificationRequested string = "auth.user.email_verification_requested"
	AuthUserEmailVerificationCompleted string = "auth.user.email_verification_completed"
	AuthUserEmailUpdated               string = "auth.user.email_updated"
	AuthUserForgotPasswordRequested    string = "auth.user.forgot_password_requested"
	AuthUserPasswordChanged            string = "auth.user.password_changed"

	// User Routing Keys
	UserProfileUpdated string = "user.profile.updated"
	UserFollowed       string = "user.followed"
	UserUnfollowed     string = "user.unfollowed"

	// Post Routing Keys
	PostCreated string = "post.created"
	PostUpdated string = "post.updated"
	PostDeleted string = "post.deleted"
	PostLiked   string = "post.liked"
	PostUnliked string = "post.unliked"

	// Comment Routing Keys
	CommentCreated string = "comment.created"
	CommentUpdated string = "comment.updated"
	CommentDeleted string = "comment.deleted"

	// Chat Routing Keys
	ChatMessageSent string = "chat.message.sent"
)
