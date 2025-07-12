package internal

import (
	"math"
)

var w = [21]float64{0.2172, 1.1771, 3.2602, 16.1507, 7.0114, 0.57, 2.0966, 0.0069, 1.5261, 0.112, 1.0178, 1.849, 
	0.1133, 0.3127, 2.2934, 0.2191, 3.0004, 0.7536, 0.3332, 0.1437, 0.2}

type T float64
type R float64
type S float64
type D float64
type Grade int64

const (
	Again Grade = iota + 1
	Hard 
	Good 
	Easy 
)

var factor = math.Pow(0.9, (-1.0 / w[20] - 1))

func Retrievability(t T, s S) R {
	return R(math.Pow((1 + factor * float64(t) / float64(s)), -w[20]))
}

func Interval(r R, s S) T {
	return T(float64(s) / factor * (math.Pow(float64(r), 1 / -w[20]) - 1))
}

func InitialStability(g Grade) S {
	return S(w[int(g) - 1])
}

func StabilityFailure(s S, d D, r R) S {
	fd := math.Pow(float64(d), -w[12])
	fs := (math.Pow(float64(s) + 1, w[13])) - 1
	fr := math.Exp(w[14] * (1 - float64(r)))
	fs = w[11] * fd * fs * fr
	return S(math.Min(fs, float64(s)))
}

func StabilitySuccess(s S, d D, r R, g Grade) S {
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

func StabilitySameDayReview(s S, g Grade) S {
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

func InitialDifficulty(g Grade) D {
	difficulty := w[4] - math.Exp(w[5] * (float64(g) - 1)) + 1
	return clampDifficulty(D(difficulty))
}

func Difficulty(d D, g Grade) D {
	deltaD := -w[6] * (float64(g) - 3)
	dUpdated := float64(d) + deltaD * ((10 - float64(d)) / 9)
	target := float64(InitialDifficulty(Good))
	difficulty := w[7] * target + (1 - w[7]) * dUpdated
	return clampDifficulty(D(difficulty))
}
