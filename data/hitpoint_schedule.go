package data

import (
	"encoding/json"
	"errors"
	"fmt"
)

type HitpointSchedule struct {
	GracePeriodInMinute uint
	IntervalInMinute    *uint
	CronSchedule        *string
}

type DurationType string

const (
	DT_Minutes DurationType = "minutes"
	DT_Hours   DurationType = "hours"
	DT_Days    DurationType = "days"
)

type Duration struct {
	DType    DurationType
	Duration uint
}

func (d *Duration) InMinutes() (uint, error) {
	if d.DType == DT_Minutes {
		return d.Duration, nil
	}
	if d.DType == DT_Hours {
		return d.Duration * 60, nil
	} else if d.DType == DT_Days {
		return d.Duration * 86400, nil
	}
	return 0, fmt.Errorf("illegal duration type found %s", d.DType)
}

func (s *HitpointSchedule) Json() (map[string]interface{}, error) {
	b, err := json.Marshal(*s)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var res map[string]interface{}
	err = json.Unmarshal(b, &res)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return res, nil
}

func (s *HitpointSchedule) JsonString() (*string, error) {
	b, err := json.Marshal(*s)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	js := string(b)
	return &js, nil
}

func CreatePeriodicSchedule(period Duration, grace Duration) (*HitpointSchedule, error) {
	period_min, p_err := period.InMinutes()
	grace_min, g_err := period.InMinutes()
	if g_err != nil {
		return nil, g_err
	}
	if p_err != nil {
		return nil, p_err
	}
	if period_min <= 0 {
		return nil, errors.New("invalid value provdided for duration")
	}
	return &HitpointSchedule{GracePeriodInMinute: grace_min, IntervalInMinute: &period_min}, nil
}
