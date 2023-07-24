package sys

import (
	"github.com/gin-gonic/gin"
)

type SealStatusSerializer struct {
	C *gin.Context
	HealthModel
}

type SealStatusResponse struct {
	Type         string `json:"type"`
	Initialized  bool   `json:"initialized"`
	Sealed       bool   `json:"sealed"`
	T            int    `json:"t"`
	N            int    `json:"n"`
	Progress     int    `json:"progress"`
	Nonce        string `json:"nonce"`
	Version      string `json:"version"`
	BuildDate    string `json:"build_date"`
	Migration    bool   `json:"migration"`
	ClusterName  string `json:"cluster_name"`
	ClusterId    string `json:"cluster_id"`
	RecoverySeal bool   `json:"recovery_seal"`
	StorageType  string `json:"storage_type"`
}

func (s *SealStatusSerializer) Response() SealStatusResponse {
	response := SealStatusResponse{
		Type:         s.Type,
		Initialized:  s.Initialized,
		Sealed:       s.Sealed,
		T:            s.T,
		N:            s.N,
		Progress:     s.Progress,
		Nonce:        s.Nonce,
		Version:      s.Version,
		BuildDate:    s.BuildDate,
		Migration:    s.Migration,
		ClusterName:  s.ClusterName,
		ClusterId:    s.ClusterId,
		RecoverySeal: s.RecoverySeal,
		StorageType:  s.StorageType,
	}
	return response
}

type HealthSerializer struct {
	C *gin.Context
	HealthModel
}

type HealthResponse struct {
	Initialized                bool   `json:"initialized"`
	Sealed                     bool   `json:"sealed"`
	Standby                    bool   `json:"standby"`
	PerformanceStandby         bool   `json:"performance_standby"`
	ReplicationPerformanceMode string `json:"replication_performance_mode"`
	ReplicationDrMode          string `json:"replication_dr_mode"`
	ServerTimeUTC              int64  `json:"server_time_utc"`
	Version                    string `json:"version"`
	ClusterName                string `json:"cluster_name"`
	ClusterId                  string `json:"cluster_id"`
}

func (s *HealthSerializer) Response() HealthResponse {
	response := HealthResponse{
		Initialized:                s.Initialized,
		Sealed:                     s.Sealed,
		Standby:                    s.Standby,
		PerformanceStandby:         s.PerformanceStandby,
		ReplicationPerformanceMode: s.ReplicationPerformanceMode,
		ReplicationDrMode:          s.ReplicationDrMode,
		ServerTimeUTC:              s.ServerTimeUTC,
		Version:                    s.Version,
		ClusterName:                s.ClusterName,
		ClusterId:                  s.ClusterId,
	}
	return response
}

type LeaderStatusSerializer struct {
	C *gin.Context
	LeaderStatusModel
}

type LeaderStatusResponse struct {
	HAEnabled                       bool   `json:"ha_enabled"`
	ISSelf                          bool   `json:"is_self"`
	ActiveTime                      string `json:"active_time"`
	LeaderAddress                   string `json:"leader_address"`
	LeaderClusterAddress            string `json:"leader_cluster_address"`
	PerformanceStandby              bool   `json:"performance_standby"`
	PerformanceStandbyLastRemoteWAL int    `json:"performance_standby_last_remote_wal"`
	RAFTCommittedIndex              int    `json:"raft_committed_index"`
	RAFTAppliedIndex                int    `json:"raft_applied_index"`
}

func (s *LeaderStatusSerializer) Response() LeaderStatusResponse {
	response := LeaderStatusResponse{
		HAEnabled:                       s.HAEnabled,
		ISSelf:                          s.ISSelf,
		ActiveTime:                      s.ActiveTime,
		LeaderAddress:                   s.LeaderAddress,
		LeaderClusterAddress:            s.LeaderClusterAddress,
		PerformanceStandby:              s.PerformanceStandby,
		PerformanceStandbyLastRemoteWAL: s.PerformanceStandbyLastRemoteWAL,
		RAFTCommittedIndex:              s.RAFTCommittedIndex,
		RAFTAppliedIndex:                s.RAFTAppliedIndex,
	}
	return response
}
