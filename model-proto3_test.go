package model_convert

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestGoModelToProto3(t *testing.T) {
	type ArrayElement struct {
	}
	type U struct {
		Arr     []ArrayElement
		Config2 []byte
		Config  json.RawMessage

		Username string
		Password string
		Age      int
		Id       int32
	}
	ps, setM, setP := GoModelToProto3(U{}, map[string]string{
		"${pb_pkg_name}":    "userProto",
		"${model_pkg_name}": "userModel",
	})
	fmt.Println(ps)
	fmt.Println(setM)
	fmt.Println(setP)
}

func TestGoModelToProto2(t *testing.T) {
	type ChallengeGameUserProcess struct {
		Id                     int             `gorm:"column:id;default:" json:"id" form:"id"`
		GameId                 int             `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		UserId                 int             `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		GameAreaId             int             `gorm:"column:game_area_id;default:" json:"game_area_id" form:"game_area_id"`
		GameTypeKey            string          `gorm:"column:game_type_key;default:" json:"game_type_key" form:"game_type_key"`
		CurrentStage           int             `gorm:"column:current_stage;default:" json:"current_stage" form:"current_stage"`
		CurrentGameNum         int             `gorm:"column:current_game_num;default:" json:"current_game_num" form:"current_game_num"`
		CurrentScore           int             `gorm:"column:current_score;default:" json:"current_score" form:"current_score"`
		CurrentAccumulateAward int             `gorm:"column:current_accumulate_award;default:" json:"current_accumulate_award" form:"current_accumulate_award"`
		PassProps              json.RawMessage `gorm:"column:pass_props;default:" json:"pass_props" form:"pass_props"`
		EasterTimes            int             `gorm:"column:easter_times;default:" json:"easter_times" form:"easter_times"`
		ChallengeState         int             `gorm:"column:challenge_state;default:" json:"challenge_state" form:"challenge_state"`
		AwardState             int             `gorm:"column:award_state;default:" json:"award_state" form:"award_state"`
		EasterTime             time.Time       `gorm:"column:easter_time;default:" json:"easter_time" form:"easter_time"`
		DefaultStartTime       time.Time       `gorm:"column:default_start_time;default:" json:"default_start_time" form:"default_start_time"`
		EndTime                time.Time       `gorm:"column:end_time;default:" json:"end_time" form:"end_time"`
		CreatedAt              time.Time       `gorm:"column:created_at;default:" json:"created_at" form:"created_at"`
		UpdatedAt              time.Time       `gorm:"column:updated_at;default:" json:"updated_at" form:"updated_at"`
	}
	ps, setM, setP := GoModelToProto2(ChallengeGameUserProcess{}, map[string]string{
		"${pb_pkg_name}":    "challengePb",
		"${model_pkg_name}": "challengeModel",
		//		"${start_index}": "1",
	})
	fmt.Println(ps)
	fmt.Println(setM)
	fmt.Println(setP)
}
