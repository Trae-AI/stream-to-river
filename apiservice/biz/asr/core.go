// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package asr

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

// ProtocolVersion represents the protocol version used in WebSocket communication, stored as a byte.
type ProtocolVersion byte

// MessageType represents the type of message in WebSocket communication, stored as a byte.
type MessageType byte

// MessageTypeSpecificFlags represents specific flags for a message type, stored as a byte.
type MessageTypeSpecificFlags byte

// SerializationType represents the serialization method of a message, stored as a byte.
type SerializationType byte

// CompressionType represents the compression method of a message, stored as a byte.
type CompressionType byte

// Constants define various codes, protocol versions, and sizes used in the application.
const (
	SuccessCode = 1000

	PROTOCOL_VERSION    = ProtocolVersion(0b0001)
	DEFAULT_HEADER_SIZE = 0b0001

	PROTOCOL_VERSION_BITS            = 4
	HEADER_BITS                      = 4
	MESSAGE_TYPE_BITS                = 4
	MESSAGE_TYPE_SPECIFIC_FLAGS_BITS = 4
	MESSAGE_SERIALIZATION_BITS       = 4
	MESSAGE_COMPRESSION_BITS         = 4
	RESERVED_BITS                    = 8
)

// Constants define different message types used in WebSocket communication.
const (
	CLIENT_FULL_REQUEST       = MessageType(0b0001)
	CLIENT_AUDIO_ONLY_REQUEST = MessageType(0b0010)
	SERVER_FULL_RESPONSE      = MessageType(0b1001)
	SERVER_ACK                = MessageType(0b1011)
	SERVER_ERROR_RESPONSE     = MessageType(0b1111)
)

// Constants define specific flags for message types.
const (
	NO_SEQUENCE    = MessageTypeSpecificFlags(0b0000)
	POS_SEQUENCE   = MessageTypeSpecificFlags(0b0001)
	NEG_SEQUENCE   = MessageTypeSpecificFlags(0b0010)
	NEG_SEQUENCE_1 = MessageTypeSpecificFlags(0b0011)
)

// Constants define different message serialization methods.
const (
	NO_SERIALIZATION = SerializationType(0b0000)
	JSON             = SerializationType(0b0001)
	THRIFT           = SerializationType(0b0011)
	CUSTOM_TYPE      = SerializationType(0b1111)
)

// Constants define different message compression methods.
const (
	NO_COMPRESSION     = CompressionType(0b0000)
	GZIP               = CompressionType(0b0001)
	CUSTOM_COMPRESSION = CompressionType(0b1111)
)

// Default WebSocket headers for different types of client requests.
var DefaultFullClientWsHeader = []byte{0x11, 0x10, 0x11, 0x00}
var DefaultAudioOnlyWsHeader = []byte{0x11, 0x20, 0x11, 0x00}
var DefaultLastAudioWsHeader = []byte{0x11, 0x22, 0x11, 0x00}

// gzipCompress compresses the input byte slice using gzip.
//
// Parameters:
//   - input: The byte slice to be compressed.
//
// Returns:
//   - []byte: The compressed byte slice.
func gzipCompress(input []byte) []byte {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(input)
	w.Close()
	return b.Bytes()
}

// gzipDecompress decompresses the input byte slice using gzip.
//
// Parameters:
//   - input: The byte slice to be decompressed.
//
// Returns:
//   - []byte: The decompressed byte slice.
func gzipDecompress(input []byte) []byte {
	b := bytes.NewBuffer(input)
	r, _ := gzip.NewReader(b)
	out, _ := ioutil.ReadAll(r)
	r.Close()
	return out
}

// AsrResponse represents the response structure from the ASR service.
type AsrResponse struct {
	Reqid    string                 `json:"reqid"`
	Code     int                    `json:"code"`
	Message  string                 `json:"message"`
	Sequence int                    `json:"sequence"`
	Results  []Result               `json:"result,omitempty"`
	Addition map[string]interface{} `json:"addition,omitempty"`
}

// Result represents the recognition result in the ASR response.
type Result struct {
	Text       string      `json:"text"`
	Confidence int         `json:"confidence"`
	Language   string      `json:"language,omitempty"`
	Utterances []Utterance `json:"utterances,omitempty"`
}

// Utterance represents an utterance in the recognition result.
type Utterance struct {
	Text      string `json:"text"`
	StartTime int    `json:"start_time"`
	EndTime   int    `json:"end_time"`
	Definite  bool   `json:"definite"`
	Words     []Word `json:"words"`
	Language  string `json:"language"`
}

// Word represents a word in the utterance.
type Word struct {
	Text          string `json:"text"`
	StartTime     int    `json:"start_time"`
	EndTime       int    `json:"end_time"`
	Pronounce     string `json:"pronounce"`
	BlankDuration int    `json:"blank_duration"`
}

// WsHeader represents the WebSocket header structure.
type WsHeader struct {
	ProtocolVersion          ProtocolVersion
	DefaultHeaderSize        int
	MessageType              MessageType
	MessageTypeSpecificFlags MessageTypeSpecificFlags
	SerializationType        SerializationType
	CompressionType          CompressionType
}

// RequestAsr is an interface for requesting ASR service.
type RequestAsr interface {
	requestAsr(audio_data []byte)
}

// AsrClient represents the client for the ASR service.
type AsrClient struct {
	Appid    string
	Token    string
	Cluster  string
	Workflow string
	Format   string
	Codec    string
	SegSize  int
	Url      string
}

// buildAsrClient creates a new instance of AsrClient with default settings.
//
// Returns:
//   - AsrClient: A new AsrClient instance.
func buildAsrClient() AsrClient {
	client := AsrClient{}
	client.Workflow = "audio_in,resample,partition,vad,fe,decode"
	client.SegSize = 160000
	client.Format = "wav"
	client.Codec = "raw"
	return client
}

// requestAsr sends an audio recognition request to the ASR service via WebSocket.
// It first sends a full client request, then sends segmented audio data, and parses the responses.
//
// Parameters:
//   - audioData: The audio data to be recognized.
//
// Returns:
//   - AsrResponse: The response from the ASR service.
//   - error: An error object if an unexpected error occurs during the process.
func (client *AsrClient) requestAsr(audioData []byte) (AsrResponse, error) {
	var tokenHeader = http.Header{"Authorization": []string{fmt.Sprintf("Bearer;%s", client.Token)}}
	c, _, err := websocket.DefaultDialer.Dial("wss://openspeech.bytedance.com/api/v2/asr", tokenHeader)
	if err != nil {
		fmt.Println(err)
		return AsrResponse{}, err
	}
	defer c.Close()

	// Send full client request
	req := client.constructRequest()
	payload := gzipCompress(req)
	payloadSize := len(payload)
	payloadSizeArr := make([]byte, 4)
	binary.BigEndian.PutUint32(payloadSizeArr, uint32(payloadSize))

	fullClientMsg := make([]byte, len(DefaultFullClientWsHeader))
	copy(fullClientMsg, DefaultFullClientWsHeader)
	fullClientMsg = append(fullClientMsg, payloadSizeArr...)
	fullClientMsg = append(fullClientMsg, payload...)
	c.WriteMessage(websocket.BinaryMessage, fullClientMsg)
	_, msg, err := c.ReadMessage()
	if err != nil {
		fmt.Println("fail to read message fail, err:", err.Error())
		return AsrResponse{}, err
	}
	asrResponse, err := client.parseResponse(msg)
	if err != nil {
		fmt.Println("fail to parse response ", err.Error())
		return AsrResponse{}, err
	}

	// Send segment audio request
	for sentSize := 0; sentSize < len(audioData); sentSize += client.SegSize {
		lastAudio := false
		if sentSize+client.SegSize >= len(audioData) {
			lastAudio = true
		}
		dataSlice := make([]byte, 0)
		audioMsg := make([]byte, len(DefaultAudioOnlyWsHeader))
		if !lastAudio {
			dataSlice = audioData[sentSize : sentSize+client.SegSize]
			copy(audioMsg, DefaultAudioOnlyWsHeader)
		} else {
			dataSlice = audioData[sentSize:]
			copy(audioMsg, DefaultLastAudioWsHeader)
		}
		payload = gzipCompress(dataSlice)
		payloadSize := len(payload)
		payloadSizeArr := make([]byte, 4)
		binary.BigEndian.PutUint32(payloadSizeArr, uint32(payloadSize))
		audioMsg = append(audioMsg, payloadSizeArr...)
		audioMsg = append(audioMsg, payload...)
		c.WriteMessage(websocket.BinaryMessage, audioMsg)
		_, msg, err := c.ReadMessage()
		if err != nil {
			fmt.Println("fail to read message fail, err:", err.Error())
			return AsrResponse{}, err
		}
		asrResponse, err = client.parseResponse(msg)
		if err != nil {
			fmt.Println("fail to parse response ", err.Error())
			return AsrResponse{}, err
		}
	}
	return asrResponse, nil
}

// constructRequest constructs the initial request for the ASR service.
// It generates a unique request ID and populates the request with client information.
//
// Returns:
//   - []byte: The JSON - marshaled request data.
func (client *AsrClient) constructRequest() []byte {
	reqid := uuid.New().String()
	req := make(map[string]map[string]interface{})
	req["app"] = make(map[string]interface{})
	req["app"]["appid"] = client.Appid
	req["app"]["cluster"] = client.Cluster
	req["app"]["token"] = client.Token
	req["user"] = make(map[string]interface{})
	req["user"]["uid"] = "uid"
	req["request"] = make(map[string]interface{})
	req["request"]["reqid"] = reqid
	req["request"]["nbest"] = 1
	req["request"]["workflow"] = client.Workflow
	req["request"]["result_type"] = "full"
	req["request"]["sequence"] = 1
	req["audio"] = make(map[string]interface{})
	req["audio"]["format"] = client.Format
	req["audio"]["codec"] = client.Codec
	reqStr, _ := json.Marshal(req)
	return reqStr
}

// parseResponse parses the response received from the ASR service.
// It extracts header information, decompresses the payload if necessary,
// and unmarshals the JSON data into an AsrResponse struct.
//
// Parameters:
//   - msg: The raw response message received from the WebSocket.
//
// Returns:
//   - AsrResponse: The parsed response.
//   - error: An error object if an unexpected error occurs during parsing.
func (client *AsrClient) parseResponse(msg []byte) (AsrResponse, error) {
	headerSize := msg[0] & 0x0f
	messageType := msg[1] >> 4
	serializationMethod := msg[2] >> 4
	messageCompression := msg[2] & 0x0f
	payload := msg[headerSize*4:]
	payloadMsg := make([]byte, 0)
	payloadSize := 0

	if messageType == byte(SERVER_FULL_RESPONSE) {
		payloadSize = int(int32(binary.BigEndian.Uint32(payload[0:4])))
		payloadMsg = payload[4:]
	} else if messageType == byte(SERVER_ACK) {
		seq := int32(binary.BigEndian.Uint32(payload[:4]))
		if len(payload) >= 8 {
			payloadSize = int(binary.BigEndian.Uint32(payload[4:8]))
			payloadMsg = payload[8:]
		}
		fmt.Println("SERVER_ACK seq: ", seq)
	} else if messageType == byte(SERVER_ERROR_RESPONSE) {
		code := int32(binary.BigEndian.Uint32(payload[:4]))
		payloadSize = int(binary.BigEndian.Uint32(payload[4:8]))
		payloadMsg = payload[8:]
		fmt.Println("SERVER_ERROR_RESPONE code: ", code)
		return AsrResponse{}, errors.New(string(payloadMsg))
	}
	if payloadSize == 0 {
		return AsrResponse{}, errors.New("payload size if 0")
	}
	if messageCompression == byte(GZIP) {
		payloadMsg = gzipDecompress(payloadMsg)
	}

	var asrResponse = AsrResponse{}
	if serializationMethod == byte(JSON) {
		err := json.Unmarshal(payloadMsg, &asrResponse)
		if err != nil {
			fmt.Println("fail to unmarshal response, ", err.Error())
			return AsrResponse{}, err
		}
	}
	return asrResponse, nil
}

// AsrService represents the ASR service.
type AsrService struct {
	appid   string
	token   string
	cluster string
	client  AsrClient
}

// NewAsrService creates a new instance of AsrService with default settings.
//
// Returns:
//   - *AsrService: A pointer to the new AsrService instance.
func NewAsrService() *AsrService {
	return &AsrService{
		appid:   appID,
		token:   token,
		cluster: cluster,
		client:  buildAsrClient(),
	}
}

// RecognizeAudio initiates an audio recognition request using the ASR service.
//
// Parameters:
//   - context: The context for the request.
//   - audioData: The audio data to be recognized.
//   - format: The format of the audio data.
//
// Returns:
//   - *AsrResponse: A pointer to the response from the ASR service.
//   - error: An error object if an unexpected error occurs during the process.
func (s *AsrService) RecognizeAudio(context context.Context, audioData []byte, format string) (*AsrResponse, error) {
	s.client.Appid = s.appid
	s.client.Token = s.token
	s.client.Cluster = s.cluster
	s.client.Format = format

	asrResponse, err := s.client.requestAsr(audioData)
	if err != nil {
		return nil, fmt.Errorf("fail to request asr: %w", err)
	}

	return &asrResponse, nil
}
