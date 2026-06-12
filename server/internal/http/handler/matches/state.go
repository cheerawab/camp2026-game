package matches

import (
	"context"
	"errors"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

func (h *Handler) findMatchByID(ctx context.Context, matchID string) (mongomodel.Match, error) {
	var match mongomodel.Match
	err := h.db.Collection(mongomodel.MatchesCollection).
		FindOne(ctx, bson.M{"_id": matchID}).
		Decode(&match)
	return match, err
}

func (h *Handler) findMatchByCode(ctx context.Context, code string) (mongomodel.Match, error) {
	var match mongomodel.Match
	err := h.db.Collection(mongomodel.MatchesCollection).
		FindOne(ctx, bson.M{"code": code, "status": mongomodel.MatchStatusWaiting}).
		Decode(&match)
	return match, err
}

func (h *Handler) saveMatch(ctx context.Context, match mongomodel.Match) error {
	_, err := h.db.Collection(mongomodel.MatchesCollection).
		ReplaceOne(ctx, bson.M{"_id": match.ID}, match)
	return err
}

func (h *Handler) findAnswers(ctx context.Context, matchID string) ([]mongomodel.MatchAnswer, error) {
	cursor, err := h.db.Collection(mongomodel.MatchAnswersCollection).Find(
		ctx,
		bson.M{"match_id": matchID},
		options.Find().SetSort(bson.D{
			{Key: "question_id", Value: 1},
			{Key: "answered_at", Value: 1},
		}),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var answers []mongomodel.MatchAnswer
	if err := cursor.All(ctx, &answers); err != nil {
		return nil, err
	}
	return answers, nil
}

func (h *Handler) findAnswer(ctx context.Context, matchID, playerID, questionID string) (mongomodel.MatchAnswer, error) {
	var answer mongomodel.MatchAnswer
	err := h.db.Collection(mongomodel.MatchAnswersCollection).
		FindOne(ctx, bson.M{
			"match_id":    matchID,
			"player_id":   playerID,
			"question_id": questionID,
		}).
		Decode(&answer)
	return answer, err
}

func (h *Handler) advanceMatch(ctx context.Context, match mongomodel.Match, now time.Time) (mongomodel.Match, []string, error) {
	if match.Status != mongomodel.MatchStatusActive {
		return match, nil, nil
	}

	var events []string
	for match.Status == mongomodel.MatchStatusActive {
		shouldAdvance, err := h.shouldAdvanceRound(ctx, match, now)
		if err != nil {
			return match, events, err
		}
		if !shouldAdvance {
			break
		}

		if match.CurrentQuestionIndex >= len(match.QuestionIDs)-1 {
			match.Status = mongomodel.MatchStatusCompleted
			match.CompletedAt = now
			if err := h.saveMatch(ctx, match); err != nil {
				return match, events, err
			}
			if err := h.writeMatchRewards(ctx, match); err != nil {
				return match, events, err
			}
			events = append(events, "match_completed")
			break
		}

		nextStart := now
		if !match.RoundEndsAt.IsZero() && now.After(match.RoundEndsAt) {
			nextStart = match.RoundEndsAt
		}
		match.CurrentQuestionIndex++
		match.RoundStartedAt = nextStart
		match.RoundEndsAt = nextStart.Add(roundDuration * time.Second)
		if err := h.saveMatch(ctx, match); err != nil {
			return match, events, err
		}
		events = append(events, "round_started")
	}

	return match, events, nil
}

func (h *Handler) shouldAdvanceRound(ctx context.Context, match mongomodel.Match, now time.Time) (bool, error) {
	if len(match.QuestionIDs) == 0 || match.CurrentQuestionIndex >= len(match.QuestionIDs) {
		return false, nil
	}
	if !match.RoundEndsAt.IsZero() && !now.Before(match.RoundEndsAt) {
		return true, nil
	}

	questionID := match.QuestionIDs[match.CurrentQuestionIndex]
	answers, err := h.findAnswers(ctx, match.ID)
	if err != nil {
		return false, err
	}

	answered := make(map[string]struct{}, len(match.Players))
	for _, answer := range answers {
		if answer.QuestionID == questionID {
			answered[answer.PlayerID] = struct{}{}
		}
	}
	for _, player := range match.Players {
		if _, ok := answered[player.PlayerID]; !ok {
			return false, nil
		}
	}
	return len(match.Players) == 2, nil
}

func (h *Handler) writeMatchRewards(ctx context.Context, match mongomodel.Match) error {
	now := match.CompletedAt
	if now.IsZero() {
		now = time.Now()
	}

	for _, player := range match.Players {
		record := mongomodel.OpenPowerRecord{
			ID:        matchRewardRecordID(match.ID, player.PlayerID),
			PlayerID:  player.PlayerID,
			Amount:    openPowerReward(player.Score),
			Reason:    "quiz_match_completed",
			Source:    matchRewardSource(match.ID, player.PlayerID),
			CreatedAt: now,
		}
		_, err := h.db.Collection(mongomodel.OpenPowerRecordsCollection).UpdateOne(
			ctx,
			bson.M{"_id": record.ID},
			bson.M{"$setOnInsert": record},
			options.UpdateOne().SetUpsert(true),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) buildMatchState(ctx context.Context, match mongomodel.Match) (MatchStateResponse, error) {
	answers, err := h.findAnswers(ctx, match.ID)
	if err != nil {
		return MatchStateResponse{}, err
	}

	answerByQuestionPlayer := make(map[string]map[string]mongomodel.MatchAnswer)
	for _, answer := range answers {
		if _, ok := answerByQuestionPlayer[answer.QuestionID]; !ok {
			answerByQuestionPlayer[answer.QuestionID] = make(map[string]mongomodel.MatchAnswer)
		}
		answerByQuestionPlayer[answer.QuestionID][answer.PlayerID] = answer
	}

	response := MatchStateResponse{
		MatchID:              match.ID,
		Code:                 match.Code,
		Status:               match.Status,
		HostPlayerID:         match.HostPlayerID,
		Players:              make([]MatchPlayerResponse, 0, len(match.Players)),
		CurrentQuestionIndex: match.CurrentQuestionIndex,
		QuestionCount:        len(match.QuestionIDs),
		CreatedAt:            match.CreatedAt,
		StartedAt:            timePtr(match.StartedAt),
		CompletedAt:          timePtr(match.CompletedAt),
	}

	if match.Status == mongomodel.MatchStatusActive {
		response.RoundStartedAt = timePtr(match.RoundStartedAt)
		response.RoundEndsAt = timePtr(match.RoundEndsAt)
		if match.CurrentQuestionIndex < len(match.QuestionIDs) {
			question, ok := h.content.GetQuizQuestion(match.QuestionIDs[match.CurrentQuestionIndex])
			if !ok {
				return MatchStateResponse{}, httpx.NewError(http.StatusInternalServerError, "match question is unavailable")
			}
			currentQuestion := questionResponse(question)
			response.CurrentQuestion = &currentQuestion
		}
	}

	currentQuestionID := ""
	if match.Status == mongomodel.MatchStatusActive && match.CurrentQuestionIndex < len(match.QuestionIDs) {
		currentQuestionID = match.QuestionIDs[match.CurrentQuestionIndex]
	}
	currentAnswers := answerByQuestionPlayer[currentQuestionID]
	for _, player := range match.Players {
		score := player.Score
		maxScore := maxScoreThroughCurrentQuestion(match, player)
		playerResponse := MatchPlayerResponse{
			PlayerID:  player.PlayerID,
			Nickname:  player.Nickname,
			Ready:     player.Ready,
			SitoneIDs: cloneStrings(player.SitoneIDs),
			Score:     &score,
			MaxScore:  &maxScore,
		}
		if match.Status == mongomodel.MatchStatusActive {
			_, playerResponse.AnsweredCurrentQuestion = currentAnswers[player.PlayerID]
		}
		if match.Status == mongomodel.MatchStatusCompleted {
			reward := openPowerReward(player.Score)
			playerResponse.OpenPowerReward = &reward
		}
		response.Players = append(response.Players, playerResponse)
	}

	if match.Status == mongomodel.MatchStatusCompleted {
		results, err := h.matchResults(match, answerByQuestionPlayer)
		if err != nil {
			return MatchStateResponse{}, err
		}
		response.Results = results
	}

	return response, nil
}

func (h *Handler) matchResults(
	match mongomodel.Match,
	answerByQuestionPlayer map[string]map[string]mongomodel.MatchAnswer,
) ([]MatchQuestionResult, error) {
	results := make([]MatchQuestionResult, 0, len(match.QuestionIDs))
	for _, questionID := range match.QuestionIDs {
		question, ok := h.content.GetQuizQuestion(questionID)
		if !ok {
			return nil, httpx.NewError(http.StatusInternalServerError, "match question is unavailable")
		}

		result := MatchQuestionResult{
			QuestionID:    question.ID,
			Prompt:        question.Prompt,
			ChoiceA:       question.ChoiceA,
			ChoiceB:       question.ChoiceB,
			ChoiceC:       question.ChoiceC,
			ChoiceD:       question.ChoiceD,
			CorrectChoice: question.CorrectChoice,
			Explanation:   question.Explanation,
			Answers:       make([]MatchAnswerResponse, 0, len(match.Players)),
		}
		for _, player := range match.Players {
			answer, ok := answerByQuestionPlayer[questionID][player.PlayerID]
			answerResponse := MatchAnswerResponse{
				PlayerID: player.PlayerID,
			}
			if ok {
				answeredAt := answer.AnsweredAt
				answerResponse.Choice = answer.Choice
				answerResponse.Correct = answer.Correct
				answerResponse.Score = answer.Score
				answerResponse.ElapsedMillis = answer.ElapsedMillis
				answerResponse.AnsweredAt = &answeredAt
			}
			result.Answers = append(result.Answers, answerResponse)
		}
		results = append(results, result)
	}
	return results, nil
}

func (h *Handler) publishState(ctx context.Context, match mongomodel.Match, eventName string) {
	if h.broker == nil {
		return
	}
	state, err := h.buildMatchState(ctx, match)
	if err != nil {
		return
	}
	h.broker.Publish(match.ID, Event{
		Name: eventName,
		Data: state,
	})
}

func writeMatchProblem(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, mongo.ErrNoDocuments) {
		httpx.WriteProblem(w, r, httpx.NotFound("match not found"))
		return
	}
	httpx.WriteProblem(w, r, err)
}
