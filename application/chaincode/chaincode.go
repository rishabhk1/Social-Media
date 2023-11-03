package chaincode

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type User struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Communities []string `json:"communities"`
	Posts       []string `json:"posts"` //list of ids
	Reputation  int
}

type Community struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Creator     string    `json:"creator"`
	CreatedAt   time.Time `json:"createdAt"`
	Moderators  []string  `json:"moderators"`
	Users       []string  `json:"users"`
	Posts       []string  `json:"posts"` //list of ids
	Appealed    []string
}

type Post struct {
	ID        string    `json:"id"` // should start with p
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"createdAt"`
	Comments  []string  `json:"comments"` //list of ids
	Hidden    bool
	Community string
	HideCount int
	ShowCount int
}

type Comment struct {
	ID        string    `json:"id"` //should syart with c
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"createdAt"`
	Parent    string    //can be comment or id
	Replies   []string  `json:"replies"` //list of ids
	Hidden    bool
	Community string
	HideCount int
	ShowCount int
}

const PostsPerPage = 20
const CommentsPerPage = 20

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func contains(slice []string, target string) bool {
	for _, element := range slice {
		if element == target {
			return true
		}
	}
	return false
}

func removeElement(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

func findIndex(slice []string, target string) int {
	for index, value := range slice {
		if value == target {
			return index
		}
	}
	return -1 // Return -1 if the element is not found
}

/*
InitLedger is used to setup initial data on the blockchain for interaction
*/
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	layout := "2006-01-02T15:04:05.000Z"

	user1 := User{
		ID:          "1",
		Username:    "john_doe",
		Email:       "john@example.com",
		Communities: []string{"co_1"},
		Posts:       []string{"p_1"},
		Reputation:  10,
	}

	user2 := User{
		ID:          "2",
		Username:    "joe_rog",
		Email:       "joe_rog@example.com",
		Communities: []string{"co_1"},
		Posts:       []string{"p_2"},
		Reputation:  11,
	}

	ctime, _ := time.Parse(layout, "2023-01-07T15:04:05.000Z")
	community := Community{
		ID:          "co_1",
		Name:        "Comm1",
		Description: "Demo Community",
		Creator:     "1",
		CreatedAt:   ctime,
		Moderators:  []string{"1"},
		Posts:       []string{"p_1", "p_2"},
		Users:       []string{"1", "2"},
		Appealed:    make([]string, 0),
	}

	ctime, _ = time.Parse(layout, "2023-04-17T15:04:05.000Z")
	post1 := Post{
		ID:        "p_1",
		Title:     "Demo Post",
		Content:   "Demo post content",
		Author:    "1",
		Score:     10,
		CreatedAt: ctime,
		Comments:  make([]string, 0),
		Hidden:    false,
		Community: "co_1",
		HideCount: 0,
		ShowCount: 0,
	}

	ctime, _ = time.Parse(layout, "2023-07-17T15:04:05.000Z")
	post2 := Post{
		ID:        "p_2",
		Title:     "Demo Post2",
		Content:   "Demo post content2",
		Author:    "2",
		Score:     11,
		CreatedAt: ctime,
		Comments:  make([]string, 0),
		Hidden:    false,
		Community: "co_1",
		HideCount: 0,
		ShowCount: 0,
	}

	user1JSON, err := json.Marshal(user1)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(user1.ID, user1JSON)
	if err != nil {
		return err
	}

	user2JSON, err := json.Marshal(user2)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(user2.ID, user2JSON)
	if err != nil {
		return err
	}

	communityJson, err := json.Marshal(community)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(community.ID, communityJson)
	if err != nil {
		return err
	}

	post1Json, err := json.Marshal(post1)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(post1.ID, post1Json)
	if err != nil {
		return err
	}

	post2Json, err := json.Marshal(post2)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(post2.ID, post2Json)
	if err != nil {
		return err
	}
	return nil
}

/*
Used to create a new user within a blockchain.
It checks if the user already exists and, if not, initializes a new user with the provided user Id, username, and email.
*/
func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, UserId string, username string, email string) error {
	existingUser, err := s.GetUser(ctx, UserId)
	if err == nil && existingUser != nil {
		return fmt.Errorf("User with ID %s already exists", UserId)
	}
	user := User{
		ID:          UserId,
		Username:    username,
		Email:       email,
		Communities: make([]string, 0),
		Posts:       make([]string, 0),
		Reputation:  0,
	}
	userJson, _ := json.Marshal(user)
	return ctx.GetStub().PutState(UserId, userJson)

}

/*
Used to retrieve User information from a  blockchain.
It takes a user Id as a parameter and queries the ledger to fetch the corresponding user data.
*/
func (s *SmartContract) GetUser(ctx contractapi.TransactionContextInterface, userId string) (*User, error) {
	userJson, err := ctx.GetStub().GetState(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to read user from ledger: %w", err)
	}
	if userJson == nil {
		return nil, nil
	}

	var user User
	err = json.Unmarshal(userJson, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

/*
Used to create a new community in a blockchain.
It takes parameters such as the community's Id, name, description, creator, and creation timestamp.
Function also makes the creator, the initial moderator of the community.
*/
func (s *SmartContract) CreateCommunity(ctx contractapi.TransactionContextInterface, id string, name string, description string, creator string, createdAt string) error {
	existingCommunity, err := s.GetCommunity(ctx, id)
	if err == nil && existingCommunity != nil {
		return fmt.Errorf("Community with ID %s already exists", id)
	}
	existingUser, err := s.GetUser(ctx, creator)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return fmt.Errorf("User with ID %s doesn't exists", creator)
	}
	layout := "2006-01-02T15:04:05.000Z"
	currentTime, _ := time.Parse(layout, createdAt)
	community := Community{
		ID:          id,
		Name:        name,
		Description: description,
		Creator:     creator,
		CreatedAt:   currentTime,
		Moderators:  []string{creator},
		Users:       make([]string, 0),
		Posts:       make([]string, 0),
		Appealed:    make([]string, 0),
	}
	community.Users = append(community.Users, creator)
	existingUser.Communities = append(existingUser.Communities, id)
	communityJson, _ := json.Marshal(community)
	ctx.GetStub().PutState(id, communityJson)
	UserJson, _ := json.Marshal(existingUser)
	ctx.GetStub().PutState(creator, UserJson)
	return nil

}

/*
Used to retrieve information about a community from blockchain.
It takes the community's Id as a parameter and retrieves the community's data from the blockchain.
*/
func (s *SmartContract) GetCommunity(ctx contractapi.TransactionContextInterface, id string) (*Community, error) {
	communityJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read community from ledger: %w", err)
	}
	if communityJson == nil {
		return nil, nil
	}

	var community Community
	err = json.Unmarshal(communityJson, &community)
	if err != nil {
		return nil, err
	}
	return &community, nil
}

/*
Used to allow a user to join a specific community.
It takes the user's Id and the community's Id as parameters and adds the user to the list of community members
*/
func (s *SmartContract) JoinCommunity(ctx contractapi.TransactionContextInterface, communityId string, userId string) error {
	existingCommunity, err := s.GetCommunity(ctx, communityId)
	if err != nil {
		return err
	}
	if existingCommunity == nil {
		return fmt.Errorf("Community with ID %s doesn't exists", communityId)
	}
	existingCommunity.Users = append(existingCommunity.Users, userId)
	communityJson, _ := json.Marshal(existingCommunity)
	ctx.GetStub().PutState(communityId, communityJson)
	return nil
}

/*
Used to allow a user to leave a specific community.
It takes the user's Id and the community's Id as parameters and removes the user to the list of community members
*/
func (s *SmartContract) UnJoinCommunity(ctx contractapi.TransactionContextInterface, communityId string, userId string) error {
	existingCommunity, err := s.GetCommunity(ctx, communityId)
	if err != nil {
		return err
	}
	if existingCommunity == nil {
		return fmt.Errorf("Community with ID %s doesn't exists", communityId)
	}
	existingCommunity.Users = removeElement(existingCommunity.Users, findIndex(existingCommunity.Users, userId))
	communityJson, _ := json.Marshal(existingCommunity)
	ctx.GetStub().PutState(communityId, communityJson)
	return nil
}

/*
Helps users to create posts within a specific community.
It takes various parameters like post's title, content, author,creation timestamp, community Id, post Id.
This function ensures that the post is associated with the relevant community and user, adding the post's Id to their respective lists.
*/
func (s *SmartContract) CreatePost(ctx contractapi.TransactionContextInterface, communityId string, id string, title string, content string, author string, createdAt string) error {
	existingCommunity, err := s.GetCommunity(ctx, communityId)
	if err != nil {
		return err
	}
	if existingCommunity == nil {
		return fmt.Errorf("Community with ID %s doesn't exists", communityId)
	}
	existingPost, err := s.GetPost(ctx, id)
	if err == nil && existingPost != nil {
		return fmt.Errorf("Post with ID %s already exists", id)
	}
	existingUser, err := s.GetUser(ctx, author)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return fmt.Errorf("User with ID %s doesn't exists", author)
	}
	layout := "2006-01-02T15:04:05.000Z"
	currentTime, _ := time.Parse(layout, createdAt)
	post := Post{
		ID:        id,
		Title:     title,
		Content:   content,
		Author:    author,
		CreatedAt: currentTime,
		Score:     0,
		Comments:  make([]string, 0),
		Hidden:    false,
		Community: communityId,
		HideCount: 0,
		ShowCount: 0,
	}
	existingCommunity.Posts = append(existingCommunity.Posts, id)
	existingUser.Posts = append(existingUser.Posts, id)
	communityJson, _ := json.Marshal(existingCommunity)
	ctx.GetStub().PutState(communityId, communityJson)
	userJson, _ := json.Marshal(existingUser)
	ctx.GetStub().PutState(author, userJson)
	postJson, _ := json.Marshal(post)
	ctx.GetStub().PutState(id, postJson)
	return nil
}

/*
Used to retrieve information about a post from blockchain.
It takes the post's Id as a parameter and retrieves the post's data from the blockchain.
*/
func (s *SmartContract) GetPost(ctx contractapi.TransactionContextInterface, id string) (*Post, error) {
	postJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read post from ledger: %w", err)
	}
	if postJson == nil {
		return nil, nil
	}

	var post Post
	err = json.Unmarshal(postJson, &post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

/*
Allows users to upvote a post or comment. It takes post or comment Id and user Id as parameters.
When a user upvotes a post or comment, it increases the item's score.
The function also manages the reputation system by incrementing the author's reputation score if the user is not the author.
*/
func (s *SmartContract) UpVotePost(ctx contractapi.TransactionContextInterface, postId string, userId string) error {
	var author string
	if postId[0] == 'p' {
		existingPost, err := s.GetPost(ctx, postId)
		if err != nil {
			return err
		}
		if existingPost == nil {
			return fmt.Errorf("Post with ID %s doesn't exists", postId)
		}
		existingPost.Score += 1
		postJson, _ := json.Marshal(existingPost)
		author = existingPost.Author
		ctx.GetStub().PutState(postId, postJson)
	} else {
		existingComment, err := s.GetComment(ctx, postId)
		if err != nil {
			return err
		}
		if existingComment == nil {
			return fmt.Errorf("Comment with ID %s doesn't exists", postId)
		}
		existingComment.Score += 1
		commentJson, _ := json.Marshal(existingComment)
		author = existingComment.Author
		ctx.GetStub().PutState(postId, commentJson)
	}
	existingUser, err := s.GetUser(ctx, author)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return fmt.Errorf("User with ID %s doesn't exists", userId)
	}
	if author != userId {
		existingUser.Reputation += 1
	}
	userJson, _ := json.Marshal(existingUser)
	ctx.GetStub().PutState(author, userJson)
	return nil
}

/*
Allows users to downvote a post or comment. It takes post or comment Id and user Id as parameters.
When a user downvotes a post or comment, it decreases the item's score.
The function also manages the reputation system by decreamenting the author's reputation score if the user is not the author.
*/
func (s *SmartContract) DownVotePost(ctx contractapi.TransactionContextInterface, postId string, userId string) error {
	var author string
	if postId[0] == 'p' {
		existingPost, err := s.GetPost(ctx, postId)
		if err != nil {
			return err
		}
		if existingPost == nil {
			return fmt.Errorf("Post with ID %s doesn't exists", postId)
		}
		existingPost.Score -= 1
		postJson, _ := json.Marshal(existingPost)
		author = existingPost.Author
		ctx.GetStub().PutState(postId, postJson)
	} else {
		existingComment, err := s.GetComment(ctx, postId)
		if err != nil {
			return err
		}
		if existingComment == nil {
			return fmt.Errorf("Comment with ID %s doesn't exists", postId)
		}
		existingComment.Score -= 1
		commentJson, _ := json.Marshal(existingComment)
		author = existingComment.Author
		ctx.GetStub().PutState(postId, commentJson)
	}
	existingUser, err := s.GetUser(ctx, author)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return fmt.Errorf("User with ID %s doesn't exists", userId)
	}
	if author != userId {
		existingUser.Reputation -= 1
	}
	userJson, _ := json.Marshal(existingUser)
	ctx.GetStub().PutState(author, userJson)
	return nil
}

/*
Used to retrieve information about a comment from blockchain.
It takes the cpmment's Id as a parameter and retrieves the comment's data from the blockchain.
*/

func (s *SmartContract) GetComment(ctx contractapi.TransactionContextInterface, id string) (*Comment, error) {
	commentJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read community from ledger: %w", err)
	}
	if commentJson == nil {
		return nil, nil
	}

	var comment Comment
	err = json.Unmarshal(commentJson, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

/*
Used to create comments on blockchain. It takes various parameters like  content, author,creation timestamp, comment Id, parent Id.
It ensures that comments are associated with their parent posts or comments, by adding comment id in comments or replies of parent post or comment respectively.
*/
func (s *SmartContract) CreateComment(ctx contractapi.TransactionContextInterface, commentId string, parentId string, content string, author string, createdAt string) error {
	existingComment, err := s.GetComment(ctx, commentId)
	if err != nil {
		return err
	}
	if err == nil && existingComment != nil {
		return fmt.Errorf("Comment with ID %s already exists", commentId)
	}
	var communityId string
	if parentId[0] == 'p' { //If parent is post
		existingPost, err := s.GetPost(ctx, parentId)
		if err != nil {
			return err
		}
		if existingPost == nil {
			return fmt.Errorf("Post with ID %s doesn't exists", parentId)
		}
		existingPost.Comments = append(existingPost.Comments, commentId)
		communityId = existingPost.Community
		postJson, _ := json.Marshal(existingPost)
		ctx.GetStub().PutState(parentId, postJson)
	} else { //If parent is comment
		existingComment, err := s.GetComment(ctx, parentId)
		if err != nil {
			return err
		}
		if existingComment == nil {
			return fmt.Errorf("Comment with ID %s doesn't exists", parentId)
		}
		existingComment.Replies = append(existingComment.Replies, commentId)
		communityId = existingComment.Community
		commentJson, _ := json.Marshal(existingComment)
		ctx.GetStub().PutState(parentId, commentJson)
	}

	layout := "2006-01-02T15:04:05.000Z"
	currentTime, _ := time.Parse(layout, createdAt)
	comment := Comment{
		ID:        commentId,
		Content:   content,
		Author:    author,
		CreatedAt: currentTime,
		Score:     0,
		Parent:    parentId,
		Replies:   make([]string, 0),
		Hidden:    false,
		Community: communityId,
		HideCount: 0,
		ShowCount: 0,
	}
	commentJson, _ := json.Marshal(comment)
	ctx.GetStub().PutState(commentId, commentJson)
	return nil
}

/*
Generates a personalized feed for a user. It takes user Id and page No as parameters.
It compiles posts from the user's joined communities, filtering out any hidden posts based on community moderation, and sorts them by timestamp in reverse chronological order.
Uses pagination for managing large feeds.
*/
func (s *SmartContract) GetUserFeed(ctx contractapi.TransactionContextInterface, userId string, pageNo int) ([]*Post, error) {
	// var userFeed []*Post
	// existingUser, err := s.GetUser(ctx, userId)
	// if err != nil {
	// 	return nil, err
	// }
	// if existingUser == nil {
	// 	return nil, fmt.Errorf("User with ID %s doesn't exists", userId)
	// }

	// for _, community := range existingUser.Communities {
	// 	communityPosts := []*Post{}
	// 	existingCommunity, err := s.GetCommunity(ctx, community)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if existingCommunity == nil {
	// 		return nil, fmt.Errorf("Community with ID %s doesn't exists", community)
	// 	}
	// 	for _, postId := range existingCommunity.Posts {
	// 		post, err := s.GetPost(ctx, postId) // Function to get a post by ID
	// 		if post.Hidden {
	// 			continue
	// 		}
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		communityPosts = append(communityPosts, post)
	// 	}

	// 	// Add community posts to the user's feed
	// 	userFeed = append(userFeed, communityPosts...)
	// }

	// // Sort the user's feed by timestamp (reverse chronological order)
	// sort.Slice(userFeed, func(i, j int) bool {
	// 	return userFeed[i].CreatedAt.After(userFeed[j].CreatedAt)
	// })
	// if len(userFeed) <= PostsPerPage*pageNo {
	// 	return []*Post{}, nil
	// }
	// return userFeed[PostsPerPage*pageNo : min(PostsPerPage*(pageNo+1), len(userFeed))], nil
	var userFeed []*Post
	existingUser, err := s.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, fmt.Errorf("User with ID %s doesn't exist", userId)
	}

	// Create a map to store community posts and channels to wait for them
	communityPosts := make(map[string][]*Post)
	var wg sync.WaitGroup

	for _, community := range existingUser.Communities {
		wg.Add(1)
		go func(communityID string) {
			defer wg.Done()

			existingCommunity, err := s.GetCommunity(ctx, communityID)
			if err != nil {
				// Handle the error
				return
			}

			// Fetch posts in parallel
			communityPosts[communityID] = s.fetchCommunityPosts(ctx, existingCommunity.Posts)
		}(community)
	}

	wg.Wait()

	// Process community posts
	for _, community := range existingUser.Communities {
		userFeed = append(userFeed, communityPosts[community]...)
	}

	// Sort the user's feed by timestamp (reverse chronological order)
	sort.Slice(userFeed, func(i, j int) bool {
		return userFeed[i].CreatedAt.After(userFeed[j].CreatedAt)
	})

	// Calculate the range for pagination
	start := pageNo * PostsPerPage
	end := min((pageNo+1)*PostsPerPage, len(userFeed))

	if start >= len(userFeed) {
		return []*Post{}, nil
	}

	return userFeed[start:end], nil
}

func (s *SmartContract) fetchCommunityPosts(ctx contractapi.TransactionContextInterface, postIDs []string) []*Post {
	communityPosts := []*Post{}
	for _, postID := range postIDs {
		post, err := s.GetPost(ctx, postID)
		if err != nil {
			// Handle the error
			continue
		}
		if !post.Hidden {
			communityPosts = append(communityPosts, post)
		}
	}
	return communityPosts
}

/*
Used fetching a feed of immediate level comments. It takes parent Id and page No as parameters.
It can either replies to a parent comment or comments on a post, in reverse chronological order.
It ensures that hidden comments, as determined by community moderation, are excluded.
Uses pagination for managing large feeds.
*/
func (s *SmartContract) GetCommentFeed(ctx contractapi.TransactionContextInterface, parentId string, pageNo int) ([]*Comment, error) {
	var commentFeed []*Comment
	var commentList []string
	if parentId[0] == 'p' { //If parent is post
		existingPost, err := s.GetPost(ctx, parentId)
		if err != nil {
			return nil, err
		}
		if existingPost == nil {
			return nil, fmt.Errorf("Post with ID %s doesn't exists", parentId)
		}
		commentList = existingPost.Comments
	} else { //If parent is comment
		existingComment, err := s.GetComment(ctx, parentId)
		if err != nil {
			return nil, err
		}
		if existingComment == nil {
			return nil, fmt.Errorf("Comment with ID %s doesn't exists", parentId)
		}
		commentList = existingComment.Replies
	}
	for i := len(commentList) - 1; i >= 0; i-- {
		comment, err := s.GetComment(ctx, commentList[i]) // Function to get a post by ID
		if comment.Hidden {
			continue
		}
		if err != nil {
			return nil, err
		}
		commentFeed = append(commentFeed, comment)
	}
	sort.Slice(commentFeed, func(i, j int) bool {
		return commentFeed[i].CreatedAt.After(commentFeed[j].CreatedAt)
	})
	if len(commentFeed) <= CommentsPerPage*pageNo {
		return []*Comment{}, nil
	}
	return commentFeed[CommentsPerPage*pageNo : min(CommentsPerPage*(pageNo+1), len(commentFeed))], nil

}

/*
Enables users to hide their own posts or comments from public by setting the "Hidden" property to true.
This function performs a verification step to ensure that the user attempting to delete the content is indeed the author, preventing unauthorized deletions.
*/
func (s *SmartContract) DeletePost(ctx contractapi.TransactionContextInterface, postId string, userId string) error {
	if postId[0] == 'p' {
		existingPost, err := s.GetPost(ctx, postId)
		if err != nil {
			return err
		}
		if existingPost == nil {
			return fmt.Errorf("Post with ID %s doesn't exists", postId)
		}
		if userId != existingPost.Author {
			return fmt.Errorf("User cannot delete post with ID %s ", postId)
		}
		existingPost.Hidden = true
		postJson, _ := json.Marshal(existingPost)
		ctx.GetStub().PutState(postId, postJson)
		return nil
	} else {
		existingComment, err := s.GetComment(ctx, postId)
		if err != nil {
			return err
		}
		if existingComment == nil {
			return fmt.Errorf("Comment with ID %s doesn't exists", postId)
		}
		if userId != existingComment.Author {
			return fmt.Errorf("User cannot delete comment with ID %s ", postId)
		}
		existingComment.Hidden = true
		commentJson, _ := json.Marshal(existingComment)
		ctx.GetStub().PutState(postId, commentJson)
		return nil
	}
}

/*
Allows users to appeal a hidden post or comment within a community.
It takes postId or comment Id and user Id as parameters
It adds the post or comment to the list of appealed items in the associated community
*/
func (s *SmartContract) AppealPost(ctx contractapi.TransactionContextInterface, postId string, userId string) error {
	var communityId string
	if postId[0] == 'p' {
		existingPost, err := s.GetPost(ctx, postId)
		if err != nil {
			return err
		}
		if existingPost == nil {
			return fmt.Errorf("Post with ID %s doesn't exists", postId)
		}

		communityId = existingPost.Community

	} else {
		existingComment, err := s.GetComment(ctx, postId)
		if err != nil {
			return err
		}
		if existingComment == nil {
			return fmt.Errorf("Comment with ID %s doesn't exists", postId)
		}
		communityId = existingComment.Community
	}
	existingCommunity, err := s.GetCommunity(ctx, communityId)
	if err != nil {
		return err
	}
	if existingCommunity == nil {
		return fmt.Errorf("Community with ID %s doesn't exists", communityId)
	}
	if !contains(existingCommunity.Users, userId) {
		return fmt.Errorf("User cannot appeal as you are not part of the community")
	}
	existingCommunity.Appealed = append(existingCommunity.Appealed, postId)
	communityJson, _ := json.Marshal(existingCommunity)
	ctx.GetStub().PutState(communityId, communityJson)
	return nil
}

/*
Allows moderators to hide posts or comments that are appealed within a community.
It verifies the moderator status of the user and increments the hide count for the post or comment.
If the hide count reaches a threshold (half of the total moderators in the community), the associated content is marked as hidden
If post is hidden then it is removed from the appeal list and post list of the community.
If comments is hidden then it is also removed from its parent's list of replies.
*/
func (s *SmartContract) HidePostModerator(ctx contractapi.TransactionContextInterface, postId string, userId string) error {
	var communityId string
	if postId[0] == 'p' {
		existingPost, err := s.GetPost(ctx, postId)
		if err != nil {
			return err
		}
		if existingPost == nil {
			return fmt.Errorf("Post with ID %s doesn't exists", postId)
		}
		communityId = existingPost.Community
		existingCommunity, err := s.GetCommunity(ctx, communityId)
		if err != nil {
			return err
		}
		if existingCommunity == nil {
			return fmt.Errorf("Community with ID %s doesn't exists", communityId)
		}
		if !contains(existingCommunity.Moderators, userId) {
			return fmt.Errorf("User cannot hide as you are not a moderator")
		}
		existingPost.HideCount += 1
		if existingPost.HideCount >= int(math.Ceil(float64(len(existingCommunity.Moderators))/2.0)) {
			existingPost.Hidden = true
			existingCommunity.Appealed = removeElement(existingCommunity.Appealed, findIndex(existingCommunity.Appealed, postId))
			existingCommunity.Posts = removeElement(existingCommunity.Posts, findIndex(existingCommunity.Posts, postId))
			communityJson, _ := json.Marshal(existingCommunity)
			ctx.GetStub().PutState(communityId, communityJson)
		}
		postJson, _ := json.Marshal(existingPost)
		ctx.GetStub().PutState(postId, postJson)
	} else {
		existingComment, err := s.GetComment(ctx, postId)
		if err != nil {
			return err
		}
		if existingComment == nil {
			return fmt.Errorf("Comment with ID %s doesn't exists", postId)
		}
		communityId = existingComment.Community
		existingCommunity, err := s.GetCommunity(ctx, communityId)
		if err != nil {
			return err
		}
		if existingCommunity == nil {
			return fmt.Errorf("Community with ID %s doesn't exists", communityId)
		}
		if !contains(existingCommunity.Moderators, userId) {
			return fmt.Errorf("User cannot hide as you are not a moderator")
		}
		existingComment.HideCount += 1
		if existingComment.HideCount >= int(math.Ceil(float64(len(existingCommunity.Moderators))/2.0)) {
			existingComment.Hidden = true
			parentId := existingComment.Parent
			if parentId[0] == 'p' {
				existingParent, _ := s.GetPost(ctx, parentId)
				existingParent.Comments = removeElement(existingParent.Comments, findIndex(existingParent.Comments, postId))
				parentJson, _ := json.Marshal(existingParent)
				ctx.GetStub().PutState(parentId, parentJson)
			} else {
				existingParent, _ := s.GetComment(ctx, parentId)
				existingParent.Replies = removeElement(existingParent.Replies, findIndex(existingParent.Replies, postId))
				parentJson, _ := json.Marshal(existingParent)
				ctx.GetStub().PutState(parentId, parentJson)
			}

			existingCommunity.Appealed = removeElement(existingCommunity.Appealed, findIndex(existingCommunity.Appealed, postId))
			communityJson, _ := json.Marshal(existingCommunity)
			ctx.GetStub().PutState(communityId, communityJson)
		}
		commentJson, _ := json.Marshal(existingComment)
		ctx.GetStub().PutState(postId, commentJson)
	}
	return nil
}

/*
Helps in selecting moderators for a given community based on the reputation.
The function  evaluates user activity by summing the scores of their posts.
Users are ranked by their aggregated scores, and the top users, up to the required number of moderators, are chosen as new moderators for the community.
*/
func (s *SmartContract) SelectModerator(ctx contractapi.TransactionContextInterface, communityId string) error {
	existingCommunity, err := s.GetCommunity(ctx, communityId)
	if err != nil {
		return err
	}
	if existingCommunity == nil {
		return fmt.Errorf("Community with ID %s doesn't exists", communityId)
	}
	sizeCommunity := len(existingCommunity.Users)
	noOfModeratorsRequired := int(math.Ceil(float64(sizeCommunity) * 0.1))
	noOfModeratorsRequired = max(noOfModeratorsRequired, 1)
	noOfModeratorsRequired = min(noOfModeratorsRequired, 100)
	userMap := make(map[string]int)
	for _, postId := range existingCommunity.Posts {
		existingPost, err := s.GetPost(ctx, postId)
		if err != nil {
			return err
		}
		if existingPost == nil {
			return fmt.Errorf("Post with ID %s doesn't exists", postId)
		}
		val, exists := userMap[existingPost.Author]
		if exists {
			userMap[existingPost.Author] = val + existingPost.Score
		} else {
			userMap[existingPost.Author] = existingPost.Score
		}
	}
	keys := make([]string, 0, len(userMap))

	for key := range userMap {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return userMap[keys[i]] > userMap[keys[j]]
	})

	var newModerators []string
	for _, k := range keys {
		newModerators = append(newModerators, k)
		if len(newModerators) == noOfModeratorsRequired {
			break
		}
	}
	existingCommunity.Moderators = newModerators
	communityJson, _ := json.Marshal(existingCommunity)
	ctx.GetStub().PutState(communityId, communityJson)
	return nil
}

/*
Allow users to retract their appeal on a post or comment.
It checks the user's membership in the community.
The function removes the user's appeal from the list of appeals for the specific post or comment.
*/
func (s *SmartContract) UnAppealPost(ctx contractapi.TransactionContextInterface, postId string, userId string) error {
	var communityId string
	if postId[0] == 'p' {
		existingPost, err := s.GetPost(ctx, postId)
		if err != nil {
			return err
		}
		if existingPost == nil {
			return fmt.Errorf("Post with ID %s doesn't exists", postId)
		}

		communityId = existingPost.Community

	} else {
		existingComment, err := s.GetComment(ctx, postId)
		if err != nil {
			return err
		}
		if existingComment == nil {
			return fmt.Errorf("Comment with ID %s doesn't exists", postId)
		}
		communityId = existingComment.Community
	}
	existingCommunity, err := s.GetCommunity(ctx, communityId)
	if err != nil {
		return err
	}
	if existingCommunity == nil {
		return fmt.Errorf("Community with ID %s doesn't exists", communityId)
	}
	if !contains(existingCommunity.Users, userId) {
		return fmt.Errorf("User cannot appeal as you are not part of the community")
	}
	if !contains(existingCommunity.Appealed, postId) {
		return fmt.Errorf("User cannot unappeal the post")
	}
	existingCommunity.Appealed = removeElement(existingCommunity.Appealed, findIndex(existingCommunity.Appealed, postId))
	communityJson, _ := json.Marshal(existingCommunity)
	ctx.GetStub().PutState(communityId, communityJson)
	return nil
}

/*
Allows moderators to uhide/Show posts or comments that are appealed within a community.
It verifies the moderator status of the user and increments the show count for the post or comment.
If the show count reaches a threshold (half of the total moderators in the community), the associated content is removed from the appealed list of that community.
*/
func (s *SmartContract) ShowPostModerator(ctx contractapi.TransactionContextInterface, postId string, userId string) error {
	var communityId string
	if postId[0] == 'p' {
		existingPost, err := s.GetPost(ctx, postId)
		if err != nil {
			return err
		}
		if existingPost == nil {
			return fmt.Errorf("Post with ID %s doesn't exists", postId)
		}
		communityId = existingPost.Community
		existingCommunity, err := s.GetCommunity(ctx, communityId)
		if err != nil {
			return err
		}
		if existingCommunity == nil {
			return fmt.Errorf("Community with ID %s doesn't exists", communityId)
		}
		if !contains(existingCommunity.Moderators, userId) {
			return fmt.Errorf("User cannot hide as you are not a moderator")
		}
		existingPost.ShowCount += 1
		if existingPost.ShowCount >= int(math.Ceil(float64(len(existingCommunity.Moderators))/2.0)) {
			//existingPost.Hidden = true
			existingCommunity.Appealed = removeElement(existingCommunity.Appealed, findIndex(existingCommunity.Appealed, postId))
			//existingCommunity.Posts = removeElement(existingCommunity.Posts, findIndex(existingCommunity.Posts, postId))
			communityJson, _ := json.Marshal(existingCommunity)
			ctx.GetStub().PutState(communityId, communityJson)
		}
		postJson, _ := json.Marshal(existingPost)
		ctx.GetStub().PutState(postId, postJson)
	} else {
		existingComment, err := s.GetComment(ctx, postId)
		if err != nil {
			return err
		}
		if existingComment == nil {
			return fmt.Errorf("Comment with ID %s doesn't exists", postId)
		}
		communityId = existingComment.Community
		existingCommunity, err := s.GetCommunity(ctx, communityId)
		if err != nil {
			return err
		}
		if existingCommunity == nil {
			return fmt.Errorf("Community with ID %s doesn't exists", communityId)
		}
		if !contains(existingCommunity.Moderators, userId) {
			return fmt.Errorf("User cannot hide as you are not a moderator")
		}
		existingComment.ShowCount += 1
		if existingComment.ShowCount >= int(math.Ceil(float64(len(existingCommunity.Moderators))/2.0)) {
			//existingComment.Hidden = true
			// parentId := existingComment.Parent
			// if parentId[0] == 'p' {
			// 	existingParent, _ := s.GetPost(ctx, parentId)
			// 	existingParent.Comments = removeElement(existingParent.Comments, findIndex(existingParent.Comments, postId))
			// 	parentJson, _ := json.Marshal(existingParent)
			// 	ctx.GetStub().PutState(parentId, parentJson)
			// } else {
			// 	existingParent, _ := s.GetComment(ctx, parentId)
			// 	existingParent.Replies = removeElement(existingParent.Replies, findIndex(existingParent.Replies, postId))
			// 	parentJson, _ := json.Marshal(existingParent)
			// 	ctx.GetStub().PutState(parentId, parentJson)
			// }
			existingCommunity.Appealed = removeElement(existingCommunity.Appealed, findIndex(existingCommunity.Appealed, postId))
			communityJson, _ := json.Marshal(existingCommunity)
			ctx.GetStub().PutState(communityId, communityJson)
		}
		commentJson, _ := json.Marshal(existingComment)
		ctx.GetStub().PutState(postId, commentJson)
	}
	return nil
}

//unappeal undo done
//if removed by moderator then remove from posts of community and also when deleted done
//Moderator selection should be dependent on score from posts and commentss not necessary
//non hiding decision made by the moderator done
//hidden can be int 0-> approved 1->normal 2->hidden
//go coroutines to parallelize the userfeed done
