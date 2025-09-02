// pkg/httpmw/logging.go
package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/bncunha/erp-api/src/infrastructure/logs"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	CtxKeyReqID        = "req_id"
	maxLoggedBodyBytes = 1 << 20 // 1 MiB por segurança
	truncateSuffix     = "...[truncated]"
	defaultContentType = "application/octet-stream"
)

// ==== Response recorder para capturar o corpo enviado ao cliente ====

type responseRecorder struct {
	http.ResponseWriter
	buf         *bytes.Buffer
	status      int
	wroteHeader bool
	limit       int64
}

func newResponseRecorder(w http.ResponseWriter, limit int64) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: w,
		buf:            &bytes.Buffer{},
		status:         http.StatusOK,
		limit:          limit,
	}
}

func (r *responseRecorder) WriteHeader(code int) {
	if r.wroteHeader {
		return
	}
	r.status = code
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	// espelha para o cliente
	n, err := r.ResponseWriter.Write(b)
	// guarda (limitado) para log
	if r.buf.Len() < int(r.limit) {
		remain := int(r.limit) - r.buf.Len()
		if remain > 0 {
			if len(b) > remain {
				r.buf.Write(b[:remain])
			} else {
				r.buf.Write(b)
			}
		}
	}
	return n, err
}

// ==== Helpers ====

func readAndRestoreBody(req *http.Request, max int64) ([]byte, string) {
	if req.Body == nil {
		return nil, ""
	}
	ct := req.Header.Get(echo.HeaderContentType)
	if ct == "" {
		ct = defaultContentType
	}

	// Leia TUDO para restaurar o body fielmente
	all, _ := io.ReadAll(req.Body)
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewBuffer(all))

	// Para LOG: mantenha só até 'max'
	if int64(len(all)) > max {
		return all[:max], ct // logará truncado, mas o handler recebe tudo
	}
	return all, ct
}

func normalizeContentType(ct string) (mt string) {
	mt, _, _ = mime.ParseMediaType(ct)
	if mt == "" {
		mt = defaultContentType
	}
	return
}

func maybeCompactJSON(raw []byte) (string, bool) {
	if !json.Valid(raw) {
		return string(raw), false
	}
	var buf bytes.Buffer
	if err := json.Compact(&buf, raw); err != nil {
		return string(raw), false
	}
	return buf.String(), true
}

func truncateIfNeeded(s string, limit int) (string, bool) {
	if len(s) <= limit {
		return s, false
	}
	return s[:limit] + truncateSuffix, true
}

// (Opcional) Mascara campos sensíveis em JSON (ex.: "password", "token").
func maskJSONFields(raw []byte, fields ...string) []byte {
	if len(raw) == 0 || len(fields) == 0 || !json.Valid(raw) {
		return raw
	}
	var data any
	if err := json.Unmarshal(raw, &data); err != nil {
		return raw
	}
	m := map[string]struct{}{}
	for _, f := range fields {
		m[strings.ToLower(f)] = struct{}{}
	}
	var walk func(any) any
	walk = func(v any) any {
		switch t := v.(type) {
		case map[string]any:
			for k, val := range t {
				if _, ok := m[strings.ToLower(k)]; ok {
					t[k] = "***"
					continue
				}
				t[k] = walk(val)
			}
			return t
		case []any:
			for i := range t {
				t[i] = walk(t[i])
			}
			return t
		default:
			return t
		}
	}
	data = walk(data)
	out, err := json.Marshal(data)
	if err != nil {
		return raw
	}
	return out
}

// ==== Middlewares ====

func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// req_id
			reqID := c.Request().Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = uuid.NewString()
			}
			c.Set(CtxKeyReqID, reqID)
			c.Response().Header().Set("X-Request-ID", reqID)

			// --- Captura Request Body (limitado) ---
			reqBodyBytes, reqCT := readAndRestoreBody(c.Request(), maxLoggedBodyBytes)
			reqMT := normalizeContentType(reqCT)

			// (opcional) mascarar campos sensíveis
			// reqBodyBytes = maskJSONFields(reqBodyBytes, "password", "token", "authorization")

			var reqBodyStr string
			if strings.Contains(reqMT, "json") {
				if pretty, ok := maybeCompactJSON(reqBodyBytes); ok {
					reqBodyStr = pretty
				} else {
					reqBodyStr = string(reqBodyBytes)
				}
			} else if strings.HasPrefix(reqMT, "text/") || strings.Contains(reqMT, "xml") || strings.Contains(reqMT, "x-www-form-urlencoded") {
				reqBodyStr = string(reqBodyBytes)
			} else if len(reqBodyBytes) > 0 {
				reqBodyStr = "[binary body omitted]"
			}

			// Truncagem explícita (string)
			if reqBodyStr != "" {
				reqBodyStr, _ = truncateIfNeeded(reqBodyStr, int(maxLoggedBodyBytes))
			}

			// --- Intercepta Response ---
			origWriter := c.Response().Writer
			rec := newResponseRecorder(origWriter, maxLoggedBodyBytes)
			c.Response().Writer = rec

			// Executa handlers
			err := next(c)

			// Dados de resposta
			latency := time.Since(start)
			status := rec.status
			if !rec.wroteHeader { // Echo pode não chamar WriteHeader
				status = c.Response().Status
			}
			resCT := c.Response().Header().Get(echo.HeaderContentType)
			resMT := normalizeContentType(resCT)

			resBodyBytes := rec.buf.Bytes()
			var resBodyStr string
			if strings.Contains(resMT, "json") {
				if pretty, ok := maybeCompactJSON(resBodyBytes); ok {
					resBodyStr = pretty
				} else {
					resBodyStr = string(resBodyBytes)
				}
			} else if strings.HasPrefix(resMT, "text/") || strings.Contains(resMT, "xml") || strings.Contains(resMT, "html") {
				resBodyStr = string(resBodyBytes)
			} else if len(resBodyBytes) > 0 {
				resBodyStr = "[binary body omitted]"
			}
			if resBodyStr != "" {
				resBodyStr, _ = truncateIfNeeded(resBodyStr, int(maxLoggedBodyBytes))
			}

			// Campos úteis
			fields := map[string]any{
				"req_id":     reqID,
				"method":     c.Request().Method,
				"path":       c.Path(),
				"status":     status,
				"ip":         c.RealIP(),
				"latency_ms": latency.Milliseconds(),
				"req_ct":     reqMT,
				"res_ct":     resMT,
			}

			// Inclui os corpos no log se existirem
			logBodies := c.Response().Status >= 400 // ou: status >= 400
			if logBodies {
				if reqBodyStr != "" {
					fields["req_body"] = reqBodyStr
				}
				if resBodyStr != "" {
					fields["res_body"] = resBodyStr
				}
			}
			// Tamanhos (não truncados) para monitorar payloads
			fields["req_size"] = len(reqBodyBytes)
			fields["res_size"] = len(resBodyBytes)

			log := logs.Logger.With(fields)

			if err != nil {
				log.Errorf("request completed with error: %v", err)
				return err
			}

			switch {
			case status >= 500:
				log.Errorf("request completed")
			case status >= 400:
				log.Warnf("request completed")
			default:
				log.Infof("request completed")
			}

			return nil
		}
	}
}

// Recovery que loga panic com req_id
func Recover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					reqID, _ := c.Get(CtxKeyReqID).(string)
					logs.Logger.With(map[string]any{
						"req_id": reqID,
						"path":   c.Path(),
					}).Errorf("panic: %v", r)
					err = echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
				}
			}()
			return next(c)
		}
	}
}
