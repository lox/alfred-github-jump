# Alfred Github Jump

A [workflow for Alfred 3](https://www.alfredapp.com/help/workflows/) for indexing your github repositories, allowing you to quickly filter them and open them in your default browser

![](http://lachlan.me/s/wvjHZ.png)


## Development

```
# Make sure the workflow dir exists
mkdir -p "$HOME/Library/Application Support/Alfred 3/Alfred.alfredpreferences/workflows"

# Then checkout the project
go get -u github.com/lox/alfred-github-jump
cd "$GOPATH/src/github.com/lox/alfred-github-jump"

# Build it and link it into Alfred
make build
ln -s "$PWD" "$HOME/Library/Application Support/Alfred 3/Alfred.alfredpreferences/workflows/alfred-github-jump"
