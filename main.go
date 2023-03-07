package main

import (
	"github.com/hashicorp/vault/sdk/physical"

	physAerospike "github.com/hashicorp/vault/physical/aerospike"
	physAliCloudOSS "github.com/hashicorp/vault/physical/alicloudoss"
	physAzure "github.com/hashicorp/vault/physical/azure"
	physCassandra "github.com/hashicorp/vault/physical/cassandra"
	physCockroachDB "github.com/hashicorp/vault/physical/cockroachdb"
	physConsul "github.com/hashicorp/vault/physical/consul"
	physCouchDB "github.com/hashicorp/vault/physical/couchdb"
	physDynamoDB "github.com/hashicorp/vault/physical/dynamodb"
	physEtcd "github.com/hashicorp/vault/physical/etcd"
	physFoundationDB "github.com/hashicorp/vault/physical/foundationdb"
	physGCS "github.com/hashicorp/vault/physical/gcs"
	physManta "github.com/hashicorp/vault/physical/manta"
	physMSSQL "github.com/hashicorp/vault/physical/mssql"
	physMySQL "github.com/hashicorp/vault/physical/mysql"
	physOCI "github.com/hashicorp/vault/physical/oci"
	physPostgreSQL "github.com/hashicorp/vault/physical/postgresql"
	physRaft "github.com/hashicorp/vault/physical/raft"
	physS3 "github.com/hashicorp/vault/physical/s3"
	physSpanner "github.com/hashicorp/vault/physical/spanner"
	physSwift "github.com/hashicorp/vault/physical/swift"
	physZooKeeper "github.com/hashicorp/vault/physical/zookeeper"
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
