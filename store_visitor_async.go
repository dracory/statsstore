package statsstore

import (
	"sync"
	"sync/atomic"
)

// VisitorCreateAsyncBufferSizeDefault bounds the pending-insert queue when
// VisitorCreateAsyncEnabled is set. Sized so a burst of requests cannot
// exhaust memory; when full, new visitor rows are dropped rather than
// blocking the request path on a database write.
const VisitorCreateAsyncBufferSizeDefault = 4096

// visitorAsyncState holds the queue and worker lifecycle shared by the
// store's async visitor insert path.
type visitorAsyncState struct {
	queue   chan map[string]any
	done    chan struct{}
	wg      sync.WaitGroup
	once    sync.Once
	closed  atomic.Bool
	dropped atomic.Uint64
}

// startVisitorAsyncWorker launches the single drain goroutine that
// performs queued INSERTs. Called from NewStore when
// VisitorCreateAsyncEnabled is true.
func (st *storeImplementation) startVisitorAsyncWorker(capacity int) {
	if capacity <= 0 {
		capacity = VisitorCreateAsyncBufferSizeDefault
	}
	st.visitorAsync = &visitorAsyncState{
		queue: make(chan map[string]any, capacity),
		done:  make(chan struct{}),
	}
	st.visitorAsync.wg.Add(1)
	go func() {
		defer st.visitorAsync.wg.Done()
		for {
			select {
			case row := <-st.visitorAsync.queue:
				if err := st.db.Query().Table(st.visitorTableName).Create(row); err != nil {
					st.logger.Error("async visitor insert failed", "error", err)
				}
			case <-st.visitorAsync.done:
				// Drain any rows still queued, then exit. The queue is
				// never closed — a late sender cannot panic, its row
				// simply sits unread and is GC'd with the queue.
				for {
					select {
					case row := <-st.visitorAsync.queue:
						if err := st.db.Query().Table(st.visitorTableName).Create(row); err != nil {
							st.logger.Error("async visitor insert failed", "error", err)
						}
					default:
						return
					}
				}
			}
		}
	}()
}

// visitorInsertAsync enqueues a prepared row for insert. Returns true when
// the row was accepted into the queue, false when the caller should fall
// back to a synchronous insert (async disabled). A full or closed queue
// drops the row and counts it — analytics must never block or fail the
// request path.
func (st *storeImplementation) visitorInsertAsync(row map[string]any) bool {
	if st.visitorAsync == nil {
		return false
	}
	if st.visitorAsync.closed.Load() {
		st.visitorAsync.dropped.Add(1)
		return true
	}
	select {
	case st.visitorAsync.queue <- row:
	default:
		st.visitorAsync.dropped.Add(1)
	}
	return true
}

// DrainVisitorWrites stops accepting rows and waits for the worker to
// flush all pending rows. The queue is never closed — rows enqueued
// after Drain are dropped and counted (the closed flag is checked before
// enqueue; a send that slips past the check lands in a queue nobody
// reads rather than panicking on a closed channel). Safe to call more
// than once and when async is disabled. Call during shutdown, before
// closing the underlying database.
func (st *storeImplementation) DrainVisitorWrites() {
	if st.visitorAsync == nil {
		return
	}
	st.visitorAsync.once.Do(func() {
		st.visitorAsync.closed.Store(true)
		close(st.visitorAsync.done)
		st.visitorAsync.wg.Wait()
	})
}

// VisitorDroppedCount returns how many visitor rows were dropped because
// the async insert queue was full or already drained.
func (st *storeImplementation) VisitorDroppedCount() uint64 {
	if st.visitorAsync == nil {
		return 0
	}
	return st.visitorAsync.dropped.Load()
}

