package helpers

import (
	"fmt"
	"os/exec"

	"github.com/brutalzinn/go-schedule-manager/utils"
	"github.com/sirupsen/logrus"
)

func CreateAudioTTS(text string) (string, error) {
	fileName := fmt.Sprintf("/static/audio/%s.wav", utils.GetMD5Hash(text))
	cmd := exec.Command("espeak-ng", "-v", "mb-br4", text, "-w", fileName)
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	logrus.Info("Audio file created: ", fileName)
	return fileName, nil
}
