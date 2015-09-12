# fascia

## Setup
### Environmets

環境変数の設定を行う必要があります．
各自で，`.bash_profile`等に記述してください．おすすめは`direnv`です．

```yml
export DB_USER="root"
export DB_PASSWORD="hogehoge"
export DB_NAME="fascia"
export DB_TEST_NAME="fascia_test"
```

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

`$ npm install`

また，ES6やscssのコンパイルのために`gulp`を使います．そのため，`gulp`コマンドを使えるようにしておく必要があります．
`$ npm install -g gulp`

`$ gulp watchify`
で，監視＆差分コンパイルが走ることを確認してください．


以上で準備は完了です．

## Development
### go
`$ gom run server.go`
これで，ブラウザから`localhost:9090`で確認できます．

`.go`のソースはコンパイルが必要になるため，サーバーの再起動無しに更新が反映されることはありえません．ソースを変更した場合は，その都度サーバを再起動してください．

ただし，テンプレートは変更分を読み直ししてくれるため，`tpl`のソースについては，再起動無しで反映されます．

### js, scss
`$ gulp watchify`
これで，jsファイルの監視と差分コンパイルが走ります．
scssについては，まだコンパイルを導入していません．
