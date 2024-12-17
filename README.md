# pwzip-creator

## build
> docker build -t pwzip-creator:latest .

## 実行
> docker run -p 8080:8080 pwzip-creator:latest

## 呼び方
> curl.exe -X POST -F "files=@C:\Users\Hitoshi Ichikawa\image.png" "http://localhost:8080/createzip?password=password&zip_filename=test.zip" --output test2.zip

querystring
password: パスワード
zip_filename: ダウンロードされるファイル名

postbody
form-data
key:files
value:ファイル