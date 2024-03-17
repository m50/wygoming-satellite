package homeassistant

// Receive:
// {"id":46,"type":"result","success":true,"result":null}
// [{"id":45,"type":"event","event":{"type":"run-start","data":{"pipeline":"01hp2qn3txcy889v4ceffyw0nj","language":"en","runner_data":{"stt_binary_handler_id":null,"timeout":300}},"timestamp":"2024-03-17T18:09:50.326565+00:00"}},{"id":45,"type":"event","event":{"type":"intent-start","data":{"engine":"da161b0bfd19f6df805dd895237432d8","language":"en","intent_input":"What is the weather today?","conversation_id":null,"device_id":null},"timestamp":"2024-03-17T18:09:50.326626+00:00"}}]
// [{"id":45,"type":"event","event":{"type":"intent-end","data":{"intent_output":{"response":{"speech":{"plain":{"speech":"The weather today is rainy with a high temperature of 11.7°C and a low temperature of 11.0°C. Would you like more details about the weather forecast?","extra_data":null,"original_speech":"The weather today is rainy with a high temperature of 11.7°C and a low temperature of 11.0°C. Would you like more details about the weather forecast?","agent_name":"ChatGPT","agent_id":"3b5d5e8ec9405e9bcbdd9b601ab8254c"}},"card":{},"language":"en","response_type":"action_done","data":{"targets":[],"success":[],"failed":[]}},"conversation_id":"01HS6SP2XD1E5HEHVDZE8R2WMR"}},"timestamp":"2024-03-17T18:09:51.844638+00:00"}},{"id":45,"type":"event","event":{"type":"run-end","data":null,"timestamp":"2024-03-17T18:09:51.844691+00:00"}}]

// Send:
// {"start_stage":"intent","input":{"text":"What is the weather today?"},"end_stage":"intent","pipeline":"01hp2qn3txcy889v4ceffyw0nj","conversation_id":null,"type":"assist_pipeline/run","id":45}
// {"type":"assist_pipeline/pipeline/get","pipeline_id":"01hp2qn3txcy889v4ceffyw0nj","id":44}
// {"type":"assist_pipeline/pipeline/list","id":45}

