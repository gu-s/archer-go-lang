/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type Question struct {
	Id       string
	Question string
	Answers  []string
}

type UserAnswer struct {
	UserId            string         `json:"user_id"`
	AnswersByQuestion map[string]int `json:"answers_by_question"` // key = questionId , value  = answerId
}
type AnswersApiResponse struct {
	Message        string `json:"message"`
	CorrectAnswers int    `json:"correct_answers"`
}

type promptContent struct {
	errorMsg string
	label    string
}

var quizData []Question

// quizCmd represents the quiz command
var quizCmd = &cobra.Command{
	Use:   "quiz",
	Short: "A brief description of your command",
	Long:  `the command allows you to receive the questions and answer.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("quiz called")
		getQuiz()

	},
}

func init() {
	rootCmd.AddCommand(quizCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// quizCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// quizCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func postAnswers(body UserAnswer) AnswersApiResponse {

	json_data, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(
		"http://127.0.0.1:8000/answers", //url
		"application/json",
		bytes.NewBuffer(json_data), //body
	)

	if err != nil {
		log.Printf("Could not request a data. %v", err)
	}

	var res AnswersApiResponse

	json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println(res)

	return res
}

func getData(baseAPI string) []byte {

	request, err := http.NewRequest(
		http.MethodGet, //method
		baseAPI,        //url
		nil,            //body
	)

	if err != nil {
		log.Printf("Could not request a data. %v", err)
	}

	request.Header.Add("Accept", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Could not make a request. %v", err)
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Could not read response body. %v", err)
	}

	return responseBytes
}

func getQuiz() {
	url := "http://127.0.0.1:8000/questions"
	responseBytes := getData(url)

	if err := json.Unmarshal(responseBytes, &quizData); err != nil {
		fmt.Printf("Could not unmarshal reponseBytes. %v", err)
	}

	UserIdPromptContent := promptContent{
		"Please provide a userName.",
		"please write your username?",
	}

	userId := promptGetInput(UserIdPromptContent)
	//userId := promptGetInput(UserIdPromptContent)

	var userAnswers UserAnswer
	userAnswers.UserId = userId
	userAnswers.AnswersByQuestion = map[string]int{}
	for _, element := range quizData {
		quizSelectPromptContent := promptContent{
			"Please provide a answer.",
			element.Question,
		}
		items := element.Answers
		var _, answerIndex = promptGetSelect(quizSelectPromptContent, items)
		userAnswers.AnswersByQuestion[element.Id] = answerIndex

	}

	confirmPostPromptContent := promptContent{
		"Please provide a answer.",
		"Do You want post your answers?",
	}
	confirmItems := []string{"NO", "YES"}
	var result, _ = promptGetSelect(confirmPostPromptContent, confirmItems)
	if result == confirmItems[1] {
		var res = postAnswers(userAnswers)
		println("Number of correct answers:", res.CorrectAnswers)
	} else {
		println("answers not sent")
	}
}

func promptGetSelect(pc promptContent, items []string) (string, int) {
	//items := q.Answers
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label: pc.label,
			Items: items,
		}

		index, result, err = prompt.Run()

		if index == -1 {
			items = append(items, result)
		}
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Input Index: %s\n", index)

	return result, index
}

func promptGetInput(pc promptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.errorMsg)
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    pc.label,
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Input: %s\n", result)

	return result
}
