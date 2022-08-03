package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Question struct {
	Id            string
	Question      string
	Answers       []string
	correctAnswer int
}

type UserAnswer struct {
	UserId            string         `json:"user_id"`
	AnswersByQuestion map[string]int `json:"answers_by_question"` // key = questionId , value  = answerId
}
type ScoreRateApiResponse struct {
	RatePercentageFromOtherUsers float32 `json:"rate_percentage_from_other_users"`
}

var correctAnswersByUser = map[string]int{} // key = UserId , value  = number of correect answers

var quiz []Question

func main() {

	quiz = []Question{
		{
			Id:            "a",
			Question:      "how many months has a year?",
			correctAnswer: 2, //index of correct answer
			Answers:       []string{"10", "14", "12", "13"},
		},
		{
			Id:            "b",
			Question:      "What is the capital of Malta?",
			correctAnswer: 1, //index of correct answer
			Answers:       []string{"Sliema", "Valletta", "Mdina", "Saint Julians"},
		},
		{
			Id:            "c",
			Question:      "Which actor played Rocky?",
			correctAnswer: 1, //index of correct answer
			Answers:       []string{"Tony Burton", "Sylvester Stallone", "Harrison Ford", "Jason Statham"},
		},
		{
			Id:            "d",
			Question:      "What is the capital city of Australia?",
			correctAnswer: 2, //index of correct answer
			Answers:       []string{"Sydney", "Melbourne", "Canberra", "Brisbane"},
		},
	}

	http.HandleFunc("/", home)

	http.HandleFunc("/questions", questions)
	http.HandleFunc("/answers", answers)
	http.HandleFunc("/score_rate", scoreRate)

	http.ListenAndServe(":8000", nil)

}

func home(w http.ResponseWriter, r *http.Request) {
	html := "<html>"
	html += "<body>"
	html += "<h1>Hello world</h1>"
	html += "</body>"
	html += "</html>"
	w.Write([]byte(html))
}

func questions(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//todo: get json
	//mockup { number_questions: 2, questions_answers : [ { question: "question 1 ? " , answers : ["1)ans a", "2","3)","4)" ] }, ... ] }

	jsonResp, err := json.Marshal(quiz)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

}

func answers(w http.ResponseWriter, r *http.Request) {

	//todo check a answer by question

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	var answers UserAnswer
	err = json.Unmarshal(body, &answers)
	if err != nil {
		panic(err)
	}
	log.Println(answers.UserId)

	//todo return how many answerts are correct
	var numberCorrectAnswers = 0
	for _, element := range quiz {
		//		fmt.Println(index, "=>", element)
		if element.correctAnswer == answers.AnswersByQuestion[element.Id] {
			numberCorrectAnswers++
		}

	}

	correctAnswersByUser[answers.UserId] = numberCorrectAnswers

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	/*
		fmt.Fprintln(w, "{ \"hola\":1 }")
	*/

	// The same json tags will be used to encode data into JSON
	type AnswersApiResponse struct {
		Message        string `json:"message"`
		CorrectAnswers int    `json:"correct_answers"`
	}
	resp := new(AnswersApiResponse)
	resp.Message = "OK"
	resp.CorrectAnswers = numberCorrectAnswers

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return

}

func scoreRate(w http.ResponseWriter, r *http.Request) {
	userIdParam, ok := r.URL.Query()["user"]
	if !ok || len(userIdParam) != 1 {
		log.Fatalf("Error happened in JSON marshal. Err: %s", ok)
	}
	var userId string

	if len(userIdParam) == 1 {
		userId = string(userIdParam[0])
		log.Println("userId: ", userId)
	}

	//todo calc the score rate by user
	var numOfUserWithLessScore = 0
	for key, element := range correctAnswersByUser {
		fmt.Println("Key:", key, "=>", "Element:", element)
		if element < correctAnswersByUser[userId] {
			numOfUserWithLessScore++
		}
	}

	var percentageOfUserWithLessRate float32 = 0.00
	if numOfUserWithLessScore > 0 && len(correctAnswersByUser) > 0 {
		percentageOfUserWithLessRate = float32(numOfUserWithLessScore) / float32(len(correctAnswersByUser))
	}

	resp := new(ScoreRateApiResponse)
	resp.RatePercentageFromOtherUsers = percentageOfUserWithLessRate

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}
