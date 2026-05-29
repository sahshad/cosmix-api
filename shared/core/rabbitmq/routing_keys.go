package rabbitmq

const (
	AuthUserEmailVerification  = "auth.user.email.verification"
	AuthUserRegistered         = "auth.user.registered"
	AuthUserEmailUpdated       = "auth.user.email.updated"
	AuthPasswordResetRequested = "auth.password.reset.requested"
	AuthPasswordChanged        = "auth.password.changed"
)

const (
	UserProfileUpdated     = "user.profile.updated"
	UserFollowed           = "user.followed"
	UserUnfollowed         = "user.unfollowed"
	UserAuthUserRegistered = "user.auth.user.registered"
)

const (
	PostAuthUserRegistered = "post.auth.user.registered"
	PostCreated            = "post.created"
	PostUpdated            = "post.updated"
	PostDeleted            = "post.deleted"
	PostLiked              = "post.liked"
	PostUnliked            = "post.unliked"

	CommentCreated = "comment.created"
	CommentUpdated = "comment.updated"
	CommentDeleted = "comment.deleted"
)

const (
	NotificationAuthUserRegistered = "notification.auth.user.registered"
	NotificationCreated            = "notification.created"
	NotificationAuthUserEmailVerification = "notification.auth.user.email.verification"
)

const (
	ChatMessageSent = "chat.message.sent"
)
