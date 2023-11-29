package llm

type GenerateTextRequest struct {
	Prompt *TextPrompt `json:"prompt,omitempty"`
	// Controls the randomness of the output.
	// Values can range from [0.0,1.0],
	Temperature *float32 `json:"temperature,omitempty"`
	// Number of generated responses to return.
	// This value must be between [1, 8], default 1
	CandidateCount *int32 `json:"candidate_count,omitempty"`
	// The maximum number of tokens to include in a candidate. Default 64
	MaxOutputTokens *int32 `json:"max_output_tokens,omitempty"`
	// The maximum cumulative probability of tokens to consider when sampling.
	TopP *float32 `json:"top_p,omitempty"`
	// The maximum number of tokens to consider when sampling.
	TopK *int32 `json:"top_k,omitempty"`
	// A list of unique `SafetySetting` instances for blocking unsafe content.
	SafetySettings []*SafetySetting `json:"safety_settings,omitempty"`
	// The set of character sequences (up to 5) that will stop output generation.
	// If specified, the API will stop at the first appearance of a stop
	// sequence. The stop sequence will not be included as part of the response.
	StopSequences []string `json:"stop_sequences,omitempty"`
}

// The category of a rating.
type HarmCategory int32

// Block at and beyond a specified harm probability.
type SafetySetting_HarmBlockThreshold int32

// Safety setting, affecting the safety-blocking behavior.
type SafetySetting struct {
	// Required. The category for this setting.
	Category HarmCategory `json:"category,omitempty"`
	// Required. Controls the probability threshold at which harm is blocked.
	Threshold SafetySetting_HarmBlockThreshold `json:"threshold,omitempty"`
}

type TextPrompt struct {
	// Required. The prompt text.
	Text string `json:"text,omitempty"`
}

// The response from the model, including candidate completions.
type GenerateTextResponse struct {
	// Candidate responses from the model.
	Candidates []*TextCompletion `json:"candidates,omitempty"`
	// A set of content filtering metadata for the prompt and response text.
	Filters []*ContentFilter `json:"filters,omitempty"`
	// Returns any safety feedback related to content filtering.
	SafetyFeedback []*SafetyFeedback `json:"safety_feedback,omitempty"`
}

// Safety feedback for an entire request.
type SafetyFeedback struct {
	// Safety rating evaluated from content.
	Rating *SafetyRating `json:"rating,omitempty"`
	// Safety settings applied to the request.
	Setting *SafetySetting `json:"setting,omitempty"`
}

// The probability that a piece of content is harmful.
type SafetyRating_HarmProbability int32

const (
	// Probability is unspecified.
	SafetyRating_HARM_PROBABILITY_UNSPECIFIED SafetyRating_HarmProbability = 0
	// Content has a negligible chance of being unsafe.
	SafetyRating_NEGLIGIBLE SafetyRating_HarmProbability = 1
	// Content has a low chance of being unsafe.
	SafetyRating_LOW SafetyRating_HarmProbability = 2
	// Content has a medium chance of being unsafe.
	SafetyRating_MEDIUM SafetyRating_HarmProbability = 3
	// Content has a high chance of being unsafe.
	SafetyRating_HIGH SafetyRating_HarmProbability = 4
)

// Safety rating for a piece of content.
type SafetyRating struct {
	// Required. The category for this rating.
	Category HarmCategory `json:"category,omitempty"`
	// Required. The probability of harm for this content.
	Probability SafetyRating_HarmProbability `json:"probability,omitempty"`
}

// A list of reasons why content may have been blocked.
type ContentFilter_BlockedReason int32

const (
	// A blocked reason was not specified.
	ContentFilter_BLOCKED_REASON_UNSPECIFIED ContentFilter_BlockedReason = 0
	// Content was blocked by safety settings.
	ContentFilter_SAFETY ContentFilter_BlockedReason = 1
	// Content was blocked, but the reason is uncategorized.
	ContentFilter_OTHER ContentFilter_BlockedReason = 2
)

type ContentFilter struct {
	// The reason content was blocked during request processing.
	Reason ContentFilter_BlockedReason `json:"reason,omitempty"`
	// A string that describes the filtering behavior in more detail.
	Message *string `json:"message,omitempty"`
}

// Output text returned from a model.
type TextCompletion struct {
	// Output only. The generated text returned from the model.
	Output string `json:"output,omitempty"`
	// Ratings for the safety of a response.
	SafetyRatings []*SafetyRating `json:"safety_ratings,omitempty"`
	// Output only. Citation information for model-generated `output` in this`TextCompletion`.
	CitationMetadata *CitationMetadata `json:"citation_metadata,omitempty"`
}

// A collection of source attributions for a piece of content.
type CitationMetadata struct {
	// Citations to sources for a specific response.
	CitationSources []*CitationSource `json:"citation_sources,omitempty"`
}

// A citation to a source for a portion of a specific response.
type CitationSource struct {
	// Optional. Start of segment of the response that is attributed to this source.
	StartIndex *int32 `json:"start_index,omitempty"`
	// Optional. End of the attributed segment, exclusive.
	EndIndex *int32 `json:"end_index,omitempty"`
	// Optional. URI that is attributed as a source for a portion of the text.
	Uri *string `json:"uri,omitempty"`
	// Optional. License for the GitHub project that is attributed as a source for segment.
	License *string `json:"license,omitempty"`
}

// Embed

// Request to get a text embedding from the model.
type EmbedTextRequest struct {
	// Required. The free-form input text that the model will turn into an embedding.
	Text string `json:"text,omitempty"`
}

// The response to a EmbedTextRequest.
type EmbedTextResponse struct {
	// Output only. The embedding generated from the input text.
	Embedding *Embedding `json:"embedding,omitempty"`
}

// A list of floats representing the embedding.
type Embedding struct {
	// The embedding values.
	Value []float32 `json:"value,omitempty"`
}

// Batch

// Request to get a text embedding from the model.
type BatchEmbedTextRequest struct {
	// Required. The free-form input text that the model will turn into an
	// embedding.
	Texts []string `json:"text,omitempty"`
}

// The response to a EmbedTextRequest.
type BatchEmbedTextResponse struct {
	// Output only. The embedding generated from the input text.
	Embeddings []Embedding `json:"embedding,omitempty"`
}
