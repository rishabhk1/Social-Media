package web

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
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

var secretKey = []byte(os.Getenv("secretKey"))

func verifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	} else {
		return "", fmt.Errorf("invalid token")
	}
}

func AuthMiddleware(handlerFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the Authorization header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		// Verify the token
		_, err := verifyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// If the token is valid, proceed with the request
		handlerFunc(w, r)
	}
}

// Serve starts http web server.
func Serve(setups OrgSetup) {
	// ctx := context.Background()
	// listener, err := startTunnel(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	http.HandleFunc("/channel", AuthMiddleware(http.HandlerFunc(setups.Query)))
	http.HandleFunc("/post", AuthMiddleware(http.HandlerFunc(setups.GetPost)))
	http.HandleFunc("/user", AuthMiddleware(http.HandlerFunc(setups.GetUser)))
	http.HandleFunc("/comment", AuthMiddleware(http.HandlerFunc(setups.GetComment)))
	http.HandleFunc("/create/channel", AuthMiddleware(http.HandlerFunc(setups.Invoke)))
	http.HandleFunc("/create/user", AuthMiddleware(http.HandlerFunc(setups.CreateUser)))
	http.HandleFunc("/channel/join", AuthMiddleware(http.HandlerFunc(setups.JoinCommunity)))
	http.HandleFunc("/channel/unjoin", AuthMiddleware(http.HandlerFunc(setups.UnJoinCommunity)))
	http.HandleFunc("/create/post", AuthMiddleware(http.HandlerFunc(setups.CreatePost)))
	http.HandleFunc("/post/upvote", AuthMiddleware(http.HandlerFunc(setups.UpVotePost)))
	http.HandleFunc("/post/downvote", AuthMiddleware(http.HandlerFunc(setups.DownVotePost)))
	http.HandleFunc("/post/undo_upvote", AuthMiddleware(http.HandlerFunc(setups.UndoUpVotePost)))
	http.HandleFunc("/post/undo_downvote", AuthMiddleware(http.HandlerFunc(setups.UndoDownVotePost)))
	http.HandleFunc("/create/comment", AuthMiddleware(http.HandlerFunc(setups.CreateComment)))
	http.HandleFunc("/feed", AuthMiddleware(http.HandlerFunc(setups.GetUserFeed)))
	http.HandleFunc("/comment/feed", AuthMiddleware(http.HandlerFunc(setups.GetCommentFeed)))
	http.HandleFunc("/user_profile/posts", AuthMiddleware(http.HandlerFunc(setups.GetUserProfilePosts)))
	http.HandleFunc("/user_profile/comments", AuthMiddleware(http.HandlerFunc(setups.GetUserProfileComments)))
	http.HandleFunc("/community/posts", AuthMiddleware(http.HandlerFunc(setups.GetCommunityPosts)))
	http.HandleFunc("/community/appealed", AuthMiddleware(http.HandlerFunc(setups.GetCommunityAppealed)))
	http.HandleFunc("/community/name", AuthMiddleware(http.HandlerFunc(setups.GetCommunityName)))
	http.HandleFunc("/delete", AuthMiddleware(http.HandlerFunc(setups.DeletePost)))
	http.HandleFunc("/appeal", AuthMiddleware(http.HandlerFunc(setups.AppealPost)))
	http.HandleFunc("/hide", AuthMiddleware(http.HandlerFunc(setups.HidePostModerator)))
	// http.HandleFunc("/moderator", setups.SelectModerator)
	http.HandleFunc("/show", AuthMiddleware(http.HandlerFunc(setups.ShowPostModerator)))
	http.HandleFunc("/unappeal", AuthMiddleware(http.HandlerFunc(setups.UnAppealPost)))
	http.HandleFunc("/login", setups.Login)
	//fmt.Printf("Listening (%s)...\n", listener.URL())
	// if err := http.Serve(listener, nil); err != nil {
	// 	fmt.Println(err)
	// }
	if err := startTunnel(context.Background()); err != nil {
		fmt.Println(err)
	}
}
