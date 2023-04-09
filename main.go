package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	lambda.Start(requestHandler)
}

func requestHandler(ctx context.Context) (string, error) {
	quote, err := generateQuote()

	if err != nil {
		log.Panicf("generateQuote error: %v", err)
	}

	err = tweetQuote(quote)

	if err != nil {
		log.Panicf("tweetQuote error: %v", err)
	}

	return "success", nil
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
					Content: `あなたは架空の歴史上の偉人の名言をつぶやくTwitterBOTと仮定してください。
					その偉人は、困難に直面したり、障害に立ち向かう人々を勇気づけるための名言をたくさん残しています。
					あなたはTwitterBOTとして、その名言を以下の形式で1つだけつぶやいてください。
					なお、文字数は140字以下として下さい。
					'''
					「${名言}」 
					ChatGPT (${職業})
					'''
					
					${名言}はあなたの名言です。「ですます」調を使ってはいけません。
					
					${職業}はあなたの職業です。なお、「架空」「歴史上」「偉人」「成功者」「先人」という言葉は使ってはいけません。
					具体的な職業名を使ってください。
					`,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	quote = resp.Choices[0].Message.Content

	fmt.Printf("generateQuote result: %v\n", quote)
	return quote, nil
}

func tweetQuote(quote string) (err error) {
	in := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           os.Getenv("GOTWI_ACCESS_TOKEN"),
		OAuthTokenSecret:     os.Getenv("GOTWI_ACCESS_TOKEN_SECRET"),
	}

	c, err := gotwi.NewClient(in)
	if err != nil {
		return err
	}

	p := &types.CreateInput{
		Text: gotwi.String(quote),
	}

	res, err := managetweet.Create(context.Background(), c, p)
	if err != nil {
		return err
	}

	fmt.Printf("tweetQuote result: [%s] %s\n", gotwi.StringValue(res.Data.ID), gotwi.StringValue(res.Data.Text))

	return nil
}
