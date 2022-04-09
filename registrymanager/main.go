package registrymanager

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"github.com/harry1453/go-common-file-dialog/cfd"
	"github.com/harry1453/go-common-file-dialog/cfdutil"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

func GetBrowserPath() string {
	log.Println("User requested browserpath")
	k, _, _ := registry.CreateKey(registry.CURRENT_USER, "Software\\downgrid", registry.QUERY_VALUE|registry.SET_VALUE)

	val, _, err := k.GetStringValue("browserpath")
	if err == windows.ERROR_FILE_NOT_FOUND || len(os.Args) == 1 {
		log.Println("Prompting user to choose a web browser")
		result, err := cfdutil.ShowOpenFileDialog(cfd.DialogConfig{
			Title: "Choose your web browser",
			Role:  "ChooseWebBrowser",
			FileFilters: []cfd.FileFilter{
				{
					DisplayName: "Executables (*.exe)",
					Pattern:     "*.exe",
				},
			},
		})
		if err == cfd.ErrorCancelled {
			log.Fatal("Dialog was cancelled by the user.")
		} else if err != nil {
			log.Fatal(err)
		}
		err = k.SetStringValue("browserpath", result)
		if err != nil {
			log.Fatal(err)
		}
		return result
	}
	log.Println("Browser path is " + val)
	return val
}

func Register() {
	log.Println("Checking if downgrid is registered as a web browser")
	p, _ := os.Executable()

	_, existing, _ := registry.CreateKey(registry.LOCAL_MACHINE, "SOFTWARE\\Classes\\downgridURL\\shell\\open\\command", registry.QUERY_VALUE)
	k, _, _ := registry.CreateKey(registry.CURRENT_USER, "Software\\downgrid", registry.QUERY_VALUE|registry.SET_VALUE)

	if existing {
		r, _, _ := k.GetStringValue("Path")
		if r == p {
			log.Println("Downgrid is registered at " + r)
			return
		}
	}

	r := []string{
		"Windows Registry Editor Version 5.00\n\n",
		"[HKEY_LOCAL_MACHINE\\SOFTWARE\\downgrid\\Capabilities]\n",
		"\"ApplicationDescription\"=\"downgrid\"\n",
		"\"ApplicationIcon\"=\"C:\\\\downgrid\\\\downgrid.exe,0\"\n",
		"\"ApplicationName\"=\"downgrid\"\n\n",
		"[HKEY_LOCAL_MACHINE\\SOFTWARE\\downgrid\\Capabilities\\FileAssociations]\n",
		"\".htm\"=\"downgridURL\"\n",
		"\".html\"=\"downgridURL\"\n",
		"\".shtml\"=\"downgridURL\"\n",
		"\".xht\"=\"downgridURL\"\n",
		"\".xhtml\"=\"downgridURL\"\n\n",
		"[HKEY_LOCAL_MACHINE\\SOFTWARE\\downgrid\\Capabilities\\URLAssociations]\n",
		"\"ftp\"=\"downgridURL\"\n",
		"\"http\"=\"downgridURL\"\n",
		"\"https\"=\"downgridURL\"\n\n",
		"[HKEY_LOCAL_MACHINE\\SOFTWARE\\RegisteredApplications]\n",
		"\"downgrid\"=\"Software\\\\downgrid\\\\Capabilities\"\n\n",
		"[HKEY_LOCAL_MACHINE\\Software\\Classes\\downgridURL]\n",
		"@=\"downgrid Document\"\n",
		"\"FriendlyTypeName\"=\"downgrid Document\"\n\n",
		"[HKEY_LOCAL_MACHINE\\Software\\Classes\\downgridURL\\shell]\n",
		"\n",
		"[HKEY_LOCAL_MACHINE\\Software\\Classes\\downgridURL\\shell\\open]\n",
		"\n",
		"[HKEY_LOCAL_MACHINE\\Software\\Classes\\downgridURL\\shell\\open\\command]\n",
		"@=\"\\\"" + strings.ReplaceAll(p, "\\", "\\\\") + "\\\" \\\"%1\\\"\"",
	}

	file, err := os.OpenFile("register.reg", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Registry file created")

	datawriter := bufio.NewWriter(file)
	for _, data := range r {
		_, _ = datawriter.WriteString(data)
	}
	err = datawriter.Flush()
	if err != nil {
		log.Fatal(err)
	}

	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Registry file written")

	cmd := exec.Command("regedit", "/s", "register.reg")
	err = cmd.Start()
	if err != nil {
		MessageBoxPlain("Error", "Looks like downgrid is not configured yet. Run this app as administrator for the initial setup.")
		log.Fatal("Requires administrative privileges")
	}

	err = cmd.Wait()
	if err != nil {
		return
	}

	log.Println("Registration succeeded")

	err = os.Remove("register.reg")
	if err != nil {
		return
	}
	log.Println("Registry file deleted")

	err = k.SetStringValue("Path", p)
	if err != nil {
		return
	}
	log.Println("Downgrid is registered at " + p)
}

func MessageBox(hwnd uintptr, caption, title string, flags uint) int {
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		hwnd,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(caption))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(flags))

	return int(ret)
}

func MessageBoxPlain(title, caption string) int {
	log.Println("Displaying a message box")
	const (
		NULL = 0
		MbOk = 0
	)
	return MessageBox(NULL, caption, title, MbOk)
}
