package apps

//https://developers.google.com/workspace/add-ons/reference/rpc/google.cloud.gsuiteaddons.v1

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
	Id    *string `json:"id,omitempty"`
	Title *string `json:"title,omitempty"`
}

type WorkspaceAppEvent struct {
	CommonEventObject        *CommonEventObject        `json:"commonEventObject,omitempty"`
	AuthorizationEventObject *AuthorizationEventObject `json:"authorizationEventObject,omitempty"`
	Docs                     *Docs                     `json:"docs,omitempty"`
}
