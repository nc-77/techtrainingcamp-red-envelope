package api

import "github.com/gofiber/fiber/v2"

type Code int

var (
	msg map[Code]string
)

const (
	SUCCESS Code = iota
	ERRPARAM
	DISABLED
	MAXCOUNT
	FAILED
)

func init() {
	msg = make(map[Code]string)
	msg[SUCCESS] = "success"
	msg[ERRPARAM] = "invalid parameter"
	msg[DISABLED] = "disabled"
	msg[MAXCOUNT] = "reach max_count"
	msg[FAILED] = "failed"

}

func Response(c *fiber.Ctx, cod Code, data interface{}) error {
	return c.JSON(fiber.Map{
		"Code": cod,
		"msg":  msg[cod],
		"data": data,
	})
}
