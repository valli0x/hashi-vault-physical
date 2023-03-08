package hashiphysical

import (
	physical "github.com/hashicorp/vault/sdk/physical"
	physAerospike "github.com/valli0x/hashi-vault-physical/aerospike"
	physAliCloudOSS "github.com/valli0x/hashi-vault-physical/alicloudoss"
	physAzure "github.com/valli0x/hashi-vault-physical/azure"
	physCassandra "github.com/valli0x/hashi-vault-physical/cassandra"
	physCockroachDB "github.com/valli0x/hashi-vault-physical/cockroachdb"
	physConsul "github.com/valli0x/hashi-vault-physical/consul"
	physCouchDB "github.com/valli0x/hashi-vault-physical/couchdb"
	physDynamoDB "github.com/valli0x/hashi-vault-physical/dynamodb"
	physEtcd "github.com/valli0x/hashi-vault-physical/etcd"
	physFoundationDB "github.com/valli0x/hashi-vault-physical/foundationdb"
	physGCS "github.com/valli0x/hashi-vault-physical/gcs"
	physManta "github.com/valli0x/hashi-vault-physical/manta"
	physMSSQL "github.com/valli0x/hashi-vault-physical/mssql"
	physMySQL "github.com/valli0x/hashi-vault-physical/mysql"
	physOCI "github.com/valli0x/hashi-vault-physical/oci"
	physPostgreSQL "github.com/valli0x/hashi-vault-physical/postgresql"
	physRaft "github.com/valli0x/hashi-vault-physical/hashi-raft"
	physS3 "github.com/valli0x/hashi-vault-physical/s3"
	physSpanner "github.com/valli0x/hashi-vault-physical/spanner"
	physSwift "github.com/valli0x/hashi-vault-physical/swift"
	physZooKeeper "github.com/valli0x/hashi-vault-physical/zookeeper"
	physFile "github.com/hashicorp/vault/sdk/physical/file"
	physInmem "github.com/hashicorp/vault/sdk/physical/inmem"
)

var (
	PhysicalBackends = map[string]physical.Factory{
		"aerospike":              physAerospike.NewAerospikeBackend,
		"alicloudoss":            physAliCloudOSS.NewAliCloudOSSBackend,
		"azure":                  physAzure.NewAzureBackend,
		"cassandra":              physCassandra.NewCassandraBackend,
		"cockroachdb":            physCockroachDB.NewCockroachDBBackend,
		"consul":                 physConsul.NewConsulBackend,
		"couchdb_transactional":  physCouchDB.NewTransactionalCouchDBBackend,
		"couchdb":                physCouchDB.NewCouchDBBackend,
		"dynamodb":               physDynamoDB.NewDynamoDBBackend,
		"etcd":                   physEtcd.NewEtcdBackend,
		"file_transactional":     physFile.NewTransactionalFileBackend,
		"file":                   physFile.NewFileBackend,
		"foundationdb":           physFoundationDB.NewFDBBackend,
		"gcs":                    physGCS.NewBackend,
		"inmem_ha":               physInmem.NewInmemHA,
		"inmem_transactional_ha": physInmem.NewTransactionalInmemHA,
		"inmem_transactional":    physInmem.NewTransactionalInmem,
		"inmem":                  physInmem.NewInmem,
		"manta":                  physManta.NewMantaBackend,
		"mssql":                  physMSSQL.NewMSSQLBackend,
		"mysql":                  physMySQL.NewMySQLBackend,
		"oci":                    physOCI.NewBackend,
		"postgresql":             physPostgreSQL.NewPostgreSQLBackend,
		"s3":                     physS3.NewS3Backend,
		"spanner":                physSpanner.NewBackend,
		"swift":                  physSwift.NewSwiftBackend,
		"raft":                   physRaft.NewRaftBackend,
		"zookeeper":              physZooKeeper.NewZooKeeperBackend,
	}
)
