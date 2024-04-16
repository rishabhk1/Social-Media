package web

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"encoding/json"
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// Invoke handles chaincode invoke requests.

func generateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // Token expires after 72 hours
	fmt.Print(secretKey)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func TodayDateTime() string {
	current := time.Now().UTC()
	formattedDate := current.Format("2006-01-02T15:04:05.000Z")
	return formattedDate
}

type ScheduledTask struct {
	Id        string    `json:"id"`
	Execution time.Time `json:"execution"`
}

func loadTasksFromFile() []ScheduledTask {
	data, err := ioutil.ReadFile("tasks.json")
	if err != nil {
		return nil
	}

	var tasks []ScheduledTask
	json.Unmarshal(data, &tasks)

	return tasks
}

func saveTaskToFile(task ScheduledTask) error {
	tasks := loadTasksFromFile()
	tasks = append(tasks, task)

	data, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	return ioutil.WriteFile("tasks.json", data, 0644)
}

func removeTaskByID(id string) error {
	loadTasks := loadTasksFromFile()
	newTasks := make([]ScheduledTask, 0)
	for _, task := range loadTasks {
		if task.Id != id { // Correctly accessing the 'Id' field
			newTasks = append(newTasks, task)
		}
	}
	data, err := json.Marshal(newTasks)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("tasks.json", data, 0644)
}

func (setup *OrgSetup) Invoke(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "CreateCommunity"
	args := r.Form["args"]
	for _, value := range args {
		fmt.Println(value)
	}
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	communityId, err := uuid.NewV7()
	if err != nil {
		fmt.Fprintf(w, "Error community id %s", err)
		return
	}
	dateTime := TodayDateTime()
	newCommunityId := "co" + dateTime + "_" + communityId.String()
	additionalArgs := []string{newCommunityId, dateTime}
	fmt.Println(newCommunityId)
	combinedArgs := append(additionalArgs, args...)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(combinedArgs...))
	if err != nil {
		http.Error(w, "Error in creating community", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error in creating community", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error in creating community", http.StatusInternalServerError)
		fmt.Printf("Error submitting transaction: %s", err)
		return
	}
	fmt.Println(txn_committed.TransactionID())
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	dateTimeObj, err := time.Parse("2006-01-02T15:04:05.000Z", dateTime)
	if err != nil {
		fmt.Printf("Error in converting datetime in community")
		return
	}

	task := ScheduledTask{
		Id:        newCommunityId,
		Execution: dateTimeObj.Add(10 * time.Minute),
	}
	saveTaskToFile(task)
	duration := task.Execution.Sub(time.Now().UTC())
	fmt.Println("duration", duration)
	time.AfterFunc(duration, func() {
		fmt.Print("inside schedule")
		setup.SelectModerator(newCommunityId)

	})
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", txn_endorsed.Result())
}

func (setup *OrgSetup) JoinCommunity(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "JoinCommunity"
	args := r.Form["args"]
	for _, value := range args {
		fmt.Println(value)
	}
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(args...))
	if err != nil {
		http.Error(w, "Error in joining community", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error in joining community", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error in joining community", http.StatusInternalServerError)
		fmt.Printf("Error submitting transaction: %s", err)
		return
	}
	fmt.Println(txn_committed.TransactionID())
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", txn_endorsed.Result())
}

func (setup *OrgSetup) UnJoinCommunity(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "UnJoinCommunity"
	args := r.Form["args"]
	for _, value := range args {
		fmt.Println(value)
	}
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(args...))
	if err != nil {
		http.Error(w, "Error in unjoining community", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error in unjoining community", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error in unjoining community", http.StatusInternalServerError)
		fmt.Printf("Error submitting transaction: %s", err)
		return
	}
	fmt.Println(txn_committed.TransactionID())
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", txn_endorsed.Result())
}

func (setup *OrgSetup) CreatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "CreatePost"
	args := r.Form["args"]
	for _, value := range args {
		fmt.Println(value)
	}
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	postId, err := uuid.NewV7()
	if err != nil {
		fmt.Fprintf(w, "Error post id %s", err)
		return
	}
	dateTime := TodayDateTime()
	newPostId := "p" + dateTime + "_" + postId.String()
	additionalArgs := []string{newPostId, dateTime}
	fmt.Println(newPostId)
	combinedArgs := append(additionalArgs, args...)
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(combinedArgs...))
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, "Error in creating post", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error in creating post", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error in creating post", http.StatusInternalServerError)
		fmt.Printf("Error submitting transaction: %s", err)
		return
	}
	fmt.Println(txn_committed.TransactionID())
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	fmt.Fprintf(w, "%s", txn_endorsed.Result())
}

func (setup *OrgSetup) CreateUser(userId string, username string, email string) error {
	fmt.Println("Received Invoke request")
	// if err := r.ParseForm(); err != nil {
	// 	fmt.Fprintf(w, "ParseForm() err: %s", err)
	// 	return
	// }
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "CreateUser"
	// args := r.Form["args"]
	// for _, value := range args {
	// 	fmt.Println(value)
	// }
	//fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(userId, username, email))
	if err != nil {
		//fmt.Fprintf(w, "Error creating txn proposal: %s", err)
		return err
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		//fmt.Fprintf(w, "Error endorsing txn: %s", err)
		return err
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		//fmt.Fprintf(w, "Error submitting transaction: %s", err)
		return err
	}
	fmt.Println(txn_committed.TransactionID())
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	// w.WriteHeader(http.StatusOK)
	fmt.Printf("%s", txn_endorsed.Result())
	return nil
	//fmt.Fprintf(w, "Transaction ID : %s Response: %s", txn_committed.TransactionID(), txn_endorsed.Result())
}

func (setup *OrgSetup) UpVotePost(w http.ResponseWriter, r *http.Request) {
	//TODO Add map in users struct to check and avoid double upvote
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "UpVotePost"
	args := r.Form["args"]
	for _, value := range args {
		fmt.Println(value)
	}
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(args...))
	if err != nil {
		http.Error(w, "Error in upvoting post", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error in upvoting post", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error in upvoting post", http.StatusInternalServerError)
		fmt.Printf("Error submitting transaction: %s", err)
		return
	}
	fmt.Println(txn_committed.TransactionID())
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	fmt.Fprintf(w, "%s", txn_endorsed.Result())
}

func (setup *OrgSetup) UndoUpVotePost(w http.ResponseWriter, r *http.Request) {
	//TODO Add map in users struct to check and avoid double upvote
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "UndoUpVotePost"
	args := r.Form["args"]
	for _, value := range args {
		fmt.Println(value)
	}
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(args...))
	if err != nil {
		http.Error(w, "Error in upvoting post", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error in upvoting post", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error in upvoting post", http.StatusInternalServerError)
		fmt.Printf("Error submitting transaction: %s", err)
		return
	}
	fmt.Println(txn_committed.TransactionID())
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	fmt.Fprintf(w, "%s", txn_endorsed.Result())
}

func (setup *OrgSetup) DownVotePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "DownVotePost"
	args := r.Form["args"]
	for _, value := range args {
		fmt.Println(value)
	}
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(args...))
	if err != nil {
		http.Error(w, "Error in downvoting post", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error in downvoting post", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error in downvoting post", http.StatusInternalServerError)
		fmt.Printf("Error submitting transaction: %s", err)
		return
	}
	fmt.Println(txn_committed.TransactionID())
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", txn_endorsed.Result())
}

func (setup *OrgSetup) UndoDownVotePost(w http.ResponseWriter, r *http.Request) {
	//TODO Add map in users struct to check and avoid double upvote
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "UndoDownVotePost"
	args := r.Form["args"]
	for _, value := range args {
		fmt.Println(value)
	}
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(args...))
	if err != nil {
		http.Error(w, "Error in downvoting post", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error in downvoting post", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error in downvoting post", http.StatusInternalServerError)
		fmt.Printf("Error submitting transaction: %s", err)
		return
	}
	fmt.Println(txn_committed.TransactionID())
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	fmt.Fprintf(w, "%s", txn_endorsed.Result())
}

func (setup *OrgSetup) CreateComment(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "CreateComment"
	args := r.Form["args"]
	for _, value := range args {
		fmt.Println(value)
	}
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	postId, err := uuid.NewV7()
	if err != nil {
		fmt.Fprintf(w, "Error post id %s", err)
		return
	}
	dateTime := TodayDateTime()
	newPostId := "c" + dateTime + "_" + postId.String()
	additionalArgs := []string{newPostId, dateTime}
	fmt.Println(newPostId)
	combinedArgs := append(additionalArgs, args...)
	w.Header().Set("Content-Type", "application/json")
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(combinedArgs...))
	if err != nil {
		http.Error(w, "Error in creating comment", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error in creating comment", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error in creating comment", http.StatusInternalServerError)
		fmt.Printf("Error submitting transaction: %s", err)
		return
	}
	fmt.Println(txn_committed.TransactionID())
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	fmt.Fprintf(w, "%s", txn_endorsed.Result())
}

func (setup *OrgSetup) DeletePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "DeletePost"
	args := r.Form["args"]
	for _, value := range args {
		fmt.Println(value)
	}
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(args...))
	if err != nil {
		http.Error(w, "Error in deleting post", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error in deleting post", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error in deleting post", http.StatusInternalServerError)
		fmt.Printf("Error submitting transaction: %s", err)
		return
	}
	fmt.Println(txn_committed.TransactionID())
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	fmt.Fprintf(w, "%s", txn_endorsed.Result())
}

func (setup *OrgSetup) AppealPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "AppealPost"
	args := r.Form["args"]
	for _, value := range args {
		fmt.Println(value)
	}
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(args...))
	if err != nil {
		http.Error(w, "Error in appealing post", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error in appealing post", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error in appealing post", http.StatusInternalServerError)
		fmt.Printf("Error submitting transaction: %s", err)
		return
	}
	fmt.Println(txn_committed.TransactionID())
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	fmt.Fprintf(w, "%s", txn_endorsed.Result())
}

func (setup *OrgSetup) HidePostModerator(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "HidePostModerator"
	args := r.Form["args"]
	for _, value := range args {
		fmt.Println(value)
	}
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(args...))
	if err != nil {
		http.Error(w, "Error in hide operation", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error in hide operation", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error in hide operation", http.StatusInternalServerError)
		fmt.Printf("Error submitting transaction: %s", err)
		return
	}
	fmt.Println(txn_committed.TransactionID())
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", txn_endorsed.Result())
	//fmt.Fprintf(w, "Transaction ID : %s Response: %s", txn_committed.TransactionID(), txn_endorsed.Result())
}

func (setup *OrgSetup) ShowPostModerator(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "ShowPostModerator"
	args := r.Form["args"]
	for _, value := range args {
		fmt.Println(value)
	}
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(args...))
	if err != nil {
		http.Error(w, "Error in show operation", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error in show operation", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err.Error())
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error in show operation", http.StatusInternalServerError)
		fmt.Printf("Error submitting transaction: %s", err)
		return
	}
	fmt.Println(txn_committed.TransactionID())
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", txn_endorsed.Result())
	//fmt.Fprintf(w, "Transaction ID : %s Response: %s", txn_committed.TransactionID(), txn_endorsed.Result())
}

func (setup *OrgSetup) SelectModerator(communityId string) {
	fmt.Println("Received mod request")
	// if err := r.ParseForm(); err != nil {
	// 	fmt.Fprintf(w, "ParseForm() err: %s", err)
	// 	return
	// }
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "SelectModerator"
	// args := r.Form["args"]
	// for _, value := range args {
	// 	fmt.Println(value)
	// }
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, communityId)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(communityId))
	if err != nil {
		//fmt.Fprintf(w, "Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		//fmt.Fprintf(w, "Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		//fmt.Fprintf(w, "Error submitting transaction: %s", err)
		return
	}
	//fmt.Fprintf(w, "Transaction ID : %s Response: %s", txn_committed.TransactionID(), txn_endorsed.Result())
	fmt.Println(txn_committed.TransactionID())
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	// w.WriteHeader(http.StatusOK)
	fmt.Printf("%s", txn_endorsed.Result())
	removeTaskByID(communityId)
	sch := loadTasksFromFile()
	fmt.Print(sch)
	task := ScheduledTask{
		Id:        communityId,
		Execution: time.Now().UTC().AddDate(0, 3, 0),
	}
	saveTaskToFile(task)
	duration := task.Execution.Sub(time.Now().UTC())
	time.AfterFunc(duration, func() {
		setup.SelectModerator(communityId)
	})
}

func (setup *OrgSetup) UnAppealPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "UnAppealPost"
	args := r.Form["args"]
	for _, value := range args {
		fmt.Println(value)
	}
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(args...))
	if err != nil {
		fmt.Fprintf(w, "Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		fmt.Fprintf(w, "Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		fmt.Fprintf(w, "Error submitting transaction: %s", err)
		return
	}
	fmt.Fprintf(w, "Transaction ID : %s Response: %s", txn_committed.TransactionID(), txn_endorsed.Result())
}

type Response struct {
	Data map[string]interface{} `json:"data"`
}

func (setup OrgSetup) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Query request")
	//queryParams := r.URL.Query()
	//chainCodeName := queryParams.Get("chaincodeid")
	// chainCodeName := "basic"
	// channelID := "mychannel"
	// function := "GetCommunityAppealed"
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	args := r.Form["args"]
	// fmt.Println(args)
	username := args[0]
	password := args[1]
	fmt.Println(username)
	fmt.Println(password)
	for _, value := range args {
		fmt.Println(value)
	}
	// fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	// network := setup.Gateway.GetNetwork(channelID)
	// contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	// evaluateResponse, err := contract.EvaluateTransaction(function, args, userId, pageNo)
	// if err != nil {
	// 	http.Error(w, "Error", http.StatusInternalServerError)
	// 	// fmt.Fprintf(w, "%s", err)
	// 	fmt.Println(err)
	// 	return
	// }

	// apiEndpoint := "https://oryx-modern-carefully.ngrok-free.app"
	// url := fmt.Sprintf("%s/%s?%s%s&%s%s", apiEndpoint, "LogIn", "RollNo=", username, "Password=", password)
	// fmt.Println(url)
	// resp, err := http.Get(url)
	// if err != nil {
	// 	http.Error(w, "Login Error", http.StatusInternalServerError)
	// 	return
	// }
	// defer resp.Body.Close()
	apiEndpoint := "https://oryx-modern-carefully.ngrok-free.app"
	endpoint := "LogIn"
	formData := url.Values{}
	formData.Set("RollNo", username)
	formData.Set("Password", password)

	// Create the request URL
	url := fmt.Sprintf("%s/%s", apiEndpoint, endpoint)

	// Create the request
	req, err := http.NewRequest("POST", url, strings.NewReader(formData.Encode()))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Login Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	// Read the response body
	loginBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Login Error", http.StatusInternalServerError)
		return
	}
	if resp.StatusCode == http.StatusOK {
		var resultLogin map[string]interface{}
		err = json.Unmarshal(loginBody, &resultLogin)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			http.Error(w, "User data error", http.StatusInternalServerError)
			return
		}
		loginToken, ok := resultLogin["jwttoken"].(string)
		if !ok {
			fmt.Println("Type assertion failed for login token")
			http.Error(w, "Login error", http.StatusInternalServerError)
			return
		}
		fmt.Println("Request was successful!")
		w.WriteHeader(http.StatusOK)
		token, err := generateToken(username)
		if err != nil {
			http.Error(w, "Token Error", http.StatusInternalServerError)
			fmt.Printf("Error submitting transaction: %s", err)
			return
		}
		// var result1 map[string]interface{}
		// err = json.Unmarshal(body1, &result1)
		// if err != nil {
		// 	fmt.Println("Error unmarshaling JSON:", err)
		// 	http.Error(w, "User data error", http.StatusInternalServerError)
		// 	return
		// }
		response := map[string]string{
			// "userId": result1["cryptoKey"].(string),
			"token": token,
		}
		// if username == "1" {
		// 	if password == "abc" {
		// 		response["userId"] = username
		// 		w.WriteHeader(http.StatusOK)
		// 		json.NewEncoder(w).Encode(response)
		// 		return
		// 	} else {
		// 		http.Error(w, "Invalid Credentials", http.StatusInternalServerError)
		// 		return
		// 	}
		// }
		// studentDetailUrl := fmt.Sprintf("%s/%s?%s%s", apiEndpoint, "GetStudentDetails", "RollNo=", username)
		// fmt.Println(studentDetailUrl)
		// resp, err := http.Get(studentDetailUrl)
		// if err != nil {
		// 	http.Error(w, "User data error", http.StatusInternalServerError)
		// 	return
		// }
		studentDetailUrl := fmt.Sprintf("%s/%s?%s%s", apiEndpoint, "GetStudentDetails", "RollNo=", username)

		// Create the request
		req, err := http.NewRequest("GET", studentDetailUrl, nil)
		if err != nil {
			http.Error(w, "User data error", http.StatusInternalServerError)
			return
		}

		// Set the JWT token in the Authorization header
		//req.Header.Set("Authorization", "Bearer "+jwtToken)

		// Add a custom header
		req.Header.Set("jwttoken", loginToken)

		// Send the request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "User data error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "User data error", http.StatusInternalServerError)
			return
		}
		if resp.StatusCode == http.StatusOK {
			fmt.Println(body)
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			if err != nil {
				fmt.Println("Error unmarshaling JSON:", err)
				http.Error(w, "User data error", http.StatusInternalServerError)
				return
			}

			// Access the data from the response
			fmt.Println("Data:", result["name"])
			email, ok := result["email"].(string)
			if !ok {
				fmt.Println("Type assertion failed for email")
				http.Error(w, "User data error", http.StatusInternalServerError)
				return
			}
			// cryptoKey, ok := result["cryptokey"].(string)
			// if !ok {
			// 	fmt.Println("Type assertion failed for email")
			// 	http.Error(w, "User data error", http.StatusInternalServerError)
			// 	return
			// }
			rollno, ok := result["rollno"].(string)
			if !ok {
				fmt.Println("Type assertion failed for email")
				http.Error(w, "User data error", http.StatusInternalServerError)
				return
			}
			response["userId"] = rollno
			err = setup.CreateUser(rollno, rollno, email)
			if err != nil {
				http.Error(w, "User data error", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "User data error", http.StatusInternalServerError)
			return
		}
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	} else if resp.StatusCode == http.StatusUnauthorized {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	} else {
		response := map[string]string{}
		if username == "1" {
			if password == "abc" {
				token, err := generateToken(username)
				if err != nil {
					http.Error(w, "Token Error", http.StatusInternalServerError)
					fmt.Printf("Error submitting transaction: %s", err)
					return
				}
				response["userId"] = username
				response["token"] = token
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(response)
				return
			} else {
				http.Error(w, "Invalid Credentials", http.StatusInternalServerError)
				return
			}
		}
		http.Error(w, "User not registered", http.StatusInternalServerError)
		return
	}

	// Print the response body
	//fmt.Println(string(body))

}
