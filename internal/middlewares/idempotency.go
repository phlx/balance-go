package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"mime"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	goredis "github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"balance/internal/redis"
)

const (
	IdempotencyKeyInHeader = "X-Idempotency-Key"
	IdempotencyKeyInQuery  = "idempotency_key"
	IdempotencyKeyInBody   = "idempotency_key"
	StorageKeyPrefix       = "idempotence:"
	StorageKeyExpiration   = 60 * time.Minute
)

type response struct {
	Status  int    `json:"status"`
	Data    string `json:"data"`
	Headers string `json:"headers"`
}

type cachedResponseWriter struct {
	gin.ResponseWriter
	cachedBody *bytes.Buffer
}

type storage struct {
	client redis.Client
	ctx    context.Context
	log    *logrus.Logger
}

func (s storage) loadResponse(key string) (*response, error) {
	result, err := s.client.Get(s.ctx, StorageKeyPrefix+key)
	if err != nil {
		return nil, err
	}
	var r *response
	err = json.Unmarshal([]byte(result), &r)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to unmarshal JSON from string '%s'", result)
	}
	return r, nil
}

func (s storage) saveResponse(key string, r *response) error {
	result, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return s.client.Set(s.ctx, StorageKeyPrefix+key, result, StorageKeyExpiration)
}

func (crw cachedResponseWriter) Write(bytes []byte) (int, error) {
	crw.cachedBody.Write(bytes)
	return crw.ResponseWriter.Write(bytes)
}

func Idempotency(ctx context.Context, client redis.Client, log *logrus.Logger) gin.HandlerFunc {
	s := storage{
		client: client,
		log:    log,
		ctx:    ctx,
	}

	return func(context *gin.Context) {
		if context.Request.Method != "POST" {
			context.Next()
			return
		}

		key := ""

		getIdempotencyKeyFunctions := []func(*gin.Context, *logrus.Logger) string{
			getIdempotencyKeyFromHeader,
			getIdempotencyKeyFromQuery,
			getIdempotencyKeyFromBody,
		}

		for _, get := range getIdempotencyKeyFunctions {
			if key = get(context, log); key != "" {
				break
			}
		}

		if key == "" {
			context.Next()
			return
		}

		resp, err := s.loadResponse(key)
		if err != nil && err != goredis.Nil {
			s.log.Warning(errors.Wrap(err, "Unable to load response from storage"))
		}
		if resp != nil {
			err := writeResponseToContext(resp, context)
			if err != nil {
				s.log.Warning(errors.Wrap(err, "Unable to write response to context"))
				context.Next()
				return
			}
			context.Abort()
			return
		}

		setCachedResponseWriterToContext(context)

		context.Next()

		writer, err := getCachedResponseWriterFromContext(context)
		if err != nil {
			return
		}
		defer writer.cachedBody.Reset()

		resp, err = makeResponseFromContext(context)
		if err != nil {
			log.Warning(errors.Wrap(err, "Unable to make response from context"))
		} else {
			err := s.saveResponse(key, resp)
			if err != nil {
				log.Warning(errors.Wrap(err, "Unable to save response"))
			}
		}
	}
}

func getIdempotencyKeyFromHeader(context *gin.Context, _ *logrus.Logger) string {
	return context.Request.Header.Get(IdempotencyKeyInHeader)
}

func getIdempotencyKeyFromQuery(context *gin.Context, _ *logrus.Logger) string {
	return context.Query(IdempotencyKeyInQuery)
}

func getIdempotencyKeyFromBody(context *gin.Context, log *logrus.Logger) string {
	if context.Request.ContentLength == 0 {
		return ""
	}
	mimeType, _, err := mime.ParseMediaType(context.ContentType())
	if err != nil {
		log.Warning(errors.Wrap(err, "Unable to parse content-type"))
		return ""
	}
	body, err := ioutil.ReadAll(context.Request.Body)
	context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	if err != nil {
		log.Warning(errors.Wrap(err, "Unable to read request body"))
		return ""
	}
	switch mimeType {
	case gin.MIMEJSON:
		var input map[string]interface{}
		err = json.Unmarshal(body, &input)
		if err != nil {
			log.WithField("body", string(body)).
				Warning(errors.Wrap(err, "Unable to unmarshal body to JSON"))

			return ""
		}
		key, _ := input[IdempotencyKeyInBody].(string)
		return key
	case gin.MIMEPOSTForm:
		v, err := url.ParseQuery(string(body))
		if err != nil {
			log.WithField("body", string(body)).
				Warning(errors.Wrap(err, "Unable to parse body as URL string"))

			return ""
		}
		return v.Get(IdempotencyKeyInBody)
	default:
		return ""
	}
}

func setCachedResponseWriterToContext(context *gin.Context) {
	context.Writer = &cachedResponseWriter{
		ResponseWriter: context.Writer,
		cachedBody:     bytes.NewBufferString(""),
	}
}

func getCachedResponseWriterFromContext(context *gin.Context) (*cachedResponseWriter, error) {
	crw, ok := context.Writer.(*cachedResponseWriter)
	if !ok {
		return nil, errors.New("Unable to get context writer as cachedResponseWriter")
	}
	return crw, nil
}

func makeResponseFromContext(context *gin.Context) (*response, error) {
	w, err := getCachedResponseWriterFromContext(context)
	if err != nil {
		return nil, err
	}
	headers, err := json.Marshal(w.Header())
	if err != nil {
		return nil, err
	}
	return &response{
		Status:  w.Status(),
		Data:    w.cachedBody.String(),
		Headers: string(headers),
	}, nil
}

func writeResponseToContext(resp *response, context *gin.Context) error {
	if resp.Headers != "" && resp.Headers != "{}" {
		var headers map[string][]string
		err := json.Unmarshal([]byte(resp.Headers), &headers)
		if err != nil {
			return errors.Wrapf(err, "Unable to unmarshal JSON from string '%s'", resp.Headers)
		}
		for header, value := range headers {
			if len(value) != 0 {
				context.Header(header, value[0])
			}
		}
		context.Writer.WriteHeaderNow()
	}
	if resp.Data != "<nil>" {
		_, err := context.Writer.Write([]byte(resp.Data))
		if err != nil {
			return err
		}
	}
	context.Status(resp.Status)
	return nil
}
