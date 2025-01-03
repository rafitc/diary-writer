package editor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/core"
	"main/models"
	"net/http"
	"strings"
)

// This is the editor package
// This edits the content using LLM model
// Currently using Groq API to get edit the content
// Same LLM uses to build the Headline and grammar correction

func EditContent(content string, date string) (string, string, string) {

	content, contentForHeadline := processContent(content, date)
	header_format, title, summary := generateTitleAndSummary(contentForHeadline, date)

	content = header_format + content

	return content, title, summary
}

func generateTitleAndSummary(content string, date string) (string, string, string) {
	var prompt string
	// Add the title to the content
	prompt = "Please generate a concise and engaging title for the following content. The title should be brief, ideally one line, and should capture the essence of the content. Provide only the title in your response without any special characters."
	title := editContentUsingLLM(content, date, prompt, "title")

	prompt = "Please generate a concise and engaging summary for the following content. The summary should be brief, ideally less than or around 35 words, and should capture the essence of the content. Provide only the summary in your response without any special characters."
	summary := editContentUsingLLM(content, date, prompt, "summary")

	// Add the static part of the content on top of the content
	header_format := fmt.Sprintf("---\ntitle: \"%s\"\npublishedAt: \"%s\"\nsummary: \"%s\"\n---\n\n", title, date, summary)
	return header_format, title, summary
}

func processContent(content string, date string) (string, string) {
	lines := strings.Split(content, "\n")
	var finalContent bytes.Buffer
	var contentForHeadline bytes.Buffer

	prompt := "I am sending you content to fix the grammar. In response, don't change the content structure; just fix the grammar. Provide only the modified content as the response. Don't change the newline characters in the content."
	buffer_string := ""
	for _, line := range lines {
		if strings.Contains(line, "![Alt Text]") {
			// Send accumulated content to the API
			respFromLLM := editContentUsingLLM(buffer_string, date, prompt, "edit")
			// now clear the buffer
			buffer_string = ""
			finalContent.WriteString(respFromLLM + "\n\n")
			finalContent.WriteString(line + "\n\n")

			contentForHeadline.WriteString(respFromLLM + "\n")
			continue
		}
		if line != "" {
			buffer_string += line + "\n"
			continue
		}
	}

	// Handle remaining content after the last image marker
	if buffer_string != "" {
		respFromLLM := editContentUsingLLM(buffer_string, date, prompt, "edit")
		buffer_string = ""
		finalContent.WriteString(respFromLLM + "\n\n")

		contentForHeadline.WriteString(respFromLLM + "\n")
	}
	finalText := finalContent.String()

	return finalText, contentForHeadline.String()
}

func editContentUsingLLM(content string, date string, prompt string, typeOfContent string) string {
	// Set the default content based on the type of content
	var defaultContent string
	if typeOfContent == "edit" {
		defaultContent = content
	} else if typeOfContent == "title" {
		defaultContent = "My day on " + date
	} else if typeOfContent == "summary" {
		defaultContent = "My day on " + date
	}

	// Call the LLM model
	apiURL := "https://api.groq.com/openai/v1/chat/completions"
	apiKey := core.Config.GEN_AI.GROQ_API_KEY

	// content = "im very poor boy. \n my friend's is good"
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": core.Config.GEN_AI.MODEL_NAME,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt + " \n\n" + content,
			},
		},
	})
	if err != nil {
		fmt.Printf("Error creating request body: %v\n", err)
		return defaultContent
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return defaultContent
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making API call: %v\n", err)
		return defaultContent
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API call failed with status: %v\n", resp.Status)
		return defaultContent
	}

	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return defaultContent
	}
	var chatCompletion models.ChatCompletion

	contentFromLLM := ""
	responseBodyBytes, err := json.Marshal(responseBody)
	if err != nil {
		fmt.Printf("Error marshalling response body: %v\n", err)
		return defaultContent
	}
	err = json.Unmarshal(responseBodyBytes, &chatCompletion)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return defaultContent
	}
	contentFromLLM = chatCompletion.Choices[0].Message.Content

	if contentFromLLM != "" {
		return contentFromLLM
	}
	return defaultContent
}
