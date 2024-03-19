package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// Invoke handles chaincode invoke requests.
func TodayDateTime() string {
	current := time.Now()
	formattedDate := current.Format("2006-01-02T15:04:05.000Z")
	return formattedDate
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
	fmt.Println(txn_committed.TransactionID())
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
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
	// fmt.Fprintf(w, "Transaction ID : %s Response: %s", txn_committed.TransactionID(), txn_endorsed.Result())
	fmt.Println(txn_committed.TransactionID())
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
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
	// fmt.Fprintf(w, "Transaction ID : %s Response: %s", txn_committed.TransactionID(), txn_endorsed.Result())
	fmt.Println(txn_committed.TransactionID())
	//fmt.Fprintf(w, "%s", txn_committed.TransactionID())
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
	fmt.Println(txn_committed.TransactionID())
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

	//fmt.Fprintf(w, "Transaction ID : %s Response: %s", txn_committed.TransactionID(), txn_endorsed.Result())
	fmt.Println(txn_committed.TransactionID())
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

	//fmt.Fprintf(w, "Transaction ID : %s Response: %s", txn_committed.TransactionID(), txn_endorsed.Result())
	fmt.Println(txn_committed.TransactionID())
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
	//fmt.Fprintf(w, "Transaction ID : %s Response: %s", txn_committed.TransactionID(), txn_endorsed.Result())
	fmt.Println(txn_committed.TransactionID())
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

	//fmt.Fprintf(w, "Transaction ID : %s Response: %s", txn_committed.TransactionID(), txn_endorsed.Result())
	fmt.Println(txn_committed.TransactionID())
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
	fmt.Println(txn_committed.TransactionID())
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
	fmt.Println(txn_committed.TransactionID())
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
	fmt.Println(txn_committed.TransactionID())
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

func (setup *OrgSetup) SelectModerator(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "SelectModerator"
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
