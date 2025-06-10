package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"go-chat/prompt"
	"log"
	"os"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func main() {
	modeFlag := flag.String("mode", "", "Select prompt mode: study, exam_prep, focus_mode")
	promptPath := flag.String("config", "prompt/system_prompts.yaml", "Path to system prompts config")
	modelName := flag.String("model", "tinyllama:latest", "Give the local model name")
	flag.Parse()

	pm, err := prompt.LoadSystemPrompts(*promptPath)
	if err != nil {
		log.Fatal("Failed to load system prompts:", err)
	}

	mode := *modeFlag
	if mode == "" {
		mode = pm.Default
	}

	fmt.Printf("🤖 Using model: %s\n", *modelName)
	fmt.Printf("🧠 Mode: %s\n", mode)
	fmt.Println("💡 Type 'exit' to end the chat, or 'clear' to clear the chat history.")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("💬 What's your study goal? ")
	initialInput, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading input:", err)
	}

	initialInput = strings.TrimSpace(initialInput)
	systemPrompt := pm.BuildPrompt(mode, initialInput)

	llm, err := ollama.New(ollama.WithModel(*modelName))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, initialInput),
	}

	for {
		fmt.Print("\nYou: ")
		userInput, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Error reading input:", err)
		}

		userInput = strings.TrimSpace(userInput)
		if userInput == "exit" {
			fmt.Println("👋 Exiting chat.")
			break
		}
		if userInput == "clear" {
			fmt.Println("🔄 History cleared!")
			messages = messages[:1]
			continue
		}

		messages = append(messages, llms.TextParts(llms.ChatMessageTypeHuman, userInput))

		fmt.Print("AI: ")
		_, err = llm.GenerateContent(ctx, messages, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}))
		if err != nil {
			log.Fatal(err)
		}
	}
}
