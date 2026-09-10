package forum

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// SeedDemo is opt-in and refuses to add sample records to a populated community.
// Demo authors cannot sign in: the random password is discarded immediately.
func SeedDemo(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var count int
	if err = tx.QueryRow("SELECT (SELECT COUNT(*) FROM Posts)+(SELECT COUNT(*) FROM Users)").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(randomToken()), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	names := []string{"MayaDemo", "OmarDemo", "LenaDemo", "SamDemo"}
	ids := []int64{}
	for _, name := range names {
		result, err := tx.Exec("INSERT INTO Users(username,email,password) VALUES(?,?,?)", name, name+"@example.invalid", string(hash))
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}
	samples := []struct {
		title, content, category string
		author                   int
	}{
		{"What are you learning just for the joy of it?", "I started learning to identify the stars without an app. I’m terrible at it, but there’s something lovely about being a beginner again. What are you learning with no deadline and no plan to turn it into a job?", "Education", 0},
		{"The games we keep coming back to", "Some games become a place rather than something you finish. For me, it’s a quiet evening in Stardew Valley. Which game feels like coming home, and what keeps you returning?", "Gaming", 1},
		{"A small habit that changed how I read", "I stopped trying to finish every book I start. Giving myself permission to move on has made reading feel like curiosity again, rather than homework. What changed your relationship with reading?", "Social", 2},
		{"What’s a scientific idea that still amazes you?", "The light from some stars started its journey long before humans existed. We’re looking at a sky full of different moments in time. What’s the fact or idea you keep thinking about?", "Science", 3},
		{"Does the best technology know when to get out of the way?", "I’ve been switching off features instead of adding new ones. Fewer notifications, simpler tools, and a little more room to think. What technology actually makes your day better, and what could you happily live without?", "Technology", 1},
		{"The internet feels smaller when you find your people", "A stranger explained something I’d been stuck on for a week. No pitch, no rush, just a thoughtful reply. That’s the kind of internet I want more of. What’s the best conversation you’ve had online lately?", "Social", 0},
	}
	replies := []string{
		"I’ve started sketching the plants on my balcony. Not being good at something yet is surprisingly freeing.",
		"Minecraft with the same friends every winter. The world is different each time; the ritual stays the same.",
		"Keeping a book in my bag helped more than setting a reading goal. A few pages while waiting really adds up.",
		"That trees can move water all the way up to their leaves without a mechanical pump. Everyday things are extraordinary.",
		"My favorite feature is Do Not Disturb. A simple calendar and fewer alerts have been much more useful than another productivity app.",
		"Someone walked me through fixing my first broken build. The patience mattered as much as the answer.",
	}
	for i, s := range samples {
		result, err := tx.Exec("INSERT INTO Posts(user_id,title,content,created_at) VALUES(?,?,?,?)", ids[s.author], s.title, s.content, time.Now().UTC().Add(-time.Duration(len(samples)-i)*time.Hour))
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		if _, err = tx.Exec("INSERT INTO Post_Categories(post_id,category_id) SELECT ?,id FROM Categories WHERE name=?", id, s.category); err != nil {
			return err
		}
		if _, err = tx.Exec("INSERT INTO Comments(post_id,user_id,content) VALUES(?,?,?)", id, ids[(s.author+1)%len(ids)], replies[i]); err != nil {
			return err
		}
		for j := 0; j < (i%3)+1; j++ {
			if _, err = tx.Exec("INSERT INTO Likes_Dislikes(user_id,post_id,like_dislike) VALUES(?,?,1)", ids[(s.author+j+1)%len(ids)], id); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}
