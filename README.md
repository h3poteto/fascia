[![Circle CI](https://circleci.com/gh/h3poteto/fascia.svg?style=shield&circle-token=	bcf712be69dd52a490bccb78f71b9637657d45c5)](https://circleci.com/gh/fascia/fascia)

# fascia

## Installation
### Environmets

環境変数を準備する必要がります．
各自で，`.bash_profile`等に記述してください．おすすめは`direnv`です．

```
export DB_USER="root"
export DB_PASSWORD="hogehoge"
export DB_NAME="fascia"
export DB_TEST_NAME="fascia_test"
export GOJIENV="development"
export GOJIROOT="/home/ubuntu/fascia"
export CLIENT_ID="hogehoge"
export CLIENT_SECRET="fugafuga"
export TEST_TOKEN="testhoge"
export SLACK_URL="https://hooks.slack.com/services/hogehoge/fugafuga"
```
`CLIENT_ID`, `CLIENT_SECRET`, `TEST_TOKEN` は適当にgithubでアプリケーションを作成して自分で用意してください．
DB関連の設定については，`db/dbconf.yml` を参考に必要項目を用意してください．

### go
goは1.5を前提としています．
パッケージ管理として`gom`を使います．
`gom`のインストールいついては下記を参照．
https://github.com/mattn/gom


```
$ gom install
$ gom exec goose up
$ gom run server.go
```
正常に起動することを確認してください．


### npm
フロント側は`React`を使う関係上，パッケージ管理は`npm`で行います．
`nodejs`と`npm`を使えるようにしておいてください．
```
$ npm install
```

以上で準備は完了です．

## Development
### go
```
$ gom run server.go
```
これで，ブラウザから`localhost:9090`で確認できます．

`.go`のソースはコンパイルが必要になるため，サーバーの再起動無しに更新が反映されることはありえません．ソースを変更した場合は，その都度サーバを再起動してください．

ただし，テンプレートは変更分を読み直ししてくれるため，`tpl`のソースについては，再起動無しで反映されます．

### js, scss
jsやcssを変更する場合は下記のコマンドによってassetsの差分コンパイルが走るようにしておいてください．
```
$ npm run-script watch
```

## Test
テストフレームワークには[Ginkgo](https://github.com/onsi/ginkgo)を採用しています．

また，マッチャーは[Gomega](https://github.com/onsi/gomega)を使用します．


以下のコマンドにより，すべてのテストを実行してくれます．
```
$ export GOJIENV=test; gom exec ginkgo -r ./
```


## Log
ログはすべてlogrusを使っています． `modules/logging` パッケージがlogrusの設定を行っているため，このパッケージを使うことでログを吐き出せます．

ログレベルの設定基準は以下のとおりです．

- Debug

  開発時のみ見えていれば良い情報で，多分に機密な情報も含む．
- Info
 
  最悪流失しても問題ないレベルの情報のみを載せておく．行動ログ的なものだと思えば良い．
- Warn

  将来的に治したい部分，本来の挙動と違うがエラーにするほではないものを出力する．
- Error

  ユーザ側にエラーを表示する際に一緒に出しておく．これは状況によっては想定できる範囲のエラーまでを扱い，それ以上のものはPanicを用いる．
- Fatal

  プログラムを終了するようなエラーを起こす予定はないため使わない．
- Panic

  プログラムの流れ的にうまく行くことを確信して良い部分のエラーとして使用する．
