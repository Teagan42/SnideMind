package extracttags

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sort"

	"github.com/teagan42/snidemind/models"
	"github.com/teagan42/snidemind/server/middleware"
	"go.uber.org/zap"
)

type Embedder struct {
	Logger       *zap.Logger
	Endpoint     string
	Model        string
	APIKey       *string
	APIKeyHeader *string
	Cache        map[string][]float64
	Client       *http.Client
}

func NewEmbedder(logger *zap.Logger, endpoint, model string, apiKey *string, apiKeyHeader *string) *Embedder {
	embedder := Embedder{
		Logger:       logger.Named("embedder"),
		Endpoint:     endpoint,
		Model:        model,
		APIKey:       apiKey,
		APIKeyHeader: apiKeyHeader,
		Cache:        make(map[string][]float64),
		Client:       &http.Client{},
	}
	embedder.Logger.Info("Generating tag embeddings", zap.String("model", model), zap.String("endpoint", endpoint))
	for _, tag := range TagTree {
		if vector, error := embedder.Embed(tag.Description); error == nil {
			embedder.Cache[tag.ID] = vector
		} else {
			embedder.Logger.Error("Failed to embed tag description", zap.String("tagID", tag.ID), zap.Error(error))
		}
	}

	return &embedder
}

func (e *Embedder) Embed(text string) ([]float64, error) {
	payload, _ := json.Marshal(map[string]interface{}{
		"input": text,
		"model": e.Model,
	})
	url, err := url.JoinPath(e.Endpoint, "embeddings")
	if err != nil {
		e.Logger.Error("Failed to join URL path", zap.Error(err))
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if body, err := middleware.PeekResponse(e.Logger, resp); err != nil {
		return nil, err
	} else {
		e.Logger.Info("Embedding response", zap.String("body", string(body)))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("embedding request failed")
	}

	var result models.EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Data) == 0 {
		return nil, errors.New("embedding data response is empty")
	}
	if len(result.Data[0].Embedding) == 0 {
		return nil, errors.New("embedding vector is empty")
	}
	return result.Data[0].Embedding, nil
}

func (e *Embedder) ExtractTagsWithWeights(userInput string) ([]ScoredTag, error) {
	e.Logger.Info("Extracting tags with weights", zap.String("input", userInput))
	inputVec, err := e.Embed(userInput)
	if err != nil {
		return nil, err
	}

	tagScores := map[string]ScoredTag{}
	tagVecs := e.Cache
	tagMap := make(map[string]TagNode)
	for _, tag := range TagTree {
		tagMap[tag.ID] = tag
	}

	// Embed all tag descriptions
	for id, vec := range tagVecs {
		score := CosineSimilarity(inputVec, vec)
		e.Logger.Info("Cosine similarity", zap.String("tagID", id), zap.Float64("score", score))
		if score >= 0.6 {
			tagScores[id] = ScoredTag{
				Tag:   tagMap[id],
				Score: score,
			}
		}
	}
	e.Logger.Info("Tag scores", zap.Any("tagScores", tagScores))
	// Bubble up parent scores
	for id, tagScore := range tagScores {
		if tag, err := getTagByID(id); err == nil {
			if tag.ParentID != nil {
				parentID := *tag.ParentID
				parentScore := tagScores[parentID]
				if tagScore.Score*0.8 > parentScore.Score {
					tagScores[parentID] = ScoredTag{
						Tag:   parentScore.Tag,
						Score: tagScore.Score * 0.8,
					}
				}
			}
		} else {
			return nil, err
		}
	}

	var result []ScoredTag
	for _, tagScore := range tagScores {
		result = append(result, tagScore)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})

	return result, nil
}
