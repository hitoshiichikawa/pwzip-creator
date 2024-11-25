package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"

    mullinsZip "github.com/alexmullins/zip" // 外部ライブラリをインポート
)

func main() {
    http.HandleFunc("/createzip", CreateZip)
    log.Println("サーバーをポート8080で起動します...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("サーバーの起動に失敗しました: ", err)
    }
}

// HTTPハンドラ関数
func CreateZip(w http.ResponseWriter, r *http.Request) {
    password := r.URL.Query().Get("password")
    zipFilename := r.URL.Query().Get("zip_filename")
    if zipFilename == "" {
        zipFilename = "protected.zip"
    }
    if password == "" {
        http.Error(w, "Password is required", http.StatusBadRequest)
        return
    }

    tempDir, err := os.MkdirTemp("", "upload")
    if err != nil {
        http.Error(w, "一時ディレクトリの作成に失敗しました", http.StatusInternalServerError)
        return
    }
    defer os.RemoveAll(tempDir)

    if err := r.ParseMultipartForm(10 << 20); err != nil {
        http.Error(w, "フォームの解析に失敗しました", http.StatusBadRequest)
        return
    }

    files := r.MultipartForm.File["files"]
    if len(files) == 0 {
        http.Error(w, "ファイルがアップロードされていません", http.StatusBadRequest)
        return
    }

    // ZIPファイルの作成
    zipPath := filepath.Join(tempDir, zipFilename)
    zipFile, err := os.Create(zipPath)
    if err != nil {
        http.Error(w, "ZIPファイルの作成に失敗しました", http.StatusInternalServerError)
        return
    }
    defer zipFile.Close()

    zipWriter := mullinsZip.NewWriter(zipFile) // パスワード対応ライブラリを使用
    defer func() {
        if err := zipWriter.Close(); err != nil {
            // エラーを標準エラー出力に記録
            fmt.Fprintf(os.Stderr, "ZIPライターのクローズに失敗しました: %v\n", err)
        }
    }()

    for _, fileHeader := range files {
        file, err := fileHeader.Open()
        if err != nil {
            http.Error(w, "ファイルのオープンに失敗しました", http.StatusInternalServerError)
            return
        }
        defer file.Close()

        // ZIP内にファイルを書き込み
        zipEntry, err := zipWriter.Encrypt(fileHeader.Filename, password)
        if err != nil {
            http.Error(w, "ZIPエントリの作成に失敗しました", http.StatusInternalServerError)
            return
        }

        if _, err := io.Copy(zipEntry, file); err != nil {
            http.Error(w, "ZIPへのファイル書き込みに失敗しました", http.StatusInternalServerError)
            return
        }
    }

    if err := zipWriter.Close(); err != nil {
        // エラーを標準エラー出力に記録
        fmt.Fprintf(os.Stderr, "ZIPファイルの最終化に失敗しました: %v\n", err)
        return
    }

    zipData, err := os.ReadFile(zipPath)
    if err != nil {
        http.Error(w, "ZIPファイルの読み込みに失敗しました", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/zip")
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipFilename))
    w.Write(zipData)
}