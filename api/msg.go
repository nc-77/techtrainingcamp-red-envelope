package api

import "github.com/gofiber/fiber/v2"

type Code int

var (
	msg map[Code]string
)

const (
	SUCCESS Code = iota
	ERRPARAM
	LIMITED
	MAXCOUNT
	FAILED
	TOOFAST
)

func init() {
	msg = make(map[Code]string)
	msg[SUCCESS] = "success"
	msg[ERRPARAM] = "invalid parameter"
	msg[LIMITED] = "limited"
	msg[MAXCOUNT] = "reach max_count"
	msg[FAILED] = "failed"
	msg[TOOFAST] = "this user request too fast"
}

func Response(c *fiber.Ctx, cod Code, data interface{}) error {
	return c.JSON(fiber.Map{
		"code": cod,
		"msg":  msg[cod],
		"data": data,
	})
}
