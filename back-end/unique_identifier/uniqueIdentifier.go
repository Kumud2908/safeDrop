package uniqueidentifier

import "github.com/dchest/uniuri"

func GenerateID() string {
	return uniuri.New()
}
