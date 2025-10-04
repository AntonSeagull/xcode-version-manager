package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	appsDir      = "/Applications"
	xcodeApp     = "/Applications/Xcode.app"
	plistBuddy   = "/usr/libexec/PlistBuddy"
	infoPlistRel = "Contents/Info.plist"
	developerRel = "Contents/Developer"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	switch os.Args[1] {
	case "list":
		handleList()
	case "current":
		handleCurrent()
	case "switch":
		handleSwitch(os.Args[2:])
	default:
		usage()
	}
}

func usage() {
	fmt.Println(`xvm — Xcode Version Manager (Go)

Команды / Commands:
  xvm list            — показать все найденные Xcode и активную / list installed Xcode versions
  xvm current         — показать активную версию Xcode / show currently active Xcode
  xvm switch <ver>    — переключиться на указанную версию / switch to specific version
     Опции / Options:
       --dry-run      — только показать, что будет сделано / show planned actions without changes

Примеры / Examples:
  xvm list
  xvm current
  sudo xvm switch 16.2
  xvm switch 16.2 --dry-run
`)
}

func handleList() {
	activePath, activeVer, _ := activeXcode()
	fmt.Printf("Активная / Active: %s (%s)\n", displayApp(activePath), valueOr(activeVer, "неизвестно / unknown"))

	entries, err := findXcodes()
	if err != nil {
		fmt.Println("Ошибка поиска / Search error:", err)
		return
	}

	if len(entries) == 0 {
		fmt.Println("❌ В /Applications Xcode не найден / No Xcode installations found in /Applications")
		return
	}

	fmt.Println("\nНайдено / Found in /Applications:")
	for _, e := range entries {
		marker := " "
		if sameFile(e.path, activePath) {
			marker = "*"
		}
		fmt.Printf(" %s %s (%s)\n", marker, displayApp(e.path), valueOr(e.version, "версия неизвестна / unknown version"))
	}
}

func handleCurrent() {
	path, ver, err := activeXcode()
	if err != nil {
		fmt.Println("Не удалось определить активный Xcode / Failed to detect active Xcode:", err)
		return
	}
	fmt.Printf("Активный Xcode / Active Xcode:\n%s\nВерсия / Version: %s\n", displayApp(path), valueOr(ver, "неизвестно / unknown"))
}

func handleSwitch(args []string) {
	fs := flag.NewFlagSet("switch", flag.ContinueOnError)
	dryRun := fs.Bool("dry-run", false, "только показать / dry run only")
	_ = fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Println("⚠️ Укажите версию / Please specify a version, например / for example: xvm switch 16.2")
		return
	}
	targetVer := strings.TrimSpace(fs.Arg(0))

	// 1) Найдём целевой Xcode
	targetApp := filepath.Join(appsDir, fmt.Sprintf("Xcode-%s.app", targetVer))
	if !pathExists(targetApp) {
		fmt.Printf("❌ %s не найден / not found.\n", targetApp)
		fmt.Println("📥 Скачайте нужную версию на / Download from:")
		fmt.Println("   🔗 https://developer.apple.com/download/all/")
		fmt.Printf("   После загрузки распакуйте и поместите в %s как \"Xcode-%s.app\"\n", appsDir, targetVer)
		fmt.Printf("   Затем выполните / Then run: sudo xvm switch %s\n", targetVer)
		return
	}

	// 2) Определяем текущий активный
	activePath, currVer, err := activeXcode()
	if err != nil {
		fmt.Println("⚠️ Не удалось определить активный Xcode / Failed to detect active Xcode:", err)
		activePath = xcodeApp
		currVer, _ = readXcodeVersion(activePath)
	}

	// 3) Готовим план действий
	var planned []string
	if strings.EqualFold(activePath, xcodeApp) && pathExists(xcodeApp) {
		if currVer == "" {
			currVer = "current"
		}
		newName := filepath.Join(appsDir, fmt.Sprintf("Xcode-%s.app", currVer))
		if sameFile(newName, targetApp) {
			// уже совпадает
		} else if pathExists(newName) {
			newName = disambiguate(newName)
		}
		planned = append(planned, fmt.Sprintf("mv %q %q", xcodeApp, newName))
		if !*dryRun {
			if err := os.Rename(xcodeApp, newName); err != nil {
				fmt.Println("Ошибка переименования / rename error:", err)
				return
			}
		}
	}

	// 4) Переименовать целевой
	if !sameFile(targetApp, xcodeApp) {
		planned = append(planned, fmt.Sprintf("mv %q %q", targetApp, xcodeApp))
		if !*dryRun {
			if err := os.Rename(targetApp, xcodeApp); err != nil {
				fmt.Println("Ошибка переименования целевого Xcode / rename target error:", err)
				return
			}
		}
	}

	// 5) Обновить xcode-select
	planned = append(planned, fmt.Sprintf("xcode-select -s %s", filepath.Join(xcodeApp, developerRel)))
	if !*dryRun {
		if err := run("xcode-select", "-s", filepath.Join(xcodeApp, developerRel)); err != nil {
			fmt.Println("Ошибка xcode-select / xcode-select error:", err)
			return
		}
	}

	// 6) Проверим результат
	newVer := ""
	if !*dryRun {
		newVer, _ = activeVersionViaXcodebuild()
	}

	if *dryRun {
		fmt.Println("\n🧪 Сухой прогон / Dry Run:")
		for _, s := range planned {
			fmt.Println("  ", s)
		}
		fmt.Println("✅ Никаких изменений не внесено / No changes applied.")
	} else {
		fmt.Printf("\n✅ Переключение завершено / Switch complete.\nБыло / Was: %s → Стало / Now: %s\n", valueOr(currVer, "неизвестно"), valueOr(newVer, "неизвестно"))
	}
}

// ===== helpers =====

type xcodeEntry struct {
	path    string
	version string
}

func findXcodes() ([]xcodeEntry, error) {
	globs := []string{
		filepath.Join(appsDir, "Xcode.app"),
		filepath.Join(appsDir, "Xcode-*.app"),
	}
	var paths []string
	for _, g := range globs {
		m, _ := filepath.Glob(g)
		paths = append(paths, m...)
	}
	seen := map[string]bool{}
	var out []xcodeEntry
	for _, p := range paths {
		rp, _ := filepath.EvalSymlinks(p)
		if rp == "" {
			rp = p
		}
		if seen[rp] {
			continue
		}
		seen[rp] = true
		ver, _ := readXcodeVersion(rp)
		out = append(out, xcodeEntry{path: rp, version: ver})
	}
	return out, nil
}

func readXcodeVersion(appPath string) (string, error) {
	if !pathExists(appPath) {
		return "", os.ErrNotExist
	}
	info := filepath.Join(appPath, infoPlistRel)
	if pathExists(plistBuddy) && pathExists(info) {
		out, err := exec.Command(plistBuddy, "-c", "Print :CFBundleShortVersionString", info).CombinedOutput()
		if err == nil {
			return strings.TrimSpace(string(out)), nil
		}
	}
	xb := filepath.Join(appPath, developerRel, "usr/bin/xcodebuild")
	if pathExists(xb) {
		out, err := exec.Command(xb, "-version").CombinedOutput()
		if err == nil {
			sc := bufio.NewScanner(strings.NewReader(string(out)))
			if sc.Scan() {
				line := sc.Text()
				re := regexp.MustCompile(`Xcode\s+([0-9.]+)`)
				m := re.FindStringSubmatch(line)
				if len(m) == 2 {
					return m[1], nil
				}
			}
		}
	}
	return "", errors.New("версия не найдена / version not found")
}

func activeXcode() (appPath, version string, err error) {
	devPathOut, err := exec.Command("xcode-select", "-p").CombinedOutput()
	if err != nil {
		return "", "", err
	}
	devPath := strings.TrimSpace(string(devPathOut))
	appPath = filepath.Clean(filepath.Join(devPath, "..", ".."))
	ver, _ := readXcodeVersion(appPath)
	return appPath, ver, nil
}

func activeVersionViaXcodebuild() (string, error) {
	out, err := exec.Command("xcodebuild", "-version").CombinedOutput()
	if err != nil {
		return "", err
	}
	sc := bufio.NewScanner(strings.NewReader(string(out)))
	if sc.Scan() {
		line := sc.Text()
		re := regexp.MustCompile(`Xcode\s+([0-9.]+)`)
		m := re.FindStringSubmatch(line)
		if len(m) == 2 {
			return m[1], nil
		}
	}
	return "", errors.New("не удалось определить / failed to parse xcodebuild output")
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func pathExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func sameFile(a, b string) bool {
	ra, _ := filepath.EvalSymlinks(a)
	rb, _ := filepath.EvalSymlinks(b)
	if ra == "" {
		ra = a
	}
	if rb == "" {
		rb = b
	}
	return filepath.Clean(ra) == filepath.Clean(rb)
}

func displayApp(p string) string {
	if p == "" {
		return "(не найден / not found)"
	}
	return p
}

func valueOr(s, alt string) string {
	if strings.TrimSpace(s) == "" {
		return alt
	}
	return s
}

func disambiguate(path string) string {
	base := path
	i := 1
	for pathExists(path) {
		path = fmt.Sprintf("%s-%d", base, i)
		i++
	}
	return path
}
