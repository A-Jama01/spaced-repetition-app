package scheduler

import (
	"errors"
	"math"
	"time"

	"github.com/A-Jama01/spaced-repetition-app/internal/store"
)

const (
	MaxDaysAllowed float64 = 99999.0
)

var (
	ErrMaxDays = errors.New("Exceeded max schedulable days")
	w = [21]float64{0.212, 1.2931, 2.3065, 8.2956, 6.4133, 0.8334, 3.0194, 0.001, 1.8722, 0.1666, 
		0.796, 1.4835, 0.0614, 0.2629, 1.6483, 0.6014, 1.8729, 0.5425, 0.0912, 0.0658, 0.1542}
	factor = math.Pow(0.9, -1.0 / w[20]) - 1
)

type T float64
type R float64
type S float64
type D float64
type Grade int64

const (
	Forgot Grade = iota + 1
	Hard 
	Good 
	Easy 
)

func retrievability(t T, s S) R {
	return R(math.Pow((1 + factor * float64(t) / float64(s)), -w[20]))
}

func interval(r R, s S) T {
	return T(float64(s) / factor * (math.Pow(float64(r), 1 / -w[20]) - 1))
}

func initialStability(g Grade) S {
	return S(w[int(g) - 1])
}

func stabilityFailure(s S, d D, r R) S {
	fd := math.Pow(float64(d), -w[12])
	fs := (math.Pow(float64(s) + 1, w[13])) - 1
	fr := math.Exp(w[14] * (1 - float64(r)))
	fs = w[11] * fd * fs * fr
	return S(math.Min(fs, float64(s)))
}

func stabilitySuccess(s S, d D, r R, g Grade) S {
	fd := 11 - float64(d)
	fs := math.Pow(float64(s), -w[9])
	fr := math.Exp(w[10] * (1 - float64(r))) - 1

	penalty, bonus := 1.0, 1.0
	if g == Hard {
		penalty = w[15]
	}
	if g == Easy {
		bonus = w[16]
	}
	sInc := 1 + bonus * penalty * math.Exp(w[8]) * fd * fs * fr

	return S(float64(s) * sInc)
}

func stability(s S, d D, r R, g Grade) S {
	if g == Forgot {
		return stabilityFailure(s, d, r)
	}

	return stabilitySuccess(s, d, r, g)
}

func stabilitySameDayReview(s S, g Grade) S {
	fg := math.Exp(w[17] * (float64(g) - 3 + w[18]))
	return S(float64(s) * fg * math.Pow(float64(s), -w[19]))
}

func clampDifficulty(d D) D {
	if float64(d) < 1 {
		return D(1.0)
	}
	if float64(d) > 10 {
		return D(10.0)
	}

	return d
}

func initialDifficulty(g Grade) D {
	difficulty := w[4] - math.Exp(w[5] * (float64(g) - 1)) + 1
	return clampDifficulty(D(difficulty))
}

func difficulty(d D, g Grade) D {
	deltaD := -w[6] * (float64(g) - 3)
	dUpdated := float64(d) + deltaD * ((10 - float64(d)) / 9)
	target := float64(initialDifficulty(Good))
	difficulty := w[7] * target + (1 - w[7]) * dUpdated
	return clampDifficulty(D(difficulty))
}

func nextDueDate(desiredRetention R, stability S) (time.Time, error) {
	interval := float64(interval(desiredRetention, stability))
	if interval > MaxDaysAllowed {
		return time.Time{}, ErrMaxDays
	}
	if interval < 0 {
		return time.Time{}, ErrMaxDays
	}

	days := time.Duration(interval * 24 * float64(time.Hour))
	due := time.Now().Add(days)

	return due, nil
}

func ScheduleCard(card *store.Card, gradeInput int64) error {
	desiredRetention := R(0.9)
	grade := Grade(gradeInput) 

	//Initial Review
	if card.Due == nil {
		initStability := initialStability(grade)	
		initDifficulty := initialDifficulty(grade)
		due, err := nextDueDate(desiredRetention, initStability)
		if err != nil {
			return err
		}

		card.Stability = float64(initStability)
		card.Difficulty = float64(initDifficulty)
		card.Due = &due
	} else { //N-th Review
		elapsedTime := T(time.Since(*card.LastReview).Hours() / 24)	
		prevStability := S(card.Stability)
		prevDifficulty := D(card.Difficulty)
		newRetrievability := retrievability(elapsedTime, prevStability)	

		var newStability S
		oneDay := T(1.0)
		if elapsedTime < oneDay {
			newStability = stabilitySameDayReview(prevStability, grade)
		} else {
			newStability = stability(prevStability, prevDifficulty, newRetrievability, grade)
		}
		newDifficulty := difficulty(prevDifficulty, grade)
		due, err := nextDueDate(desiredRetention, newStability)
		if err != nil {
			return err
		}

		card.Retrievability = float64(newRetrievability)
		card.Stability = float64(newStability)
		card.Difficulty = float64(newDifficulty)
		card.Due = &due
	}
	
	currentTime := time.Now()
	card.LastReview = &currentTime

	return nil
}
