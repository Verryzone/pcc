package ai

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func MulaiChatAi() {
	client := openai.NewClient(
		option.WithBaseURL("http://localhost:8080/v1/"),
		option.WithAPIKey("saksake-karena-gak-butuh-api-key"),
	)
	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)

	var chatHistory []openai.ChatCompletionMessageParamUnion

	for {
		fmt.Print("\nAnda: ")

		userInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\n[Eror membaca input: %v]\n", err)
			continue
		}

		userInput = strings.TrimSpace(userInput)
		if strings.ToLower(userInput) == "keluar" {
			break
		}

		chatHistory = append(chatHistory, openai.UserMessage(userInput))

		fmt.Print("FikomBot: ")

		stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
			Model:    openai.ChatModel("diisi-saksake"),
			Messages: chatHistory,
		})

		var fullResponse strings.Builder

		for stream.Next() {
			chunk := stream.Current()

			if len(chunk.Choices) == 0 {
				continue
			}

			token := chunk.Choices[0].Delta.Content
			fmt.Print(token)
			fullResponse.WriteString(token)
		}

		if err := stream.Err(); err != nil {
			fmt.Printf("\n[Eror saat streaming: %v]\n", err)
			continue
		}

		chatHistory = append(chatHistory, openai.AssistantMessage(fullResponse.String()))
		fmt.Println()
	}
}
