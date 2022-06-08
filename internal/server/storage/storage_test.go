package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemoryRepoRW(t *testing.T) {
	memoryRepo := NewMemoryRepo()

	usernameExpect := "Efim"
	memoryRepo.Write("username", usernameExpect)
	usernameReal, err := memoryRepo.Read("username")

	require.NoError(t, err)
	require.Equal(t, usernameExpect, usernameReal)
}

func TestMemoryRepoReadEmpty(t *testing.T) {
	memoryRepo := NewMemoryRepo()

	_, err := memoryRepo.Read("username")

	require.Error(t, err)
}

func TestUpdateCounterValue(t *testing.T) {
	memStatsStorage := NewMemStatsMemoryRepo()

	err := memStatsStorage.UpdateCounterValue("PollCount", 7)
	require.NoError(t, err)
	err = memStatsStorage.UpdateCounterValue("PollCount", 22)
	require.NoError(t, err)
	PollCount, err := memStatsStorage.ReadValue("PollCount")
	require.NoError(t, err)

	require.Equal(t, "29", PollCount)
}
