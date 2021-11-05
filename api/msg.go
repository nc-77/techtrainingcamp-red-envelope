package api

import "github.com/gofiber/fiber/v2"

type code int

var (
	msg map[code]string
)

const (
	SUCCESS code = iota
	ERRPARAM
	DISABLED
	MAXCOUNT
	FAILED
)

func init() {
	msg = make(map[code]string)
	msg[SUCCESS] = "success"
	msg[ERRPARAM] = "invalid parameter"
	msg[DISABLED] = "disabled"
	msg[MAXCOUNT] = "reach max_count"
	msg[FAILED] = "failed"
}

func Response(c *fiber.Ctx, cod code, data interface{}) error {
	return c.JSON(fiber.Map{
		"code": cod,
		"msg":  msg[cod],
		"data": data,
	})
}
