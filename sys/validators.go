// Actually everything described in this file isn't used right now
// We're not saving anything from this package to the DB
package sys

import (
	"github.com/gin-gonic/gin"
	"github.com/miknikif/vault-auto-unseal/common"
)

type SealStatusModelValidator struct {
	Health struct {
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
	} `json:"key"`
	healthModel HealthModel `json:"-"`
}

func (s *SealStatusModelValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, s)
	if err != nil {
		return err
	}

	s.healthModel.Type = s.Health.Type
	s.healthModel.Initialized = s.Health.Initialized
	s.healthModel.Sealed = s.Health.Sealed
	s.healthModel.T = s.Health.T
	s.healthModel.N = s.Health.N
	s.healthModel.Progress = s.Health.Progress
	s.healthModel.Nonce = s.Health.Nonce
	s.healthModel.Version = s.Health.Version
	s.healthModel.BuildDate = s.Health.BuildDate
	s.healthModel.Migration = s.Health.Migration
	s.healthModel.ClusterName = s.Health.ClusterName
	s.healthModel.ClusterId = s.Health.ClusterId
	s.healthModel.RecoverySeal = s.Health.RecoverySeal
	s.healthModel.StorageType = s.Health.StorageType
	return nil
}

func NewSealStatusModelValidator() SealStatusModelValidator {
	return SealStatusModelValidator{}
}

func NewSealStatusModelValidatorFillWith(healthModel HealthModel) SealStatusModelValidator {
	sealStatusModelValidator := NewSealStatusModelValidator()

	sealStatusModelValidator.Health.Type = healthModel.Type
	sealStatusModelValidator.Health.Initialized = healthModel.Initialized
	sealStatusModelValidator.Health.Sealed = healthModel.Sealed
	sealStatusModelValidator.Health.T = healthModel.T
	sealStatusModelValidator.Health.N = healthModel.N
	sealStatusModelValidator.Health.Progress = healthModel.Progress
	sealStatusModelValidator.Health.Nonce = healthModel.Nonce
	sealStatusModelValidator.Health.Version = healthModel.Version
	sealStatusModelValidator.Health.BuildDate = healthModel.BuildDate
	sealStatusModelValidator.Health.Migration = healthModel.Migration
	sealStatusModelValidator.Health.ClusterName = healthModel.ClusterName
	sealStatusModelValidator.Health.ClusterId = healthModel.ClusterId
	sealStatusModelValidator.Health.RecoverySeal = healthModel.RecoverySeal
	sealStatusModelValidator.Health.StorageType = healthModel.StorageType
	return sealStatusModelValidator
}

type HealthModelValidator struct {
	Health struct {
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
	} `json:"key"`
	healthModel HealthModel `json:"-"`
}

func (s *HealthModelValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, s)
	if err != nil {
		return err
	}
	s.healthModel.Initialized = s.Health.Initialized
	s.healthModel.Sealed = s.Health.Sealed
	s.healthModel.Standby = s.Health.Standby
	s.healthModel.PerformanceStandby = s.Health.PerformanceStandby
	s.healthModel.ReplicationPerformanceMode = s.Health.ReplicationPerformanceMode
	s.healthModel.ReplicationDrMode = s.Health.ReplicationDrMode
	s.healthModel.ServerTimeUTC = s.Health.ServerTimeUTC
	s.healthModel.Version = s.Health.Version
	s.healthModel.ClusterName = s.Health.ClusterName
	s.healthModel.ClusterId = s.Health.ClusterId
	return nil
}

func NewHealthModelValidator() HealthModelValidator {
	return HealthModelValidator{}
}

func NewHealthModelValidatorFillWith(healthModel HealthModel) HealthModelValidator {
	healthModelValidator := NewHealthModelValidator()
	healthModelValidator.Health.Initialized = healthModel.Initialized
	healthModelValidator.Health.Sealed = healthModel.Sealed
	healthModelValidator.Health.Standby = healthModel.Standby
	healthModelValidator.Health.PerformanceStandby = healthModel.PerformanceStandby
	healthModelValidator.Health.ReplicationPerformanceMode = healthModel.ReplicationPerformanceMode
	healthModelValidator.Health.ReplicationDrMode = healthModel.ReplicationDrMode
	healthModelValidator.Health.ServerTimeUTC = healthModel.ServerTimeUTC
	healthModelValidator.Health.Version = healthModel.Version
	healthModelValidator.Health.ClusterName = healthModel.ClusterName
	healthModelValidator.Health.ClusterId = healthModel.ClusterId
	return healthModelValidator
}

type LeaderStatusModelValidator struct {
	LeaderStatus struct {
		HAEnabled                       bool   `json:"ha_enabled"`
		ISSelf                          bool   `json:"is_self"`
		ActiveTime                      string `json:"active_time"`
		LeaderAddress                   string `json:"leader_address"`
		LeaderClusterAddress            string `json:"leader_cluster_address"`
		PerformanceStandby              bool   `json:"performance_standby"`
		PerformanceStandbyLastRemoteWAL int    `json:"performance_standby_last_remote_wal"`
		RAFTCommittedIndex              int    `json:"raft_committed_index"`
		RAFTAppliedIndex                int    `json:"raft_applied_index"`
	} `json:"key"`
	leaderStatusModel LeaderStatusModel `json:"-"`
}

func (s *LeaderStatusModelValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, s)
	if err != nil {
		return err
	}
	s.leaderStatusModel.HAEnabled = s.LeaderStatus.HAEnabled
	s.leaderStatusModel.ISSelf = s.LeaderStatus.ISSelf
	s.leaderStatusModel.ActiveTime = s.LeaderStatus.ActiveTime
	s.leaderStatusModel.LeaderAddress = s.LeaderStatus.LeaderAddress
	s.leaderStatusModel.LeaderClusterAddress = s.LeaderStatus.LeaderClusterAddress
	s.leaderStatusModel.PerformanceStandby = s.LeaderStatus.PerformanceStandby
	s.leaderStatusModel.PerformanceStandbyLastRemoteWAL = s.LeaderStatus.PerformanceStandbyLastRemoteWAL
	s.leaderStatusModel.RAFTCommittedIndex = s.LeaderStatus.RAFTCommittedIndex
	s.leaderStatusModel.RAFTAppliedIndex = s.LeaderStatus.RAFTAppliedIndex
	return nil
}

func NewLeaderStatusModelValidator() LeaderStatusModelValidator {
	return LeaderStatusModelValidator{}
}

func NewLeaderStatusModelValidatorFillWith(leaderStatusModel LeaderStatusModel) LeaderStatusModelValidator {
	leaderStatusModelValidator := NewLeaderStatusModelValidator()
	leaderStatusModelValidator.LeaderStatus.HAEnabled = leaderStatusModel.HAEnabled
	leaderStatusModelValidator.LeaderStatus.ISSelf = leaderStatusModel.ISSelf
	leaderStatusModelValidator.LeaderStatus.ActiveTime = leaderStatusModel.ActiveTime
	leaderStatusModelValidator.LeaderStatus.LeaderAddress = leaderStatusModel.LeaderAddress
	leaderStatusModelValidator.LeaderStatus.LeaderClusterAddress = leaderStatusModel.LeaderClusterAddress
	leaderStatusModelValidator.LeaderStatus.PerformanceStandby = leaderStatusModel.PerformanceStandby
	leaderStatusModelValidator.LeaderStatus.PerformanceStandbyLastRemoteWAL = leaderStatusModel.PerformanceStandbyLastRemoteWAL
	leaderStatusModelValidator.LeaderStatus.RAFTCommittedIndex = leaderStatusModel.RAFTCommittedIndex
	leaderStatusModelValidator.LeaderStatus.RAFTAppliedIndex = leaderStatusModel.RAFTAppliedIndex
	return leaderStatusModelValidator
}
