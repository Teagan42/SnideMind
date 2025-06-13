package extracttags

import (
	"bytes"
	"encoding/json"
	"errors"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"sort"

	"github.com/teagan42/snidemind/models"
	"go.uber.org/zap"
)

type Embedder struct {
	Logger       *zap.Logger
	Endpoint     string
	Model        string
	APIKey       *string
	APIKeyHeader *string
	Cache        map[string]TagNode
	Client       *http.Client
}

func NewEmbedder(logger *zap.Logger, endpoint, model string, apiKey *string, apiKeyHeader *string) *Embedder {
	embedder := Embedder{
		Logger:       logger.Named("embedder"),
		Endpoint:     endpoint,
		Model:        model,
		APIKey:       apiKey,
		APIKeyHeader: apiKeyHeader,
		Cache:        make(map[string]TagNode),
		Client:       &http.Client{},
	}
	embedder.Logger.Info("Generating tag embeddings", zap.String("model", model), zap.String("endpoint", endpoint))
	tags := make([]TagNode, 0, len(TagTree))
	batch := make([]string, 0, len(TagTree))
	for _, tag := range TagTree {
		tags = append(tags, tag)
		batch = append(batch, tag.Description)
	}

	if vectors, err := embedder.Embed(batch...); err != nil {
		embedder.Logger.Error("Failed to embed tag descriptions", zap.Error(err))
		return nil
	} else {
		for i, vector := range vectors {
			embedder.Cache[tags[i].ID] = TagNode{
				ID:          tags[i].ID,
				Name:        tags[i].Name,
				Description: tags[i].Description,
				Vector:      &vector,
				ParentID:    tags[i].ParentID,
			}
		}
	}

	return &embedder
}

func (e *Embedder) Embed(text ...string) ([][]float64, error) {
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

	vectors := make([][]float64, len(result.Data))
	for i, data := range result.Data {
		vectors[i] = data.Embedding
	}

	return vectors, nil
}

func (e *Embedder) ExtractTagsWithWeights(userInput string) ([]ScoredTag, error) {
	e.Logger.Info("Extracting tags with weights", zap.String("input", userInput))
	var inputVec []float64
	if vec, err := e.Embed(userInput); err != nil {
		return nil, err
	} else {
		inputVec = vec[0]
	}

	tagScores := map[string]ScoredTag{}
	// Embed all tag descriptions
	for id, vec := range e.Cache {
		score := CosineSimilarity(inputVec, *vec.Vector)
		e.Logger.Info("Cosine similarity", zap.String("tagID", id), zap.Float64("score", score))
		if score >= 0.6 {
			tagScores[id] = ScoredTag{
				Tag:   vec,
				Score: score,
			}
		}
	}

	e.Logger.Info("Tag scores", zap.Any("tagScores", tagScores))
	// Bubble up parent scores
	for _, tagScore := range tagScores {
		if tagScore.Tag.ParentID != nil {
			parentID := *tagScore.Tag.ParentID
			parentScore := tagScores[parentID]
			if tagScore.Score*0.8 > parentScore.Score {
				tagScores[parentID] = ScoredTag{
					Tag:   parentScore.Tag,
					Score: tagScore.Score * 0.8,
				}
			}
		}
	}
	e.Logger.Info("Updated tag scores with parent scores", zap.Any("tagScores", tagScores))
	values := slices.Collect(maps.Values(tagScores))
	sort.Slice(values, func(i, j int) bool {
		return values[i].Score > values[j].Score
	})

	return values, nil
}
