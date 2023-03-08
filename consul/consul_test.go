package consul

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/consul"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
)

func TestConsul_newConsulBackend(t *testing.T) {
	tests := []struct {
		name            string
		consulConfig    map[string]string
		fail            bool
		redirectAddr    string
		checkTimeout    time.Duration
		path            string
		service         string
		address         string
		scheme          string
		token           string
		max_parallel    int
		disableReg      bool
		consistencyMode string
	}{
		{
			name:            "Valid default config",
			consulConfig:    map[string]string{},
			checkTimeout:    5 * time.Second,
			redirectAddr:    "http://127.0.0.1:8200",
			path:            "vault/",
			service:         "vault",
			address:         "127.0.0.1:8500",
			scheme:          "http",
			token:           "",
			max_parallel:    4,
			disableReg:      false,
			consistencyMode: "default",
		},
		{
			name: "Valid modified config",
			consulConfig: map[string]string{
				"path":                 "seaTech/",
				"service":              "astronomy",
				"redirect_addr":        "http://127.0.0.2:8200",
				"check_timeout":        "6s",
				"address":              "127.0.0.2",
				"scheme":               "https",
				"token":                "deadbeef-cafeefac-deadc0de-feedface",
				"max_parallel":         "4",
				"disable_registration": "false",
				"consistency_mode":     "strong",
			},
			checkTimeout:    6 * time.Second,
			path:            "seaTech/",
			service:         "astronomy",
			redirectAddr:    "http://127.0.0.2:8200",
			address:         "127.0.0.2",
			scheme:          "https",
			token:           "deadbeef-cafeefac-deadc0de-feedface",
			max_parallel:    4,
			consistencyMode: "strong",
		},
		{
			name: "Unix socket",
			consulConfig: map[string]string{
				"address": "unix:///tmp/.consul.http.sock",
			},
			address: "/tmp/.consul.http.sock",
			scheme:  "http", // Default, not overridden?

			// Defaults
			checkTimeout:    5 * time.Second,
			redirectAddr:    "http://127.0.0.1:8200",
			path:            "vault/",
			service:         "vault",
			token:           "",
			max_parallel:    4,
			disableReg:      false,
			consistencyMode: "default",
		},
		{
			name: "Scheme in address",
			consulConfig: map[string]string{
				"address": "https://127.0.0.2:5000",
			},
			address: "127.0.0.2:5000",
			scheme:  "https",

			// Defaults
			checkTimeout:    5 * time.Second,
			redirectAddr:    "http://127.0.0.1:8200",
			path:            "vault/",
			service:         "vault",
			token:           "",
			max_parallel:    4,
			disableReg:      false,
			consistencyMode: "default",
		},
	}

	for _, test := range tests {
		logger := logging.NewVaultLogger(log.Debug)

		be, err := NewConsulBackend(test.consulConfig, logger)
		if test.fail {
			if err == nil {
				t.Fatalf(`Expected config "%s" to fail`, test.name)
			} else {
				continue
			}
		} else if !test.fail && err != nil {
			t.Fatalf("Expected config %s to not fail: %v", test.name, err)
		}

		c, ok := be.(*ConsulBackend)
		if !ok {
			t.Fatalf("Expected ConsulBackend: %s", test.name)
		}

		if test.path != c.path {
			t.Errorf("bad: %s %v != %v", test.name, test.path, c.path)
		}

		if test.consistencyMode != c.consistencyMode {
			t.Errorf("bad consistency_mode value: %v != %v", test.consistencyMode, c.consistencyMode)
		}

		// The configuration stored in the Consul "client" object is not exported, so
		// we either have to skip validating it, or add a method to export it, or use reflection.
		consulConfig := reflect.Indirect(reflect.ValueOf(c.client)).FieldByName("config")
		consulConfigScheme := consulConfig.FieldByName("Scheme").String()
		consulConfigAddress := consulConfig.FieldByName("Address").String()

		if test.scheme != consulConfigScheme {
			t.Errorf("bad scheme value: %v != %v", test.scheme, consulConfigScheme)
		}

		if test.address != consulConfigAddress {
			t.Errorf("bad address value: %v != %v", test.address, consulConfigAddress)
		}

		// FIXME(sean@): Unable to test max_parallel
		// if test.max_parallel != cap(c.permitPool) {
		// 	t.Errorf("bad: %v != %v", test.max_parallel, cap(c.permitPool))
		// }
	}
}

func TestConsulBackend(t *testing.T) {
	cleanup, config := consul.PrepareTestContainer(t, "1.4.4", false, true)
	defer cleanup()

	client, err := api.NewClient(config.APIConfig())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	randPath := fmt.Sprintf("vault-%d/", time.Now().Unix())
	defer func() {
		client.KV().DeleteTree(randPath, nil)
	}()

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewConsulBackend(map[string]string{
		"address":      config.Address(),
		"token":        config.Token,
		"path":         randPath,
		"max_parallel": "256",
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

func TestConsul_TooLarge(t *testing.T) {
	cleanup, config := consul.PrepareTestContainer(t, "1.4.4", false, true)
	defer cleanup()

	client, err := api.NewClient(config.APIConfig())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	randPath := fmt.Sprintf("vault-%d/", time.Now().Unix())
	defer func() {
		client.KV().DeleteTree(randPath, nil)
	}()

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewConsulBackend(map[string]string{
		"address":      config.Address(),
		"token":        config.Token,
		"path":         randPath,
		"max_parallel": "256",
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	zeros := make([]byte, 600000)
	n, err := rand.Read(zeros)
	if n != 600000 {
		t.Fatalf("expected 500k zeros, read %d", n)
	}
	if err != nil {
		t.Fatal(err)
	}

	err = b.Put(context.Background(), &physical.Entry{
		Key:   "foo",
		Value: zeros,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), physical.ErrValueTooLarge) {
		t.Fatalf("expected value too large error, got %v", err)
	}

	err = b.(physical.Transactional).Transaction(context.Background(), []*physical.TxnEntry{
		{
			Operation: physical.PutOperation,
			Entry: &physical.Entry{
				Key:   "foo",
				Value: zeros,
			},
		},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), physical.ErrValueTooLarge) {
		t.Fatalf("expected value too large error, got %v", err)
	}
}

func TestConsul_TransactionalBackend_GetTransactionsForNonExistentValues(t *testing.T) {
	// TODO: unskip this after Consul releases 1.14 and we update our API dep. It currently fails but should pass with Consul 1.14
	t.SkipNow()

	cleanup, config := consul.PrepareTestContainer(t, "1.4.4", false, true)
	defer cleanup()

	client, err := api.NewClient(config.APIConfig())
	if err != nil {
		t.Fatal(err)
	}

	txns := make([]*physical.TxnEntry, 0)
	ctx := context.Background()
	logger := logging.NewVaultLogger(log.Debug)
	backendConfig := map[string]string{
		"address":      config.Address(),
		"token":        config.Token,
		"path":         "vault/",
		"max_parallel": "-1",
	}

	be, err := NewConsulBackend(backendConfig, logger)
	if err != nil {
		t.Fatal(err)
	}
	b := be.(*ConsulBackend)

	defer func() {
		_, _ = client.KV().DeleteTree("foo/", nil)
	}()

	txns = append(txns, &physical.TxnEntry{
		Operation: physical.GetOperation,
		Entry: &physical.Entry{
			Key: "foo/bar",
		},
	})
	txns = append(txns, &physical.TxnEntry{
		Operation: physical.PutOperation,
		Entry: &physical.Entry{
			Key:   "foo/bar",
			Value: []byte("baz"),
		},
	})

	err = b.Transaction(ctx, txns)
	if err != nil {
		t.Fatal(err)
	}

	// This should return nil, because the key foo/bar didn't exist when we ran that transaction, so the get
	// should return nil, and the put always returns nil
	for _, txn := range txns {
		if txn.Operation == physical.GetOperation {
			if txn.Entry.Value != nil {
				t.Fatalf("expected txn.entry.value to be nil but it was %q", string(txn.Entry.Value))
			}
		}
	}
}

// TestConsul_TransactionalBackend_GetTransactions tests that passing a slice of transactions to the
// consul backend will populate values for any transactions that are Get operations.
func TestConsul_TransactionalBackend_GetTransactions(t *testing.T) {
	// TODO: unskip this after Consul releases 1.14 and we update our API dep. It currently fails but should pass with Consul 1.14
	t.SkipNow()

	cleanup, config := consul.PrepareTestContainer(t, "1.4.4", false, true)
	defer cleanup()

	client, err := api.NewClient(config.APIConfig())
	if err != nil {
		t.Fatal(err)
	}

	txns := make([]*physical.TxnEntry, 0)
	ctx := context.Background()
	logger := logging.NewVaultLogger(log.Debug)
	backendConfig := map[string]string{
		"address":      config.Address(),
		"token":        config.Token,
		"path":         "vault/",
		"max_parallel": "-1",
	}

	be, err := NewConsulBackend(backendConfig, logger)
	if err != nil {
		t.Fatal(err)
	}
	b := be.(*ConsulBackend)

	defer func() {
		_, _ = client.KV().DeleteTree("foo/", nil)
	}()

	// Add some seed values to consul, and prepare our slice of transactions at the same time
	for i := 0; i < 64; i++ {
		key := fmt.Sprintf("foo/lol-%d", i)
		err := b.Put(ctx, &physical.Entry{Key: key, Value: []byte(fmt.Sprintf("value-%d", i))})
		if err != nil {
			t.Fatal(err)
		}

		txns = append(txns, &physical.TxnEntry{
			Operation: physical.GetOperation,
			Entry: &physical.Entry{
				Key: key,
			},
		})
	}

	for i := 0; i < 64; i++ {
		key := fmt.Sprintf("foo/lol-%d", i)
		if i%2 == 0 {
			txns = append(txns, &physical.TxnEntry{
				Operation: physical.PutOperation,
				Entry: &physical.Entry{
					Key:   key,
					Value: []byte("lmao"),
				},
			})
		} else {
			txns = append(txns, &physical.TxnEntry{
				Operation: physical.DeleteOperation,
				Entry: &physical.Entry{
					Key: key,
				},
			})
		}
	}

	if len(txns) != 128 {
		t.Fatal("wrong number of transactions")
	}

	err = b.Transaction(ctx, txns)
	if err != nil {
		t.Fatal(err)
	}

	// Check that our Get operations were populated with their values
	for i, txn := range txns {
		if txn.Operation == physical.GetOperation {
			val := []byte(fmt.Sprintf("value-%d", i))
			if !bytes.Equal(val, txn.Entry.Value) {
				t.Fatalf("expected %s to equal %s but it didn't", hex.EncodeToString(val), hex.EncodeToString(txn.Entry.Value))
			}
		}
	}
}

func TestConsulHABackend(t *testing.T) {
	cleanup, config := consul.PrepareTestContainer(t, "1.4.4", false, true)
	defer cleanup()

	client, err := api.NewClient(config.APIConfig())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	randPath := fmt.Sprintf("vault-%d/", time.Now().Unix())
	defer func() {
		client.KV().DeleteTree(randPath, nil)
	}()

	logger := logging.NewVaultLogger(log.Debug)
	backendConfig := map[string]string{
		"address":      config.Address(),
		"token":        config.Token,
		"path":         randPath,
		"max_parallel": "-1",
	}

	b, err := NewConsulBackend(backendConfig, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	b2, err := NewConsulBackend(backendConfig, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseHABackend(t, b.(physical.HABackend), b2.(physical.HABackend))

	detect, ok := b.(physical.RedirectDetect)
	if !ok {
		t.Fatalf("consul does not implement RedirectDetect")
	}
	host, err := detect.DetectHostAddr()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if host == "" {
		t.Fatalf("bad addr: %v", host)
	}
}
