package web

import (
	"fmt"
	"net/http"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// OrgSetup contains organization's config to interact with the network.
type OrgSetup struct {
	OrgName      string
	MSPID        string
	CryptoPath   string
	CertPath     string
	KeyPath      string
	TLSCertPath  string
	PeerEndpoint string
	GatewayPeer  string
	Gateway      client.Gateway
}

// Serve starts http web server.
func Serve(setups OrgSetup) {
	http.HandleFunc("/channel", setups.Query)
	http.HandleFunc("/post", setups.GetPost)
	http.HandleFunc("/user", setups.GetUser)
	http.HandleFunc("/comment", setups.GetComment)
	http.HandleFunc("/create/channel", setups.Invoke)
	http.HandleFunc("/create/user", setups.CreateUser)
	http.HandleFunc("/channel/join", setups.JoinCommunity)
	http.HandleFunc("/channel/unjoin", setups.UnJoinCommunity)
	http.HandleFunc("/create/post", setups.CreatePost)
	http.HandleFunc("/post/upvote", setups.UpVotePost)
	http.HandleFunc("/post/downvote", setups.DownVotePost)
	http.HandleFunc("/create/comment", setups.CreateComment)
	http.HandleFunc("/feed", setups.GetUserFeed)
	http.HandleFunc("/comment/feed", setups.GetCommentFeed)
	http.HandleFunc("/delete", setups.DeletePost)
	http.HandleFunc("/appeal", setups.AppealPost)
	http.HandleFunc("/hide", setups.HidePostModerator)
	http.HandleFunc("/moderator", setups.SelectModerator)
	http.HandleFunc("/show", setups.ShowPostModerator)
	http.HandleFunc("/unappeal", setups.UnAppealPost)
	fmt.Println("Listening (http://localhost:3000/)...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Println(err)
	}
}
