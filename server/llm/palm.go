package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	generativelanguage "cloud.google.com/go/ai/generativelanguage/apiv1beta2"
	pb "cloud.google.com/go/ai/generativelanguage/apiv1beta2/generativelanguagepb"
	"github.com/rcleveng/assistant/server/env"
	"google.golang.org/api/option"
)

type GenerateTextRequest struct {
	Prompt *TextPrompt `json:"prompt,omitempty"`
	// Controls the randomness of the output.
	// Note: The default value varies by model, see the `Model.temperature`
	// attribute of the `Model` returned the `getModel` function.
	//
	// Values can range from [0.0,1.0],
	// inclusive. A value closer to 1.0 will produce responses that are more
	// varied and creative, while a value closer to 0.0 will typically result in
	// more straightforward responses from the model.
	Temperature *float32 `json:"temperature,omitempty"`
	// Number of generated responses to return.
	//
	// This value must be between [1, 8], inclusive. If unset, this will default
	// to 1.
	CandidateCount *int32 `json:"candidate_count,omitempty"`
	// The maximum number of tokens to include in a candidate.
	//
	// If unset, this will default to 64.
	MaxOutputTokens *int32 `json:"max_output_tokens,omitempty"`
	// The maximum cumulative probability of tokens to consider when sampling.
	//
	// The model uses combined Top-k and nucleus sampling.
	//
	// Tokens are sorted based on their assigned probabilities so that only the
	// most liekly tokens are considered. Top-k sampling directly limits the
	// maximum number of tokens to consider, while Nucleus sampling limits number
	// of tokens based on the cumulative probability.
	//
	// Note: The default value varies by model, see the `Model.top_p`
	// attribute of the `Model` returned the `getModel` function.
	TopP *float32 `json:"top_p,omitempty"`
	// The maximum number of tokens to consider when sampling.
	//
	// The model uses combined Top-k and nucleus sampling.
	//
	// Top-k sampling considers the set of `top_k` most probable tokens.
	// Defaults to 40.
	//
	// Note: The default value varies by model, see the `Model.top_k`
	// attribute of the `Model` returned the `getModel` function.
	TopK *int32 `json:"top_k,omitempty"`
	// A list of unique `SafetySetting` instances for blocking unsafe content.
	//
	// that will be enforced on the `GenerateTextRequest.prompt` and
	// `GenerateTextResponse.candidates`. There should not be more than one
	// setting for each `SafetyCategory` type. The API will block any prompts and
	// responses that fail to meet the thresholds set by these settings. This list
	// overrides the default settings for each `SafetyCategory` specified in the
	// safety_settings. If there is no `SafetySetting` for a given
	// `SafetyCategory` provided in the list, the API will use the default safety
	// setting for that category.
	SafetySettings []*SafetySetting `json:"safety_settings,omitempty"`
	// The set of character sequences (up to 5) that will stop output generation.
	// If specified, the API will stop at the first appearance of a stop
	// sequence. The stop sequence will not be included as part of the response.
	StopSequences []string `json:"stop_sequences,omitempty"`
}

// The category of a rating.
//
// These categories cover various kinds of harms that developers
// may wish to adjust.
type HarmCategory int32

// Block at and beyond a specified harm probability.
type SafetySetting_HarmBlockThreshold int32

// Safety setting, affecting the safety-blocking behavior.
//
// Passing a safety setting for a category changes the allowed proability that
// content is blocked.
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
	// A set of content filtering metadata for the prompt and response
	// text.
	//
	// This indicates which `SafetyCategory`(s) blocked a
	// candidate from this response, the lowest `HarmProbability`
	// that triggered a block, and the HarmThreshold setting for that category.
	// This indicates the smallest change to the `SafetySettings` that would be
	// necessary to unblock at least 1 response.
	//
	// The blocking is configured by the `SafetySettings` in the request (or the
	// default `SafetySettings` of the API).
	Filters []*ContentFilter `json:"filters,omitempty"`
	// Returns any safety feedback related to content filtering.
	SafetyFeedback []*SafetyFeedback `json:"safety_feedback,omitempty"`
}

// Safety feedback for an entire request.
//
// This field is populated if content in the input and/or response is blocked
// due to safety settings. SafetyFeedback may not exist for every HarmCategory.
// Each SafetyFeedback will return the safety settings used by the request as
// well as the lowest HarmProbability that should be allowed in order to return
// a result.
type SafetyFeedback struct {
	// Safety rating evaluated from content.
	Rating *SafetyRating `json:"rating,omitempty"`
	// Safety settings applied to the request.
	Setting *SafetySetting `json:"setting,omitempty"`
}

// The probability that a piece of content is harmful.
//
// The classification system gives the probability of the content being
// unsafe. This does not indicate the severity of harm for a piece of content.
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
//
// The safety rating contains the category of harm and the
// harm probability level in that category for a piece of content.
// Content is classified for safety across a number of
// harm categories and the probability of the harm classification is included
// here.
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
	//
	// There is at most one rating per category.
	SafetyRatings []*SafetyRating `json:"safety_ratings,omitempty"`
	// Output only. Citation information for model-generated `output` in this
	// `TextCompletion`.
	//
	// This field may be populated with attribution information for any text
	// included in the `output`.
	CitationMetadata *CitationMetadata `json:"citation_metadata,omitempty"`
}

// A collection of source attributions for a piece of content.
type CitationMetadata struct {
	// Citations to sources for a specific response.
	CitationSources []*CitationSource `json:"citation_sources,omitempty"`
}

// A citation to a source for a portion of a specific response.
type CitationSource struct {
	// Optional. Start of segment of the response that is attributed to this
	// source.
	//
	// Index indicates the start of the segment, measured in bytes.
	StartIndex *int32 `json:"start_index,omitempty"`
	// Optional. End of the attributed segment, exclusive.
	EndIndex *int32 `json:"end_index,omitempty"`
	// Optional. URI that is attributed as a source for a portion of the text.
	Uri *string `json:"uri,omitempty"`
	// Optional. License for the GitHub project that is attributed as a source for
	// segment.
	//
	// License info is required for code citations.
	License *string `json:"license,omitempty"`
}

// Embed

// Request to get a text embedding from the model.
type EmbedTextRequest struct {
	// Required. The free-form input text that the model will turn into an
	// embedding.
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

type PalmLLMClient struct {
	c           *generativelanguage.TextClient
	environment *env.ServerEnvironment
	endpoint    string
}

func (c *PalmLLMClient) Close() error {
	return c.c.Close()
}

func (c *PalmLLMClient) Post(model, fn string, body io.Reader) ([]byte, error) {
	url := fmt.Sprintf("%s/%s:%s", c.endpoint, model, fn)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-goog-api-key", c.environment.PalmApiKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	response, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return response, nil
}

//model = models/text-bison-001

func (c *PalmLLMClient) GenerateText(ctx context.Context, prompt string) (string, error) {

	req := &GenerateTextRequest{
		Prompt: &TextPrompt{
			Text: prompt,
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	b, err := c.Post("models/text-bison-001", "generateText", bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	resp := GenerateTextResponse{}
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) > 0 {
		s := resp.Candidates[0].Output
		//fmt.Println("LLM Response: " + s)
		return s, nil
	}
	return "", fmt.Errorf("no candidate response, just %#v", resp)
}

func (c *PalmLLMClient) EmbedText(ctx context.Context, text string) ([]float32, error) {
	req := &pb.EmbedTextRequest{
		Model: "models/embedding-gecko-001",
		Text:  text,
	}

	resp, err := c.c.EmbedText(ctx, req)
	if err != nil {
		return nil, err
	}

	emb := resp.GetEmbedding().GetValue()
	return emb, nil
}

// TODO - use the batchEmbedText endpoint that's part of v1beta3 once available in the
// client libraries or just give up on the client libraries and call the rest apis
// manually.
func (c *PalmLLMClient) BatchEmbedText(ctx context.Context, text []string) ([][]float32, error) {
	emb := make([][]float32, len(text))

	for _, t := range text {
		ce, err := c.EmbedText(ctx, t)
		if err == nil {
			emb = append(emb, ce)
		} else {
			emb = append(emb, []float32{})
		}
	}

	return emb, nil
}

func NewPalmLLMClient(ctx context.Context, environment *env.ServerEnvironment) (*PalmLLMClient, error) {
	apiKey := option.WithAPIKey(environment.PalmApiKey)
	c, err := generativelanguage.NewTextRESTClient(ctx, apiKey)
	if err != nil {
		return nil, err
	}
	return &PalmLLMClient{
		c:           c,
		environment: environment,
		endpoint:    "https://generativelanguage.googleapis.com/v1beta3",
	}, nil
}
