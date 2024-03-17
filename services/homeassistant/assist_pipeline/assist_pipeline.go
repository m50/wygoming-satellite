package assist_pipeline

import (
	"encoding/json"

	"github.com/m50/wygoming-satellite/services/homeassistant"
)

type PipelineManager struct {
	ha                *homeassistant.HomeAssistant
	PreferredPipeline string
}

type Pipeline struct {
	ID                   string `json:"id"`
	ConversationEngine   string `json:"conversation_engine"`
	ConversationLanguage string `json:"conversation_language"`
	Language             string `json:"language"`
	Name                 string `json:"name"`
	STTEngine            string `json:"stt_engine"`
	SSTLanguage          string `json:"stt_language"`
	TTSEngine            string `json:"tts_engine"`
	TTSLanguage          string `json:"tts_language"`
	WakeWordEntity       string `json:"wake_word_entity"`
	WakeWordID           string `json:"wake_word_id"`
}

func NewPipelineManager(ha *homeassistant.HomeAssistant) PipelineManager {
	return PipelineManager{
		ha: ha,
		PreferredPipeline: "",
	}
}

func (pm *PipelineManager) ListPipelines() ([]Pipeline, error) {
	reqID := pm.ha.NextRequestId()
	if _, err := pm.ha.Request(reqID, map[string]interface{}{
		"id": reqID,
		"type": "assist_pipeline/pipeline/list",
	}); err != nil {
		return []Pipeline{}, err
	}
	resp, err := pm.ha.AwaitResponse(reqID)
	if err != nil {
		return []Pipeline{}, err
	}
	var r struct{
		Result struct{
			Pipelines []Pipeline `json:"pipelines"`
			PreferredPipeline string `json:"preferred_pipeline"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(resp), &r); err != nil {
		return []Pipeline{}, err
	}
	pm.PreferredPipeline = r.Result.PreferredPipeline

	return r.Result.Pipelines, nil
}

func (pm *PipelineManager) GetPipeline(id string) (Pipeline, error) {
	reqID := pm.ha.NextRequestId()
	if _, err := pm.ha.Request(reqID, map[string]interface{}{
		"id": reqID,
		"type": "assist_pipeline/pipeline/get",
		"pipeline_id": id,
	}); err != nil {
		return Pipeline{}, err
	}

	resp, err := pm.ha.AwaitResponse(reqID)
	if err != nil {
		return Pipeline{}, err
	}

	var r struct{
		Result Pipeline `json:"result"`
	}
	if err := json.Unmarshal([]byte(resp), &r); err != nil {
		return Pipeline{}, err
	}

	return r.Result, nil
}

