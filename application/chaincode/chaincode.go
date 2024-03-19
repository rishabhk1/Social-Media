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

type MetaData struct {
	ID   string
	Name []CommunityName `json:"name"`
}

type User struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Communities []string `json:"communities"`
	Posts       []string `json:"posts"` //list of ids
	Comments    []string `json:"comments"`
	Reputation  int
}

type UserModified struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Reputation int    `json:"reputation"`
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

type CommunityModified struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Creator     string         `json:"creator"`
	CreatedAt   time.Time      `json:"createdAt"`
	Moderators  []UserModified `json:"moderators"`
	Users       []UserModified `json:"users"`
}

type CommunityName struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
	UpVote    []string
	DownVote  []string
	HideVote  []string
	ShowVote  []string
}

type PostModified struct {
	ID            string    `json:"id"` // should start with p
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Author        string    `json:"author"`
	Score         int       `json:"score"`
	CreatedAt     time.Time `json:"createdAt"`
	Comments      []string  `json:"comments"` //list of ids
	Hidden        bool      `json:"hidden"`
	Community     string    `json:"community"`
	HideCount     int       `json:"hideCount"`
	ShowCount     int       `json:"showCount"`
	AuthorName    string    `json:"authorName"`
	CommunityName string    `json:"communityName"`
	HasUpvoted    bool      `json:"hasUpvoted"`
	HasDownvoted  bool      `json:"hasDownvoted"`
	HasHideVoted  bool      `json:"hasHidevoted"`
	HasShowVoted  bool      `json:"hasShowvoted"`
	IsAppealed    bool      `json:"isAppealed"`
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
	UpVote    []string
	DownVote  []string
	HideVote  []string
	ShowVote  []string
}

type CommentModified struct {
	ID            string    `json:"id"` //should syart with c
	Content       string    `json:"content"`
	Author        string    `json:"author"`
	Score         int       `json:"score"`
	CreatedAt     time.Time `json:"createdAt"`
	Parent        string    `json:"parentId"`
	Replies       []string  `json:"replies"` //list of ids
	Hidden        bool      `json:"hidden"`
	Community     string    `json:"community"`
	HideCount     int       `json:"hideCount"`
	ShowCount     int       `json:"showCount"`
	AuthorName    string    `json:"authorName"`
	CommunityName string    `json:"communityName"`
	HasUpvoted    bool      `json:"hasUpvoted"`
	HasDownvoted  bool      `json:"hasDownvoted"`
	IsAppealed    bool      `json:"isAppealed"`
	HasHideVoted  bool      `json:"hasHidevoted"`
	HasShowVoted  bool      `json:"hasShowvoted"`
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
		Comments:    []string{"c_1", "c_2"},
	}

	user2 := User{
		ID:          "2",
		Username:    "joe_rog",
		Email:       "joe_rog@example.com",
		Communities: []string{"co_1"},
		Posts:       []string{"p_2"},
		Comments:    make([]string, 0),
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
		Appealed:    []string{"p_2"},
	}

	ctime, _ = time.Parse(layout, "2023-04-17T15:04:05.000Z")
	post1 := Post{
		ID:        "p_1",
		Title:     "Demo Post",
		Content:   "Demo post content",
		Author:    "1",
		Score:     10,
		CreatedAt: ctime,
		Comments:  []string{"c_1"},
		Hidden:    false,
		Community: "co_1",
		HideCount: 0,
		ShowCount: 0,
		UpVote:    make([]string, 0),
		DownVote:  make([]string, 0),
		HideVote:  make([]string, 0),
		ShowVote:  make([]string, 0),
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
		UpVote:    make([]string, 0),
		DownVote:  make([]string, 0),
		HideVote:  make([]string, 0),
		ShowVote:  make([]string, 0),
	}

	comment1 := Comment{
		ID:        "c_1",
		Content:   "comment of p_1",
		Author:    "1",
		Score:     0,
		CreatedAt: ctime,
		Parent:    "p_1",
		Replies:   []string{"c_2"},
		Hidden:    false,
		Community: "co_1",
		HideCount: 0,
		ShowCount: 0,
		UpVote:    make([]string, 0),
		DownVote:  make([]string, 0),
		HideVote:  make([]string, 0),
		ShowVote:  make([]string, 0),
	}
	comment2 := Comment{
		ID:        "c_2",
		Content:   "comment of c_1",
		Author:    "1",
		Score:     0,
		CreatedAt: ctime,
		Parent:    "c_1",
		Replies:   make([]string, 0),
		Hidden:    false,
		Community: "co_1",
		HideCount: 0,
		ShowCount: 0,
		UpVote:    make([]string, 0),
		DownVote:  make([]string, 0),
		HideVote:  make([]string, 0),
		ShowVote:  make([]string, 0),
	}
	communityName1 := CommunityName{
		ID:   "co_1",
		Name: "Comm1",
	}
	metaData := MetaData{
		ID:   "md",
		Name: []CommunityName{communityName1},
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
	comment1Json, err := json.Marshal(comment1)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(comment1.ID, comment1Json)
	if err != nil {
		return err
	}

	comment2Json, err := json.Marshal(comment2)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(comment2.ID, comment2Json)
	if err != nil {
		return err
	}
	metaDataJson, err := json.Marshal(metaData)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(metaData.ID, metaDataJson)
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
		Comments:    make([]string, 0),
	}
	userJson, _ := json.Marshal(user)
	return ctx.GetStub().PutState(UserId, userJson)

}

func (s *SmartContract) GetMetaData(ctx contractapi.TransactionContextInterface, metaDataId string) (*MetaData, error) {
	metaJson, err := ctx.GetStub().GetState(metaDataId)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from ledger: %w", err)
	}
	if metaJson == nil {
		return nil, nil
	}

	var metaData MetaData
	err = json.Unmarshal(metaJson, &metaData)
	if err != nil {
		return nil, err
	}
	return &metaData, nil
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
func (s *SmartContract) GetUserModified(ctx contractapi.TransactionContextInterface, userId string) (*UserModified, error) {
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
	var userModified *UserModified
	userModified, err = s.convertToUserModified(ctx, &user)
	if err != nil {
		return nil, err
	}
	return userModified, nil
}

func (s *SmartContract) GetCommunityModified(ctx contractapi.TransactionContextInterface, communityId string) (*CommunityModified, error) {
	communityJson, err := ctx.GetStub().GetState(communityId)
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
	var communityModified *CommunityModified
	communityModified, err = s.convertToCommunityModified(ctx, &community)
	if err != nil {
		return nil, err
	}
	return communityModified, nil
}

/*
Used to create a new community in a blockchain.
It takes parameters such as the community's Id, name, description, creator, and creation timestamp.
Function also makes the creator, the initial moderator of the community.
*/
func (s *SmartContract) CreateCommunity(ctx contractapi.TransactionContextInterface, id string, createdAt string, name string, description string, creator string) error {
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
	existingMetaData, err := s.GetMetaData(ctx, "md")
	if err != nil {
		return err
	}
	if existingMetaData == nil {
		return fmt.Errorf("data  doesn't exists")
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
	communityName := CommunityName{
		ID:   id,
		Name: name,
	}
	existingMetaData.Name = append(existingMetaData.Name, communityName)
	metaDataJson, _ := json.Marshal(existingMetaData)
	ctx.GetStub().PutState("1", metaDataJson)
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
func (s *SmartContract) JoinCommunity(ctx contractapi.TransactionContextInterface, communityId string, userId string) (*UserModified, error) {
	existingCommunity, err := s.GetCommunity(ctx, communityId)
	if err != nil {
		return nil, err
	}
	if existingCommunity == nil {
		return nil, fmt.Errorf("Community with ID %s doesn't exists", communityId)
	}
	existingCommunity.Users = append(existingCommunity.Users, userId)
	communityJson, _ := json.Marshal(existingCommunity)
	ctx.GetStub().PutState(communityId, communityJson)
	user, err := s.GetUserModified(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

/*
Used to allow a user to leave a specific community.
It takes the user's Id and the community's Id as parameters and removes the user to the list of community members
*/
func (s *SmartContract) UnJoinCommunity(ctx contractapi.TransactionContextInterface, communityId string, userId string) (bool, error) {
	existingCommunity, err := s.GetCommunity(ctx, communityId)
	if err != nil {
		return false, err
	}
	if existingCommunity == nil {
		return false, fmt.Errorf("Community with ID %s doesn't exists", communityId)
	}
	existingCommunity.Users = removeElement(existingCommunity.Users, findIndex(existingCommunity.Users, userId))
	communityJson, _ := json.Marshal(existingCommunity)
	ctx.GetStub().PutState(communityId, communityJson)
	return true, nil
}

/*
Helps users to create posts within a specific community.
It takes various parameters like post's title, content, author,creation timestamp, community Id, post Id.
This function ensures that the post is associated with the relevant community and user, adding the post's Id to their respective lists.
*/
func (s *SmartContract) CreatePost(ctx contractapi.TransactionContextInterface, id string, createdAt string, communityId string, title string, content string, author string) error {
	existingCommunity, err := s.GetCommunity(ctx, communityId)
	fmt.Println(existingCommunity)
	fmt.Println(err)
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
		UpVote:    make([]string, 0),
		DownVote:  make([]string, 0),
		HideVote:  make([]string, 0),
		ShowVote:  make([]string, 0),
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

func (s *SmartContract) GetPostModified(ctx contractapi.TransactionContextInterface, postId string, userId string) (*PostModified, error) {
	postJson, err := ctx.GetStub().GetState(postId)
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
	var postModified *PostModified
	postModified, err = s.convertToPostModified(ctx, &post, userId)
	if err != nil {
		return nil, err
	}
	return postModified, nil
}

/*
Allows users to upvote a post or comment. It takes post or comment Id and user Id as parameters.
When a user upvotes a post or comment, it increases the item's score.
The function also manages the reputation system by incrementing the author's reputation score if the user is not the author.
*/
func (s *SmartContract) UpVotePost(ctx contractapi.TransactionContextInterface, postId string, userId string) (bool, error) {
	var author string
	var upVotedDiff = 0
	if postId[0] == 'p' {
		existingPost, err := s.GetPost(ctx, postId)
		if err != nil {
			return false, err
		}
		if existingPost == nil {
			return false, fmt.Errorf("Post with ID %s doesn't exists", postId)
		}
		if !contains(existingPost.UpVote, userId) {
			existingPost.Score += 1
			upVotedDiff += 1
			existingPost.UpVote = append(existingPost.UpVote, userId)
			if contains(existingPost.DownVote, userId) {
				existingPost.DownVote = removeElement(existingPost.DownVote, findIndex(existingPost.DownVote, userId))
				existingPost.Score += 1
				upVotedDiff += 1
			}
		}
		postJson, _ := json.Marshal(existingPost)
		author = existingPost.Author
		ctx.GetStub().PutState(postId, postJson)
	} else {
		existingComment, err := s.GetComment(ctx, postId)
		if err != nil {
			return false, err
		}
		if existingComment == nil {
			return false, fmt.Errorf("Comment with ID %s doesn't exists", postId)
		}
		if !contains(existingComment.UpVote, userId) {
			existingComment.Score += 1
			upVotedDiff += 1
			existingComment.UpVote = append(existingComment.UpVote, userId)
			if contains(existingComment.DownVote, userId) {
				existingComment.DownVote = removeElement(existingComment.DownVote, findIndex(existingComment.DownVote, userId))
				existingComment.Score += 1
				upVotedDiff += 1
			}
		}
		commentJson, _ := json.Marshal(existingComment)
		author = existingComment.Author
		ctx.GetStub().PutState(postId, commentJson)
	}
	existingUser, err := s.GetUser(ctx, author)
	if err != nil {
		return false, err
	}
	if existingUser == nil {
		return false, fmt.Errorf("User with ID %s doesn't exists", userId)
	}
	if author != userId {
		existingUser.Reputation += upVotedDiff
	}
	userJson, _ := json.Marshal(existingUser)
	ctx.GetStub().PutState(author, userJson)
	return upVotedDiff > 0, nil
}

func (s *SmartContract) UndoUpVotePost(ctx contractapi.TransactionContextInterface, postId string, userId string) (bool, error) {
	var author string
	var upVotedDiff = 0
	if postId[0] == 'p' {
		existingPost, err := s.GetPost(ctx, postId)
		if err != nil {
			return false, err
		}
		if existingPost == nil {
			return false, fmt.Errorf("Post with ID %s doesn't exists", postId)
		}
		if contains(existingPost.UpVote, userId) {
			existingPost.Score -= 1
			upVotedDiff -= 1
			existingPost.UpVote = removeElement(existingPost.UpVote, findIndex(existingPost.UpVote, userId))
		}
		postJson, _ := json.Marshal(existingPost)
		author = existingPost.Author
		ctx.GetStub().PutState(postId, postJson)
	} else {
		existingComment, err := s.GetComment(ctx, postId)
		if err != nil {
			return false, err
		}
		if existingComment == nil {
			return false, fmt.Errorf("Comment with ID %s doesn't exists", postId)
		}
		if contains(existingComment.UpVote, userId) {
			existingComment.Score -= 1
			upVotedDiff -= 1
			existingComment.UpVote = removeElement(existingComment.UpVote, findIndex(existingComment.UpVote, userId))
		}
		commentJson, _ := json.Marshal(existingComment)
		author = existingComment.Author
		ctx.GetStub().PutState(postId, commentJson)
	}
	existingUser, err := s.GetUser(ctx, author)
	if err != nil {
		return false, err
	}
	if existingUser == nil {
		return false, fmt.Errorf("User with ID %s doesn't exists", userId)
	}
	if author != userId {
		existingUser.Reputation += upVotedDiff
	}
	userJson, _ := json.Marshal(existingUser)
	ctx.GetStub().PutState(author, userJson)
	return upVotedDiff < 0, nil
}

/*
Allows users to downvote a post or comment. It takes post or comment Id and user Id as parameters.
When a user downvotes a post or comment, it decreases the item's score.
The function also manages the reputation system by decreamenting the author's reputation score if the user is not the author.
*/
func (s *SmartContract) DownVotePost(ctx contractapi.TransactionContextInterface, postId string, userId string) (bool, error) {
	var author string
	var downVotedDiff = 0
	if postId[0] == 'p' {
		existingPost, err := s.GetPost(ctx, postId)
		if err != nil {
			return false, err
		}
		if existingPost == nil {
			return false, fmt.Errorf("Post with ID %s doesn't exists", postId)
		}
		if !contains(existingPost.DownVote, userId) {
			existingPost.Score -= 1
			downVotedDiff -= 1
			existingPost.DownVote = append(existingPost.DownVote, userId)
			if contains(existingPost.UpVote, userId) {
				existingPost.UpVote = removeElement(existingPost.UpVote, findIndex(existingPost.UpVote, userId))
				existingPost.Score -= 1
				downVotedDiff -= 1
			}
		}
		// existingPost.Score -= 1
		postJson, _ := json.Marshal(existingPost)
		author = existingPost.Author
		ctx.GetStub().PutState(postId, postJson)
	} else {
		existingComment, err := s.GetComment(ctx, postId)
		if err != nil {
			return false, err
		}
		if existingComment == nil {
			return false, fmt.Errorf("Comment with ID %s doesn't exists", postId)
		}
		if !contains(existingComment.DownVote, userId) {
			existingComment.Score -= 1
			downVotedDiff -= 1
			existingComment.DownVote = append(existingComment.DownVote, userId)
			if contains(existingComment.UpVote, userId) {
				existingComment.UpVote = removeElement(existingComment.UpVote, findIndex(existingComment.UpVote, userId))
				existingComment.Score -= 1
				downVotedDiff -= 1
			}
		}
		commentJson, _ := json.Marshal(existingComment)
		author = existingComment.Author
		ctx.GetStub().PutState(postId, commentJson)
	}
	existingUser, err := s.GetUser(ctx, author)
	if err != nil {
		return false, err
	}
	if existingUser == nil {
		return false, fmt.Errorf("User with ID %s doesn't exists", userId)
	}
	if author != userId {
		existingUser.Reputation += downVotedDiff
	}
	userJson, _ := json.Marshal(existingUser)
	ctx.GetStub().PutState(author, userJson)
	return downVotedDiff < 0, nil
}

func (s *SmartContract) UndoDownVotePost(ctx contractapi.TransactionContextInterface, postId string, userId string) (bool, error) {
	var author string
	var downVotedDiff = 0
	if postId[0] == 'p' {
		existingPost, err := s.GetPost(ctx, postId)
		if err != nil {
			return false, err
		}
		if existingPost == nil {
			return false, fmt.Errorf("Post with ID %s doesn't exists", postId)
		}
		if contains(existingPost.DownVote, userId) {
			existingPost.Score += 1
			downVotedDiff += 1
			existingPost.DownVote = removeElement(existingPost.DownVote, findIndex(existingPost.DownVote, userId))
		}
		// existingPost.Score -= 1
		postJson, _ := json.Marshal(existingPost)
		author = existingPost.Author
		ctx.GetStub().PutState(postId, postJson)
	} else {
		existingComment, err := s.GetComment(ctx, postId)
		if err != nil {
			return false, err
		}
		if existingComment == nil {
			return false, fmt.Errorf("Comment with ID %s doesn't exists", postId)
		}
		if contains(existingComment.DownVote, userId) {
			existingComment.Score += 1
			downVotedDiff += 1
			existingComment.DownVote = removeElement(existingComment.DownVote, findIndex(existingComment.DownVote, userId))
		}
		commentJson, _ := json.Marshal(existingComment)
		author = existingComment.Author
		ctx.GetStub().PutState(postId, commentJson)
	}
	existingUser, err := s.GetUser(ctx, author)
	if err != nil {
		return false, err
	}
	if existingUser == nil {
		return false, fmt.Errorf("User with ID %s doesn't exists", userId)
	}
	if author != userId {
		existingUser.Reputation += downVotedDiff
	}
	userJson, _ := json.Marshal(existingUser)
	ctx.GetStub().PutState(author, userJson)
	return downVotedDiff > 0, nil
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

func (s *SmartContract) GetCommentModified(ctx contractapi.TransactionContextInterface, commentId string, userId string) (*CommentModified, error) {
	commentJson, err := ctx.GetStub().GetState(commentId)
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
	var commentModified *CommentModified
	commentModified, err = s.convertToCommentModified(ctx, &comment, userId)
	if err != nil {
		return nil, err
	}
	return commentModified, nil
}

/*
Used to create comments on blockchain. It takes various parameters like  content, author,creation timestamp, comment Id, parent Id.
It ensures that comments are associated with their parent posts or comments, by adding comment id in comments or replies of parent post or comment respectively.
*/
func (s *SmartContract) CreateComment(ctx contractapi.TransactionContextInterface, commentId string, createdAt string, parentId string, content string, author string) error {
	existingComment, err := s.GetComment(ctx, commentId)
	if err != nil {
		return err
	}
	if err == nil && existingComment != nil {
		return fmt.Errorf("Comment with ID %s already exists", commentId)
	}
	existingUser, err := s.GetUser(ctx, author)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return fmt.Errorf("User with ID %s doesn't exists", author)
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
		UpVote:    make([]string, 0),
		DownVote:  make([]string, 0),
		HideVote:  make([]string, 0),
		ShowVote:  make([]string, 0),
	}
	existingUser.Comments = append(existingUser.Comments, commentId)
	commentJson, _ := json.Marshal(comment)
	ctx.GetStub().PutState(commentId, commentJson)
	userJson, _ := json.Marshal(existingUser)
	ctx.GetStub().PutState(author, userJson)
	return nil
}

/*
Generates a personalized feed for a user. It takes user Id and page No as parameters.
It compiles posts from the user's joined communities, filtering out any hidden posts based on community moderation, and sorts them by timestamp in reverse chronological order.
Uses pagination for managing large feeds.
*/
func (s *SmartContract) convertToPostModified(ctx contractapi.TransactionContextInterface, original *Post, userId string) (*PostModified, error) {

	// Create a new PostModified instance
	existingAuthor, err := s.GetUser(ctx, original.Author)
	if err != nil {
		return nil, err
	}
	if existingAuthor == nil {
		return nil, fmt.Errorf("User with ID %s doesn't exists", original.Author)
	}
	existingCommunity, err := s.GetCommunity(ctx, original.Community)
	if err != nil {
		return nil, err
	}
	if existingAuthor == nil {
		return nil, fmt.Errorf("Community with ID %s doesn't exists", original.Community)
	}
	val, err := s.isAppealed(ctx, original.ID)
	if err != nil {
		return nil, err
	}
	modified := PostModified{
		ID:            original.ID,
		Title:         original.Title,
		Content:       original.Content,
		Author:        original.Author,
		Score:         original.Score,
		CreatedAt:     original.CreatedAt,
		Comments:      original.Comments,
		Hidden:        original.Hidden,
		Community:     original.Community,
		HideCount:     original.HideCount,
		ShowCount:     original.ShowCount,
		AuthorName:    existingAuthor.Username, // Set your desired value for AuthorName
		CommunityName: existingCommunity.Name,  // Set your desired value for CommunityName
		HasUpvoted:    contains(original.UpVote, userId),
		HasDownvoted:  contains(original.DownVote, userId),
		IsAppealed:    val,
		HasHideVoted:  contains(original.HideVote, userId),
		HasShowVoted:  contains(original.ShowVote, userId),
	}
	fmt.Println(original)
	return &modified, nil
}

func (s *SmartContract) convertToCommentModified(ctx contractapi.TransactionContextInterface, original *Comment, userId string) (*CommentModified, error) {

	// Create a new PostModified instance
	existingAuthor, err := s.GetUser(ctx, original.Author)
	if err != nil {
		return nil, err
	}
	if existingAuthor == nil {
		return nil, fmt.Errorf("User with ID %s doesn't exists", original.Author)
	}
	existingCommunity, err := s.GetCommunity(ctx, original.Community)
	if err != nil {
		return nil, err
	}
	if existingCommunity == nil {
		return nil, fmt.Errorf("Community with ID %s doesn't exists", original.Community)
	}
	val, err := s.isAppealed(ctx, original.ID)
	if err != nil {
		return nil, err
	}
	modified := CommentModified{
		ID:            original.ID,
		Content:       original.Content,
		Author:        original.Author,
		Score:         original.Score,
		CreatedAt:     original.CreatedAt,
		Replies:       original.Replies,
		Hidden:        original.Hidden,
		Community:     original.Community,
		HideCount:     original.HideCount,
		ShowCount:     original.ShowCount,
		AuthorName:    existingAuthor.Username, // Set your desired value for AuthorName
		CommunityName: existingCommunity.Name,  // Set your desired value for CommunityName
		HasUpvoted:    contains(original.UpVote, userId),
		HasDownvoted:  contains(original.DownVote, userId),
		Parent:        original.Parent,
		IsAppealed:    val,
		HasHideVoted:  contains(original.HideVote, userId),
		HasShowVoted:  contains(original.ShowVote, userId),
	}
	fmt.Println(original)
	return &modified, nil
}

func (s *SmartContract) convertToUserModified(ctx contractapi.TransactionContextInterface, original *User) (*UserModified, error) {

	// Create a new PostModified instance
	// existingAuthor, err := s.GetUser(ctx, original.Author)
	// if err != nil {
	// 	return nil, err
	// }
	// if existingAuthor == nil {
	// 	return nil, fmt.Errorf("User with ID %s doesn't exists", original.Author)
	// }
	// existingCommunity, err := s.GetCommunity(ctx, original.Community)
	// if err != nil {
	// 	return nil, err
	// }
	// if existingAuthor == nil {
	// 	return nil, fmt.Errorf("Community with ID %s doesn't exists", original.Community)
	// }
	modified := UserModified{
		ID:         original.ID,
		Reputation: original.Reputation,
		Email:      original.Email,
		Username:   original.Username,
	}
	//fmt.Println(original)
	return &modified, nil
}

func (s *SmartContract) convertToCommunityModified(ctx contractapi.TransactionContextInterface, original *Community) (*CommunityModified, error) {

	// Create a new PostModified instance
	// existingAuthor, err := s.GetUser(ctx, original.Author)
	// if err != nil {
	// 	return nil, err
	// }
	// if existingAuthor == nil {
	// 	return nil, fmt.Errorf("User with ID %s doesn't exists", original.Author)
	// }
	// existingCommunity, err := s.GetCommunity(ctx, original.Community)
	// if err != nil {
	// 	return nil, err
	// }
	// if existingAuthor == nil {
	// 	return nil, fmt.Errorf("Community with ID %s doesn't exists", original.Community)
	// }

	var Moderators []UserModified
	var Users []UserModified
	for id := range original.Users {
		// var userModified *UserModified
		userModified, err := s.GetUserModified(ctx, original.Users[id])
		if err != nil {
			return nil, err
		}
		Users = append(Users, *userModified)
	}
	for id := range original.Moderators {
		// var userModified *UserModified
		userModified, err := s.GetUserModified(ctx, original.Moderators[id])
		if err != nil {
			return nil, err
		}
		Moderators = append(Moderators, *userModified)
	}
	modified := CommunityModified{
		ID:          original.ID,
		Name:        original.Name,
		Description: original.Description,
		Creator:     original.Creator,
		CreatedAt:   original.CreatedAt,
		Moderators:  Moderators,
		Users:       Users,
	}
	//fmt.Println(original)
	return &modified, nil
}

func (s *SmartContract) GetUserFeed(ctx contractapi.TransactionContextInterface, userId string, pageNo int) ([]*PostModified, error) {
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
	var userFeedModified []*PostModified
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
		return []*PostModified{}, nil
	}

	for _, originalPost := range userFeed[start:end] {
		modifiedPost, error := s.convertToPostModified(ctx, originalPost, userId)
		if error != nil {
			return nil, error
		}
		userFeedModified = append(userFeedModified, modifiedPost)
	}
	return userFeedModified, nil
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
func (s *SmartContract) GetCommentFeed(ctx contractapi.TransactionContextInterface, parentId string, pageNo int, userId string) ([]*CommentModified, error) {
	var commentFeed []*Comment
	var commentFeedModified []*CommentModified
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
		return []*CommentModified{}, nil
	}

	start := CommentsPerPage * pageNo
	end := min(CommentsPerPage*(pageNo+1), len(commentFeed))
	for _, originalComment := range commentFeed[start:end] {
		modifiedComment, error := s.convertToCommentModified(ctx, originalComment, userId)
		if error != nil {
			return nil, error
		}
		commentFeedModified = append(commentFeedModified, modifiedComment)
	}
	return commentFeedModified, nil
	// return commentFeed[CommentsPerPage*pageNo : min(CommentsPerPage*(pageNo+1), len(commentFeed))], nil

}

func (s *SmartContract) GetUserProfileComments(ctx contractapi.TransactionContextInterface, targetUserId string, userId string, pageNo int) ([]*CommentModified, error) {
	var commentFeed []*Comment
	var commentFeedModified []*CommentModified
	var commentList []string
	// if parentId[0] == 'p' { //If parent is post
	// 	existingPost, err := s.GetPost(ctx, parentId)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if existingPost == nil {
	// 		return nil, fmt.Errorf("Post with ID %s doesn't exists", parentId)
	// 	}
	// 	commentList = existingPost.Comments
	// } else { //If parent is comment
	// 	existingComment, err := s.GetComment(ctx, parentId)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if existingComment == nil {
	// 		return nil, fmt.Errorf("Comment with ID %s doesn't exists", parentId)
	// 	}
	// 	commentList = existingComment.Replies
	// }
	targetUser, err := s.GetUser(ctx, targetUserId)
	commentList = targetUser.Comments
	if err != nil {
		return nil, err
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
	// sort.Slice(commentFeed, func(i, j int) bool {
	// 	return commentFeed[i].CreatedAt.After(commentFeed[j].CreatedAt)
	// })
	if len(commentFeed) <= CommentsPerPage*pageNo {
		return []*CommentModified{}, nil
	}

	start := CommentsPerPage * pageNo
	end := min(CommentsPerPage*(pageNo+1), len(commentFeed))
	for _, originalComment := range commentFeed[start:end] {
		modifiedComment, error := s.convertToCommentModified(ctx, originalComment, userId)
		if error != nil {
			return nil, error
		}
		commentFeedModified = append(commentFeedModified, modifiedComment)
	}
	return commentFeedModified, nil
	// return commentFeed[CommentsPerPage*pageNo : min(CommentsPerPage*(pageNo+1), len(commentFeed))], nil

}

func (s *SmartContract) GetUserProfilePosts(ctx contractapi.TransactionContextInterface, targetUserId string, userId string, pageNo int) ([]*PostModified, error) {
	var postFeed []*Post
	var postFeedModified []*PostModified
	var postList []string
	// if parentId[0] == 'p' { //If parent is post
	// 	existingPost, err := s.GetPost(ctx, parentId)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if existingPost == nil {
	// 		return nil, fmt.Errorf("Post with ID %s doesn't exists", parentId)
	// 	}
	// 	commentList = existingPost.Comments
	// } else { //If parent is comment
	// 	existingComment, err := s.GetComment(ctx, parentId)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if existingComment == nil {
	// 		return nil, fmt.Errorf("Comment with ID %s doesn't exists", parentId)
	// 	}
	// 	commentList = existingComment.Replies
	// }
	targetUser, err := s.GetUser(ctx, targetUserId)
	postList = targetUser.Posts
	if err != nil {
		return nil, err
	}
	for i := len(postList) - 1; i >= 0; i-- {
		post, err := s.GetPost(ctx, postList[i]) // Function to get a post by ID
		if post.Hidden {
			continue
		}
		if err != nil {
			return nil, err
		}
		postFeed = append(postFeed, post)
	}
	// sort.Slice(commentFeed, func(i, j int) bool {
	// 	return commentFeed[i].CreatedAt.After(commentFeed[j].CreatedAt)
	// })
	if len(postFeed) <= PostsPerPage*pageNo {
		return []*PostModified{}, nil
	}

	start := PostsPerPage * pageNo
	end := min(PostsPerPage*(pageNo+1), len(postFeed))
	for _, originalpost := range postFeed[start:end] {
		modifiedpost, error := s.convertToPostModified(ctx, originalpost, userId)
		if error != nil {
			return nil, error
		}
		postFeedModified = append(postFeedModified, modifiedpost)
	}
	return postFeedModified, nil
	// return commentFeed[CommentsPerPage*pageNo : min(CommentsPerPage*(pageNo+1), len(commentFeed))], nil

}

func (s *SmartContract) GetCommunityPosts(ctx contractapi.TransactionContextInterface, communityId string, userId string, pageNo int) ([]*PostModified, error) {
	var postFeed []*Post
	var postFeedModified []*PostModified
	var postList []string
	// if parentId[0] == 'p' { //If parent is post
	// 	existingPost, err := s.GetPost(ctx, parentId)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if existingPost == nil {
	// 		return nil, fmt.Errorf("Post with ID %s doesn't exists", parentId)
	// 	}
	// 	commentList = existingPost.Comments
	// } else { //If parent is comment
	// 	existingComment, err := s.GetComment(ctx, parentId)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if existingComment == nil {
	// 		return nil, fmt.Errorf("Comment with ID %s doesn't exists", parentId)
	// 	}
	// 	commentList = existingComment.Replies
	// }
	targetCommunity, err := s.GetCommunity(ctx, communityId)
	postList = targetCommunity.Posts
	if err != nil {
		return nil, err
	}
	for i := len(postList) - 1; i >= 0; i-- {
		post, err := s.GetPost(ctx, postList[i]) // Function to get a post by ID
		if post.Hidden {
			continue
		}
		if err != nil {
			return nil, err
		}
		postFeed = append(postFeed, post)
	}
	// sort.Slice(commentFeed, func(i, j int) bool {
	// 	return commentFeed[i].CreatedAt.After(commentFeed[j].CreatedAt)
	// })
	if len(postFeed) <= PostsPerPage*pageNo {
		return []*PostModified{}, nil
	}

	start := PostsPerPage * pageNo
	end := min(PostsPerPage*(pageNo+1), len(postFeed))
	for _, originalpost := range postFeed[start:end] {
		modifiedpost, error := s.convertToPostModified(ctx, originalpost, userId)
		if error != nil {
			return nil, error
		}
		postFeedModified = append(postFeedModified, modifiedpost)
	}
	return postFeedModified, nil
	// return commentFeed[CommentsPerPage*pageNo : min(CommentsPerPage*(pageNo+1), len(commentFeed))], nil

}

type PostOrComment struct {
	Post    *PostModified
	Comment *CommentModified
}

func (s *SmartContract) GetCommunityAppealed(ctx contractapi.TransactionContextInterface, communityId string, userId string, pageNo int) ([]*PostOrComment, error) {
	//var postFeed []*Post
	var PostOrCommentArray []*PostOrComment
	var postList []string
	// if parentId[0] == 'p' { //If parent is post
	// 	existingPost, err := s.GetPost(ctx, parentId)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if existingPost == nil {
	// 		return nil, fmt.Errorf("Post with ID %s doesn't exists", parentId)
	// 	}
	// 	commentList = existingPost.Comments
	// } else { //If parent is comment
	// 	existingComment, err := s.GetComment(ctx, parentId)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if existingComment == nil {
	// 		return nil, fmt.Errorf("Comment with ID %s doesn't exists", parentId)
	// 	}
	// 	commentList = existingComment.Replies
	// }
	targetCommunity, err := s.GetCommunity(ctx, communityId)
	postList = targetCommunity.Appealed
	if err != nil {
		return nil, err
	}
	for i := len(postList) - 1; i >= 0; i-- {
		if postList[i][0] == 'c' {
			post, err := s.GetComment(ctx, postList[i]) // Function to get a post by ID
			if post.Hidden {
				continue
			}
			if err != nil {
				return nil, err
			}
			// postFeed = append(postFeed, post)
			modifiedpost, error := s.convertToCommentModified(ctx, post, userId)
			if error != nil {
				return nil, error
			}
			PostOrCommentArray = append(PostOrCommentArray, &PostOrComment{Comment: modifiedpost})
		} else {
			post, err := s.GetPost(ctx, postList[i]) // Function to get a post by ID
			if post.Hidden {
				continue
			}
			if err != nil {
				return nil, err
			}
			// postFeed = append(postFeed, post)
			modifiedpost, error := s.convertToPostModified(ctx, post, userId)
			if error != nil {
				return nil, error
			}
			PostOrCommentArray = append(PostOrCommentArray, &PostOrComment{Post: modifiedpost})
		}
	}
	// sort.Slice(commentFeed, func(i, j int) bool {
	// 	return commentFeed[i].CreatedAt.After(commentFeed[j].CreatedAt)
	// })
	if len(postList) <= PostsPerPage*pageNo {
		return []*PostOrComment{}, nil
	}

	start := PostsPerPage * pageNo
	end := min(PostsPerPage*(pageNo+1), len(postList))
	// for _, originalpost := range postFeed[start:end] {
	// 	modifiedpost, error := s.convertToPostModified(ctx, originalpost, userId)
	// 	if error != nil {
	// 		return nil, error
	// 	}
	// 	postFeedModified = append(postFeedModified, modifiedpost)
	// }
	return PostOrCommentArray[start:end], nil
	// return commentFeed[CommentsPerPage*pageNo : min(CommentsPerPage*(pageNo+1), len(commentFeed))], nil

}

func (s *SmartContract) GetCommunityAppealedComments(ctx contractapi.TransactionContextInterface, communityId string, userId string, pageNo int) ([]*CommentModified, error) {
	var postFeed []*Comment
	var postFeedModified []*CommentModified
	var postList []string
	// if parentId[0] == 'p' { //If parent is post
	// 	existingPost, err := s.GetPost(ctx, parentId)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if existingPost == nil {
	// 		return nil, fmt.Errorf("Post with ID %s doesn't exists", parentId)
	// 	}
	// 	commentList = existingPost.Comments
	// } else { //If parent is comment
	// 	existingComment, err := s.GetComment(ctx, parentId)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if existingComment == nil {
	// 		return nil, fmt.Errorf("Comment with ID %s doesn't exists", parentId)
	// 	}
	// 	commentList = existingComment.Replies
	// }
	targetCommunity, err := s.GetCommunity(ctx, communityId)
	postList = targetCommunity.Appealed
	if err != nil {
		return nil, err
	}
	for i := len(postList) - 1; i >= 0; i-- {
		if postList[i][0] == 'p' {
			continue
		}
		post, err := s.GetComment(ctx, postList[i]) // Function to get a post by ID
		if post.Hidden {
			continue
		}
		if err != nil {
			return nil, err
		}
		postFeed = append(postFeed, post)
	}
	// sort.Slice(commentFeed, func(i, j int) bool {
	// 	return commentFeed[i].CreatedAt.After(commentFeed[j].CreatedAt)
	// })
	if len(postFeed) <= PostsPerPage*pageNo {
		return []*CommentModified{}, nil
	}

	start := PostsPerPage * pageNo
	end := min(PostsPerPage*(pageNo+1), len(postFeed))
	for _, originalpost := range postFeed[start:end] {
		modifiedpost, error := s.convertToCommentModified(ctx, originalpost, userId)
		if error != nil {
			return nil, error
		}
		postFeedModified = append(postFeedModified, modifiedpost)
	}
	return postFeedModified, nil
	// return commentFeed[CommentsPerPage*pageNo : min(CommentsPerPage*(pageNo+1), len(commentFeed))], nil

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
	// if !contains(existingCommunity.Users, userId) {
	// 	return fmt.Errorf("User cannot appeal as you are not part of the community")
	// }
	existingCommunity.Appealed = append(existingCommunity.Appealed, postId)
	communityJson, _ := json.Marshal(existingCommunity)
	ctx.GetStub().PutState(communityId, communityJson)
	return nil
}

func (s *SmartContract) isAppealed(ctx contractapi.TransactionContextInterface, postId string) (bool, error) {
	var communityId string
	if postId[0] == 'p' {
		existingPost, err := s.GetPost(ctx, postId)
		if err != nil {
			return false, err
		}
		if existingPost == nil {
			return false, fmt.Errorf("Post with ID %s doesn't exists", postId)
		}

		communityId = existingPost.Community

	} else {
		existingComment, err := s.GetComment(ctx, postId)
		if err != nil {
			return false, err
		}
		if existingComment == nil {
			return false, fmt.Errorf("Comment with ID %s doesn't exists", postId)
		}
		communityId = existingComment.Community
	}
	existingCommunity, err := s.GetCommunity(ctx, communityId)
	if err != nil {
		return false, err
	}
	if existingCommunity == nil {
		return false, fmt.Errorf("Community with ID %s doesn't exists", communityId)
	}
	return contains(existingCommunity.Appealed, postId), nil
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
		if contains(existingPost.HideVote, userId) {
			return fmt.Errorf("User already voted")
		}
		existingPost.HideCount += 1
		existingPost.HideVote = append(existingPost.HideVote, userId)
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
		if contains(existingComment.HideVote, userId) {
			return fmt.Errorf("User already voted")
		}
		existingComment.HideCount += 1
		existingComment.HideVote = append(existingComment.HideVote, userId)
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
		if contains(existingPost.ShowVote, userId) {
			return fmt.Errorf("User already voted")
		}
		existingPost.ShowCount += 1
		existingPost.ShowVote = append(existingPost.ShowVote, userId)
		if existingPost.ShowCount >= int(math.Ceil(float64(len(existingCommunity.Moderators))/2.0)) {
			//existingPost.Hidden = true
			existingPost.ShowCount = -100
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
		if contains(existingComment.ShowVote, userId) {
			return fmt.Errorf("User already voted")
		}
		existingComment.ShowCount += 1
		existingComment.ShowVote = append(existingComment.ShowVote, userId)
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
			existingComment.ShowCount = -100
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
