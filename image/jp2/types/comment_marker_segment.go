package types

import "gitee.com/kayi-cloud/gdimse/example/image/jp2/types/comment_registration_value"

type CommentMarkerSegment struct {
	registrationValue []uint8
	comment           []byte
}

func (c *CommentMarkerSegment) RegistrationValue() comment_registration_value.CommentRegistrationValue {
	panic("")
}

func (c *CommentMarkerSegment) Comment() string {
	panic("")
}
