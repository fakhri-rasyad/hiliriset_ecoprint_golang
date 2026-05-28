package mqttpackage

import (
	"log"
	"sync"
	"time"

	"hiliriset_ecoprint_golang/models"

	"github.com/google/uuid"
)

const (
	batchFlushInterval = 5 * time.Second
	batchMaxSize       = 50
)

type batchEntry struct {
	records []models.SessionRecordInput
}

type TelemetryBatcher struct {
	mu      sync.Mutex
	batches map[uuid.UUID]*batchEntry
	service SessionRecordCreator
	stopCh  chan struct{}
}

func NewTelemetryBatcher(service SessionRecordCreator) *TelemetryBatcher {
	b := &TelemetryBatcher{
		batches: make(map[uuid.UUID]*batchEntry),
		service: service,
		stopCh:  make(chan struct{}),
	}
	go b.flushLoop()
	return b
}

func (b *TelemetryBatcher) Add(sessionID uuid.UUID, record models.SessionRecordInput) {
	b.mu.Lock()
	defer b.mu.Unlock()

	entry, ok := b.batches[sessionID]
	if !ok {
		entry = &batchEntry{}
		b.batches[sessionID] = entry
	}

	entry.records = append(entry.records, record)

	if len(entry.records) >= batchMaxSize {
		go b.flushSession(sessionID, entry.records)
		delete(b.batches, sessionID)
	}
}

func (b *TelemetryBatcher) flushLoop() {
	ticker := time.NewTicker(batchFlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			b.flushAll()
		case <-b.stopCh:
			b.flushAll()
			return
		}
	}
}

func (b *TelemetryBatcher) flushAll() {
	b.mu.Lock()
	if len(b.batches) == 0 {
		b.mu.Unlock()
		return
	}

	snapshot := b.batches
	b.batches = make(map[uuid.UUID]*batchEntry)
	b.mu.Unlock()

	for sessionID, entry := range snapshot {
		go b.flushSession(sessionID, entry.records)
	}
}

func (b *TelemetryBatcher) flushSession(sessionID uuid.UUID, records []models.SessionRecordInput) {
	for _, record := range records {
		if err := b.service.CreateRecord(sessionID, record); err != nil {
			log.Printf("Batcher: failed to save record for session %s: %v", sessionID, err)
		}
	}
}

func (b *TelemetryBatcher) Stop() {
	close(b.stopCh)
}
