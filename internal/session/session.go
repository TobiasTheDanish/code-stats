package session

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"tobiasthedanish/code-stats/internal/database"
	"tobiasthedanish/code-stats/internal/kvs"
	"tobiasthedanish/code-stats/internal/viewmodel"

	"go.mongodb.org/mongo-driver/mongo"
)

type CodingSessions interface {
	ToViewModel() viewmodel.CodingSessions
}

// Period represents the time period for which the coding sessions have been aggregated.
type Period int8

const (
	Day Period = iota
	Week
	Month
	Year
)

type codingSessions []CodingSession

func (c codingSessions) ToViewModel() viewmodel.CodingSessions {
	return viewmodel.CodingSessions{
		TimeSpentData: c.TimeSpentChartData(),
		LanguageData:  c.LanguageChartData(),
	}
}

func (c codingSessions) TimeSpentChartData() viewmodel.ChartData {
	timeMap := make(map[string]float64)

	for _, s := range c {
		val, ok := timeMap[s.DateString]
		if !ok {
			timeMap[s.DateString] = float64(s.TotalTimeMs) / (1000 * 60 * 60)
		} else {
			timeMap[s.DateString] = val + float64(s.TotalTimeMs)/(1000*60*60)
		}
	}
	pairs := kvs.KeySortedPairs(timeMap)

	for k, v := range timeMap {
		pairs.Append(kvs.Pair[string, float64]{Key: k, Val: v})
	}
	sort.Sort(pairs)

	return viewmodel.ChartData{
		Labels: pairs.Keys(),
		Datasets: []viewmodel.Dataset{
			{
				Label: "Hours Spent",
				Data:  pairs.Values(),
			},
		},
	}
}

func (c codingSessions) LanguageChartData() viewmodel.ChartData {
	langMap := make(map[string]float64)

	for _, s := range c {
		for _, repo := range s.Repositories {
			for _, file := range repo.Files {
				val, ok := langMap[file.Filetype]
				if !ok {
					langMap[file.Filetype] = float64(file.DurationMs) / (1000 * 60 * 60)
				} else {
					langMap[file.Filetype] = val + float64(file.DurationMs)/(1000*60*60)
				}
			}
		}
	}

	pairs := kvs.ValueSortedPairs(langMap)

	pairs = pairs.Filter(func(p kvs.Pair[string, float64], i int) bool { return p.Val < 1.0 })
	sort.Sort(pairs)

	return viewmodel.ChartData{
		Labels: pairs.Keys(),
		Datasets: []viewmodel.Dataset{
			{
				Label: "Hours spent",
				Data:  pairs.Values(),
			},
		},
	}
}

// CodingSession represents a coding session that has been aggregated
// for a given time period (day, week, month, year).
type CodingSession struct {
	ID           string       `bson:"_id,omitempty"`
	Period       Period       `bson:"period"`
	EpochDateMs  int64        `bson:"date"`
	DateString   string       `bson:"date_string"`
	TotalTimeMs  int64        `bson:"total_time_ms"`
	Repositories []Repository `bson:"repositories"`
}

type Repository struct {
	Name       string `bson:"name"`
	Files      []File `bson:"files"`
	DurationMs int64  `bson:"duration_ms"`
}

type File struct {
	Name       string `bson:"name"`
	Path       string `bson:"path"`
	Filetype   string `bson:"filetype"`
	DurationMs int64  `bson:"duration_ms"`
}

type SessionStore interface {
	ForPeriod(Period) (CodingSessions, error)
}

type sessionStore struct {
	dbService database.Service
}

func NewStore() SessionStore {
	return &sessionStore{
		dbService: database.New(),
	}
}

func (s *sessionStore) ForPeriod(period Period) (CodingSessions, error) {
	cursor, err := s.GetCursorForPeriod(period)
	if err != nil {
		return nil, err
	}

	return unmarshallCodingSessions(cursor)
}

func (s *sessionStore) GetCursorForPeriod(period Period) (*mongo.Cursor, error) {
	switch period {
	case Day:
		return s.dbService.GetAllDaily()
	case Week:
		return s.dbService.GetAllWeekly()
	case Month:
		return s.dbService.GetAllMonthly()
	case Year:
		return s.dbService.GetAllYearly()

	default:
		return nil, errors.New(fmt.Sprintf("Invalid period: %d", period))
	}
}

func unmarshallCodingSessions(cursor *mongo.Cursor) (CodingSessions, error) {
	result := make(codingSessions, 0, 0)
	for cursor.Next(context.TODO()) {
		var res CodingSession
		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}

		result = append(result, res)
	}

	return result, nil
}
