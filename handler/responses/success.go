package responses

import (
	"mime"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
)

func Success(
	fiberCtx *fiber.Ctx,
	data any,
) error {
	jsonData := make(map[string]string)

	if data == nil {
		data = jsonData
	}

	return fiberCtx.Status(fiber.StatusOK).
		JSON(fiber.Map{
			"data": data,
		})
}

func Download(
	fiberCtx *fiber.Ctx,
	filename string,
	result *s3.GetObjectOutput,
) error {
	extensions, err := mime.ExtensionsByType(*result.ContentType)
	if err == nil && len(extensions) > 0 {
		filename += extensions[0]
	}

	fiberCtx.Set("Content-Disposition", "attachment; filename="+filename)
	fiberCtx.Set("Content-Type", *result.ContentType)

	return fiberCtx.SendStream(result.Body)
}
