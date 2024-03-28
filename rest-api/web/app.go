package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
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

func startTunnel(ctx context.Context) error {
	listener, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(
			config.WithDomain("raptor-rested-evidently.ngrok-free.app"),
		),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		return err
	}

	// Return the public URL of the tunnel
	fmt.Println("App URL", listener.URL())
	return http.Serve(listener, nil)
}

// Serve starts http web server.
func Serve(setups OrgSetup) {
	// ctx := context.Background()
	// listener, err := startTunnel(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// }

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
	http.HandleFunc("/post/undo_upvote", setups.UndoUpVotePost)
	http.HandleFunc("/post/undo_downvote", setups.UndoDownVotePost)
	http.HandleFunc("/create/comment", setups.CreateComment)
	http.HandleFunc("/feed", setups.GetUserFeed)
	http.HandleFunc("/comment/feed", setups.GetCommentFeed)
	http.HandleFunc("/user_profile/posts", setups.GetUserProfilePosts)
	http.HandleFunc("/user_profile/comments", setups.GetUserProfileComments)
	http.HandleFunc("/community/posts", setups.GetCommunityPosts)
	http.HandleFunc("/community/appealed", setups.GetCommunityAppealed)
	http.HandleFunc("/community/name", setups.GetCommunityName)
	http.HandleFunc("/delete", setups.DeletePost)
	http.HandleFunc("/appeal", setups.AppealPost)
	http.HandleFunc("/hide", setups.HidePostModerator)
	// http.HandleFunc("/moderator", setups.SelectModerator)
	http.HandleFunc("/show", setups.ShowPostModerator)
	http.HandleFunc("/unappeal", setups.UnAppealPost)
	//fmt.Printf("Listening (%s)...\n", listener.URL())
	// if err := http.Serve(listener, nil); err != nil {
	// 	fmt.Println(err)
	// }
	if err := startTunnel(context.Background()); err != nil {
		fmt.Println(err)
	}
}
