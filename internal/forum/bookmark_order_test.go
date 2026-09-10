package forum

import "testing"

func TestSavedDiscussionsFollowSaveOrder(t *testing.T) {
	s := newSite(t)
	s.signup(t)
	createDiscussion(t, s, "An older discussion")
	createDiscussion(t, s, "A newer discussion")
	for _, id := range []string{"2", "1", "2"} {
		r, b := bookmarkJSON(t, s, id, "save")
		checkStatus(t, r, 200, b)
	}
	app := New(s.db, false)
	posts, err := app.posts(postQuery{Author: 1, Viewer: 1, Saved: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 2 || posts[0].ID != 1 {
		t.Fatal("repeat save reordered the private collection")
	}
	r, b := bookmarkJSON(t, s, "2", "remove")
	checkStatus(t, r, 200, b)
	r, b = bookmarkJSON(t, s, "2", "save")
	checkStatus(t, r, 200, b)
	posts, err = app.posts(postQuery{Author: 1, Viewer: 1, Saved: true})
	if err != nil {
		t.Fatal(err)
	}
	if posts[0].ID != 2 {
		t.Fatal("a fresh save did not move the discussion first")
	}
}
