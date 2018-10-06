package handler

import (
	"clip/app/entity"
	"clip/app/form"
	"clip/app/param"
	"clip/utilities/uer"
	"fmt"

	"github.com/gin-gonic/gin"
)

type clipHandler struct {
	clip entity.Clip
}

var lock = 0

func init() {
	lock = 1
}

func (h clipHandler) New(c *gin.Context) {
	presenter := Presenter{}

	filename, err := h.clip.ExtractFromStream()
	if err != nil {
		presenter.Flashes = "something went wrong, try retry?"
	}

	presenter.VideoSrc = fmt.Sprintf("/r/%s.mkv", filename)
	presenter.VideoName = filename

	c.HTML(200, "index.html", presenter)
}

func (h clipHandler) Watch(c *gin.Context) {
	slug := c.Param("slug")
	if len(slug) == 0 {
		c.Redirect(302, "/")
		return
	}

	presenter := Presenter{}
	clip, err := h.clip.GetClipFromSlug(slug)
	if err != nil {
		presenter.Flashes = err.Error()
		c.HTML(200, "clip.html", presenter)
		return
	}
	if clip == nil {
		uer.HandleErrorGin(err, c)
		return
	}

	presenter.VideoSrc = clip.Url
	presenter.VideoName = clip.Name
	c.HTML(200, "clip.html", presenter)
}

func (h clipHandler) Cut(c *gin.Context) {
	if lock == 0 {
		presenter := Presenter{
			Flashes: "there is someone clipping right now, our server can only handle 1 clipping operation at a time",
		}
		c.HTML(200, "index.html", presenter)
		return
	}
	lock = 0
	defer func() {
		lock = 1
	}()
	cutForm := form.CutForm{}
	err := c.Bind(&cutForm)
	if err != nil {
		presenter := Presenter{
			Flashes: err.Error(),
		}
		c.HTML(200, "index.html", presenter)
		return
	}

	mf := param.GetIntFromStrWithDefault(cutForm.MinuteFrom, 0)
	mt := param.GetIntFromStrWithDefault(cutForm.MinuteTo, 0)

	sf := param.GetIntFromStrWithDefault(cutForm.SecondFrom, 0)
	st := param.GetIntFromStrWithDefault(cutForm.SecondTo, 0)

	if mf < 0 || mt < 0 || sf < 0 || st < 0 {
		presenter := Presenter{
			Flashes: "invalid clip timestamp",
		}
		c.HTML(200, "index.html", presenter)
		return
	}
	if len(cutForm.Name) == 0 {
		presenter := Presenter{
			Flashes: "clip must has a name",
		}
		c.HTML(200, "index.html", presenter)
		return
	}

	if (mf > mt) || (mf == mt && st-sf < 5) {
		presenter := Presenter{
			Flashes: "clip length must be longer than 5 seconds",
		}
		c.HTML(200, "index.html", presenter)
		return
	}

	clip, err := h.clip.CreateClip(mf, mt, sf, st, cutForm.Filename, cutForm.Name, "admin")
	if err != nil {
		presenter := Presenter{
			Flashes: "something went wrong",
		}
		c.HTML(200, "index.html", presenter)
		return

	}
	c.Redirect(302, fmt.Sprintf("/s/%s", clip.Slug))
}
