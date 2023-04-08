package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
)

func main() {
	quote, err := generateQuote()

	if err != nil {
		log.Panic(err)
	}

	fmt.Println(quote)
}

func generateQuote() (quote string, err error) {
	client := openai.NewClient(os.Getenv("OPEN_AI_API_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					Content: `あなたは歴史上の偉人と仮定してください。
					あなたは、困難に直面したり、障害に立ち向かう人々を勇気づけるための名言をたくさん残しています。
					それを以下の形式で1つだけ出力してください。
					また、この形式以外の説明や解説を出力しないでください。（重要）
					'''
					${名言} - ChatGPT (${職業})
					'''

					${名言}はあなたの名言です。100-150文字程度で記入してください。「ですます」調を使ってはいけません。
					${職業}はあなたの職業です。
					`,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
