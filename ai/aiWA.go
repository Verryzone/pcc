package ai

import (
	"context"
	"fmt"
	"sync"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var client openai.Client

var aiReady bool

var userHistories = make(map[string][]openai.ChatCompletionMessageParamUnion)

var mu sync.Mutex

func InitAi() {
	client = openai.NewClient(
		option.WithBaseURL("http://localhost:8080/v1/"),
		option.WithAPIKey("saksake-karena-gak-butuh-api-key"),
	)
	aiReady = true

	fmt.Println("AI Engine berhasil diinisialisasi.")
}

func TanyaAi(userID string, userInput string) string {
	if !aiReady {
		return "Mohon maaf, AI belum diinisialisasi."
	}

	ctx := context.Background()
	instruksiSistem := "Anda adalah FikomBot, asisten virtual resmi Fakultas Ilmu Komputer\nUDB Surakarta"

	mu.Lock()
	chatHistory := userHistories[userID]
	mu.Unlock()

	currentPayload := make([]openai.ChatCompletionMessageParamUnion, 0, len(chatHistory)+2)
	currentPayload = append(currentPayload, openai.SystemMessage(instruksiSistem))
	currentPayload = append(currentPayload, chatHistory...)
	currentPayload = append(currentPayload, openai.UserMessage(userInput))

	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    openai.ChatModel("diisi-saksake"),
		Messages: currentPayload,
	})
	if err != nil {
		return "Mohon maaf, terjadi gangguan saat memproses jawaban."
	}

	jawabanAi := "Mohon maaf, AI tidak memberikan respon"
	if len(resp.Choices) > 0 {
		jawabanAi = resp.Choices[0].Message.Content
	}

	mu.Lock()
	if len(userHistories[userID]) > 10 {
		userHistories[userID] = userHistories[userID][2:]
	}
	userHistories[userID] = append(
		userHistories[userID],
		openai.UserMessage(userInput),
		openai.AssistantMessage(jawabanAi),
	)
	mu.Unlock()

	return jawabanAi
}
