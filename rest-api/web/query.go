package web

import (
	"fmt"
	"net/http"
)

// Query handles chaincode query requests.
func (setup OrgSetup) Query(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Query request")
	//queryParams := r.URL.Query()
	//chainCodeName := queryParams.Get("chaincodeid")
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "GetCommunity"
	args := r.URL.Query()["id"]
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	evaluateResponse, err := contract.EvaluateTransaction(function, args...)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}

func (setup OrgSetup) GetPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Query request")
	//queryParams := r.URL.Query()
	//chainCodeName := queryParams.Get("chaincodeid")
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "GetPost"
	args := r.URL.Query()["id"]
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	evaluateResponse, err := contract.EvaluateTransaction(function, args...)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}

func (setup OrgSetup) GetComment(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Query request")
	//queryParams := r.URL.Query()
	//chainCodeName := queryParams.Get("chaincodeid")
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "GetComment"
	args := r.URL.Query()["id"]
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	evaluateResponse, err := contract.EvaluateTransaction(function, args...)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}

func (setup OrgSetup) GetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Query request")
	//queryParams := r.URL.Query()
	//chainCodeName := queryParams.Get("chaincodeid")
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "GetUser"
	args := r.URL.Query()["id"]
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	evaluateResponse, err := contract.EvaluateTransaction(function, args...)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}

func (setup OrgSetup) GetUserFeed(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Query request")
	//queryParams := r.URL.Query()
	//chainCodeName := queryParams.Get("chaincodeid")
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "GetUserFeed"
	args := r.URL.Query().Get("id")
	pageNo := r.URL.Query().Get("pageNo")
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	evaluateResponse, err := contract.EvaluateTransaction(function, args, pageNo)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}

func (setup OrgSetup) GetCommentFeed(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Query request")
	//queryParams := r.URL.Query()
	//chainCodeName := queryParams.Get("chaincodeid")
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "GetCommentFeed"
	args := r.URL.Query().Get("id")
	pageNo := r.URL.Query().Get("pageNo")
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	evaluateResponse, err := contract.EvaluateTransaction(function, args, pageNo)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}
