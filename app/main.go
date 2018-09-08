package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/line/clova-cek-sdk-go/cek"
	"github.com/sugyan/clova-hatena-hotentry/hatena"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func main() {
	http.HandleFunc("/callback", callbackHandler)
	appengine.Main()
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	response, err := hotentry(ctx, r)
	if err != nil {
		log.Errorf(ctx, err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if response != nil {
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Errorf(ctx, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func hotentry(ctx context.Context, r *http.Request) (*cek.ResponseMessage, error) {
	ext := cek.NewExtension(os.Getenv("EXTENSION_ID"))
	message, err := ext.ParseRequest(r)
	if err != nil {
		return nil, err
	}
	switch request := message.Request.(type) {
	case *cek.LaunchRequest:
		return cek.NewResponseBuilder().
			OutputSpeech(cek.NewOutputSpeechBuilder().
				AddSpeechText("人気エントリーを読み上げます。カテゴリーを指定してください。", cek.SpeechInfoLangJA).
				Build()).
			Build(), nil
	case *cek.IntentRequest:
		for _, slot := range request.Intent.Slots {
			client := hatena.NewClient(hatena.WithHTTPClient(urlfetch.Client(ctx)))
			entries, err := client.Fetch(hatena.Category(slot.Value))
			if err != nil {
				return nil, err
			}
			if len(entries) > 0 {
				outputSpeech := cek.NewOutputSpeechBuilder().
					AddSpeechText(fmt.Sprintf("「%s」の人気エントリーです。", slot.Value), cek.SpeechInfoLangJA)
				for _, entry := range entries[0:3] {
					speechText := fmt.Sprintf("「%s」。 %d ブックマーク。", entry.Title, entry.BookmarkCount)
					outputSpeech.AddSpeechText(speechText, cek.SpeechInfoLangJA)
				}
				return cek.NewResponseBuilder().
					OutputSpeech(outputSpeech.Build()).
					Build(), nil
			}
		}
		return cek.NewResponseBuilder().
			OutputSpeech(cek.NewOutputSpeechBuilder().
				AddSpeechText("よく分かりませんでした。もう一回言ってください。", cek.SpeechInfoLangJA).
				Build()).
			Build(), nil
	default:
		log.Warningf(ctx, "request: %v", request)
		return nil, errors.New("Unknown request type")
	}
}
