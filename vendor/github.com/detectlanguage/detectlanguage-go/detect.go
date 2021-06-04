package detectlanguage

// DetectRequest contains language detection request params
type DetectRequest struct {
	Query string `json:"q"`
}

// DetectResponse is a resource containing language detection response
type DetectResponse struct {
	Data *DetectResponseData `json:"data"`
}

// DetectResponseData contains language detection response data
type DetectResponseData struct {
	Detections []*DetectionResult `json:"detections"`
}

// DetectionResult is single language detection result
type DetectionResult struct {
	Language   string  `json:"language"`
	Reliable   bool    `json:"isReliable"`
	Confidence float32 `json:"confidence"`
}

// DetectBatchRequest contains batch language detection request params
type DetectBatchRequest struct {
	Query []string `json:"q"`
}

// DetectBatchResponse is a resource batch containing language detection response
type DetectBatchResponse struct {
	Data *DetectBatchResponseData `json:"data"`
}

// DetectBatchResponseData contains batch language detection response data
type DetectBatchResponseData struct {
	Detections [][]*DetectionResult `json:"detections"`
}

// Detect executes language detection for a single text
func (c *Client) Detect(in string) (out []*DetectionResult, err error) {
	var response DetectResponse
	err = c.post(nil, "detect", &DetectRequest{Query: in}, &response)

	if err != nil {
		return nil, err
	}

	return response.Data.Detections, err
}

// DetectCode executes language detection for a single text and returns detected language code
func (c *Client) DetectCode(in string) (out string, err error) {
	detections, err := c.Detect(in)

	if err != nil {
		return "", err
	}

	if len(detections) == 0 {
		return "", &DetectionError{"Language not detected"}
	}

	return detections[0].Language, err
}

// DetectBatch executes language detection with multiple texts.
// It is significantly faster than doing a separate request for each text indivdually.
func (c *Client) DetectBatch(in []string) (out [][]*DetectionResult, err error) {
	var response DetectBatchResponse
	err = c.post(nil, "detect", &DetectBatchRequest{Query: in}, &response)

	if err != nil {
		return nil, err
	}

	return response.Data.Detections, err
}
