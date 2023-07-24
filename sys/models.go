package sys

import "time"

type HealthModel struct {
	Type                       string
	Initialized                bool
	Sealed                     bool
	T                          int
	N                          int
	Progress                   int
	Nonce                      string
	Version                    string
	BuildDate                  string
	Migration                  bool
	ClusterName                string
	ClusterId                  string
	RecoverySeal               bool
	StorageType                string
	Standby                    bool
	PerformanceStandby         bool
	ReplicationPerformanceMode string
	ReplicationDrMode          string
	ServerTimeUTC              int64
}

type LeaderStatusModel struct {
	HAEnabled                       bool
	ISSelf                          bool
	ActiveTime                      string
	LeaderAddress                   string
	LeaderClusterAddress            string
	PerformanceStandby              bool
	PerformanceStandbyLastRemoteWAL int
	RAFTCommittedIndex              int
	RAFTAppliedIndex                int
}

func GetSealStatus() (HealthModel, error) {
	return HealthModel{
		Type:        "shamir",
		Initialized: true,
		Sealed:      false,
		T:           3,
		N:           5,
		Progress:    0,
		Nonce:       "",
		// TODO: read following values from the app
		Version:     "1.14.0",
		BuildDate:   "2023-06-19T11:40:23Z",
		Migration:   false,
		ClusterName: "vau-prod-01",
		// TODO: generate UUID on the first startup
		ClusterId:    "64f762ad-3841-a6e2-5165-7803cd169d6c",
		RecoverySeal: false,
		StorageType:  "raft",
	}, nil
}

func GetHealthStatus() (HealthModel, error) {
	return HealthModel{
		Initialized:                true,
		Sealed:                     false,
		Standby:                    true,
		PerformanceStandby:         false,
		ReplicationPerformanceMode: "disabled",
		ReplicationDrMode:          "disabled",
		ServerTimeUTC:              time.Now().Unix(),
		Version:                    "1.14.0",
		ClusterName:                "vau-prod-01",
		ClusterId:                  "64f762ad-3841-a6e2-5165-7803cd169d6c",
	}, nil
}

func GetLeaderStatus() (LeaderStatusModel, error) {
	return LeaderStatusModel{
		HAEnabled:                       true,
		ISSelf:                          false,
		ActiveTime:                      "0001-01-01T00:00:00Z",
		LeaderAddress:                   "https://127.0.0.1:8200",
		LeaderClusterAddress:            "https://127.0.0.1:8201",
		PerformanceStandby:              false,
		PerformanceStandbyLastRemoteWAL: 0,
		RAFTCommittedIndex:              108,
		RAFTAppliedIndex:                108,
	}, nil
}
