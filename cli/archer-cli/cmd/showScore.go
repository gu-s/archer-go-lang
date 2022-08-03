/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

type ScoreRateApiResponse struct {
	RatePercentageFromOtherUsers float32 `json:"rate_percentage_from_other_users"`
}

var scoreResponse ScoreRateApiResponse

// showScoreCmd represents the showScore command
var showScoreCmd = &cobra.Command{
	Use:   "showScore",
	Short: "A brief description of your command",
	Long:  `Shows the rated of a user compared to others that have taken the quiz.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("showScore called")
		getScore()
	},
}

func init() {
	quizCmd.AddCommand(showScoreCmd)

}

func getScore() {

	UserIdPromptContent := promptContent{
		"Please provide a userName.",
		"please enter the username to show the score",
	}

	userId := promptGetInput(UserIdPromptContent)

	fmt.Printf("Input: %s\n", userId)

	url := "http://127.0.0.1:8000/score_rate?user=" + userId
	responseBytes := getData(url)

	if err := json.Unmarshal(responseBytes, &scoreResponse); err != nil {
		fmt.Printf("Could not unmarshal reponseBytes. %v", err)
	}

	fmt.Printf("You scored higher than %.2f %% of all quizzers\n", scoreResponse.RatePercentageFromOtherUsers*100)

}
