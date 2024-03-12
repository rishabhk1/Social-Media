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
	function := "GetCommunityModified"
	args := r.URL.Query()["id"]
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	evaluateResponse, err := contract.EvaluateTransaction(function, args...)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}
	fmt.Fprintf(w, "%s", evaluateResponse)
}

func (setup OrgSetup) GetPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Query request")
	//queryParams := r.URL.Query()
	//chainCodeName := queryParams.Get("chaincodeid")
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "GetPostModified"
	postId := r.URL.Query().Get("id")
	userId := r.URL.Query().Get("userId")
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args:\n", channelID, chainCodeName, function)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	evaluateResponse, err := contract.EvaluateTransaction(function, postId, userId)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}
	fmt.Fprintf(w, "%s", evaluateResponse)
}

func (setup OrgSetup) GetComment(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Query request")
	//queryParams := r.URL.Query()
	//chainCodeName := queryParams.Get("chaincodeid")
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "GetCommentModified"
	commentId := r.URL.Query().Get("id")
	userId := r.URL.Query().Get("userId")
	fmt.Printf("channel: %s, chaincode: %s, function: %s\n", channelID, chainCodeName, function)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	evaluateResponse, err := contract.EvaluateTransaction(function, commentId, userId)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}
	fmt.Fprintf(w, "%s", evaluateResponse)
}

func (setup OrgSetup) GetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Query request")
	//queryParams := r.URL.Query()
	//chainCodeName := queryParams.Get("chaincodeid")
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "GetUserModified"
	args := r.URL.Query()["id"]
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	evaluateResponse, err := contract.EvaluateTransaction(function, args...)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}
	fmt.Fprintf(w, "%s", evaluateResponse)
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
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s page: %s\n", channelID, chainCodeName, function, args, pageNo)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	evaluateResponse, err := contract.EvaluateTransaction(function, args, pageNo)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}
	fmt.Fprintf(w, "%s", evaluateResponse)
}

func (setup OrgSetup) GetCommentFeed(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Query request")
	//queryParams := r.URL.Query()
	//chainCodeName := queryParams.Get("chaincodeid")
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "GetCommentFeed"
	args := r.URL.Query().Get("parentId")
	pageNo := r.URL.Query().Get("pageNo")
	userId := r.URL.Query().Get("userId")
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	evaluateResponse, err := contract.EvaluateTransaction(function, args, pageNo, userId)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}
	fmt.Fprintf(w, "%s", evaluateResponse)
}

func (setup OrgSetup) GetUserProfilePosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Query request")
	//queryParams := r.URL.Query()
	//chainCodeName := queryParams.Get("chaincodeid")
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "GetUserProfilePosts"
	args := r.URL.Query().Get("targetId")
	userId := r.URL.Query().Get("userId")
	pageNo := r.URL.Query().Get("pageNo")
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	evaluateResponse, err := contract.EvaluateTransaction(function, args, userId, pageNo)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}
	fmt.Fprintf(w, "%s", evaluateResponse)
}

func (setup OrgSetup) GetUserProfileComments(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Query request")
	//queryParams := r.URL.Query()
	//chainCodeName := queryParams.Get("chaincodeid")
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "GetUserProfileComments"
	args := r.URL.Query().Get("targetId")
	userId := r.URL.Query().Get("userId")
	pageNo := r.URL.Query().Get("pageNo")
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	evaluateResponse, err := contract.EvaluateTransaction(function, args, userId, pageNo)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}
	fmt.Fprintf(w, "%s", evaluateResponse)
}

func (setup OrgSetup) GetCommunityPosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Query request")
	//queryParams := r.URL.Query()
	//chainCodeName := queryParams.Get("chaincodeid")
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "GetCommunityPosts"
	args := r.URL.Query().Get("communityId")
	userId := r.URL.Query().Get("userId")
	pageNo := r.URL.Query().Get("pageNo")
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	evaluateResponse, err := contract.EvaluateTransaction(function, args, userId, pageNo)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}
	fmt.Fprintf(w, "%s", evaluateResponse)
}

func (setup OrgSetup) GetCommunityAppealed(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Query request")
	//queryParams := r.URL.Query()
	//chainCodeName := queryParams.Get("chaincodeid")
	chainCodeName := "basic"
	channelID := "mychannel"
	function := "GetCommunityAppealed"
	args := r.URL.Query().Get("communityId")
	userId := r.URL.Query().Get("userId")
	pageNo := r.URL.Query().Get("pageNo")
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	w.Header().Set("Content-Type", "application/json")
	evaluateResponse, err := contract.EvaluateTransaction(function, args, userId, pageNo)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}
	fmt.Fprintf(w, "%s", evaluateResponse)
}
