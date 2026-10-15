package logz

import (
	"errors"
	"testing"

	"nova/errorn"
)

func Test_LogErrorN(t *testing.T) {
	err := errors.New("err msg")
	errN0 := errorn.NewRedis(err)
	errN1 := errorn.NewRedis(err).
		SetMetadata("mt1", "query Str").
		SetMetadata("mt2", "st Str")
	errN2 := errorn.NewSls(errN1).SetMetadata("mt3", "neStr").WithStackS()

	Error().Err(errN0).Msg("0")
	Error().Err(errN1).Msg("1")
	Error().Err(errN2).Msg("2")
	Error().Stack().Err(errN2).Msg("2 with stack")

	Stack(Info()).Msg("only stack")

	//logz.Info().Msgf("%v,%v", errorn.IsCode(errN2, coden.DbSlsError), errorn.IsCode(errN2, coden.DbRedisError))
}
