package elasticsearch

import (
	"context"
	"fmt"
	"gopkg.in/olivere/elastic.v5"
	"gopkg.in/olivere/elastic.v5/config"
	"io"
	"o.o/backend/pkg/common/cmenv"
	"o.o/common/l"
)

var ll = l.New()

const (
	DefaultIndex       = "portsip_pbx"
	DefaultMappingType = "cdr"
	DefaultSize        = 50
)

type ElasticSearch struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (c ElasticSearch) ConnectionString() string {
	s := ""
	if c.Username == "" || c.Password == "" {
		s = fmt.Sprintf("http://%s:%s?sniff=false", c.Host, c.Port)
	} else {
		s = fmt.Sprintf("http://%s:%s@%s:%s?sniff=false", c.Username, c.Password, c.Host, c.Port)
	}
	return s
}

func Connect(cfg ElasticSearch) Store {
	store := New(cfg.ConnectionString())
	return store
}

type Store interface {
	Ping(url string) (string, int, error)
	GetIndex(id string) string
	Scroll(index string, query elastic.Query, sort elastic.Sorter, scrollId string) (*elastic.SearchResult, error)
}

type elasticSearchStore struct {
	client *elastic.Client
	index  string
}

func New(elasticURL string) Store {
	cfg, err := config.Parse(elasticURL)
	if err != nil {
		ll.Fatal("Unable to connect to Elasticsearch", l.Error(err), l.String("ConnectionString", elasticURL))
	}

	client, err := elastic.NewClientFromConfig(cfg)
	if err != nil {
		ll.Fatal("Unable to connect to Elasticsearch", l.Error(err), l.String("endpoint", elasticURL))
	}

	store := &elasticSearchStore{
		client: client,
	}

	switch cmenv.Env() {
	case cmenv.EnvDev:
		store.index = fmt.Sprintf("d_%s", DefaultIndex)
	case cmenv.EnvSandbox, cmenv.EnvStag:
		store.index = fmt.Sprintf("s_%s", DefaultIndex)
	case cmenv.EnvProd:
		store.index = DefaultIndex
	default:
		ll.Fatal("NTX: Invalid env")
	}

	return store
}

func (e elasticSearchStore) Ping(url string) (string, int, error) {
	info, code, err := e.client.Ping(url).Do(context.Background())
	if err != nil {
		panic(err)
	}

	return info.ClusterName, code, nil
}

func (e elasticSearchStore) GetIndex(id string) string {
	return fmt.Sprintf("%s_%v", e.index, id)
}

func (e elasticSearchStore) Scroll(index string, query elastic.Query, sort elastic.Sorter, scrollId string) (*elastic.SearchResult, error) {
	scroller := e.client.Scroll(index).
		Type(DefaultMappingType).
		Routing(index).
		Query(query).
		SortBy(sort).
		Size(DefaultSize)

	if scrollId != "" {
		scroller.ScrollId(scrollId)
	}

	results, err := scroller.Do(context.Background())
	if err == io.EOF {
		return results, nil
	}

	return results, nil
}
