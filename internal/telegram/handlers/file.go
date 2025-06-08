package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func (h *handlers) FileHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	client, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("admin_user", "admin_pass", ""),
		Secure: false,
		Region: "us-east-1",
	})
	if err != nil {
		h.logger.Error("Failed to create MinIO client", zap.Error(err))
		return nil
	}

	var (
		fileInfo *gotgbot.File
		mimeType string
	)
	if ctx.EffectiveMessage.Audio != nil {
		fileInfo, _ = b.GetFile(ctx.EffectiveMessage.Audio.FileId, nil)
		mimeType = ctx.EffectiveMessage.Audio.MimeType
	} else {
		fileInfo, _ = b.GetFile(ctx.EffectiveMessage.Voice.FileId, nil)
		mimeType = ctx.EffectiveMessage.Voice.MimeType
	}
	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", b.Token, fileInfo.FilePath)

	resp, _ := http.Get(fileURL)
	defer resp.Body.Close()

	if _, err := client.PutObject(context.Background(), "dawgbot", fileInfo.FileId, resp.Body, resp.ContentLength, minio.PutObjectOptions{
		ContentType: mimeType,
	}); err != nil {
		h.logger.Error("Failed to upload file to MinIO", zap.Error(err))
	}

	b.SendMessage(ctx.EffectiveChat.Id, "File added to playlist!", nil)

	messageJSON, _ := json.Marshal(map[string]string{
		"file_url": "http://localhost:9000/dawgbot/" + fileInfo.FileId,
	})
	h.rabbit.Publish("", "playlist", false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageJSON,
		},
	)

	return nil
}
