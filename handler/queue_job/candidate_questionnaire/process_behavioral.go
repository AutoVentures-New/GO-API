package candidate_questionnaire

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	questionnaire_adp "github.com/hubjob/api/app/adapters/queue_job/candidate/questionnaire"
	"github.com/hubjob/api/config"
	"github.com/hubjob/api/model"
	"github.com/hubjob/api/pkg"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
)

const spreadsheetBehavioralID = "1T4XesNEU4bTJwS_jJQX9mWkbJcRTpzIDCJbpnc7ak0Y"
const spreadsheetBehavioralTextID = "1NKZkgyoCZ3VB8HNXfSJX2yvpwLaaZCJsc8YXXExxFIQ"

func ExecuteQuestionnaireBehavioral(
	ctx context.Context,
	queueJob model.QueueJob,
) error {
	logrus.WithFields(logrus.Fields{
		"queue_job": queueJob,
	}).Info("Executing questionnaire behavioral")

	questionnaire, err := questionnaire_adp.GetCandidateQuestionnaire(ctx, queueJob.Configurations.CandidateID, model.BEHAVIORAL)
	if err != nil {
		return err
	}

	candidate, err := questionnaire_adp.GetCandidate(ctx, queueJob.Configurations.CandidateID)
	if err != nil {
		return err
	}

	categoryResult := map[string]int64{
		"E": 0,
		"I": 0,
		"S": 0,
		"N": 0,
		"T": 0,
		"F": 0,
		"J": 0,
		"P": 0,
	}

	for _, answer := range questionnaire.Answers {
		category := model.QuestionnaireBehavioralCategories[answer.QuestionID]
		index := category.AnswerA

		if answer.Answer == "answer_b" {
			index = category.AnswerB
		}

		categoryResult[index]++
	}

	result1 := "E"
	result2 := "S"
	result3 := "T"
	result4 := "J"

	if (categoryResult["E"] - categoryResult["I"]) < 0 {
		result1 = "I"
	}

	if (categoryResult["S"] - categoryResult["N"]) < 0 {
		result2 = "N"
	}

	if (categoryResult["T"] - categoryResult["F"]) < 0 {
		result3 = "F"
	}

	if (categoryResult["J"] - categoryResult["P"]) < 0 {
		result4 = "P"
	}

	resultFinal := result1 + result2 + result3 + result4

	questionnaire.ResultFilePath, err = createBehavioralResultFile(ctx, candidate, questionnaire, resultFinal)
	if err != nil {
		logrus.WithError(err).
			WithField("category_result", categoryResult).
			WithField("candidate", candidate).
			Error("Error to create result file")

		return err
	}

	questionnaire.BucketName = config.Config.S3.Bucket

	err = questionnaire_adp.UpdateCandidateQuestionnaire(
		ctx,
		questionnaire,
	)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"queue_job": queueJob,
	}).Info("Questionnaire behavioral executed with success")

	return nil
}

func createBehavioralResultFile(
	ctx context.Context,
	candidate model.Candidate,
	questionnaire model.CandidateQuestionnaire,
	resultFinal string,
) (string, error) {
	candidateFileName := fmt.Sprintf("result_behavioral_%d_%d", candidate.ID, questionnaire.ID)

	candidateFileId, err := createBehavioralFileCopy(ctx, candidateFileName)
	if err != nil {
		return "", err
	}

	_, err = pkg.SheetsService.Spreadsheets.Values.Update(
		candidateFileId,
		"result!B3",
		&sheets.ValueRange{Values: [][]interface{}{{candidate.Name}}},
	).ValueInputOption("RAW").Do()
	if err != nil {
		logrus.WithError(err).Error("Error to update value")

		return "", err
	}

	_, err = pkg.SheetsService.Spreadsheets.Values.Update(
		candidateFileId,
		"result!H2",
		&sheets.ValueRange{Values: [][]interface{}{{time.Now().UTC().Format("02/01/2006")}}},
	).ValueInputOption("RAW").Do()
	if err != nil {
		logrus.WithError(err).Error("Error to update value")

		return "", err
	}

	_, err = pkg.SheetsService.Spreadsheets.Values.Update(
		candidateFileId,
		"result!H3",
		&sheets.ValueRange{Values: [][]interface{}{{resultFinal}}},
	).ValueInputOption("RAW").Do()
	if err != nil {
		logrus.WithError(err).Error("Error to update value")

		return "", err
	}

	err = updateCellValue(ctx, candidateFileId, "A19", model.CaracteristicasResult[resultFinal])
	if err != nil {
		return "", err
	}

	err = updateCellValue(ctx, candidateFileId, "A21", model.EstiloDeInteracaoResult[resultFinal])
	if err != nil {
		return "", err
	}

	err = updateCellValue(ctx, candidateFileId, "A23", model.TemperamentoResult[resultFinal])
	if err != nil {
		return "", err
	}

	err = updateCellValue(ctx, candidateFileId, "A25", model.ComoMembroDeEquipeResult[resultFinal])
	if err != nil {
		return "", err
	}

	err = updateCellValue(ctx, candidateFileId, "A27", model.PreferenciasComoLiderResult[resultFinal])
	if err != nil {
		return "", err
	}

	err = updateCellValue(ctx, candidateFileId, "A29", model.PreferenciasNaturaisDeTrabalhoResult[resultFinal])
	if err != nil {
		return "", err
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

func createBehavioralFileCopy(ctx context.Context, newFileName string) (string, error) {
	copyFile := &drive.File{
		Name: newFileName,
	}

	copiedFile, err := pkg.DriveService.Files.Copy(spreadsheetBehavioralID, copyFile).Do()
	if err != nil {
		logrus.WithError(err).Error("Error to copy file behavioral")

		return "", err
	}

	return copiedFile.Id, nil
}

func updateCellValue(ctx context.Context, candidateFileId, cellToId, cellFromId string) error {
	cellFrom := fmt.Sprintf("text!%s", cellFromId)
	cellTo := fmt.Sprintf("result!%s", cellToId)

	resp, err := pkg.SheetsService.Spreadsheets.Values.Get(spreadsheetBehavioralTextID, cellFrom).Do()
	if err != nil {
		logrus.WithError(err).Error("Error to get value")

		return err
	}

	if len(resp.Values) == 0 {
		return fmt.Errorf("no data found in the specified range")
	}

	_, err = pkg.SheetsService.Spreadsheets.Values.Update(
		candidateFileId,
		cellTo,
		&sheets.ValueRange{Values: resp.Values},
	).ValueInputOption("RAW").Do()
	if err != nil {
		logrus.WithError(err).Error("Error to update value")

		return err
	}

	return nil
}
