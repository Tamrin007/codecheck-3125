# CLIアプリケーション作成用テンプレート(Go)

GoでCLIアプリケーションを作成するためのテンプレートです。

[app.go](app.go)を編集することでCLIアプリケーションを作成することができます。

このテンプレートではCLI作成のユーティリティとして、cli.goを使用していsます。  
cli.goについての詳細は[cli.goのリポジトリ](https://github.com/codegangsta/cli)をごらんください。

## コマンドライン引数の取得方法
app.goファイルには`doMain`というメソッドが定義されています。

``` go
func doMain(c *cli.Context) {
}
```

通常のパラメータは`c.Args()`で配列として取得できます。
オプション形式の引数を使用する場合は[main.go](main.rb)でcliをカスタマイズしてください。

## コマンド実行結果の標準出力への出力
`fmt.Println`メソッドを使用してください。

``` go
  fmt.Println(c.Argus[0])
```
