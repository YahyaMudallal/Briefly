package models

type article struct {
	id int
	name string
	description string
	upvote_number int	
	downvote_number int
	tl_dr string
}