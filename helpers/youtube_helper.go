package helpers

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"

	"github.com/kkdai/youtube/v2"
	"github.com/sirupsen/logrus"
)

func DownloadYouTubeAudio(link string) string {
	client := youtube.Client{}
	video, err := client.GetVideo(link)
	if err != nil {
		logrus.Println("Error getting YouTube video:", err)
		return ""
	}

	formats := video.Formats.WithAudioChannels()
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		logrus.Println("Error getting YouTube stream:", err)
		return ""
	}
	regex := regexp.MustCompile(`([&$\+,:;=\?@#\s<>\[\]\{\}[\/]|\\\^%])+`)
	fileName := fmt.Sprintf("static/audio/%s.mp3", regex.ReplaceAllString(video.Title, "_"))
	file, err := os.Create(fileName)
	if err != nil {
		logrus.Println("Error creating file:", err)
		return ""
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		logrus.Println("Error saving video:", err)
		return ""
	}

	audioFileName := fmt.Sprintf("static/audio/%s_audio.mp3", regex.ReplaceAllString(video.Title, "_"))
	cmd := exec.Command("ffmpeg", "-i", fileName, audioFileName)
	err = cmd.Run()
	if err != nil {
		logrus.Println("Error converting to mp3:", err)
		return ""
	}
	err = os.Remove(fileName)
	if err != nil {
		logrus.Println("Error delete video:", err)
		return ""
	}
	return audioFileName
}
