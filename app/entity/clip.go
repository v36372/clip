package entity

import (
	"clip/models"
	"clip/repo"
	"clip/utilities/uer"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type clipEntity struct {
	clipRepo repo.IClip
}

type Clip interface {
	ExecuteCutCommand(mf, mt, sf, st int, filename, name string) error
	CreateClip(mf, mt, sf, st int, filename, name, user string) (clip *models.Clip, err error)
	GetClipFromSlug(slug string) (clip *models.Clip, err error)
	ExtractFromStream() (filename string, err error)
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

func (c clipEntity) CreateClip(mf, mt, sf, st int, filename, name, user string) (clip *models.Clip, err error) {
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
	clip.Url = fmt.Sprintf("/vids/%s.mkv", encodedId)
	err = c.clipRepo.Update(clip)
	if err != nil {
		err = uer.InternalError(err)
		return
	}

	err = c.ExecuteCutCommand(mf, mt, sf, st, filename, clip.Slug)
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

func (c clipEntity) ExtractFromStream() (filename string, err error) {
	b, err := exec.Command("sh", "-c", "echo `ls -Art ./stream | tail -n 1 | cut -d'-' -f 2 | cut -d'.' -f 1`").CombinedOutput()
	if err != nil {
		return
	}

	filename = strings.TrimSpace(string(b))

	cmd := fmt.Sprintf("`./script-extract-mkv.sh %s` &>> logs/script.log", filename)
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return
	}

	return
}

func (c clipEntity) ExecuteCutCommand(mf, mt, sf, st int, filename, name string) error {
	cmd := fmt.Sprintf("`./script-cut.sh %s %d %d %d %d %s` &>> logs.script.log", filename, mf, sf, mt, st, name)
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return err
	}

	return nil
}
