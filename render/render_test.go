package render

import (
    "log"
	// "encoding/xml"
	// "errors"
	// "html/template"
	// "net/http"
	"net/http/httptest"
	// "strconv"
	// "strings"
	"testing"

	// testdata "github.com/gin-gonic/gin/testdata/protoexample"
    "github.com/stretchr/testify/assert"
	// "google.golang.org/protobuf/proto"
)

// TODO unit tests
// test errors

func TestRenderJSON(t *testing.T) {
    log.Printf("%v", "TestRenderJSON")

	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"foo":  "bar",
		"html": "<b>",
	}

	(JSON{data}).WriteContentType(w)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	err := (JSON{data}).Render(w)

	assert.NoError(t, err)
	assert.Equal(t, "{\"foo\":\"bar\",\"html\":\"\\u003cb\\u003e\"}", w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

