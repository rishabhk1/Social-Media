package chaincode

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"sort"
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

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// user := User{
	// 	ID:       "1",
	// 	Username: "john_doe",
	// 	Email:    "john@example.com",
	// }

	// comment := Comment{
	// 	ID:        "3",
	// 	Content:   "Demo Comment",
	// 	Author:    "1",
	// 	Score:     5,
	// 	CreatedAt: time.Now(),
	// }

	// post := Post{
	// 	ID:        "4",
	// 	Title:     "Demo Post",
	// 	Content:   "Demo post content",
	// 	Author:    "1",
	// 	Score:     10,
	// 	CreatedAt: time.Now(),
	// 	Comments:  []Comment{comment},
	// }

	// community := Community{
	// 	ID:          "2",
	// 	Name:        "Comm1",
	// 	Description: "Demo Community",
	// 	Creator:     "1",
	// 	CreatedAt:   time.Now(),
	// 	Moderators:  []string{"1"},
	// 	Posts:       []Post{post},
	// }

	// userJSON, err := json.Marshal(user)
	// if err != nil {
	// 	return err
	// }

	// err = ctx.GetStub().PutState(user.ID, userJSON)
	// if err != nil {
	// 	return err
	// }

	// communityJson, err := json.Marshal(community)
	// if err != nil {
	// 	return err
	// }

	// err = ctx.GetStub().PutState(community.ID, communityJson)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func findIndex(slice []string, target string) int {
	for index, value := range slice {
		if value == target {
			return index
		}
	}
	return -1 // Return -1 if the element is not found
}

// Define a function to remove an element from a slice by index
func removeElement(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

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
	}
	commentJson, _ := json.Marshal(comment)
	ctx.GetStub().PutState(commentId, commentJson)
	return nil
}

func (s *SmartContract) GetUserFeed(ctx contractapi.TransactionContextInterface, userId string, pageNo int) ([]*Post, error) {
	var userFeed []*Post
	existingUser, err := s.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, fmt.Errorf("User with ID %s doesn't exists", userId)
	}

	for _, community := range existingUser.Communities {
		communityPosts := []*Post{}
		existingCommunity, err := s.GetCommunity(ctx, community)
		if err != nil {
			return nil, err
		}
		if existingCommunity == nil {
			return nil, fmt.Errorf("Community with ID %s doesn't exists", community)
		}
		for _, postId := range existingCommunity.Posts {
			post, err := s.GetPost(ctx, postId) // Function to get a post by ID
			if post.Hidden {
				continue
			}
			if err != nil {
				return nil, err
			}
			communityPosts = append(communityPosts, post)
		}

		// Add community posts to the user's feed
		userFeed = append(userFeed, communityPosts...)
	}

	// Sort the user's feed by timestamp (reverse chronological order)
	sort.Slice(userFeed, func(i, j int) bool {
		return userFeed[i].CreatedAt.After(userFeed[j].CreatedAt)
	})
	if len(userFeed) <= PostsPerPage*pageNo {
		return []*Post{}, nil
	}
	return userFeed[PostsPerPage*pageNo : min(PostsPerPage*(pageNo+1), len(userFeed))], nil
}

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
	if !slices.Contains(existingCommunity.Users, userId) {
		return fmt.Errorf("User cannot appeal as you are not part of the community")
	}
	existingCommunity.Appealed = append(existingCommunity.Appealed, postId)
	communityJson, _ := json.Marshal(existingCommunity)
	ctx.GetStub().PutState(communityId, communityJson)
	return nil
}

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
		if !slices.Contains(existingCommunity.Moderators, userId) {
			return fmt.Errorf("User cannot hide as you are not a moderator")
		}
		existingPost.HideCount += 1
		if existingPost.HideCount >= int(math.Ceil(float64(len(existingCommunity.Moderators))/2.0)) {
			existingPost.Hidden = true
			existingCommunity.Appealed = removeElement(existingCommunity.Appealed, findIndex(existingCommunity.Appealed, postId))
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
		if !slices.Contains(existingCommunity.Moderators, userId) {
			return fmt.Errorf("User cannot hide as you are not a moderator")
		}
		existingComment.HideCount += 1
		if existingComment.HideCount >= int(math.Ceil(float64(len(existingCommunity.Moderators))/2.0)) {
			existingComment.Hidden = true
			existingCommunity.Appealed = removeElement(existingCommunity.Appealed, findIndex(existingCommunity.Appealed, postId))
			communityJson, _ := json.Marshal(existingCommunity)
			ctx.GetStub().PutState(communityId, communityJson)
		}
		commentJson, _ := json.Marshal(existingComment)
		ctx.GetStub().PutState(postId, commentJson)
	}
	return nil
}

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
	var userMap map[string]int
	for _, postId := range existingCommunity.Posts {
		existingPost, err := s.GetPost(ctx, postId)
		if err != nil {
			return err
		}
		if existingPost == nil {
			return fmt.Errorf("Post with ID %s doesn't exists", postId)
		}
		userMap[existingPost.Author] += existingPost.Score
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
	}
	existingCommunity.Moderators = newModerators
	communityJson, _ := json.Marshal(existingCommunity)
	ctx.GetStub().PutState(communityId, communityJson)
	return nil
}
