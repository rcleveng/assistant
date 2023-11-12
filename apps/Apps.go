package apps

type CommonEventObject struct {
	HostApp  string `json:"hostApp,omitempty"`
	Platform string `json:"platform,omitempty"`
}

type AuthorizationEventObject struct {
	UserOAuthToken string `json:"userOAuthToken,omitempty"`
	SystemIdToken  string `json:"systemIdToken,omitempty"`
	UserIdToken    string `json:"userIdToken,omitempty"`
}

type Docs struct {
}

type Event struct {
	CommonEventObject        *CommonEventObject        `json:"commonEventObject,omitempty"`
	AuthorizationEventObject *AuthorizationEventObject `json:"authorizationEventObject,omitempty"`
	Docs                     *Docs                     `json:"docs,omitempty"`
}
