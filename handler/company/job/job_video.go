package job

import (
	"mime"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/app/adapters/company/job"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

func SaveJobVideo(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	video, err := fiberCtx.FormFile("video")
	if err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	contentType := video.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(video.Filename))
	}

	jobVideo, err := job.SaveJobVideo(fiberCtx.UserContext(), user.CompanyID, video, contentType)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, map[string]string{
		"video_link": jobVideo,
	})
}
