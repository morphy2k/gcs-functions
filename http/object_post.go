package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/morphy2k/gcs-functions/event"

	"cloud.google.com/go/functions/metadata"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func PostFinalizedObject(ctx context.Context, e event.GCSEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}

	if meta.EventType != event.StorageObjectFinalize {
		return fmt.Errorf("bad event type \"%s\"", meta.EventType)
	}

	opts := option.WithScopes(storage.ScopeReadOnly)
	client, err := storage.NewClient(ctx, opts)
	if err != nil {
		return err
	}

	r, err := client.Bucket(e.Bucket).Object(e.Name).NewReader(ctx)
	if err != nil {
		return err
	}
	defer r.Close()

	name := e.Name
	if trimPath {
		parts := strings.Split(e.Name, "/")
		name = parts[len(parts)-1]
	}

	return postForm(ctx, name, e.ContentType, r)
}

func postForm(ctx context.Context, name, contentType string, r io.ReadCloser) error {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, name))
	header.Set("Content-Type", contentType)

	fw, err := w.CreatePart(header)
	if err != nil {
		return err
	}

	if _, err := io.Copy(fw, r); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	u := postURL.String()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		err := fmt.Errorf("error while posting file: %s", res.Status)
		bodyToStderr(res.Body)
		return err
	}

	return nil
}
