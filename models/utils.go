package models

type UtilsService struct{}

type StringDetails struct {
	Length int
}

func (u *UtilsService) StringDetails(s string) StringDetails {
	d := StringDetails{}
	d.Length = len(s)
	return d
}
