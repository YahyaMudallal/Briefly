package models

import "time"

type comment struct {
	id int
	author int			// id of the user who written the comment
	article int			// id of the commented article  
	date time.Time		// date on which the comment was published
	content string		// content of the comment
}