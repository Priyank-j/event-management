package conf

import "strings"

/*
Redis Configurations
*/

func getRedisAddr() string {
	if ENV == ENV_PROD {
		return getSSMKey("REDIS_MASTER_ADDRESS_PROD")
	} else if strings.HasPrefix(ENV, ENV_UAT) {
		return getSSMKey("REDIS_MASTER_ADDRESS_UAT")
	} else if strings.HasPrefix(ENV, ENV_DEV) {
		return "redis-server:6379"
	}

	if IsRunningInDockerContainer() {
		return "host.docker.internal:6379"
	}

	return "0.0.0.0:6379"
}

func getRedisReplicaAddr() string {
	if ENV == ENV_PROD {
		return getSSMKey("REDIS_REPLICA_ADDRESS_PROD")
	} else if ENV == ENV_UAT_TDL {
		return getSSMKey("REDIS_REPLICA_ADDRESS_UAT")
	}

	return getRedisAddr()
}

var RedisConf = map[string]interface{}{
	"Addr":        getRedisAddr(),
	"ReplicaAddr": getRedisReplicaAddr(),
	"SSL":         getSSLInfo(),
}

func getSSLInfo() bool {
	switch ENV {
	case ENV_PROD:
		return true
	case ENV_UAT_TDL:
		return true
	default:
		return false
	}
}
