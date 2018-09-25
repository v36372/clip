package entity

import (
	"clip/models"
	"clip/repo"
	"clip/utilities/uer"
	"errors"
	"fmt"
	"master/utils"
	"os/exec"
	"strings"
)

const (
	padLength            = 3
	maxClipIdCanGenerate = 181
	clipCodeLength       = 3
	flagClip             = "12"
)

var (
	padDict = [4]string{
		"Six",
		"Two",
		"One",
		"Ten",
	}
	padHash = map[string]int{
		"Six": 0,
		"Two": 1,
		"One": 2,
	}
	orderDict = [14]string{
		"Dota2",
		"Nhson",
		"Trung",
		"Nttin",
		"Ntien",
		"QTrux",
		"LMinh",
		"Quang",
		"PHang",
		"SCoii",
		"BMinh",
		"AKhoa",
		"Huann",
		"Phong",
	}
	ranDict = [37]string{
		"Absent",
		"Messy",
		"Wasteful",
		"Super",
		"Obese",
		"Disgusted",
		"Smiling",
		"Tired",
		"Remarkable",
		"Undesirable",
		"Fantastic",
		"Modern",
		"Friendly",
		"Shut",
		"Tricky",
		"Dead",
		"Lazy",
		"Pink",
		"Yellow",
		"Short",
		"Sudden",
		"Lethal",
		"Sincere",
		"Present",
		"Bright",
		"Fabulous",
		"Precious",
		"Poor",
		"Weak",
		"Ugly",
		"Mad",
		"Old",
		"Nice",
		"Rare",
		"Tall",
		"Odd",
		"Tiny",
	}
)

type clipEntity struct {
	clipRepo repo.IClip
}

type Clip interface {
	ExecuteCutCommand(mf, mt, sf, st int, name string) error
	CreateClip(mf, mt, sf, st int, name, user string) (clip *models.Clip, err error)
	GetClipFromSlug(slug string) (clip *models.Clip, err error)
}

func NewClip(clipRepo repo.IClip) Clip {
	return &clipEntity{
		clipRepo: clipRepo,
	}
}

func (c clipEntity) GetClipFromSlug(slug string) (clip *models.Clip, err error) {
	clipId := decode(slug)
	if clipId > maxClipIdCanGenerate {
		err = uer.NotFoundError(errors.New("Clip not found"))
		return
	}

	clip, err = c.clipRepo.GetById(clipId)
	if err != nil {
		err = uer.InternalError(err)
		return
	}

	return
}

func (c clipEntity) CreateClip(mf, mt, sf, st int, name, user string) (clip *models.Clip, err error) {
	clip = &models.Clip{
		Name:      name,
		CreatedBy: user,
	}
	clip, err = c.clipRepo.Create(clip)
	if err != nil {
		err = uer.InternalError(err)
		return
	}

	if clip.Id > maxClipIdCanGenerate {
		err = c.clipRepo.Delete(clip)
		if err != nil {
			err = uer.InternalError(err)
			return
		}

		err = uer.BadRequestError(errors.New("Maximum clip reached!"))
		return
	}

	encodedId := encode(clip.Id)
	clip.Slug = encodedId
	clip.Url = fmt.Sprintf("/vids/%s.mp4", encodedId)
	err = c.clipRepo.Update(clip)
	if err != nil {
		err = uer.InternalError(err)
		return
	}

	err = c.ExecuteCutCommand(mf, mt, sf, st, clip.Slug)
	if err != nil {
		err = c.clipRepo.Delete(clip)
		if err != nil {
			err = uer.InternalError(err)
			return
		}

		err = uer.InternalError(err)
		return
	}

	return clip, nil
}

func (c clipEntity) ExecuteCutCommand(mf, mt, sf, st int, name string) error {
	cmd := fmt.Sprintf("ffmpeg -i ./vids/output5.mp4 -ss 00:%d:%d -t 00:%d:%d -async 1 -y -strict -2 ./vids/%s.mp4", mf, sf, mt, st, name)
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return err
	}

	return nil
}

func encode(value int) string {
	base := len(orderDict) - 1
	var a [4]string
	index := len(a) //4

	// general case
	b := base //13
	for value > b {
		index--                                 // 63 //62
		q := (value / b)                        // 14/13
		a[index] = "And" + orderDict[value-q*b] // 14-0=14
		value = q                               // 102 // 3
	}
	// value < base
	index--                     // 61
	a[index] = orderDict[value] // 3

	// Padding
	c := index - 1 // 3 - 1 = 2
	used := map[int]bool{
		0: true,
	}
	for index > 0 {
		index--
		var ranIndex int
		for used[ranIndex] == true {
			ranIndex = utils.GetRandomNumber(len(orderDict))
		}
		a[index] = ranDict[ranIndex]
		used[ranIndex] = true
	}

	a[index] = padDict[c]

	return strings.Join(a[index:], "")
}

func decode(code string) int {
	flag := code[0:3]
	pad := padLength - padHash[flag]
	rawCode := code[len(code)-pad*5-(pad-1)*3 : len(code)]
	bits := strings.Split(rawCode, "And")

	raw := 0

	for i := len(bits) - 1; i >= 0; i-- {
		// q: value in base len(dictionary)
		for q, digit := range orderDict {
			if digit == bits[i] {
				t := q
				for j := 0; j < len(bits)-1-i; j++ {
					t = t * (len(orderDict) - 1)
				}
				raw += t
				break
			}
		}
	}

	return raw
}
