// cr is short for control routes
package cr

import (
	"github.com/flex-smartlock/barikata-backend/libs"
	"github.com/labstack/echo/v4"
)

func Test(c echo.Context) error {
	libs.MQTT_Test()
	return c.String(200, "ok")
}
