# XVM - Xcode Version Manager

A simple and powerful command-line tool for managing multiple Xcode versions on macOS.

Простое и мощное командное приложение для управления несколькими версиями Xcode на macOS.

## Installation / Установка

### Quick Install / Быстрая установка

```bash
curl -sSL https://raw.githubusercontent.com/AntonSeagull/xcode-version-manager/main/install-xvm.sh | bash
```

### Manual Install / Ручная установка

1. Download the latest binary from [releases](https://github.com/AntonSeagull/xcode-version-manager/releases)
2. Make it executable: `chmod +x xvm-darwin-amd64` (or `xvm-darwin-arm64` for Apple Silicon)
3. Move to PATH: `sudo mv xvm-darwin-amd64 /usr/local/bin/xvm`

4. Скачайте последний бинарник из [релизов](https://github.com/AntonSeagull/xcode-version-manager/releases)
5. Сделайте исполняемым: `chmod +x xvm-darwin-amd64` (или `xvm-darwin-arm64` для Apple Silicon)
6. Переместите в PATH: `sudo mv xvm-darwin-amd64 /usr/local/bin/xvm`

## Features / Возможности

- **List all Xcode versions** / Показывает все установленные версии Xcode
- **Show current active version** / Показывает текущую активную версию
- **Switch between versions** / Переключение между версиями
- **Dry-run mode** / Режим предварительного просмотра
- **Bilingual interface** / Двуязычный интерфейс (English/Русский)

## Usage / Использование

### List all Xcode installations / Показать все установки Xcode

```bash
xvm list
```

Example output / Пример вывода:

```
Активная / Active: /Applications/Xcode.app (16.2)

Найдено / Found in /Applications:
 * /Applications/Xcode.app (16.2)
   /Applications/Xcode-15.4.app (15.4)
   /Applications/Xcode-14.3.1.app (14.3.1)
```

### Show current active version / Показать текущую активную версию

```bash
xvm current
```

Example output / Пример вывода:

```
Активный Xcode / Active Xcode:
/Applications/Xcode.app
Версия / Version: 16.2
```

### Switch to a specific version / Переключиться на конкретную версию

```bash
sudo xvm switch 15.4
```

### Dry-run mode / Режим предварительного просмотра

```bash
xvm switch 15.4 --dry-run
```

This will show what actions would be performed without making any changes.

Это покажет, какие действия будут выполнены, без внесения изменений.

## How it works / Как это работает

XVM manages Xcode versions by:

1. **Renaming the current active Xcode** (usually `/Applications/Xcode.app`) to a versioned name (e.g., `Xcode-16.2.app`)
2. **Renaming the target Xcode** (e.g., `Xcode-15.4.app`) to `/Applications/Xcode.app`
3. **Updating xcode-select** to point to the new active version

XVM управляет версиями Xcode следующим образом:

1. **Переименовывает текущий активный Xcode** (обычно `/Applications/Xcode.app`) в версионное имя (например, `Xcode-16.2.app`)
2. **Переименовывает целевой Xcode** (например, `Xcode-15.4.app`) в `/Applications/Xcode.app`
3. **Обновляет xcode-select** для указания на новую активную версию

## Requirements / Требования

- **macOS** (tested on macOS 10.15+)
- **Multiple Xcode installations** in `/Applications/` directory
- **Administrator privileges** for switching versions (sudo required)

- **macOS** (протестировано на macOS 10.15+)
- **Несколько установок Xcode** в директории `/Applications/`
- **Права администратора** для переключения версий (требуется sudo)

## Setting up multiple Xcode versions / Настройка нескольких версий Xcode

1. **Download Xcode versions** from [Apple Developer Downloads](https://developer.apple.com/download/all/)
2. **Install each version** to `/Applications/` with a versioned name:
   - `Xcode-16.2.app`
   - `Xcode-15.4.app`
   - `Xcode-14.3.1.app`
3. **Use XVM** to switch between them

4. **Скачайте версии Xcode** с [Apple Developer Downloads](https://developer.apple.com/download/all/)
5. **Установите каждую версию** в `/Applications/` с версионным именем:
   - `Xcode-16.2.app`
   - `Xcode-15.4.app`
   - `Xcode-14.3.1.app`
6. **Используйте XVM** для переключения между ними

## Commands Reference / Справочник команд

| Command / Команда                | Description / Описание                                                       |
| -------------------------------- | ---------------------------------------------------------------------------- |
| `xvm list`                       | Show all installed Xcode versions / Показать все установленные версии Xcode  |
| `xvm current`                    | Show currently active Xcode version / Показать текущую активную версию Xcode |
| `xvm switch <version>`           | Switch to specified version / Переключиться на указанную версию              |
| `xvm switch <version> --dry-run` | Preview switch operation / Предварительный просмотр операции переключения    |

## Examples / Примеры

```bash
# List all versions / Показать все версии
xvm list

# Show current version / Показать текущую версию
xvm current

# Switch to Xcode 15.4 / Переключиться на Xcode 15.4
sudo xvm switch 15.4

# Preview switch to Xcode 14.3.1 / Предварительный просмотр переключения на Xcode 14.3.1
xvm switch 14.3.1 --dry-run
```

## Troubleshooting / Решение проблем

### Permission denied / Ошибка прав доступа

Make sure you're using `sudo` when switching versions:

Убедитесь, что используете `sudo` при переключении версий:

```bash
sudo xvm switch 15.4
```

### Xcode not found / Xcode не найден

Ensure your Xcode installations are named correctly:

Убедитесь, что ваши установки Xcode названы правильно:

```bash
ls -la /Applications/Xcode*
```

Should show files like `/Applications/Xcode-16.2.app`, `/Applications/Xcode-15.4.app`, etc.

Должны отображаться файлы типа `/Applications/Xcode-16.2.app`, `/Applications/Xcode-15.4.app` и т.д.

### Version detection issues / Проблемы с определением версии

XVM tries multiple methods to detect Xcode versions:

- Reading Info.plist
- Using xcodebuild -version

XVM пытается несколько методов для определения версий Xcode:

- Чтение Info.plist
- Использование xcodebuild -version

## Contributing / Участие в разработке

Contributions are welcome! Please feel free to submit a Pull Request.

Участие в разработке приветствуется! Пожалуйста, отправляйте Pull Request.

## License / Лицензия

This project is open source and available under the [MIT License](LICENSE).

Этот проект с открытым исходным кодом и доступен под [лицензией MIT](LICENSE).

## Author / Автор

Created by [AntonSeagull](https://github.com/AntonSeagull)

Создано [AntonSeagull](https://github.com/AntonSeagull)
