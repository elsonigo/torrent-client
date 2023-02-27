package magnet

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const uri = "magnet:?xt=urn:btih:fc9f6cccfc2f1839f1bc1d3a8f63a9247c55f608&dn=%E7%B5%90%E6%9D%9F%E3%83%90%E3%83%B3%E3%83%89%20%2F%20kessoku%20band%20%28Bocchi%20the%20Rock%21%29%20-%20kessoku%20band%20%2F%20%E7%B5%90%E6%9D%9F%E3%83%90%E3%83%B3%E3%83%89%20%5BFLAC%5D&tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce&xl=12387987"

const uri_no_tr = "magnet:?xt=urn:btih:fc9f6cccfc2f1839f1bc1d3a8f63a9247c55f608&dn=%E7%B5%90%E6%9D%9F%E3%83%90%E3%83%B3%E3%83%89%20%2F%20kessoku%20band%20%28Bocchi%20the%20Rock%21%29%20-%20kessoku%20band%20%2F%20%E7%B5%90%E6%9D%9F%E3%83%90%E3%83%B3%E3%83%89%20%5BFLAC%5D&xl=12387987"

const uri_no_xt = "magnet:&dn=%E7%B5%90%E6%9D%9F%E3%83%90%E3%83%B3%E3%83%89%20%2F%20kessoku%20band%20%28Bocchi%20the%20Rock%21%29%20-%20kessoku%20band%20%2F%20%E7%B5%90%E6%9D%9F%E3%83%90%E3%83%B3%E3%83%89%20%5BFLAC%5D&xl=12387987"

const uri_malformed_xt = "magnet:?xt=ur:btih:fc9f6cccfc2f1839f1bc1d3a8f63a9247c55f608&dn=%E7%B5%90%E6%9D%9F%E3%83%90%E3%83%B3%E3%83%89%20%2F%20kessoku%20band%20%28Bocchi%20the%20Rock%21%29%20-%20kessoku%20band%20%2F%20%E7%B5%90%E6%9D%9F%E3%83%90%E3%83%B3%E3%83%89%20%5BFLAC%5D&tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce&xl=12387987"

const uri_invalid_xt = "magnet:?xt=urn:btih:fc9f6cccfc2f1839f1bc1d3a8f63a9247c55f60"

const uri_invalid_scheme = "magent:?xt=urn:btih:fc9f6cccfc2f1839f1bc1d3a8f63a9247c55f608"

func TestParseMagnet(t *testing.T) {
	m, err := ParseMagnet(uri)
	require.Nil(t, err)

	assert.Equal(
		t,
		"結束バンド / kessoku band (Bocchi the Rock!) - kessoku band / 結束バンド [FLAC]",
		m.DisplayName,
	)
	assert.True(t, len(m.Trackers) > 0)
	assert.True(t, len(m.InfoHash) == 20)
}

func TestParseInvalidMagnet(t *testing.T) {
	_, err := ParseMagnet(uri_no_tr)
	assert.Equal(t, err, fmt.Errorf("magnet link is missing 'address trackers' parameter"))
}

func TestMissingXt(t *testing.T) {
	_, err := ParseMagnet(uri_no_xt)
	assert.Equal(t, err, fmt.Errorf("magnet link is missing the 'exact topic' parameter"))
}

func TestMalformedXt(t *testing.T) {
	_, err := ParseMagnet(uri_malformed_xt)
	assert.Equal(t, err, fmt.Errorf("malformed 'exact topic'"))
}

func TestInvalidXt(t *testing.T) {
	_, err := ParseMagnet(uri_invalid_xt)
	assert.Equal(t, err, fmt.Errorf("error decoding 'exact topic'"))
}

func TestInvalidScheme(t *testing.T) {
	_, err := ParseMagnet(uri_invalid_scheme)
	assert.Equal(t, err, fmt.Errorf("unexpected url scheme"))
}
