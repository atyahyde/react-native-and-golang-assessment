package cache

import (
	"context"
	"time"
)

// Cache adalah abstraksi di atas Redis supaya TaskService bisa diuji dengan
// mock tanpa Redis sungguhan.
type Cache interface {
	Get(ctx context.Context, key string) (value string, hit bool, err error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error

	// TrackListKey mendaftarkan sebuah key cache list (GET /api/tasks) supaya
	// nanti bisa diinvalidasi sekaligus lewat InvalidateListCache. Diperlukan
	// karena cache key list mengandung query parameter (status/keyword/page/dst),
	// jadi ada banyak variasi key yang harus dihapus bersamaan saat ada
	// create/update/delete.
	TrackListKey(ctx context.Context, key string) error

	// InvalidateListCache menghapus seluruh cache list task yang pernah
	// ditrack. Dipanggil setelah create/update/delete berhasil.
	InvalidateListCache(ctx context.Context) error
}
