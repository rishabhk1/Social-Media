package web

import (
	"fmt"
	"net/http"
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
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
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
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
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
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
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
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error submitting transaction: %s", err)
		return
	}
	fmt.Println(txn_committed.TransactionID())
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
	fmt.Fprintf(w, "%s", txn_endorsed.Result())
}

func (setup *OrgSetup) CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "CreateUser"
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
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
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
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
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
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
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
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
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
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
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
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
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
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
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
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
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
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Printf("Error endorsing txn: %s", err.Error())
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
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

func (setup OrgSetup) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Query request")
	//queryParams := r.URL.Query()
	//chainCodeName := queryParams.Get("chaincodeid")
	// chainCodeName := "basic"
	// channelID := "mychannel"
	// function := "GetCommunityAppealed"
	args := r.Form["args"]
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
	token, err := generateToken("1")
	if err != nil {
		http.Error(w, "Token Error", http.StatusInternalServerError)
		fmt.Printf("Error submitting transaction: %s", err)
		return
	}
	response := map[string]string{
		"userId": "1",
		"token":  token,
	}

	w.WriteHeader(http.StatusOK)

	// Send the response
	json.NewEncoder(w).Encode(response)
}
