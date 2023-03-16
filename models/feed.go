package ktn

type Feed struct {
	WebUrl        string
	FeedReference string
	FeedTitle     string
	EmailDomain   string
	Entries       []Entry
}

type NewFeed struct {
	Title     string
	Reference string
}

type ORMFeed struct {
	Id        int
	Reference string
	Title     string
	CreatedAt string
	UpdatedAt string
}

func (f Feed) FeedUpdatedAtRfc3339() string {
	return "some date"
}

func (f Feed) CreatedAtRfc3339() string {
	return "some date"
}
