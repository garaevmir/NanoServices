package scenarios

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nanoservices/tests/api"
	"github.com/nanoservices/tests/models"
	"github.com/nanoservices/tests/utils"
)

func Scenario3(user models.User) {
	token, err := api.LoginUser(user)
	if err != nil {
		fmt.Println("[ ] Login failed:", err)
		return
	}
	fmt.Println("[*] Login successful")

	updateReq := models.ProfileUpdateRequest{
		FirstName: utils.RandomString(8),
		LastName:  utils.RandomString(10),
		Bio:       fmt.Sprintf("Updated bio at %s", time.Now().Format(time.RFC3339)),
	}
	if err := api.UpdateProfile(token, updateReq); err != nil {
		fmt.Println("[ ] Update profile failed:", err)
		return
	}
	fmt.Println("[*] Profile updated")

	profile, err := api.GetProfile(token)
	if err != nil {
		fmt.Println("[ ] Get profile failed:", err)
		return
	}
	fmt.Printf("[*] Updated profile: %s %s, Bio: %s\n",
		profile.FirstName, profile.LastName, profile.Bio)

	posts, total, err := api.ListPosts(token)
	if err != nil {
		fmt.Println("[ ] List posts failed:", err)
		return
	}
	fmt.Printf("[*] Found %d posts (total: %d)\n", len(posts), total)

	if len(posts) == 0 {
		fmt.Println("[ ] No posts available for commenting")
		return
	}

	rand.Seed(time.Now().UnixNano())
	randomPost := posts[rand.Intn(len(posts))]
	comment := fmt.Sprintf("Great post! Commented at %s", time.Now().Format(time.RFC3339))

	if err := api.AddComment(token, randomPost.ID, comment); err != nil {
		fmt.Println("[ ] Add comment failed:", err)
		return
	}
	fmt.Printf("[*] Comment added to post '%s'\n", randomPost.ID)

	fmt.Println("[...] Waiting for stats processing (2s)...")
	time.Sleep(2 * time.Second)

	topPosts, err := api.GetTopPosts(token, "comments", "7d")
	if err != nil {
		fmt.Println("[ ] Get top posts failed:", err)
		return
	}

	fmt.Println("[*] Top posts by comments (7 days):")
	for i, post := range topPosts {
		fmt.Printf("%d. %s (ID: %s) - %d comments\n",
			i+1, post.Title, post.PostID, post.Count)
	}
}
