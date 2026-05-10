package rabbitmq

const (
	AuthUserRegistered         = "auth.user.registered"
	AuthUserEmailUpdated       = "auth.user.email.updated"
	AuthPasswordResetRequested = "auth.password.reset.requested"
	AuthPasswordChanged        = "auth.password.changed"
)

const (
	UserProfileUpdated = "user.profile.updated"
	UserFollowed       = "user.followed"
	UserUnfollowed     = "user.unfollowed"

	UserAuthUserRegistered = "user.auth.user.registered"
)

const (
	PostCreated = "post.created"
	PostUpdated = "post.updated"
	PostDeleted = "post.deleted"
	PostLiked   = "post.liked"
	PostUnliked = "post.unliked"
)

const (
	CommentCreated = "comment.created"
	CommentUpdated = "comment.updated"
	CommentDeleted = "comment.deleted"
)

const (
	NotificationCreated = "notification.created"

	NotificationAuthUserRegistered = "notification.auth.user.registered"
)

const (
	ChatMessageSent = "chat.message.sent"
)
