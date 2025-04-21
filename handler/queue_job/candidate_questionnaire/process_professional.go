package candidate_questionnaire

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/hubjob/api/app/adapters/queue_job/candidate/questionnaire/professional"
	"github.com/hubjob/api/config"
	"github.com/hubjob/api/model"
	"github.com/hubjob/api/pkg"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
)

const spreadsheetID = "1qn0wJtyU6q6AIdWBjt1rSe1wGUcBi85QEKTCCt178dE"

type PdfResponse struct {
	URL              string `json:"url"`
	Error            bool   `json:"error"`
	StatusCode       int    `json:"status"`
	Name             string `json:"name"`
	RemainingCredits int64  `json:"remainingCredits"`
}

func ExecuteQuestionnaireProfessional(
	ctx context.Context,
	queueJob model.QueueJob,
) error {
	logrus.WithFields(logrus.Fields{
		"queue_job": queueJob,
	}).Info("Executing questionnaire professional")

	questionnaire, err := professional.GetCandidateQuestionnaire(ctx, queueJob.Configurations.CandidateID)
	if err != nil {
		return err
	}

	candidate, err := professional.GetCandidate(ctx, queueJob.Configurations.CandidateID)
	if err != nil {
		return err
	}

	categoryResult := make(map[int64]int64)

	for i := 1; i <= 20; i++ {
		categoryResult[int64(i)] = int64(0)
	}

	for _, answer := range questionnaire.Answers {
		category := model.QuestionsCategoryMap[answer.QuestionID]
		result := category.AnswerA

		if answer.Answer == "answer_b" {
			result = category.AnswerB
		}

		categoryResult[result]++
	}

	questionnaire.ResultFilePath, err = createResultFile(ctx, categoryResult, candidate, questionnaire)
	if err != nil {
		logrus.WithError(err).
			WithField("category_result", categoryResult).
			WithField("candidate", candidate).
			Error("Error to create result file")

		return err
	}

	questionnaire.BucketName = config.Config.S3.Bucket

	err = professional.UpdateCandidateQuestionnaire(
		ctx,
		questionnaire,
	)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"queue_job": queueJob,
	}).Info("Questionnaire professional executed with success")

	return nil
}

func createResultFile(
	ctx context.Context,
	categoryResult map[int64]int64,
	candidate model.Candidate,
	questionnaire model.CandidateQuestionnaire,
) (string, error) {
	candidateFileName := fmt.Sprintf("result_professional_%d_%d", candidate.ID, questionnaire.ID)

	candidateFileId, err := createFileCopy(ctx, candidateFileName)
	if err != nil {
		return "", err
	}

	_, err = pkg.SheetsService.Spreadsheets.Values.Update(
		candidateFileId,
		"result!B4",
		&sheets.ValueRange{Values: [][]interface{}{{candidate.Name}}},
	).ValueInputOption("RAW").Do()
	if err != nil {
		logrus.WithError(err).Error("Error to update value")

		return "", err
	}

	_, err = pkg.SheetsService.Spreadsheets.Values.Update(
		candidateFileId,
		"result!E4",
		&sheets.ValueRange{Values: [][]interface{}{{time.Now().UTC().Format("02/01/2006")}}},
	).ValueInputOption("RAW").Do()
	if err != nil {
		logrus.WithError(err).Error("Error to update value")

		return "", err
	}

	for category, result := range categoryResult {
		_, err = pkg.SheetsService.Spreadsheets.Values.Update(
			candidateFileId,
			"result!"+model.ResultCellLocation[category],
			&sheets.ValueRange{Values: [][]interface{}{{model.ResultScore[category][result]}}},
		).ValueInputOption("RAW").Do()
		if err != nil {
			logrus.WithError(err).Error("Error to update value")

			return "", err
		}
	}

	resp, err := pkg.DriveService.Files.Export(
		candidateFileId,
		"application/pdf",
	).Download()
	if err != nil {
		logrus.WithError(err).Error("Error to export to pdf")

		return "", err
	}

	defer resp.Body.Close()

	resultFileName := fmt.Sprintf("candidates/%d/%s.pdf", candidate.ID, candidateFileName)

	_, err = pkg.S3Uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(config.Config.S3.Bucket),
		Key:         aws.String(resultFileName),
		Body:        resp.Body,
		ContentType: aws.String("application/pdf"),
	})
	if err != nil {
		logrus.WithError(err).Error("Error to upload result file")

		return "", err
	}

	err = pkg.DriveService.Files.Delete(candidateFileId).Do()
	if err != nil {
		logrus.WithError(err).Error("Error to delete candidate file copy")
	}

	return resultFileName, nil
}

func createFileCopy(ctx context.Context, newFileName string) (string, error) {
	copyFile := &drive.File{
		Name: newFileName,
	}

	copiedFile, err := pkg.DriveService.Files.Copy(spreadsheetID, copyFile).Do()
	if err != nil {
		logrus.WithError(err).Error("Error to copy file")

		return "", err
	}

	return copiedFile.Id, nil
}
