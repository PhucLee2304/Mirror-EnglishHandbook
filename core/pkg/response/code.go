package response

const (
	MessageCodeBadRequest   string = "bad-request"
	MessageCodeNotFound     string = "not-found"
	MessageCodeUnauthorized string = "unauthorized"

	MessageCodeInvalidIDToken               string = "invalid-id-token"
	MessageCodeInvalidRefreshToken          string = "invalid-refresh-token"
	MessageCodeFailedToGetFirebaseUser      string = "failed-to-get-firebase-user"
	MessageCodeFailedToLogin                string = "failed-to-login"
	MessageCodeFailedToGenerateAccessToken  string = "failed-to-generate-access-token"
	MessageCodeFailedToGenerateRefreshToken string = "failed-to-generate-refresh-token"
	MessageCodeUserNotFound                 string = "user-not-found"

	MessageCodeFailedToGetWords string = "failed-to-get-words"

	MessageCodeFailedToGetBooks     string = "failed-to-get-books"
	MessageCodeFailedToGetLessons   string = "failed-to-get-lessons"
	MessageCodeFailedToGetQuestions string = "failed-to-get-questions"
)
