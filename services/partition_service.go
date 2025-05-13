package services

import (
	"errors"
	"hash/fnv"
	"sync"
)

// PartitionService maneja el particionamiento de datos
type PartitionService struct {
	mu         sync.RWMutex
	partitions map[string]*Partition
	shardCount int
}

// Partition representa una partición de datos
type Partition struct {
	ID        string
	ShardKey  string
	Data      map[string]interface{}
	Replicas  []string
	IsPrimary bool
}

// ShardKeyGenerator genera claves de partición
type ShardKeyGenerator interface {
	GenerateKey(data interface{}) string
}

// HashShardKeyGenerator implementa generación de claves por hash
type HashShardKeyGenerator struct {
	shardCount int
}

// RangeShardKeyGenerator implementa generación de claves por rango
type RangeShardKeyGenerator struct {
	ranges []Range
}

// Range define un rango para particionamiento
type Range struct {
	Start interface{}
	End   interface{}
	Shard string
}

// NewPartitionService crea una nueva instancia del servicio de particionamiento
func NewPartitionService(shardCount int) *PartitionService {
	return &PartitionService{
		partitions: make(map[string]*Partition),
		shardCount: shardCount,
	}
}

// CreatePartition crea una nueva partición
func (s *PartitionService) CreatePartition(id, shardKey string, replicas []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.partitions[id]; exists {
		return ErrPartitionExists
	}

	s.partitions[id] = &Partition{
		ID:        id,
		ShardKey:  shardKey,
		Data:      make(map[string]interface{}),
		Replicas:  replicas,
		IsPrimary: true,
	}

	return nil
}

// GetPartition obtiene una partición por ID
func (s *PartitionService) GetPartition(id string) (*Partition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	partition, exists := s.partitions[id]
	if !exists {
		return nil, ErrPartitionNotFound
	}

	return partition, nil
}

// AddData agrega datos a una partición
func (s *PartitionService) AddData(partitionID string, key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	partition, exists := s.partitions[partitionID]
	if !exists {
		return ErrPartitionNotFound
	}

	partition.Data[key] = value
	return nil
}

// GetData obtiene datos de una partición
func (s *PartitionService) GetData(partitionID string, key string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	partition, exists := s.partitions[partitionID]
	if !exists {
		return nil, ErrPartitionNotFound
	}

	value, exists := partition.Data[key]
	if !exists {
		return nil, ErrDataNotFound
	}

	return value, nil
}

// GenerateShardKey genera una clave de partición usando hash
func (g *HashShardKeyGenerator) GenerateKey(data interface{}) string {
	h := fnv.New32a()
	h.Write([]byte(data.(string)))
	hash := h.Sum32()
	return string(rune(hash % uint32(g.shardCount)))
}

// GenerateShardKey genera una clave de partición usando rangos
func (g *RangeShardKeyGenerator) GenerateKey(data interface{}) string {
	value := data.(int)
	for _, r := range g.ranges {
		if value >= r.Start.(int) && value < r.End.(int) {
			return r.Shard
		}
	}
	return ""
}

// NewHashShardKeyGenerator crea un nuevo generador de claves por hash
func NewHashShardKeyGenerator(shardCount int) *HashShardKeyGenerator {
	return &HashShardKeyGenerator{
		shardCount: shardCount,
	}
}

// NewRangeShardKeyGenerator crea un nuevo generador de claves por rango
func NewRangeShardKeyGenerator(ranges []Range) *RangeShardKeyGenerator {
	return &RangeShardKeyGenerator{
		ranges: ranges,
	}
}

// BalancePartitions balancea las particiones entre shards
func (s *PartitionService) BalancePartitions() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Implementar lógica de balanceo
	// ...
	return nil
}

// ReplicatePartition replica una partición
func (s *PartitionService) ReplicatePartition(partitionID string, replicaID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	partition, exists := s.partitions[partitionID]
	if !exists {
		return ErrPartitionNotFound
	}

	// Verificar si la réplica ya existe
	for _, r := range partition.Replicas {
		if r == replicaID {
			return ErrReplicaExists
		}
	}

	partition.Replicas = append(partition.Replicas, replicaID)
	return nil
}

// RemoveReplica elimina una réplica de una partición
func (s *PartitionService) RemoveReplica(partitionID string, replicaID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	partition, exists := s.partitions[partitionID]
	if !exists {
		return ErrPartitionNotFound
	}

	// Encontrar y eliminar la réplica
	for i, r := range partition.Replicas {
		if r == replicaID {
			partition.Replicas = append(partition.Replicas[:i], partition.Replicas[i+1:]...)
			return nil
		}
	}

	return ErrReplicaNotFound
}

// Errores del servicio
var (
	ErrPartitionExists   = errors.New("partition already exists")
	ErrPartitionNotFound = errors.New("partition not found")
	ErrDataNotFound      = errors.New("data not found")
	ErrReplicaExists     = errors.New("replica already exists")
	ErrReplicaNotFound   = errors.New("replica not found")
)
