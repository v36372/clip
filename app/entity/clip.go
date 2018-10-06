package entity

import (
	"clip/models"
	"clip/repo"
	"clip/utilities/uer"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
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

	if clip.IsReady == false {
		err = uer.NotFoundError(errors.New("Sorry but your clip is not ready yet, wait about 1 minute for the cutting operation to complete"))
		return
	}

	return
}

func (c clipEntity) CreateClip(mf, mt, sf, st int, filename, name, user string) (clip *models.Clip, err error) {
	clip = &models.Clip{
		Name:      name,
		CreatedBy: user,
		IsReady:   false,
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

	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("Cut command failed: %s", e)
			}
		}()

		err = c.ExecuteCutCommand(mf, mt, sf, st, filename, clip.Slug)
		if err != nil {
			fmt.Println(err)
			err = c.clipRepo.Delete(clip)
			if err != nil {
				fmt.Println(err)
				return
			}

			return
		}

		clip.IsReady = true
		err = c.clipRepo.Update(clip)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	return clip, nil
}

func existFile(path string) bool {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	if stat.IsDir() {
		return false
	}
	return true
}

func (c clipEntity) ExtractFromStream() (filename string, err error) {
	b, err := exec.Command("sh", "-c", "echo `ls -Art ./stream/hls | tail -n 1 | cut -d'-' -f 2 | cut -d'.' -f 1`").CombinedOutput()
	if err != nil {
		return
	}

	filename = strings.TrimSpace(string(b))

	b, err = exec.Command("sh", "-c", "echo `date +\"%Y%m%d\"`").CombinedOutput()
	if err != nil {
		return
	}

	dirname := strings.TrimSpace(string(b))
	filename = path.Join(dirname, filename)

	if existFile(filename) {
		return filename, nil
	}

	cmd := fmt.Sprintf("`./script-extract-mkv.sh %s` &>> ./log/script.log", filename)
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return
	}

	return
}

func (c clipEntity) ExecuteCutCommand(mf, mt, sf, st int, filename, name string) error {
	cmd := fmt.Sprintf("`./script-cut.sh %s %d %d %d %d %s` &>> ./log/script.log", filename, mf, sf, mt, st, name)
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return err
	}

	return nil
}
