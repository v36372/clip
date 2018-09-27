package form

import (
	"github.com/gin-gonic/gin"
)

type CutForm struct {
	MinuteFrom string `form:"minute_from"`
	MinuteTo   string `form:"minute_to"`
	SecondFrom string `form:"second_from"`
	SecondTo   string `form:"second_to"`
	Name       string `form:"name"`
	Filename   string `form:"filename"`
}

func (cutForm *CutForm) FromCtx(c *gin.Context) error {
	if err := c.Bind(cutForm); err != nil {
		return err
	}

	return nil
}
