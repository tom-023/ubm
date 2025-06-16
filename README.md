# ubm - URL Bookmark Manager

[English](./README_en.md) | 日本語

ubm (URL Bookmark Manager) は、ツリー構造でブックマークを整理できる対話的なコマンドラインツールです。

## 機能

- 📁 **階層的なカテゴリ管理**: ブックマークをツリー構造で整理
- 🔍 **対話的なナビゲーション**: 矢印キーで簡単にブックマークを探索
- 🌐 **ブラウザ統合**: 選択したブックマークを自動的にブラウザで開く
- ✏️ **編集機能**: ブックマークの情報を後から変更可能
- 📂 **カテゴリ間の移動**: ブックマークを別のカテゴリに移動

## インストール

### Goを使用したインストール

```bash
go install github.com/tom-023/ubm/cmd/ubm@latest
```

### ソースからビルド

```bash
git clone https://github.com/tom-023/ubm.git
cd ubm
go build -o ubm ./cmd/ubm
```

## 使い方

### ブックマークの追加

```bash
# 対話的に追加（URL、タイトル、カテゴリを順番に入力）
ubm add
```

### ブックマークの閲覧

```bash
# 対話的なナビゲーション（ブックマークを選択するとブラウザで開いて終了）
ubm list

# ツリー形式で全体を表示
ubm show

# フラットリストで表示
ubm show --flat
```

### カテゴリ管理

```bash
# カテゴリの作成
ubm category create

# カテゴリ一覧
ubm category list

# 空のカテゴリを削除
ubm category delete
```

### ブックマークの編集

```bash
# タイトルを指定して編集
ubm edit "ブックマークのタイトル"

# ブックマークを別のカテゴリに移動
ubm move

# ブックマークの削除
ubm delete
```

## キーボードショートカット

対話的なモードでは以下のキーが使用できます：

- `↑` `↓`: 項目の選択
- `Enter`: 選択/確定
- `Backspace` `Esc`: 親ディレクトリに戻る
- `/`: 検索モードの切り替え
- `q` `Ctrl+C`: 終了

## データの保存場所

ブックマークデータは以下の場所に保存されます：

- **Linux/macOS**: `~/.config/ubm/bookmarks.json`
- **Windows**: `%APPDATA%\ubm\bookmarks.json`

## 開発

### 必要な環境

- Go 1.23以上

### 依存関係

- [cobra](https://github.com/spf13/cobra) - CLIフレームワーク
- [promptui](https://github.com/manifoldco/promptui) - 対話的プロンプト
- [browser](https://github.com/pkg/browser) - ブラウザ統合

### ビルド

```bash
go build -o ubm ./cmd/ubm
```

### テスト

```bash
go test ./...
```

## ライセンス

MIT License
