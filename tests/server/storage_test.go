package server

import (
	"github.com/stretchr/testify/require"
	"metrics/internal/server/storage"
	"testing"
)

func TestMemoryRepoRW(t *testing.T) {
	memoryRepo := storage.NewMemoryRepo()

	usernameExpect := "Efim"
	memoryRepo.Write("username", usernameExpect)
	usernameReal, err := memoryRepo.Read("username")

	require.NoError(t, err)
	require.Equal(t, usernameExpect, usernameReal)
}

func TestMemoryRepoReadEmpty(t *testing.T) {
	memoryRepo := storage.NewMemoryRepo()

	_, err := memoryRepo.Read("username")

	require.Error(t, err)
}

func TestUpdateCounterValue(t *testing.T) {
	memStatsStorage := storage.NewMemStatsMemoryRepo()

	err := memStatsStorage.UpdateCounterValue("PollCount", 7)
	require.NoError(t, err)
	err = memStatsStorage.UpdateCounterValue("PollCount", 22)
	require.NoError(t, err)
	PollCount, err := memStatsStorage.ReadValue("PollCount")
	require.NoError(t, err)

	require.Equal(t, "29", PollCount)
}
