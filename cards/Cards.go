package cards

import "encoding/json"

// https://developers.google.com/workspace/add-ons/reference/rpc/google.apps.card.v1

type Icon struct {
	AltText   string `json:"altText,omitempty"`
	ImageType string `json:"imageType,omitempty"` // SQUARE or CIRCLE

	KnownIcon string `json:"knownIcon,omitempty"`
	IconUrl   string `json:"iconUrl,omitempty"`
	// AltText string `json:"altText,omitempty"`
}

type Color struct {
	Red   float64 `json:"red,omitempty"`
	Green float64 `json:"green,omitempty"`
	Blue  float64 `json:"blue,omitempty"`
	Alpha float64 `json:"alpha,omitempty"`
}

type OpenLink struct {
	OnClose string `json:"onClose,omitempty"`
	OpenAs  string `json:"openAs,omitempty"`
	Url     string `json:"url,omitempty"`
}

type Parameter struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type Action struct {
	Function      string       `json:"function,omitempty"`
	Interaction   string       `json:"interaction,omitempty"`
	LoadIndicator string       `json:"loadIndicator,omitempty"`
	Parameters    []*Parameter `json:"parameters,omitempty"`
	PersistValues bool         `json:"persistValues,omitempty"`
}

type CardAction struct {
	// ActionLabel: The label used to be displayed in the action menu item.
	ActionLabel string `json:"actionLabel,omitempty"`

	// OnClick: The onclick action for this action item.
	OnClick *OnClick `json:"onClick,omitempty"`
}

type CardHeader struct {
	// ImageStyle: The image's type (for example, square border or circular
	// border).
	//
	// Possible values:
	//   "IMAGE_STYLE_UNSPECIFIED"
	//   "IMAGE" - Square border.
	//   "AVATAR" - Circular border.
	ImageStyle string `json:"imageStyle,omitempty"`

	// ImageUrl: The URL of the image in the card header.
	ImageUrl string `json:"imageUrl,omitempty"`

	// Subtitle: The subtitle of the card header.
	Subtitle string `json:"subtitle,omitempty"`

	// Title: The title must be specified. The header has a fixed height: if
	// both a title and subtitle is specified, each takes up one line. If
	// only the title is specified, it takes up both lines.
	Title string `json:"title,omitempty"`
}

type Section struct {
	// Header: The header of the section. Formatted text is supported. For
	// more information about formatting text, see Formatting text in Google
	// Chat apps
	// (https://developers.google.com/chat/format-messages#card-formatting)
	// and Formatting text in Google Workspace Add-ons
	// (https://developers.google.com/apps-script/add-ons/concepts/widgets#text_formatting).
	Header string `json:"header,omitempty"`

	// Widgets: A section must contain at least one widget.
	Widgets []*WidgetMarkup `json:"widgets,omitempty"`
}

type Image struct {
	// AspectRatio: The aspect ratio of this image (width and height). This
	// field lets you reserve the right height for the image while waiting
	// for it to load. It's not meant to override the built-in aspect ratio
	// of the image. If unset, the server fills it by prefetching the image.
	AspectRatio float64 `json:"aspectRatio,omitempty"`

	// ImageUrl: The URL of the image.
	ImageUrl string `json:"imageUrl,omitempty"`

	// OnClick: The `onclick` action.
	OnClick *OnClick `json:"onClick,omitempty"`
}

type KeyValue struct {
	// BottomLabel: The text of the bottom label. Formatted text supported.
	// For more information about formatting text, see Formatting text in
	// Google Chat apps
	// (https://developers.google.com/chat/format-messages#card-formatting)
	// and Formatting text in Google Workspace Add-ons
	// (https://developers.google.com/apps-script/add-ons/concepts/widgets#text_formatting).
	BottomLabel string `json:"bottomLabel,omitempty"`

	// Button: A button that can be clicked to trigger an action.
	Button *Button `json:"button,omitempty"`

	// Content: The text of the content. Formatted text supported and always
	// required. For more information about formatting text, see Formatting
	// text in Google Chat apps
	// (https://developers.google.com/chat/format-messages#card-formatting)
	// and Formatting text in Google Workspace Add-ons
	// (https://developers.google.com/apps-script/add-ons/concepts/widgets#text_formatting).
	Content string `json:"content,omitempty"`

	// ContentMultiline: If the content should be multiline.
	ContentMultiline bool `json:"contentMultiline,omitempty"`

	// Icon: An enum value that's replaced by the Chat API with the
	// corresponding icon image.
	//
	// Possible values:
	//   "ICON_UNSPECIFIED"
	//   "AIRPLANE"
	//   "BOOKMARK"
	//   "BUS"
	//   "CAR"
	//   "CLOCK"
	//   "CONFIRMATION_NUMBER_ICON"
	//   "DOLLAR"
	//   "DESCRIPTION"
	//   "EMAIL"
	//   "EVENT_PERFORMER"
	//   "EVENT_SEAT"
	//   "FLIGHT_ARRIVAL"
	//   "FLIGHT_DEPARTURE"
	//   "HOTEL"
	//   "HOTEL_ROOM_TYPE"
	//   "INVITE"
	//   "MAP_PIN"
	//   "MEMBERSHIP"
	//   "MULTIPLE_PEOPLE"
	//   "OFFER"
	//   "PERSON"
	//   "PHONE"
	//   "RESTAURANT_ICON"
	//   "SHOPPING_CART"
	//   "STAR"
	//   "STORE"
	//   "TICKET"
	//   "TRAIN"
	//   "VIDEO_CAMERA"
	//   "VIDEO_PLAY"
	Icon string `json:"icon,omitempty"`

	// IconUrl: The icon specified by a URL.
	IconUrl string `json:"iconUrl,omitempty"`

	// OnClick: The `onclick` action. Only the top label, bottom label, and
	// content region are clickable.
	OnClick *OnClick `json:"onClick,omitempty"`

	// TopLabel: The text of the top label. Formatted text supported. For
	// more information about formatting text, see Formatting text in Google
	// Chat apps
	// (https://developers.google.com/chat/format-messages#card-formatting)
	// and Formatting text in Google Workspace Add-ons
	// (https://developers.google.com/apps-script/add-ons/concepts/widgets#text_formatting).
	TopLabel string `json:"topLabel,omitempty"`
}

type TextParagraph struct {
	Text string `json:"text,omitempty"`
}

type ButtonList struct {
	// Buttons: An array of buttons.
	Buttons []*Button `json:"buttons,omitempty"`
}

type WidgetMarkup struct {
	HorizontalAlignment string `json:"horizontalAlignment,omitempty"`

	// The following are all a oneof:
	ButtonList *ButtonList `json:"buttonList,omitempty"`

	/*
		SelectionInput *GoogleAppsCardV1SelectionInput `json:"selectionInput,omitempty"`
		TextInput *GoogleAppsCardV1TextInput `json:"textInput,omitempty"`
	*/

	// Image: Display an image in this widget.
	Image *Image `json:"image,omitempty"`

	// TextParagraph: Display a text paragraph in this widget.
	TextParagraph *TextParagraph `json:"textParagraph,omitempty"`
}

type CardFixedFooter struct {
	PrimaryButton   *Button `json:"primaryButton,omitempty"`
	SecondaryButton *Button `json:"secondaryButton,omitempty"`
}

type Card struct {
	CardActions []*CardAction `json:"cardActions,omitempty"`
	Header      *CardHeader   `json:"header,omitempty"`

	// Name: Name of the card.
	Name string `json:"name,omitempty"`

	// Sections: Sections are separated by a line divider.
	Sections []*Section `json:"sections,omitempty"`

	FixedFooter *CardFixedFooter `json:"fixedFooter,omitempty"`
}

type OnClick struct {
	// Action: If specified, an action is triggered by this `onClick`.
	Action *Action `json:"action,omitempty"`

	// Card: A new card is pushed to the card stack after clicking if
	// specified. Supported by Google Workspace Add-ons, but not Google Chat
	// apps.
	Card *Card `json:"card,omitempty"`

	// OpenDynamicLinkAction: An add-on triggers this action when the action
	// needs to open a link. This differs from the `open_link` above in that
	// this needs to talk to server to get the link. Thus some preparation
	// work is required for web client to do before the open link action
	// response comes back. Supported by Google Workspace Add-ons, but not
	// Google Chat apps.
	OpenDynamicLinkAction *Action `json:"openDynamicLinkAction,omitempty"`

	// OpenLink: If specified, this `onClick` triggers an open link action.
	OpenLink *OpenLink `json:"openLink,omitempty"`
}

type Button struct {
	Text     string   `json:"text,omitempty"`
	Icon     *Icon    `json:"icon,omitempty"`
	Color    *Color   `json:"color,omitempty"`
	OnClick  *OnClick `json:"onClick,omitempty"`
	Disabled bool     `json:"disabled,omitempty"`
	AltText  string   `json:"altText,omitempty"`
}

type HostAppAction struct {
	EditorAction *EditorClientAction `json:"editorAction,omitempty"`
}

type EditorClientAction struct {
	RequestFileScopeForActiveDocument RequestFileScopeForActiveDocument `json:"requestFileScopeForActiveDocument,omitempty"`
}

type RequestFileScopeForActiveDocument struct{}

type SubmitFormResponse struct {
	RenderAction *RenderActions `json:"renderActions,omitempty"`
	StateChanged *bool          `json:"stateChanged,omitempty"`
}

func (s *SubmitFormResponse) MarshalJSON() ([]byte, error) {
	type NoMethod SubmitFormResponse
	raw := NoMethod(*s)
	return json.Marshal(raw)
}

type RenderActions struct {
	HostAppAction *HostAppAction       `json:"hostAppAction,omitempty"`
	Action        *RenderActionsAction `json:"action,omitempty"`
	Schema        string               `json:"schema,omitempty"` // unused
}

func (s *RenderActions) MarshalJSON() ([]byte, error) {
	type NoMethod RenderActions
	raw := NoMethod(*s)
	return json.Marshal(raw)
}

// TODO implement me
type RenderActionsAction struct {
	Navigations  *[]Navigation `json:"navigations,omitempty"`
	Link         *OpenLink     `json:"link,omitempty"`
	Notification *Notification `json:"notification,omitempty"`

	//LinkPreview *LinkPreview `json:"linkPreview,omitempty"`
	/*
	   navigations[]
	   Navigation

	   Push, pop, or update displayed cards.

	   link
	   OpenLink

	   Immediately open the target link in a new tab or a pop-up.

	   notification
	   Notification

	   Display a notification to the end-user.

	   linkPreview
	   LinkPreview

	   Display a link preview to the end user.*/
}

type Notification struct {
	Text string `json:"text,omitempty"`
}

// Card action that manipulates the card stack.
// Can only have one of the following
type Navigation struct {
	PopToRoot  *bool   `json:"popToRoot,omitempty"`
	Pop        *bool   `json:"pop,omitempty"`
	PopToCard  *string `json:"popToCard,omitempty"`
	PushCard   *Card   `json:"pushCard,omitempty"`
	UpdateCard *Card   `json:"updateCard,omitempty"`
}

type RenderCard struct {
	Card *Card `json:"card,omitempty"`
}
